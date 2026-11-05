package utils

import (
	"context"
	"errors"
	"testing"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestSentryHookFire(t *testing.T) {
	transport := &sentry.MockTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              "https://public@example.com/1",
		Transport:        transport,
		AttachStacktrace: true,
	})
	require.NoError(t, err)
	hub := sentry.NewHub(client, sentry.NewScope())

	logger := log.New()
	logger.AddHook(newSentryHook([]log.Level{log.ErrorLevel}))
	ctx := sentry.SetHubOnContext(context.Background(), hub)

	logger.WithContext(ctx).Info("not reported")
	logger.WithContext(ctx).WithError(errors.New("boom")).WithField("user", 42).Error("something failed")

	events := transport.Events()
	require.Len(t, events, 1)
	event := events[0]
	require.Equal(t, sentry.LevelError, event.Level)
	require.Equal(t, "something failed", event.Message)
	require.Equal(t, 42, event.Contexts["logrus"]["user"])
	require.Equal(t, "boom", event.Contexts["logrus"][log.ErrorKey])
	require.Len(t, event.Exception, 1)
	require.Equal(t, "boom", event.Exception[0].Value)
}
