package migration

import (
	_ "embed"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

//go:embed static_data/source-schema.sql
var sourceSchema string

// migrate20200000000000
// adds the source shema
func (m TumDBMigrator) migrate20200000000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20200000000000",
		Migrate: func(tx *gorm.DB) error {
			tx = tx.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})
			for _, line := range strings.Split(sourceSchema, ";") {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				if err := tx.Exec(line).Error; err != nil {
					log.WithError(err).WithField("line", line).Error("failed to execute line")
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			tables := []string{"actions", "alarm_ban", "alarm_log", "barrierFree_moreInfo", "barrierFree_persons", "card_type", "chat_room", "crontab", "curricula", "dish", "dishflags", "dish2dishflags", "faculty", "feedback", "files", "kino", "lecture", "location", "member", "card", "card_box", "card_comment", "card_option", "chat_message", "chat_room2members", "devices", "device2stats", "members_card", "members_card_answer_history", "mensa", "dish2mensa", "mensaplan_mensa", "mensaprices", "migrations", "newsSource", "news", "news_alert", "notification_type", "notification", "notification_confirmation", "openinghours", "question", "question2answer", "question2faculty", "questionAnswers", "reports", "rights", "menu", "modules", "roles", "roles2rights", "roomfinder_building2area", "roomfinder_buildings", "roomfinder_buildings2gps", "roomfinder_buildings2maps", "roomfinder_maps", "roomfinder_rooms", "roomfinder_rooms2maps", "roomfinder_schedules", "sessions", "tag", "card2tag", "ticket_admin", "ticket_group", "event", "ticket_admin2group", "ticket_payment", "ticket_type", "ticket_history", "update_note", "users", "log", "recover", "users2info", "users2roles", "wifi_measurement"}
			for _, table := range tables {
				if err := tx.Migrator().DropTable(table); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
