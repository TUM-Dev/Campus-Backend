package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20240103000000
// made sure that question2answer have the correct fk-relationships
func migrate20240103000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240103000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`alter table question2answer
					add constraint question2answer_member_member_fk
					foreign key (member) references member (member)
            		on update cascade on delete cascade;`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`alter table question2answer
					add constraint question2answer_questionAnswers_answer_fk
					foreign key (answer) references questionAnswers (answer)
            		on update cascade on delete cascade;`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`alter table question2answer
					add constraint question2answer_question_question_fk
					foreign key (question) references question (question)
            		on update cascade on delete cascade;`).Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec(`alter table question2answer drop constraint question2answer_member_member_fk;`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`alter table question2answer drop constraint question2answer_questionAnswers_answer_fk;`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`alter table question2answer drop constraint question2answer_question_question_fk;`).Error; err != nil {
				return err
			}
			return nil
		},
	}
}
