package backend

import (
	"context"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const soldOutConfirmThreshold int32 = 10

func berlinDateString(now time.Time) string {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		loc = time.UTC
	}
	return now.In(loc).Format("2006-01-02")
}

// ReportDishSoldOut records a device-deduplicated "sold out" report for today in Europe/Berlin.
// Reports reset automatically at Berlin midnight; a dish is marked confirmed once at least
// soldOutConfirmThreshold distinct devices report it today.
func (s *CampusServer) ReportDishSoldOut(ctx context.Context, req *pb.ReportDishSoldOutRequest) (*pb.ReportDishSoldOutReply, error) {
	s.recordDeviceIfPresent(ctx)
	return s.reportDishSoldOutAt(ctx, req, time.Now())
}

func (s *CampusServer) reportDishSoldOutAt(ctx context.Context, req *pb.ReportDishSoldOutRequest, now time.Time) (*pb.ReportDishSoldOutReply, error) {
	deviceHash := callerDeviceHash(ctx)
	if !deviceHash.Valid {
		return nil, status.Error(codes.Unauthenticated, "Device identification (x-device-id) is required to report a sold-out dish.")
	}
	if req.GetCanteenId() == "" || req.GetDish() == "" {
		return nil, status.Error(codes.InvalidArgument, "canteen_id and dish are required.")
	}

	tx := s.db.WithContext(ctx)
	canteenID, err := getIDForCanteenName(req.CanteenId, tx)
	if err != nil {
		return nil, err
	}
	dishID, err := getIDForDishName(req.Dish, canteenID, tx)
	if err != nil {
		return nil, err
	}

	reportDate := berlinDateString(now)
	entry := model.DishSoldOutReport{
		DishID:     int64(dishID),
		DeviceHash: deviceHash.String,
		ReportDate: reportDate,
		ReportedAt: now,
	}
	res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&entry)
	if res.Error != nil {
		log.WithError(res.Error).Error("while creating dish sold-out report")
		return nil, status.Error(codes.Internal, "Could not record sold-out report.")
	}
	alreadyReported := res.RowsAffected == 0

	var count int64
	if err := tx.Model(&model.DishSoldOutReport{}).
		Where("dish_id = ? AND report_date = ?", dishID, reportDate).
		Count(&count).Error; err != nil {
		log.WithError(err).Error("while counting dish sold-out reports")
		return nil, status.Error(codes.Internal, "Could not count sold-out reports.")
	}

	return &pb.ReportDishSoldOutReply{
		ReportCount:     int32(count),
		Confirmed:       int32(count) >= soldOutConfirmThreshold,
		AlreadyReported: alreadyReported,
	}, nil
}

type soldOutCountRow struct {
	DishID int64 `gorm:"column:dish_id"`
	Count  int32 `gorm:"column:cnt"`
}

func batchLoadTodaySoldOutCounts(dishIDs []int64, now time.Time, tx *gorm.DB) map[int64]int32 {
	out := make(map[int64]int32)
	if len(dishIDs) == 0 {
		return out
	}
	var rows []soldOutCountRow
	if err := tx.Table("dish_sold_out_reports").
		Select("dish_id, COUNT(*) AS cnt").
		Where("dish_id IN ? AND report_date = ?", dishIDs, berlinDateString(now)).
		Group("dish_id").
		Scan(&rows).Error; err != nil {
		log.WithError(err).Warn("while batch loading sold-out report counts")
		return out
	}
	for _, r := range rows {
		out[r.DishID] = r.Count
	}
	return out
}
