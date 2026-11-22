package backend

import (
	"context"
	"time"
	_ "time/tzdata"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	busynessWindowDays    = 84   // 12 weeks
	busynessHalfLifeDays  = 14.0 // 2-week exponential decay half-life so semester/break shifts adapt quickly
	lunchStartBucket      = 11 * 4
	lunchEndBucket        = 14*4 + 2 // 14:30
	liveStalenessLimitMin = 30
)

type busynessBucketAgg struct {
	Bucket         int32   `gorm:"column:bucket"`
	TypicalPercent float64 `gorm:"column:typical_percent"`
	TypicalCount   float64 `gorm:"column:typical_count"`
	SampleDays     int32   `gorm:"column:sample_days"`
}

// GetCanteenBusyness returns live canteen head count together with typical 15-minute
// weekday busyness curves and a quietest upcoming lunch slot hint.
func (s *CampusServer) GetCanteenBusyness(ctx context.Context, input *pb.GetCanteenBusynessRequest) (*pb.GetCanteenBusynessReply, error) {
	_ = s.checkDevice(ctx)
	if input.GetCanteenId() == "" {
		return nil, status.Error(codes.InvalidArgument, "canteen_id is required.")
	}
	return s.getCanteenBusynessAt(ctx, input, time.Now())
}

func (s *CampusServer) getCanteenBusynessAt(ctx context.Context, input *pb.GetCanteenBusynessRequest, now time.Time) (*pb.GetCanteenBusynessReply, error) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		loc = time.UTC
	}
	localNow := now.In(loc)
	todayWeekday := int32((int(localNow.Weekday()) + 6) % 7)
	currentBucket := int32(localNow.Hour()*4 + localNow.Minute()/15)

	weekday := todayWeekday
	if input.Weekday != nil {
		if *input.Weekday < 0 || *input.Weekday > 6 {
			return nil, status.Error(codes.InvalidArgument, "weekday must be in the range 0 (Monday) - 6 (Sunday).")
		}
		weekday = *input.Weekday
	}

	canteenID := model.CanteenEatAPIID(input.CanteenId)
	current, err := s.GetCanteenHeadCount(ctx, &pb.GetCanteenHeadCountRequest{CanteenId: canteenID})
	if err != nil {
		return nil, err
	}

	todayStr := localNow.Format("2006-01-02")
	cutoff := localNow.AddDate(0, 0, -busynessWindowDays).Format("2006-01-02")
	var aggs []busynessBucketAgg
	// Weight recent weeks exponentially higher (14-day half-life) so transitions between
	// lecture periods and semester breaks adapt quickly, exclude today's partial data,
	// and ignore closed/holiday days where daily peak count stayed at 0.
	if err := s.db.WithContext(ctx).Raw(`
		SELECT
			h.bucket AS bucket,
			SUM(h.avg_percent * EXP(-LN(2) * GREATEST(DATEDIFF(?, h.sample_date), 0) / ?)) /
				SUM(EXP(-LN(2) * GREATEST(DATEDIFF(?, h.sample_date), 0) / ?)) AS typical_percent,
			SUM(h.avg_count * EXP(-LN(2) * GREATEST(DATEDIFF(?, h.sample_date), 0) / ?)) /
				SUM(EXP(-LN(2) * GREATEST(DATEDIFF(?, h.sample_date), 0) / ?)) AS typical_count,
			COUNT(DISTINCT h.sample_date) AS sample_days
		FROM canteen_head_count_history h
		JOIN (
			SELECT sample_date
			FROM canteen_head_count_history
			WHERE canteen_id = ? AND weekday = ? AND sample_date >= ? AND sample_date < ?
			GROUP BY sample_date
			HAVING MAX(avg_count) > 0
		) open_days ON open_days.sample_date = h.sample_date
		WHERE h.canteen_id = ? AND h.weekday = ? AND h.sample_date >= ? AND h.sample_date < ?
		GROUP BY h.bucket
		ORDER BY h.bucket ASC
	`,
		todayStr, busynessHalfLifeDays,
		todayStr, busynessHalfLifeDays,
		todayStr, busynessHalfLifeDays,
		todayStr, busynessHalfLifeDays,
		canteenID, weekday, cutoff, todayStr,
		canteenID, weekday, cutoff, todayStr,
	).Scan(&aggs).Error; err != nil {
		log.WithError(err).WithField("canteen_id", canteenID).Error("while querying canteen busyness history")
		return nil, status.Error(codes.Internal, "Could not query canteen busyness history.")
	}

	slots := make([]*pb.BusynessSlot, 0, len(aggs))
	var currentTypical *busynessBucketAgg
	var bestUpcoming *pb.BusynessSlot

	for i := range aggs {
		a := aggs[i]
		slot := &pb.BusynessSlot{
			Hour:           a.Bucket / 4,
			Minute:         (a.Bucket % 4) * 15,
			TypicalPercent: a.TypicalPercent,
			TypicalCount:   a.TypicalCount,
			SampleDays:     a.SampleDays,
		}
		slots = append(slots, slot)
		if weekday == todayWeekday && a.Bucket == currentBucket {
			currentTypical = &aggs[i]
		}
		// Suggest the quietest upcoming lunch slot (11:00–14:30 Berlin).
		minBucket := int32(lunchStartBucket)
		if weekday == todayWeekday && currentBucket+1 > minBucket {
			minBucket = currentBucket + 1
		}
		if a.Bucket >= minBucket && a.Bucket <= lunchEndBucket {
			if bestUpcoming == nil || slot.TypicalPercent < bestUpcoming.TypicalPercent {
				bestUpcoming = slot
			}
		}
	}

	reply := &pb.GetCanteenBusynessReply{
		Current:          current,
		Weekday:          weekday,
		TypicalSlots:     slots,
		BestUpcomingSlot: bestUpcoming,
	}
	if current.Percent >= 0 && currentTypical != nil && current.Timestamp != nil &&
		now.Sub(current.Timestamp.AsTime()) <= liveStalenessLimitMin*time.Minute {
		reply.HasTypicalComparison = true
		reply.VsTypicalPercent = float64(current.Percent) - currentTypical.TypicalPercent
	}
	return reply, nil
}
