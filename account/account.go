package account

import (
	"fmt"

	"github.com/iotaledger/iota.go/account"
	"github.com/iotaledger/iota.go/account/builder"
	"github.com/iotaledger/iota.go/account/store/badger"
	"github.com/iotaledger/iota.go/account/timesrc"
	"github.com/iotaledger/iota.go/api"
	"github.com/vivian-tangle/vivian-client/config"
	"github.com/vivian-tangle/vivian-client/tools"
)

// Account is the structure for storing the account info
type Account struct {
	Seed        string
	IotaAccount account.Account
	Config      *config.Config
}

func (ac *Account) Init() {
	// Define the node to connect to
	apiSettings := api.HTTPClientSettings{URI: "https://nodes.devnet.iota.org:443"}

	iotaAPI, err := api.ComposeAPI(apiSettings)
	tools.HandleErr(err)

	// Define a database in which to store the seed state
	store, err := badger.NewBadgerStore("seed-state-database")
	tools.HandleErr(err)

	// Make sure the database closes when the code stops
	defer store.Close()

	// Use the Google NTP servers as a reliable source of time to check CDA timeouts
	timesource := timesrc.NewNTPTimeSource("time.google.com")

	account, err := builder.NewBuilder().
		// Connect to a node
		WithAPI(iotaAPI).
		// Connect to the database
		WithStore(store).
		// Load the seed
		WithSeed(ac.Seed).
		// Set the minimum weight magnitude for the Devnet (default is 14)
		WithMWM(9).
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
