package migration

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type tableWithWrongfield struct {
	table string
	field string
}
type fkNeedingMigration struct {
	fromTable      string
	fromColumn     string
	nullability    string
	toTable        string
	toColumn       string
	constraintName string
}

func tablesWithWrongId() []tableWithWrongfield {
	return []tableWithWrongfield{
		// PKs which are int32
		{"alarm_ban", "ban"},
		{"alarm_log", "alarm"},
		{"barrierFree_moreInfo", "id"},
		{"barrierFree_persons", "id"},
		{"chat_message", "message"},
		{"chat_room2members", "room2members"},
		{"crontab", "cron"},
		{"curricula", "curriculum"},
		{"dish2dishflags", "dish2dishflags"},
		{"dish2mensa", "dish2mensa"},
		{"feedback", "id"},
		{"log", "log"},
		{"mensaplan_mensa", "id"},
		{"mensaprices", "price"},
		{"modules", "module"},
		{"news_alert", "news_alert"},
		{"openinghours", "id"},
		{"recover", "recover"},
		{"reports", "report"},
		{"ticket_admin2group", "ticket_admin2group"},
		{"ticket_history", "ticket_history"},
		{"wifi_measurement", "id"},
		{"update_note", "version_code"},
		// should have been a fk, but is an index
		{"question", "member"},
		// multi-pk-indexes
		{"roomfinder_buildings2maps", "map_id"},
		{"roomfinder_maps", "map_id"},
		{"roomfinder_rooms", "room_id"},
		{"roomfinder_rooms2maps", "room_id"},
		{"roomfinder_rooms2maps", "map_id"},
		{"roomfinder_schedules", "room_id"},
	}
}

// tablesWithWrongFk tells a user which FKs exist that have a int-ish type
// can be generated via:
// ```sql
// with fks as (select fks.table_name            as from_table,
//
//	               group_concat(kcu.COLUMN_NAME
//	                            order by position_in_unique_constraint separator ', ')
//	                                         as from_columns,
//	               fks.referenced_table_name as to_table,
//	               group_concat(kcu.REFERENCED_COLUMN_NAME
//	                            order by position_in_unique_constraint separator ', ')
//	                                         as to_columns,
//	               fks.constraint_name
//	        from information_schema.referential_constraints fks
//	                 join information_schema.key_column_usage kcu
//	                      on fks.constraint_schema = kcu.table_schema
//	                          and fks.table_name = kcu.table_name
//	                          and fks.constraint_name = kcu.constraint_name
//	        where fks.constraint_schema = 'campus_db'
//	        group by fks.constraint_schema,
//	                 fks.table_name,
//	                 fks.unique_constraint_schema,
//	                 fks.referenced_table_name,
//	                 fks.constraint_name
//	        order by fks.constraint_schema,
//	                 fks.table_name),
//	tables_with_matching_type as (SELECT TABLE_NAME, COLUMN_NAME
//	                              from information_schema.columns
//	                              WHERE DATA_TYPE like '%int%'
//	                                and TABLE_SCHEMA = 'campus_db')
//
// SELECT f.*
// from fks f
// WHERE EXISTS(SELECT *
//
//	FROM tables_with_matching_type t
//	where t.TABLE_NAME = f.from_table
//	  and t.COLUMN_NAME = f.from_columns);
func tablesWithWrongFk() []fkNeedingMigration {
	return []fkNeedingMigration{
		{"chat_message", "member", "NOT NULL", "member", "member", "chat_message_ibfk_1"},
		{"chat_message", "room", "NOT NULL", "chat_room", "room", "FK_chat_message_chat_room"},
		{"chat_room2members", "room", "NOT NULL", "chat_room", "room", "FK_chat_room2members_chat_room"},
		{"chat_room2members", "member", "NOT NULL", "member", "member", "chat_room2members_ibfk_2"},
		{"device2stats", "device", "NOT NULL", "devices", "device", "device2stats_ibfk_2"},
		{"devices", "member", "NULL", "member", "member", "devices_ibfk_1"},
		{"dish2dishflags", "dish", "NOT NULL", "dish", "dish", "dish2dishflags_ibfk_1"},
		{"dish2dishflags", "flag", "NOT NULL", "dishflags", "flag", "dish2dishflags_ibfk_2"},
		{"dish2mensa", "mensa", "NOT NULL", "mensa", "mensa", "dish2mensa_ibfk_1"},
		{"dish2mensa", "dish", "NOT NULL", "dish", "dish", "dish2mensa_ibfk_2"},
		{"dish_rating", "dishID", "NOT NULL", "dish", "dish", "dish_rating_dish_dish_fk"},
		{"event", "news", "NULL", "news", "news", "fkNews"},
		{"event", "kino", "NULL", "kino", "kino", "fkKino"},
		{"event", "file", "NULL", "files", "file", "fkEventFile"},
		{"event", "ticket_group", "NULL", "ticket_group", "ticket_group", "fkEventGroup"},
		{"kino", "cover", "NULL", "files", "file", "kino_ibfk_1"},
		{"log", "user_executed", "NULL", "users", "user", "fkLog2UsersEx"},
		{"log", "user_affected", "NULL", "users", "user", "fkLog2UsersAf"},
		{"log", "action", "NULL", "actions", "action", "fkLog2Actions"},
		{"menu", "right", "NULL", "rights", "right", "menu_ibfk_1"},
		{"menu", "parent", "NULL", "menu", "menu", "menu_ibfk_2"},
		{"modules", "right", "NULL", "rights", "right", "fkMod2Rights"},
		{"news", "src", "NOT NULL", "newsSource", "source", "news_ibfk_1"},
		{"news", "file", "NULL", "files", "file", "news_ibfk_2"},
		{"newsSource", "icon", "NOT NULL", "files", "file", "newsSource_ibfk_1"},
		{"notification", "type", "NOT NULL", "notification_type", "type", "notification_ibfk_1"},
		{"notification", "location", "NULL", "location", "location", "notification_ibfk_2"},
		{"notification_confirmation", "notification", "NOT NULL", "notification", "notification", "notification_confirmation_ibfk_1"},
		{"notification_confirmation", "device", "NOT NULL", "devices", "device", "notification_confirmation_ibfk_2"},
		{"question2answer", "question", "NOT NULL", "question", "question", "question2answer_question_question_fk"},
		{"question2answer", "answer", "NOT NULL", "questionAnswers", "answer", "question2answer_questionAnswers_answer_fk"},
		{"question2answer", "member", "NOT NULL", "member", "member", "question2answer_member_member_fk"},
		{"question2faculty", "question", "NOT NULL", "question", "question", "question2faculty_ibfk_1"},
		{"question2faculty", "faculty", "NOT NULL", "faculty", "faculty", "question2faculty_ibfk_2"},
		{"recover", "user", "NOT NULL", "users", "user", "fkRecover2User"},
		{"reports", "device", "NULL", "devices", "device", "reports_ibfk_3"},
		{"roles2rights", "role", "NOT NULL", "roles", "role", "fkRole"},
		{"roles2rights", "right", "NOT NULL", "rights", "right", "fkRight"},
		{"ticket_admin2group", "ticket_admin", "NOT NULL", "ticket_admin", "ticket_admin", "fkTicketAdmin"},
		{"ticket_admin2group", "ticket_group", "NOT NULL", "ticket_group", "ticket_group", "fkTicketGroup"},
		{"ticket_history", "member", "NOT NULL", "member", "member", "fkMember"},
		{"ticket_history", "ticket_payment", "NULL", "ticket_payment", "ticket_payment", "fkTicketPayment"},
		{"ticket_history", "ticket_type", "NOT NULL", "ticket_type", "ticket_type", "fkTicketType"},
		{"ticket_type", "event", "NOT NULL", "event", "event", "fkEvent"},
		{"ticket_type", "ticket_payment", "NOT NULL", "ticket_payment", "ticket_payment", "fkPayment"},
		{"users2info", "user", "NOT NULL", "users", "user", "fkUsers"},
		{"users2roles", "user", "NOT NULL", "users", "user", "fkUser2RolesUser"},
		{"users2roles", "role", "NOT NULL", "roles", "role", "fkUser2RolesRole"},
	}
}

func migrateField(tx *gorm.DB, table string, field string, typeDefiniton string) error {
	// change both the origin of the fk and the destination to be a bigint
	if err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s` MODIFY `%s` %s", table, field, typeDefiniton)).Error; err != nil {
		return err
	}
	// data is still stored as int32, but we can change this
	if err := tx.Exec(fmt.Sprintf("UPDATE `%s` SET `%s` = CAST(`%s` AS UNSIGNED INTEGER)", table, field, field)).Error; err != nil {
		return err
	}
	return nil
}

// migrate20240316000000
// made sure that all ids are int64
func migrate20240316000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240316000000",
		Migrate: func(tx *gorm.DB) error {
			log.Info("started migrating all keys to be i64 based by removing the FK")
			for _, f := range tablesWithWrongFk() {
				if err := tx.Exec(fmt.Sprintf("alter table `%s` DROP FOREIGN KEY `%s`", f.fromTable, f.constraintName)).Error; err != nil {
					return err
				}
			}
			log.Info("changing both the source and destination collumn to the same type")
			for _, f := range tablesWithWrongFk() {
				if err := migrateField(tx, f.fromTable, f.fromColumn, "BIGINT "+f.nullability); err != nil {
					return err
				}
				// mysql does not allow primary keys to have anything other than not null
				if err := migrateField(tx, f.toTable, f.toColumn, "BIGINT NOT NULL autoIncrement"); err != nil {
					return err
				}
			}
			log.Info("re-adding the FK")
			for _, f := range tablesWithWrongFk() {
				if err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s` ADD CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s` (`%s`)", f.fromTable, f.constraintName, f.fromColumn, f.toTable, f.toColumn)).Error; err != nil {
					return err
				}
			}
			log.Info("migrated added all FK-relationships to be i64 based")
			// because we have migrated all fk relationships, this does not mean that we have migrated all primary keys => this is done this way
			for _, t := range tablesWithWrongId() {
				if err := migrateField(tx, t.table, t.field, "BIGINT NOT NULL"); err != nil {
					return err
				}
				log.WithField("table", t.table).Info("migrated PK-field")
			}
			return nil
		},
		Rollback: func(_ *gorm.DB) error {
			log.Fatal("intentionally no rollback function as this would be lossy!")
			return nil
		},
	}
}
