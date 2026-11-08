package admin

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/stretchr/testify/require"
)

func TestAdminCafeteriaRatings(t *testing.T) {
	db := utils.SetupTestContainer(context.Background(), t)
	canteens := []model.Canteen{
		{Name: "MENSA_GARCHING", Address: "Garching"},
		{Name: "MENSA_ARCISSTR", Address: "Munich"},
		{Name: "MENSA_WEIHENSTEPHAN", Address: "Freising"},
	}
	require.NoError(t, db.Create(&canteens).Error)
	now := time.Now().UTC()
	ratings := []model.CanteenRating{
		{CanteenID: canteens[0].CanteenID, Points: 5, Comment: "<script>rating</script>", Timestamp: now},
		{CanteenID: canteens[0].CanteenID, Points: 3, Comment: "Good choice of meals.", Timestamp: now.Add(-time.Minute)},
		{CanteenID: canteens[1].CanteenID, Points: 1, Comment: "Long queue.", Timestamp: now.Add(-time.Hour)},
	}
	require.NoError(t, db.Create(&ratings).Error)
	option := model.CanteenRatingTagOption{EN: "Cleanliness", DE: "Sauberkeit"}
	require.NoError(t, db.Create(&option).Error)
	require.NoError(t, db.Create(&model.CanteenRatingTag{
		RatingID: ratings[0].CanteenRating, TagID: option.CanteenRatingTagOption, Points: 4,
	}).Error)
	cfg, client := testOIDC(t)
	handler, err := newHandler(db, cron.New(db), "test", cfg, client)
	require.NoError(t, err)
	response, _, _ := loginRequest(t, handler, cfg, nil)
	require.Equal(t, http.StatusSeeOther, response.Code)
	cookie := findCookie(t, response, sessionCookie)
	require.Equal(t, http.StatusSeeOther, adminRequest(handler, http.MethodGet, "/admin/ratings", nil, nil).Code)
	require.Equal(t, http.StatusMethodNotAllowed, adminRequest(handler, http.MethodPost, "/admin/ratings", nil, cookie).Code)

	all := adminRequest(handler, http.MethodGet, "/admin/ratings", nil, cookie)
	require.Equal(t, http.StatusOK, all.Code, all.Body.String())
	require.Contains(t, all.Body.String(), "3.00")
	require.Contains(t, all.Body.String(), "&lt;script&gt;rating&lt;/script&gt;")
	require.NotContains(t, all.Body.String(), "<script>rating</script>")
	require.Contains(t, all.Body.String(), "Cleanliness")
	filter := fmt.Sprintf("/admin/ratings?canteen=%d", canteens[0].CanteenID)
	filtered := adminRequest(handler, http.MethodGet, filter, nil, cookie)
	require.Equal(t, http.StatusOK, filtered.Code, filtered.Body.String())
	require.Contains(t, filtered.Body.String(), "4.00")
	require.NotContains(t, filtered.Body.String(), "Long queue.")
	empty := adminRequest(handler, http.MethodGet, fmt.Sprintf("/admin/ratings?canteen=%d", canteens[2].CanteenID), nil, cookie)
	require.Equal(t, http.StatusOK, empty.Code)
	require.Contains(t, empty.Body.String(), "No ratings yet")
	require.NotContains(t, empty.Body.String(), "0.00")
	require.Equal(t, http.StatusBadRequest, adminRequest(handler, http.MethodGet, "/admin/ratings?canteen=-1", nil, cookie).Code)
	require.Equal(t, http.StatusNotFound, adminRequest(handler, http.MethodGet, "/admin/ratings?canteen=2147483647", nil, cookie).Code)

	var extra []model.CanteenRating
	for i := range pageSize {
		extra = append(extra, model.CanteenRating{CanteenID: canteens[0].CanteenID, Points: 4, Comment: fmt.Sprintf("Extra rating %d", i), Timestamp: now.Add(-time.Duration(i+2) * time.Hour)})
	}
	require.NoError(t, db.Create(&extra).Error)
	first := adminRequest(handler, http.MethodGet, filter, nil, cookie)
	require.Contains(t, first.Body.String(), fmt.Sprintf(`href="?canteen=%d&amp;page=2"`, canteens[0].CanteenID))
	second := adminRequest(handler, http.MethodGet, filter+"&page=2", nil, cookie)
	require.Equal(t, http.StatusOK, second.Code)
	require.Contains(t, second.Body.String(), "Extra rating 49")
	require.NotContains(t, second.Body.String(), "Long queue.")
}
