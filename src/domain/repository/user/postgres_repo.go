package user_repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/andy74139/webserver/src/database"
	"github.com/andy74139/webserver/src/domain/entity/user"
)

// postgresql account repository
type postgresRepo struct {
	db *bun.DB
}

func NewPostgresRepo(db *bun.DB) (user.Repository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &postgresRepo{
		db: db,
	}, nil
}

func (r *postgresRepo) Create(ctx context.Context, platform string, deviceID string, name string) error {
	user1 := &database.User{
		Platform: &platform,
		DeviceID: &deviceID,
		Name:     name,
	}

	if _, err := r.db.NewInsert().Model(user1).Exec(ctx); err != nil {
		return fmt.Errorf("insert user error: %w", err)
	}
	return nil
}

func (r *postgresRepo) Get(ctx context.Context, id uuid.UUID) (*user.User, error) {
	user1 := &database.User{}
	if err := r.db.NewSelect().Model(user1).Where("id = ?", id).Scan(ctx, user1); err != nil {
		return nil, fmt.Errorf("select error: %w", err)
	}

	return &user.User{
		ID:   user1.ID,
		Name: user1.Name,
	}, nil
}

func (r *postgresRepo) GetIDByDevice(ctx context.Context, platform string, deviceID string) (uuid.UUID, error) {
	user1 := &database.User{}
	query := r.db.NewSelect().Model(user1).Column("id").Where("platform = ? and device_id = ?", platform, deviceID)
	if err := query.Scan(ctx, user1); errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, user.ErrNotFound
	} else if err != nil {
		return uuid.Nil, fmt.Errorf("select error: %w", err)
	}

	return user1.ID, nil
}

func (r *postgresRepo) GetIDBySSO(ctx context.Context, ssoProvider string, ssoAccountID string) (uuid.UUID, error) {
	user1 := &database.User{}
	query := r.db.NewSelect().Model(user1).Column("id").Where("sso_provider = ? and sso_account_id = ?", ssoProvider, ssoAccountID)
	if err := query.Scan(ctx, user1); err != nil {
		return uuid.Nil, fmt.Errorf("select error: %w", err)
	}

	return user1.ID, nil
}

func (r *postgresRepo) CheckValidLoginUser(ctx context.Context, id uuid.UUID) (bool, error) {
	user1 := &database.User{}
	isExist, err := r.db.NewSelect().Model(user1).Where("id = ?", id).Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("select error: %w", err)
	}

	return isExist, nil
}

func (r *postgresRepo) Update(ctx context.Context, user1 *user.User) error {
	query := r.db.NewUpdate().Model(user1).Where("id = ?", user1.ID)
	if user1.Name != "" {
		query = query.SetColumn("name", "?", user1.Name)
	}

	if result, err := query.Exec(ctx); err != nil {
		return fmt.Errorf("update error: %w", err)
	} else if rows, err2 := result.RowsAffected(); err2 != nil {
		// SHOULD NOT BE HERE for postgres
		return fmt.Errorf("RowsAffected error: %w", err2)
	} else if rows == 0 {
		return fmt.Errorf("no user %s", user1.ID)
	}
	return nil
}

func (r *postgresRepo) AddSSO(ctx context.Context, id uuid.UUID, provider string, providerAccountID string) error {
	query := r.db.NewUpdate().Model(&database.User{}).Where("id = ? AND platform is not NULL", id)
	query = query.SetColumn("sso_provider", "?", provider).SetColumn("sso_account_id", "?", providerAccountID)
	query = query.SetColumn("platform", "?", nil).SetColumn("device_id", "", nil)

	if result, err := query.Exec(ctx); err != nil {
		return fmt.Errorf("update error: %w", err)
	} else if rows, err2 := result.RowsAffected(); err2 != nil {
		// SHOULD NOT BE HERE for postgres
		return fmt.Errorf("RowsAffected error: %w", err2)
	} else if rows == 0 {
		isExists, err := r.db.NewSelect().Model(&database.User{}).Where("id = ?", id).Exists(ctx)
		if err != nil {
			return fmt.Errorf("select error: %w", err)
		} else if !isExists {
			return fmt.Errorf("no user: %s", id.String())
		}
		return fmt.Errorf("user is not device login: %s", id.String())
	}
	return nil
}

func (r *postgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if result, err := r.db.NewDelete().Model(&database.User{}).Where("id = ?", id).Exec(ctx); err != nil {
		return fmt.Errorf("delete error: %w", err)
	} else if rows, err2 := result.RowsAffected(); err2 != nil {
		// SHOULD NOT BE HERE for postgres
		return fmt.Errorf("RowsAffected error: %w", err2)
	} else if rows == 0 {
		return fmt.Errorf("no user: %s", id.String())
	}
	return nil
}
