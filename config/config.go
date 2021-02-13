package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	// DefaultNetwork is the default parameter for choosing IOTA network (mainnet/devnet)
	DefaultNetwork = "devnet"
	// DefaultSeedPath is the directory for storing the seeds
	DefaultSeedPath = "seeds"
	// DefaultDatabasePath is directory for storing database files
	DefaultDatabasePath = "db"
	configDir           = "config.json"
)

// Config is the struct for storing config parameters
type Config struct {
	Network      string
	SeedPath     string
	DatabasePath string
}

// LoadConfig loads the configures from default config json
func (c *Config) LoadConfig() {
	c.Network = DefaultNetwork
	c.SeedPath = DefaultSeedPath
	c.DatabasePath = DefaultDatabasePath

	configJSON, err := os.Open(configDir)
	if err != nil {
		fmt.Printf("Cannot load %s, using default configure values...", configDir)
		return
	}
	defer configJSON.Close()

	// Read the opened config json as a byte array
	byteValue, _ := ioutil.ReadAll(configJSON)
	var data map[string]string
	err = json.Unmarshal(byteValue, &data)
	handleErr(err)
	if val, ok := data["network"]; ok {
		c.Network = val
	}
	if val, ok := data["seedPath"]; ok {
		c.SeedPath = val
	}
	if val, ok := data["databasePath"]; ok {
		c.DatabasePath = val
	}

	fmt.Println("Configuration loaded")
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
