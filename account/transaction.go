package account

import (
	"fmt"

	iotaAPI "github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/bundle"
	"github.com/iotaledger/iota.go/converter"
	"github.com/vivian-tangle/vivian-client/tools"
)

// GetNewAddress generates a new address based on the seed
func (ac *Account) GetNewAddress(api *iotaAPI.API) []string {
	// Generate an unspent address with the security level
	// If this address is spent, this method returns the next unspent address with the lowest index
	address, err := api.GetNewAddress(ac.Seed, ac.MakeNewAddressOptions())
	tools.HandleErr(err)

	return address
}

// MakeNewAddressOptions generates GetNewAddressOptions struct by security level
func (ac *Account) MakeNewAddressOptions() iotaAPI.GetNewAddressOptions {
	var newAddressOption iotaAPI.GetNewAddressOptions
	switch sl := ac.Config.SecurityLevel; sl {
	case 1:
		newAddressOption = iotaAPI.GetNewAddressOptions{Security: 1}
	case 2:
		newAddressOption = iotaAPI.GetNewAddressOptions{Security: 2}
	case 3:
		newAddressOption = iotaAPI.GetNewAddressOptions{Security: 3}
	}

	return newAddressOption
}

// HelloWorldTx sends a "Hello World" transaction
func (ac *Account) HelloWorldTx() {
	// Connect to a node
	node := ac.Config.Node
	api, err := iotaAPI.ComposeAPI(iotaAPI.HTTPClientSettings{URI: node})
	tools.HandleErr(err)
	address := ac.GetNewAddress(api)
	var data = "{'message' : 'Hello world'}"
	message, err := converter.ASCIIToTrytes(data)
	tools.HandleErr(err)
	transfers := bundle.Transfers{
		{
			Address: address[0],
			Value:   0,
			Message: message,
		},
	}
	trytes, err := api.PrepareTransfers(ac.Seed, transfers, iotaAPI.PrepareTransfersOptions{})
	tools.HandleErr(err)

	depth := ac.Config.Depth
	minimumWeightMagnitude := ac.Config.MinimumWeightMagnitude
	myBundle, err := api.SendTrytes(trytes, depth, minimumWeightMagnitude)
	tools.HandleErr(err)

	fmt.Println(bundle.TailTransactionHash(myBundle))
}
