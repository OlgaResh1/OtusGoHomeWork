package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/app"
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/server/http"
	"github.com/spf13/pflag"
)

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	pflag.Parse()

	if pflag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	cfg := config.NewConfig()
	if err := cfg.ReadConfig(configFile); err != nil {
		log.Fatal("read config failed: ", err)
		return
	}
	logg := logger.New(cfg.Logger.Level, cfg.Logger.Format, cfg.Logger.AddSource)

	storage, err := setupStorage(ctx, cfg)
	if err != nil {
		log.Fatal("error create storage: ", err)
		return
	}
	defer func() {
		err = closeStorage(ctx, storage)
		if err != nil {
			logg.Error("error close storage: " + err.Error())
			return
		}
	}()
	calendar := app.New(logg, storage)
	serverHTTP := internalhttp.NewServer(cfg, logg, calendar)
	serverGrpc := internalgrpc.NewServer(cfg, logg, calendar)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		if err := serverHTTP.Start(ctx); err != nil {
			if !serverHTTP.IsClosed {
				logg.Error("failed to start HTTP-server: "+err.Error(), "source", "http")
			}
			cancel()
		}
	}()
	go func() {
		defer wg.Done()
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverHTTP.Stop(ctx); err != nil {
			logg.Error("failed to stop HTTP-server: "+err.Error(), "source", "http")
		} else {
			logg.Info("HTTP-server stopped ok", "source", "http")
		}
	}()
	go func() {
		if err := serverGrpc.Start(ctx); err != nil {
			logg.Error("failed to start GRPC-server: "+err.Error(), "source", "grpc")
			cancel()
		}
	}()
	go func() {
		defer wg.Done()
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverGrpc.Stop(ctx); err != nil {
			logg.Error("failed to stop GRPC-server: "+err.Error(), "source", "grpc")
		} else {
			logg.Info("GRPC-server stopped ok", "source", "grpc")
		}
	}()
	wg.Wait()
}
