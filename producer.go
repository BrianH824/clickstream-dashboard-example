package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"clickstream-dashboard-example/models"
	"github.com/segmentio/kafka-go"
)

func main() {
	const topic = "clickstream_events"
	const broker = "web-user-clickstream-brian-3077.g.aivencloud.com:26239"

	producer := createKafkaConnection(broker, topic)
	defer producer.Close()

	rand.Seed(time.Now().UnixNano())

	populateInitialUsers()
	emitEvents(producer, 100)
}

func createKafkaConnection(broker string, topic string) *kafka.Writer {
	keypair, err := tls.LoadX509KeyPair("./conf/service.cert", "./conf/service.key")
	if err != nil {
		log.Fatalf("Failed to load TLS key pair: %s", err)
	}

	caCert, err := ioutil.ReadFile("./conf/ca.pem")
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %s", err)
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatalf("Failed to parse CA certificate")
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS: &tls.Config{
			Certificates: []tls.Certificate{keypair},
			RootCAs:      caCertPool,
		},
	}

	return kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   topic,
		Dialer:  dialer,
	})
}

func emitEvents(producer *kafka.Writer, messageLimit int) {
	for i := 0; i < messageLimit; i++ {
		event := models.GenerateRandomEvent()

		jsonValue, err := json.Marshal(event)
		if err != nil {
			log.Printf("Failed to marshal JSON: %v", err)
			continue
		}

		err = producer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(event.UserID),
			Value: jsonValue,
		})
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		} else {
			log.Printf("Sent event: %s", jsonValue)
		}

		time.Sleep(5 * time.Second)
	}
}

func populateInitialUsers() {
	for i := 0; i < 10; i++ {
		_ = models.CreateUser()
	}
}
