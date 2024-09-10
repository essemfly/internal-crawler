package repository

import (
	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"gorm.io/gorm"
)

type WishketService struct {
	db *gorm.DB
}

func NewWishketService() *WishketService {
	db := config.DB()
	db.AutoMigrate(&domain.ProjectInfo{})
	return &WishketService{
		db,
	}
}
