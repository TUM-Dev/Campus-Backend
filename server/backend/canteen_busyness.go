package backend

import (
	"context"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	busynessWindowDays    = 56 // 8 weeks
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

	cutoff := localNow.AddDate(0, 0, -busynessWindowDays).Format("2006-01-02")
	var aggs []busynessBucketAgg
	if err := s.db.WithContext(ctx).
		Table("canteen_head_count_history").
		Select("bucket, AVG(avg_percent) AS typical_percent, AVG(avg_count) AS typical_count, COUNT(DISTINCT sample_date) AS sample_days").
		Where("canteen_id = ? AND weekday = ? AND sample_date >= ?", canteenID, weekday, cutoff).
		Group("bucket").
		Order("bucket ASC").
		Scan(&aggs).Error; err != nil {
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
