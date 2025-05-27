package service

import (
	"fmt"

	"github.com/ShekleinAleksey/top-places/internal/entity"
	"github.com/ShekleinAleksey/top-places/internal/repository"
)

type PlaceService struct {
	repo *repository.PlaceRepository
}

func NewPlaceService(repo *repository.PlaceRepository) *PlaceService {
	return &PlaceService{repo: repo}
}

func (s *PlaceService) Create(place *entity.Place) (*entity.Place, error) {
	if place.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	// if len(place.PhotoURLs) == 0 {
	// 	return nil, fmt.Errorf("at least one photo is required")
	// }
	return s.repo.Create(place)
}

func (s *PlaceService) GetByID(id int) (*entity.Place, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid ID")
	}
	return s.repo.GetByID(id)
}

func (s *PlaceService) GetAll() ([]*entity.Place, error) {
	return s.repo.GetAll()
}

func (s *PlaceService) Update(place *entity.Place) (*entity.Place, error) {
	if place.ID <= 0 {
		return nil, fmt.Errorf("invalid ID")
	}
	if place.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	return s.repo.Update(place)
}

func (s *PlaceService) Delete(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid ID")
	}
	return s.repo.Delete(id)
}

func (s *PlaceService) GetPlacesByCountry(countryID int) ([]entity.Place, error) {
	// Проверяем существование страны
	// if _, err := s.repo.GetCountryByID(countryID); err != nil {
	// 	return nil, fmt.Errorf("country not found")
	// }

	return s.repo.GetPlacesByCountryID(countryID)
}
