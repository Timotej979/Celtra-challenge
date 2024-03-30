package dal_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Timotej979/Celtra-challenge/api/dal"
)

func TestMigrationAndCRUD(t *testing.T) {
	// Run tests for PostgreSQL
	t.Run("PostgreSQL", func(t *testing.T) {
		testMigrationAndCRUD(t, "postgres")
	})

	// Run tests for MongoDB
	t.Run("MongoDB", func(t *testing.T) {
		testMigrationAndCRUD(t, "mongo")
	})

	// Run tests for MySQL
	t.Run("MySQL", func(t *testing.T) {
		testMigrationAndCRUD(t, "mysql")
	})
}

// Helper function to perform migration and CRUD operations for a specific database type
func testMigrationAndCRUD(t *testing.T, dbType string) {
	ctx := context.Background()

	// Create test container based on database type
	var container testcontainers.Container
	var port string

	switch dbType {
	case "postgres":
		container, port = createPostgresContainer(ctx)
	case "mongo":
		container, port = createMongoContainer(ctx)
	case "mysql":
		container, port = createMySQLContainer(ctx)
	default:
		t.Fatalf("Unsupported database type: %s", dbType)
	}

	defer container.Terminate(ctx)

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		t.Fatal(err)
	}

	config := &dal.DALConfig{
		DbType: dbType,
		DbHost: host,
		DbPort: portInt,
		DbUser: "testuser",
		DbPass: "testpassword",
		DbName: "testdb",
	}

	dalInstance, err := dal.NewDAL(config)
	if err != nil {
		t.Fatal(err)
	}
	defer dalInstance.DbDriver.Close()

	// Connect to the database
	err = dalInstance.DbDriver.Connect()
	if err != nil {
		t.Fatal(err)
	}

	// Perform migration
	err = dalInstance.DbDriver.Migrate()
	assert.NoError(t, err)

	// Perform CRUD operations
	// Test create operation
	accountID := "test123"
	data := "test data"
	err = dalInstance.DbDriver.InsertUserData(accountID, data)
	assert.NoError(t, err)

	// Test get operation
	retrievedData, err := dalInstance.DbDriver.GetUserData(accountID)
	assert.NoError(t, err)
	assert.Equal(t, data, retrievedData)

	// Test delete operation
	err = dalInstance.DbDriver.DeleteUserData(accountID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = dalInstance.DbDriver.GetUserData(accountID)
	assert.Error(t, err)

	// Disconnect from the database
	err = dalInstance.DbDriver.Close()
	assert.NoError(t, err)
}

// Helper function to create a PostgreSQL container
func createPostgresContainer(ctx context.Context) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		panic(err)
	}

	return container, port.Port()
}

// Helper function to create a MongoDB container
func createMongoContainer(ctx context.Context) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "testuser",
			"MONGO_INITDB_ROOT_PASSWORD": "testpassword",
			"MONGO_INITDB_DATABASE":      "testdb",
		},
		WaitingFor: wait.ForListeningPort("27017/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	port, err := container.MappedPort(ctx, "27017/tcp")
	if err != nil {
		panic(err)
	}

	return container, port.Port()
}

// Helper function to create a MySQL container
func createMySQLContainer(ctx context.Context) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "mysql:latest",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_USER":          "testuser",
			"MYSQL_ROOT_PASSWORD": "testpassword",
			"MYSQL_DATABASE":      "testdb",
		},
		WaitingFor: wait.ForListeningPort("3306/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	port, err := container.MappedPort(ctx, "3306/tcp")
	if err != nil {
		panic(err)
	}

	return container, port.Port()
}
