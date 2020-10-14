package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/fuczak/go-worker/src/shared/config"
	handler "github.com/fuczak/go-worker/src/shared/error"
	"github.com/fuczak/go-worker/src/shared/message"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var MessageHandlers = map[string]func(m message.Message) (bool, error){
	"Compute:Dots": handleComputeDots,
}

type ProcessingError struct {
	code int
}

func (m *ProcessingError) Error() string {
	return strconv.Itoa(m.code)
}

type args struct {
	config *string
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	values := args{
		config: flag.String("config", "./config.json", "Path to the config file"),
	}
	flag.Parse()

	err := validateInput(values)
	handler.Error(err)

	if err != nil {
		flag.PrintDefaults()
	}

	cfg := config.InitConfig(*values.config)
	//conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.RabbitMq.Username, cfg.RabbitMq.Password, cfg.RabbitMq.Host, cfg.RabbitMq.Port))
	conn, err := amqp.Dial("amqp://go_test:go_test@172.16.0.235:5672/")
	handler.Error(err)

	ch, err := conn.Channel()
	handler.Error(err)

	_ = ch.Qos(1, 0, false)

	q, err := ch.QueueDeclare(cfg.RabbitMq.Queue.Name, true, false, false, false, nil)
	handler.Error(err)

	msgs, err := ch.Consume(q.Name, strconv.Itoa(os.Getpid()), false, false, false, false, nil)
	handler.Error(err)

	loop := make(chan bool)

	go func() {
		for d := range msgs {
			ack, err := handleMessage(d)
			handler.Error(err)

			if ack == true {
				_ = d.Ack(false)
			} else {
				_ = d.Nack(false, false)
			}
		}
	}()

	log.Infof("Consumer initialized PID:%d", os.Getpid())
	log.Infof("Waiting for messages. To exit press CTRL+C")
	<-loop
}

func handleMessage(d amqp.Delivery) (bool, error) {
	log.Infof("Received a message: '%s'", d.Body)

	var m message.Message
	if unmarshalErr := json.Unmarshal(d.Body, &m); unmarshalErr != nil {
		handler.Error(unmarshalErr)

		return false, unmarshalErr
	}

	if MessageHandlers[m.Type] == nil {
		typeErr := errors.New("invalid message type")
		handler.Error(typeErr)

		return false, typeErr
	}

	log.Infof("Detected message of type: %s", m.Type)

	ack, err := MessageHandlers[m.Type](m)

	log.Printf("Message '%s' processed", d.Body)

	return ack, err
}

func validateInput(values args) error {
	var e[] string

	if _, err := os.Stat(*values.config); err != nil {
		e = append(e, err.Error())
	}

	if len(e) > 0 {
		return errors.New(fmt.Sprintf("Input validation errors: %s", strings.Join(e, ", ")))
	}

	return nil
}

func handleComputeDots(m message.Message) (bool, error) {
	dotCount := bytes.Count([]byte(m.Data), []byte("."))
	for i := 0; i < dotCount; i++ {
		log.Printf("Processing...")
		time.Sleep(200 * time.Millisecond)
		rand.Seed(time.Now().UnixNano())
		errorCode := rand.Intn(10)
		if errorCode > 7 {
			return false, &ProcessingError{
				errorCode,
			}
		}
	}

	return true, nil
}
