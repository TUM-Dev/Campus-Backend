package cron

import (
	"os"
	"testing"
	"time"

	"github.com/guregu/null"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		EmailId:    "magic-id",
		Recipient:  "tca",
		ReplyTo:    null.StringFrom("test@example.de"),
		Feedback:   "This is a Test",
		ImageCount: 1,
		Latitude:   null.FloatFrom(0),
		Longitude:  null.FloatFrom(0),
		AppVersion: null.StringFrom("TCA 10.2"),
		OsVersion:  null.StringFrom("Android 10.0"),
		Timestamp:  null.TimeFrom(time.Now()),
	}
}

func emptyFeedback() *model.Feedback {
	return &model.Feedback{
		EmailId:    "",
		Recipient:  "",
		ReplyTo:    null.String{},
		Feedback:   "",
		ImageCount: 0,
		Latitude:   null.Float{},
		Longitude:  null.Float{},
		AppVersion: null.String{},
		OsVersion:  null.String{},
		Timestamp:  null.Time{},
	}
}

func TestHeaderInstantiationWithFullFeedback(t *testing.T) {
	require.NoError(t, os.Setenv("SMTP_USERNAME", "outgoing@example.de"))
	require.NoError(t, os.Setenv("SMTP_FROM", "from@example.de"))
	fb := fullFeedback()
	m := messageWithHeaders(fb)
	assert.Equal(t, []string{`"TUM Campus App" <from@example.de>`}, m.GetHeader("From"))
	assert.Equal(t, []string{fb.Recipient}, m.GetHeader("To"))
	assert.Equal(t, []string{"test@example.de"}, m.GetHeader("Reply-To"))
	assert.Equal(t, []string{fb.Timestamp.Time.Format(time.RFC1123Z)}, m.GetHeader("Date"))
	assert.Equal(t, []string{"Feedback via the TUM Campus App"}, m.GetHeader("Subject"))
}

func TestHeaderInstantiationWithEmptyFeedback(t *testing.T) {
	require.NoError(t, os.Setenv("SMTP_USERNAME", "outgoing@example.de"))
	require.NoError(t, os.Setenv("SMTP_FROM", "from@example.de"))
	m := messageWithHeaders(emptyFeedback())
	assert.Equal(t, []string{`"TUM Campus App" <from@example.de>`}, m.GetHeader("From"))
	assert.Equal(t, []string{"app@tum.de"}, m.GetHeader("To"))
	assert.Equal(t, []string(nil), m.GetHeader("Reply-To"))
	// Date is set to now in messageWithHeaders => checking that this is actually now is a bit tricker
	dates := m.GetHeader("Date")
	assert.Equal(t, 1, len(dates))
	date, err := time.Parse(time.RFC1123Z, dates[0])
	require.NoError(t, err)
	assert.WithinDuration(t, time.Now(), date, time.Minute)
	assert.Equal(t, []string{"Feedback via the TUM Campus App"}, m.GetHeader("Subject"))
}

func TestTrunaction(t *testing.T) {
	require.Equal(t, truncate("abc", -1), "")
	require.Equal(t, truncate("abc", 0), "")
	require.Equal(t, truncate("abc", 1), "a..")
	require.Equal(t, truncate("abc", 2), "ab..")
	require.Equal(t, truncate("abc", 3), "abc")
	require.Equal(t, truncate("abc", 200), "abc")
}

func TestTemplatingResultsWithFullFeedback(t *testing.T) {
	html, txt, err := parseTemplates()
	require.NoError(t, err)
	htmlBody, txtBody, err := generateTemplatedMail(html, txt, fullFeedback())
	require.NoError(t, err)
	assert.Equal(t, `<b>Feedback via the TUM Campus App:</b>
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
	assert.Equal(t, `Feedback via the TUM Campus App:

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
	assert.Equal(t, `<b>Feedback via the TUM Campus App:</b>
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
	assert.Equal(t, `Feedback via the TUM Campus App:

no feedback provided

Metadata:
- OS-Version: unknown
- App-Version: unknown`, txtBody)
}
