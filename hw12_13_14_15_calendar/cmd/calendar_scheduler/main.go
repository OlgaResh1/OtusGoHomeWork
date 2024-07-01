package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"    //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/logger"    //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/rabbit"    //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/scheduler" //nolint:depguard
	"github.com/spf13/pflag"                                                        //nolint:depguard
)

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.toml", "Path to configuration file")
}

func main() {
	pflag.Parse()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := config.NewConfig()
	if err := cfg.ReadConfig(configFile); err != nil {
		log.Fatal("read config failed: ", err)
		return
	}
	logg := logger.New(cfg.Logger.Level, cfg.Logger.Format, cfg.Logger.AddSource)
	r, err := rabbit.New(cfg.RMQ, *logg)
	if err != nil {
		logg.Error("failed to start: "+err.Error(), "source", "scheduler")
		return
	}
	defer func() {
		err := r.Close()
		if err != nil {
			logg.Error("rabbit closed with error : "+err.Error(), "source", "scheduler")
		} else {
			logg.Info("rabbit closed ok", "source", "scheduler")
		}
	}()
	s, err := scheduler.New(cfg, *logg, r)
	if err != nil {
		logg.Error("failed to start: "+err.Error(), "source", "scheduler")
		return
	}
	go s.RequestNotifications(ctx)
	go s.SendNotifiction(ctx)

	logg.Info("Scheduler service started", "source", "scheduler")
	<-ctx.Done()
}
