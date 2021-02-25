package domain

import (
	"github.com/dgraph-io/badger"
	"github.com/vivian-tangle/vivian-client/account"
	"github.com/vivian-tangle/vivian-client/config"
)

const (
	// ReservedDatabaseName is the name of the badger DB of reserved domain names
	ReservedDatabaseName = "rv"
	// PendingDatabaseName is the name of the badfer DB of pending domain names
	PendingDatabaseName = "pd"
	// RegisteredDatabaseName is the name of the badger DB of registered domain names
	RegisteredDatabaseName = "rs"
)

// Domain is the structure for storing the info of domain names belong to the account
type Domain struct {
	Config     *config.Config
	Account    *account.Account
	ReserveDB  *badger.DB
	RegisterDB *badger.DB
}
