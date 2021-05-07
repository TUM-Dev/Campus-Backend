package firebase

import "gorm.io/gorm"

type MessagingService struct {
	db *gorm.DB
}

func New(db *gorm.DB) *MessagingService {

	return &MessagingService{db: db}
}

func (c *MessagingService) Run() error {

	return nil
}
