package repository

import (
	"database/sql"
	"fmt"

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
		INSERT INTO places (name, description, longitude, latitude)
		VALUES (:name, :description, :longitude, :latitude)
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

	// if len(place.PhotoURLs) > 0 {
	// 	if err := r.addPhotos(place.ID, place.PhotoURLs); err != nil {
	// 		return nil, fmt.Errorf("failed to add photos: %w", err)
	// 	}
	// }

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

	// photos, err := r.getPhotos(id)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get photos: %w", err)
	// }
	// place.PhotoURLs = photos

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

	// for _, place := range places {
	// 	photos, err := r.getPhotos(place.ID)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get photos for place %d: %w", place.ID, err)
	// 	}
	// 	place.PhotoURLs = photos
	// }

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
