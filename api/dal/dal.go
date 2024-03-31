package dal

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/Timotej979/Celtra-challenge/api/dal/mongo"
	"github.com/Timotej979/Celtra-challenge/api/dal/mysql"
	"github.com/Timotej979/Celtra-challenge/api/dal/postgres"
)

// DALConfig represents the configuration for the Data Access Layer
type DALConfig struct {
	DbType string
	DbHost string
	DbPort int
	DbUser string
	DbPass string
	DbName string
}

type DatabaseDriver interface {
	Connect() error
	Close() error
	Migrate() error
	InsertUserData(accountID string, data string) error
	GetUserData(accountID string) (string, time.Time, error)
	DeleteUserData(accountID string) error
}

// DAL is the Data Access Layer
type DAL struct {
	// Configuration
	DbConfig *DALConfig

	// Database driver interface
	DbDriver DatabaseDriver
}

// NewDAL creates a new instance of the Data Access Layer
func NewDAL(config *DALConfig) (*DAL, error) {
	// Create the Data Access Layer
	dal := &DAL{
		DbConfig: config,
	}

	// Create the database driver
	switch config.DbType {
	case "postgres":
		dal.DbDriver = postgres.NewPostgresDriver(config.DbHost, config.DbPort, config.DbUser, config.DbPass, config.DbName)
	case "mysql":
		dal.DbDriver = mysql.NewMySQLDriver(config.DbHost, config.DbPort, config.DbUser, config.DbPass, config.DbName)
	case "mongo":
		dal.DbDriver = mongo.NewMongoDBDriver(config.DbHost, config.DbPort, config.DbUser, config.DbPass, config.DbName)
	default:
		err := fmt.Errorf("invalid database type: %s", config.DbType)
		log.Error().Err(err).Msg("failed to create DAL")
		return nil, err
	}

	return dal, nil
}
