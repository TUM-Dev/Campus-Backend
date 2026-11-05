package main

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/dynamicpb"
)

// emptyRequestSafe lists the RPCs that may be called with an empty request on
// production. Each entry was verified to have no side effects for an empty request:
// reads, or writes that fail validation before touching the database.
// New RPCs are skipped until someone checks them and adds them here.
var emptyRequestSafe = map[string]string{
	"ListNewsAlerts":           "read",
	"ListNewsSources":          "read",
	"ListNews":                 "read",
	"ListCanteenRatings":       "read",
	"GetDishRatings":           "read",
	"ListAvailableDishTags":    "read",
	"ListNameTags":             "read",
	"ListAvailableCanteenTags": "read",
	"ListCanteens":             "read",
	"ListDishes":               "read (rejected: year/week out of range)",
	"GetUpdateNote":            "read",
	"ListMovies":               "read",
	"GetCanteenHeadCount":      "read",
	"ListStudentClub":          "read",
	"CreateCanteenRating":      "rejected: canteen \"\" does not exist",
	"CreateDishRating":         "rejected: canteen \"\" does not exist",
	"CreateDevice":             "rejected: deviceId is empty",
	"DeleteDevice":             "rejected: deviceId is empty",
	"CreateFeedback":           "rejected: neither message nor attachment",
}

// acceptableForEmpty are the codes a well-behaved handler returns for an empty request.
var acceptableForEmpty = []codes.Code{codes.OK, codes.InvalidArgument, codes.NotFound, codes.FailedPrecondition, codes.OutOfRange}

func sweepChecks() []check {
	var out []check
	methods := campusService().Methods()
	for i := range methods.Len() {
		m := methods.Get(i)
		name := string(m.Name())
		fullName := fmt.Sprintf("/%s/%s", campusService().FullName(), name)
		out = append(out, check{group: "sweep", name: name, fn: func(ctx context.Context, e *env, t *T) {
			why, ok := emptyRequestSafe[name]
			if !ok {
				t.Warnf("not in the sweep allowlist; verify an empty request has no side effects and add it to emptyRequestSafe")
				return
			}
			var err error
			if m.IsStreamingClient() {
				err = t.timed(name, func() error { return emptyFeedbackStream(ctx, e) })
			} else {
				in, out := dynamicpb.NewMessage(m.Input()), dynamicpb.NewMessage(m.Output())
				err = t.timed(name, func() error { return e.conn.Invoke(ctx, fullName, in, out) })
			}
			judgeSweep(t, name, why, err)
		}})
	}
	return out
}

func judgeSweep(t *T, name, why string, err error) {
	st := status.Convert(err)
	code := st.Code()
	// the code grpc-go generates for methods that the server struct does not implement
	if code == codes.Unimplemented && strings.HasPrefix(st.Message(), "method ") && strings.HasSuffix(st.Message(), " not implemented") {
		t.Failf("%s is declared in the proto but the server falls back to Unimplemented (%q); is the handler misnamed?", name, st.Message())
		return
	}
	for _, ok := range acceptableForEmpty {
		if code == ok {
			t.Infof("empty request -> %s (%s)", code, why)
			return
		}
	}
	switch code {
	case codes.Internal, codes.Unknown:
		t.Warnf("empty request -> %s %q; validation errors should be InvalidArgument", code, st.Message())
	case codes.Unimplemented:
		t.Warnf("empty request -> Unimplemented %q", st.Message())
	default:
		t.Failf("empty request -> %s %q", code, st.Message())
	}
}

func emptyFeedbackStream(ctx context.Context, e *env) error {
	stream, err := e.c.CreateFeedback(ctx)
	if err != nil {
		return err
	}
	_, err = stream.CloseAndRecv()
	return err
}
