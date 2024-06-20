package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/app"                      //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"                   //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/logger"                   //nolint:depguard
	internalgrpc "github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/server/grpc" //nolint:depguard
	internalhttp "github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/server/http" //nolint:depguard
	"github.com/spf13/pflag"                                                                       //nolint:depguard
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
	defer closeStorage(ctx, storage)

	calendar := app.New(logg, storage)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		serverHTTP := internalhttp.NewServer(cfg, logg, calendar)
		if err := serverHTTP.Start(ctx); err != nil {
			logg.Error("failed to start http server: "+err.Error(), "source", "http")
			cancel()
			os.Exit(1)
		}
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverHTTP.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: "+err.Error(), "source", "http")
		}
	}()
	go func() {
		defer wg.Done()
		serverGrpc := internalgrpc.NewServer(cfg, logg, calendar)
		if err := serverGrpc.Start(ctx); err != nil {
			logg.Error("failed to start grpc server: "+err.Error(), "source", "grpc")
			cancel()
			os.Exit(1)
		}
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverGrpc.Stop(ctx); err != nil {
			logg.Error("failed to stop grpc server: "+err.Error(), "source", "grpc")
		}
	}()
	wg.Wait()
}
