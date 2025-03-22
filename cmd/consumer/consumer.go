package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// OpenSearchConfig TODO: when secrets get moved out of repo, this will get refactored.
type OpenSearchConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Index    string `json:"index"`
}

var openSearchConfig OpenSearchConfig

const (
	kafkaBroker = "web-user-clickstream-brian-3077.g.aivencloud.com:26239"
	kafkaTopic  = "clickstream_events"
)

func main() {
	loadOpenSearchConfig("conf/opensearch.json")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaBroker},
		Topic:       kafkaTopic,
		StartOffset: kafka.FirstOffset,
		Dialer:      getAivenDialer(),
	})

	log.Println("Starting Kafka to OpenSearch consumer...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Printf("Consumed from Kafka: key=%s", string(msg.Key))

		if err := indexToOpenSearch(msg.Value); err != nil {
			log.Printf("Failed to index to OpenSearch: %v", err)
		}
	}
}

func loadOpenSearchConfig(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open OpenSearch config file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Failed to close OpenSearch config file: %v", err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&openSearchConfig); err != nil {
		log.Fatalf("Failed to parse OpenSearch config file: %v", err)
	}
}

func indexToOpenSearch(jsonData []byte) error {
	url := openSearchConfig.Host + "/" + openSearchConfig.Index + "/_doc"

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	req.SetBasicAuth(openSearchConfig.Username, openSearchConfig.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("OpenSearch response (%d): %s", resp.StatusCode, string(body))
		return err
	}

	log.Println("Indexed to OpenSearch")
	return nil
}

func getAivenDialer() *kafka.Dialer {
	cert, err := tls.LoadX509KeyPair("conf/service.cert", "conf/service.key")
	if err != nil {
		log.Fatalf("Failed to load client cert/key: %v", err)
	}

	caCert, err := ioutil.ReadFile("conf/ca.pem")
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("Failed to parse CA cert")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	return &kafka.Dialer{
		Timeout:   10 * time.Second,
		TLS:       tlsConfig,
		DualStack: true,
	}
}
