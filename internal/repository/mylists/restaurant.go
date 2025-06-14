package mylists

import (
	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain/mylists"
	"gorm.io/gorm"
)

type RestaurantService struct {
	db *gorm.DB
}

func NewRestaurantService() *RestaurantService {
	db := config.DB()
	db.AutoMigrate(&mylists.Restaurant{})
	return &RestaurantService{
		db,
	}
}

func (s *RestaurantService) CreateRestaurant(restaurant *mylists.Restaurant) error {
	return s.db.Create(restaurant).Error
}

func (s *RestaurantService) GetRestaurantByID(id uint) (*mylists.Restaurant, error) {
	var restaurant mylists.Restaurant
	return &restaurant, s.db.First(&restaurant, id).Error
}

func (s *RestaurantService) GetAllRestaurants() ([]*mylists.Restaurant, error) {
	var restaurants []*mylists.Restaurant
	return restaurants, s.db.Find(&restaurants).Error
}

func (s *RestaurantService) UpdateRestaurant(restaurant *mylists.Restaurant) error {
	return s.db.Save(restaurant).Error
}
