package main

import (
	"fmt"

	"github.com/vivian-tangle/vivian-client/account"
	"github.com/vivian-tangle/vivian-client/config"
)

func main() {
	fmt.Println("Hello world!")
	c := config.Config{}
	c.LoadConfig()
	ac := account.Account{Seed: "", Config: c}
	// ac.GetSeed()
	// ac.HelloWorldTx()
	ac.ReadTxTagMsg("LGKZQJGPLRGRQQAQTVIWSRNBBUWNQBHGGCHQJNRPVPNBWXQXGFPSFMJKKFTIQCARNDEJI9FGDGSWVA999")
}
