package backend

import (
	"context"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CampusServer) ListStudentClub(ctx context.Context, req *pb.ListCampusRequest) (*pb.ListCampusReply, error) {
	studentClubs, err := s.getAllStudentClubs(ctx, req.GetLanguage())
	if err != nil {
		return nil, status.Error(codes.Internal, "could not query the student clubs. Please retry later")
	}
	studentCouncils, err := s.getAllStudentCouncils(ctx, req.GetLanguage())
	if err != nil {
		return nil, status.Error(codes.Internal, "could not query the student clubs. Please retry later")
	}

	return &pb.ListCampusReply{StudentClubs: studentClubs, StudentCouncils: studentCouncils}, nil
}

func (s *CampusServer) getAllStudentClubs(ctx context.Context, lang pb.Language) ([]*pb.StudentClubCollection, error) {
	var dbClubs []model.StudentClub
	if err := s.db.WithContext(ctx).
		Where(&model.StudentClub{Language: lang.String()}).
		Where(&model.StudentClubCollection{Language: lang.String()}).
		Joins("Image").
		Joins("StudentClubCollection").
		Find(&dbClubs).Error; err != nil {
		log.WithError(err).Error("Error while querying student clubs")
		return nil, err
	}

	var dbClubCollections []model.StudentClubCollection
	if err := s.db.WithContext(ctx).
		Where(&model.StudentClubCollection{Language: lang.String()}).
		Find(&dbClubCollections).Error; err != nil {
		log.WithError(err).Error("Error while querying student club collections")
		return nil, err
	}
	// map from the db to the response
	collections := make([]*pb.StudentClubCollection, 0)
	for _, dbCollection := range dbClubCollections {
		collections = append(collections, &pb.StudentClubCollection{
			Title:                dbCollection.Name,
			Description:          dbCollection.Description,
			Clubs:                make([]*pb.StudentClub, 0),
			UnstableCollectionId: uint64(dbCollection.ID),
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
				break
			}
		}
	}
	return collections, nil
}

func (s *CampusServer) getAllStudentCouncils(ctx context.Context, lang pb.Language) ([]*pb.StudentCouncilCollection, error) {
	var dbCouncils []model.StudentCouncil
	if err := s.db.WithContext(ctx).
		Where(&model.StudentCouncil{Language: lang.String()}).
		Where(&model.StudentCouncilCollection{Language: lang.String()}).
		Joins("Image").
		Joins("StudentCouncilCollection").
		Find(&dbCouncils).Error; err != nil {
		log.WithError(err).Error("Error while querying student councils")
		return nil, err
	}

	var dbCouncilCollections []model.StudentCouncilCollection
	if err := s.db.WithContext(ctx).
		Where(&model.StudentCouncilCollection{Language: lang.String()}).
		Find(&dbCouncilCollections).Error; err != nil {
		log.WithError(err).Error("Error while querying student council collections")
		return nil, err
	}
	// map from the db to the response
	collections := make([]*pb.StudentCouncilCollection, 0)
	for _, dbCollection := range dbCouncilCollections {
		collections = append(collections, &pb.StudentCouncilCollection{
			Title:                dbCollection.Name,
			Description:          dbCollection.Description,
			Councils:             make([]*pb.StudentCouncil, 0),
			UnstableCollectionId: uint64(dbCollection.ID),
		})
	}
	for _, dbCouncil := range dbCouncils {
		resCouncil := &pb.StudentCouncil{
			Name:        dbCouncil.Name,
			Description: dbCouncil.Description.Ptr(),
			LinkUrl:     dbCouncil.LinkUrl.Ptr(),
		}
		if dbCouncil.Image != nil {
			cov := dbCouncil.Image.FullExternalUrl()
			resCouncil.CoverUrl = &cov // go does not allow inlining here
		}

		for _, collection := range collections {
			if collection.UnstableCollectionId == uint64(dbCouncil.StudentCouncilCollectionID) {
				collection.Councils = append(collection.Councils, resCouncil)
				break
			}
		}
	}
	return collections, nil
}
