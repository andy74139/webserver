package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Reference: https://bun.uptrace.dev/
// Hooks for created_at & updated_at: https://bun.uptrace.dev/guide/hooks.html
// Soft delete: https://bun.uptrace.dev/guide/soft-deletes.html

type User struct {
	bun.BaseModel `bun:"table:user"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	ID           uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Name         string    `bun:"name,notnull,type:varchar(256)"`
	Platform     *string   `bun:"platform,type:varchar(256)"`
	DeviceID     *string   `bun:"device_id,type:varchar(256)"`
	SSOProvider  *string   `bun:"sso_provider,type:varchar(256)"`
	SSOAccountID *string   `bun:"sso_account_id,type:varchar(256)"`
}

var _ bun.BeforeAppendModelHook = (*User)(nil)

func (m *User) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
