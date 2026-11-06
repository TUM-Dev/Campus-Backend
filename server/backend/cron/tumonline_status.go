package cron

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
	_ "time/tzdata" // The runtime image does not include the zoneinfo database.

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/oauth2/google"
)

const tumonlineStatusURL = "https://status.it.tum.de/"

type tumonlineStatus struct {
	Maintenance bool
	Degraded    bool
	MessageDE   string
	MessageEN   string
}

var tumonlineService = regexp.MustCompile(`(?i)\b(tumonline|campus\.tum\.de)\b`)
var maintenanceNotice = regexp.MustCompile(`(?i)\b(wartung\w*|maintenance)\b`)

// Read the authoritative current section rather than estimating recovery from
// the advertised end time. Future and past notices must not enable a banner.
func parseTUMOnlineStatus(r io.Reader) (tumonlineStatus, error) {
	return parseTUMOnlineStatusAt(r, time.Now())
}

func parseTUMOnlineStatusAt(r io.Reader, now time.Time) (tumonlineStatus, error) {
	var result tumonlineStatus
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return result, err
	}
	var current *goquery.Selection
	var future, past bool
	doc.Find("h2").Each(func(_ int, heading *goquery.Selection) {
		switch strings.Join(strings.Fields(heading.Text()), " ") {
		case "Aktuelle Meldungen", "Current notifications":
			current = heading.NextUntil("h2")
		case "Angekündigte Meldungen", "Future notifications":
			future = true
		case "Beendete Meldungen ↑", "Past notifications ↑":
			past = true
		}
	})
	// A login/error page or changed markup must never clear an active warning.
	if current == nil || !future || !past || current.Find(".table-headers").Length() != 1 {
		return result, fmt.Errorf("unrecognized TUM status page structure")
	}
	var parseErr error
	var maintenanceEnds, incidentEnds []time.Time
	berlin, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return result, err
	}
	current.Find("a.table-entry").Each(func(_ int, entry *goquery.Selection) {
		service := entry.Find(".entry-service .mobile-bold")
		title := entry.Find(".entry-problemtype")
		if service.Length() != 1 || title.Length() != 1 || strings.TrimSpace(service.Text()) == "" || strings.TrimSpace(title.Text()) == "" {
			parseErr = fmt.Errorf("unrecognized TUM status entry")
			return
		}
		if !tumonlineService.MatchString(service.Text()) {
			return
		}
		end := parseStatusEnd(entry.Find(".entry-end.entry-date").Text(), berlin)
		if maintenanceNotice.MatchString(title.Text()) {
			result.Maintenance = true
			maintenanceEnds = append(maintenanceEnds, end)
		} else {
			// The app has one flag for all unplanned service incidents, including outages.
			result.Degraded = true
			incidentEnds = append(incidentEnds, end)
		}
	})
	if result.Maintenance {
		result.Degraded = false
		result.setMessagesWithEstimatedEnd(maintenanceEnds, now)
	} else if result.Degraded {
		result.setMessagesWithEstimatedEnd(incidentEnds, now)
	}
	return result, parseErr
}

func parseStatusEnd(text string, location *time.Location) time.Time {
	text = strings.Join(strings.Fields(text), " ")
	for _, layout := range []string{"02.01.2006 15:04", "1/2/2006 3:04 PM"} {
		if end, err := time.ParseInLocation(layout, text, location); err == nil {
			return end
		}
	}
	return time.Time{} // An unknown end must not prevent publishing the warning.
}

func (s *tumonlineStatus) setMessagesWithEstimatedEnd(ends []time.Time, now time.Time) {
	var latest time.Time
	for _, end := range ends {
		// A missing or elapsed estimate on any relevant notice makes a combined
		// recovery estimate unreliable. Leave both messages empty so the app
		// displays its translated default for the active Boolean flag.
		if !end.After(now) {
			return
		}
		if end.After(latest) {
			latest = end
		}
	}
	if latest.IsZero() {
		return
	}
	if s.Maintenance {
		s.MessageDE = "TUMonline wird derzeit gewartet."
		s.MessageEN = "TUMonline is currently undergoing maintenance."
	} else {
		s.MessageDE = "Die Dienste von TUMonline sind derzeit beeinträchtigt."
		s.MessageEN = "TUMonline services are currently degraded."
	}
	s.MessageDE += " Voraussichtliches Ende: " + latest.Format("02.01.2006, 15:04 MST") + "."
	s.MessageEN += " Estimated end: " + latest.Format("2 Jan 2006, 15:04 MST") + "."
}

func fetchTUMOnlineStatus(ctx context.Context, client *http.Client) (tumonlineStatus, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tumonlineStatusURL, nil)
	if err != nil {
		return tumonlineStatus{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return tumonlineStatus{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return tumonlineStatus{}, fmt.Errorf("TUM status returned HTTP %d", resp.StatusCode)
	}
	// Do not parse a truncated document as a valid all-clear.
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20+1))
	if err != nil {
		return tumonlineStatus{}, err
	}
	if len(body) > 2<<20 {
		return tumonlineStatus{}, fmt.Errorf("TUM status page exceeds size limit")
	}
	return parseTUMOnlineStatus(strings.NewReader(string(body)))
}

func (c *CronService) tumonlineStatusCron() error {
	project := os.Getenv("FIREBASE_PROJECT_ID")
	if project == "" || os.Getenv("TUMONLINE_STATUS_SYNC_ENABLED") != "true" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	state, err := fetchTUMOnlineStatus(ctx, c.client)
	if err != nil {
		return err
	}
	client, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/firebase.remoteconfig")
	if err != nil {
		return fmt.Errorf("remote config credentials: %w", err)
	}
	client.Timeout = 30 * time.Second
	endpoint := "https://firebaseremoteconfig.googleapis.com/v1/projects/" + url.PathEscape(project) + "/remoteConfig"
	return publishTUMOnlineStatus(ctx, client, endpoint, state)
}
