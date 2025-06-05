package service

import (
	"fmt"

	"github.com/ShekleinAleksey/top-places/internal/entity"
	"github.com/ShekleinAleksey/top-places/internal/repository"
)

type PlaceService struct {
	placeRepo   *repository.PlaceRepository
	countryRepo *repository.CountryRepository
}

func NewPlaceService(placeRepo *repository.PlaceRepository, countryRepo *repository.CountryRepository) *PlaceService {
	return &PlaceService{
		placeRepo:   placeRepo,
		countryRepo: countryRepo,
	}
}

func (s *PlaceService) Create(place *entity.Place) (*entity.Place, error) {
	if place.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	// if len(place.PhotoURLs) == 0 {
	// 	return nil, fmt.Errorf("at least one photo is required")
	// }
	return s.placeRepo.Create(place)
}

func (s *PlaceService) GetByID(id int) (*entity.Place, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid ID")
	}
	place, err := s.placeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	country, err := s.countryRepo.GetCountryByID(place.CountryID)
	if err != nil {
		return nil, err
	}
	place.Country = country

	return place, nil
}

func (s *PlaceService) GetAll() ([]*entity.Place, error) {
	places, err := s.placeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, place := range places {
		country, err := s.countryRepo.GetCountryByID(place.CountryID)
		if err != nil {
			return nil, err
		}
		place.Country = country
	}

	return places, nil

}

func (s *PlaceService) Update(place *entity.Place) (*entity.Place, error) {
	if place.ID <= 0 {
		return nil, fmt.Errorf("invalid ID")
	}
	if place.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	return s.placeRepo.Update(place)
}

func (s *PlaceService) Delete(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid ID")
	}
	return s.placeRepo.Delete(id)
}

func (s *PlaceService) GetPlacesByCountry(countryID int) ([]*entity.Place, error) {
	// Проверяем существование страны
	// if _, err := s.repo.GetCountryByID(countryID); err != nil {
	// 	return nil, fmt.Errorf("country not found")
	// }
	places, err := s.placeRepo.GetPlacesByCountryID(countryID)
	if err != nil {
		return nil, err
	}

	for _, place := range places {
		country, err := s.countryRepo.GetCountryByID(place.CountryID)
		if err != nil {
			return nil, err
		}
		place.Country = country
	}

	return places, nil
}

func (s *PlaceService) SearchPlaces(query string, limit int) ([]*entity.Place, error) {
	places, err := s.placeRepo.SearchByName(query, limit)
	if err != nil {
		return nil, err
	}

	if places == nil {
		return []*entity.Place{}, nil
	}
	for _, place := range places {
		country, err := s.countryRepo.GetCountryByID(place.CountryID)
		if err != nil {
			return nil, err
		}
		place.Country = country
	}

	return places, nil
}
