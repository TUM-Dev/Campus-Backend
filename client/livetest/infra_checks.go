package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc/codes"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func infraChecks() []check {
	return []check{
		{group: "infra", name: "health", fn: checkHealth},
		{group: "infra", name: "tls-certificate", fn: checkTLS},
		{group: "infra", name: "https-redirect", fn: checkHTTPSRedirect},
		{group: "infra", name: "grpc-reflection-drift", fn: checkReflectionDrift},
		{group: "infra", name: "swagger", fn: checkSwagger},
		{group: "infra", name: "metrics-exposure", fn: checkMetricsExposure},
		{group: "infra", name: "files-listing", fn: checkFilesListing},
		{group: "infra", name: "unknown-route", fn: checkUnknownRoute},
	}
}

func campusService() protoreflect.ServiceDescriptor {
	return pb.File_tumdev_campus_backend_proto.Services().ByName("Campus")
}

func checkHealth(ctx context.Context, e *env, t *T) {
	r, ok := e.restOK(ctx, t, http.MethodGet, "/health", nil)
	if ok && strings.TrimSpace(string(r.body)) != "healthy" {
		t.Warnf("/health returned %q, want \"healthy\"", snippet(r.body))
	}
}

func checkTLS(ctx context.Context, e *env, t *T) {
	hosts := map[string]bool{}
	if u, err := url.Parse(e.cfg.restBase); err == nil && u.Scheme == "https" {
		host := u.Host
		if u.Port() == "" {
			host = net.JoinHostPort(u.Hostname(), "443")
		}
		hosts[host] = true
	}
	if !e.cfg.plaintext {
		hosts[e.cfg.grpcAddr] = true
	}
	if len(hosts) == 0 {
		t.Skipf("no TLS endpoints configured")
		return
	}
	for host := range hosts {
		d := tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}}
		conn, err := d.DialContext(ctx, "tcp", host)
		if err != nil {
			t.Failf("%s: TLS handshake failed: %v", host, err)
			continue
		}
		state := conn.(*tls.Conn).ConnectionState()
		_ = conn.Close()
		leaf := state.PeerCertificates[0]
		left := time.Until(leaf.NotAfter)
		switch {
		case left < 7*24*time.Hour:
			t.Failf("%s: certificate expires in %s (%s)", host, left.Round(time.Hour), leaf.NotAfter.Format(time.DateOnly))
		case left < 21*24*time.Hour:
			t.Warnf("%s: certificate expires in %s (%s)", host, left.Round(time.Hour), leaf.NotAfter.Format(time.DateOnly))
		default:
			t.Infof("%s: certificate valid until %s, %s", host, leaf.NotAfter.Format(time.DateOnly), tls.VersionName(state.Version))
		}
	}
}

func checkHTTPSRedirect(ctx context.Context, e *env, t *T) {
	u, err := url.Parse(e.cfg.restBase)
	if err != nil || u.Scheme != "https" {
		t.Skipf("REST base is not https")
		return
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+u.Host+"/health", nil)
	client := &http.Client{Timeout: 10 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	res, err := client.Do(req)
	if err != nil {
		t.Warnf("plain HTTP not reachable: %v", err)
		return
	}
	_ = res.Body.Close()
	loc := res.Header.Get("Location")
	if res.StatusCode >= 300 && res.StatusCode < 400 && strings.HasPrefix(loc, "https://") {
		t.Infof("http:// redirects with %d to %s", res.StatusCode, loc)
		return
	}
	t.Warnf("http://%s/health returned %d (Location %q) instead of redirecting to https", u.Host, res.StatusCode, loc)
}

// checkReflectionDrift compares the service the server actually exposes with
// the proto this harness was built against.
func checkReflectionDrift(ctx context.Context, e *env, t *T) {
	stream, err := rpb.NewServerReflectionClient(e.conn).ServerReflectionInfo(ctx)
	if err != nil {
		t.Skipf("reflection unavailable: %v", err)
		return
	}
	defer func() { _ = stream.CloseSend() }()
	svc := campusService()
	if err := stream.Send(&rpb.ServerReflectionRequest{
		MessageRequest: &rpb.ServerReflectionRequest_FileContainingSymbol{FileContainingSymbol: string(svc.FullName())},
	}); err != nil {
		t.Skipf("reflection request failed: %v", err)
		return
	}
	res, err := stream.Recv()
	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			t.Skipf("server does not support reflection v1")
		} else {
			t.Skipf("reflection response failed: %v", err)
		}
		return
	}
	if er := res.GetErrorResponse(); er != nil {
		t.Failf("server does not know %s: %s", svc.FullName(), er.GetErrorMessage())
		return
	}

	remote := map[string]*descriptorpb.MethodDescriptorProto{}
	for _, raw := range res.GetFileDescriptorResponse().GetFileDescriptorProto() {
		fd := &descriptorpb.FileDescriptorProto{}
		if err := proto.Unmarshal(raw, fd); err != nil {
			t.Failf("could not decode remote descriptor: %v", err)
			return
		}
		for _, s := range fd.GetService() {
			if fd.GetPackage()+"."+s.GetName() == string(svc.FullName()) {
				for _, m := range s.GetMethod() {
					remote[m.GetName()] = m
				}
			}
		}
	}
	if len(remote) == 0 {
		t.Failf("reflection returned no methods for %s", svc.FullName())
		return
	}

	served := len(remote)
	methods := svc.Methods()
	for i := range methods.Len() {
		m := methods.Get(i)
		rm, ok := remote[string(m.Name())]
		if !ok {
			t.Failf("%s is in the proto but not served (deployment older than proto?)", m.Name())
			continue
		}
		if strings.TrimPrefix(rm.GetInputType(), ".") != string(m.Input().FullName()) ||
			strings.TrimPrefix(rm.GetOutputType(), ".") != string(m.Output().FullName()) {
			t.Failf("%s has a different signature on the server: %s -> %s", m.Name(), rm.GetInputType(), rm.GetOutputType())
		}
		delete(remote, string(m.Name()))
	}
	for name := range remote {
		t.Warnf("server exposes %s which this harness doesn't know: bump the server module in client/go.mod", name)
	}
	t.Infof("server exposes %d methods, proto has %d", served, methods.Len())
}

var pathParam = regexp.MustCompile(`\{[^}]+\}`)

// checkSwagger verifies that /swagger/swagger.json is served and documents every HTTP route of the proto.
func checkSwagger(ctx context.Context, e *env, t *T) {
	r, ok := e.restOK(ctx, t, http.MethodGet, "/swagger/swagger.json", nil)
	if !ok {
		return
	}
	var doc struct {
		BasePath string                    `json:"basePath"`
		Paths    map[string]map[string]any `json:"paths"`
	}
	if err := json.Unmarshal(r.body, &doc); err != nil {
		t.Failf("swagger.json is not valid JSON: %v", err)
		return
	}
	normalize := func(p string) string { return pathParam.ReplaceAllString(p, "{}") }
	documented := map[string]bool{}
	for p, ops := range doc.Paths {
		for verb := range ops {
			documented[strings.ToUpper(verb)+" "+normalize(p)] = true
		}
	}
	methods := campusService().Methods()
	missing := 0
	for i := range methods.Len() {
		m := methods.Get(i)
		rule, _ := proto.GetExtension(m.Options(), annotations.E_Http).(*annotations.HttpRule)
		if rule == nil {
			continue
		}
		verb, path := httpRule(rule)
		if !documented[verb+" "+normalize(path)] {
			missing++
			t.Warnf("%s (%s %s) is missing from the served swagger.json", m.Name(), verb, path)
		}
	}
	if missing == 0 {
		t.Infof("swagger.json documents all %d HTTP routes", len(documented))
	}
}

func httpRule(r *annotations.HttpRule) (string, string) {
	switch p := r.GetPattern().(type) {
	case *annotations.HttpRule_Get:
		return "GET", p.Get
	case *annotations.HttpRule_Post:
		return "POST", p.Post
	case *annotations.HttpRule_Put:
		return "PUT", p.Put
	case *annotations.HttpRule_Delete:
		return "DELETE", p.Delete
	case *annotations.HttpRule_Patch:
		return "PATCH", p.Patch
	}
	return "", ""
}

func checkMetricsExposure(ctx context.Context, e *env, t *T) {
	r, err := e.rest(ctx, http.MethodGet, "/metrics", nil)
	if err != nil {
		t.Failf("GET /metrics: %v", err)
		return
	}
	if r.code == http.StatusOK && strings.Contains(string(r.body), "# HELP") {
		t.Warnf("/metrics (Prometheus) is publicly reachable (%d bytes); consider restricting it at the proxy", len(r.body))
		return
	}
	t.Infof("/metrics not public (HTTP %d)", r.code)
}

var listingEntry = regexp.MustCompile(`<a href="([^"]+)">`)

// checkFilesListing flags directory indexes under /files/. It never downloads files,
// it only counts the entries of the index pages.
func checkFilesListing(ctx context.Context, e *env, t *T) {
	r, err := e.rest(ctx, http.MethodGet, "/files/", nil)
	if err != nil {
		t.Failf("GET /files/: %v", err)
		return
	}
	entries := listingEntry.FindAllStringSubmatch(string(r.body), -1)
	if r.code != http.StatusOK || len(entries) == 0 {
		t.Infof("/files/ has no directory index (HTTP %d)", r.code)
		return
	}
	t.Warnf("/files/ serves a directory index (%d entries); http.FileServer lists directories by default", len(entries))

	// directories that hold user uploads must never be enumerable
	sensitive := map[string]string{
		"feedback/": "feedback attachments sent by users",
		"rating/":   "images attached to ratings",
		"tmp/":      "temporary files",
		"storage/":  "raw storage",
	}
	for _, m := range entries {
		dir := m[1]
		why, ok := sensitive[dir]
		if !ok {
			continue
		}
		sub, err := e.rest(ctx, http.MethodGet, "/files/"+dir, nil)
		if err != nil {
			continue
		}
		n := len(listingEntry.FindAllString(string(sub.body), -1))
		if sub.code == http.StatusOK && n > 0 {
			t.Failf("/files/%s is publicly enumerable (%d entries of %s)", dir, n, why)
		}
	}
}

func checkUnknownRoute(ctx context.Context, e *env, t *T) {
	for _, p := range []string{"/v1/livetest-does-not-exist", "/livetest-does-not-exist"} {
		r, err := e.rest(ctx, http.MethodGet, p, nil)
		if err != nil {
			t.Failf("GET %s: %v", p, err)
			continue
		}
		switch {
		case r.code == http.StatusNotFound:
			t.Infof("GET %s -> 404", p)
		case r.code >= 300 && r.code < 400, r.code == http.StatusOK && !strings.HasPrefix(p, "/v1/"):
			// "/" redirects everything unknown to /swagger/
			t.Infof("GET %s -> %d", p, r.code)
		default:
			t.Warnf("GET %s -> %d, want 404: %s", p, r.code, snippet(r.body))
		}
	}
}
