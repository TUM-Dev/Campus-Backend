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
		Id:          1,
		Date:        timestamppb.New(time.Now()),
		Created:     timestamppb.New(time.Now()),
		Title:       "Mission Impossible 4 - Ghost Protocol",
		Year:        "2011",
		Runtime:     "133 min",
		Genre:       "Action, Adventure, Thriller",
		Director:    "Brad Bird",
		Actors:      "Tom Cruise, Jeremy Renner, Simon Pegg, Paula Patton",
		ImdbRating:  "7.4",
		Description: "The IMF is shut down when it's implicated in the bombing of the Kremlin, causing Ethan Hunt and his new team to go rogue to clear their organization's name.",
		CoverName:   "mission_impossible_4.jpg",
		CoverPath:   "movie/mission_impossible_4.jpg",
		CoverID:     1,
		Link:        "https://www.imdb.com/title/tt1229238/",
	}
	movie2 = pb.Movie{
		Id:          2,
		Date:        timestamppb.New(time.Now()),
		Created:     timestamppb.New(time.Now()),
		Title:       "Mission Impossible 5 - Rogue Nation",
		Year:        "2015",
		Runtime:     "131 min",
		Genre:       "Action, Adventure, Thriller",
		Director:    "Christopher McQuarrie",
		Actors:      "Tom Cruise, Jeremy Renner, Simon Pegg, Rebecca Ferguson",
		ImdbRating:  "7.4",
		Description: "Ethan and his team take on their most impossible mission yet when they have to eradicate an international rogue organization as highly skilled as they are and committed to destroying the IMF.",
		CoverName:   "mission_impossible_5.jpg",
		CoverPath:   "movie/mission_impossible_5.jpg",
		CoverID:     2,
		Link:        "https://www.imdb.com/title/tt2381249/",
	}
)

func (s *MovieSuite) Test_GetMoviesAll() {
	server := CampusServer{db: s.DB}
	s.mock.ExpectQuery("SELECT `kino`.`kino`,`kino`.`date`,`kino`.`created`,`kino`.`title`,`kino`.`year`,`kino`.`runtime`,`kino`.`genre`,`kino`.`director`,`kino`.`actors`,`kino`.`rating`,`kino`.`description`,`kino`.`trailer`,`kino`.`cover`,`kino`.`link`,`Files`.`file` AS `Files__file`,`Files`.`name` AS `Files__name`,`Files`.`path` AS `Files__path`,`Files`.`downloads` AS `Files__downloads`,`Files`.`url` AS `Files__url`,`Files`.`downloaded` AS `Files__downloaded` FROM `kino` LEFT JOIN `files` `Files` ON `kino`.`cover` = `Files`.`file` WHERE kino > ?").
		WithArgs(-1).
		WillReturnRows(sqlmock.NewRows([]string{"kino", "date", "created", "title", "year", "runtime", "genre", "director", "actors", "rating", "description", "trailer", "cover", "link", "Files__file", "Files__name", "Files__path", "Files__downloads", "Files__url", "Files__downloaded"}).
			AddRow(movie2.Id, movie2.Date.AsTime(), movie2.Created.AsTime(), movie2.Title, movie2.Year, movie2.Runtime, movie2.Genre, movie2.Director, movie2.Actors, movie2.ImdbRating, movie2.Description, nil, movie2.CoverID, movie2.Link, movie2.CoverID, movie2.CoverName, movie2.CoverPath, 1, "", 1).
			AddRow(movie1.Id, movie1.Date.AsTime(), movie1.Created.AsTime(), movie1.Title, movie1.Year, movie1.Runtime, movie1.Genre, movie1.Director, movie1.Actors, movie1.ImdbRating, movie1.Description, nil, movie1.CoverID, movie1.Link, movie1.CoverID, movie1.CoverName, movie1.CoverPath, 1, "", 1))
	response, err := server.GetMovies(context.Background(), &pb.GetMoviesRequest{LastId: -1})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.GetMoviesReply{Movies: []*pb.Movie{&movie2, &movie1}}, response)
}

func (s *MovieSuite) Test_GetMoviesOne() {
	server := CampusServer{db: s.DB}
	s.mock.ExpectQuery("SELECT `kino`.`kino`,`kino`.`date`,`kino`.`created`,`kino`.`title`,`kino`.`year`,`kino`.`runtime`,`kino`.`genre`,`kino`.`director`,`kino`.`actors`,`kino`.`rating`,`kino`.`description`,`kino`.`trailer`,`kino`.`cover`,`kino`.`link`,`Files`.`file` AS `Files__file`,`Files`.`name` AS `Files__name`,`Files`.`path` AS `Files__path`,`Files`.`downloads` AS `Files__downloads`,`Files`.`url` AS `Files__url`,`Files`.`downloaded` AS `Files__downloaded` FROM `kino` LEFT JOIN `files` `Files` ON `kino`.`cover` = `Files`.`file` WHERE kino > ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"kino", "date", "created", "title", "year", "runtime", "genre", "director", "actors", "rating", "description", "trailer", "cover", "link", "Files__file", "Files__name", "Files__path", "Files__downloads", "Files__url", "Files__downloaded"}).
			AddRow(movie1.Id, movie1.Date.AsTime(), movie1.Created.AsTime(), movie1.Title, movie1.Year, movie1.Runtime, movie1.Genre, movie1.Director, movie1.Actors, movie1.ImdbRating, movie1.Description, nil, movie1.CoverID, movie1.Link, movie1.CoverID, movie1.CoverName, movie1.CoverPath, 1, "", 1))
	response, err := server.GetMovies(context.Background(), &pb.GetMoviesRequest{LastId: 1})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.GetMoviesReply{Movies: []*pb.Movie{&movie1}}, response)
}

func (s *MovieSuite) Test_GetMoviesNone() {
	server := CampusServer{db: s.DB}
	s.mock.ExpectQuery("SELECT `kino`.`kino`,`kino`.`date`,`kino`.`created`,`kino`.`title`,`kino`.`year`,`kino`.`runtime`,`kino`.`genre`,`kino`.`director`,`kino`.`actors`,`kino`.`rating`,`kino`.`description`,`kino`.`trailer`,`kino`.`cover`,`kino`.`link`,`Files`.`file` AS `Files__file`,`Files`.`name` AS `Files__name`,`Files`.`path` AS `Files__path`,`Files`.`downloads` AS `Files__downloads`,`Files`.`url` AS `Files__url`,`Files`.`downloaded` AS `Files__downloaded` FROM `kino` LEFT JOIN `files` `Files` ON `kino`.`cover` = `Files`.`file` WHERE kino > ?").
		WithArgs(42).
		WillReturnRows(sqlmock.NewRows([]string{"kino", "date", "created", "title", "year", "runtime", "genre", "director", "actors", "rating", "description", "trailer", "cover", "link", "Files__file", "Files__name", "Files__path", "Files__downloads", "Files__url", "Files__downloaded"}))
	response, err := server.GetMovies(context.Background(), &pb.GetMoviesRequest{LastId: 42})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.GetMoviesReply{Movies: []*pb.Movie(nil)}, response)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(MovieSuite))
}
