package models

type ListNotesFilter struct {
	Target           string `json:"target" validate:"required,oneof=company user"`
	PublicIdentifier string `json:"publicIdentifier" validate:"required,max=255"`
}

type ListNotes struct {
	Filters  []ListNotesFilter `json:"filters" validate:"required,dive"`
	AuthorID string            `json:"authorID" validate:"required,max=255"`
}
