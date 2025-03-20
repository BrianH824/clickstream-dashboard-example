package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	TOPIC_NAME := "clickstream_events"

	keypair, err := tls.LoadX509KeyPair("./conf/service.cert", "./conf/service.key")
	if err != nil {
		log.Fatalf("Failed to load access key and/or access certificate: %s", err)
	}

	caCert, err := ioutil.ReadFile("./conf/ca.pem")
	if err != nil {
		log.Fatalf("Failed to read CA certificate file: %s", err)
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatalf("Failed to parse CA certificate file: %s", err)
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS: &tls.Config{
			Certificates: []tls.Certificate{keypair},
			RootCAs:      caCertPool,
		},
	}

	// init producer
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"web-user-clickstream-brian-3077.g.aivencloud.com:26239"},
		Topic:   TOPIC_NAME,
		Dialer:  dialer,
	})

	// produce 100 messages
	for i := 0; i < 1; i++ {
		message := fmt.Sprint("Hello from Go using SSL ", i+1, "!")
		producer.WriteMessages(context.Background(), kafka.Message{Value: []byte(message)})
		log.Printf("Message sent: " + message)
		time.Sleep(time.Second)
	}

	producer.Close()
}
