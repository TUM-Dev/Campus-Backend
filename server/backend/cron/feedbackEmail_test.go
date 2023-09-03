package cron

import (
	"database/sql"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestIterate(t *testing.T) {
	assert.Equal(t, []int32(nil), iterate(0))
	assert.Equal(t, []int32{0}, iterate(1))
	assert.Equal(t, []int32{0, 1}, iterate(2))
	assert.Equal(t, []int32{0, 1, 2}, iterate(3))
	assert.Equal(t, []int32{0, 1, 2, 3}, iterate(4))
	assert.Equal(t, 42, len(iterate(42)))
}

func fullFeedback() *model.Feedback {
	return &model.Feedback{
		EmailId:    sql.NullString{String: "magic-id", Valid: true},
		Receiver:   sql.NullString{String: "tca", Valid: true},
		ReplyTo:    sql.NullString{String: "test@example.de", Valid: true},
		Feedback:   sql.NullString{String: "This is a Test", Valid: true},
		ImageCount: 1,
		Latitude:   sql.NullFloat64{Float64: 0, Valid: true},
		Longitude:  sql.NullFloat64{Float64: 0, Valid: true},
		AppVersion: sql.NullString{String: "TCA 10.2", Valid: true},
		OsVersion:  sql.NullString{String: "Android 10.0", Valid: true},
		Timestamp:  sql.NullTime{Time: time.Time{}, Valid: true},
	}
}

func emptyFeedback() *model.Feedback {
	return &model.Feedback{
		EmailId:    sql.NullString{Valid: false},
		Receiver:   sql.NullString{Valid: false},
		ReplyTo:    sql.NullString{Valid: false},
		Feedback:   sql.NullString{Valid: false},
		ImageCount: 0,
		Latitude:   sql.NullFloat64{Valid: false},
		Longitude:  sql.NullFloat64{Valid: false},
		AppVersion: sql.NullString{Valid: false},
		OsVersion:  sql.NullString{Valid: false},
		Timestamp:  sql.NullTime{Valid: false},
	}
}

func TestHeaderInstantiationWithFullFeedback(t *testing.T) {
	require.NoError(t, os.Setenv("SMTP_USERNAME", "outgoing@example.de"))
	fb := fullFeedback()
	m := messageWithHeaders(fb)
	assert.Equal(t, []string{`"TUM Campus App" <outgoing@example.de>`}, m.GetHeader("From"))
	assert.Equal(t, []string{fb.Receiver.String}, m.GetHeader("To"))
	assert.Equal(t, []string{"test@example.de"}, m.GetHeader("Reply-To"))
	assert.Equal(t, []string{fb.Timestamp.Time.Format(time.RFC1123Z)}, m.GetHeader("Date"))
	assert.Equal(t, []string{"Feedback via Tum Campus App"}, m.GetHeader("Subject"))
}

func TestHeaderInstantiationWithEmptyFeedback(t *testing.T) {
	require.NoError(t, os.Setenv("SMTP_USERNAME", "outgoing@example.de"))
	m := messageWithHeaders(emptyFeedback())
	assert.Equal(t, []string{`"TUM Campus App" <outgoing@example.de>`}, m.GetHeader("From"))
	assert.Equal(t, []string{"app@tum.de"}, m.GetHeader("To"))
	assert.Equal(t, []string(nil), m.GetHeader("Reply-To"))
	// Date is set to now in messageWithHeaders => checking that this is actually now is a bit tricker
	dates := m.GetHeader("Date")
	assert.Equal(t, 1, len(dates))
	date, err := time.Parse(time.RFC1123Z, dates[0])
	require.NoError(t, err)
	assert.WithinDuration(t, time.Now(), date, time.Minute)
	assert.Equal(t, []string{"Feedback via Tum Campus App"}, m.GetHeader("Subject"))
}

func TestTemplatingResultsWithFullFeedback(t *testing.T) {
	html, txt, err := parseTemplates()
	require.NoError(t, err)
	htmlBody, txtBody, err := generateTemplatedMail(html, txt, fullFeedback())
	require.NoError(t, err)
	assert.Equal(t, `<h1>Feedback via TumCampusApp:</h1>
<blockquote>This is a Test</blockquote>
<table>
    <tr>
        <th>Inforation type</th>
        <th>Details</th>
    </tr>
    <tr>
        <th>Nutzer-Standort</th>
        <td>
            <a href="https://www.google.com/maps/search/?api=1&query=0,0">
                latitude: 0, longitude: 0
            </a>
        </td>
    </tr>
    <tr>
        <th>OS-Version</th>
        <td>Android 10.0</td>
    </tr>
    <tr>
        <th>App-Version</th>
        <td>TCA 10.2</td>
    </tr>
</table>
<h2>Fotos:</h2><br/>
<ol>
    <li>
        <a href="https://app.tum.de/File/feedback/0/0.png">Foto 0</a>
    </li>
</ol>`, htmlBody)
	assert.Equal(t, `Feedback via TumCampusApp:

This is a Test

Metadata:
- Nutzer-Standort: 0,0 (latitude,longitude)
  https://www.google.com/maps/search/?api=1&query=0,0
- OS-Version: Android 10.0
- App-Version: TCA 10.2

Photos:
- Photo 0: https://app.tum.de/File/feedback/0/0.png`, txtBody)
}

func TestTemplatingResultsWithEmptyFeedback(t *testing.T) {
	html, txt, err := parseTemplates()
	require.NoError(t, err)
	htmlBody, txtBody, err := generateTemplatedMail(html, txt, emptyFeedback())
	require.NoError(t, err)
	assert.Equal(t, `<h1>Feedback via TumCampusApp:</h1>
<i>no feedback provided</i>
<table>
    <tr>
        <th>Inforation type</th>
        <th>Details</th>
    </tr>
    <tr>
        <th>OS-Version</th>
        <td>unknown</td>
    </tr>
    <tr>
        <th>App-Version</th>
        <td>unknown</td>
    </tr>
</table>`, htmlBody)
	assert.Equal(t, `Feedback via TumCampusApp:

no feedback provided

Metadata:
- OS-Version: unknown
- App-Version: unknown`, txtBody)
}
