package service

import (
	"github.com/ShekleinAleksey/top-places/internal/app/repository"
	"github.com/ShekleinAleksey/top-places/internal/entity"
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

func (s *CountryService) AddCountry(country *entity.Country) (int, error) {
	return s.repo.AddCountry(country)
}
