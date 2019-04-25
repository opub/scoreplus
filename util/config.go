package util

import (
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//Config application settings
type Config struct {
	Salt string
	DB   Database
	Log  Log
	Path Path
}

//Path config settings
type Path struct {
	StaticFiles string
	Templates   string
}

//Log config settings
type Log struct {
	Level   string // debug, info, warn, error, fatal, panic, none
	Console bool
	Caller  bool
}

//Database config settings
type Database struct {
	Name     string
	Host     string
	Port     int
	Username string
	Password string
}

var (
	config     Config
	onceConfig sync.Once
	testRE     = regexp.MustCompile("_test|.test.exe$|(\\.test$)")
)

//GetConfig gets application configuration settings based on current environment
func GetConfig() Config {
	onceConfig.Do(func() {
		setupConfig()
		config = loadConfig()
	})
	return config
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

func setupConfig() {
	env := getEnvironment()
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("config settings changed in %s\n", e.Name)
		config = loadConfig()
	})

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("failed to read %s config: %s\n", env, err)
		panic("config setup failed")
	}
}

func loadConfig() Config {
	settings := Config{}
	err := viper.Unmarshal(&settings)
	if err != nil {
		fmt.Printf("failed to unmarshal config: %s\n", err)
		panic("config loading failed")
	}
	return settings
}
