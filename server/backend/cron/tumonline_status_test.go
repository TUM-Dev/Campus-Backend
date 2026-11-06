package cron

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func statusPage(current, future string) string {
	return `<h2>Aktuelle Meldungen</h2><div class="outer-grid"><div class="table-headers"></div>` + current +
		`</div><h2>Angekündigte Meldungen</h2>` + future + `<h2>Beendete Meldungen ↑</h2>`
}

func statusEntry(service, title string) string {
	return `<a class="table-entry"><div class="entry-service"><span class="mobile-bold">` + service +
		`</span></div><div class="entry-problemtype">` + title + `</div></a>`
}

func TestParseTUMOnlineStatus(t *testing.T) {
	for _, tc := range []struct {
		name, current, future string
		want                  tumonlineStatus
	}{
		{name: "all clear"},
		{name: "unrelated service", current: statusEntry("EduVPN", "Wartung")},
		{name: "title only is not affected service", current: statusEntry("EduVPN", "TUMonline Wartung")},
		{name: "upcoming ignored", future: statusEntry("TUMonline", "Wartung")},
		{name: "maintenance", current: statusEntry("TUMonline", "Wartungsarbeiten"), want: tumonlineStatus{Maintenance: true}},
		{name: "English maintenance", current: statusEntry("campus.tum.de", "Maintenance"), want: tumonlineStatus{Maintenance: true}},
		{name: "degradation", current: statusEntry("TUMonline", "Performance eingeschränkt"), want: tumonlineStatus{Degraded: true}},
		{name: "outage", current: statusEntry("TUMonline &amp; angeschlossene Systeme", "Nicht erreichbar"), want: tumonlineStatus{Degraded: true}},
		{name: "maintenance priority", current: statusEntry("TUMonline", "Wartung") + statusEntry("TUMonline", "Ausfall"), want: tumonlineStatus{Maintenance: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseTUMOnlineStatus(strings.NewReader(statusPage(tc.current, tc.future)))
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestParseTUMOnlineStatusRealPage(t *testing.T) {
	body, err := os.ReadFile("testdata/tumonline_status.html")
	require.NoError(t, err)
	got, err := parseTUMOnlineStatusAt(strings.NewReader(string(body)), time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Equal(t, tumonlineStatus{
		Degraded:  true,
		MessageDE: "Die Dienste von TUMonline sind derzeit beeinträchtigt. Voraussichtliches Ende: 08.11.2026, 09:55 CET.",
		MessageEN: "TUMonline services are currently degraded. Estimated end: 8 Nov 2026, 09:55 CET.",
	}, got)
	// An elapsed estimate clears custom text while retaining the published warning.
	got, err = parseTUMOnlineStatusAt(strings.NewReader(string(body)), time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Equal(t, tumonlineStatus{Degraded: true}, got)
}

func TestTUMOnlineStatusMessageEstimates(t *testing.T) {
	now := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	entry := func(title, end string) string {
		return strings.Replace(statusEntry("TUMonline", title), "</a>", `<div class="entry-end entry-date">`+end+`</div></a>`, 1)
	}
	for _, tc := range []struct {
		name, entries, wantDE, wantEN string
	}{
		{"summer time", entry("Wartung", "06.07.2026 18:00"),
			"TUMonline wird derzeit gewartet. Voraussichtliches Ende: 06.07.2026, 18:00 CEST.",
			"TUMonline is currently undergoing maintenance. Estimated end: 6 Jul 2026, 18:00 CEST."},
		{"English date", entry("Maintenance", "7/6/2026 6:00\u202fPM"),
			"TUMonline wird derzeit gewartet. Voraussichtliches Ende: 06.07.2026, 18:00 CEST.",
			"TUMonline is currently undergoing maintenance. Estimated end: 6 Jul 2026, 18:00 CEST."},
		{"elapsed estimate", entry("Wartung", "06.07.2026 13:00"),
			"", ""},
		{"unrecognized estimate", entry("Wartung", "unknown"),
			"", ""},
		{"multiple notices", entry("Wartung", "06.07.2026 18:00") + entry("Maintenance", "07.07.2026 08:00"),
			"TUMonline wird derzeit gewartet. Voraussichtliches Ende: 07.07.2026, 08:00 CEST.",
			"TUMonline is currently undergoing maintenance. Estimated end: 7 Jul 2026, 08:00 CEST."},
		{"one unknown end", entry("Wartung", "06.07.2026 18:00") + entry("Maintenance", "unknown"),
			"", ""},
		{"maintenance priority uses maintenance estimate", entry("Ausfall", "07.07.2026 08:00") + entry("Wartung", "06.07.2026 18:00"),
			"TUMonline wird derzeit gewartet. Voraussichtliches Ende: 06.07.2026, 18:00 CEST.",
			"TUMonline is currently undergoing maintenance. Estimated end: 6 Jul 2026, 18:00 CEST."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseTUMOnlineStatusAt(strings.NewReader(statusPage(tc.entries, "")), now)
			require.NoError(t, err)
			require.True(t, got.Maintenance)
			require.False(t, got.Degraded)
			require.Equal(t, tc.wantDE, got.MessageDE)
			require.Equal(t, tc.wantEN, got.MessageEN)
		})
	}
}

func TestParseTUMOnlineStatusRejectsChangedMarkup(t *testing.T) {
	for _, body := range []string{
		"<html>Sign in</html>",
		"<h2>Aktuelle Meldungen</h2>",
		statusPage(`<a class="table-entry">changed markup</a>`, ""),
		strings.ReplaceAll(statusPage("", ""), "table-headers", "changed"),
	} {
		_, err := parseTUMOnlineStatus(strings.NewReader(body))
		require.Error(t, err)
	}
}

func TestTUMOnlineStatusCronDisabled(t *testing.T) {
	t.Setenv("FIREBASE_PROJECT_ID", "example")
	t.Setenv("TUMONLINE_STATUS_SYNC_ENABLED", "false")
	require.NoError(t, (&CronService{}).tumonlineStatusCron())
	t.Setenv("TUMONLINE_STATUS_SYNC_ENABLED", "true")
	t.Setenv("FIREBASE_PROJECT_ID", "")
	require.NoError(t, (&CronService{}).tumonlineStatusCron())
}

func TestTUMOnlineStatusFetchFailures(t *testing.T) {
	for _, tc := range []struct {
		code int
		body string
		err  error
	}{
		{code: 503, body: "unavailable"},
		{code: 200, body: "<html>Sign in</html>"},
		{code: 200, body: strings.Repeat("x", (2<<20)+1)},
		{err: errors.New("network failure")},
	} {
		client := &http.Client{Transport: httpTestTransport(func(*http.Request) (*http.Response, error) {
			if tc.err != nil {
				return nil, tc.err
			}
			return &http.Response{StatusCode: tc.code, Body: io.NopCloser(strings.NewReader(tc.body))}, nil
		})}
		_, err := fetchTUMOnlineStatus(context.Background(), client)
		require.Error(t, err)
	}
}
