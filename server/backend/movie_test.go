package backend

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MovieSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *MovieSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	dialector := mysql.New(mysql.Config{
		Conn:       db,
		DriverName: "mysql",
	})
	s.mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("10.11.4-MariaDB"))
	s.DB, err = gorm.Open(dialector, &gorm.Config{})
	require.NoError(s.T(), err)
}

var (
	movie1 = pb.Movie{
		MovieId:                  1,
		Date:                     timestamppb.New(time.Now()),
		Created:                  timestamppb.New(time.Now()),
		Title:                    "Mission Impossible 4 - Ghost Protocol",
		ReleaseYear:              "2011",
		Runtime:                  "133 min",
		Genre:                    "Action, Adventure, Thriller",
		Director:                 "Brad Bird",
		Actors:                   "Tom Cruise, Jeremy Renner, Simon Pegg, Paula Patton",
		ImdbRating:               "7.4",
		Description:              "The IMF is shut down when it's implicated in the bombing of the Kremlin, causing Ethan Hunt and his new team to go rogue to clear their organization's name.",
		CoverUrl:                 "https://api.tum.app/files/movie/mission_impossible_4.jpg",
		CoverId:                  1,
		AdditionalInformationUrl: "https://www.imdb.com/title/tt1229238/",
		Location:                 "Hörsaal 1200, Stammgelände",
	}
	movie2 = pb.Movie{
		MovieId:                  2,
		Date:                     timestamppb.New(time.Now()),
		Created:                  timestamppb.New(time.Now()),
		Title:                    "Mission Impossible 5 - Rogue Nation",
		ReleaseYear:              "2015",
		Runtime:                  "131 min",
		Genre:                    "Action, Adventure, Thriller",
		Director:                 "Christopher McQuarrie",
		Actors:                   "Tom Cruise, Jeremy Renner, Simon Pegg, Rebecca Ferguson",
		ImdbRating:               "7.4",
		Description:              "Ethan and his team take on their most impossible mission yet when they have to eradicate an international rogue organization as highly skilled as they are and committed to destroying the IMF.",
		CoverUrl:                 "https://api.tum.app/files/movie/mission_impossible_5.jpg",
		CoverId:                  2,
		AdditionalInformationUrl: "https://www.imdb.com/title/tt2381249/",
		Location:                 "Hörsaal MW1801, Campus Garching",
	}
)

const ListMoviesQuery = "SELECT `kino`.`kino`,`kino`.`date`,`kino`.`created`,`kino`.`title`,`kino`.`year`,`kino`.`runtime`,`kino`.`genre`,`kino`.`director`,`kino`.`actors`,`kino`.`rating`,`kino`.`description`,`kino`.`trailer`,`kino`.`cover`,`kino`.`link`,`kino`.`location`,`File`.`file` AS `File__file`,`File`.`name` AS `File__name`,`File`.`path` AS `File__path`,`File`.`downloads` AS `File__downloads`,`File`.`url` AS `File__url`,`File`.`downloaded` AS `File__downloaded` FROM `kino` LEFT JOIN `files` `File` ON `kino`.`cover` = `File`.`file` WHERE kino > ?"

func (s *MovieSuite) Test_ListMoviesAll() {
	server := CampusServer{db: s.DB}
	s.mock.ExpectQuery(ListMoviesQuery).
		WithArgs(-1).
		WillReturnRows(sqlmock.NewRows([]string{"kino", "date", "created", "title", "year", "runtime", "genre", "director", "actors", "rating", "description", "trailer", "cover", "link", "location", "File__file", "File__name", "File__path", "File__downloads", "File__url", "File__downloaded"}).
			AddRow(movie2.MovieId, movie2.Date.AsTime(), movie2.Created.AsTime(), movie2.Title, movie2.ReleaseYear, movie2.Runtime, movie2.Genre, movie2.Director, movie2.Actors, movie2.ImdbRating, movie2.Description, nil, movie2.CoverId, movie2.AdditionalInformationUrl, movie2.Location, movie2.CoverId, "mission_impossible_5.jpg", "movie/", 1, "", 1).
			AddRow(movie1.MovieId, movie1.Date.AsTime(), movie1.Created.AsTime(), movie1.Title, movie1.ReleaseYear, movie1.Runtime, movie1.Genre, movie1.Director, movie1.Actors, movie1.ImdbRating, movie1.Description, nil, movie1.CoverId, movie1.AdditionalInformationUrl, movie1.Location, movie1.CoverId, "mission_impossible_4.jpg", "movie/", 1, "", 1))
	response, err := server.ListMovies(context.Background(), &pb.ListMoviesRequest{LastId: -1})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.ListMoviesReply{Movies: []*pb.Movie{&movie2, &movie1}}, response)
}

func (s *MovieSuite) Test_ListMoviesOne() {
	server := CampusServer{db: s.DB}
	s.mock.ExpectQuery(ListMoviesQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"kino", "date", "created", "title", "year", "runtime", "genre", "director", "actors", "rating", "description", "trailer", "cover", "link", "location", "File__file", "File__name", "File__path", "File__downloads", "File__url", "File__downloaded"}).
			AddRow(movie1.MovieId, movie1.Date.AsTime(), movie1.Created.AsTime(), movie1.Title, movie1.ReleaseYear, movie1.Runtime, movie1.Genre, movie1.Director, movie1.Actors, movie1.ImdbRating, movie1.Description, nil, movie1.CoverId, movie1.AdditionalInformationUrl, movie1.Location, movie1.CoverId, "mission_impossible_4.jpg", "movie/", 1, "", 1))
	response, err := server.ListMovies(context.Background(), &pb.ListMoviesRequest{LastId: 1})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.ListMoviesReply{Movies: []*pb.Movie{&movie1}}, response)
}

func (s *MovieSuite) Test_ListMoviesNone() {
	server := CampusServer{db: s.DB}
	s.mock.ExpectQuery(ListMoviesQuery).
		WithArgs(42).
		WillReturnRows(sqlmock.NewRows([]string{"kino", "date", "created", "title", "year", "runtime", "genre", "director", "actors", "rating", "description", "trailer", "cover", "link", "location", "File__file", "File__name", "File__path", "File__downloads", "File__url", "File__downloaded"}))
	response, err := server.ListMovies(context.Background(), &pb.ListMoviesRequest{LastId: 42})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.ListMoviesReply{Movies: []*pb.Movie(nil)}, response)
}

func (s *MovieSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMovieSuite(t *testing.T) {
	suite.Run(t, new(MovieSuite))
}
