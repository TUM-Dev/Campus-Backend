package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20240102000000
// removes the option to create learning cards from the db as it has been removed from the api a long time ago
func migrate20240102000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240102000000",
		Migrate: func(tx *gorm.DB) error {
			//history+comments
			if err := tx.Exec("DROP table members_card").Error; err != nil {
				return err
			}
			if err := tx.Exec("DROP table members_card_answer_history").Error; err != nil {
				return err
			}
			if err := tx.Exec("DROP table card_box").Error; err != nil {
				return err
			}
			if err := tx.Exec("DROP table card_comment").Error; err != nil {
				return err
			}
			// tags
			if err := tx.Exec("DROP table card2tag").Error; err != nil {
				return err
			}
			if err := tx.Exec("DROP table tag").Error; err != nil {
				return err
			}
			//cards
			if err := tx.Exec("DROP table card_option").Error; err != nil {
				return err
			}
			if err := tx.Exec("DROP table card").Error; err != nil {
				return err
			}
			// things that dangle of card
			if err := tx.Exec("DROP table lecture").Error; err != nil {
				return err
			}
			if err := tx.Exec("DROP table card_type").Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec(`create table if not exists card_type(
				card_type int auto_increment primary key,
				title     varchar(255) null
) auto_increment = 2;`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`create table if not exists lecture(
				lecture int auto_increment primary key,
				title   varchar(255) null
			);`).Error; err != nil {
				return err
			}
			// cards
			if err := tx.Exec(`create table if not exists card(
				card           int auto_increment primary key,
				member         int                  null,
				lecture        int                  null,
				card_type      int                  null,
				title          varchar(255)         null,
				front_text     varchar(2000)        null,
				front_image    varchar(2000)        null,
				back_text      varchar(2000)        null,
				back_image     varchar(2000)        null,
				can_shift      tinyint(1) default 0 null,
				created_at     date                 not null,
				updated_at     date                 not null,
				duplicate_card int                  null,
				aggr_rating    float      default 0 null,
				constraint card_ibfk_1 foreign key (member) references member (member) on delete set null,
				constraint card_ibfk_2 foreign key (lecture) references lecture (lecture) on delete set null,
				constraint card_ibfk_3 foreign key (card_type) references card_type (card_type) on delete set null,
				constraint card_ibfk_4 foreign key (duplicate_card) references card (card) on delete set null
			);`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`create table if not exists card_option(
				card_option       int auto_increment primary key,
				card              int                      not null,
				text              varchar(2000) default '' null,
				is_correct_answer tinyint(1)    default 0  null,
				sort_order        int           default 0  not null,
				image             varchar(2000)            null,
				constraint card_option_ibfk_1 foreign key (card) references card (card) on delete cascade
);`).Error; err != nil {
				return err
			}
			// tags
			if err := tx.Exec(`create table if not exists tag(
				tag   int auto_increment primary key,
				title varchar(255) null
);`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`create table if not exists card2tag(
				tag  int not null,
				card int not null,
				primary key (tag, card),
				constraint card2tag_ibfk_1 foreign key (tag) references tag (tag) on delete cascade,
				constraint card2tag_ibfk_2 foreign key (card) references card (card) on delete cascade
);`).Error; err != nil {
				return err
			}
			//comments + history
			if err := tx.Exec(`create table if not exists card_comment(
				card_comment int auto_increment primary key,
				member       int           null,
				card         int           not null,
				rating       int default 0 null,
				created_at   date          not null,
				constraint card_comment_ibfk_1 foreign key (member) references member (member) on delete set null,
				constraint card_comment_ibfk_2 foreign key (card) references card (card) on delete cascade
		);`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`create table if not exists card_box(
				card_box int auto_increment primary key,
				member   int          null,
				title    varchar(255) null,
				duration int          not null,
				constraint card_box_ibfk_1 foreign key (member) references member (member) on delete cascade
			) auto_increment = 6;`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`create table if not exists members_card(
				member                   int not null,
				card                     int not null,
				card_box                 int null,
				last_answered_active_day int null,
				primary key (member, card),
				constraint members_card_ibfk_1 foreign key (member) references member (member) on delete cascade,
				constraint members_card_ibfk_2 foreign key (card) references card (card) on delete cascade,
				constraint members_card_ibfk_3 foreign key (card_box) references card_box (card_box) on delete set null
			);`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`create table if not exists members_card_answer_history(
				members_card_answer_history int auto_increment primary key,
				member                      int                       not null,
				card                        int                       null,
				card_box                    int                       null,
				answer                      varchar(2000)             null,
				answer_score                float(10, 2) default 0.00 null,
				created_at                  date                      not null,
				constraint members_card_answer_history_ibfk_1 foreign key (member) references member (member) on delete cascade,
				constraint members_card_answer_history_ibfk_2 foreign key (card) references card (card) on delete set null,
				constraint members_card_answer_history_ibfk_3 foreign key (card_box) references card_box (card_box) on delete set null
		);`).Error; err != nil {
				return err
			}
			return nil
		},
	}
}
