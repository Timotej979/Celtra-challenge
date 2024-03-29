package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Modified source from "Sting of the Viper" (https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/) project for production use.
// Why use this implementation? This implementation represents a good practice how any application should be configured for fast DevOps integration and deployment.

// CLI argument parser static settings
const (
	// The name of our config file, without the file extension because viper supports many different config file languages.
	defaultConfigFilename = "config"

	// The environment variable prefix of all environment variables bound to our command line flags.
	// Example conversion: --number -> API_NUMBER
	envPrefix = "API"

	// Replace hyphenated flag names with camelCase in the config file
	replaceHyphenWithCamelCase = false
)

// Environment variable custom storage type
type EnvVarStore struct {
	AppConfig  string
	DbUsername string
	DbPassword string
	DbName     string
}

// Get the root command for our CLI tool and extract the values from it.\
func GetEnvVars() (*EnvVarStore, error) {
	// Create the root command
	cmd := NewRootCommand()

	// Execute the command and check for errors
	if err := cmd.Execute(); err != nil {
		return nil, err
	}

	// Extract the values from the command
	vars := &EnvVarStore{
		AppConfig:  cmd.Flag("app-config").Value.String(),
		DbUsername: cmd.Flag("db-username").Value.String(),
		DbPassword: cmd.Flag("db-password").Value.String(),
		DbName:     cmd.Flag("db-name").Value.String(),
	}

	return vars, nil
}

// Build the cobra command that handles our command line tool.
func NewRootCommand() *cobra.Command {
	// Define our variables
	variables := &EnvVarStore{
		AppConfig:  "dev",
		DbUsername: "Celtra",
		DbPassword: "C3ltr4Ch4ll3ng3",
		DbName:     "UserData",
	}

	// Define our command
	rootCmd := &cobra.Command{
		Use:   "userapi",
		Short: "User API is a simple API that fetches user data from a database.",
		Long: `User API is a simple API that fetches user data from a database.
			    It demonstrates how to use cobra and viper to bind command line flags to configuration files and environment variables.
				The command hierarchy is as follows:
					- flags > environment variables > configuration files and the defaults set by the tool`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			out := cmd.OutOrStdout()

			// Print the final resolved value from binding cobra flags and viper config
			log.Printf("AppConfig: %s\n", variables.AppConfig)
			log.Printf("DbUsername: %s\n", variables.DbUsername)
			log.Printf("DbPassword: REDACTED\n")
			log.Printf("DbName: %s\n", variables.DbName)
		},
	}

	// Define cobra flags, the default value has the lowest (least significant) precedence
	rootCmd.Flags().StringVarP(&variables.AppConfig, "app-config", "c", "dev", "The application configuration")
	rootCmd.Flags().StringVarP(&variables.DbUsername, "db-username", "u", "Celtra", "The database username")
	rootCmd.Flags().StringVarP(&variables.DbPassword, "db-password", "p", "C3ltr4Ch4ll3ng3", "The database password")
	rootCmd.Flags().StringVarP(&variables.DbName, "db-name", "n", "UserData", "The database name")

	return rootCmd
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	// Set the base name of the config file, without the file extension.
	v.SetConfigName(defaultConfigFilename)

	// Set as many paths as you like where viper should look for the
	// config file. We are only looking in the current working directory.
	v.AddConfigPath(".")

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	v.SetEnvPrefix(envPrefix)

	// Environment variables can't have dashes in them, so bind them to their equivalent
	// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, v)

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Determine the naming convention of the flags when represented in the config file
		configName := f.Name
		// If using camelCase in the config file, replace hyphens with a camelCased string.
		// Since viper does case-insensitive comparisons, we don't need to bother fixing the case, and only need to remove the hyphens.
		if replaceHyphenWithCamelCase {
			configName = strings.ReplaceAll(f.Name, "-", "")
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
