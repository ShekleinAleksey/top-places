package repository

import (
	"fmt"
	"strings"

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
	query := fmt.Sprintf("SELECT * FROM countries")
	err := r.db.Select(&countries, query)

	return countries, err
}

func (r *CountryRepository) GetCountryByID(id int) (entity.Country, error) {
	query := `
        SELECT * 
        FROM countries 
        WHERE id = $1
    `

	var country entity.Country

	err := r.db.Get(&country, query, id)
	if err != nil {
		return entity.Country{}, err
	}

	return country, nil
}

func (r *CountryRepository) AddCountry(country *entity.Country) (int, error) {
	query := `
        INSERT INTO countries (
            name, 
            capital, 
            language, 
            currency, 
            description, 
            photo_url
        ) 
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `

	var countryID int
	err := r.db.QueryRow(
		query,
		country.Name,
		country.Capital,
		country.Language,
		country.Currency,
		country.Description,
		country.PhotoURL,
	).Scan(&countryID)

	if err != nil {
		return 0, fmt.Errorf("failed to add country: %v", err)
	}

	return countryID, nil
}

func (r *CountryRepository) UpdateCountry(country *entity.Country) (*entity.Country, error) {
	query := `
        UPDATE countries 
        SET name = :name,
            capital = :capital,
            language = :language,
            currency = :currency,
            description = :description,
        WHERE id = :id
        RETURNING id
    `

	result, err := r.db.NamedExec(query, country)
	if err != nil {
		return nil, fmt.Errorf("failed to update country: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("country not found")
	}

	return country, nil
}

func (r *CountryRepository) DeleteCountry(id int) (int, error) {
	// Проверяем существование страны перед удалением
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM countries WHERE id = $1)"
	err := r.db.Get(&exists, checkQuery, id)
	if err != nil {
		return 0, fmt.Errorf("failed to check country existence: %w", err)
	}
	if !exists {
		return 0, fmt.Errorf("country with ID %d not found", id)
	}

	// Выполняем удаление
	query := "DELETE FROM countries WHERE id = $1 RETURNING id"
	var deletedID int

	err = r.db.QueryRow(query, id).Scan(&deletedID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete country: %w", err)
	}

	return deletedID, nil
}

func (r *CountryRepository) SearchByName(query string, limit int) ([]entity.Country, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil
	}

	rows, err := r.db.Query(`
		SELECT id, name, capital, language, currency, description, photo_url
		FROM countries
		WHERE name ILIKE '%' || $1 || '%'
		ORDER BY name
		LIMIT $2
	`, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search countries: %w", err)
	}
	defer rows.Close()

	var countries []entity.Country
	for rows.Next() {
		var c entity.Country
		if err := rows.Scan(&c.ID, &c.Name, &c.Capital, &c.Language, &c.Currency, &c.Description, &c.PhotoURL); err != nil {
			return nil, fmt.Errorf("failed to scan country: %w", err)
		}
		countries = append(countries, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return countries, nil
}
