package models

type GetAllNotes struct {
	Limit  int64 `json:"limit" validate:"required,min=1,max=1000"`
	Offset int64 `json:"offset" validate:"min=0"`
}
