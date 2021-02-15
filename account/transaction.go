package account

import (
	"fmt"

	iotaAPI "github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/bundle"
	"github.com/iotaledger/iota.go/converter"
	"github.com/iotaledger/iota.go/transaction"
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
	// Load configurations
	node := ac.Config.Node
	depth := ac.Config.Depth
	minimumWeightMagnitude := ac.Config.MinimumWeightMagnitude

	// Connect to a node
	api, err := iotaAPI.ComposeAPI(iotaAPI.HTTPClientSettings{URI: node})
	tools.HandleErr(err)

	address := ac.GetNewAddress(api)
	var data = "{'message' : 'Hello world'}"
	tag, err := converter.ASCIIToTrytes("HELLO")
	message, err := converter.ASCIIToTrytes(data)
	tools.HandleErr(err)

	transfers := bundle.Transfers{
		{
			Address: address[0],
			Value:   0,
			Tag:     tag,
			Message: message,
		},
	}

	trytes, err := api.PrepareTransfers(ac.Seed, transfers, iotaAPI.PrepareTransfersOptions{})
	tools.HandleErr(err)

	myBundle, err := api.SendTrytes(trytes, depth, minimumWeightMagnitude)
	tools.HandleErr(err)

	fmt.Printf("Transaction sent with tail tx hash:\n%s\n", bundle.TailTransactionHash(myBundle))
}

// ZeroValueTx sends a zero value transaction for conveying the message
func (ac *Account) ZeroValueTx(message, tag string) error {
	// Load configurations
	node := ac.Config.Node
	depth := ac.Config.Depth
	minimumWeightMagnitude := ac.Config.MinimumWeightMagnitude

	// Connect to a node
	api, err := iotaAPI.ComposeAPI(iotaAPI.HTTPClientSettings{URI: node})
	if err != nil {
		return err
	}

	address := ac.GetNewAddress(api)
	messageTrytes, err := converter.ASCIIToTrytes(message)
	if err != nil {
		return err
	}
	tagTrytes, err := converter.ASCIIToTrytes(tag)
	if err != nil {
		return err
	}

	transfers := bundle.Transfers{
		{
			Address: address[0],
			Value:   0,
			Tag:     tagTrytes,
			Message: messageTrytes,
		},
	}

	trytes, err := api.PrepareTransfers(ac.Seed, transfers, iotaAPI.PrepareTransfersOptions{})
	if err != nil {
		return err
	}

	myBundle, err := api.SendTrytes(trytes, depth, minimumWeightMagnitude)
	if err != nil {
		return err
	}

	fmt.Printf("Transaction sent with tail tx hash:\n%s\n", bundle.TailTransactionHash(myBundle))

	return nil
}

// ReadTxTagMsg reads the transaction tag and message by tail transaction hash
func (ac *Account) ReadTxTagMsg(tailTxHash string) (string, string) {
	// load configurations
	node := ac.Config.Node

	// Connect to a node
	api, err := iotaAPI.ComposeAPI(iotaAPI.HTTPClientSettings{URI: node})
	tools.HandleErr(err)

	bundle, err := api.GetBundle(tailTxHash)
	tools.HandleErr(err)

	jsonMsg, err := transaction.ExtractJSON(bundle)
	tools.HandleErr(err)

	fmt.Println(bundle[0].Tag)
	fmt.Println(jsonMsg)
	return bundle[0].Tag, jsonMsg
}
