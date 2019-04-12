package util

import (
	"fmt"
	"os"
	"regexp"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//Config application settings
type Config struct {
	Salt string
	DB   Database
}

//Database config settings
type Database struct {
	Name     string
	Host     string
	Port     int
	Username string
	Password string
}

var config Config
var testRE = regexp.MustCompile("_test|.test.exe$|(\\.test$)")

//GetConfig gets application configuration settings based on current environment
func GetConfig() (Config, error) {
	if config == (Config{}) {
		var err error
		config, err = loadConfig()
		if err != nil {
			return config, err
		}
	}
	return config, nil
}

func getEnvironment() string {
	if env := os.Getenv("SCOREPLUS_ENV"); env != "" {
		return env
	}
	if testRE.MatchString(os.Args[0]) {
		return "test"
	}
	return "dev"
}

func loadConfig() (Config, error) {
	env := getEnvironment()
	fmt.Printf("ENV: %s\n", env)
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed: ", e.Name)
		viper.Unmarshal(&config)
	})

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("failed to read %s config: %s", env, err)
	}

	settings := Config{}
	viper.Unmarshal(&settings)
	return settings, err
}
