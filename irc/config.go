package irc

import (
	"encoding/json"
	"os"
)

// Config contains all of the configuration settings required to bring up a
// local irc server.
type Config struct {
	Name    string
	Network string
	Port    int
	MOTD    string

	DefaultChannelMode string
	DefaultUserMode    string

	SpoofHostName string

	PingFrequency  int
	PongMaxLatency int
}

// ConfigFromJSONFile reads a Config struct from a file containing a JSON
// encoded value.
func ConfigFromJSONFile(path string) Config {
	var cfg Config
	file, err := os.Open(path)
	if err != nil {
		logf(fatal, "Could not open config file: %v.", err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		logf(fatal, "Problem parsing config data: %v", err)
	}

	return setConfigDefaults(cfg)
}

// setConfigDefaults fills in the default values of the Config if no value is
// specified for a field.
func setConfigDefaults(cfg Config) Config {
	if cfg.PingFrequency == 0 {
		cfg.PingFrequency = 30
	}

	if cfg.PongMaxLatency == 0 {
		cfg.PongMaxLatency = 5
	}

	return cfg
}
