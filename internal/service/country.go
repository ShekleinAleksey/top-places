package service

import (
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
