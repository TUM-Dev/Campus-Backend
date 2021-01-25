package model

import "fmt"

// Action CRUD actions
type Action int32

var (
	// Create action when record is created
	Create = Action(0)

	// RetrieveOne action when a record is retrieved from db
	RetrieveOne = Action(1)

	// RetrieveMany action when record(s) are retrieved from db
	RetrieveMany = Action(2)

	// Update action when record is updated in db
	Update = Action(3)

	// Delete action when record is deleted in db
	Delete = Action(4)

	// FetchDDL action when fetching ddl info from db
	FetchDDL = Action(5)

	tables map[string]*TableInfo
)

func init() {
	tables = make(map[string]*TableInfo)

	tables["actions"] = actionsTableInfo
	tables["alarm_ban"] = alarm_banTableInfo
	tables["alarm_log"] = alarm_logTableInfo
	tables["barrierFree_moreInfo"] = barrierFree_moreInfoTableInfo
	tables["barrierFree_persons"] = barrierFree_personsTableInfo
	tables["card"] = cardTableInfo
	tables["card2tag"] = card2tagTableInfo
	tables["card_box"] = card_boxTableInfo
	tables["card_comment"] = card_commentTableInfo
	tables["card_option"] = card_optionTableInfo
	tables["card_type"] = card_typeTableInfo
	tables["chat_message"] = chat_messageTableInfo
	tables["chat_room"] = chat_roomTableInfo
	tables["chat_room2members"] = chat_room2membersTableInfo
	tables["crontab"] = crontabTableInfo
	tables["curricula"] = curriculaTableInfo
	tables["device2stats"] = device2statsTableInfo
	tables["devices"] = devicesTableInfo
	tables["dish"] = dishTableInfo
	tables["dish2dishflags"] = dish2dishflagsTableInfo
	tables["dish2mensa"] = dish2mensaTableInfo
	tables["dishflags"] = dishflagsTableInfo
	tables["event"] = eventTableInfo
	tables["faculty"] = facultyTableInfo
	tables["feedback"] = feedbackTableInfo
	tables["files"] = filesTableInfo
	tables["kino"] = kinoTableInfo
	tables["lecture"] = lectureTableInfo
	tables["location"] = locationTableInfo
	tables["log"] = logTableInfo
	tables["member"] = memberTableInfo
	tables["members_card"] = members_cardTableInfo
	tables["members_card_answer_history"] = members_card_answer_historyTableInfo
	tables["mensa"] = mensaTableInfo
	tables["mensaplan_mensa"] = mensaplan_mensaTableInfo
	tables["mensaprices"] = mensapricesTableInfo
	tables["menu"] = menuTableInfo
	tables["modules"] = modulesTableInfo
	tables["news"] = newsTableInfo
	tables["newsSource"] = newsSourceTableInfo
	tables["news_alert"] = news_alertTableInfo
	tables["notification"] = notificationTableInfo
	tables["notification_confirmation"] = notification_confirmationTableInfo
	tables["notification_type"] = notification_typeTableInfo
	tables["openinghours"] = openinghoursTableInfo
	tables["question"] = questionTableInfo
	tables["question2answer"] = question2answerTableInfo
	tables["question2faculty"] = question2facultyTableInfo
	tables["questionAnswers"] = questionAnswersTableInfo
	tables["recover"] = recoverTableInfo
	tables["reports"] = reportsTableInfo
	tables["rights"] = rightsTableInfo
	tables["roles"] = rolesTableInfo
	tables["roles2rights"] = roles2rightsTableInfo
	tables["roomfinder_building2area"] = roomfinder_building2areaTableInfo
	tables["roomfinder_buildings"] = roomfinder_buildingsTableInfo
	tables["roomfinder_buildings2gps"] = roomfinder_buildings2gpsTableInfo
	tables["roomfinder_buildings2maps"] = roomfinder_buildings2mapsTableInfo
	tables["roomfinder_maps"] = roomfinder_mapsTableInfo
	tables["roomfinder_rooms"] = roomfinder_roomsTableInfo
	tables["roomfinder_rooms2maps"] = roomfinder_rooms2mapsTableInfo
	tables["roomfinder_schedules"] = roomfinder_schedulesTableInfo
	tables["sessions"] = sessionsTableInfo
	tables["tag"] = tagTableInfo
	tables["ticket_admin"] = ticket_adminTableInfo
	tables["ticket_admin2group"] = ticket_admin2groupTableInfo
	tables["ticket_group"] = ticket_groupTableInfo
	tables["ticket_history"] = ticket_historyTableInfo
	tables["ticket_payment"] = ticket_paymentTableInfo
	tables["ticket_type"] = ticket_typeTableInfo
	tables["update_note"] = update_noteTableInfo
	tables["users"] = usersTableInfo
	tables["users2info"] = users2infoTableInfo
	tables["users2roles"] = users2rolesTableInfo
	tables["wifi_measurement"] = wifi_measurementTableInfo
}

// String describe the action
func (i Action) String() string {
	switch i {
	case Create:
		return "Create"
	case RetrieveOne:
		return "RetrieveOne"
	case RetrieveMany:
		return "RetrieveMany"
	case Update:
		return "Update"
	case Delete:
		return "Delete"
	case FetchDDL:
		return "FetchDDL"
	default:
		return fmt.Sprintf("unknown action: %d", int(i))
	}
}

// Model interface methods for database structs generated
type Model interface {
	TableName() string
	BeforeSave() error
	Prepare()
	Validate(action Action) error
	TableInfo() *TableInfo
}

// TableInfo describes a table in the database
type TableInfo struct {
	Name    string        `json:"name"`
	Columns []*ColumnInfo `json:"columns"`
}

// ColumnInfo describes a column in the database table
type ColumnInfo struct {
	Index              int    `json:"index"`
	GoFieldName        string `json:"go_field_name"`
	GoFieldType        string `json:"go_field_type"`
	JSONFieldName      string `json:"json_field_name"`
	ProtobufFieldName  string `json:"protobuf_field_name"`
	ProtobufType       string `json:"protobuf_field_type"`
	ProtobufPos        int    `json:"protobuf_field_pos"`
	Comment            string `json:"comment"`
	Notes              string `json:"notes"`
	Name               string `json:"name"`
	Nullable           bool   `json:"is_nullable"`
	DatabaseTypeName   string `json:"database_type_name"`
	DatabaseTypePretty string `json:"database_type_pretty"`
	IsPrimaryKey       bool   `json:"is_primary_key"`
	IsAutoIncrement    bool   `json:"is_auto_increment"`
	IsArray            bool   `json:"is_array"`
	ColumnType         string `json:"column_type"`
	ColumnLength       int64  `json:"column_length"`
	DefaultValue       string `json:"default_value"`
}

// GetTableInfo retrieve TableInfo for a table
func GetTableInfo(name string) (*TableInfo, bool) {
	val, ok := tables[name]
	return val, ok
}
