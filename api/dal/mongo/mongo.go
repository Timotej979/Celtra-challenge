package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserData model
type UserData struct {
	ID        uint      `bson:"_id,omitempty"`
	AccountID string    `bson:"account_id"`
	Timestamp time.Time `bson:"timestamp"`
	Data      string    `bson:"data"`
}

// MongoDBDriver represents the MongoDB database driver
type MongoDBDriver struct {
	// External data
	DbHost string
	DbPort int
	DbUser string
	DbPass string
	DbName string

	// Internal data
	Db *mongo.Client
}

// NewMongoDBDriver creates a new instance of MongoDBDriver
func NewMongoDBDriver(host string, port int, user string, pass string, dbName string, logger zerolog.Logger) *MongoDBDriver {
	// Assign the logger
	log.Logger = logger

	return &MongoDBDriver{
		DbHost: host,
		DbPort: port,
		DbUser: user,
		DbPass: pass,
		DbName: dbName,
	}
}

// Connect establishes a connection to the database
func (m *MongoDBDriver) Connect() error {
	log.Info().Msg("Connecting to the database...")

	// Construct connection string
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", m.DbUser, m.DbPass, m.DbHost, m.DbPort)
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
		return err
	}

	// Ping the MongoDB server to verify that the client is connected
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ping the database server")
		return err
	}

	// Assign the database client to the driver instance
	m.Db = client
	return nil
}

// Close closes the connection to the database
func (m *MongoDBDriver) Close() error {
	log.Info().Msg("Closing the database connection...")

	// Disconnect from MongoDB
	err := m.Db.Disconnect(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to disconnect from the database")
		return err
	}
	return nil
}

// Migrate the schema to the database
func (m *MongoDBDriver) Migrate() error {
	log.Info().Msg("Migrating the database schema...")

	// Get the database and collection
	collection := m.Db.Database(m.DbName).Collection("user_data")

	// Insert a dummy record to create the collection
	_, err := collection.InsertOne(context.Background(), UserData{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to migrate the database schema")
		return err
	}
	return nil
}

// InsertUserData inserts a new user data record into the database
func (m *MongoDBDriver) InsertUserData(accountID string, data string) error {
	log.Info().Msg("Inserting user data record...")

	// Get the database and collection
	collection := m.Db.Database(m.DbName).Collection("user_data")

	// Create the user data document
	userData := UserData{
		AccountID: accountID,
		Timestamp: time.Now(),
		Data:      data,
	}

	// Insert the user data document
	_, err := collection.InsertOne(context.Background(), userData)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert user data record")
		return err
	}
	return nil
}

// GetUserData retrieves a user data record from the database
func (m *MongoDBDriver) GetUserData(accountID string) (string, time.Time, error) {
	log.Info().Msg("Getting user data record...")
	var userData UserData

	// Get the database and collection
	collection := m.Db.Database(m.DbName).Collection("user_data")

	// Find the user data record
	err := collection.FindOne(context.Background(), UserData{AccountID: accountID}).Decode(&userData)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user data record")
		return "", time.Time{}, err
	}
	return userData.Data, userData.Timestamp, nil
}

// DeleteUserData deletes a user data record from the database
func (m *MongoDBDriver) DeleteUserData(accountID string) error {
	log.Info().Msg("Deleting user data record...")

	// Get the database and collection
	collection := m.Db.Database(m.DbName).Collection("user_data")

	// Delete the user data record
	_, err := collection.DeleteOne(context.Background(), UserData{AccountID: accountID})
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete user data record")
		return err
	}
	return nil
}
