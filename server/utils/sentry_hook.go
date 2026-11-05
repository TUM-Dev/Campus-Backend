package utils

import (
	"fmt"
	"reflect"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

var sentryLevels = map[log.Level]sentry.Level{
	log.TraceLevel: sentry.LevelDebug,
	log.DebugLevel: sentry.LevelDebug,
	log.InfoLevel:  sentry.LevelInfo,
	log.WarnLevel:  sentry.LevelWarning,
	log.ErrorLevel: sentry.LevelError,
	log.FatalLevel: sentry.LevelFatal,
	log.PanicLevel: sentry.LevelFatal,
}

// sentryHook reports logrus entries as sentry events (issues).
// sentry-go's own logrus integration only emits sentry logs since v0.48.0, so we keep a small hook here.
type sentryHook struct {
	levels []log.Level
}

func newSentryHook(levels []log.Level) sentryHook {
	return sentryHook{levels: levels}
}

func (h sentryHook) Levels() []log.Level {
	return h.levels
}

func (h sentryHook) Fire(entry *log.Entry) error {
	hub := sentry.CurrentHub()
	if entry.Context != nil {
		if ctxHub := sentry.GetHubFromContext(entry.Context); ctxHub != nil {
			hub = ctxHub
		}
	}
	client := hub.Client()
	if client == nil || !client.IsEnabled() {
		return nil
	}

	event := sentry.NewEvent()
	event.Level = sentryLevels[entry.Level]
	event.Message = entry.Message

	if len(entry.Data) > 0 {
		fields := make(sentry.Context, len(entry.Data))
		for k, v := range entry.Data {
			if err, ok := v.(error); ok {
				// errors usually marshal to {}, which is useless in sentry
				fields[k] = err.Error()
			} else {
				fields[k] = v
			}
		}
		event.Contexts["logrus"] = fields
	}

	if err, ok := entry.Data[log.ErrorKey].(error); ok {
		exception := sentry.Exception{
			Type:  reflect.TypeOf(err).String(),
			Value: fmt.Sprint(err),
		}
		if client.Options().AttachStacktrace {
			exception.Stacktrace = sentry.ExtractStacktrace(err)
		}
		event.Exception = []sentry.Exception{exception}
	}

	// hub.CaptureEvent does not accept a hint, so go through the client to keep the entry's context
	client.CaptureEvent(event, &sentry.EventHint{Context: entry.Context}, hub.Scope())
	return nil
}
