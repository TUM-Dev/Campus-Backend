package cron

import (
	"context"
	"strings"
	"testing"

	"github.com/TUM-Dev/Campus-Backend/server/internal/fcm"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/guregu/null"
	"github.com/stretchr/testify/require"
)

type recordingFCMSender struct {
	sent []fcm.Message
}

func (r *recordingFCMSender) Send(_ context.Context, msg fcm.Message) error {
	r.sent = append(r.sent, msg)
	return nil
}

func TestParseTrackedServiceStatusesAndTransitions(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)

	require.NoError(t, db.Create(&model.PushDevice{
		DeviceID: "dev-it", DeviceType: "ios", FCMToken: "tok-it",
	}).Error)
	require.NoError(t, db.Create(&model.PushPreference{
		DeviceID: "dev-it", CanteenIDs: "[]", FavoriteDishes: "[]", DietaryLabels: "[]",
		NotifyITStatus: true, Language: "en",
	}).Error)

	const outageHTML = `
	<html><body>
		<h2>Aktuelle Meldungen</h2>
		<div>
			<div class="table-headers"></div>
			<a class="table-entry">
				<span class="entry-service"><span class="mobile-bold">Moodle</span></span>
				<span class="entry-problemtype">Login Störung</span>
			</a>
			<a class="table-entry">
				<span class="entry-service"><span class="mobile-bold">eduroam WLAN</span></span>
				<span class="entry-problemtype">Wartungsarbeiten</span>
			</a>
		</div>
		<h2>Angekündigte Meldungen</h2>
		<h2>Beendete Meldungen ↑</h2>
	</body></html>`

	statuses, err := parseTrackedServiceStatuses(strings.NewReader(outageHTML))
	require.NoError(t, err)

	sender := &recordingFCMSender{}
	require.NoError(t, syncAndNotifyServiceStatuses(ctx, db, sender, statuses))
	// Moodle (degraded) and eduroam (maintenance) transitioned from ok -> 2 push alerts.
	require.Len(t, sender.sent, 2)
	require.Equal(t, "moodle", sender.sent[0].Data["service"])
	require.Equal(t, model.ServiceStatusDegraded, sender.sent[0].Data["status"])
	require.Equal(t, "eduroam", sender.sent[1].Data["service"])
	require.Equal(t, model.ServiceStatusMaintenance, sender.sent[1].Data["status"])

	// Polling the same outage again does not duplicate alerts.
	sender.sent = nil
	require.NoError(t, syncAndNotifyServiceStatuses(ctx, db, sender, statuses))
	require.Empty(t, sender.sent)

	// Recovery back to ok sends recovery notifications.
	const allClearHTML = `
	<html><body>
		<h2>Aktuelle Meldungen</h2>
		<div><div class="table-headers"></div></div>
		<h2>Angekündigte Meldungen</h2>
		<h2>Beendete Meldungen ↑</h2>
	</body></html>`
	recovered, err := parseTrackedServiceStatuses(strings.NewReader(allClearHTML))
	require.NoError(t, err)
	require.NoError(t, syncAndNotifyServiceStatuses(ctx, db, sender, recovered))
	require.Len(t, sender.sent, 2)
	require.Equal(t, model.ServiceStatusOK, sender.sent[0].Data["status"])
	require.Contains(t, sender.sent[0].Title, "recovered")
}

func TestNotifyFavoriteDishAlertsWithDietaryFilterAndDedup(t *testing.T) {
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)

	canteen := model.Canteen{Name: "MENSA_GARCHING"}
	require.NoError(t, db.Create(&canteen).Error)
	veganCurry := model.Dish{
		Name: "Veganes Thai Curry", Type: "Vegan", CanteenID: canteen.CanteenID,
		Labels: null.StringFrom("VEGAN"),
	}
	meatCurry := model.Dish{
		Name: "Hähnchen Curry", Type: "Fleisch", CanteenID: canteen.CanteenID,
		Labels: null.StringFrom("MEAT"),
	}
	require.NoError(t, db.Create(&veganCurry).Error)
	require.NoError(t, db.Create(&meatCurry).Error)
	require.NoError(t, db.Create(&model.DishesOfTheWeek{DishID: veganCurry.Dish, Year: 2026, Week: 46, Day: 1}).Error)
	require.NoError(t, db.Create(&model.DishesOfTheWeek{DishID: meatCurry.Dish, Year: 2026, Week: 46, Day: 1}).Error)

	require.NoError(t, db.Create(&model.PushDevice{
		DeviceID: "dev-vegan", DeviceType: "android", FCMToken: "tok-vegan",
	}).Error)
	require.NoError(t, db.Create(&model.PushPreference{
		DeviceID:             "dev-vegan",
		CanteenIDs:           `["mensa-garching"]`,
		FavoriteDishes:       `["curry"]`,
		DietaryLabels:        `["VEGAN"]`,
		NotifyFavoriteDishes: true,
		Language:             "de",
	}).Error)

	sender := &recordingFCMSender{}
	require.NoError(t, notifyFavoriteDishAlerts(ctx, db, sender, 2026, 46))
	// Only the vegan curry matches both keyword "curry" and dietary label "VEGAN".
	require.Len(t, sender.sent, 1)
	require.Equal(t, "tok-vegan", sender.sent[0].Token)
	require.Equal(t, "Veganes Thai Curry", sender.sent[0].Data["dish"])

	// Running again for the same week is deduplicated via push_delivery_logs.
	sender.sent = nil
	require.NoError(t, notifyFavoriteDishAlerts(ctx, db, sender, 2026, 46))
	require.Empty(t, sender.sent)
}
