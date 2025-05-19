package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	CountryRepository *CountryRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		CountryRepository: NewCountryRepository(db),
	}
}
