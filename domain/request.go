package domain

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dgraph-io/badger"
	"github.com/vivian-tangle/vivian-client/tools"
)

const (
	tagSuffix = "99999999999999999999999"

	// TagPreorderTrytes is the trytes for pre-order tag
	TagPreorderTrytes = "ZBYB" + tagSuffix
	// TagRegisterTrytes is the trytes for register tag
	TagRegisterTrytes = "ACQB" + tagSuffix
	// TagRenewTrytes is the trytes for renew tag
	TagRenewTrytes = "ACXB" + tagSuffix
	// TagUpdateTrytes is the trytes for update tag
	TagUpdateTrytes = "DCNB" + tagSuffix
	// TagTransferTrytes is the trytes for transfer tag
	TagTransferTrytes = "CCPB" + tagSuffix
	// TagRevokeTrytes is the trytes for revoke tag
	TagRevokeTrytes = "ACEC" + tagSuffix

	// TagPreorder is the tag for pre-ordering the domain
	TagPreorder = "PO"
	// TagRegister is the tag for registering the domain
	TagRegister = "RG"
	// TagRenew is the tag for renewing the domain
	TagRenew = "RN"
	// TagUpdate is the tag for updating the domain
	TagUpdate = "UD"
	// TagTransfer is the tag for transfering the domain
	TagTransfer = "TF"
	// TagRevoke is the tag for revoking the domain
	TagRevoke = "RV"
)

// PreorderName sends the transaction for preordering a name
func (d *Domain) PreorderName(name string) error {
	// Open the database
	_, err := os.Stat(d.Config.DatabasePath)
	if os.IsNotExist(err) {
		os.MkdirAll(d.Config.DatabasePath, os.ModePerm)
	}
	reserveDBPath := filepath.Join(d.Config.DatabasePath, ReservedDatabaseName)
	opts := badger.DefaultOptions(reserveDBPath)
	opts.Logger = nil // disable the message from badger log

	db, err := tools.OpenDB(reserveDBPath, opts)
	if err != nil {
		return err
	}
	defer db.Close()

	// Check if the domain name is already pre-ordered
	err = db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(name))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			return err
		}
		errMsg := fmt.Sprintf("Name: %s already reserved!", name)
		return errors.New(errMsg)
	})
	if err != nil {
		return err
	}

	fmt.Printf("Preordering name: %s\n", name)
	// Pederson commitment
	g, h := tools.GenerateParametersToString()
	r := tools.GenerateRandomToString()
	commit, err := tools.CommitByString(g, h, r, []byte(name))
	if err != nil {
		return err
	}
	fmt.Printf("Pederson commitment: %s\n", commit)

	txHash, err := d.Account.ZeroValueTx(commit, TagPreorder)
	if err != nil {
		return err
	}

	pc := tools.PedersonCommit{
		Content: name,
		G:       g,
		H:       h,
		R:       r,
		Commit:  commit,
		TxHash:  txHash,
	}

	fmt.Println(pc)

	// Put the information to the database
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(name), pc.Serialize())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// RegisterName sends the transaction for registering a name
func (d *Domain) RegisterName(name, value string) error {
	// Open the pre-order database to check if the domain name is pre-ordered and retrive the information of pre-ordered domain name
	_, err := os.Stat(d.Config.DatabasePath)
	if os.IsNotExist(err) {
		os.MkdirAll(d.Config.DatabasePath, os.ModePerm)
		errMsg := fmt.Sprintf("Database path %s does not exist!", d.Config.DatabasePath)
		return errors.New(errMsg)
	}
	reserveDBPath := filepath.Join(d.Config.DatabasePath, ReservedDatabaseName)
	reserveDBOpts := badger.DefaultOptions(reserveDBPath)
	reserveDBOpts.Logger = nil // disable the message from badger log

	reserveDB, err := tools.OpenDB(reserveDBPath, reserveDBOpts)
	if err != nil {
		return err
	}
	defer reserveDB.Close()

	var preorderInfoByte []byte
	err = reserveDB.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(name))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				errMsg := fmt.Sprintf("Name: %s hasn't been pre-ordered yet!", name)
				return errors.New(errMsg)
			}
			return err
		}
		preorderInfoByte, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return err
	}

	var preorderInfo tools.PedersonCommit
	msg := encodeRegisterJSON(preorderInfo.Deserialize(preorderInfoByte), name, value)
	fmt.Println(msg)

	_, err = d.Account.ZeroValueTx(msg, TagRegister)

	return err
}

// RenewName sends the transaction for renewing a name
func (d *Domain) RenewName() {}

// UpdateName sends the transaction for updating a name
func (d *Domain) UpdateName() {}

// TransferName sends the transaction for transfering a name
func (d *Domain) TransferName() {}

// RevokeName sends the transaction for recoking a name
func (d *Domain) RevokeName() {}

func encodeRegisterJSON(pc *tools.PedersonCommit, name, value string) string {
	var res string
	res = fmt.Sprintf("{'hash': '%s', 'g': '%s', 'h': '%s','r': '%s', 'name': '%s', 'value': '%s'}",
		pc.TxHash, pc.G, pc.H, pc.R, name, value)

	return res
}
