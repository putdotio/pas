package main

import (
	"github.com/BurntSushi/toml"
	"github.com/kelseyhightower/envconfig"
)

// Config for application
type Config struct {
	// Listen address for HTTP server.
	ListenAddress string
	// Give some time to unfinished HTTP requests before shutting down the server (milliseconds).
	ShutdownTimeout uint
	// MySQL database DSN.
	MySQLDSN string
}

func (c *Config) Read() error {
	_, err := toml.DecodeFile(*configPath, c)
	if err != nil {
		return err
	}
	err = envconfig.Process("", c)
	if err != nil {
		return err
	}
	c.setDefaults()
	return nil
}

func (c *Config) setDefaults() {
	if c.ListenAddress == "" {
		c.ListenAddress = "127.0.0.1:8080"
	}
	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = 5000
	}
}
