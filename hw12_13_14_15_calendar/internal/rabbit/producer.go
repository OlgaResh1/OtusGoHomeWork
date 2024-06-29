package rabbit

import (
	"context"

	"github.com/streadway/amqp" //nolint:depguard
)

func (r *Rabbit) PublishMessage(_ context.Context, message []byte) (err error) {
	if r.reliable {
		if err := r.channel.Confirm(false); err != nil {
			return err
		}
		confirms := r.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		defer r.confirmOne(confirms)
	}

	if err = r.channel.Publish(
		r.exchangeName, // publish to an exchange
		r.routingKey,   // routing to 0 or more queues
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            message,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return err
	}

	return nil
}

func (r *Rabbit) confirmOne(confirms <-chan amqp.Confirmation) {
	if confirmed := <-confirms; confirmed.Ack {
		r.logger.Info("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		r.logger.Info("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
