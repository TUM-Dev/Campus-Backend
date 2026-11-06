package main

import (
	"context"
	"slices"
	"strings"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"google.golang.org/grpc/codes"
)

func miscChecks() []check {
	return []check{
		{group: "movies", name: "list", fn: checkMovies},
		{group: "movies", name: "pagination", fn: checkMoviesPagination},
		{group: "student-clubs", name: "german", fn: func(ctx context.Context, e *env, t *T) { checkStudentClubs(ctx, e, t, pb.Language_German) }},
		{group: "student-clubs", name: "english", fn: func(ctx context.Context, e *env, t *T) { checkStudentClubs(ctx, e, t, pb.Language_English) }},
		{group: "student-clubs", name: "language-parity", fn: checkStudentClubParity},
		{group: "update-note", name: "unknown-version", fn: checkUpdateNote},
	}
}

func (e *env) allMovies() ([]*pb.Movie, error) {
	return memo(e, "movies", func(ctx context.Context) ([]*pb.Movie, error) {
		r, err := e.c.ListMovies(ctx, &pb.ListMoviesRequest{LastId: -1})
		return r.GetMovies(), err
	})
}

func checkMovies(ctx context.Context, e *env, t *T) {
	var movies []*pb.Movie
	if err := t.timed("ListMovies(-1)", func() (err error) { movies, err = e.allMovies(); return }); err != nil {
		t.Failf("ListMovies(last_id=-1): %v", err)
		return
	}
	if len(movies) == 0 {
		t.Failf("no movies")
		return
	}
	size := sizeOf(&pb.ListMoviesReply{Movies: movies})
	t.Infof("%d movies, %d KiB", len(movies), size>>10)
	if size > 512<<10 {
		t.Warnf("ListMovies(last_id=-1) returns %d KiB; apps download the whole archive unless they pass oldest_date_at", size>>10)
	}

	now := time.Now()
	var upcoming []*pb.Movie
	ids := map[int64]bool{}
	var covers []string
	for i, m := range movies {
		if ids[m.GetMovieId()] {
			t.Failf("duplicate movie id %d", m.GetMovieId())
		}
		ids[m.GetMovieId()] = true
		if i > 0 && m.GetDate().AsTime().Before(movies[i-1].GetDate().AsTime()) {
			t.Warnf("movies are not sorted by date at index %d", i)
		}
		if m.GetTitle() == "" {
			t.Warnf("movie %d has no title", m.GetMovieId())
		}
		if m.GetDate().AsTime().After(now) {
			upcoming = append(upcoming, m)
			covers = append(covers, m.GetCoverUrl())
		}
	}
	if len(upcoming) == 0 {
		last := movies[len(movies)-1].GetDate().AsTime()
		t.Warnf("no upcoming movies; last screening was %s (semester break, or tu-film import broken?)", last.Format(time.DateOnly))
		return
	}
	noCover, noLocation := 0, 0
	for _, m := range upcoming {
		if m.GetCoverUrl() == "" {
			noCover++
		}
		if m.GetLocation() == "" {
			noLocation++
		}
	}
	t.Infof("%d upcoming movies", len(upcoming))
	if noCover > 0 {
		t.Warnf("%d/%d upcoming movies have no cover_url (cover id set but file not downloaded?)", noCover, len(upcoming))
	}
	if noLocation > 0 {
		t.Warnf("%d/%d upcoming movies have no location", noLocation, len(upcoming))
	}
	e.checkURLs(ctx, t, "upcoming movie covers", covers, 20, Fail)
}

func checkMoviesPagination(ctx context.Context, e *env, t *T) {
	movies, err := e.allMovies()
	if err != nil || len(movies) == 0 {
		t.Skipf("no movies: %v", err)
		return
	}
	maxID := slices.MaxFunc(movies, func(a, b *pb.Movie) int { return int(a.GetMovieId() - b.GetMovieId()) }).GetMovieId()
	r, err := e.c.ListMovies(ctx, &pb.ListMoviesRequest{LastId: int32(maxID)})
	if err != nil {
		t.Failf("ListMovies(last_id=%d): %v", maxID, err)
	} else if n := len(r.GetMovies()); n != 0 {
		t.Failf("ListMovies(last_id=%d) returned %d movies, want 0", maxID, n)
	}

	r0, err := e.c.ListMovies(ctx, &pb.ListMoviesRequest{LastId: 0})
	if err != nil {
		t.Failf("ListMovies(last_id=0): %v", err)
	} else if len(r0.GetMovies()) != len(movies) {
		t.Infof("last_id=0 returns %d movies, last_id=-1 returns %d", len(r0.GetMovies()), len(movies))
	}
}

func (e *env) studentClubs(lang pb.Language) ([]*pb.StudentClubCollection, error) {
	return memo(e, "clubs-"+lang.String(), func(ctx context.Context) ([]*pb.StudentClubCollection, error) {
		r, err := e.c.ListStudentGroup(ctx, &pb.ListStudentGroupRequest{Language: &lang})
		return r.GetStudentClubs(), err
	})
}

type clubStats struct{ collections, clubs, links, covers, descriptions int }

func statsOf(cols []*pb.StudentClubCollection) clubStats {
	s := clubStats{collections: len(cols)}
	for _, c := range cols {
		for _, club := range c.GetClubs() {
			s.clubs++
			if club.GetLinkUrl() != "" {
				s.links++
			}
			if club.GetCoverUrl() != "" {
				s.covers++
			}
			if club.GetDescription() != "" {
				s.descriptions++
			}
		}
	}
	return s
}

func checkStudentClubs(ctx context.Context, e *env, t *T, lang pb.Language) {
	var cols []*pb.StudentClubCollection
	if err := t.timed("ListStudentClub", func() (err error) { cols, err = e.studentClubs(lang); return }); err != nil {
		t.Failf("ListStudentClub(%s): %v", lang, err)
		return
	}
	s := statsOf(cols)
	if s.clubs == 0 {
		t.Failf("ListStudentClub(%s) returned no clubs", lang)
		return
	}
	t.Infof("%d collections, %d clubs: %.0f%% with link, %.0f%% with cover, %.0f%% with description",
		s.collections, s.clubs, pct(s.links, s.clubs), pct(s.covers, s.clubs), pct(s.descriptions, s.clubs))

	var covers, links []string
	names := map[string]int{}
	for _, c := range cols {
		if strings.TrimSpace(c.GetTitle()) == "" {
			t.Warnf("collection %d has no title", c.GetUnstableCollectionId())
		}
		if len(c.GetClubs()) == 0 {
			t.Warnf("collection %q is empty", c.GetTitle())
		}
		for _, club := range c.GetClubs() {
			names[club.GetName()]++
			covers = append(covers, club.GetCoverUrl())
			links = append(links, club.GetLinkUrl())
			if l := club.GetLinkUrl(); l != "" && !strings.HasPrefix(l, "http://") && !strings.HasPrefix(l, "https://") {
				t.Warnf("club %q has a non-http link %q", club.GetName(), l)
			}
		}
	}
	for name, n := range names {
		if n > 1 {
			t.Warnf("club %q appears %d times", name, n)
		}
	}
	if pct(s.links, s.clubs) < 50 {
		t.Warnf("only %d/%d clubs have a link", s.links, s.clubs)
	}
	e.checkURLs(ctx, t, "club covers", covers, 30, Fail)
	e.checkURLs(ctx, t, "club links (external)", links, 15, Warn)
}

// checkStudentClubParity compares the German and English scrape, which come from different pages.
func checkStudentClubParity(_ context.Context, e *env, t *T) {
	de, errDE := e.studentClubs(pb.Language_German)
	en, errEN := e.studentClubs(pb.Language_English)
	if errDE != nil || errEN != nil {
		t.Skipf("need both languages: de=%v en=%v", errDE, errEN)
		return
	}
	sd, se := statsOf(de), statsOf(en)
	t.Infof("de: %d clubs/%d links/%d covers, en: %d clubs/%d links/%d covers", sd.clubs, sd.links, sd.covers, se.clubs, se.links, se.covers)
	ratio := func(a, b int) float64 {
		if b == 0 {
			return 1
		}
		return float64(a) / float64(b)
	}
	if r := ratio(se.clubs, sd.clubs); r < 0.8 || r > 1.25 {
		t.Warnf("English has %d clubs, German %d", se.clubs, sd.clubs)
	}
	if ratio(se.links, sd.links) < 0.8 {
		t.Warnf("English clubs have far fewer links (%d) than German ones (%d)", se.links, sd.links)
	}
	if ratio(se.covers, sd.covers) < 0.8 {
		t.Warnf("English clubs have far fewer covers (%d) than German ones (%d)", se.covers, sd.covers)
	}
	if sd.collections != se.collections {
		t.Warnf("%d German collections but %d English ones", sd.collections, se.collections)
	}
}

func checkUpdateNote(ctx context.Context, e *env, t *T) {
	_, err := e.c.GetUpdateNote(ctx, &pb.GetUpdateNoteRequest{Version: 999999999})
	expectCode(t, "GetUpdateNote(999999999)", err, codes.NotFound)
}
