package models

type GetNote struct {
	Target           string `json:"target" validate:"required,oneof=company user"`
	PublicIdentifier string `json:"publicIdentifier" validate:"required,max=255"`
	AuthorID         string `json:"authorID" validate:"required,max=255"`
}
