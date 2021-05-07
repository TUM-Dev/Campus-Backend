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


CREATE TABLE `notification` (
  `notification` int NOT NULL AUTO_INCREMENT,
  `type` int NOT NULL,
  `location` int DEFAULT NULL,
  `title` text NOT NULL,
  `description` text NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `signature` text,
  `silent` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`notification`),
  KEY `type` (`type`),
  KEY `location` (`location`),
  CONSTRAINT `notification_ibfk_1` FOREIGN KEY (`type`) REFERENCES `notification_type` (`type`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `notification_ibfk_2` FOREIGN KEY (`location`) REFERENCES `location` (`location`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=104 DEFAULT CHARSET=utf8

JSON Sample
-------------------------------------
{    "notification": 83,    "type": 86,    "location": 80,    "title": "eeFZRyYFdrOGdBWtiqBMTyFAL",    "description": "qxcTwnwojWpCjEgSuHuWKuAOi",    "created": "2242-06-24T15:42:16.206870536+01:00",    "signature": "oPGvTLfahDoWYSaOZRxoVORRU",    "silent": 88}



*/

// Notification struct is a row record of the notification table in the tca database
type Notification struct {
	//[ 0] notification                                   int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Notification int32 `gorm:"primary_key;AUTO_INCREMENT;column:notification;type:int;" json:"notification"`
	//[ 1] type                                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Type int32 `gorm:"column:type;type:int;" json:"type"`
	//[ 2] location                                       int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Location null.Int `gorm:"column:location;type:int;" json:"location"`
	//[ 3] title                                          text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Title string `gorm:"column:title;type:text;size:65535;" json:"title"`
	//[ 4] description                                    text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Description string `gorm:"column:description;type:text;size:65535;" json:"description"`
	//[ 5] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 6] signature                                      text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Signature null.String `gorm:"column:signature;type:text;size:65535;" json:"signature"`
	//[ 7] silent                                         tinyint              null: false  primary: false  isArray: false  auto: false  col: tinyint         len: -1      default: [0]
	Silent int32 `gorm:"column:silent;type:tinyint;default:0;" json:"silent"`
}
