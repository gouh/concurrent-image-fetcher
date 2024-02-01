package models

import "time"

// AppFile representa la estructura de tu tabla AppFile en la base de datos.
type AppFile struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	MimeType  string    `db:"mime_type"`
	Size      int64     `db:"size"`
	Path      string    `db:"path"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
