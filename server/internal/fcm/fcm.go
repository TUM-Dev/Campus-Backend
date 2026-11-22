// Package fcm sends push notifications via the Firebase Cloud Messaging HTTP v1 API.
package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/oauth2/google"
)

const messagingScope = "https://www.googleapis.com/auth/firebase.messaging"

// Message represents a single FCM notification addressed to a device registration token.
type Message struct {
	Token string
	Title string
	Body  string
	Data  map[string]string
}

// Sender sends push notifications; implementations may be mocked in tests.
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

// Client sends messages over the FCM HTTP v1 REST API.
type Client struct {
	httpClient *http.Client
	endpoint   string
}

// NewFromEnv returns a configured FCM Client when FCM_ENABLED=true and FIREBASE_PROJECT_ID is set.
// When disabled, it returns (nil, nil) so callers can treat push dispatch as an opt-in no-op.
func NewFromEnv(ctx context.Context) (*Client, error) {
	project := os.Getenv("FIREBASE_PROJECT_ID")
	if project == "" || os.Getenv("FCM_ENABLED") != "true" {
		return nil, nil
	}
	httpClient, err := google.DefaultClient(ctx, messagingScope)
	if err != nil {
		return nil, fmt.Errorf("initializing FCM OAuth client: %w", err)
	}
	httpClient.Timeout = 15 * time.Second
	endpoint := "https://fcm.googleapis.com/v1/projects/" + url.PathEscape(project) + "/messages:send"
	return &Client{httpClient: httpClient, endpoint: endpoint}, nil
}

// NewClient constructs a Client with an explicit HTTP client and endpoint (useful for tests).
func NewClient(httpClient *http.Client, endpoint string) *Client {
	return &Client{httpClient: httpClient, endpoint: endpoint}
}

type sendRequest struct {
	Message wireMessage `json:"message"`
}

type wireMessage struct {
	Token        string            `json:"token"`
	Notification *wireNotification `json:"notification,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
}

type wireNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Send dispatches a single FCM notification.
func (c *Client) Send(ctx context.Context, msg Message) error {
	if c == nil {
		return nil
	}
	if msg.Token == "" {
		return fmt.Errorf("fcm token is empty")
	}
	payload := sendRequest{
		Message: wireMessage{
			Token: msg.Token,
			Notification: &wireNotification{
				Title: msg.Title,
				Body:  msg.Body,
			},
			Data: msg.Data,
		},
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("fcm returned HTTP %d: %s", resp.StatusCode, string(body))
	}
	return nil
}
