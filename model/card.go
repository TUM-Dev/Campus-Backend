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


CREATE TABLE `card` (
  `card` int NOT NULL AUTO_INCREMENT,
  `member` int DEFAULT NULL,
  `lecture` int DEFAULT NULL,
  `card_type` int DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `front_text` varchar(2000) DEFAULT NULL,
  `front_image` varchar(2000) DEFAULT NULL,
  `back_text` varchar(2000) DEFAULT NULL,
  `back_image` varchar(2000) DEFAULT NULL,
  `can_shift` tinyint(1) DEFAULT '0',
  `created_at` date NOT NULL,
  `updated_at` date NOT NULL,
  `duplicate_card` int DEFAULT NULL,
  `aggr_rating` float DEFAULT '0',
  PRIMARY KEY (`card`),
  KEY `member` (`member`),
  KEY `lecture` (`lecture`),
  KEY `card_type` (`card_type`),
  KEY `duplicate_card` (`duplicate_card`),
  CONSTRAINT `card_ibfk_1` FOREIGN KEY (`member`) REFERENCES `member` (`member`) ON DELETE SET NULL,
  CONSTRAINT `card_ibfk_2` FOREIGN KEY (`lecture`) REFERENCES `lecture` (`lecture`) ON DELETE SET NULL,
  CONSTRAINT `card_ibfk_3` FOREIGN KEY (`card_type`) REFERENCES `card_type` (`card_type`) ON DELETE SET NULL,
  CONSTRAINT `card_ibfk_4` FOREIGN KEY (`duplicate_card`) REFERENCES `card` (`card`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "card": 54,    "lecture": 59,    "can_shift": 54,    "back_text": "ybIMMXmVmsjxrUSRbHrKdrwuM",    "back_image": "MjNsNHJtFstikeEocgHOGFrad",    "aggr_rating": 0.48434657,    "member": 82,    "card_type": 80,    "title": "rvbYkWBDxJwlBZmeGYUhbpEci",    "created_at": "2313-02-03T00:26:32.154957768+01:00",    "duplicate_card": 70,    "front_text": "RvUkamiJsYGftimTrxdgFBFNZ",    "front_image": "EdFWPrkTWUUfbmfODqBfpxQrm",    "updated_at": "2287-11-23T21:16:41.171077176+01:00"}



*/

// Card struct is a row record of the card table in the tca database
type Card struct {
	//[ 0] card                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Card int32 `gorm:"primary_key;AUTO_INCREMENT;column:card;type:int;" json:"card"`
	//[ 1] member                                         int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Member null.Int `gorm:"column:member;type:int;" json:"member"`
	//[ 2] lecture                                        int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Lecture null.Int `gorm:"column:lecture;type:int;" json:"lecture"`
	//[ 3] card_type                                      int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	CardType null.Int `gorm:"column:card_type;type:int;" json:"card_type"`
	//[ 4] title                                          varchar(255)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 255     default: []
	Title null.String `gorm:"column:title;type:varchar;size:255;" json:"title"`
	//[ 5] front_text                                     varchar(2000)        null: true   primary: false  isArray: false  auto: false  col: varchar         len: 2000    default: []
	FrontText null.String `gorm:"column:front_text;type:varchar;size:2000;" json:"front_text"`
	//[ 6] front_image                                    varchar(2000)        null: true   primary: false  isArray: false  auto: false  col: varchar         len: 2000    default: []
	FrontImage null.String `gorm:"column:front_image;type:varchar;size:2000;" json:"front_image"`
	//[ 7] back_text                                      varchar(2000)        null: true   primary: false  isArray: false  auto: false  col: varchar         len: 2000    default: []
	BackText null.String `gorm:"column:back_text;type:varchar;size:2000;" json:"back_text"`
	//[ 8] back_image                                     varchar(2000)        null: true   primary: false  isArray: false  auto: false  col: varchar         len: 2000    default: []
	BackImage null.String `gorm:"column:back_image;type:varchar;size:2000;" json:"back_image"`
	//[ 9] can_shift                                      tinyint              null: true   primary: false  isArray: false  auto: false  col: tinyint         len: -1      default: [0]
	CanShift null.Int `gorm:"column:can_shift;type:tinyint;default:0;" json:"can_shift"`
	//[10] created_at                                     date                 null: false  primary: false  isArray: false  auto: false  col: date            len: -1      default: []
	CreatedAt time.Time `gorm:"column:created_at;type:date;" json:"created_at"`
	//[11] updated_at                                     date                 null: false  primary: false  isArray: false  auto: false  col: date            len: -1      default: []
	UpdatedAt time.Time `gorm:"column:updated_at;type:date;" json:"updated_at"`
	//[12] duplicate_card                                 int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	DuplicateCard null.Int `gorm:"column:duplicate_card;type:int;" json:"duplicate_card"`
	//[13] aggr_rating                                    float                null: true   primary: false  isArray: false  auto: false  col: float           len: -1      default: [0]
	AggrRating null.Float `gorm:"column:aggr_rating;type:float;default:0;" json:"aggr_rating"`
}

var cardTableInfo = &TableInfo{
	Name: "card",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "card",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       true,
			IsAutoIncrement:    true,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "Card",
			GoFieldType:        "int32",
			JSONFieldName:      "card",
			ProtobufFieldName:  "card",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "member",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "Member",
			GoFieldType:        "null.Int",
			JSONFieldName:      "member",
			ProtobufFieldName:  "member",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "lecture",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "Lecture",
			GoFieldType:        "null.Int",
			JSONFieldName:      "lecture",
			ProtobufFieldName:  "lecture",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "card_type",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "CardType",
			GoFieldType:        "null.Int",
			JSONFieldName:      "card_type",
			ProtobufFieldName:  "card_type",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "title",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       255,
			GoFieldName:        "Title",
			GoFieldType:        "null.String",
			JSONFieldName:      "title",
			ProtobufFieldName:  "title",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "front_text",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(2000)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       2000,
			GoFieldName:        "FrontText",
			GoFieldType:        "null.String",
			JSONFieldName:      "front_text",
			ProtobufFieldName:  "front_text",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "front_image",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(2000)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       2000,
			GoFieldName:        "FrontImage",
			GoFieldType:        "null.String",
			JSONFieldName:      "front_image",
			ProtobufFieldName:  "front_image",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "back_text",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(2000)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       2000,
			GoFieldName:        "BackText",
			GoFieldType:        "null.String",
			JSONFieldName:      "back_text",
			ProtobufFieldName:  "back_text",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "back_image",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(2000)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       2000,
			GoFieldName:        "BackImage",
			GoFieldType:        "null.String",
			JSONFieldName:      "back_image",
			ProtobufFieldName:  "back_image",
			ProtobufType:       "string",
			ProtobufPos:        9,
		},

		&ColumnInfo{
			Index:              9,
			Name:               "can_shift",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "tinyint",
			DatabaseTypePretty: "tinyint",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "tinyint",
			ColumnLength:       -1,
			GoFieldName:        "CanShift",
			GoFieldType:        "null.Int",
			JSONFieldName:      "can_shift",
			ProtobufFieldName:  "can_shift",
			ProtobufType:       "int32",
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
			Name:               "created_at",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "date",
			DatabaseTypePretty: "date",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "date",
			ColumnLength:       -1,
			GoFieldName:        "CreatedAt",
			GoFieldType:        "time.Time",
			JSONFieldName:      "created_at",
			ProtobufFieldName:  "created_at",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        11,
		},

		&ColumnInfo{
			Index:              11,
			Name:               "updated_at",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "date",
			DatabaseTypePretty: "date",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "date",
			ColumnLength:       -1,
			GoFieldName:        "UpdatedAt",
			GoFieldType:        "time.Time",
			JSONFieldName:      "updated_at",
			ProtobufFieldName:  "updated_at",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        12,
		},

		&ColumnInfo{
			Index:              12,
			Name:               "duplicate_card",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "DuplicateCard",
			GoFieldType:        "null.Int",
			JSONFieldName:      "duplicate_card",
			ProtobufFieldName:  "duplicate_card",
			ProtobufType:       "int32",
			ProtobufPos:        13,
		},

		&ColumnInfo{
			Index:              13,
			Name:               "aggr_rating",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "float",
			DatabaseTypePretty: "float",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "float",
			ColumnLength:       -1,
			GoFieldName:        "AggrRating",
			GoFieldType:        "null.Float",
			JSONFieldName:      "aggr_rating",
			ProtobufFieldName:  "aggr_rating",
			ProtobufType:       "float",
			ProtobufPos:        14,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *Card) TableName() string {
	return "card"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *Card) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *Card) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *Card) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *Card) TableInfo() *TableInfo {
	return cardTableInfo
}
