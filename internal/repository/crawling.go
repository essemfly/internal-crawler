package repository

import (
	"errors"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"gorm.io/gorm"
)

type CrawlingService struct {
	db *gorm.DB
}

func NewCrawlingService() *CrawlingService {
	db := config.DB()
	db.AutoMigrate(&domain.CrawlingSource{})
	return &CrawlingService{
		db,
	}
}

// CreateCrawlingSource adds a new CrawlingSource to the database
func (cs *CrawlingService) CreateCrawlingSource(source *domain.CrawlingSource) error {
	if source.SourceName == "" {
		return errors.New("source name is required")
	}
	var existingSource domain.CrawlingSource
	if cs.db.Where("source_name = ?", source.SourceName).First(&existingSource).RowsAffected > 0 {
		return errors.New("source name already exists")
	}
	result := cs.db.Create(source)
	return result.Error
}

// GetCrawlingSourceByID retrieves a CrawlingSource from the database by ID
func (cs *CrawlingService) GetCrawlingSourceByID(id uint) (*domain.CrawlingSource, error) {
	var source domain.CrawlingSource
	result := cs.db.First(&source, id)
	return &source, result.Error
}

// GetAllCrawlingSources retrieves all CrawlingSource records from the database
func (cs *CrawlingService) ListAllCrawlingSources() ([]*domain.CrawlingSource, error) {
	var sources []*domain.CrawlingSource
	result := cs.db.Where("is_active = ?", true).Find(&sources) // pass a pointer to the slice
	return sources, result.Error
}

// UpdateCrawlingSource updates an existing CrawlingSource in the database
func (cs *CrawlingService) UpdateCrawlingSource(source *domain.CrawlingSource) error {
	result := cs.db.Save(source)
	return result.Error
}
