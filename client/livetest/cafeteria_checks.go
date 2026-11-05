package main

import (
	"context"
	"fmt"
	"math"
	"time"
	_ "time/tzdata" // Europe/Berlin must resolve on any machine

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// canteens checked for dishes and ratings, in the eat-api id scheme.
var probeCanteens = []string{"mensa-garching", "mensa-arcisstr", "mensa-leopoldstr"}

const unknownCanteen = "LIVETEST_DOES_NOT_EXIST"

func cafeteriaChecks() []check {
	return []check{
		{group: "cafeteria", name: "canteens", fn: checkCanteens},
		{group: "cafeteria", name: "rating-tags", fn: checkRatingTags},
		{group: "cafeteria", name: "dishes-this-week", fn: checkDishesThisWeek},
		{group: "cafeteria", name: "dishes-validation", fn: checkDishesValidation},
		{group: "cafeteria", name: "canteen-ratings", fn: checkCanteenRatings},
		{group: "cafeteria", name: "dish-ratings", fn: checkDishRatings},
		{group: "cafeteria", name: "head-count", fn: checkHeadCount},
	}
}

func berlinNow() time.Time {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return time.Now()
	}
	return time.Now().In(loc)
}

func checkCanteens(ctx context.Context, e *env, t *T) {
	var r *pb.ListCanteensReply
	if err := t.timed("ListCanteens", func() (err error) {
		r, err = e.c.ListCanteens(ctx, &pb.ListCanteensRequest{})
		return
	}); err != nil {
		t.Failf("ListCanteens: %v", err)
		return
	}
	if len(r.GetCanteen()) == 0 {
		t.Failf("ListCanteens returned no canteens")
	}
	for _, c := range r.GetCanteen() {
		// everything the eat-api covers is in Bavaria
		if c.GetLatitude() < 47 || c.GetLatitude() > 50.6 || c.GetLongitude() < 8.9 || c.GetLongitude() > 13.9 {
			t.Warnf("canteen %s has coordinates outside Bavaria: %f,%f", c.GetId(), c.GetLatitude(), c.GetLongitude())
		}
		if c.GetAddress() == "" {
			t.Warnf("canteen %s has no address", c.GetId())
		}
	}
	t.Infof("%d canteens", len(r.GetCanteen()))
}

func checkRatingTags(ctx context.Context, e *env, t *T) {
	calls := []struct {
		name string
		fn   func() ([]*pb.TagsOverview, error)
	}{
		{"ListAvailableDishTags", func() ([]*pb.TagsOverview, error) {
			r, err := e.c.ListAvailableDishTags(ctx, &pb.ListAvailableDishTagsRequest{})
			return r.GetRatingTags(), err
		}},
		{"ListNameTags", func() ([]*pb.TagsOverview, error) {
			r, err := e.c.ListNameTags(ctx, &pb.ListNameTagsRequest{})
			return r.GetRatingTags(), err
		}},
		{"ListAvailableCanteenTags", func() ([]*pb.TagsOverview, error) {
			r, err := e.c.ListAvailableCanteenTags(ctx, &pb.ListAvailableCanteenTagsRequest{})
			return r.GetRatingTags(), err
		}},
	}
	for _, c := range calls {
		var tags []*pb.TagsOverview
		if err := t.timed(c.name, func() (err error) { tags, err = c.fn(); return }); err != nil {
			t.Failf("%s: %v", c.name, err)
			continue
		}
		if len(tags) == 0 {
			t.Failf("%s returned no tags", c.name)
			continue
		}
		ids := map[int32]bool{}
		for _, tag := range tags {
			if ids[tag.GetTagId()] {
				t.Warnf("%s: duplicate tag id %d", c.name, tag.GetTagId())
			}
			ids[tag.GetTagId()] = true
			if tag.GetDe() == "" || tag.GetEn() == "" {
				t.Warnf("%s: tag %d is missing a translation (de=%q en=%q)", c.name, tag.GetTagId(), tag.GetDe(), tag.GetEn())
			}
		}
		t.Infof("%s: %d tags", c.name, len(tags))
	}
}

// dishesThisWeek maps canteen -> weekday (0 = Monday) -> dishes.
func (e *env) dishesThisWeek() (map[string]map[int32][]string, error) {
	return memo(e, "dishes", func(ctx context.Context) (map[string]map[int32][]string, error) {
		year, week := berlinNow().ISOWeek()
		out := map[string]map[int32][]string{}
		var lastErr error
		for _, c := range probeCanteens {
			out[c] = map[int32][]string{}
			for day := int32(0); day <= 4; day++ {
				r, err := e.c.ListDishes(ctx, &pb.ListDishesRequest{CanteenId: c, Year: int32(year), Week: int32(week), Day: day})
				if err != nil {
					lastErr = fmt.Errorf("ListDishes(%s, %d-W%02d, day %d): %w", c, year, week, day, err)
					continue
				}
				out[c][day] = r.GetDish()
			}
		}
		for _, days := range out {
			if len(days) > 0 {
				return out, nil
			}
		}
		return nil, lastErr // every single call failed
	})
}

func checkDishesThisWeek(ctx context.Context, e *env, t *T) {
	var dishes map[string]map[int32][]string
	if err := t.timed("ListDishes x15", func() (err error) { dishes, err = e.dishesThisWeek(); return }); err != nil {
		t.Failf("%v", err)
		return
	}
	year, week := berlinNow().ISOWeek()
	total := 0
	for _, c := range probeCanteens {
		n, days := 0, 0
		for _, d := range dishes[c] {
			n += len(d)
			if len(d) > 0 {
				days++
			}
		}
		total += n
		if n == 0 {
			t.Warnf("%s: no dishes for %d-W%02d", c, year, week)
		} else {
			t.Infof("%s: %d dishes on %d weekdays", c, n, days)
		}
	}
	if total == 0 {
		t.Failf("no dishes at all for %d-W%02d in %v; is the dish import running?", year, week, probeCanteens)
	}
}

func checkDishesValidation(ctx context.Context, e *env, t *T) {
	year, week := berlinNow().ISOWeek()
	bad := []struct {
		what string
		req  *pb.ListDishesRequest
	}{
		{"year 2000", &pb.ListDishesRequest{CanteenId: probeCanteens[0], Year: 2000, Week: int32(week), Day: 0}},
		{"week 0", &pb.ListDishesRequest{CanteenId: probeCanteens[0], Year: int32(year), Week: 0, Day: 0}},
		{"day 7", &pb.ListDishesRequest{CanteenId: probeCanteens[0], Year: int32(year), Week: int32(week), Day: 7}},
	}
	for _, b := range bad {
		_, err := e.c.ListDishes(ctx, b.req)
		expectCode(t, "ListDishes("+b.what+")", err, codes.InvalidArgument)
	}

	// ISO years have 53 weeks when Dec 28 falls into week 53.
	for _, y := range []int{year, year + 1} {
		if _, w := time.Date(y, 12, 28, 0, 0, 0, 0, time.UTC).ISOWeek(); w == 53 {
			_, err := e.c.ListDishes(ctx, &pb.ListDishesRequest{CanteenId: probeCanteens[0], Year: int32(y), Week: 53, Day: 0})
			if err != nil {
				t.Warnf("ListDishes rejects week 53, but %d has an ISO week 53 (starts %s): %v", y, isoWeekStart(y, 53).Format(time.DateOnly), err)
			}
		}
	}

	_, err := e.c.ListDishes(ctx, &pb.ListDishesRequest{CanteenId: unknownCanteen, Year: int32(year), Week: int32(week), Day: 0})
	expectCode(t, "ListDishes(unknown canteen)", err, codes.OK, codes.NotFound, codes.InvalidArgument)
}

func isoWeekStart(year, week int) time.Time {
	jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, time.UTC)
	monday := jan4.AddDate(0, 0, -((int(jan4.Weekday()) + 6) % 7))
	return monday.AddDate(0, 0, (week-1)*7)
}

func checkRatingStats(t *T, what string, ratings []*pb.SingleRatingReply, avg float64, minV, maxV int32) {
	for _, r := range ratings {
		if r.GetPoints() < 1 || r.GetPoints() > 5 {
			t.Warnf("%s: rating with %d points (documented range 1-5)", what, r.GetPoints())
		}
		if len(r.GetComment()) > 256 {
			t.Warnf("%s: comment longer than 256 chars", what)
		}
	}
	if len(ratings) > 0 && (avg < float64(minV)-1e-9 || avg > float64(maxV)+1e-9 || math.IsNaN(avg)) {
		t.Warnf("%s: avg %.2f outside [min %d, max %d]", what, avg, minV, maxV)
	}
}

func checkCanteenRatings(ctx context.Context, e *env, t *T) {
	for _, c := range probeCanteens[:2] {
		id := eatAPIEnum(c)
		var r *pb.ListCanteenRatingsReply
		if err := t.timed("ListCanteenRatings", func() (err error) {
			r, err = e.c.ListCanteenRatings(ctx, &pb.ListCanteenRatingsRequest{CanteenId: id, Limit: 5})
			return
		}); err != nil {
			t.Failf("ListCanteenRatings(%s): %v", id, err)
			continue
		}
		if len(r.GetRating()) > 5 {
			t.Failf("ListCanteenRatings(%s, limit 5) returned %d ratings", id, len(r.GetRating()))
		}
		checkRatingStats(t, "ListCanteenRatings("+id+")", r.GetRating(), r.GetAvg(), r.GetMin(), r.GetMax())
		t.Infof("%s: %d ratings returned, avg %.2f", id, len(r.GetRating()), r.GetAvg())
	}
	_, err := e.c.ListCanteenRatings(ctx, &pb.ListCanteenRatingsRequest{CanteenId: unknownCanteen, Limit: 1})
	expectCode(t, "ListCanteenRatings(unknown canteen)", err, codes.InvalidArgument, codes.NotFound)
}

func checkDishRatings(ctx context.Context, e *env, t *T) {
	dishes, err := e.dishesThisWeek()
	if err != nil {
		t.Skipf("no dishes to look up ratings for: %v", err)
		return
	}
	checked := 0
	for _, c := range probeCanteens {
		for day := int32(0); day <= 4 && checked < 3; day++ {
			if len(dishes[c][day]) == 0 {
				continue
			}
			dish := dishes[c][day][0]
			var r *pb.GetDishRatingsReply
			err := t.timed("GetDishRatings", func() (err error) {
				r, err = e.c.GetDishRatings(ctx, &pb.GetDishRatingsRequest{CanteenId: eatAPIEnum(c), Dish: dish, Limit: 3})
				return
			})
			checked++
			if status.Code(err) == codes.NotFound {
				// Unrated dishes are answered with NotFound ("This dish has not yet been rated.").
				t.Infof("GetDishRatings(%s, %q): not rated yet", c, dish)
				continue
			}
			if err != nil {
				t.Failf("GetDishRatings(%s, %q): %v", c, dish, err)
				continue
			}
			checkRatingStats(t, fmt.Sprintf("GetDishRatings(%s, %q)", c, dish), r.GetRating(), r.GetAvg(), r.GetMin(), r.GetMax())
		}
	}
	if checked == 0 {
		t.Skipf("no dishes this week")
	} else {
		t.Infof("looked up ratings of %d dishes", checked)
	}
	_, err = e.c.GetDishRatings(ctx, &pb.GetDishRatingsRequest{CanteenId: eatAPIEnum(probeCanteens[0]), Dish: "livetest dish that does not exist", Limit: 1})
	expectCode(t, "GetDishRatings(unknown dish)", err, codes.InvalidArgument, codes.NotFound)
}

func checkHeadCount(ctx context.Context, e *env, t *T) {
	now := berlinNow()
	var r *pb.GetCanteenHeadCountReply
	if err := t.timed("GetCanteenHeadCount", func() (err error) {
		r, err = e.c.GetCanteenHeadCount(ctx, &pb.GetCanteenHeadCountRequest{CanteenId: probeCanteens[0]})
		return
	}); err != nil {
		t.Failf("GetCanteenHeadCount(%s): %v", probeCanteens[0], err)
		return
	}
	switch p := r.GetPercent(); {
	case p == -1:
		t.Warnf("%s: no head count data", probeCanteens[0])
	case p < 0 || p > 100:
		t.Failf("%s: percent %.1f outside [0,100]", probeCanteens[0], p)
	default:
		age := time.Since(r.GetTimestamp().AsTime())
		t.Infof("%s: %d/%d (%.0f%%), %s old", probeCanteens[0], r.GetCount(), r.GetMaxCount(), p, age.Round(time.Minute))
		weekday := now.Weekday() >= time.Monday && now.Weekday() <= time.Friday
		if weekday && now.Hour() >= 11 && now.Hour() < 14 && age > 30*time.Minute {
			t.Warnf("%s: head count is %s old during lunch time; is the head count import running?", probeCanteens[0], age.Round(time.Minute))
		}
		if age > 7*24*time.Hour {
			t.Warnf("%s: head count not updated for %d days", probeCanteens[0], int(age.Hours()/24))
		}
	}

	r, err := e.c.GetCanteenHeadCount(ctx, &pb.GetCanteenHeadCountRequest{CanteenId: unknownCanteen})
	if err == nil && r.GetPercent() != -1 {
		t.Warnf("unknown canteen: percent %.1f, documented to be -1 when no data is available", r.GetPercent())
	}
}

// eatAPIEnum turns "mensa-garching" into "MENSA_GARCHING".
func eatAPIEnum(id string) string {
	out := []byte(id)
	for i, b := range out {
		switch {
		case b == '-':
			out[i] = '_'
		case b >= 'a' && b <= 'z':
			out[i] = b - 'a' + 'A'
		}
	}
	return string(out)
}
