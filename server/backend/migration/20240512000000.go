package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20240512000000
// Removes ticketsales from the db
func migrate20240512000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240512000000",
		Migrate: func(tx *gorm.DB) error {
			// order intentional to avoid foreign key constraint errors
			for _, tbl := range []string{"ticket_history", "ticket_type", "ticket_payment", "ticket_admin2group", "ticket_admin", "event", "ticket_group"} {
				if err := tx.Migrator().DropTable(tbl); err != nil {
					return err
				}
			}
			return nil
		},

		Rollback: func(tx *gorm.DB) error {
			//ticket_group
			if err := tx.Exec(`create table ticket_group (ticket_group int auto_increment primary key,description  text not null);`).Error; err != nil {
				return err
			}
			//event
			if err := tx.Exec("create table event(event int auto_increment primary key, news int null, kino int null, file int null, title varchar(100) not null, description text not null, locality varchar(200) not null, link varchar(200) null, start datetime null, end datetime null, ticket_group int default 1 null, constraint fkEventFile foreign key (file) references files (file) on update cascade on delete set null, constraint fkEventGroup foreign key (ticket_group) references ticket_group (ticket_group), constraint fkKino foreign key (kino) references kino (kino) on update cascade on delete set null, constraint fkNews foreign key (news) references news (news) on update cascade on delete set null);").Error; err != nil {
				return err
			}
			if err := tx.Exec("create index file on event (file);").Error; err != nil {
				return err
			}
			if err := tx.Exec("create fulltext index searchTitle on event (title);").Error; err != nil {
				return err
			}
			//ticket_admin
			if err := tx.Exec("create table ticket_admin(ticket_admin int auto_increment primary key, `key` text not null, created timestamp default current_timestamp() not null, active tinyint(1) default 0 not null, comment text null);").Error; err != nil {
				return err
			}
			//ticket_admin2group
			if err := tx.Exec("create table ticket_admin2group( ticket_admin2group int auto_increment primary key, ticket_admin int not null, ticket_group int not null, constraint fkTicketAdmin foreign key (ticket_admin) references ticket_admin (ticket_admin) on update cascade on delete cascade, constraint fkTicketGroup foreign key (ticket_group) references ticket_group (ticket_group) on update cascade on delete cascade);").Error; err != nil {
				return err
			}
			if err := tx.Exec("create index ticket_admin on ticket_admin2group (ticket_admin);").Error; err != nil {
				return err
			}
			if err := tx.Exec("create index ticket_group on ticket_admin2group (ticket_group);").Error; err != nil {
				return err
			}
			//ticket_payment
			if err := tx.Exec("create table ticket_payment( ticket_payment int auto_increment primary key, name varchar(50) not null, min_amount int null, max_amount int null, config text not null);").Error; err != nil {
				return err
			}
			//ticket_type
			if err := tx.Exec("create table ticket_type( ticket_type int auto_increment primary key, event int not null, ticket_payment int not null, price double not null, contingent int not null, description varchar(100) not null, constraint fkEvent foreign key (event) references event (event) on update cascade on delete cascade, constraint fkPayment foreign key (ticket_payment) references ticket_payment (ticket_payment) on update cascade);").Error; err != nil {
				return err
			}
			if err := tx.Exec("create index event on ticket_type (event);").Error; err != nil {
				return err
			}
			if err := tx.Exec("create index ticket_payment on ticket_type (ticket_payment);").Error; err != nil {
				return err
			}
			//ticket_history
			if err := tx.Exec("create table ticket_history( ticket_history int auto_increment primary key, member int not null, ticket_payment int null, ticket_type int not null, purchase datetime null, redemption datetime null, created timestamp default current_timestamp() not null, code char(128) not null, constraint fkMember foreign key (member) references member (member) on update cascade on delete cascade, constraint fkTicketPayment foreign key (ticket_payment) references ticket_payment (ticket_payment) on update cascade, constraint fkTicketType foreign key (ticket_type) references ticket_type (ticket_type) on update cascade);").Error; err != nil {
				return err
			}
			if err := tx.Exec("create index member on ticket_history (member);").Error; err != nil {
				return err
			}
			if err := tx.Exec("create index ticket_payment on ticket_history (ticket_payment);").Error; err != nil {
				return err
			}
			return tx.Exec("create index ticket_type on ticket_history (ticket_type);").Error
		},
	}
}
