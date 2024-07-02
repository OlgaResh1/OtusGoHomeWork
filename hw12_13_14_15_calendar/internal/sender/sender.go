package sender

import (
	"context"
	"fmt"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/logger"
)

type Consumer interface {
	Consume(ctx context.Context) (msg chan []byte, err error)
	Close() error
}

type Sender struct {
	logger   logger.Logger
	consumer Consumer
}

func New(logger logger.Logger, consumer Consumer) (*Sender, error) {
	return &Sender{
		logger:   logger,
		consumer: consumer,
	}, nil
}

func (s *Sender) SendNotification(ctx context.Context) error {
	msg, err := s.consumer.Consume(ctx)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Sender cicle ended")
			return nil
		case m := <-msg:
			fmt.Printf("Send message: %s\n", m)
		}
	}
}
