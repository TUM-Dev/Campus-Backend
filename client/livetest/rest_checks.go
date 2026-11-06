package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sort"
	"sync"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
)

func restChecks() []check {
	return []check{
		{group: "rest", name: "parity", fn: checkRESTParity},
		{group: "rest", name: "error-mapping", fn: checkRESTErrors},
		{group: "rest", name: "cors", fn: checkCORS},
	}
}

func perfChecks() []check {
	return []check{
		{group: "perf", name: "burst", serial: true, fn: checkBurst},
	}
}

// checkRESTParity calls the gateway and compares item counts with the gRPC answers.
func checkRESTParity(ctx context.Context, e *env, t *T) {
	type pair struct {
		path string
		key  string // JSON key of the list, "" if the body is the list itself (response_body)
		grpc func() (int, error)
	}
	pairs := []pair{
		{"/v1/news/sources", "", func() (int, error) { s, err := e.newsSources(); return len(s), err }},
		{"/v1/news", "", func() (int, error) { n, err := e.allNews(); return len(n), err }},
		{"/v1/movies/-1", "movies", func() (int, error) { m, err := e.allMovies(); return len(m), err }},
		{"/v1/student_clubs", "collections", func() (int, error) { c, err := e.studentClubs(pb.Language_German); return len(c), err }},
		{"/v1/student_clubs?language=English", "collections", func() (int, error) { c, err := e.studentClubs(pb.Language_English); return len(c), err }},
		{"/v1/dish/rating/allRatingTags", "", func() (int, error) {
			r, err := e.c.ListAvailableDishTags(ctx, &pb.ListAvailableDishTagsRequest{})
			return len(r.GetRatingTags()), err
		}},
		{"/v1/dish/rating/allDishTags", "", func() (int, error) {
			r, err := e.c.ListNameTags(ctx, &pb.ListNameTagsRequest{})
			return len(r.GetRatingTags()), err
		}},
		{"/v1/canteen/rating/allRatingTags", "", func() (int, error) {
			r, err := e.c.ListAvailableCanteenTags(ctx, &pb.ListAvailableCanteenTagsRequest{})
			return len(r.GetRatingTags()), err
		}},
		{"/v1/canteen/allCanteens", "", func() (int, error) {
			r, err := e.c.ListCanteens(ctx, &pb.ListCanteensRequest{})
			return len(r.GetCanteen()), err
		}},
	}
	for _, p := range pairs {
		want, gerr := p.grpc()
		r, ok := e.restOK(ctx, t, http.MethodGet, p.path, nil)
		if !ok {
			continue
		}
		got, err := jsonLen(r.body, p.key)
		if err != nil {
			t.Failf("GET %s: unexpected JSON shape: %v: %s", p.path, err, snippet(r.body))
			continue
		}
		if gerr != nil {
			t.Warnf("GET %s works (%d items) but the gRPC call fails: %v", p.path, got, gerr)
			continue
		}
		if got != want {
			t.Warnf("GET %s returned %d items, gRPC %d (caching, or a mapping bug)", p.path, got, want)
		} else {
			t.Infof("GET %s: %d items, matches gRPC", p.path, got)
		}
	}
}

// checkRESTErrors makes sure client errors don't surface as 5xx and errors have a JSON body.
func checkRESTErrors(ctx context.Context, e *env, t *T) {
	year, week := berlinNow().ISOWeek()
	cases := []struct {
		method, path string
		body         any
		want         []int
	}{
		{http.MethodGet, "/v1/updatenote/999999999", nil, []int{404}},
		{http.MethodGet, "/v1/updatenote/not-a-number", nil, []int{400}},
		{http.MethodGet, "/v1/movies/not-a-number", nil, []int{400}},
		{http.MethodGet, fmt.Sprintf("/v1/dishes?canteenId=mensa-garching&year=%d&week=99&day=0", year), nil, []int{400}},
		{http.MethodGet, fmt.Sprintf("/v1/dishes?canteenId=mensa-garching&year=%d&week=%d&day=0", year, week), nil, []int{200}},
		{http.MethodGet, "/v1/student_clubs?language=Klingon", nil, []int{400}},
		{http.MethodPost, "/v1/canteen/rating/get", map[string]any{"canteenId": unknownCanteen, "limit": 1}, []int{400, 404}},
		{http.MethodPost, "/v1/canteen/rating/new", map[string]any{"canteenId": unknownCanteen, "points": 3}, []int{400, 404}},
		{http.MethodPost, "/v1/canteen/rating/new", map[string]any{"canteenId": unknownCanteen, "points": 99}, []int{400}},
		{http.MethodPost, "/v1/dish/rating/new", map[string]any{"canteenId": unknownCanteen, "dish": "x", "points": 3}, []int{400, 404}},
		{http.MethodPost, "/v1/device", map[string]any{}, []int{400}},
		{http.MethodPost, "/v1/news/sources", nil, []int{404, 405, 501}},
	}
	for _, c := range cases {
		r, err := e.rest(ctx, c.method, c.path, c.body)
		what := c.method + " " + c.path
		if err != nil {
			t.Failf("%s: %v", what, err)
			continue
		}
		if slices.Contains(c.want, r.code) {
			t.Infof("%s -> %d", what, r.code)
		} else if r.code >= 500 && !slices.Contains(c.want, http.StatusOK) {
			t.Warnf("%s -> %d, want %v: %s", what, r.code, c.want, snippet(r.body))
		} else {
			t.Failf("%s -> %d, want %v: %s", what, r.code, c.want, snippet(r.body))
		}
		if r.code >= 400 && !json.Valid(r.body) {
			t.Warnf("%s: error body is not JSON: %s", what, snippet(r.body))
		}
	}
}

func checkCORS(ctx context.Context, e *env, t *T) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodOptions, e.cfg.restBase+"/v1/news/sources", nil)
	req.Header.Set("Origin", "https://app.tum.de")
	req.Header.Set("Access-Control-Request-Method", "GET")
	res, err := e.http.Do(req)
	if err != nil {
		t.Failf("OPTIONS /v1/news/sources: %v", err)
		return
	}
	_ = res.Body.Close()
	if acao := res.Header.Get("Access-Control-Allow-Origin"); acao != "" {
		t.Infof("CORS preflight -> %d, Access-Control-Allow-Origin: %s", res.StatusCode, acao)
	} else {
		t.Infof("no CORS headers (HTTP %d); browsers (e.g. a Flutter web build) can't call the REST API", res.StatusCode)
	}
}

// checkBurst fires concurrent requests and reports error rate and latency percentiles.
func checkBurst(ctx context.Context, e *env, t *T) {
	n := e.cfg.burst
	if n <= 0 {
		t.Skipf("disabled with -burst 0")
		return
	}
	calls := map[string]func(context.Context) error{
		"ListNewsSources": func(ctx context.Context) error {
			_, err := e.c.ListNewsSources(ctx, &pb.ListNewsSourcesRequest{})
			return err
		},
		"ListStudentClub": func(ctx context.Context) error {
			_, err := e.c.ListStudentGroup(ctx, &pb.ListStudentGroupRequest{})
			return err
		},
		"GetCanteenHeadCount": func(ctx context.Context) error {
			_, err := e.c.GetCanteenHeadCount(ctx, &pb.GetCanteenHeadCountRequest{CanteenId: probeCanteens[0]})
			return err
		},
	}
	names := make([]string, 0, len(calls))
	for name := range calls {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		var mu sync.Mutex
		var durs []time.Duration
		var errs []error
		var wg sync.WaitGroup
		for range n {
			wg.Add(1)
			go func() {
				defer wg.Done()
				start := time.Now()
				err := calls[name](ctx)
				mu.Lock()
				defer mu.Unlock()
				durs = append(durs, time.Since(start))
				if err != nil {
					errs = append(errs, err)
				}
			}()
		}
		wg.Wait()
		slices.Sort(durs)
		p50, p95 := durs[len(durs)/2], durs[min(len(durs)-1, len(durs)*95/100)]
		t.Infof("%s x%d: p50 %s, p95 %s, %d errors", name, n, p50.Round(time.Millisecond), p95.Round(time.Millisecond), len(errs))
		if len(errs) > 0 {
			t.Failf("%s: %d/%d concurrent requests failed, e.g. %v", name, len(errs), n, errs[0])
		}
		if p95 > e.cfg.slow {
			t.Warnf("%s: p95 %s under %d concurrent requests", name, p95.Round(time.Millisecond), n)
		}
	}
}
