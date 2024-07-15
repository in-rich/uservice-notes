package models

import "time"

type Note struct {
	Target           string
	PublicIdentifier string
	AuthorID         string
	Content          string
	UpdatedAt        *time.Time
}
