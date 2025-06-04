package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ShekleinAleksey/top-places/internal/entity"
	"github.com/jmoiron/sqlx"
)

type PlaceRepository struct {
	db *sqlx.DB
}

func NewPlaceRepository(db *sqlx.DB) *PlaceRepository {
	return &PlaceRepository{db: db}
}

func (r *PlaceRepository) Create(place *entity.Place) (*entity.Place, error) {
	query := `
		INSERT INTO places (name, description, longitude, latitude, country_id)
		VALUES (:name, :description, :longitude, :latitude, :country_id)
		RETURNING id
	`

	rows, err := r.db.NamedQuery(query, place)
	if err != nil {
		return nil, fmt.Errorf("failed to create place: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(place); err != nil {
			return nil, fmt.Errorf("failed to scan created place: %w", err)
		}
	}

	if len(place.PhotoURLs) > 0 {
		for _, url := range place.PhotoURLs {
			if err := r.addPhotos(place.ID, url); err != nil {
				return nil, fmt.Errorf("failed to add photo %s: %w", url, err)
			}
		}
	}

	return place, nil
}

func (r *PlaceRepository) GetByID(id int) (*entity.Place, error) {
	place := &entity.Place{}
	query := `
		SELECT *
		FROM places
		WHERE id = $1
	`

	err := r.db.Get(place, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("place not found")
		}
		return nil, fmt.Errorf("failed to get place: %w", err)
	}

	photos, err := r.getPhotos(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get photos: %w", err)
	}
	place.PhotoURLs = photos

	return place, nil
}

func (r *PlaceRepository) GetAll() ([]*entity.Place, error) {
	var places []*entity.Place
	query := `
		SELECT *
		FROM places
	`

	err := r.db.Select(&places, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get places: %w", err)
	}

	for _, place := range places {
		photos, err := r.getPhotos(place.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get photos for place %d: %w", place.ID, err)
		}
		place.PhotoURLs = photos
	}

	return places, nil
}

func (r *PlaceRepository) Update(place *entity.Place) (*entity.Place, error) {
	query := `
		UPDATE places
		SET name = :name,
			description = :description,
			longitude = :longitude,
			latitude = :latitude,
		WHERE id = :id
		RETURNING id
	`

	result, err := r.db.NamedExec(query, place)
	if err != nil {
		return nil, fmt.Errorf("failed to update place: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("place not found")
	}

	// if err := r.updatePhotos(place.ID, place.PhotoURLs); err != nil {
	// 	return nil, fmt.Errorf("failed to update photos: %w", err)
	// }

	return place, nil
}

func (r *PlaceRepository) Delete(id int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM place_photos WHERE place_id = $1", id); err != nil {
		return fmt.Errorf("failed to delete photos: %w", err)
	}

	result, err := tx.Exec("DELETE FROM places WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete place: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("place not found")
	}

	return tx.Commit()
}

func (r *PlaceRepository) GetPlacesByCountryID(countryID int) ([]entity.Place, error) {
	var places []entity.Place
	query := `
        SELECT id, name, description, longitude, latitude 
        FROM places 
        WHERE country_id = $1
    `
	err := r.db.Select(&places, query, countryID)

	for _, place := range places {
		photoUrl, err := r.getPhotos(place.CountryID)
		if err != nil {
			return nil, fmt.Errorf("failed to get photos for place %d: %w", place.ID, err)
		}
		place.PhotoURLs = photoUrl
	}
	return places, err
}

func (r *PlaceRepository) SearchByName(query string, limit int) ([]entity.Place, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []entity.Place{}, nil
	}

	rows, err := r.db.Query(`
		SELECT id, name, description, longitude, latitude, country_id
		FROM places
		WHERE name ILIKE '%' || $1 || '%'
		ORDER BY name
		LIMIT $2
	`, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search places: %w", err)
	}
	defer rows.Close()

	var places []entity.Place
	for rows.Next() {
		var p entity.Place
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Longitude, &p.Latitude, &p.CountryID); err != nil {
			return nil, fmt.Errorf("failed to scan place: %w", err)
		}
		places = append(places, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return places, nil
}

// Фото мест
func (r *PlaceRepository) addPhotos(placeID int, url string) error {
	query := `
		INSERT INTO place_photos (place_id, url)
		VALUES ($1, $2)
		RETURNING id
	`
	_, err := r.db.Exec(query, placeID, url)
	return err
}

func (r *PlaceRepository) getPhotos(placeID int) ([]string, error) {
	var photos []string
	err := r.db.Select(&photos, "SELECT url FROM place_photos WHERE place_id = $1", placeID)
	return photos, err
}

func (r *PlaceRepository) updatePhoto(photoID int, placeID int, url string) error {
	query := `
        UPDATE place_photos
        SET url = $3
        WHERE id = $1 AND place_id = $2
    `
	_, err := r.db.Exec(query, photoID, placeID, url)
	return err
}

func (r *PlaceRepository) deletePhotos(id int) error {
	_, err := r.db.Exec("DELETE FROM place_photos WHERE id = $1", id)
	return err
}
