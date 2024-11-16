package cassandra

import (
	"myapp/src/config"
	"time"

	"github.com/gocql/gocql"
	"github.com/rs/zerolog"
)

func NewCassandraConnection(config config.CassandraConfig, log zerolog.Logger) (*gocql.Session,error){
	cluster:=gocql.NewCluster(config.Hosts...)
	// cluster.Keyspace=config.Keyspace
	cluster.Port=config.Port
	cluster.Consistency=gocql.Quorum
	if config.Username!="" && config.Password!="" {
		cluster.Authenticator=gocql.PasswordAuthenticator{
			Username: config.Username,
			Password: config.Password,
		}
	}
	cluster.ProtoVersion=config.ProtoVersion
	
	var session *gocql.Session
	var err error
	
	maxRetries := 10
	retryDelay := 5 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		session, err = cluster.CreateSession()
		if err == nil {
			log.Info().Msgf("Connected to Cassandra on attempt %d", attempt)
			return session, nil
		}

		log.Error().Err(err).Msgf("Failed to connect to Cassandra (attempt %d/%d)", attempt, maxRetries)
		time.Sleep(retryDelay)
	}

	log.Error().Err(err).Msg("Exhausted all retries to connect to Cassandra")
	return nil, err

}