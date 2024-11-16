package cassandra

import (
	"fmt"
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
	
	maxRetries := 20
	retryDelay := 5 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		session, err = cluster.CreateSession()
		if err == nil {
			log.Info().Msgf("Connected to Cassandra on attempt %d", attempt)
			if err=ensureKeyspace(session, config, log);err!=nil{
				return nil,err
			}
			session.Close()
			
			cluster.Keyspace = config.Keyspace
			session, err = cluster.CreateSession()
			if err!=nil{
				return nil,err
			}


			return session, nil
		}

		log.Error().Err(err).Msgf("Failed to connect to Cassandra (attempt %d/%d)", attempt, maxRetries)
		time.Sleep(retryDelay)
	}

	log.Error().Err(err).Msg("Exhausted all retries to connect to Cassandra")
	return nil, err

}


// EnsureKeyspace checks if the keyspace exists and creates it if it does not
func ensureKeyspace(session *gocql.Session, config config.CassandraConfig, logger zerolog.Logger) error {
	// Query to check if the keyspace exists
	query := fmt.Sprintf("SELECT keyspace_name FROM system_schema.keyspaces WHERE keyspace_name='%s'", config.Keyspace)
	iter := session.Query(query).Iter()

	var name string
	found := iter.Scan(&name)
	iter.Close()

	if found {
		logger.Info().Msgf("Cassandra Keyspace '%s' already exists", config.Keyspace)
		return nil
	}

	// Create the keyspace if not found
	createKeyspaceQuery := fmt.Sprintf(`
		CREATE KEYSPACE %s
		WITH replication = {
			'class': '%s',
			'replication_factor': %d
		}`,
		config.Keyspace,config.Replication.Strategy,config.Replication.Factor)

	if err := session.Query(createKeyspaceQuery).Exec(); err != nil {
		logger.Error().Err(err).Msg("Failed to create cassandra keyspace")
		return fmt.Errorf("failed to create cassandra keyspace: %w", err)
	}

	logger.Info().Msgf("Cassandra Keyspace '%s' created successfully", config.Keyspace)
	return nil
}
