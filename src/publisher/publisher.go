package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	handler "github.com/fuczak/go-worker/src/shared/error"
	"github.com/fuczak/go-worker/src/shared/message"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"strings"
)

type args struct {
	config *string
	queue  *string
	mType  *string
	mData  *string
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	values := args{
		config: flag.String("config", "./config.json", "Path to the config file"),
		queue:  flag.String("queue", "go_test", "Name of the queue to publish the message to"),
		mType:  flag.String("type", "Compute:Dots", "Type of the message to publish"),
		mData:  flag.String("data", "....", "Data of the message to publish"),
	}
	flag.Parse()

	err := validateInput(values)
	handler.Error(err)

	if err != nil {
		flag.PrintDefaults()
	}

	//cfg := config.InitConfig(*values.config)
	//conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.RabbitMq.Username, cfg.RabbitMq.Password, cfg.RabbitMq.Host, cfg.RabbitMq.Port))
	conn, err := amqp.Dial("amqp://go_test:go_test@172.16.0.235:5672/")
	handler.Error(err)

	ch, err := conn.Channel()
	handler.Error(err)

	messageBody, err := json.Marshal(message.Message{
		Type: *values.mType,
		Data: *values.mData,
	})
	handler.Error(err)

	err = ch.Publish("", *values.queue, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/string",
		Body:         messageBody,
	})
	handler.Error(err)

	if err == nil {
		log.Printf("Message '%s' published to queue '%s'", messageBody, *values.queue)
	}
}

func validateInput(values args) error {
	var e []string

	if _, err := os.Stat(*values.config); err != nil {
		e = append(e, err.Error())
	}

	if len(strings.TrimSpace(*values.queue)) == 0 {
		e = append(e, "Invalid queue name")
	}

	if len(strings.TrimSpace(*values.mType)) == 0 {
		e = append(e, "Invalid message type")
	}

	if len(strings.TrimSpace(*values.mData)) == 0 {
		e = append(e, "Invalid message data")
	}

	if len(e) > 0 {
		return errors.New(fmt.Sprintf("Input validation errors: %s", strings.Join(e, ", ")))
	}

	return nil
}
