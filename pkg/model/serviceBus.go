package model

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/config"
)

type EventBusHandlerFn func(event string, args []byte)

type ServiceBusMessage struct {
	Event      string     `json:"event"`
	ExecutorId *uuid.UUID `json:"executor_id"`
	Arguments  any        `json:"args"`
}

type ServiceBus struct {
	conn     *amqp.Connection
	handlers map[string][]*EventBusHandlerFn
}

func NewServiceBus() *ServiceBus {
	return &ServiceBus{
		handlers: make(map[string][]*EventBusHandlerFn),
	}
}

func (s *ServiceBus) Publish(key string, message *ServiceBusMessage) error {
	msg, err := json.Marshal(message)
	if err != nil {
		logrus.Error("Failed to marshal message: ", err)
		return err
	}
	conn, err := amqp.Dial(config.GetInstance().RabbitMQ.Url)

	if err != nil {
		logrus.Error("Failed to connect to RabbitMQ: ", err)
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		logrus.Error("Failed to open a channel: ", err)
		return err
	}
	defer ch.Close()
	// если exchange нет, то создаем
	err = ch.ExchangeDeclare("autocarat.login.service",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Error("Failed to declare a exchange: ", err)
		return err
	}
	// если очереди нет, то создаем
	q, err := ch.QueueDeclare(
		config.GetInstance().RabbitMQ.QueueName+"."+key,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Error("Failed to declare a queue: ", err)
		return err
	}
	err = ch.QueueBind(q.Name, "#."+key, "autocarat.login.service", false, nil)
	if err != nil {
		logrus.Error("Failed to bind queue to exchange: ", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = ch.PublishWithContext(ctx,
		"autocarat.login.service",
		config.GetInstance().RabbitMQ.QueueName+"."+key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})
	if err != nil {
		logrus.Error("Failed to publish a message: ", err)
		return err
	}
	return nil
}

func (s *ServiceBus) runHandler(event string, handlersPointer *map[string][]*EventBusHandlerFn, args []byte) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error("panic in service bus handler:", event, args, err)
		}
	}()
	handlers := *handlersPointer
	if handlers[event] != nil {
		for _, handler := range handlers[event] {
			(*handler)(event, args)
		}
	} else {
		logrus.Warn("No handlers message bus for event '", event, "', this message will be ignored")
	}
}

func (s *ServiceBus) addHandler(event string, handler EventBusHandlerFn, handlersPointer *map[string][]*EventBusHandlerFn) {
	handlers := *handlersPointer
	if handlers[event] == nil {
		handlers[event] = make([]*EventBusHandlerFn, 0)
	}
	handlers[event] = append(handlers[event], &handler)
	logrus.Info("Add new handler message bus for event '", event, "' : ", handler)
}

func (s *ServiceBus) removeHandler(event string, handler EventBusHandlerFn, handlersPointer *map[string][]*EventBusHandlerFn) {
	handlers := *handlersPointer
	for i, h := range handlers[event] {
		if h == &handler {
			handlers[event] = append(handlers[event][:i], handlers[event][i+1:]...)
			break
		}
	}
	logrus.Info("Remove handler message bus for event '", event, "' : ", handler)
}

func (s *ServiceBus) On(event string, handler EventBusHandlerFn) {
	s.addHandler(event, handler, &s.handlers)
}
func (s *ServiceBus) Off(event string, handler EventBusHandlerFn) {
	s.removeHandler(event, handler, &s.handlers)
}

func (s *ServiceBus) Listen(key string) error {
	conn, err := amqp.Dial(config.GetInstance().RabbitMQ.Url)

	if err != nil {
		logrus.Error("Failed to connect to RabbitMQ: ", err)
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		logrus.Error("Failed to open a channel: ", err)
		return err
	}
	// если exchange нет, то создаем
	err = ch.ExchangeDeclare("autocarat.login.service",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Error("Failed to declare a exchange: ", err)
		return err
	}
	// если очереди нет, то создаем
	q, err := ch.QueueDeclare(
		config.GetInstance().RabbitMQ.QueueName+"."+key,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Error("Failed to declare a queue: ", err)
		return err
	}
	//err = ch.QueueBind(q.Name, "#."+key, "autocarat.login.service", false, nil)
	//if err != nil {
	//	logrus.Error("Failed to bind queue to exchange: ", err)
	//	return err
	//}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	go func() {
		logrus.Info("Start listen queue: ", q.Name)
		defer conn.Close()
		defer ch.Close()

		for {
			for d := range msgs {
				msg := &ServiceBusMessage{}
				err := json.Unmarshal(d.Body, msg)
				if err != nil {
					logrus.Error("Failed to unmarshal message: ", err)
					continue
				} else {
					logrus.Info("Get message from queue, ", msg)
					args, err := json.Marshal(msg.Arguments)
					if err != nil {
						logrus.Error("Failed to marshal arguments: ", err)
						continue
					}
					s.runHandler(msg.Event, &s.handlers, args)
				}
			}
		}
		logrus.Info("Stop listen queue: ", q.Name, ch.IsClosed())
	}()

	return nil
}
