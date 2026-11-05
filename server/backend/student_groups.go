package backend

import (
	"context"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CampusServer) ListStudentGroup(ctx context.Context, req *pb.ListStudentGroupRequest) (*pb.ListStudentGroupReply, error) {
	clubs, err := s.getStudentClubCollections(ctx, req.GetLanguage())
	if err != nil {
		return nil, err
	}
	councils, err := s.getStudentCouncilCollections(ctx, req.GetLanguage())
	if err != nil {
		return nil, err
	}
	return &pb.ListStudentGroupReply{StudentClubs: clubs, StudentCouncils: councils}, nil
}

func (s *CampusServer) getStudentCouncilCollections(ctx context.Context, lang pb.Language) ([]*pb.StudentCouncilCollection, error) {
	var dbCouncils []model.StudentCouncil
	if err := s.db.WithContext(ctx).
		Where(&model.StudentCouncil{Language: lang.String()}).
		Where(&model.StudentCouncilCollection{Language: lang.String()}).
		Joins("Image").
		Joins("StudentCouncilCollection").
		Order("student_councils.id ASC").
		Find(&dbCouncils).Error; err != nil {
		log.WithError(err).Error("Error while querying student councils")
		return nil, status.Error(codes.Internal, "could not query the student councils. Please retry later")
	}

	var dbCollections []model.StudentCouncilCollection
	if err := s.db.WithContext(ctx).
		Where(&model.StudentCouncilCollection{Language: lang.String()}).
		Order("id ASC").
		Find(&dbCollections).Error; err != nil {
		log.WithError(err).Error("Error while querying student council collections")
		return nil, status.Error(codes.Internal, "could not query the student council collections. Please retry later")
	}

	// map from the db to the response
	collections := make([]*pb.StudentCouncilCollection, 0, len(dbCollections))
	byID := make(map[uint]*pb.StudentCouncilCollection, len(dbCollections))
	for _, dbCollection := range dbCollections {
		collection := &pb.StudentCouncilCollection{
			Title:                dbCollection.Name,
			Description:          dbCollection.Description,
			Councils:             make([]*pb.StudentCouncil, 0),
			UnstableCollectionId: uint64(dbCollection.ID),
		}
		collections = append(collections, collection)
		byID[dbCollection.ID] = collection
	}
	for _, dbCouncil := range dbCouncils {
		collection, ok := byID[dbCouncil.StudentCouncilCollectionID]
		if !ok {
			continue
		}
		resCouncil := &pb.StudentCouncil{
			Name:        dbCouncil.Name,
			Description: dbCouncil.Description.Ptr(),
			LinkUrl:     dbCouncil.LinkUrl.Ptr(),
		}
		if dbCouncil.Image != nil {
			cov := dbCouncil.Image.FullExternalUrl()
			resCouncil.CoverUrl = &cov
		}
		collection.Councils = append(collection.Councils, resCouncil)
	}
	return collections, nil
}
