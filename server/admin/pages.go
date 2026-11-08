package admin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/guregu/null"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const pageSize = 50

type page struct {
	Title, Nav, Subject, CSRF string
}

func basePage(title, nav string, session model.AdminSession) page {
	return page{Title: title, Nav: nav, Subject: session.Subject, CSRF: session.CSRFToken}
}

func (s *server) render(w http.ResponseWriter, name string, data any) {
	var body bytes.Buffer
	if err := s.templates.ExecuteTemplate(&body, name, data); err != nil {
		s.internalError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(body.Bytes())
}

type jobView struct {
	ID          int64
	Type        string
	SourceID    null.Int
	Interval    int32
	LastStarted time.Time
	Observation cron.JobStatus
}

type dashboardPage struct {
	page
	Version         string
	Uptime          time.Duration
	Goroutines      int
	HeapMiB         uint64
	DBOpen, DBInUse int
	PollTime        time.Time
	PollOK          bool
	Counts          map[string]int64
	Jobs            []jobView
}

func (s *server) dashboard(w http.ResponseWriter, r *http.Request, session model.AdminSession) {
	data := dashboardPage{
		page: basePage("System status", "status", session), Version: s.version,
		Uptime: time.Since(s.started).Round(time.Second), Goroutines: runtime.NumGoroutine(),
		Counts: make(map[string]int64),
	}
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	data.HeapMiB = memory.HeapAlloc / (1 << 20)
	db := s.db.WithContext(r.Context())
	sqlDB, err := db.DB()
	if err != nil {
		s.internalError(w, err)
		return
	}
	if err := sqlDB.PingContext(r.Context()); err != nil {
		s.internalError(w, err)
		return
	}
	stats := sqlDB.Stats()
	data.DBOpen, data.DBInUse = stats.OpenConnections, stats.InUse
	for name, table := range map[string]any{
		"News sources": &model.NewsSource{}, "News articles": &model.News{},
		"Student groups":   &model.StudentClub{},
		"Student councils": &model.StudentCouncil{},
	} {
		var count int64
		if err := db.Model(table).Count(&count).Error; err != nil {
			s.internalError(w, err)
			return
		}
		data.Counts[name] = count
	}
	var jobs []model.Crontab
	if err := db.Order("cron").Limit(501).Find(&jobs).Error; err != nil {
		s.internalError(w, err)
		return
	}
	if len(jobs) > 500 {
		http.Error(w, "Too many cron jobs for this view; inspect the scheduler database.", http.StatusConflict)
		return
	}
	observation := s.cron.Status()
	data.PollTime, data.PollOK = observation.LastPoll, observation.PollOK
	for _, job := range jobs {
		last := time.Time{}
		if job.LastRun > 0 {
			last = time.Unix(int64(job.LastRun), 0).UTC()
		}
		data.Jobs = append(data.Jobs, jobView{
			ID: job.Cron, Type: job.Type.String, SourceID: job.ID, Interval: job.Interval,
			LastStarted: last, Observation: observation.Jobs[job.Cron],
		})
	}
	s.render(w, "dashboard", data)
}

type pagination struct {
	Number, Previous, Next int
	HasNext                bool
	PreviousURL, NextURL   string
}

func requestedPage(r *http.Request) (pagination, error) {
	number := 1
	if value := r.URL.Query().Get("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 1000000 {
			return pagination{}, fmt.Errorf("invalid page number")
		}
		number = parsed
	}
	query := r.URL.Query()
	query.Set("page", strconv.Itoa(number-1))
	previous := "?" + query.Encode()
	query.Set("page", strconv.Itoa(number+1))
	return pagination{Number: number, Previous: number - 1, Next: number + 1, PreviousURL: previous, NextURL: "?" + query.Encode()}, nil
}

type feedView struct {
	model.NewsSource
	Jobs            []model.Crontab
	State, Schedule string
}

type feedsPage struct {
	page
	pagination
	Feeds []feedView
	Icons []model.NewsSource
}

func (s *server) feeds(w http.ResponseWriter, r *http.Request, session model.AdminSession) {
	paging, err := requestedPage(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db := s.db.WithContext(r.Context())
	var sources []model.NewsSource
	if err := db.Order("source").Offset((paging.Number - 1) * pageSize).Limit(pageSize + 1).Find(&sources).Error; err != nil {
		s.internalError(w, err)
		return
	}
	paging.HasNext = len(sources) > pageSize
	if paging.HasNext {
		sources = sources[:pageSize]
	}
	var icons []model.NewsSource
	if err := db.Select("source", "title", "icon").Order("source").Limit(501).Find(&icons).Error; err != nil {
		s.internalError(w, err)
		return
	}
	if len(icons) > 500 {
		http.Error(w, "Too many icons for the feed editor.", http.StatusConflict)
		return
	}
	data := feedsPage{page: basePage("News feeds", "feeds", session), pagination: paging}
	seenIcons := make(map[int64]bool)
	for _, icon := range icons {
		if !seenIcons[icon.FileID] {
			data.Icons = append(data.Icons, icon)
			seenIcons[icon.FileID] = true
		}
	}
	ids := make([]int64, 0, len(sources))
	for _, source := range sources {
		ids = append(ids, source.Source)
	}
	var jobs []model.Crontab
	if len(ids) > 0 {
		if err := db.Where("type = ? AND id IN ?", cron.NewsType, ids).Order("cron").Find(&jobs).Error; err != nil {
			s.internalError(w, err)
			return
		}
	}
	for _, source := range sources {
		view := feedView{NewsSource: source}
		for _, job := range jobs {
			if job.ID.Valid && job.ID.Int64 == source.Source {
				view.Jobs = append(view.Jobs, job)
			}
		}
		switch len(view.Jobs) {
		case 0:
			view.State, view.Schedule = "unscheduled", "No import schedule"
		case 1:
			view.State, view.Schedule = "paused", "Imports paused"
			if view.Jobs[0].Interval > 0 {
				view.State, view.Schedule = "active", "Every "+intervalLabel(view.Jobs[0].Interval)
			}
		default:
			view.State, view.Schedule = "attention", "Multiple import schedules"
		}
		data.Feeds = append(data.Feeds, view)
	}
	s.render(w, "feeds", data)
}

type feedInput struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	IconID   int64  `json:"icon_id"`
	Interval int32  `json:"interval"`
}

type feedAuditDetails struct {
	Before feedInput `json:"before"`
	After  feedInput `json:"after"`
}

func intervalLabel(seconds int32) string {
	if seconds <= 0 {
		return "Not scheduled"
	}
	for _, unit := range []struct {
		seconds int32
		name    string
	}{{86400, "day"}, {3600, "hour"}, {60, "minute"}, {1, "second"}} {
		if seconds%unit.seconds == 0 {
			count := seconds / unit.seconds
			if count == 1 {
				return "1 " + unit.name
			}
			return fmt.Sprintf("%d %ss", count, unit.name)
		}
	}
	panic("unreachable interval unit")
}

func parseFeed(r *http.Request) (feedInput, error) {
	id, err := strconv.ParseInt(r.PostForm.Get("id"), 10, 32)
	if err != nil || id < 0 {
		return feedInput{}, fmt.Errorf("invalid source ID")
	}
	icon, err := strconv.ParseInt(r.PostForm.Get("icon"), 10, 32)
	if err != nil || icon <= 0 {
		return feedInput{}, fmt.Errorf("choose an existing icon")
	}
	interval, err := strconv.ParseInt(r.PostForm.Get("interval"), 10, 32)
	if err != nil || interval < 60 || interval > 31536000 {
		return feedInput{}, fmt.Errorf("import interval must be between 60 and 31536000 seconds")
	}
	title := strings.TrimSpace(r.PostForm.Get("title"))
	if title == "" || len(title) > 255 {
		return feedInput{}, fmt.Errorf("title must contain between 1 and 255 bytes")
	}
	feedURL := strings.TrimSpace(r.PostForm.Get("url"))
	u, err := url.Parse(feedURL)
	if err != nil || len(feedURL) > 2048 || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.Fragment != "" {
		return feedInput{}, fmt.Errorf("feed URL must be an absolute HTTPS URL without credentials or fragment")
	}
	switch r.PostForm.Get("enabled") {
	case "on":
	case "":
		interval = 0
	default:
		return feedInput{}, fmt.Errorf("invalid import enabled value")
	}
	return feedInput{ID: id, Title: title, URL: feedURL, IconID: icon, Interval: int32(interval)}, nil
}

var errMultipleSchedules = errors.New("source has multiple news schedules; resolve them before editing")

func (s *server) saveFeed(w http.ResponseWriter, r *http.Request, session model.AdminSession) {
	input, err := parseFeed(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.db.WithContext(r.Context()).Transaction(func(tx *gorm.DB) error {
		var source model.NewsSource
		if input.ID != 0 {
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&source, input.ID).Error; err != nil {
				return err
			}
		}
		var icon model.File
		if err := tx.First(&icon, input.IconID).Error; err != nil {
			return err
		}
		var jobs []model.Crontab
		if input.ID != 0 {
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("type = ? AND id = ?", cron.NewsType, input.ID).Find(&jobs).Error; err != nil {
				return err
			}
		}
		if len(jobs) > 1 {
			return errMultipleSchedules
		}
		before := feedInput{ID: source.Source, Title: source.Title, URL: source.URL.String, IconID: source.FileID}
		if len(jobs) == 1 {
			before.Interval = jobs[0].Interval
		}
		action := "feed.update"
		if input.ID == 0 {
			action = "feed.create"
			source = model.NewsSource{Title: input.Title, URL: null.StringFrom(input.URL), FileID: input.IconID}
			if err := tx.Omit(clause.Associations).Create(&source).Error; err != nil {
				return err
			}
			input.ID = source.Source
		} else if err := tx.Model(&source).Updates(map[string]any{
			"title": input.Title, "url": input.URL, "icon": input.IconID,
		}).Error; err != nil {
			return err
		}
		if len(jobs) == 0 {
			job := model.Crontab{Type: null.StringFrom(cron.NewsType), ID: null.IntFrom(source.Source), Interval: input.Interval}
			if err := tx.Create(&job).Error; err != nil {
				return err
			}
			// GORM's default tag treats a zero interval as unspecified on insert.
			if input.Interval == 0 {
				if err := tx.Model(&job).UpdateColumn("interval", 0).Error; err != nil {
					return err
				}
			}
		} else if err := tx.Model(&jobs[0]).UpdateColumn("interval", input.Interval).Error; err != nil {
			return err
		}
		details, err := json.Marshal(feedAuditDetails{Before: before, After: input})
		if err != nil {
			return err
		}
		return tx.Create(&model.AdminAudit{
			Subject: session.Subject, Action: action, SourceID: source.Source, Details: string(details),
		}).Error
	})
	if errors.Is(err, errMultipleSchedules) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Source or icon not found.", http.StatusNotFound)
		return
	}
	if err != nil {
		s.internalError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/feeds", http.StatusSeeOther)
}

type directoryEntry struct {
	Name, Language, Collection, ImageURL string
	Description, LinkURL, ImageCaption   null.String
}

type directoryPage struct {
	page
	pagination
	Councils bool
	Entries  []directoryEntry
}

func localImageURL(file *model.File) string {
	if file == nil || !file.Downloaded.Valid || !file.Downloaded.Bool {
		return ""
	}
	return (&url.URL{Path: "/files/" + file.Path + file.Name}).String()
}

func (s *server) groups(w http.ResponseWriter, r *http.Request, session model.AdminSession) {
	paging, err := requestedPage(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var groups []model.StudentClub
	if err := s.db.WithContext(r.Context()).Preload("StudentClubCollection").Preload("Image").
		Order("name").Order("id").Offset((paging.Number - 1) * pageSize).Limit(pageSize + 1).Find(&groups).Error; err != nil {
		s.internalError(w, err)
		return
	}
	paging.HasNext = len(groups) > pageSize
	if paging.HasNext {
		groups = groups[:pageSize]
	}
	data := directoryPage{page: basePage("Student groups", "groups", session), pagination: paging}
	for _, group := range groups {
		data.Entries = append(data.Entries, directoryEntry{
			Name: group.Name, Language: group.Language, Collection: group.StudentClubCollection.Name,
			Description: group.Description, LinkURL: group.LinkUrl, ImageCaption: group.ImageCaption,
			ImageURL: localImageURL(group.Image),
		})
	}
	s.render(w, "groups", data)
}

func (s *server) councils(w http.ResponseWriter, r *http.Request, session model.AdminSession) {
	paging, err := requestedPage(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var councils []model.StudentCouncil
	if err := s.db.WithContext(r.Context()).Preload("StudentCouncilCollection").Preload("Image").
		Order("name").Order("id").Offset((paging.Number - 1) * pageSize).Limit(pageSize + 1).Find(&councils).Error; err != nil {
		s.internalError(w, err)
		return
	}
	paging.HasNext = len(councils) > pageSize
	if paging.HasNext {
		councils = councils[:pageSize]
	}
	data := directoryPage{page: basePage("Student councils", "groups", session), pagination: paging, Councils: true}
	for _, council := range councils {
		data.Entries = append(data.Entries, directoryEntry{
			Name: council.Name, Language: council.Language, Collection: council.StudentCouncilCollection.Name,
			Description: council.Description, LinkURL: council.LinkUrl, ImageCaption: council.ImageCaption,
			ImageURL: localImageURL(council.Image),
		})
	}
	s.render(w, "groups", data)
}

type auditPage struct {
	page
	pagination
	Entries []auditEntry
}

type auditChange struct {
	Field, Before, After string
}

type auditEntry struct {
	model.AdminAudit
	Title, ActionLabel      string
	Created, InvalidDetails bool
	Changes                 []auditChange
}

func auditView(entry model.AdminAudit) auditEntry {
	view := auditEntry{AdminAudit: entry, ActionLabel: "Source updated", Created: entry.Action == "feed.create"}
	if view.Created {
		view.ActionLabel = "Source added"
	}
	var details feedAuditDetails
	if err := json.Unmarshal([]byte(entry.Details), &details); err != nil {
		log.WithError(err).WithField("audit_id", entry.ID).Error("could not decode admin change details")
		view.InvalidDetails = true
		return view
	}
	if (entry.Action != "feed.create" && entry.Action != "feed.update") || details.After.ID != entry.SourceID || details.After.Title == "" {
		log.WithField("audit_id", entry.ID).Error("unsupported or inconsistent admin change details")
		view.InvalidDetails = true
		return view
	}
	view.Title = details.After.Title
	state := func(interval int32) string {
		if interval > 0 {
			return "Enabled"
		}
		return "Paused"
	}
	changes := []auditChange{
		{"Title", details.Before.Title, details.After.Title},
		{"Feed URL", details.Before.URL, details.After.URL},
		{"Icon", fmt.Sprintf("File #%d", details.Before.IconID), fmt.Sprintf("File #%d", details.After.IconID)},
		{"Imports", state(details.Before.Interval), state(details.After.Interval)},
		{"Import interval", intervalLabel(details.Before.Interval), intervalLabel(details.After.Interval)},
	}
	for _, change := range changes {
		if view.Created || change.Before != change.After {
			view.Changes = append(view.Changes, change)
		}
	}
	return view
}

func (s *server) audit(w http.ResponseWriter, r *http.Request, session model.AdminSession) {
	paging, err := requestedPage(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var entries []model.AdminAudit
	if err := s.db.WithContext(r.Context()).Order("id DESC").Offset((paging.Number - 1) * pageSize).Limit(pageSize + 1).Find(&entries).Error; err != nil {
		s.internalError(w, err)
		return
	}
	paging.HasNext = len(entries) > pageSize
	if paging.HasNext {
		entries = entries[:pageSize]
	}
	data := auditPage{page: basePage("Change history", "audit", session), pagination: paging}
	for _, entry := range entries {
		data.Entries = append(data.Entries, auditView(entry))
	}
	s.render(w, "audit", data)
}
