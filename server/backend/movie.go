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

func (s *CampusServer) GetMovies(ctx context.Context, req *pb.GetMoviesRequest) (*pb.GetMoviesReply, error) {
	var movies []model.Kino
	if err := s.db.WithContext(ctx).Joins("Files").Find(&movies, "kino > ?", req.LastId).Error; err != nil {
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
			CoverName:   movie.Files.Name,
			CoverPath:   movie.Files.Path,
			CoverId:     movie.Files.File,
			Link:        movie.Link,
		})
	}
	return &pb.GetMoviesReply{
		Movies: movieResponse,
	}, nil
}
