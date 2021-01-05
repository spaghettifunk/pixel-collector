package kafka

import (
	"errors"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

var (
	ackMessages = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "continuum",
		Subsystem: "events",
		Name:      "kafka_ack_counter",
		Help:      "Number of Acked messages sent to Kafka.",
	})
	errorMessages = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "continuum",
		Subsystem: "events",
		Name:      "kafka_errors_counter",
		Help:      "Number of error messages received from Kafka.",
	})
)

// Client is the object for sending the logs to Kafa
type Client struct {
	Producer sarama.AsyncProducer
}

// SASLMechanism functional option for the kafka configuration
// Values accepted: "OAUTHBEARER", "PLAIN", "SCRAM-SHA-256", "SCRAM-SHA-512", "GSSAPI"
func SASLMechanism(m string) func(*sarama.Config) {
	return func(cfg *sarama.Config) {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.Mechanism = sarama.SASLMechanism(m)
	}
}

// Credentials functional option for the kafka configuration
// For functional options see:
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func Credentials(username, password string) func(*sarama.Config) {
	return func(cfg *sarama.Config) {
		cfg.Net.SASL.User = username
		cfg.Net.SASL.Password = password
	}
}

// NewClient create a new object for interacting with Kafka
func NewClient(brokers string, options ...func(*sarama.Config)) (*Client, error) {
	bs := strings.Split(brokers, ",")

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true

	for _, opt := range options {
		opt(cfg)
	}

	// register prometheus metrics
	prometheus.MustRegister(ackMessages)
	prometheus.MustRegister(errorMessages)

	producer, err := sarama.NewAsyncProducer(bs, cfg)
	if err != nil {
		return &Client{}, err
	}
	return &Client{
		Producer: producer,
	}, nil
}

// Write sends the message to the kafka broker
func (c *Client) Write(topic string, event []byte) error {
	// if the event is empty ignore it
	if event == nil {
		log.Debug().Msg("Ignored event")
		return errors.New("event is empty")
	}

	// send to kafka topic
	c.Producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(event),
	}
	return nil
}

// ProcessResponse grabs results and errors from the producer asynchronously
func (c *Client) ProcessResponse() {
	for {
		select {
		case result := <-c.Producer.Successes():
			// increase ACK for Kafka in prometheus
			ackMessages.Inc()
			log.Debug().Msgf("> message: '%s' sent to partition  %d at offset %d\n", result.Value, result.Partition, result.Offset)
		case err := <-c.Producer.Errors():
			// increase Error for Kafka in prometheus
			errorMessages.Inc()
			log.Error().Msgf("Failed to produce message - error: %v", err)
		}
	}
}

// Close closes the producer object
func (c *Client) Close() error {
	return c.Producer.Close()
}
