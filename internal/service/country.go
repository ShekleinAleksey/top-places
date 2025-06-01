package service

import (
	"fmt"

	"github.com/ShekleinAleksey/top-places/internal/entity"
	"github.com/ShekleinAleksey/top-places/internal/repository"
)

type CountryService struct {
	repo repository.CountryRepository
}

func NewCountryService(repo repository.CountryRepository) *CountryService {
	return &CountryService{repo: repo}
}

func (s *CountryService) GetCountries() ([]entity.Country, error) {
	return s.repo.GetCountries()
}

func (s *CountryService) GetCountryByID(id int) (entity.Country, error) {
	return s.repo.GetCountryByID(id)
}

func (s *CountryService) AddCountry(country *entity.Country) (int, error) {
	return s.repo.AddCountry(country)
}

func (s *CountryService) DeleteCountry(id int) (int, error) {
	return s.repo.DeleteCountry(id)
}

func (s *CountryService) UpdateCountry(country *entity.Country) (*entity.Country, error) {
	// Проверяем существование страны
	if _, err := s.repo.GetCountryByID(country.ID); err != nil {
		return nil, fmt.Errorf("country not found")
	}

	// Валидация данных
	if country.Name == "" {
		return nil, fmt.Errorf("country name is required")
	}
	if country.Capital == "" {
		return nil, fmt.Errorf("capital is required")
	}

	return s.repo.UpdateCountry(country)
}

func (s *CountryService) SearchCountries(query string, limit int) ([]entity.Country, error) {
	return s.repo.SearchByName(query, limit)
}
