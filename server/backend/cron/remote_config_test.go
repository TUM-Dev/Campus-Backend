package cron

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublishTUMOnlineStatusMergesAndRetries(t *testing.T) {
	gets, puts := 0, 0
	client := &http.Client{Transport: httpTestTransport(func(req *http.Request) (*http.Response, error) {
		if req.Method == http.MethodGet {
			gets++
			require.Equal(t, "gzip", req.Header.Get("Accept-Encoding"))
			// The second read includes a concurrent console edit, which must survive.
			value := "old"
			if gets > 1 {
				value = "concurrent edit"
			}
			raw := `{"parameters":{"unrelated":{"defaultValue":{"value":"` + value + `"}}},"parameterGroups":{"status":{"description":"keep group","parameters":{"isTUMOnlineDegraded":{"defaultValue":{"value":"true"},"valueType":"BOOLEAN","description":"keep description"}}}},"conditions":[{"name":"test","expression":"true"}],"version":{"versionNumber":"1"}}`
			var compressed bytes.Buffer
			gz := gzip.NewWriter(&compressed)
			_, err := gz.Write([]byte(raw))
			require.NoError(t, err)
			require.NoError(t, gz.Close())
			return &http.Response{StatusCode: 200, Header: http.Header{"Etag": {`"v` + string(rune('0'+gets)) + `"`}, "Content-Encoding": {"gzip"}}, Body: io.NopCloser(&compressed)}, nil
		}
		puts++
		require.Equal(t, http.MethodPut, req.Method)
		require.Equal(t, `"v`+string(rune('0'+gets))+`"`, req.Header.Get("If-Match"))
		var template map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&template))
		require.NotContains(t, template, "version")
		parameters := template["parameters"].(map[string]any)
		require.Equal(t, "true", parameters["isTUMOnlineMaintenanceMode"].(map[string]any)["defaultValue"].(map[string]any)["value"])
		for key, want := range map[string]string{
			"tumOnlineStatusMessage_de": "TUMonline wird derzeit gewartet.",
			"tumOnlineStatusMessage_en": "TUMonline is currently undergoing maintenance.",
		} {
			message := parameters[key].(map[string]any)
			require.Equal(t, "STRING", message["valueType"])
			require.Equal(t, want, message["defaultValue"].(map[string]any)["value"])
		}
		group := template["parameterGroups"].(map[string]any)["status"].(map[string]any)
		require.Equal(t, "keep group", group["description"])
		degraded := group["parameters"].(map[string]any)["isTUMOnlineDegraded"].(map[string]any)
		require.Equal(t, "false", degraded["defaultValue"].(map[string]any)["value"])
		require.Equal(t, "keep description", degraded["description"])
		require.Contains(t, template, "conditions")
		code := http.StatusConflict
		if puts > 1 {
			code = 200
			require.Equal(t, "concurrent edit", parameters["unrelated"].(map[string]any)["defaultValue"].(map[string]any)["value"])
		}
		return &http.Response{StatusCode: code, Body: io.NopCloser(strings.NewReader("{}"))}, nil
	})}
	require.NoError(t, publishTUMOnlineStatus(context.Background(), client, "https://firebase.example/config", tumonlineStatus{
		Maintenance: true,
		MessageDE:   "TUMonline wird derzeit gewartet.",
		MessageEN:   "TUMonline is currently undergoing maintenance.",
	}))
	require.Equal(t, 2, gets)
	require.Equal(t, 2, puts)
}

func TestPublishTUMOnlineStatusSkipsUnchanged(t *testing.T) {
	client := &http.Client{Transport: httpTestTransport(func(req *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodGet, req.Method)
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"parameters":{"isTUMOnlineMaintenanceMode":{"defaultValue":{"value":"false"},"valueType":"BOOLEAN"},"isTUMOnlineDegraded":{"defaultValue":{"value":"false"},"valueType":"BOOLEAN"},"tumOnlineStatusMessage_de":{"defaultValue":{"value":""},"valueType":"STRING"},"tumOnlineStatusMessage_en":{"defaultValue":{"value":""},"valueType":"STRING"}}}`))}, nil
	})}
	require.NoError(t, publishTUMOnlineStatus(context.Background(), client, "https://firebase.example/config", tumonlineStatus{}))
}

func TestPublishTUMOnlineStatusFailures(t *testing.T) {
	for _, tc := range []struct {
		name, body, etag string
		getCode, putCode int
		wantCalls        int
	}{
		{"fetch failure", `{}`, "", 403, 0, 1},
		{"invalid JSON", `{`, "", 200, 0, 1},
		{"missing ETag", `{}`, "", 200, 0, 1},
		{"wildcard ETag", `{}`, "*", 200, 0, 1},
		{"publish failure", `{}`, `"v1"`, 200, 403, 2},
		{"conflicts bounded", `{}`, `"v1"`, 200, 409, 6},
		{"preconditions bounded", `{}`, `"v1"`, 200, 412, 6},
		{"conditional managed value", `{"parameters":{"isTUMOnlineDegraded":{"conditionalValues":{"test":{"value":"true"}}}}}`, `"v1"`, 200, 0, 1},
		{"conditional message", `{"parameters":{"tumOnlineStatusMessage_de":{"conditionalValues":{"test":{"value":"manual message"}}}}}`, `"v1"`, 200, 0, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			client := &http.Client{Transport: httpTestTransport(func(req *http.Request) (*http.Response, error) {
				calls++
				code := tc.getCode
				if req.Method == http.MethodPut {
					code = tc.putCode
				}
				return &http.Response{StatusCode: code, Header: http.Header{"Etag": {tc.etag}}, Body: io.NopCloser(strings.NewReader(tc.body))}, nil
			})}
			require.Error(t, publishTUMOnlineStatus(context.Background(), client, "https://firebase.example/config", tumonlineStatus{}))
			require.Equal(t, tc.wantCalls, calls)
		})
	}
}

func TestStatusMessagesUpdateAndClearInGroups(t *testing.T) {
	var template configObject
	require.NoError(t, json.Unmarshal([]byte(`{"parameterGroups":{"status":{"parameters":{"tumOnlineStatusMessage_de":{"defaultValue":{"value":"old DE"},"description":"keep"},"tumOnlineStatusMessage_en":{"defaultValue":{"value":"old EN"}}}}}}`), &template))
	state := tumonlineStatus{Degraded: true, MessageDE: "Beeinträchtigt.", MessageEN: "Degraded."}
	for _, want := range []tumonlineStatus{state, {Degraded: true}, {}} {
		changed, err := setStatusParameters(template, want)
		require.NoError(t, err)
		require.True(t, changed)
		var decoded map[string]any
		body, err := json.Marshal(template)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(body, &decoded))
		group := decoded["parameterGroups"].(map[string]any)["status"].(map[string]any)["parameters"].(map[string]any)
		for key, text := range map[string]string{"tumOnlineStatusMessage_de": want.MessageDE, "tumOnlineStatusMessage_en": want.MessageEN} {
			message := group[key].(map[string]any)
			require.Equal(t, "STRING", message["valueType"])
			require.Equal(t, text, message["defaultValue"].(map[string]any)["value"])
		}
		require.Equal(t, "keep", group["tumOnlineStatusMessage_de"].(map[string]any)["description"])
		changed, err = setStatusParameters(template, want)
		require.NoError(t, err)
		require.False(t, changed)
	}
}
