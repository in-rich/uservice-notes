package models

import "time"

type UpsertNote struct {
	Target           string     `json:"target" validate:"required,oneof=company user"`
	PublicIdentifier string     `json:"publicIdentifier" validate:"required,max=255"`
	AuthorID         string     `json:"authorID" validate:"required,max=255"`
	Content          string     `json:"content" validate:"max=15000"`
	UpdatedAt        *time.Time `json:"updatedAt"`
}
