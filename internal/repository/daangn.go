package repository

import (
	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"gorm.io/gorm"
)

type DaangnService struct {
	db *gorm.DB
}

// NewDaangnService creates a new service with a DB connection
func NewDaangnService() *DaangnService {
	db := config.DB()
	db.AutoMigrate(&domain.DaangnProduct{}, &domain.DaangnKeyword{}, &domain.DaangnConfig{})
	return &DaangnService{
		db,
	}
}

// CreateDaangnProduct adds a new DaangnProduct to the database
func (ds *DaangnService) CreateDaangnProduct(product *domain.DaangnProduct) error {
	result := ds.db.Create(product)
	return result.Error
}

// GetDaangnProductByID retrieves a DaangnProduct from the database by ID
func (ds *DaangnService) GetDaangnProductByID(id uint) (*domain.DaangnProduct, error) {
	var product domain.DaangnProduct
	result := ds.db.First(&product, id)
	return &product, result.Error
}

// UpdateDaangnProduct updates an existing DaangnProduct in the database
func (ds *DaangnService) UpdateDaangnProduct(product *domain.DaangnProduct) error {
	result := ds.db.Save(product)
	return result.Error
}

// CreateDaangnKeyword adds a new DaangnKeyword to the database
func (ds *DaangnService) CreateDaangnKeyword(keyword *domain.DaangnKeyword) error {
	result := ds.db.Create(keyword)
	return result.Error
}

// GetDaangnKeywordByID retrieves a DaangnKeyword from the database by ID
func (ds *DaangnService) GetDaangnKeywordByID(id uint) (*domain.DaangnKeyword, error) {
	var keyword domain.DaangnKeyword
	result := ds.db.First(&keyword, id)
	return &keyword, result.Error
}

// UpdateDaangnKeyword updates an existing DaangnKeyword in the database
func (ds *DaangnService) UpdateDaangnKeyword(keyword *domain.DaangnKeyword) error {
	result := ds.db.Save(keyword)
	return result.Error
}

func (ds *DaangnService) ListLiveDaangnKeywords() ([]*domain.DaangnKeyword, error) {
	var keywords []*domain.DaangnKeyword
	result := ds.db.Where("is_live = ?", true).Find(&keywords)
	return keywords, result.Error
}

func (ds *DaangnService) TurnOffDaangnKeyword(keyword *domain.DaangnKeyword) error {
	keyword.IsLive = false
	result := ds.db.Save(keyword)
	return result.Error
}

func (ds *DaangnService) CreateDaangnConfig(config *domain.DaangnConfig) error {
	result := ds.db.Create(config)
	return result.Error
}

func (ds *DaangnService) GetLastDaangnConfig() (*domain.DaangnConfig, error) {
	var config domain.DaangnConfig
	result := ds.db.Order("id desc").First(&config)
	return &config, result.Error
}

func (ds *DaangnService) UpdateDaangnConfig(config *domain.DaangnConfig) error {
	result := ds.db.Save(config)
	return result.Error
}
