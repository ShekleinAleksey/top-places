package entity

type Country struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name" binding:"required"`
	Capital     string `json:"capital" db:"capital" binding:"required"`
	Language    string `json:"language" db:"language"`
	Currency    string `json:"currency" db:"currency"`
	Description string `json:"description" db:"description"`
	PhotoURL    string `json:"url" db:"photo_url"`
}
