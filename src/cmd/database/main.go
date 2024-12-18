package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/zap"

	"github.com/andy74139/webserver/src/config"
	"github.com/andy74139/webserver/src/database"
	"github.com/andy74139/webserver/src/infra"
)

// init
func main() {
	// args
	ctx := context.Background()
	args := os.Args

	cmd := "migrate"
	if len(args) > 2 {
		panic("too many arguments")
	}
	if len(args) == 2 {
		cmd = args[1]
	}

	// logger
	// TODO: logger config, save logs
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("zap.NewDevelopment error: %w", err))
	}
	logger := log.Sugar()
	infra.SetDefaultLogger(logger)
	ctx = infra.SetLogger(ctx, logger)

	switch cmd {
	case "migrate":
		panic("not implemented")
	case "init":
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.GetLocalDSN())))
		db := bun.NewDB(sqldb, pgdialect.New())
		if _, err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"); err != nil {
			logger.Fatal("db.Exec error", zap.Error(err))
		}
		if err := createSchema(ctx, db); err != nil {
			logger.Fatal("createSchema error", zap.Error(err))
		}
		if err := createDefaultData(ctx, db); err != nil {
			logger.Fatal("createDefaultData error", zap.Error(err))
		}
	default:
		panic("unknown command")
	}
}

// DB processes
// * Initialization: create settings, schema, and default data, for new db instance
// * Migration: update schema, for updating to new schema
// * Backup: backup for snapshot of db, includes settings, schema, and data
// * Recovery: recover to specific snapshot, includes settings, schema, and data

// createSchema creates database schema for models.
func createSchema(ctx context.Context, db *bun.DB) error {
	models := []interface{}{
		(*database.User)(nil),
	}

	for _, model := range models {
		_, err := db.NewCreateTable().Model(model).Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func createDefaultData(ctx context.Context, db *bun.DB) error {
	id, err := uuid.Parse("00000000-0000-0000-0000-000000000001")
	if err != nil {
		return fmt.Errorf("parse default user ID error: %w", err)
	}

	users := []*database.User{
		{ID: id, Name: "TestUserCuteCapoo"},
	}
	if _, err := db.NewInsert().Model(&users).Exec(ctx); err != nil {
		return fmt.Errorf("insert users error: %w", err)
	}

	//migrate.NewMigrations()

	return nil
}
