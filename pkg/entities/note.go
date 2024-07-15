package entities

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Target string

const (
	TargetUser    Target = "user"
	TargetCompany Target = "company"
)

var _ sql.Scanner = (*Target)(nil)
var _ driver.Valuer = (*Target)(nil)

func (target Target) Valid() bool {
	switch target {
	case TargetUser, TargetCompany:
		return true
	default:
		return false
	}
}

func (target *Target) Scan(src interface{}) error {
	switch tsrc := src.(type) {
	case string:
		*target = Target(tsrc)
		if !target.Valid() {
			return fmt.Errorf("invalid target: %q", tsrc)
		}
		return nil
	case []byte:
		*target = Target(tsrc)
		if !target.Valid() {
			return fmt.Errorf("invalid target: %q", tsrc)
		}
		return nil
	case nil:
		return fmt.Errorf("scanning nil into Target")
	default:
		return fmt.Errorf("unsupported data type for Target: %T", src)
	}
}

func (target Target) Value() (driver.Value, error) {
	if !target.Valid() {
		return nil, fmt.Errorf("invalid target: %q", target)
	}
	return string(target), nil
}

type Note struct {
	bun.BaseModel `bun:"table:notes"`

	ID *uuid.UUID `bun:"id,pk,type:uuid"`

	AuthorID         string `bun:"author_id,notnull"`
	PublicIdentifier string `bun:"public_identifier,notnull"`
	Target           Target `bun:"target,notnull"`

	Content string `bun:"content,notnull"`

	UpdatedAt *time.Time `bun:"updated_at,notnull"`
}
