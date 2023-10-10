package backend

import (
	"context"
	"fmt"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *CampusServer) ListMovies(ctx context.Context, req *pb.ListMoviesRequest) (*pb.ListMoviesReply, error) {
	var movies []model.Kino
	tx := s.db.WithContext(ctx).Joins("File")
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
			MovieId:     movie.Id,
			Date:        timestamppb.New(movie.Date),
			Created:     timestamppb.New(movie.Created),
			Title:       movie.Title,
			ReleaseYear: movie.Year,
			Runtime:     movie.Runtime,
			Genre:       movie.Genre,
			Director:    movie.Director,
			Actors:      movie.Actors,
			ImdbRating:  movie.ImdbRating,
			Description: movie.Description,
			CoverUrl:    fmt.Sprintf("https://api.tum.app/files/%s/%s", movie.File.Path, movie.File.Name),
			CoverId:     movie.File.File,
			Link:        movie.Link,
		})
	}
	return &pb.ListMoviesReply{
		Movies: movieResponse,
	}, nil
}
