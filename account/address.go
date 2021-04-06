package account

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/iotaledger/iota.go/address"
	iotaAPI "github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/consts"
	"github.com/vivian-tangle/vivian-client/tools"
)

const (
	// AddressesDatabaseName is the name of the badger DB for storing the addresses
	AddressesDatabaseName = "addr"
	// LastIndex is the key of the last index in addresses database
	LastIndex = "LH"
)

// GetNewAddressFromAPI generates a new address based on the seed through IOTA API
func (ac *Account) GetNewAddressFromAPI(api *iotaAPI.API) []string {
	// Generate an unspent address with the security level
	// If this address is spent, this method returns the next unspent address with the lowest index
	secLvl := consts.SecurityLevel(ac.Config.SecurityLevel)
	address, err := api.GetNewAddress(ac.Seed, iotaAPI.GetNewAddressOptions{Security: secLvl})
	tools.HandleErr(err)

	return address
}

// GetNewAddressLocal generates a new unspent address based on the record of local database
func (ac *Account) GetNewAddressLocal() (string, error) {
	// Open the database for addresses
	_, err := os.Stat(ac.Config.DatabasePath)
	if os.IsNotExist(err) {
		os.MkdirAll(ac.Config.DatabasePath, os.ModePerm)
	}
	addressesDBPath := filepath.Join(ac.Config.DatabasePath, AddressesDatabaseName+strconv.Itoa(ac.Config.SecurityLevel))
	opts := badger.DefaultOptions(addressesDBPath)
	opts.Logger = nil // disable the message from badger log

	db, err := tools.OpenDB(addressesDBPath, opts)
	if err != nil {
		return "", err
	}
	defer db.Close()

	// Get the last index if exists
	var index uint64 = 0
	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(LastIndex))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				txn.Set([]byte(LastIndex), tools.Int2Byte(uint64(0)))
				return nil
			}
			return err
		}
		lastIndexByte, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		index = tools.ByteToInt(lastIndexByte) + 1
		err = txn.Set([]byte(LastIndex), tools.Int2Byte(uint64(index)))

		return err
	})
	if err != nil {
		return "", err
	}

	secLvl := consts.SecurityLevel(ac.Config.SecurityLevel)
	addr, err := address.GenerateAddress(ac.Seed, index, secLvl, true)
	tools.HandleErr(err)
	fmt.Println(addr)

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set(tools.Int2Byte(index), []byte(addr))
		return err
	})

	fmt.Printf("Index: %d, address:  %s\n", index, addr)

	return addr, err
}
