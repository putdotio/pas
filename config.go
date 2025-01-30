package main

import (
	"github.com/BurntSushi/toml"
	"github.com/kelseyhightower/envconfig"
	"github.com/putdotio/pas/internal/event"
	"github.com/putdotio/pas/internal/property"
)

// Config for application
type Config struct {
	// Listen address for HTTP server.
	ListenAddress string
	// Give some time to unfinished HTTP requests before shutting down the server (milliseconds).
	ShutdownTimeout uint
	// MySQL database DSN.
	MySQLDSN string
	// Secret for signing user IDs.
	Secret string
	// Corresponds to the schema of user table.
	User property.Types
	// Corresponds to the schema of event tables.
	Events map[event.Name]property.Types
}

func NewConfig() (*Config, error) {
	c := new(Config)
	_, err := toml.DecodeFile(*configPath, c)
	if err != nil {
		return nil, err
	}
	err = envconfig.Process("PAS", c)
	if err != nil {
		return nil, err
	}
	c.setDefaults()
	return c, nil
}

func (c *Config) setDefaults() {
	if c.ListenAddress == "" {
		c.ListenAddress = "127.0.0.1:8080"
	}
	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = 5000
	}
}
