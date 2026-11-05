package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const maxBody = 64 << 20

type restResp struct {
	code   int
	header http.Header
	body   []byte
	dur    time.Duration
}

// rest performs an HTTP request against the server. path is relative to the base URL.
func (e *env) rest(ctx context.Context, method, path string, body any) (*restResp, error) {
	var rd io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rd = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, e.cfg.restBase+path, rd)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-device-id", e.cfg.deviceID)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	start := time.Now()
	res, err := e.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	b, err := io.ReadAll(io.LimitReader(res.Body, maxBody))
	if err != nil {
		return nil, err
	}
	return &restResp{code: res.StatusCode, header: res.Header, body: b, dur: time.Since(start)}, nil
}

// restOK performs a request and records a failure unless it returns 200.
func (e *env) restOK(ctx context.Context, t *T, method, path string, body any) (*restResp, bool) {
	var r *restResp
	err := t.timed("REST "+method+" "+path, func() (err error) {
		r, err = e.rest(ctx, method, path, body)
		return err
	})
	if err != nil {
		t.Failf("REST %s %s: %v", method, path, err)
		return nil, false
	}
	if r.code != http.StatusOK {
		t.Failf("REST %s %s: HTTP %d: %s", method, path, r.code, snippet(r.body))
		return r, false
	}
	return r, true
}

// jsonLen returns the number of elements of a JSON array, or of the array
// stored under key if the body is an object.
func jsonLen(body []byte, key string) (int, error) {
	if key == "" {
		var arr []json.RawMessage
		err := json.Unmarshal(body, &arr)
		return len(arr), err
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(body, &obj); err != nil {
		return 0, err
	}
	raw, ok := obj[key]
	if !ok {
		return 0, nil // proto3 JSON omits empty repeated fields
	}
	var arr []json.RawMessage
	err := json.Unmarshal(raw, &arr)
	return len(arr), err
}

func snippet(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 160 {
		s = s[:160] + "…"
	}
	return s
}

// expectCode checks that err carries one of the wanted gRPC codes.
// An Internal/Unknown code for a client mistake is reported as a warning
// (REST clients see a 500, and it ends up in Sentry); anything else unexpected fails.
func expectCode(t *T, what string, err error, want ...codes.Code) {
	got := status.Code(err)
	if slices.Contains(want, got) {
		t.Infof("%s: %s as expected", what, got)
		return
	}
	msg := fmt.Sprintf("%s: got %s %q, want one of %v", what, got, status.Convert(err).Message(), want)
	switch got {
	case codes.Internal, codes.Unknown:
		t.Warnf("%s (client error reported as server error)", msg)
	default:
		t.Failf("%s", msg)
	}
}

// checkURLs verifies that a sample of the given URLs is reachable. Broken URLs
// are reported at the given level.
func (e *env) checkURLs(ctx context.Context, t *T, what string, urls []string, sample int, level Status) {
	urls = sampleOf(dedupe(urls), sample)
	if len(urls) == 0 {
		return
	}
	var mu sync.Mutex
	var broken []string
	var wg sync.WaitGroup
	sem := make(chan struct{}, 6)
	for _, u := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			if code, err := e.probeURL(ctx, u); err != nil || code >= 400 {
				reason := fmt.Sprintf("HTTP %d", code)
				if err != nil {
					reason = err.Error()
				}
				mu.Lock()
				broken = append(broken, fmt.Sprintf("%s (%s)", u, reason))
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if len(broken) == 0 {
		t.Infof("%s: %d/%d sampled URLs reachable", what, len(urls), len(urls))
		return
	}
	slices.Sort(broken)
	t.add(level, "%s: %d/%d sampled URLs broken, e.g. %s", what, len(broken), len(urls), strings.Join(broken[:min(3, len(broken))], ", "))
}

func (e *env) probeURL(ctx context.Context, u string) (int, error) {
	do := func(method string) (int, error) {
		req, err := http.NewRequestWithContext(ctx, method, u, nil)
		if err != nil {
			return 0, err
		}
		req.Header.Set("User-Agent", "campus-backend-livetest")
		if method == http.MethodGet {
			req.Header.Set("Range", "bytes=0-0")
		}
		res, err := e.http.Do(req)
		if err != nil {
			return 0, err
		}
		_, _ = io.Copy(io.Discard, io.LimitReader(res.Body, 1<<16))
		_ = res.Body.Close()
		return res.StatusCode, nil
	}
	code, err := do(http.MethodHead)
	if err != nil || code == http.StatusMethodNotAllowed || code == http.StatusForbidden || code == http.StatusNotImplemented {
		return do(http.MethodGet)
	}
	return code, nil
}

func dedupe(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if s != "" && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// sampleOf returns up to n evenly spaced elements.
func sampleOf[S any](in []S, n int) []S {
	if n <= 0 || len(in) <= n {
		return in
	}
	out := make([]S, 0, n)
	for i := range n {
		out = append(out, in[i*len(in)/n])
	}
	return out
}

func pct(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return 100 * float64(part) / float64(total)
}

func sizeOf(m proto.Message) int { return proto.Size(m) }
