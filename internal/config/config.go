package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address"`
}

type Config struct {
	Env         string               `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string               `yaml:"storage_path" env-required:"true" `
	HTTPServer  `yaml:"http_server"` //embedding of HTTPServer structure in Config Structure so that we can use it in Congif only
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		//it is possible that the onfig path is not given int the env and it is given as an argument for example in cmd line: go run main.go -config-path xyz then in this case we have to retrive that path also from the command

		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags //so flags will be parsed on the configPath variable and also it will be parsed as the pointer because it configPath is a string type variable and flags is a pointer not a string so we have to dereference it

		if configPath == "" {
			log.Fatal("Path not set") // that means config path is not given from both that is flag in the command line interface and from the ENV, so just return a Fatal error that path is not set yet
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) { //in this, os.Stat(configPath) will return some information about that file on that path, but we dont want that information so we haven't stored it in a particular so we used '_' identifier. so if error exists about file not exists then it will be stored in the err variable and then we will check if the path exists or not , if it doesn't exists then we will return the fatal and stop the application by returning a message that config path doesn't exists and the provided file path
		log.Fatalf("Config path does not exists: %s", configPath)
	}

	//now as of this step we have our valid configPath so we have to serialize config structure in a variable so we will do it using following:

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg) //it will read the path of the config file and will assign the configuration at the address of cfg which is of type Config only and it will return an error if any

	if err != nil {
		log.Fatalf("can not read config file: %s", err.Error()) //err.Error() will concat error message part in the log message
	}

	//so if we have reached till this step then we have our configuration ready and serialized in cfg variable so we can return its address

	return &cfg
}
