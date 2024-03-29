package config

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrecedence(t *testing.T) {
	// Run the tests in a temporary directory
	tmpDir, err := os.MkdirTemp("", "userapi")
	require.NoError(t, err, "error creating a temporary test directory")
	testDir, err := os.Getwd()
	require.NoError(t, err, "error getting the current working directory")
	defer os.Chdir(testDir)
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "error changing to the temporary test directory")

	// Set arguments with the TOML config file
	t.Run("TOML config file", func(t *testing.T) {
		testcases := []struct {
			name                       string
			configFile                 string
			replaceHyphenWithCamelCase bool
		}{
			{name: "hyphen", configFile: "testdata/config-hyphen.toml"},
			{name: "camelCase", configFile: "testdata/config-camel.toml", replaceHyphenWithCamelCase: true},
		}
		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				replaceHyphenWithCamelCase = tc.replaceHyphenWithCamelCase
				defer func() { replaceHyphenWithCamelCase = false }()

				// Copy the config file into our temporary test directory
				configB, err := os.ReadFile(filepath.Join(testDir, tc.configFile))
				require.NoError(t, err, "error reading test config file")
				err = os.WriteFile(filepath.Join(tmpDir, "config.toml"), configB, 0644)
				require.NoError(t, err, "error writing test config file")
				defer os.Remove(filepath.Join(tmpDir, "config.toml"))

				// Run ./userapi
				cmd := NewRootCommand()
				output := &bytes.Buffer{}
				cmd.SetOut(output)
				cmd.Execute()

				gotOutput := output.String()
				wantOutput := "AppConfig: dev" +
					"\nDbUsername: root" +
					"\nDbPassword: REDACTED" +
					"\nDbName: data\n"

				assert.Equal(t, wantOutput, gotOutput, "expected the output to match the configuration settings")
			})
		}
	})

	// Set arguments with the TOML config file
	t.Run("JSON config file", func(t *testing.T) {
		testcases := []struct {
			name                           string
			configFile                     string
			replaceHyphenWithCamelCaseJSON bool
		}{
			{name: "hyphen", configFile: "testdata/config-hyphen.json"},
			{name: "camelCase", configFile: "testdata/config-camel.json", replaceHyphenWithCamelCaseJSON: true},
		}
		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				replaceHyphenWithCamelCase = tc.replaceHyphenWithCamelCaseJSON
				defer func() { replaceHyphenWithCamelCase = false }()

				// Copy the config file into our temporary test directory
				configB, err := os.ReadFile(filepath.Join(testDir, tc.configFile))
				require.NoError(t, err, "error reading test config file")
				err = os.WriteFile(filepath.Join(tmpDir, "config.json"), configB, 0644)
				require.NoError(t, err, "error writing test config file")
				defer os.Remove(filepath.Join(tmpDir, "config.json"))

				// Run ./userapi
				cmd := NewRootCommand()
				output := &bytes.Buffer{}
				cmd.SetOut(output)
				cmd.Execute()

				gotOutput := output.String()
				wantOutput := "AppConfig: dev" +
					"\nDbUsername: root" +
					"\nDbPassword: REDACTED" +
					"\nDbName: data\n"

				assert.Equal(t, wantOutput, gotOutput, "expected the output to match the configuration settings")
			})
		}
	})

	// Set arguments with the TOML config file
	t.Run("YAML config file", func(t *testing.T) {
		testcases := []struct {
			name                       string
			configFile                 string
			replaceHyphenWithCamelCase bool
		}{
			{name: "hyphen", configFile: "testdata/config-hyphen.yaml"},
			{name: "camelCase", configFile: "testdata/config-camel.yaml", replaceHyphenWithCamelCase: true},
		}
		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				replaceHyphenWithCamelCase = tc.replaceHyphenWithCamelCase
				defer func() { replaceHyphenWithCamelCase = false }()

				// Copy the config file into our temporary test directory
				configB, err := os.ReadFile(filepath.Join(testDir, tc.configFile))
				require.NoError(t, err, "error reading test config file")
				err = os.WriteFile(filepath.Join(tmpDir, "config.yaml"), configB, 0644)
				require.NoError(t, err, "error writing test config file")
				defer os.Remove(filepath.Join(tmpDir, "config.yaml"))

				// Run ./userapi
				cmd := NewRootCommand()
				output := &bytes.Buffer{}
				cmd.SetOut(output)
				cmd.Execute()

				gotOutput := output.String()
				wantOutput := "AppConfig: dev" +
					"\nDbUsername: root" +
					"\nDbPassword: REDACTED" +
					"\nDbName: data\n"

				assert.Equal(t, wantOutput, gotOutput, "expected the output to match the configuration settings")
			})
		}
	})

	// Set arguments with an environment variable
	t.Run("env var", func(t *testing.T) {
		// Run API_APP_CONFIG=prod API_DB_USERNAME=admin API_DB_PASSWORD=admin API_DB_NAME=database ./userapi
		os.Setenv("API_APP_CONFIG", "prod")
		os.Setenv("API_DB_USERNAME", "admin")
		os.Setenv("API_DB_PASSWORD", "admin")
		os.Setenv("API_DB_NAME", "database")
		defer os.Unsetenv("API_APP_CONFIG")
		defer os.Unsetenv("API_DB_USERNAME")
		defer os.Unsetenv("API_DB_PASSWORD")
		defer os.Unsetenv("API_DB_NAME")

		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.Execute()

		gotOutput := output.String()
		wantOutput := "AppConfig: prod" +
			"\nDbUsername: admin" +
			"\nDbPassword: REDACTED" +
			"\nDbName: database\n"

		assert.Equal(t, wantOutput, gotOutput, "expected the output to match the environment variables")
	})

	// Set arguments with a full flag
	t.Run("full flag", func(t *testing.T) {
		// Run ./userapi --app-config prod --db-username admin --db-password admin --db-name database
		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"--app-config", "prod", "--db-username", "admin", "--db-password", "admin", "--db-name", "database"})
		cmd.Execute()

		gotOutput := output.String()
		wantOutput := "AppConfig: prod" +
			"\nDbUsername: admin" +
			"\nDbPassword: REDACTED" +
			"\nDbName: database\n"

		assert.Equal(t, wantOutput, gotOutput, "expected the output to match the flag values")
	})

	// Set arguments with a shorthand flag
	t.Run("shorthand flag", func(t *testing.T) {
		// Run ./userapi -c prod -u admin -p admin -n database
		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"-c", "prod", "-u", "admin", "-p", "admin", "-n", "database"})
		cmd.Execute()

		gotOutput := output.String()
		wantOutput := "AppConfig: prod" +
			"\nDbUsername: admin" +
			"\nDbPassword: REDACTED" +
			"\nDbName: database\n"

		assert.Equal(t, wantOutput, gotOutput, "expected the output to match the flag values")
	})
}
