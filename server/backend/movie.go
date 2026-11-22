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
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}
	movies, err := s.getMovies(ctx, req.LastId, req.OldestDateAt.AsTime(), req.Cinema)
	if err != nil {
		return nil, err
	}
	var movieResponse []*pb.Movie
	for _, movie := range movies {
		cinema := movieCinema(movie)
		if req.Cinema != pb.Cinema_CINEMA_UNSPECIFIED && cinema != req.Cinema {
			continue
		}
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
			TrailerUrl:               movie.Trailer.String,
			CoverUrl:                 movie.File.FullExternalUrl(),
			CoverId:                  movie.File.File,
			AdditionalInformationUrl: movie.Link,
			Location:                 movie.Location.String,
			Cinema:                   cinema,
		})
	}
	return &pb.ListMoviesReply{
		Movies: movieResponse,
	}, nil
}

func movieCinema(movie model.Movie) pb.Cinema {
	if movie.IsCampuskinoHeilbronn() {
		return pb.Cinema_CAMPUSKINO_HEILBRONN
	}
	return pb.Cinema_TU_FILM
}

func (s *CampusServer) getMovies(ctx context.Context, lastID int32, oldestDateAt time.Time, cinema pb.Cinema) ([]model.Movie, error) {
	cacheKey := fmt.Sprintf("%d-%d", lastID, oldestDateAt.UnixNano())
	if movies, ok := s.moviesCache.Get(cacheKey); ok {
		if cinema != pb.Cinema_CINEMA_UNSPECIFIED && !hasMovieForCinema(movies, cinema) {
			s.moviesCache.Remove(cacheKey)
		} else {
			return movies, nil
		}
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
	if len(movies) > 0 && (cinema == pb.Cinema_CINEMA_UNSPECIFIED || hasMovieForCinema(movies, cinema)) {
		s.moviesCache.Add(cacheKey, movies)
	}
	return movies, nil
}

func hasMovieForCinema(movies []model.Movie, cinema pb.Cinema) bool {
	for _, m := range movies {
		if movieCinema(m) == cinema {
			return true
		}
	}
	return false
}
