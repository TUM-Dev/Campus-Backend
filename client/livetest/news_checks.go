package main

import (
	"context"
	"slices"
	"strconv"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newsChecks() []check {
	return []check{
		{group: "news", name: "sources", fn: checkNewsSources},
		{group: "news", name: "list", fn: checkNewsList},
		{group: "news", name: "filters", fn: checkNewsFilters},
		{group: "news", name: "alerts", fn: checkNewsAlerts},
	}
}

func (e *env) newsSources() ([]*pb.NewsSource, error) {
	return memo(e, "newsSources", func(ctx context.Context) ([]*pb.NewsSource, error) {
		r, err := e.c.ListNewsSources(ctx, &pb.ListNewsSourcesRequest{})
		return r.GetSources(), err
	})
}

func (e *env) allNews() ([]*pb.News, error) {
	return memo(e, "allNews", func(ctx context.Context) ([]*pb.News, error) {
		r, err := e.c.ListNews(ctx, &pb.ListNewsRequest{})
		return r.GetNews(), err
	})
}

func checkNewsSources(ctx context.Context, e *env, t *T) {
	var sources []*pb.NewsSource
	if err := t.timed("ListNewsSources", func() (err error) { sources, err = e.newsSources(); return }); err != nil {
		t.Failf("ListNewsSources: %v", err)
		return
	}
	if len(sources) == 0 {
		t.Failf("no news sources")
		return
	}
	var icons []string
	for _, s := range sources {
		if s.GetSource() == "" || s.GetTitle() == "" {
			t.Warnf("news source with empty id or title: %v", s)
		}
		if s.GetIconUrl() == "" {
			t.Warnf("news source %s (%s) has no icon", s.GetSource(), s.GetTitle())
		}
		icons = append(icons, s.GetIconUrl())
	}
	t.Infof("%d news sources", len(sources))
	e.checkURLs(ctx, t, "source icons", icons, 30, Fail)
}

func checkNewsList(ctx context.Context, e *env, t *T) {
	var news []*pb.News
	if err := t.timed("ListNews", func() (err error) { news, err = e.allNews(); return }); err != nil {
		t.Failf("ListNews: %v", err)
		return
	}
	if len(news) == 0 {
		t.Failf("ListNews returned no news")
		return
	}
	sources, _ := e.newsSources()
	known := map[string]bool{}
	for _, s := range sources {
		known[s.GetSource()] = true
	}

	now := time.Now()
	ids := map[int64]bool{}
	newestPerSource := map[string]time.Time{}
	var images []string
	var dupes, unknownSource, future, emptyTitle, noLink int
	for _, n := range news {
		if ids[n.GetId()] {
			dupes++
		}
		ids[n.GetId()] = true
		if len(known) > 0 && !known[n.GetSourceId()] {
			unknownSource++
		}
		if n.GetTitle() == "" {
			emptyTitle++
		}
		if n.GetLink() == "" {
			noLink++
		}
		d := n.GetDate().AsTime()
		if d.After(now.Add(90 * 24 * time.Hour)) {
			future++
		}
		if d.Before(now) && d.After(newestPerSource[n.GetSourceId()]) {
			newestPerSource[n.GetSourceId()] = d
		}
		images = append(images, n.GetImageUrl())
	}
	size := sizeOf(&pb.ListNewsReply{News: news})
	t.Infof("%d news items, %d KiB", len(news), size>>10)
	if dupes > 0 {
		t.Failf("%d duplicate news ids", dupes)
	}
	if unknownSource > 0 {
		t.Warnf("%d news items reference a source not returned by ListNewsSources", unknownSource)
	}
	if future > 0 {
		t.Warnf("%d news items are dated more than 90 days in the future", future)
	}
	if emptyTitle > 0 {
		t.Warnf("%d news items have an empty title", emptyTitle)
	}
	if noLink > 0 {
		t.Warnf("%d news items have no link", noLink)
	}
	if size > 2<<20 {
		t.Warnf("ListNews without filters returns %d KiB; the full history is sent every time", size>>10)
	}

	var newest time.Time
	for src, d := range newestPerSource {
		newest = maxTime(newest, d)
		if age := now.Sub(d); age > 60*24*time.Hour {
			t.Warnf("news source %s: newest item is %d days old (scraper broken or source dead?)", src, int(age.Hours()/24))
		}
	}
	if age := now.Sub(newest); age > 7*24*time.Hour {
		t.Failf("newest news item is %d days old; news import seems stuck", int(age.Hours()/24))
	}

	// thumbnails are served by us, so broken ones are our bug
	e.checkURLs(ctx, t, "news images", lastN(images, 40), 40, Fail)
	var links []string
	for _, n := range lastN(news, 10) {
		links = append(links, n.GetLink())
	}
	e.checkURLs(ctx, t, "news links (external)", links, 10, Warn)
}

func checkNewsFilters(ctx context.Context, e *env, t *T) {
	all, err := e.allNews()
	if err != nil || len(all) == 0 {
		t.Skipf("no news to filter: %v", err)
		return
	}
	maxID := slices.MaxFunc(all, func(a, b *pb.News) int { return int(a.GetId() - b.GetId()) }).GetId()

	// last_news_id = max id -> nothing newer
	var r *pb.ListNewsReply
	if err := t.timed("ListNews(last_news_id)", func() (err error) {
		r, err = e.c.ListNews(ctx, &pb.ListNewsRequest{LastNewsId: int32(maxID)})
		return
	}); err != nil {
		t.Failf("ListNews(last_news_id=%d): %v", maxID, err)
	} else if len(r.GetNews()) != 0 {
		t.Failf("ListNews(last_news_id=%d) returned %d items, want 0", maxID, len(r.GetNews()))
	}

	// last_news_id in the middle -> only greater ids
	mid := all[len(all)/2].GetId()
	if r, err := e.c.ListNews(ctx, &pb.ListNewsRequest{LastNewsId: int32(mid)}); err != nil {
		t.Failf("ListNews(last_news_id=%d): %v", mid, err)
	} else {
		for _, n := range r.GetNews() {
			if n.GetId() <= mid {
				t.Failf("ListNews(last_news_id=%d) returned id %d", mid, n.GetId())
				break
			}
		}
	}

	// news_source filter
	src := all[len(all)-1].GetSourceId()
	srcID, _ := strconv.Atoi(src)
	if r, err := e.c.ListNews(ctx, &pb.ListNewsRequest{NewsSource: int32(srcID)}); err != nil {
		t.Failf("ListNews(news_source=%s): %v", src, err)
	} else if len(r.GetNews()) == 0 {
		t.Failf("ListNews(news_source=%s) returned nothing although that source has news", src)
	} else {
		for _, n := range r.GetNews() {
			if n.GetSourceId() != src {
				t.Failf("ListNews(news_source=%s) returned an item of source %s", src, n.GetSourceId())
				break
			}
		}
	}

	// oldest_date_at filter
	cutoff := time.Now().Add(-30 * 24 * time.Hour).Truncate(time.Hour)
	if r, err := e.c.ListNews(ctx, &pb.ListNewsRequest{OldestDateAt: timestamppb.New(cutoff)}); err != nil {
		t.Failf("ListNews(oldest_date_at): %v", err)
	} else {
		for _, n := range r.GetNews() {
			if !n.GetDate().AsTime().After(cutoff) {
				t.Failf("ListNews(oldest_date_at=%s) returned an item dated %s", cutoff.Format(time.DateOnly), n.GetDate().AsTime().Format(time.DateOnly))
				break
			}
		}
		t.Infof("%d news items in the last 30 days", len(r.GetNews()))
	}
}

func checkNewsAlerts(ctx context.Context, e *env, t *T) {
	var r *pb.ListNewsAlertsReply
	if err := t.timed("ListNewsAlerts", func() (err error) {
		r, err = e.c.ListNewsAlerts(ctx, &pb.ListNewsAlertsRequest{})
		return
	}); err != nil {
		t.Failf("ListNewsAlerts: %v", err)
		return
	}
	now := time.Now()
	var images []string
	for _, a := range r.GetAlerts() {
		if a.GetFrom().AsTime().After(a.GetTo().AsTime()) {
			t.Warnf("alert %s: from %s is after to %s", a.GetLink(), a.GetFrom().AsTime(), a.GetTo().AsTime())
		}
		if a.GetTo().AsTime().Before(now.Add(-time.Minute)) {
			t.Failf("expired alert returned (to=%s)", a.GetTo().AsTime())
		}
		images = append(images, a.GetImageUrl())
	}
	t.Infof("%d active news alerts", len(r.GetAlerts()))
	e.checkURLs(ctx, t, "alert images", images, 10, Fail)
}

func lastN[S any](s []S, n int) []S {
	return s[max(0, len(s)-n):]
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
