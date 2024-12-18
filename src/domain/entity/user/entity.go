package user

// User domain handles for CRUD of account login info.
// It is base of auth and other account-related domains.
// It doesn't responsible for auth and user info.

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("not found")
)

type Service interface {
	Create(ctx context.Context, platform string, deviceID string, name string) error
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	GetIDByDevice(ctx context.Context, platform string, deviceID string) (uuid.UUID, error) // TODO: check if only get user ID
	GetIDBySSO(ctx context.Context, ssoProvider string, ssoAccountID string) (uuid.UUID, error)
	CheckValidLoginUser(ctx context.Context, id uuid.UUID) (bool, error)
	Update(ctx context.Context, user *User) error
	AddSSO(ctx context.Context, id uuid.UUID, provider string, providerAccountID string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Repository interface {
	Create(ctx context.Context, platform string, deviceID string, name string) error
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	GetIDByDevice(ctx context.Context, platform string, deviceID string) (uuid.UUID, error)
	GetIDBySSO(ctx context.Context, ssoProvider string, ssoAccountID string) (uuid.UUID, error)
	CheckValidLoginUser(ctx context.Context, id uuid.UUID) (bool, error)
	Update(ctx context.Context, user *User) error
	AddSSO(ctx context.Context, id uuid.UUID, provider string, providerAccountID string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type User struct {
	ID   uuid.UUID
	Name string
}
