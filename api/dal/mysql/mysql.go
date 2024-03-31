package mysql

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// UserData model
type UserData struct {
	ID        uint      `gorm:"primaryKey"`
	AccountID string    `gorm:"column:accountid;"`
	Timestamp time.Time `gorm:"column:timestamp;"`
	Data      string    `gorm:"column:data;"`
}

// MySQLDriver represents the MySQL database driver
type MySQLDriver struct {
	// External data
	DbHost string
	DbPort int
	DbUser string
	DbPass string
	DbName string

	// Internal data
	Db *gorm.DB
}

// NewMySQLDriver creates a new instance of MySQLDriver
func NewMySQLDriver(host string, port int, user string, pass string, dbName string) *MySQLDriver {
	return &MySQLDriver{
		DbHost: host,
		DbPort: port,
		DbUser: user,
		DbPass: pass,
		DbName: dbName,
	}
}

// Connect establishes a connection to the MySQL database
func (m *MySQLDriver) Connect() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", m.DbUser, m.DbPass, m.DbHost, m.DbPort, m.DbName)

	// Open a connection to the database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
		return err
	}

	// Assign the database connection to the driver instance
	m.Db = db
	return nil
}

// Close closes the connection to the MySQL database
func (m *MySQLDriver) Close() error {
	// Close the database connection
	sqlDB, err := m.Db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to close the database connection")
		return err
	}
	sqlDB.Close()
	return nil
}

// Migrate the schema to database
func (m *MySQLDriver) Migrate() error {
	// Migrate the schema
	err := m.Db.AutoMigrate(&UserData{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate the database schema")
		return err
	}
	return nil
}

// InsertUserData inserts a new user data record into the database
func (m *MySQLDriver) InsertUserData(accountID string, data string) error {
	log.Info().Msg("Inserting user data record...")

	userData := UserData{
		AccountID: accountID,
		Timestamp: time.Now(),
		Data:      data,
	}

	// Insert the user data record
	result := m.Db.Create(&userData)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to insert user data record")
		return result.Error
	}
	return nil
}

// GetUserData retrieves a user data record from the database
func (m *MySQLDriver) GetUserData(accountID string) (string, time.Time, error) {
	log.Info().Msg("Retrieving user data record...")
	var userData UserData

	// Get the user data record
	result := m.Db.Where("account_id = ?", accountID).First(&userData)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to retrieve user data record")
		return "", time.Time{}, result.Error
	}
	return userData.Data, userData.Timestamp, nil
}

// DeleteUserData deletes a user data record from the database
func (m *MySQLDriver) DeleteUserData(accountID string) error {
	log.Info().Msg("Deleting user data record...")

	// Delete the user data record
	result := m.Db.Where("account_id = ?", accountID).Delete(&UserData{})
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to delete user data record")
		return result.Error
	}
	return nil
}
