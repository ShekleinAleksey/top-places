package service

import "github.com/ShekleinAleksey/top-places/internal/app/repository"

type Service struct {
	CountryService *CountryService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		CountryService: NewCountryService(*repo.CountryRepository),
	}
}
