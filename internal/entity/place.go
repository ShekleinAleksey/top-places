package entity

type Place struct {
	ID          int      `json:"id" db:"id"`
	Name        string   `json:"name" db:"name" binding:"required"`
	Description string   `json:"description" db:"description"`
	Longitude   float64  `json:"longitude" db:"longitude"`
	Latitude    float64  `json:"latitude" db:"latitude"`
	CountryID   int      `json:"-" db:"country_id"`
	Country     Country  `json:"country" db:"-"`
	PhotoURLs   []string `json:"url" db:"-"`
}

type PlacePhoto struct {
	ID      int    `json:"id" db:"id"`
	PlaceID int    `json:"place_id" db:"place_id"`
	URL     string `json:"url" db:"url"`
}
