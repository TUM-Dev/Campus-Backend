package fcm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientSend(t *testing.T) {
	var received sendRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.NoError(t, json.NewDecoder(r.Body).Decode(&received))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"projects/test/messages/1"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL)
	err := client.Send(context.Background(), Message{
		Token: "tok-123",
		Title: "Kaiserschmarrn is back!",
		Body:  "Mensa Garching serves Kaiserschmarrn tomorrow.",
		Data:  map[string]string{"type": "favorite_dish"},
	})
	require.NoError(t, err)
	require.Equal(t, "tok-123", received.Message.Token)
	require.Equal(t, "Kaiserschmarrn is back!", received.Message.Notification.Title)
	require.Equal(t, "favorite_dish", received.Message.Data["type"])
}

func TestClientSendReturnsErrorOnNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "invalid token", http.StatusBadRequest)
	}))
	defer srv.Close()

	client := NewClient(srv.Client(), srv.URL)
	err := client.Send(context.Background(), Message{Token: "bad-tok", Title: "Hi"})
	require.ErrorContains(t, err, "HTTP 400")
}
