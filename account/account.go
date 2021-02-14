package account

import "github.com/vivian-tangle/vivian-client/config"

// Account is the structure for storing the account info
type Account struct {
	Seed   string
	Config config.Config
}
