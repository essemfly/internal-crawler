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

func (w *WishketService) SaveProject(project *domain.ProjectInfo) {
	w.db.Create(project)
}

func (w *WishketService) GetLastProject() domain.ProjectInfo {
	var project domain.ProjectInfo
	w.db.Last(&project)
	return project
}

func (w *WishketService) FindProjectByUrl(url string) []domain.ProjectInfo {
	var projects []domain.ProjectInfo
	w.db.Where("url = ?", url).Find(&projects)
	return projects
}
