package backend

import (
	"context"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CampusServer) ListStudentClubs(ctx context.Context, _ *pb.ListStudentClubRequest) (*pb.ListStudentClubReply, error) {
	var dbClubs []model.StudentClub
	if err := s.db.WithContext(ctx).
		Joins("Image").
		Joins("StudentClubCollection").
		Find(&dbClubs).Error; err != nil {
		log.WithError(err).Error("Error while querying student clubs")
		return nil, status.Error(codes.Internal, "could not query the student clubs. Please retry later")
	}
	// map from the db to the response
	collections := make([]*pb.StudentClubCollection, 0)
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
			if collection.Title == dbClub.StudentClubCollectionID {
				collection.Clubs = append(collection.Clubs, resClub)
				resClub = nil
				break
			}
		}
		// no collection matched => we need to insert a new collection
		if resClub != nil {
			collection := &pb.StudentClubCollection{
				Title:       dbClub.StudentClubCollection.ID,
				Description: dbClub.StudentClubCollection.Description,
				Clubs:       make([]*pb.StudentClub, 1),
			}
			collection.Clubs[0] = resClub
			collections = append(collections, collection)
		}
	}

	return &pb.ListStudentClubReply{Collections: collections}, nil
}
