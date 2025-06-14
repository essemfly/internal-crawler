package mylists

import (
	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain/mylists"
	"gorm.io/gorm"
)

type CollectionService struct {
	db *gorm.DB
}

func NewCollectionService() *CollectionService {
	db := config.DB()
	db.AutoMigrate(&mylists.Collection{})
	return &CollectionService{
		db,
	}
}

func (s *CollectionService) CreateCollection(collection *mylists.Collection) error {
	return s.db.Create(collection).Error
}

func (s *CollectionService) GetCollectionByID(id uint) (*mylists.Collection, error) {
	var collection mylists.Collection
	return &collection, s.db.First(&collection, id).Error
}

func (s *CollectionService) GetAllCollections() ([]*mylists.Collection, error) {
	var collections []*mylists.Collection
	return collections, s.db.Find(&collections).Error
}

func (s *CollectionService) UpdateCollection(collection *mylists.Collection) error {
	return s.db.Save(collection).Error
}
