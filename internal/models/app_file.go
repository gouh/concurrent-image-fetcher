package models

import "time"

// AppFile representa la estructura de tu tabla AppFile en la base de datos.
type AppFile struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	MimeType  string    `db:"mime_type" json:"mime_type"`
	Size      int64     `db:"size" json:"size"`
	Path      string    `db:"path" json:"path"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
