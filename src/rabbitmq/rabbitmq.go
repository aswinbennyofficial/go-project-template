package rabbitmq

import (
	"myapp/src/config"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

// NewRabbitMQClientWrapper selects the appropriate RabbitMQ client based on config
func NewRabbitMQClientWrapper(config config.RabbitMQConfig, log zerolog.Logger) (*amqp091.Connection, error) {
    if config.Mode == "cluster" {
        return NewClusterRabbitMQClient(config, log)
    }
    return NewStandaloneRabbitMQClient(config, log)
}


// NewStandaloneRabbitMQClient connects to a single RabbitMQ instance
func NewStandaloneRabbitMQClient(config config.RabbitMQConfig, log zerolog.Logger) (*amqp091.Connection, error) {
    var conn *amqp091.Connection
    var err error

    for i := 0; i < 30; i++ {
        conn, err = amqp091.Dial(config.Address)
        if err == nil {
            log.Info().
                Str("address", config.Address).
                Msg("Connected to RabbitMQ (Standalone) successfully")
            return conn, nil
        }

        log.Warn().
            Err(err).
            Int("attempt", i+1).
            Int("max_attempts", 30).
            Msg("Failed to connect to RabbitMQ (Standalone), retrying...")

        time.Sleep(5 * time.Second)
    }

    return nil, err
}



// NewClusterRabbitMQClient connects to a RabbitMQ cluster (failover mechanism)
func NewClusterRabbitMQClient(config config.RabbitMQConfig, log zerolog.Logger) (*amqp091.Connection, error) {
    var conn *amqp091.Connection
    var err error

    for i := 0; i < 30; i++ {
        for _, address := range config.ClusterAddresses {
            conn, err = amqp091.Dial(address)
            if err == nil {
                log.Info().
                    Str("connected_to", address).
                    Msg("Connected to RabbitMQ (Cluster) successfully")
                return conn, nil
            }

            log.Warn().
                Err(err).
                Str("attempted_address", address).
                Msg("Failed to connect to RabbitMQ node, trying next node...")
        }

        log.Warn().
            Int("attempt", i+1).
            Int("max_attempts", 30).
            Msg("Failed to connect to RabbitMQ Cluster, retrying...")

        time.Sleep(5 * time.Second)
    }

    
    return nil, err
}
