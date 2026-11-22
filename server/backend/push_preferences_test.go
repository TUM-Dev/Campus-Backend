package backend

import (
	"context"
	"testing"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/guregu/null"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPushDeviceAndPreferencesLifecycle(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	server := CampusServer{db: db}

	// Updating preferences before registering device returns NotFound.
	_, err := server.UpdatePushPreferences(ctx, &pb.UpdatePushPreferencesRequest{DeviceId: "dev-1"})
	require.Equal(t, codes.NotFound, status.Code(err))

	// Register device via CreateDevice (supports fcm_token or legacy public_key).
	created, err := server.CreateDevice(ctx, &pb.CreateDeviceRequest{
		DeviceId:   "dev-1",
		DeviceType: pb.DeviceType_IOS,
		FcmToken:   "fcm-token-1",
	})
	require.NoError(t, err)
	require.Equal(t, "dev-1", created.DeviceId)

	// Re-registering updates the token cleanly.
	_, err = server.CreateDevice(ctx, &pb.CreateDeviceRequest{
		DeviceId:   "dev-1",
		DeviceType: pb.DeviceType_IOS,
		FcmToken:   "fcm-token-2",
	})
	require.NoError(t, err)
	var stored model.PushDevice
	require.NoError(t, db.First(&stored, "device_id = ?", "dev-1").Error)
	require.Equal(t, "fcm-token-2", stored.FCMToken)

	// Update and read push preferences.
	updated, err := server.UpdatePushPreferences(ctx, &pb.UpdatePushPreferencesRequest{
		DeviceId:             "dev-1",
		CanteenIds:           []string{"MENSA_GARCHING", "mensa-garching"},
		FavoriteDishes:       []string{"Kaiserschmarrn", "  Currywurst  "},
		DietaryLabels:        []string{"vegan"},
		NotifyFavoriteDishes: true,
		NotifyItStatus:       true,
		Language:             pb.Language_German,
	})
	require.NoError(t, err)
	require.Equal(t, []string{"mensa-garching"}, updated.CanteenIds)
	require.Equal(t, []string{"Kaiserschmarrn", "Currywurst"}, updated.FavoriteDishes)
	require.Equal(t, []string{"VEGAN"}, updated.DietaryLabels)
	require.True(t, updated.NotifyFavoriteDishes)
	require.True(t, updated.NotifyItStatus)
	require.Equal(t, pb.Language_German, updated.Language)

	loaded, err := server.GetPushPreferences(ctx, &pb.GetPushPreferencesRequest{DeviceId: "dev-1"})
	require.NoError(t, err)
	require.Equal(t, updated, loaded)

	// DeleteDevice removes device and cascades to push_preferences.
	_, err = server.DeleteDevice(ctx, &pb.DeleteDeviceRequest{DeviceId: "dev-1", DeviceType: pb.DeviceType_IOS})
	require.NoError(t, err)
	_, err = server.GetPushPreferences(ctx, &pb.GetPushPreferencesRequest{DeviceId: "dev-1"})
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestTopRatedDishesDietaryAndPriceFilters(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	server := CampusServer{db: db}

	canteen := model.Canteen{Name: "MENSA_GARCHING"}
	require.NoError(t, db.Create(&canteen).Error)

	veganCheap := model.Dish{
		Name: "Linsen Dal", Type: "Vegan", CanteenID: canteen.CanteenID,
		Labels: null.StringFrom("VEGAN,GLUTEN"), PriceStudentCents: null.IntFrom(320),
	}
	veganExpensive := model.Dish{
		Name: "Truffle Bowl", Type: "Vegan", CanteenID: canteen.CanteenID,
		Labels: null.StringFrom("VEGAN"), PriceStudentCents: null.IntFrom(650),
	}
	meatCheap := model.Dish{
		Name: "Schnitzel", Type: "Fleisch", CanteenID: canteen.CanteenID,
		Labels: null.StringFrom("MEAT"), PriceStudentCents: null.IntFrom(350),
	}
	for _, d := range []*model.Dish{&veganCheap, &veganExpensive, &meatCheap} {
		require.NoError(t, db.Create(d).Error)
		require.NoError(t, db.Create(&model.DishesOfTheWeek{DishID: d.Dish, Year: 2026, Week: 45, Day: 2}).Error)
		require.NoError(t, db.Create(&model.DishRating{DishID: d.Dish, Points: 5, Status: model.RatingStatusVisible}).Error)
	}

	reply, err := server.ListTopRatedDishes(ctx, &pb.ListTopRatedDishesRequest{
		Year:                 2026,
		Week:                 45,
		Day:                  2,
		RequiredLabels:       []string{"vegan"},
		MaxPriceStudentCents: 400,
	})
	require.NoError(t, err)
	require.Len(t, reply.Dishes, 1)
	require.Equal(t, "Linsen Dal", reply.Dishes[0].Dish)
	require.Equal(t, int32(320), reply.Dishes[0].PriceStudentCents)
	require.Equal(t, []string{"VEGAN", "GLUTEN"}, reply.Dishes[0].Labels)
}
