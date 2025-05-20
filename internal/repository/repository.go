package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	CountryRepository *CountryRepository
	PlaceRepository   *PlaceRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		CountryRepository: NewCountryRepository(db),
		PlaceRepository:   NewPlaceRepository(db),
	}
}
