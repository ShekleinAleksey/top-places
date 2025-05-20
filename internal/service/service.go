package service

import "github.com/ShekleinAleksey/top-places/internal/repository"

type Service struct {
	CountryService *CountryService
	PlaceService   *PlaceService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		CountryService: NewCountryService(*repo.CountryRepository),
		PlaceService:   NewPlaceService(*&repo.PlaceRepository),
	}
}
