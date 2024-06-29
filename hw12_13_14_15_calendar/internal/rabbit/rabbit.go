package rabbit

import (
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config" //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/logger" //nolint:depguard
	"github.com/streadway/amqp"                                                  //nolint:depguard
)

type Rabbit struct {
	exchangeName string
	queueName    string
	consumerTag  string
	routingKey   string
	reliable     bool

	conn    *amqp.Connection
	channel *amqp.Channel
	logger  logger.Logger
}

func New(cfg config.RMQConfig, logger logger.Logger) (r *Rabbit, err error) {
	r = &Rabbit{
		exchangeName: cfg.ExchangeName,
		queueName:    cfg.QueueName,
		consumerTag:  cfg.ConsumerTag,
		routingKey:   cfg.RoutingKey,
		reliable:     false, // cfg.reliable
		logger:       logger,
	}

	r.conn, err = amqp.Dial(cfg.URI)
	if err != nil {
		return nil, err
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		return nil, err
	}

	if err = r.channel.ExchangeDeclare(
		cfg.ExchangeName, // name of the exchange
		cfg.ExchangeKind, //  "direct", "fanout", "topic", "headers"
		true,             // durable
		false,            // delete when complete
		false,            // internal
		false,            // noWait
		nil,              // arguments
	); err != nil {
		return nil, err
	}

	queue, err := r.channel.QueueDeclare(
		cfg.QueueName, // name of the queue
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // noWait
		nil,           // arguments
	)
	if err != nil {
		return nil, err
	}

	if err = r.channel.QueueBind(
		queue.Name,       // name of the queue
		cfg.BindingKey,   // bindingKey
		cfg.ExchangeName, // sourceExchange
		false,            // noWait
		nil,              // arguments
	); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Rabbit) Close() error {
	err := r.channel.Close()
	if err != nil {
		return err
	}
	return r.conn.Close()
}
