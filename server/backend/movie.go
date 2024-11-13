package backend

import (
	"context"
	"fmt"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *CampusServer) ListMovies(ctx context.Context, req *pb.ListMoviesRequest) (*pb.ListMoviesReply, error) {
	movies, err := s.getMovies(ctx, req.LastId, req.OldestDateAt.AsTime())
	if err != nil {
		return nil, err
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

func (s *CampusServer) getMovies(ctx context.Context, lastID int32, oldestDateAt time.Time) ([]model.Movie, error) {
	cacheKey := fmt.Sprintf("%d-%d", lastID, oldestDateAt.Second())
	if movies, ok := s.moviesCache.Get(cacheKey); ok {
		return movies, nil
	}
	var movies []model.Movie
	tx := s.db.WithContext(ctx).
		Joins("File").
		Order("date ASC")
	if oldestDateAt.Unix() != 0 {
		tx = tx.Where("date > ?", oldestDateAt)
	}
	if err := tx.Find(&movies, "kino > ?", lastID).Error; err != nil {
		log.WithError(err).Error("Error while fetching movies from database")
		return nil, status.Error(codes.Internal, "Error while fetching movies from database")
	}
	s.moviesCache.Add(cacheKey, movies)
	return movies, nil
}
