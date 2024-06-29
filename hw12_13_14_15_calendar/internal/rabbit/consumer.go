package rabbit

import (
	"context"
	"log"
)

func (r *Rabbit) Consume(ctx context.Context) (msg chan []byte, err error) {
	msg = make(chan []byte)

	deliveries, err := r.channel.Consume(
		r.queueName,   // name
		r.consumerTag, // consumerTag,
		false,         // noAck
		false,         // exclusive
		false,         // noLocal
		false,         // noWait
		nil,           // arguments
	)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			close(msg)
			log.Println("close messages channel")
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case del := <-deliveries:
				if err := del.Ack(false); err != nil {
					log.Println(err)
				}

				select {
				case <-ctx.Done():
					return
				case msg <- del.Body:
				}
			}
		}
	}()
	return msg, nil
}
