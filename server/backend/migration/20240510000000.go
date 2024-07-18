package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type device2stats struct {
	ChatRoomsActivity  int `gorm:"column:ChatRoomsActivity;default 0;not null"`
	ChatActivity       int `gorm:"column:ChatActivity;default 0;not null"`
	WizNavChatActivity int `gorm:"column:WizNavChatActivity;default 0;not null"`
}

// migrate20240510000000
// - Removes all traces of the chat from the database
func migrate20240510000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240510000000",
		Migrate: func(tx *gorm.DB) error {
			// Remove tracking from device2stats
			if err := tx.Migrator().DropColumn(&device2stats{}, "ChatRoomsActivity"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&device2stats{}, "ChatActivity"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&device2stats{}, "WizNavChatActivity"); err != nil {
				return err
			}

			// Delete all tables
			if err := tx.Migrator().DropTable("chat_message"); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable("chat_room2members"); err != nil {
				return err
			}
			return tx.Migrator().DropTable("chat_room")
		},
		Rollback: func(tx *gorm.DB) error {
			// Restore chat_room
			if err := tx.Exec("create table chat_room(room int auto_increment primary key, name varchar(100) not null, semester varchar(3) null, constraint `Index 2` unique (semester, name));").Error; err != nil {
				return err
			}

			// Add tracking from device2stats
			if err := tx.Migrator().AutoMigrate(&device2stats{}); err != nil {
				return err
			}

			// Restore chat_message
			if err := tx.Exec("create table chat_message (message int auto_increment primary key, member int not null, room int not null, text longtext not null, created datetime not null, signature longtext not null, constraint FK_chat_message_chat_room foreign key (room) references chat_room (room) on update cascade on delete cascade, constraint chat_message_ibfk_1 foreign key (member) references member (member) on update cascade on delete cascade);").Error; err != nil {
				return err
			}
			if err := tx.Exec("create index chat_message_b3c09425 on chat_message (member);").Error; err != nil {
				return err
			}
			if err := tx.Exec("create index chat_message_ca20ebca on chat_message (room);").Error; err != nil {
				return err
			}

			// Restore chat_room2members
			if err := tx.Exec("create table chat_room2members(room2members int auto_increment primary key, room int not null, member int not null, constraint chatroom_id unique (room, member), constraint FK_chat_room2members_chat_room foreign key (room) references chat_room (room) on update cascade on delete cascade, constraint chat_room2members_ibfk_2 foreign key (member) references member (member) on update cascade on delete cascade );").Error; err != nil {
				return err
			}
			if err := tx.Exec("create index chat_chatroom_members_29801a33 on chat_room2members (room);").Error; err != nil {
				return err
			}
			return tx.Exec("create index chat_chatroom_members_b3c09425 on chat_room2members (member);").Error
		},
	}
}
