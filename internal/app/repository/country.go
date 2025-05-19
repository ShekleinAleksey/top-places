package repository

import (
	"fmt"

	"github.com/ShekleinAleksey/top-places/internal/entity"
	"github.com/jmoiron/sqlx"
)

type CountryRepository struct {
	db *sqlx.DB
}

func NewCountryRepository(db *sqlx.DB) *CountryRepository {
	return &CountryRepository{db: db}
}

func (r *CountryRepository) GetCountries() ([]entity.Country, error) {
	var countries []entity.Country
	query := fmt.Sprintf("SELECT * FROM country")
	err := r.db.Select(&countries, query)

	return countries, err
}

func (r *CountryRepository) AddCountry(country *entity.Country) (int, error) {
	var countryCreatedID int
	query := fmt.Sprintf("INSERT INTO country (name, description) values ($1, $2) RETURNING id")

	row := r.db.QueryRow(query, country.Name, country.Description)
	if err := row.Scan(&countryCreatedID); err != nil {
		return 0, err
	}

	return countryCreatedID, nil
}
