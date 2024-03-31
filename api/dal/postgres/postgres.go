package postgres

import (
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// UserData model
type UserData struct {
	ID        uint      `gorm:"primaryKey"`
	AccountID string    `gorm:"column:accountid;"`
	Timestamp time.Time `gorm:"column:timestamp;"`
	Data      string    `gorm:"column:data;"`
}

type PostgresDriver struct {
	// External data
	DbHost string
	DbPort int
	DbUser string
	DbPass string
	DbName string

	// Internal data
	Db *gorm.DB
}

// NewPostgresDriver creates a new instance of PostgresDriver
func NewPostgresDriver(host string, port int, user string, pass string, dbName string) *PostgresDriver {
	return &PostgresDriver{
		DbHost: host,
		DbPort: port,
		DbUser: user,
		DbPass: pass,
		DbName: dbName,
	}
}

// Connect establishes a connection to the database
func (p *PostgresDriver) Connect() error {
	log.Info().Msg("Connecting to the database...")

	// Construct the connection string
	dsn := "host=" + p.DbHost +
		" port=" + strconv.Itoa(p.DbPort) +
		" user=" + p.DbUser +
		" password=" + p.DbPass +
		" dbname=" + p.DbName +
		" sslmode=disable"

	// Open a connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
		return err
	}

	// Assign the database connection to the driver instance
	p.Db = db
	return nil
}

// Close closes the connection to the database
func (p *PostgresDriver) Close() error {
	log.Info().Msg("Closing the database connection...")

	// Close the connection to the database
	sqlDB, err := p.Db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to close the database connection")
		return err
	}
	sqlDB.Close()
	return nil
}

// Migrate the schema to database
func (p *PostgresDriver) Migrate() error {
	log.Info().Msg("Migrating the database...")

	// Perform the database migration
	err := p.Db.AutoMigrate(&UserData{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate the database")
		return err
	}
	return nil
}

// InsertUserData inserts a new user data record into the database
func (p *PostgresDriver) InsertUserData(accountID string, data string) error {
	log.Info().Msg("Inserting user data record...")

	userData := UserData{
		AccountID: accountID,
		Timestamp: time.Now(),
		Data:      data,
	}

	// Insert the user data record
	result := p.Db.Create(&userData)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to insert user data record")
		return result.Error
	}
	return nil
}

// GetUserData retrieves a user data record from the database
func (p *PostgresDriver) GetUserData(accountID string) (string, time.Time, error) {
	log.Info().Msg("Getting user data record...")
	var userData UserData

	// Get the user data record
	result := p.Db.Where("account_id = ?", accountID).First(&userData)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to get user data record")
		return "", time.Time{}, result.Error
	}
	return userData.Data, userData.Timestamp, nil
}

// DeleteUserData deletes a user data record from the database
func (p *PostgresDriver) DeleteUserData(accountID string) error {
	log.Info().Msg("Deleting user data record...")

	// Delete the user data record
	result := p.Db.Where("account_id = ?", accountID).Delete(&UserData{})
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to delete user data record")
		return result.Error
	}
	return nil
}
