package main

import (
	"context"
	"fmt"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/app"                          //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"                       //nolint:depguard
	memorystorage "github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage/memory" //nolint:depguard
	sqlstorage "github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage/sql"       //nolint:depguard
)

func setupStorage(ctx context.Context, cfg config.Config) (storage app.Storage, err error) {
	switch cfg.Storage.Type {
	case "inmemory":
		{
			storage, err = memorystorage.New(cfg)
		}
	case "sql":
		{
			storage, err = sqlstorage.New(ctx, cfg)
		}
	default:
		return nil, fmt.Errorf("storage type error %s", cfg.Storage.Type)
	}
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func closeStorage(ctx context.Context, storage app.Storage) error {
	return storage.Close(ctx)
}
