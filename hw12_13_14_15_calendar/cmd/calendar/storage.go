package main

import (
	"context"
	"fmt"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/app"
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage/sql"
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
