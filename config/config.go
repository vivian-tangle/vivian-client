package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var (
	// DefaultNetwork is the default parameter for choosing IOTA network (mainnet/devnet)
	DefaultNetwork = "devnet"
	// DefaultSeedPath is the directory for storing the seeds
	DefaultSeedPath = "seeds"
	// DefaultDatabasePath is directory for storing database files
	DefaultDatabasePath = "db"
	// DefaultNode is the default node for connecting to IOTA network
	DefaultNode = "https://nodes.devnet.iota.org"
	// DefaultSecurityLevel is the default security level you want to use for your address
	DefaultSecurityLevel = "2"
	// DefaultDepth is the default depth of IOTA network
	DefaultDepth = "3"
	// DefaultMinimumWeightMagnitude is the default minimum magnitude of IOTA network
	DefaultMinimumWeightMagnitude = "9"
	configDir                     = "config.json"
)

// Config is the struct for storing config parameters
type Config struct {
	Network                string
	SeedPath               string
	DatabasePath           string
	Node                   string
	SecurityLevel          int
	Depth                  uint64
	MinimumWeightMagnitude uint64
}

// LoadConfig loads the configures from default config json
func (c *Config) LoadConfig() {
	// Load default configurations
	c.Network = DefaultNetwork
	c.SeedPath = DefaultSeedPath
	c.DatabasePath = DefaultDatabasePath
	c.Node = DefaultNode
	c.SecurityLevel, _ = strconv.Atoi(DefaultSecurityLevel)
	c.Depth, _ = strconv.ParseUint(DefaultDepth, 0, 64)
	c.MinimumWeightMagnitude, _ = strconv.ParseUint(DefaultMinimumWeightMagnitude, 0, 64)

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
	if val, ok := data["node"]; ok {
		c.Node = val
	}
	if val, ok := data["securityLevel"]; ok {
		c.SecurityLevel, _ = strconv.Atoi(val)
	}
	if val, ok := data["depth"]; ok {
		c.Depth, _ = strconv.ParseUint(val, 0, 64)
	}
	if val, ok := data["minimumWeightMagnitude"]; ok {
		c.MinimumWeightMagnitude, _ = strconv.ParseUint(val, 0, 64)
	}

	fmt.Println("Configuration loaded")
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
