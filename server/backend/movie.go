package backend

import (
	"context"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *CampusServer) ListMovies(ctx context.Context, req *pb.ListMoviesRequest) (*pb.ListMoviesReply, error) {
	var movies []model.Kino
	tx := s.db.WithContext(ctx).
		Joins("File").
		Order("date ASC")
	if req.OldestDateAt.GetSeconds() != 0 || req.OldestDateAt.GetNanos() != 0 {
		tx = tx.Where("date > ?", req.OldestDateAt.AsTime())
	}
	if err := tx.Find(&movies, "kino > ?", req.LastId).Error; err != nil {
		log.WithError(err).Error("Error while fetching movies from database")
		return nil, status.Error(codes.Internal, "Error while fetching movies from database")
	}
	var movieResponse []*pb.Movie
	for _, movie := range movies {
		movieResponse = append(movieResponse, &pb.Movie{
			MovieId:                  movie.Id,
			Date:                     timestamppb.New(movie.Date),
			Created:                  timestamppb.New(movie.Created),
			Title:                    movie.Title,
			ReleaseYear:              movie.Year.String,
			Runtime:                  movie.Runtime.String,
			Genre:                    movie.Genre.String,
			Director:                 movie.Director.String,
			Actors:                   movie.Actors.String,
			ImdbRating:               movie.ImdbRating.String,
			Description:              movie.Description,
			CoverUrl:                 movie.File.FullExternalUrl(),
			CoverId:                  movie.File.File,
			AdditionalInformationUrl: movie.Link,
			Location:                 movie.Location.String,
		})
	}
	return &pb.ListMoviesReply{
		Movies: movieResponse,
	}, nil
}
