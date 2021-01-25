package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

/*
DB Table Details
-------------------------------------


CREATE TABLE `device2stats` (
  `device` int NOT NULL DEFAULT '0',
  `LecturesPersonalActivity` int NOT NULL DEFAULT '0',
  `CafeteriaActivity` int NOT NULL DEFAULT '0',
  `WizNavStartActivity` int NOT NULL DEFAULT '0',
  `NewsActivity` int NOT NULL DEFAULT '0',
  `StartupActivity` int NOT NULL DEFAULT '0',
  `MainActivity` int NOT NULL DEFAULT '0',
  `ChatRoomsActivity` int NOT NULL DEFAULT '0',
  `CalendarActivity` int NOT NULL DEFAULT '0',
  `WizNavCheckTokenActivity` int NOT NULL DEFAULT '0',
  `ChatActivity` int NOT NULL DEFAULT '0',
  `CurriculaActivity` int NOT NULL DEFAULT '0',
  `CurriculaDetailsActivity` int NOT NULL DEFAULT '0',
  `GradeChartActivity` int NOT NULL DEFAULT '0',
  `GradesActivity` int NOT NULL DEFAULT '0',
  `InformationActivity` int NOT NULL DEFAULT '0',
  `LecturesAppointmentsActivity` int NOT NULL DEFAULT '0',
  `LecturesDetailsActivity` int NOT NULL DEFAULT '0',
  `OpeningHoursDetailActivity` int NOT NULL DEFAULT '0',
  `OpeningHoursListActivity` int NOT NULL DEFAULT '0',
  `OrganisationActivity` int NOT NULL DEFAULT '0',
  `OrganisationDetailsActivity` int NOT NULL DEFAULT '0',
  `PersonsDetailsActivity` int NOT NULL DEFAULT '0',
  `PersonsSearchActivity` int NOT NULL DEFAULT '0',
  `PlansActivity` int NOT NULL DEFAULT '0',
  `PlansDetailsActivity` int NOT NULL DEFAULT '0',
  `RoomFinderActivity` int NOT NULL DEFAULT '0',
  `RoomFinderDetailsActivity` int NOT NULL DEFAULT '0',
  `SetupEduroamActivity` int NOT NULL DEFAULT '0',
  `TransportationActivity` int NOT NULL DEFAULT '0',
  `TransportationDetailsActivity` int NOT NULL DEFAULT '0',
  `TuitionFeesActivity` int NOT NULL DEFAULT '0',
  `UserPreferencesActivity` int NOT NULL DEFAULT '0',
  `WizNavExtrasActivity` int NOT NULL DEFAULT '0',
  `WizNavChatActivity` int NOT NULL DEFAULT '0',
  `TuitionFeesCard` int NOT NULL DEFAULT '0',
  `NextLectureCard` int NOT NULL DEFAULT '0',
  `CafeteriaMenuCard` int NOT NULL DEFAULT '0',
  `NewsCard1` int NOT NULL DEFAULT '0',
  `NewsCard2` int NOT NULL DEFAULT '0',
  `NewsCard3` int NOT NULL DEFAULT '0',
  `NewsCard7` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`device`),
  CONSTRAINT `device2stats_ibfk_2` FOREIGN KEY (`device`) REFERENCES `devices` (`device`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "main_activity": 89,    "information_activity": 35,    "opening_hours_list_activity": 99,    "plans_details_activity": 51,    "room_finder_details_activity": 51,    "wiz_nav_chat_activity": 36,    "chat_activity": 72,    "curricula_details_activity": 42,    "grades_activity": 51,    "lectures_details_activity": 97,    "tuition_fees_activity": 92,    "news_card_3": 24,    "news_card_7": 69,    "lectures_personal_activity": 81,    "wiz_nav_start_activity": 79,    "opening_hours_detail_activity": 50,    "setup_eduroam_activity": 83,    "transportation_details_activity": 30,    "news_card_1": 37,    "room_finder_activity": 90,    "transportation_activity": 69,    "device": 84,    "cafeteria_activity": 38,    "startup_activity": 16,    "calendar_activity": 60,    "curricula_activity": 7,    "plans_activity": 8,    "chat_rooms_activity": 77,    "wiz_nav_check_token_activity": 43,    "organisation_details_activity": 96,    "user_preferences_activity": 80,    "grade_chart_activity": 1,    "lectures_appointments_activity": 68,    "persons_details_activity": 64,    "next_lecture_card": 58,    "news_card_2": 37,    "news_activity": 29,    "organisation_activity": 97,    "persons_search_activity": 67,    "wiz_nav_extras_activity": 74,    "tuition_fees_card": 76,    "cafeteria_menu_card": 28}



*/

// Device2stats struct is a row record of the device2stats table in the tca database
type Device2stats struct {
	//[ 0] device                                         int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: [0]
	Device int32 `gorm:"primary_key;column:device;type:int;default:0;" json:"device"`
	//[ 1] LecturesPersonalActivity                       int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	LecturesPersonalActivity int32 `gorm:"column:LecturesPersonalActivity;type:int;default:0;" json:"lectures_personal_activity"`
	//[ 2] CafeteriaActivity                              int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	CafeteriaActivity int32 `gorm:"column:CafeteriaActivity;type:int;default:0;" json:"cafeteria_activity"`
	//[ 3] WizNavStartActivity                            int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	WizNavStartActivity int32 `gorm:"column:WizNavStartActivity;type:int;default:0;" json:"wiz_nav_start_activity"`
	//[ 4] NewsActivity                                   int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	NewsActivity int32 `gorm:"column:NewsActivity;type:int;default:0;" json:"news_activity"`
	//[ 5] StartupActivity                                int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	StartupActivity int32 `gorm:"column:StartupActivity;type:int;default:0;" json:"startup_activity"`
	//[ 6] MainActivity                                   int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	MainActivity int32 `gorm:"column:MainActivity;type:int;default:0;" json:"main_activity"`
	//[ 7] ChatRoomsActivity                              int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	ChatRoomsActivity int32 `gorm:"column:ChatRoomsActivity;type:int;default:0;" json:"chat_rooms_activity"`
	//[ 8] CalendarActivity                               int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	CalendarActivity int32 `gorm:"column:CalendarActivity;type:int;default:0;" json:"calendar_activity"`
	//[ 9] WizNavCheckTokenActivity                       int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	WizNavCheckTokenActivity int32 `gorm:"column:WizNavCheckTokenActivity;type:int;default:0;" json:"wiz_nav_check_token_activity"`
	//[10] ChatActivity                                   int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	ChatActivity int32 `gorm:"column:ChatActivity;type:int;default:0;" json:"chat_activity"`
	//[11] CurriculaActivity                              int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	CurriculaActivity int32 `gorm:"column:CurriculaActivity;type:int;default:0;" json:"curricula_activity"`
	//[12] CurriculaDetailsActivity                       int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	CurriculaDetailsActivity int32 `gorm:"column:CurriculaDetailsActivity;type:int;default:0;" json:"curricula_details_activity"`
	//[13] GradeChartActivity                             int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	GradeChartActivity int32 `gorm:"column:GradeChartActivity;type:int;default:0;" json:"grade_chart_activity"`
	//[14] GradesActivity                                 int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	GradesActivity int32 `gorm:"column:GradesActivity;type:int;default:0;" json:"grades_activity"`
	//[15] InformationActivity                            int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	InformationActivity int32 `gorm:"column:InformationActivity;type:int;default:0;" json:"information_activity"`
	//[16] LecturesAppointmentsActivity                   int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	LecturesAppointmentsActivity int32 `gorm:"column:LecturesAppointmentsActivity;type:int;default:0;" json:"lectures_appointments_activity"`
	//[17] LecturesDetailsActivity                        int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	LecturesDetailsActivity int32 `gorm:"column:LecturesDetailsActivity;type:int;default:0;" json:"lectures_details_activity"`
	//[18] OpeningHoursDetailActivity                     int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	OpeningHoursDetailActivity int32 `gorm:"column:OpeningHoursDetailActivity;type:int;default:0;" json:"opening_hours_detail_activity"`
	//[19] OpeningHoursListActivity                       int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	OpeningHoursListActivity int32 `gorm:"column:OpeningHoursListActivity;type:int;default:0;" json:"opening_hours_list_activity"`
	//[20] OrganisationActivity                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	OrganisationActivity int32 `gorm:"column:OrganisationActivity;type:int;default:0;" json:"organisation_activity"`
	//[21] OrganisationDetailsActivity                    int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	OrganisationDetailsActivity int32 `gorm:"column:OrganisationDetailsActivity;type:int;default:0;" json:"organisation_details_activity"`
	//[22] PersonsDetailsActivity                         int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	PersonsDetailsActivity int32 `gorm:"column:PersonsDetailsActivity;type:int;default:0;" json:"persons_details_activity"`
	//[23] PersonsSearchActivity                          int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	PersonsSearchActivity int32 `gorm:"column:PersonsSearchActivity;type:int;default:0;" json:"persons_search_activity"`
	//[24] PlansActivity                                  int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	PlansActivity int32 `gorm:"column:PlansActivity;type:int;default:0;" json:"plans_activity"`
	//[25] PlansDetailsActivity                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	PlansDetailsActivity int32 `gorm:"column:PlansDetailsActivity;type:int;default:0;" json:"plans_details_activity"`
	//[26] RoomFinderActivity                             int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	RoomFinderActivity int32 `gorm:"column:RoomFinderActivity;type:int;default:0;" json:"room_finder_activity"`
	//[27] RoomFinderDetailsActivity                      int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	RoomFinderDetailsActivity int32 `gorm:"column:RoomFinderDetailsActivity;type:int;default:0;" json:"room_finder_details_activity"`
	//[28] SetupEduroamActivity                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	SetupEduroamActivity int32 `gorm:"column:SetupEduroamActivity;type:int;default:0;" json:"setup_eduroam_activity"`
	//[29] TransportationActivity                         int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	TransportationActivity int32 `gorm:"column:TransportationActivity;type:int;default:0;" json:"transportation_activity"`
	//[30] TransportationDetailsActivity                  int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	TransportationDetailsActivity int32 `gorm:"column:TransportationDetailsActivity;type:int;default:0;" json:"transportation_details_activity"`
	//[31] TuitionFeesActivity                            int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	TuitionFeesActivity int32 `gorm:"column:TuitionFeesActivity;type:int;default:0;" json:"tuition_fees_activity"`
	//[32] UserPreferencesActivity                        int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	UserPreferencesActivity int32 `gorm:"column:UserPreferencesActivity;type:int;default:0;" json:"user_preferences_activity"`
	//[33] WizNavExtrasActivity                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	WizNavExtrasActivity int32 `gorm:"column:WizNavExtrasActivity;type:int;default:0;" json:"wiz_nav_extras_activity"`
	//[34] WizNavChatActivity                             int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	WizNavChatActivity int32 `gorm:"column:WizNavChatActivity;type:int;default:0;" json:"wiz_nav_chat_activity"`
	//[35] TuitionFeesCard                                int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	TuitionFeesCard int32 `gorm:"column:TuitionFeesCard;type:int;default:0;" json:"tuition_fees_card"`
	//[36] NextLectureCard                                int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	NextLectureCard int32 `gorm:"column:NextLectureCard;type:int;default:0;" json:"next_lecture_card"`
	//[37] CafeteriaMenuCard                              int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	CafeteriaMenuCard int32 `gorm:"column:CafeteriaMenuCard;type:int;default:0;" json:"cafeteria_menu_card"`
	//[38] NewsCard1                                      int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	NewsCard1 int32 `gorm:"column:NewsCard1;type:int;default:0;" json:"news_card_1"`
	//[39] NewsCard2                                      int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	NewsCard2 int32 `gorm:"column:NewsCard2;type:int;default:0;" json:"news_card_2"`
	//[40] NewsCard3                                      int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	NewsCard3 int32 `gorm:"column:NewsCard3;type:int;default:0;" json:"news_card_3"`
	//[41] NewsCard7                                      int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	NewsCard7 int32 `gorm:"column:NewsCard7;type:int;default:0;" json:"news_card_7"`
}

var device2statsTableInfo = &TableInfo{
	Name: "device2stats",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "device",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "Device",
			GoFieldType:        "int32",
			JSONFieldName:      "device",
			ProtobufFieldName:  "device",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "LecturesPersonalActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "LecturesPersonalActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "lectures_personal_activity",
			ProtobufFieldName:  "lectures_personal_activity",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "CafeteriaActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "CafeteriaActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "cafeteria_activity",
			ProtobufFieldName:  "cafeteria_activity",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "WizNavStartActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "WizNavStartActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "wiz_nav_start_activity",
			ProtobufFieldName:  "wiz_nav_start_activity",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "NewsActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "NewsActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "news_activity",
			ProtobufFieldName:  "news_activity",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "StartupActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "StartupActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "startup_activity",
			ProtobufFieldName:  "startup_activity",
			ProtobufType:       "int32",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "MainActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "MainActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "main_activity",
			ProtobufFieldName:  "main_activity",
			ProtobufType:       "int32",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "ChatRoomsActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "ChatRoomsActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "chat_rooms_activity",
			ProtobufFieldName:  "chat_rooms_activity",
			ProtobufType:       "int32",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "CalendarActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "CalendarActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "calendar_activity",
			ProtobufFieldName:  "calendar_activity",
			ProtobufType:       "int32",
			ProtobufPos:        9,
		},

		&ColumnInfo{
			Index:              9,
			Name:               "WizNavCheckTokenActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "WizNavCheckTokenActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "wiz_nav_check_token_activity",
			ProtobufFieldName:  "wiz_nav_check_token_activity",
			ProtobufType:       "int32",
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
			Name:               "ChatActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "ChatActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "chat_activity",
			ProtobufFieldName:  "chat_activity",
			ProtobufType:       "int32",
			ProtobufPos:        11,
		},

		&ColumnInfo{
			Index:              11,
			Name:               "CurriculaActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "CurriculaActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "curricula_activity",
			ProtobufFieldName:  "curricula_activity",
			ProtobufType:       "int32",
			ProtobufPos:        12,
		},

		&ColumnInfo{
			Index:              12,
			Name:               "CurriculaDetailsActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "CurriculaDetailsActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "curricula_details_activity",
			ProtobufFieldName:  "curricula_details_activity",
			ProtobufType:       "int32",
			ProtobufPos:        13,
		},

		&ColumnInfo{
			Index:              13,
			Name:               "GradeChartActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "GradeChartActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "grade_chart_activity",
			ProtobufFieldName:  "grade_chart_activity",
			ProtobufType:       "int32",
			ProtobufPos:        14,
		},

		&ColumnInfo{
			Index:              14,
			Name:               "GradesActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "GradesActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "grades_activity",
			ProtobufFieldName:  "grades_activity",
			ProtobufType:       "int32",
			ProtobufPos:        15,
		},

		&ColumnInfo{
			Index:              15,
			Name:               "InformationActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "InformationActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "information_activity",
			ProtobufFieldName:  "information_activity",
			ProtobufType:       "int32",
			ProtobufPos:        16,
		},

		&ColumnInfo{
			Index:              16,
			Name:               "LecturesAppointmentsActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "LecturesAppointmentsActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "lectures_appointments_activity",
			ProtobufFieldName:  "lectures_appointments_activity",
			ProtobufType:       "int32",
			ProtobufPos:        17,
		},

		&ColumnInfo{
			Index:              17,
			Name:               "LecturesDetailsActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "LecturesDetailsActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "lectures_details_activity",
			ProtobufFieldName:  "lectures_details_activity",
			ProtobufType:       "int32",
			ProtobufPos:        18,
		},

		&ColumnInfo{
			Index:              18,
			Name:               "OpeningHoursDetailActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "OpeningHoursDetailActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "opening_hours_detail_activity",
			ProtobufFieldName:  "opening_hours_detail_activity",
			ProtobufType:       "int32",
			ProtobufPos:        19,
		},

		&ColumnInfo{
			Index:              19,
			Name:               "OpeningHoursListActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "OpeningHoursListActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "opening_hours_list_activity",
			ProtobufFieldName:  "opening_hours_list_activity",
			ProtobufType:       "int32",
			ProtobufPos:        20,
		},

		&ColumnInfo{
			Index:              20,
			Name:               "OrganisationActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "OrganisationActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "organisation_activity",
			ProtobufFieldName:  "organisation_activity",
			ProtobufType:       "int32",
			ProtobufPos:        21,
		},

		&ColumnInfo{
			Index:              21,
			Name:               "OrganisationDetailsActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "OrganisationDetailsActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "organisation_details_activity",
			ProtobufFieldName:  "organisation_details_activity",
			ProtobufType:       "int32",
			ProtobufPos:        22,
		},

		&ColumnInfo{
			Index:              22,
			Name:               "PersonsDetailsActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "PersonsDetailsActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "persons_details_activity",
			ProtobufFieldName:  "persons_details_activity",
			ProtobufType:       "int32",
			ProtobufPos:        23,
		},

		&ColumnInfo{
			Index:              23,
			Name:               "PersonsSearchActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "PersonsSearchActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "persons_search_activity",
			ProtobufFieldName:  "persons_search_activity",
			ProtobufType:       "int32",
			ProtobufPos:        24,
		},

		&ColumnInfo{
			Index:              24,
			Name:               "PlansActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "PlansActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "plans_activity",
			ProtobufFieldName:  "plans_activity",
			ProtobufType:       "int32",
			ProtobufPos:        25,
		},

		&ColumnInfo{
			Index:              25,
			Name:               "PlansDetailsActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "PlansDetailsActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "plans_details_activity",
			ProtobufFieldName:  "plans_details_activity",
			ProtobufType:       "int32",
			ProtobufPos:        26,
		},

		&ColumnInfo{
			Index:              26,
			Name:               "RoomFinderActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "RoomFinderActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "room_finder_activity",
			ProtobufFieldName:  "room_finder_activity",
			ProtobufType:       "int32",
			ProtobufPos:        27,
		},

		&ColumnInfo{
			Index:              27,
			Name:               "RoomFinderDetailsActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "RoomFinderDetailsActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "room_finder_details_activity",
			ProtobufFieldName:  "room_finder_details_activity",
			ProtobufType:       "int32",
			ProtobufPos:        28,
		},

		&ColumnInfo{
			Index:              28,
			Name:               "SetupEduroamActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "SetupEduroamActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "setup_eduroam_activity",
			ProtobufFieldName:  "setup_eduroam_activity",
			ProtobufType:       "int32",
			ProtobufPos:        29,
		},

		&ColumnInfo{
			Index:              29,
			Name:               "TransportationActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "TransportationActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "transportation_activity",
			ProtobufFieldName:  "transportation_activity",
			ProtobufType:       "int32",
			ProtobufPos:        30,
		},

		&ColumnInfo{
			Index:              30,
			Name:               "TransportationDetailsActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "TransportationDetailsActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "transportation_details_activity",
			ProtobufFieldName:  "transportation_details_activity",
			ProtobufType:       "int32",
			ProtobufPos:        31,
		},

		&ColumnInfo{
			Index:              31,
			Name:               "TuitionFeesActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "TuitionFeesActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "tuition_fees_activity",
			ProtobufFieldName:  "tuition_fees_activity",
			ProtobufType:       "int32",
			ProtobufPos:        32,
		},

		&ColumnInfo{
			Index:              32,
			Name:               "UserPreferencesActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "UserPreferencesActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "user_preferences_activity",
			ProtobufFieldName:  "user_preferences_activity",
			ProtobufType:       "int32",
			ProtobufPos:        33,
		},

		&ColumnInfo{
			Index:              33,
			Name:               "WizNavExtrasActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "WizNavExtrasActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "wiz_nav_extras_activity",
			ProtobufFieldName:  "wiz_nav_extras_activity",
			ProtobufType:       "int32",
			ProtobufPos:        34,
		},

		&ColumnInfo{
			Index:              34,
			Name:               "WizNavChatActivity",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "WizNavChatActivity",
			GoFieldType:        "int32",
			JSONFieldName:      "wiz_nav_chat_activity",
			ProtobufFieldName:  "wiz_nav_chat_activity",
			ProtobufType:       "int32",
			ProtobufPos:        35,
		},

		&ColumnInfo{
			Index:              35,
			Name:               "TuitionFeesCard",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "TuitionFeesCard",
			GoFieldType:        "int32",
			JSONFieldName:      "tuition_fees_card",
			ProtobufFieldName:  "tuition_fees_card",
			ProtobufType:       "int32",
			ProtobufPos:        36,
		},

		&ColumnInfo{
			Index:              36,
			Name:               "NextLectureCard",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "NextLectureCard",
			GoFieldType:        "int32",
			JSONFieldName:      "next_lecture_card",
			ProtobufFieldName:  "next_lecture_card",
			ProtobufType:       "int32",
			ProtobufPos:        37,
		},

		&ColumnInfo{
			Index:              37,
			Name:               "CafeteriaMenuCard",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "CafeteriaMenuCard",
			GoFieldType:        "int32",
			JSONFieldName:      "cafeteria_menu_card",
			ProtobufFieldName:  "cafeteria_menu_card",
			ProtobufType:       "int32",
			ProtobufPos:        38,
		},

		&ColumnInfo{
			Index:              38,
			Name:               "NewsCard1",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "NewsCard1",
			GoFieldType:        "int32",
			JSONFieldName:      "news_card_1",
			ProtobufFieldName:  "news_card_1",
			ProtobufType:       "int32",
			ProtobufPos:        39,
		},

		&ColumnInfo{
			Index:              39,
			Name:               "NewsCard2",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "NewsCard2",
			GoFieldType:        "int32",
			JSONFieldName:      "news_card_2",
			ProtobufFieldName:  "news_card_2",
			ProtobufType:       "int32",
			ProtobufPos:        40,
		},

		&ColumnInfo{
			Index:              40,
			Name:               "NewsCard3",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "NewsCard3",
			GoFieldType:        "int32",
			JSONFieldName:      "news_card_3",
			ProtobufFieldName:  "news_card_3",
			ProtobufType:       "int32",
			ProtobufPos:        41,
		},

		&ColumnInfo{
			Index:              41,
			Name:               "NewsCard7",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "NewsCard7",
			GoFieldType:        "int32",
			JSONFieldName:      "news_card_7",
			ProtobufFieldName:  "news_card_7",
			ProtobufType:       "int32",
			ProtobufPos:        42,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *Device2stats) TableName() string {
	return "device2stats"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *Device2stats) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *Device2stats) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *Device2stats) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *Device2stats) TableInfo() *TableInfo {
	return device2statsTableInfo
}
