package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jinzhu/configor"
)

// Config is a config :)
type Config struct {
	PgURL             string `env:"WAKEUP_PG_URL"`
	PgProto           string `env:"WAKEUP_PG_PROTO"`
	PgAddr            string `env:"WAKEUP_PG_ADDR"`
	PgDb              string `env:"WAKEUP_PG_DB"`
	PgUser            string `env:"WAKEUP_PG_USER"`
	PgPassword        string `env:"WAKEUP_PG_PASSWORD"`
	CallRoomEndpoint  string `env:"WAKEUP_CALL_ROOM_ENDPOINT"`
	SchedulerLocation string `env:"WAKEUP_SCHEDULER_LOCATION"`
	SchedulerTimeZone string `env:"WAKEUP_SCHEDULER_TIMEZONE"`
	CallEndpoint      string `env:"WAKEUP_CALL_ENDPOINT"`
}

var (
	config Config
	once   sync.Once
)

// Get reads config from environment
func Get() *Config {
	once.Do(func() {
		envType := os.Getenv("WAKEUP_ENV")
		if envType == "" {
			envType = "dev"
		}
		if err := configor.New(&configor.Config{Environment: envType}).Load(&config, "config.json"); err != nil {
			log.Fatal(err)
		}
		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Configuration:", string(configBytes))
	})
	return &config
}
