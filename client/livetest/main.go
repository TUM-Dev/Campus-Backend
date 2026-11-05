// Command livetest probes a running Campus-Backend (by default the production
// deployment at api.tum.app) over gRPC and REST and reports anything that looks
// broken, inconsistent or risky.
//
// The harness is read-only: write RPCs are only ever called with requests that
// the server rejects during validation, before anything is persisted.
//
//	go run ./livetest                 # against production
//	go run ./livetest -local          # against a local server on :50051
//	go run ./livetest -run 'news/'    # only checks whose "group/name" matches
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type config struct {
	grpcAddr  string
	restBase  string
	plaintext bool
	timeout   time.Duration
	run       *regexp.Regexp
	jsonOut   bool
	strict    bool
	verbose   bool
	burst     int
	deviceID  string
	slow      time.Duration
}

func main() {
	var cfg config
	var runPattern string
	var local bool
	flag.StringVar(&cfg.grpcAddr, "grpc", "api.tum.app:443", "gRPC address (host:port)")
	flag.StringVar(&cfg.restBase, "rest", "https://api.tum.app", "base URL of the HTTP server (REST lives under /v1)")
	flag.BoolVar(&cfg.plaintext, "plaintext", false, "connect to gRPC without TLS")
	flag.BoolVar(&local, "local", false, "shorthand for -grpc 127.0.0.1:50051 -rest http://127.0.0.1:50051 -plaintext")
	flag.DurationVar(&cfg.timeout, "timeout", 30*time.Second, "timeout per check")
	flag.DurationVar(&cfg.slow, "slow", 2*time.Second, "warn about single requests slower than this")
	flag.StringVar(&runPattern, "run", "", "only run checks whose \"group/name\" matches this regexp")
	flag.BoolVar(&cfg.jsonOut, "json", false, "print results as JSON")
	flag.BoolVar(&cfg.strict, "strict", false, "exit non-zero on warnings as well as failures")
	flag.BoolVar(&cfg.verbose, "v", false, "also print the notes of passing checks")
	flag.IntVar(&cfg.burst, "burst", 10, "concurrent requests in the burst check (0 disables it)")
	flag.StringVar(&cfg.deviceID, "device-id", "campus-backend-livetest", "x-device-id sent with every request")
	flag.Parse()

	if local {
		cfg.grpcAddr, cfg.restBase, cfg.plaintext = "127.0.0.1:50051", "http://127.0.0.1:50051", true
	}
	cfg.restBase = strings.TrimSuffix(cfg.restBase, "/")
	if runPattern != "" {
		re, err := regexp.Compile(runPattern)
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid -run pattern:", err)
			os.Exit(2)
		}
		cfg.run = re
	}

	e, err := newEnv(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "setup failed:", err)
		os.Exit(2)
	}
	defer func() { _ = e.conn.Close() }()

	results := runChecks(e, allChecks())
	if cfg.jsonOut {
		printJSON(results)
	} else {
		printText(results, cfg.verbose)
	}

	worst := Pass
	for _, r := range results {
		worst = max(worst, r.Status)
	}
	if worst == Fail || (cfg.strict && worst == Warn) {
		os.Exit(1)
	}
}

func allChecks() []check {
	var all []check
	all = append(all, infraChecks()...)
	all = append(all, sweepChecks()...)
	all = append(all, newsChecks()...)
	all = append(all, cafeteriaChecks()...)
	all = append(all, miscChecks()...)
	all = append(all, restChecks()...)
	all = append(all, perfChecks()...)
	return all
}

// ---------------------------------------------------------------------------
// environment shared by all checks

type env struct {
	cfg   config
	conn  *grpc.ClientConn
	c     pb.CampusClient
	http  *http.Client
	cache sync.Map // string -> *lazy
}

type lazy struct {
	once sync.Once
	v    any
	err  error
}

func newEnv(cfg config) (*env, error) {
	var creds credentials.TransportCredentials
	if cfg.plaintext {
		creds = insecure.NewCredentials()
	} else {
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("loading system cert pool: %w", err)
		}
		creds = credentials.NewTLS(&tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12})
	}
	conn, err := grpc.NewClient(cfg.grpcAddr,
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(64<<20)))
	if err != nil {
		return nil, fmt.Errorf("creating gRPC client: %w", err)
	}
	return &env{
		cfg:  cfg,
		conn: conn,
		c:    pb.NewCampusClient(conn),
		http: &http.Client{
			Timeout: cfg.timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}, nil
}

// md attaches the metadata the apps send with every request.
func (e *env) md(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		"x-device-id", e.cfg.deviceID,
		"x-os-version", "livetest",
		"x-app-version", "livetest")
}

// memo computes a value once per run and shares it between checks. It runs f
// with its own context so that one check's cancellation can't poison the cache.
func memo[V any](e *env, key string, f func(ctx context.Context) (V, error)) (V, error) {
	l, _ := e.cache.LoadOrStore(key, &lazy{})
	lz := l.(*lazy)
	lz.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), e.cfg.timeout)
		defer cancel()
		lz.v, lz.err = f(e.md(ctx))
	})
	if lz.err != nil {
		var zero V
		return zero, lz.err
	}
	return lz.v.(V), nil
}

// ---------------------------------------------------------------------------
// checks and results

type Status int

const (
	Pass Status = iota
	Skip
	Warn
	Fail
)

func (s Status) String() string {
	return [...]string{"PASS", "SKIP", "WARN", "FAIL"}[s]
}

func (s Status) MarshalText() ([]byte, error) { return []byte(s.String()), nil }

type note struct {
	Status Status `json:"status"`
	Msg    string `json:"msg"`
}

// T collects the outcome of one check. It is safe for concurrent use.
type T struct {
	mu     sync.Mutex
	status Status
	notes  []note
	slow   time.Duration
}

func (t *T) add(s Status, format string, args ...any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.status = max(t.status, s)
	t.notes = append(t.notes, note{s, fmt.Sprintf(format, args...)})
}

func (t *T) Infof(format string, args ...any) { t.add(Pass, format, args...) }
func (t *T) Skipf(format string, args ...any) { t.add(Skip, format, args...) }
func (t *T) Warnf(format string, args ...any) { t.add(Warn, format, args...) }
func (t *T) Failf(format string, args ...any) { t.add(Fail, format, args...) }

// timed runs f and complains if it was slow.
func (t *T) timed(what string, f func() error) error {
	start := time.Now()
	err := f()
	if d := time.Since(start); d > t.slow {
		t.Warnf("%s took %s (threshold %s)", what, d.Round(time.Millisecond), t.slow)
	}
	return err
}

type check struct {
	group, name string
	// serial checks run alone after all others, so they don't skew each other's timing.
	serial bool
	fn     func(ctx context.Context, e *env, t *T)
}

func (c check) id() string { return c.group + "/" + c.name }

type result struct {
	Check    string        `json:"check"`
	Status   Status        `json:"status"`
	Duration time.Duration `json:"duration_ns"`
	Notes    []note        `json:"notes,omitempty"`
}

func runChecks(e *env, checks []check) []result {
	var parallel, serial []check
	for _, c := range checks {
		if e.cfg.run != nil && !e.cfg.run.MatchString(c.id()) {
			continue
		}
		if c.serial {
			serial = append(serial, c)
		} else {
			parallel = append(parallel, c)
		}
	}

	results := make([]result, 0, len(parallel)+len(serial))
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 4) // be gentle with production
	for _, c := range parallel {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			r := runOne(e, c)
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}()
	}
	wg.Wait()
	for _, c := range serial {
		results = append(results, runOne(e, c))
	}

	order := map[string]int{}
	for i, c := range checks {
		order[c.id()] = i
	}
	sort.Slice(results, func(i, j int) bool { return order[results[i].Check] < order[results[j].Check] })
	return results
}

func runOne(e *env, c check) (r result) {
	t := &T{slow: e.cfg.slow}
	ctx, cancel := context.WithTimeout(e.md(context.Background()), e.cfg.timeout)
	defer cancel()
	start := time.Now()
	defer func() {
		if p := recover(); p != nil {
			t.Failf("check panicked: %v", p)
		}
		if ctx.Err() != nil {
			t.Failf("check did not finish within %s", e.cfg.timeout)
		}
		r = result{Check: c.id(), Status: t.status, Duration: time.Since(start), Notes: t.notes}
	}()
	c.fn(ctx, e, t)
	return
}

func printText(results []result, verbose bool) {
	counts := map[Status]int{}
	for _, r := range results {
		counts[r.Status]++
		fmt.Printf("%s  %-40s %6s\n", r.Status, r.Check, r.Duration.Round(time.Millisecond))
		for _, n := range r.Notes {
			if verbose || n.Status >= Skip {
				fmt.Printf("        %s %s\n", marker(n.Status), n.Msg)
			}
		}
	}
	fmt.Printf("\n%d checks: %d passed, %d warnings, %d failed, %d skipped\n",
		len(results), counts[Pass], counts[Warn], counts[Fail], counts[Skip])
}

func marker(s Status) string {
	switch s {
	case Fail:
		return "✗"
	case Warn:
		return "!"
	case Skip:
		return "-"
	default:
		return "·"
	}
}

func printJSON(results []result) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(results)
}
