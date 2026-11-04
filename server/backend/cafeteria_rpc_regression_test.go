package backend

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func cafeteriaTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, mock.ExpectationsWereMet()); _ = conn.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: conn, SkipInitializeWithVersion: true}), &gorm.Config{})
	require.NoError(t, err)
	return db, mock
}

func TestCanteenRPCsDispatchToImplementations(t *testing.T) {
	db, mock := cafeteriaTestDB(t)
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pb.RegisterCampusServer(server, &CampusServer{db: db})
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(func() { server.Stop(); _ = listener.Close() })
	conn, err := grpc.NewClient("passthrough:///test", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	client := pb.NewCampusClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery("SELECT name as id,address,latitude,longitude FROM `canteens`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "latitude", "longitude"}).AddRow("MENSA_GARCHING", "Garching", 48.2, 11.6))
	canteens, err := client.ListCanteens(ctx, &pb.ListCanteensRequest{})
	require.NoError(t, err)
	require.Len(t, canteens.Canteen, 1)
	require.Equal(t, "Garching", canteens.Canteen[0].Address)
	require.Equal(t, "MENSA_GARCHING", canteens.Canteen[0].Id)
	mock.ExpectQuery("SELECT d.name FROM dishes_of_the_weeks weekly JOIN dishes d ON d.dish = weekly.dishID JOIN canteens c ON c.cafeteria = d.cafeteriaID").
		WithArgs(int32(2), int32(45), int32(2026), canteens.Canteen[0].Id).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Pasta"))
	dishes, err := client.ListDishes(ctx, &pb.ListDishesRequest{CanteenId: canteens.Canteen[0].Id, Year: 2026, Week: 45, Day: 2})
	require.NoError(t, err)
	require.Equal(t, []string{"Pasta"}, dishes.Dish)

	mock.ExpectQuery("SELECT DE as de, EN as en, cafeteriaRatingTagOption as TagId FROM `canteen_rating_tag_options`").
		WillReturnRows(sqlmock.NewRows([]string{"de", "en", "TagId"}).AddRow("Service", "Service", 1))
	tags, err := client.ListAvailableCanteenTags(ctx, &pb.ListAvailableCanteenTagsRequest{})
	require.NoError(t, err)
	require.Len(t, tags.RatingTags, 1)
	require.Equal(t, int32(1), tags.RatingTags[0].TagId)
}

func TestListDishesUsesRenamedTablesAndSupportsWeek53(t *testing.T) {
	db, mock := cafeteriaTestDB(t)
	server := &CampusServer{db: db}
	mock.ExpectQuery("SELECT d.name FROM dishes_of_the_weeks weekly JOIN dishes d ON d.dish = weekly.dishID JOIN canteens c ON c.cafeteria = d.cafeteriaID").
		WithArgs(int32(2), int32(53), int32(2026), "MENSA_GARCHING").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Pasta"))
	reply, err := server.ListDishes(context.Background(), &pb.ListDishesRequest{CanteenId: "mensa-garching", Year: 2026, Week: 53, Day: 2})
	require.NoError(t, err)
	require.Equal(t, []string{"Pasta"}, reply.Dish)
	for _, req := range []*pb.ListDishesRequest{{Year: 2021, Week: 1}, {Year: 2026, Week: 54}, {Year: 2026, Week: 1, Day: 5}} {
		_, err := server.ListDishes(context.Background(), req)
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	}
}

func TestCanteenRatingTagsUseCurrentTablesAndMapAverage(t *testing.T) {
	db, mock := cafeteriaTestDB(t)
	mock.ExpectQuery("FROM canteen_rating_tag_options options JOIN cafeteria_rating_tag_statistics").
		WithArgs(int32(7)).WillReturnRows(sqlmock.NewRows([]string{"tagId", "avg", "std", "min", "max"}).AddRow(1, 4.25, 0.5, 3, 5))
	tags := queryTags(7, -1, model.CAFETERIA, db)
	require.Len(t, tags, 1)
	require.Equal(t, 4.25, tags[0].Avg)
	mock.ExpectQuery("FROM canteen_rating_tag_options options JOIN canteen_rating_tags rating").
		WithArgs(int64(12)).WillReturnRows(sqlmock.NewRows([]string{"tagId", "points"}).AddRow(1, 4))
	ratings := queryTagRatingsOverviewForRating(12, model.CAFETERIA, db)
	require.Len(t, ratings, 1)
	require.Equal(t, int32(4), ratings[0].Points)
}

func TestCanteenNameAcceptsEatAPIIdentifiers(t *testing.T) {
	db, mock := cafeteriaTestDB(t)
	mock.ExpectQuery("SELECT `cafeteria` FROM `canteens`").WithArgs("MENSA_GARCHING").WillReturnRows(sqlmock.NewRows([]string{"cafeteria"}).AddRow(7))
	require.Equal(t, int32(7), getIDForCafeteriaName("mensa-garching", db))
}
