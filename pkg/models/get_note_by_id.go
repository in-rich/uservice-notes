package models

type GetNoteByID struct {
	AuthorID string `json:"authorID" validate:"required,max=255"`
	NoteID   string `json:"noteID" validate:"required"`
}
