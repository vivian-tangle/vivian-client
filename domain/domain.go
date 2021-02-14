package domain

import (
	"github.com/vivian-tangle/vivian-client/account"
	"github.com/vivian-tangle/vivian-client/config"
)

// Domain is the structure for storing the info of domain names belong to the account
type Domain struct {
	Config  *config.Config
	Account *account.Account
}
