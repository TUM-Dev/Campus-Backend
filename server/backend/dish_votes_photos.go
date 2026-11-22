package backend

import (
	"context"
	"errors"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/guregu/null"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	defaultDishPhotosLimit = 12
	maxDishPhotosLimit     = 50
)

// VoteDishRating casts (helpful=true) or removes (helpful=false) a device's helpful vote on a dish rating.
func (s *CampusServer) VoteDishRating(ctx context.Context, req *pb.VoteDishRatingRequest) (*pb.VoteDishRatingReply, error) {
	s.recordDeviceIfPresent(ctx)
	deviceHash := callerDeviceHash(ctx)
	if !deviceHash.Valid {
		return nil, status.Error(codes.Unauthenticated, "Device identification (x-device-id) is required to vote on a rating.")
	}
	if req.RatingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "A valid rating_id is required.")
	}

	tx := s.db.WithContext(ctx)
	var rating model.DishRating
	if err := tx.Select("dishRating").First(&rating, "dishRating = ? AND status = ?", req.RatingId, model.RatingStatusVisible).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Rating not found.")
		}
		log.WithError(err).Error("while loading dish rating for vote")
		return nil, status.Error(codes.Internal, "Could not load rating.")
	}

	if req.Helpful {
		vote := model.DishRatingVote{
			RatingID:   rating.DishRating,
			DeviceHash: deviceHash.String,
			CreatedAt:  time.Now(),
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&vote).Error; err != nil {
			log.WithError(err).Error("while recording helpful vote")
			return nil, status.Error(codes.Internal, "Could not save vote.")
		}
	} else {
		if err := tx.Where("rating_id = ? AND device_hash = ?", rating.DishRating, deviceHash.String).
			Delete(&model.DishRatingVote{}).Error; err != nil {
			log.WithError(err).Error("while removing helpful vote")
			return nil, status.Error(codes.Internal, "Could not remove vote.")
		}
	}

	var totalVotes int64
	if err := tx.Model(&model.DishRatingVote{}).Where("rating_id = ?", rating.DishRating).Count(&totalVotes).Error; err != nil {
		log.WithError(err).Error("while counting helpful votes")
		return nil, status.Error(codes.Internal, "Could not count votes.")
	}
	return &pb.VoteDishRatingReply{
		HelpfulCount: totalVotes,
		UserVoted:    req.Helpful,
	}, nil
}

type voteCountRow struct {
	RatingID int64 `gorm:"column:rating_id"`
	Count    int64 `gorm:"column:cnt"`
}

func batchLoadHelpfulVoteCounts(ratingIDs []int64, tx *gorm.DB) map[int64]int64 {
	out := make(map[int64]int64)
	if len(ratingIDs) == 0 {
		return out
	}
	var rows []voteCountRow
	if err := tx.Table("dish_rating_votes").
		Select("rating_id, COUNT(*) AS cnt").
		Where("rating_id IN ?", ratingIDs).
		Group("rating_id").
		Scan(&rows).Error; err != nil {
		log.WithError(err).Warn("while batch loading helpful vote counts")
		return out
	}
	for _, r := range rows {
		out[r.RatingID] = r.Count
	}
	return out
}

func batchLoadUserVotedRatings(ratingIDs []int64, deviceHash null.String, tx *gorm.DB) map[int64]bool {
	out := make(map[int64]bool)
	if len(ratingIDs) == 0 || !deviceHash.Valid {
		return out
	}
	var votedIDs []int64
	if err := tx.Table("dish_rating_votes").
		Where("rating_id IN ? AND device_hash = ?", ratingIDs, deviceHash.String).
		Pluck("rating_id", &votedIDs).Error; err != nil {
		log.WithError(err).Warn("while loading user helpful votes")
		return out
	}
	for _, id := range votedIDs {
		out[id] = true
	}
	return out
}

type dishPhotoRow struct {
	RatingID     int64     `gorm:"column:rating_id"`
	Points       int32     `gorm:"column:points"`
	Comment      string    `gorm:"column:comment"`
	Timestamp    time.Time `gorm:"column:timestamp"`
	HelpfulCount int64     `gorm:"column:helpful_count"`
}

// ListDishPhotos returns a paginated gallery of visible photos for a dish, ranked by helpful votes and recency.
func (s *CampusServer) ListDishPhotos(ctx context.Context, req *pb.ListDishPhotosRequest) (*pb.ListDishPhotosReply, error) {
	s.recordDeviceIfPresent(ctx)
	if req.GetCanteenId() == "" || req.GetDish() == "" {
		return nil, status.Error(codes.InvalidArgument, "canteen_id and dish are required.")
	}
	if req.Offset < 0 {
		return nil, status.Error(codes.InvalidArgument, "offset must be non-negative.")
	}
	limit := int(req.Limit)
	if limit <= 0 {
		limit = defaultDishPhotosLimit
	}
	if limit > maxDishPhotosLimit {
		limit = maxDishPhotosLimit
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

	var total int64
	if err := tx.Model(&model.DishRating{}).
		Where("dishID = ? AND status = ? AND image <> ''", dishID, model.RatingStatusVisible).
		Count(&total).Error; err != nil {
		log.WithError(err).Error("while counting dish photos")
		return nil, status.Error(codes.Internal, "Could not count dish photos.")
	}
	if total == 0 {
		return &pb.ListDishPhotosReply{Photos: []*pb.DishPhoto{}, Total: 0, HasMore: false}, nil
	}

	var rows []dishPhotoRow
	if err := tx.Table("dish_ratings r").
		Joins("LEFT JOIN (SELECT rating_id, COUNT(*) AS cnt FROM dish_rating_votes GROUP BY rating_id) v ON v.rating_id = r.dishRating").
		Select("r.dishRating AS rating_id, r.points AS points, r.comment AS comment, r.timestamp AS timestamp, COALESCE(v.cnt, 0) AS helpful_count").
		Where("r.dishID = ? AND r.status = ? AND r.image <> ''", dishID, model.RatingStatusVisible).
		Order("helpful_count DESC, r.timestamp DESC, r.dishRating DESC").
		Limit(limit).
		Offset(int(req.Offset)).
		Scan(&rows).Error; err != nil {
		log.WithError(err).Error("while querying dish photos")
		return nil, status.Error(codes.Internal, "Could not load dish photos.")
	}

	photos := make([]*pb.DishPhoto, 0, len(rows))
	for _, r := range rows {
		photos = append(photos, &pb.DishPhoto{
			RatingId:     r.RatingID,
			ImageUrl:     dishRatingImageURL(r.RatingID),
			Points:       r.Points,
			Comment:      r.Comment,
			HelpfulCount: r.HelpfulCount,
			Visited:      timestamppb.New(r.Timestamp),
		})
	}
	return &pb.ListDishPhotosReply{
		Photos:  photos,
		Total:   total,
		HasMore: int64(req.Offset)+int64(len(photos)) < total,
	}, nil
}
