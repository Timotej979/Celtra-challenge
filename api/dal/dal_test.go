package dal_test

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Timotej979/Celtra-challenge/api/dal"
)

var disableRyuk string

func init() {
	flag.StringVar(&disableRyuk, "disable-ryuk", "", "Description of custom flag")
}

func TestMigrationAndCRUD(t *testing.T) {
	// Parse the command line flags
	flag.Parse()

	// Check if tests are running in verbose mode
	verbose := testing.Verbose()
	if verbose {
		// Check if the -disable-ryuk flag was set
		if disableRyuk != "y" {
			// Display a message and wait for rerun of the -v command with 'y' to continue
			fmt.Print("\n!!!!!!!!!!!!!! WARNING !!!!!!!!!!!!!!\n" +
				"!!! Running tests in verbose mode !!!\n" +
				"!!!-------------------------------!!!\n" +
				"!!  This means that testcontainer !!!\n" +
				"!!  module garbage collection will !!\n" +
				"!!  be disabled and DB containers  !!\n" +
				"!!  will keep on running until you !!\n" +
				"!!  remove them manually.          !!\n" +
				"!!!-------------------------------!!!\n" +
				"!!! To disable this message and   !!!\n" +
				"!!! disable garbage collection,   !!!\n" +
				"!!! please rerun the tests with   !!!\n" +
				"!!! the following command:        !!!\n" +
				"!!!-------------------------------!!!\n" +
				"!!! go test -v -disable-ryuk=y    !!!\n" +
				"!!!-------------------------------!!!\n\n")

			// Set environment variable to disable Ryuk otherwise
			t.Skip("Skipping tests in verbose mode")
		}

		// Set environment variable to disable Ryuk
		t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	}

	// Run tests for PostgreSQL
	t.Run("PostgreSQL", func(t *testing.T) {
		testMigrationAndCRUD(t, "postgres")
	})

	// Run tests for MongoDB
	//t.Run("MongoDB", func(t *testing.T) {
	//	testMigrationAndCRUD(t, "mongo")
	//})

	// Run tests for MySQL
	//t.Run("MySQL", func(t *testing.T) {
	//	testMigrationAndCRUD(t, "mysql")
	//})
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

	// Get the host of the container
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Get the port of the container
	portInt, err := strconv.Atoi(port)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new instance of the Data Access Layer
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

	// Connect to the database
	err = dalInstance.DbDriver.Connect()
	if err != nil {
		t.Fatal(err)
	}

	// Perform migration
	err = dalInstance.DbDriver.Migrate()
	assert.NoError(t, err)

	// Check if relation exists in Database
	switch dbType {
	case "postgres":
		out, err := exec.Command("psql", "postgres://testuser:testpassword@"+host+":"+port+"/testdb", "-c", "SELECT EXISTS (SELECT * FROM pg_tables WHERE tablename  = 'user_data');").Output()
		assert.NoError(t, err)
		if !strings.Contains(string(out), "exists") {
			t.Fatalf("Relation user_data does not exist in the database")
		}

	case "mongo":
		out, err := exec.Command("mongo", "mongo://testuser:testpassword@"+host+":"+port+"/testdb", "--eval", "db.getCollectionNames().includes('user_data');").Output()
		assert.NoError(t, err)
		if !strings.Contains(string(out), "true") {
			t.Fatalf("Collection user_data does not exist in the database")
		}

	case "mysql":
		cmd := exec.Command("mysql", "mysql://testuser:testpassword@"+host+":"+port+"/testdb", "-e", "SHOW TABLES LIKE 'user_data';")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		assert.NoError(t, err)
		assert.Contains(t, out.String(), "user_data")
		t.Log(out.String())
	}

	// Perform CRUD operations
	// Test create operation
	accountID := "test123"
	data := "test data"
	err = dalInstance.DbDriver.InsertUserData(accountID, data)
	assert.NoError(t, err)

	// Check if data was inserted
	switch dbType {
	case "postgres":
		// Check if column account_id exists in the user_data table
		out, err := exec.Command("psql", "postgres://testuser:testpassword@"+host+":"+port+"/testdb", "-c", "SELECT column_name FROM information_schema.columns WHERE table_name='user_data' and column_name='account_id';").Output()
		assert.NoError(t, err)
		if !strings.Contains(string(out), "account_id") {
			t.Fatalf("Column account_id does not exist in the user_data table")
		}

		// Check if column data exists in the user_data table
		out, err = exec.Command("psql", "postgres://testuser:testpassword@"+host+":"+port+"/testdb", "-c", "SELECT column_name FROM information_schema.columns WHERE table_name='user_data' and column_name='data';").Output()
		assert.NoError(t, err)
		if !strings.Contains(string(out), "data") {
			t.Fatalf("Column data does not exist in the user_data table")
		}

	case "mongo":
		cmd := exec.Command("mongo", "mongo://testuser:testpassword@"+host+":"+port+"/testdb", "--eval", "db.user_data.find({account_id: 'test123'});")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		assert.NoError(t, err)
		assert.Contains(t, out.String(), "test data")
		t.Logf(out.String())

	case "mysql":
		cmd := exec.Command("mysql", "mysql://testuser:testpassword@"+host+":"+port+"/testdb", "-e", "SELECT * FROM user_data WHERE account_id = 'test123';")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		assert.NoError(t, err)
		assert.Contains(t, out.String(), "test data")
		t.Logf(out.String())
	}

	// Test get operation
	retrievedData, _, err := dalInstance.DbDriver.GetUserData(accountID)
	assert.NoError(t, err)
	assert.Equal(t, data, retrievedData)

	// Test delete operation
	err = dalInstance.DbDriver.DeleteUserData(accountID)
	assert.NoError(t, err)

	// Verify deletion
	_, _, err = dalInstance.DbDriver.GetUserData(accountID)
	assert.Error(t, err)

	// Disconnect from the database
	err = dalInstance.DbDriver.Close()
	assert.NoError(t, err)
}

// Helper function to create a PostgreSQL container
func createPostgresContainer(ctx context.Context) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:alpine",
		ExposedPorts: []string{"5432"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432"),
			wait.ForExec([]string{"pg_isready", "-U", "testuser", "-d", "testdb"}),
		).WithDeadline(time.Minute * 2),
		AlwaysPullImage: true,
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}

	return container, port.Port()
}

// Helper function to create a MongoDB container
func createMongoContainer(ctx context.Context) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "testuser",
			"MONGO_INITDB_ROOT_PASSWORD": "testpassword",
			"MONGO_INITDB_DATABASE":      "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort("27017"),
			wait.ForExec([]string{"mongo", "mongodb://testuser:testpassword@localhost:27017/testdb", "--eval", "db.getCollectionNames();"}),
		).WithDeadline(time.Minute * 2),
		AlwaysPullImage: true,
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	port, err := container.MappedPort(ctx, "27017")
	if err != nil {
		panic(err)
	}

	return container, port.Port()
}

// Helper function to create a MySQL container
func createMySQLContainer(ctx context.Context) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306"},
		Env: map[string]string{
			"MYSQL_USER":          "testuser",
			"MYSQL_ROOT_PASSWORD": "testpassword",
			"MYSQL_DATABASE":      "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog(" MySQL Community Server - GPL"),
			wait.ForListeningPort("3306"),
			wait.ForExec([]string{"mysql", "-u", "testuser", "-h", "localhost", "-P", "3306", "-p", "testpassword", "testdb", "-e", "SHOW TABLES;"}),
		).WithDeadline(time.Minute * 2),
		AlwaysPullImage: true,
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	port, err := container.MappedPort(ctx, "3306")
	if err != nil {
		panic(err)
	}

	return container, port.Port()
}
