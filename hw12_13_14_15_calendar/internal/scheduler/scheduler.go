package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"  //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/logger"  //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type Notification struct {
	ID            storage.EventID      `json:"id,omitempty"`
	OwnerID       storage.EventOwnerID `json:"userid,omitempty"`
	Title         string               `json:"title,omitempty"`
	StartDateTime time.Time            `json:"time,omitempty"`
}

type Producer interface {
	PublishMessage(ctx context.Context, message []byte) error
	Close() error
}

type Scheduler struct {
	logger           logger.Logger
	produser         Producer
	timePeriod       time.Duration
	eventsExpiration time.Duration
	notifyQueue      chan Notification
	calendarAddr     string
}

func New(_ context.Context, cfg config.Config, logger logger.Logger, p Producer) (*Scheduler, error) {
	d, err := time.ParseDuration(cfg.Scheduler.TimePeriod)
	if err != nil {
		return nil, err
	}
	exp, err := time.ParseDuration(cfg.Scheduler.EventsExpiration)
	if err != nil {
		return nil, err
	}
	return &Scheduler{
		logger:           logger,
		produser:         p,
		timePeriod:       d,
		eventsExpiration: exp,
		notifyQueue:      make(chan Notification),
		calendarAddr:     cfg.GrpcClient.Address,
	}, nil
}

func (s *Scheduler) RequestNotifications(ctx context.Context) {
	timer := time.NewTicker(s.timePeriod)
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Scheduler timer cicle finished!")
			return

		case <-timer.C:
			err := s.checkEventsToNotify(ctx)
			if err != nil {
				s.logger.Error(err.Error())
			}
			err = s.checkOldEvents(ctx)
			if err != nil {
				s.logger.Error(err.Error())
			}
		}
	}
}

func (s *Scheduler) SendNotifiction(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Scheduler finished!")
			return

		case n := <-s.notifyQueue:

			b, err := json.Marshal(n)
			if err != nil {
				s.logger.Error(err.Error())
				continue
			}
			err = s.produser.PublishMessage(ctx, b)
			if err != nil {
				s.logger.Error(err.Error())
				continue
			}
		}
	}
}
