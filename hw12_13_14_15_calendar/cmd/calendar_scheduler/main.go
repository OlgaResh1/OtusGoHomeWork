package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/logger"
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/spf13/pflag"
)

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.toml", "Path to configuration file")
}

func main() {
	pflag.Parse()

	cfg := config.NewConfig()
	if err := cfg.ReadConfig(configFile); err != nil {
		log.Fatal("read config failed: ", err)
		return
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logg := logger.New(cfg.Logger.Level, cfg.Logger.Format, cfg.Logger.AddSource)
	r, err := rabbit.New(cfg.RMQ, *logg)
	if err != nil {
		logg.Error("failed to start: rqm.uri: "+cfg.RMQ.URI+", error:"+err.Error(), "source", "scheduler")
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
