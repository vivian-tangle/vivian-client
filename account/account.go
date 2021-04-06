package account

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/iotaledger/iota.go/account"
	"github.com/iotaledger/iota.go/account/builder"
	"github.com/iotaledger/iota.go/account/store/badger"
	"github.com/iotaledger/iota.go/account/timesrc"
	"github.com/iotaledger/iota.go/api"
	"github.com/vivian-tangle/vivian-client/config"
	"github.com/vivian-tangle/vivian-client/tools"
)

const (
	// SeedStateDatabasePath is the name of the badger DB for storing the seed state
	SeedStateDatabaseName = "seed_state"
)

// Account is the structure for storing the account info
type Account struct {
	Seed        string
	IotaAccount account.Account
	Config      *config.Config
}

func (ac *Account) Init() {
	// Define the node to connect to
	apiSettings := api.HTTPClientSettings{URI: ac.Config.Node}

	iotaAPI, err := api.ComposeAPI(apiSettings)
	tools.HandleErr(err)

	// Create the database path if it does not exist
	_, err = os.Stat(ac.Config.DatabasePath)
	if os.IsNotExist(err) {
		os.MkdirAll(ac.Config.DatabasePath, os.ModePerm)
	}
	// Define a database in which to store the seed state
	seedStateDatabaseName := filepath.Join(ac.Config.DatabasePath, SeedStateDatabaseName)
	store, err := badger.NewBadgerStore(seedStateDatabaseName)
	tools.HandleErr(err)

	// Make sure the database closes when the code stops
	defer store.Close()

	// Use a reliable source of time to check CDA timeouts
	timesource := timesrc.NewNTPTimeSource(ac.Config.NTPTimeSource)

	account, err := builder.NewBuilder().
		// Connect to a node
		WithAPI(iotaAPI).
		// Connect to the database
		WithStore(store).
		// Load the seed
		WithSeed(ac.Seed).
		// Set the minimum weight magnitude for the Devnet (default is 14)
		WithMWM(ac.Config.MinimumWeightMagnitude).
		// Use a reliable time source
		WithTimeSource(timesource).
		// Load the default plugins that enhance the functionality of the account
		WithDefaultPlugins().
		Build()
	tools.HandleErr(err)

	tools.HandleErr(account.Start())

	// Make sure the account shuts down when the code stops
	defer account.Shutdown()

	balance, err := account.AvailableBalance()
	tools.HandleErr(err)
	fmt.Println("Total available balance: ")
	fmt.Println(balance)
}

// CheckBalance gets the total available balance of the account
func (ac *Account) CheckBalance() {

}
