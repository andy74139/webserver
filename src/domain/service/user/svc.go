package user_svc

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/andy74139/webserver/src/domain/entity/user"
)

type service struct {
	userRepo user.Repository
}

func New(userRepo user.Repository) (user.Service, error) {
	if userRepo == nil {
		return nil, errors.New("userRepo is nil")
	}

	return &service{
		userRepo: userRepo,
	}, nil
}

func (s *service) Create(ctx context.Context, platform string, deviceID string, name string) error {
	return s.userRepo.Create(ctx, platform, deviceID, name)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*user.User, error) {
	return s.userRepo.Get(ctx, id)
}

func (s *service) GetIDByDevice(ctx context.Context, platform string, deviceID string) (uuid.UUID, error) {
	return s.userRepo.GetIDByDevice(ctx, platform, deviceID)
}

func (s *service) GetIDBySSO(ctx context.Context, ssoProvider string, ssoAccountID string) (uuid.UUID, error) {
	return s.userRepo.GetIDBySSO(ctx, ssoProvider, ssoAccountID)
}

func (s *service) CheckValidLoginUser(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.userRepo.CheckValidLoginUser(ctx, id)
}

func (s *service) Update(ctx context.Context, user1 *user.User) error {
	return s.userRepo.Update(ctx, user1)
}

func (s *service) AddSSO(ctx context.Context, id uuid.UUID, provider string, providerAccountID string) error {
	return s.userRepo.AddSSO(ctx, id, provider, providerAccountID)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}
