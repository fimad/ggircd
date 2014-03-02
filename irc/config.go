package irc

import (
  "encoding/json"
  "log"
  "os"
)

// Config contains all of the configuration settings required to bring up a
// local irc server.
type Config struct {
  Name    string
  Network string
  Port    int
  MOTD    string
}

// ConfigFromJSONFile reads a Config struct from a file containing a JSON
// encoded value.
func ConfigFromJSONFile(path string) Config {
  var cfg Config
  file, err := os.Open(path)
  if err != nil {
    log.Fatalf("Could not open config file: %v.", err)
  }

  decoder := json.NewDecoder(file)
  err = decoder.Decode(&cfg)
  if err != nil {
    log.Fatalf("Problem parsing config data: %v", err)
  }

  return cfg
}
