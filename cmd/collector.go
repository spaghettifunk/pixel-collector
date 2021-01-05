package cmd

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/pixel-collector/collector"
	"github.com/spaghettifunk/pixel-collector/pkg/kafka"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serverHostFlag            = "server-host"
	serverPortFlag            = "server-port"
	prometheusHostURLFlag     = "prometheus-host"
	kafkaBrokersFlag          = "kafka-brokers"
	kafkaSASLMechanism        = "kafka-sasl-mechanism"
	kafkaProducerUsernameFlag = "kafka-producer-username"
	kafkaProducerPasswordFlag = "kafka-producer-password"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "collect",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		addr := viper.GetString(serverHostFlag)
		port := viper.GetString(serverPortFlag)
		promHost := viper.GetString(prometheusHostURLFlag)
		logDebug := viper.GetBool(logDebugFlag)

		// log level debug
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if logDebug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}

		// kafka client
		kc, err := getKafkaClient()
		if err != nil {
			log.Fatal().Err(err).Msgf("could not create Kafka client: %s", err.Error())
		}

		// process the responses asynchronous
		go func() {
			kc.ProcessResponse()
		}()

		// start prometheus service
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			if err := http.ListenAndServe(promHost, nil); err != nil {
				panic(err)
			}
		}()

		// create new Server object
		s, err := collector.NewServer(kc)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}

		// Start server
		go func() {
			if err := s.ListenAndServe(addr, port); err != nil {
				log.Info().Msg("server is shutting down...")
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server
		// Use a buffered channel to avoid missing signals as recommended for signal.Notify
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		if err := s.Shutdown(); err != nil {
			log.Fatal().Msg(err.Error())
		}

		kc.Close()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	f := serverCmd.PersistentFlags()

	f.String(serverHostFlag, "", "server address")
	f.String(serverPortFlag, "8082", "server port")
	f.String(prometheusHostURLFlag, ":9090", "The address to listen on for HTTP requests.")
	f.String(kafkaBrokersFlag, "localhost:29092", "Kafka brokers separated by ','")
	f.String(kafkaSASLMechanism, "", "Kafka authentication mechanism. Values accepted are ['SASL_PLAINTEXT','SASL_SSL']")
	f.String(kafkaProducerUsernameFlag, "", "Kafka producer username")
	f.String(kafkaProducerPasswordFlag, "", "Kafka producer password")

	viper.BindEnv(prometheusHostURLFlag, "PROMETHEUS_HOST")
	viper.BindEnv(kafkaBrokersFlag, "KAFKA_BROKERS")
	viper.BindEnv(kafkaSASLMechanism, "KAFKA_SASL_MECHANISM")
	viper.BindEnv(kafkaProducerUsernameFlag, "KAFKA_PRODUCER_USERNAME")
	viper.BindEnv(kafkaProducerPasswordFlag, "KAFKA_PRODUCER_PASSWORD")

	viper.BindPFlags(f)
}

func getKafkaClient() (*kafka.Client, error) {
	var options []func(*sarama.Config)

	brokers := viper.GetString(kafkaBrokersFlag)
	username := viper.GetString(kafkaProducerUsernameFlag)
	password := viper.GetString(kafkaProducerPasswordFlag)
	saslMechanism := viper.GetString(kafkaSASLMechanism)

	if username != "" && password != "" {
		options = append(options, kafka.Credentials(username, password))
	}

	if saslMechanism != "" {
		options = append(options, kafka.SASLMechanism(saslMechanism))
	}

	return kafka.NewClient(brokers, options...)
}
