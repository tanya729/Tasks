package amqp

import (
	"fmt"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/config"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/logger"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Ampq struct
type Ampq struct {
	Client  *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
}

// NewAmpq Create new message broker
func NewAmpq(conf *config.Ampq) (*Ampq, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", conf.User, conf.Password, conf.Host, conf.Port))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to AMQP broker")
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare AMQP channel")
	}

	queue, err := channel.QueueDeclare(conf.Queue, true, false, false, false, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to declare AMQP queue")
	}
	return &Ampq{conn, channel, queue}, nil
}

// Publish Send message to bus
func (am *Ampq) Publish(message []byte) error {
	err := am.Channel.Publish("", am.Queue.Name, false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         message,
		})
	if err != nil {
		return errors.Wrap(err, "failed to publish message")
	}
	logger.Debug(fmt.Sprintf("Message send: %s", message))
	return nil
}

// Subscribe to message queue
func (am *Ampq) Subscribe(consumerName string, handlerFunc func(amqp.Delivery)) error {
	messages, err := am.Channel.Consume(am.Queue.Name, consumerName, true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "failed to init consumer")
	}
	waitingChan := make(chan struct{})
	go func() {
		for msg := range messages {
			handlerFunc(msg)
		}
	}()
	<-waitingChan

	return nil
}
