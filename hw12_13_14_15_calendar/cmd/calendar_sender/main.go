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
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/sender"
	"github.com/spf13/pflag"
)

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "/etc/calendar/sender_config.toml", "Path to configuration file")
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
		logg.Error("failed to start: "+err.Error(), "source", "sender")
		return
	}
	defer func() {
		err := r.Close()
		if err != nil {
			logg.Error("rabbit closed with error : "+err.Error(), "source", "sender")
		} else {
			logg.Info("rabbit closed ok", "source", "sender")
		}
	}()
	s, err := sender.New(*logg, r)
	if err != nil {
		logg.Error("failed to create sender: "+err.Error(), "source", "sender")
		return
	}
	go func() {
		err := s.SendNotification(ctx)
		if err != nil {
			logg.Error("failed to send notification: "+err.Error(), "source", "sender")
			return
		}
	}()
	logg.Info("Sender service started", "source", "sender")
	<-ctx.Done()
}
