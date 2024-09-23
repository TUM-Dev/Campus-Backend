package backend

import (
	"context"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CampusServer) ListStudentClub(ctx context.Context, req *pb.ListStudentClubRequest) (*pb.ListStudentClubReply, error) {
	var dbClubs []model.StudentClub
	if err := s.db.WithContext(ctx).
		Where("language = ?", req.GetLanguage().String()).
		Joins("Image").
		Joins("StudentClubCollection").
		Find(&dbClubs).Error; err != nil {
		log.WithError(err).Error("Error while querying student clubs")
		return nil, status.Error(codes.Internal, "could not query the student clubs. Please retry later")
	}

	var dbClubCollections []model.StudentClubCollection
	if err := s.db.WithContext(ctx).
		Where("language = ?", req.GetLanguage().String()).
		Find(&dbClubCollections).Error; err != nil {
		log.WithError(err).Error("Error while querying student club collections")
		return nil, status.Error(codes.Internal, "could not query the student club collections. Please retry later")
	}
	// map from the db to the response
	collections := make([]*pb.StudentClubCollection, 0)
	for _, dbCollection := range dbClubCollections {
		collections = append(collections, &pb.StudentClubCollection{
			Title:                dbCollection.Name,
			Description:          dbCollection.Description,
			Clubs:                make([]*pb.StudentClub, 0),
			UnstableCollectionId: 0,
		})
	}
	for _, dbClub := range dbClubs {
		resClub := &pb.StudentClub{
			Name:        dbClub.Name,
			Description: dbClub.Description.Ptr(),
			LinkUrl:     dbClub.LinkUrl.Ptr(),
		}
		if dbClub.Image != nil {
			cov := dbClub.Image.FullExternalUrl()
			resClub.CoverUrl = &cov // go does not allow inlining here
		}

		for _, collection := range collections {
			if collection.UnstableCollectionId == uint64(dbClub.StudentClubCollectionID) {
				collection.Clubs = append(collection.Clubs, resClub)
				resClub = nil
				break
			}
		}
	}

	return &pb.ListStudentClubReply{Collections: collections}, nil
}
