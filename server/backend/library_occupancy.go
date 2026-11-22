package backend

import (
	"context"
	"errors"
	"strings"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const (
	libraryStudyStartBucket = 8 * 4  // 08:00 Berlin
	libraryStudyEndBucket   = 21 * 4 // 21:00 Berlin
)

// ListLibraries returns all tracked TUM libraries with their latest live occupancy estimates.
func (s *CampusServer) ListLibraries(ctx context.Context, req *pb.ListLibrariesRequest) (*pb.ListLibrariesReply, error) {
	_ = s.checkDevice(ctx)
	campusFilter := strings.ToLower(strings.TrimSpace(req.GetCampus()))

	var rows []model.LibraryOccupancy
	q := s.db.WithContext(ctx).Order("name ASC")
	if campusFilter != "" {
		q = q.Where("LOWER(campus) = ?", campusFilter)
	}
	if err := q.Find(&rows).Error; err != nil {
		log.WithError(err).Error("while listing library occupancies")
		return nil, status.Error(codes.Internal, "Could not load library occupancies.")
	}

	byID := make(map[string]model.LibraryOccupancy, len(rows))
	for _, r := range rows {
		byID[r.LibraryID] = r
	}

	result := make([]*pb.LibraryOccupancy, 0, len(cron.TrackedLibraries))
	for _, def := range cron.TrackedLibraries {
		if campusFilter != "" && !strings.EqualFold(def.Campus, campusFilter) {
			continue
		}
		if live, ok := byID[def.LibraryID]; ok {
			result = append(result, &pb.LibraryOccupancy{
				LibraryId: live.LibraryID,
				Name:      live.Name,
				Campus:    live.Campus,
				Count:     live.Count,
				MaxCount:  live.MaxCount,
				Percent:   live.Percent,
				Timestamp: timestamppb.New(live.Timestamp),
			})
		} else {
			result = append(result, &pb.LibraryOccupancy{
				LibraryId: def.LibraryID,
				Name:      def.Name,
				Campus:    def.Campus,
				MaxCount:  def.MaxCount,
				Percent:   -1,
			})
		}
	}
	return &pb.ListLibrariesReply{Libraries: result}, nil
}

// GetLibraryBusyness returns live library occupancy alongside typical 15-minute weekday
// busyness curves and a quietest upcoming study slot hint.
func (s *CampusServer) GetLibraryBusyness(ctx context.Context, req *pb.GetLibraryBusynessRequest) (*pb.GetLibraryBusynessReply, error) {
	_ = s.checkDevice(ctx)
	if req.GetLibraryId() == "" {
		return nil, status.Error(codes.InvalidArgument, "library_id is required.")
	}
	return s.getLibraryBusynessAt(ctx, req, time.Now())
}

func (s *CampusServer) getLibraryBusynessAt(ctx context.Context, req *pb.GetLibraryBusynessRequest, now time.Time) (*pb.GetLibraryBusynessReply, error) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		loc = time.UTC
	}
	localNow := now.In(loc)
	todayWeekday := int32((int(localNow.Weekday()) + 6) % 7)
	currentBucket := int32(localNow.Hour()*4 + localNow.Minute()/15)

	weekday := todayWeekday
	if req.Weekday != nil {
		if *req.Weekday < 0 || *req.Weekday > 6 {
			return nil, status.Error(codes.InvalidArgument, "weekday must be in the range 0 (Monday) - 6 (Sunday).")
		}
		weekday = *req.Weekday
	}

	var def *cron.LibraryApInformation
	for i := range cron.TrackedLibraries {
		if cron.TrackedLibraries[i].LibraryID == req.LibraryId {
			def = &cron.TrackedLibraries[i]
			break
		}
	}

	var live model.LibraryOccupancy
	err = s.db.WithContext(ctx).First(&live, "library_id = ?", req.LibraryId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if def == nil {
			return nil, status.Error(codes.NotFound, "Library not found.")
		}
		live = model.LibraryOccupancy{
			LibraryID: def.LibraryID,
			Name:      def.Name,
			Campus:    def.Campus,
			MaxCount:  def.MaxCount,
			Percent:   -1,
		}
	} else if err != nil {
		log.WithError(err).WithField("library_id", req.LibraryId).Error("while querying library occupancy")
		return nil, status.Error(codes.Internal, "Could not query library occupancy.")
	}

	current := &pb.LibraryOccupancy{
		LibraryId: live.LibraryID,
		Name:      live.Name,
		Campus:    live.Campus,
		Count:     live.Count,
		MaxCount:  live.MaxCount,
		Percent:   live.Percent,
	}
	if !live.Timestamp.IsZero() {
		current.Timestamp = timestamppb.New(live.Timestamp)
	}

	cutoff := localNow.AddDate(0, 0, -busynessWindowDays).Format("2006-01-02")
	var aggs []busynessBucketAgg
	if err := s.db.WithContext(ctx).
		Table("library_occupancy_history").
		Select("bucket, AVG(avg_percent) AS typical_percent, AVG(avg_count) AS typical_count, COUNT(DISTINCT sample_date) AS sample_days").
		Where("library_id = ? AND weekday = ? AND sample_date >= ?", req.LibraryId, weekday, cutoff).
		Group("bucket").
		Order("bucket ASC").
		Scan(&aggs).Error; err != nil {
		log.WithError(err).WithField("library_id", req.LibraryId).Error("while querying library busyness history")
		return nil, status.Error(codes.Internal, "Could not query library busyness history.")
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
		minBucket := int32(libraryStudyStartBucket)
		if weekday == todayWeekday && currentBucket+1 > minBucket {
			minBucket = currentBucket + 1
		}
		if a.Bucket >= minBucket && a.Bucket <= libraryStudyEndBucket {
			if bestUpcoming == nil || slot.TypicalPercent < bestUpcoming.TypicalPercent {
				bestUpcoming = slot
			}
		}
	}

	reply := &pb.GetLibraryBusynessReply{
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
