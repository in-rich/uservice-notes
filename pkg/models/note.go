package models

import "time"

type Note struct {
	ID               string
	Target           string
	PublicIdentifier string
	AuthorID         string
	Content          string
	UpdatedAt        *time.Time
}
