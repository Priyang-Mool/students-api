package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// HTTPServer represents the configuration for an HTTP server.
type HTTPServer struct {
	// Addr is the address of the HTTP server.
	Addr string `yaml:"address"`
}

// Config represents the application configuration.
type Config struct {
	// Env is the environment in which the application is running.
	Env string `yaml:"env" env:"ENV" env-required:"true"`
	// StoragePath is the path to the storage directory.
	StoragePath string `yaml:"storage_path" env-required:"true" `
	// HTTPServer is the embedded HTTP server configuration.
	HTTPServer `yaml:"http_server"` //embedding of HTTPServer structure in Config Structure so that we can use it in Congif only
}

// MustLoad loads the application configuration from a file.
// It returns a pointer to the loaded configuration.
// If the configuration cannot be loaded, it logs a fatal error and exits the application.
func MustLoad() *Config {
	// Get the configuration file path from the environment variable CONFIG_PATH.
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	// If the configuration file path is not set in the environment variable,
	// try to get it from the command-line flag -config.
	if configPath == "" {
		// Define a command-line flag for the configuration file path.
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		// Get the configuration file path from the command-line flag.
		configPath = *flags //so flags will be parsed on the configPath variable and also it will be parsed as the pointer because it configPath is a string type variable and flags is a pointer not a string so we have to dereference it

		// If the configuration file path is still not set, log a fatal error and exit the application.
		if configPath == "" {
			log.Fatal("Path not set") // that means config path is not given from both that is flag in the command line interface and from the ENV, so just return a Fatal error that path is not set yet
		}
	}

	// Check if the configuration file exists.
	if _, err := os.Stat(configPath); os.IsNotExist(err) { //in this, os.Stat(configPath) will return some information about that file on that path, but we dont want that information so we haven't stored it in a particular so we used '_' identifier. so if error exists about file not exists then it will be stored in the err variable and then we will check if the path exists or not , if it doesn't exists then we will return the fatal and stop the application by returning a message that config path doesn't exists and the provided file path
		// If the configuration file does not exist, log a fatal error and exit the application.
		log.Fatalf("Config path does not exists: %s", configPath)
	}

	// Create a new Config instance to store the loaded configuration.
	var cfg Config

	// Load the configuration from the file using the cleanenv package.
	err := cleanenv.ReadConfig(configPath, &cfg) //it will read the path of the config file and will assign the configuration at the address of cfg which is of type Config only and it will return an error if any

	// If an error occurs while loading the configuration, log a fatal error and exit the application.
	if err != nil {
		log.Fatalf("can not read config file: %s", err.Error()) //err.Error() will concat error message part in the log message
	}

	// Return the loaded configuration.
	return &cfg
}