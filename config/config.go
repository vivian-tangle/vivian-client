package config

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/vivian-tangle/vivian-client/tools"
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
	DefaultSecurityLevel = 2
	// DefaultDepth is the default depth of IOTA network
	DefaultDepth = 3
	// DefaultMinimumWeightMagnitude is the default minimum magnitude of IOTA network
	DefaultMinimumWeightMagnitude = 9
	// DefaultNodeClientVersion is the default node client version for lib-p2p network
	DefaultNodeClientVersion = "go-p2p-node/0.0.1"
	// DefaultPingRequestVersion is the default ping request version for lib-p2p network
	DefaultPingRequestVersion = "/ping/pingreq/0.0.1"
	// DefaultPingResponseVersion is the default ping response version for lip-p2p network
	DefaultPingResponseVersion = "/ping/pingresp/0.0.1"
	configDir                  = "config.json"
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
	NodeClientVersion      string
	PingRequestVersion     string
	PingResponseVersion    string
}

// LoadConfig loads the configures from default config json
func (c *Config) LoadConfig() {
	viper.AddConfigPath("./")
	viper.SetConfigFile(configDir)

	// Load default configurations
	viper.SetDefault("network", DefaultNetwork)
	viper.SetDefault("seedPath", DefaultSeedPath)
	viper.SetDefault("databasePath", DefaultDatabasePath)
	viper.SetDefault("node", DefaultNode)
	viper.SetDefault("securityLevel", DefaultSecurityLevel)
	viper.SetDefault("depth", DefaultDepth)
	viper.SetDefault("minimumWeightMagnitude", DefaultMinimumWeightMagnitude)
	viper.SetDefault("nodeClientVersion", DefaultNodeClientVersion)
	viper.SetDefault("pingRequestVersion", DefaultPingRequestVersion)
	viper.SetDefault("pingResponseVersion", DefaultPingResponseVersion)

	// Load the configurations from config file
	err := viper.ReadInConfig()
	tools.HandleErr(err)

	// unmarshal config
	err = viper.Unmarshal(&c)
	tools.HandleErr(err)

	fmt.Println("Configuration loaded")
}
