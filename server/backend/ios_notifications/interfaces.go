package ios_notifications

import "gorm.io/gorm"

type BaseRepository interface {
	NewRepository(db *gorm.DB) *BaseRepository
}

type BaseService interface {
	NewService(repository *BaseRepository) *BaseService
}
