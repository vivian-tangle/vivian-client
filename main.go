package main

import (
	"fmt"

	"github.com/vivian-tangle/vivian-client/account"
	"github.com/vivian-tangle/vivian-client/config"
	"github.com/vivian-tangle/vivian-client/domain"
)

func main() {
	fmt.Println("Hello world!")
	c := config.Config{}
	c.LoadConfig()
	ac := account.Account{Seed: "", Config: &c}
	ac.GetSeed()
	ac.HelloWorldTx()
	// ac.ReadTxTagMsg("LGKZQJGPLRGRQQAQTVIWSRNBBUWNQBHGGCHQJNRPVPNBWXQXGFPSFMJKKFTIQCARNDEJI9FGDGSWVA999")
	dm := domain.Domain{Config: &c, Account: &ac}
	err := dm.PreorderName("vivian.vi")
	fmt.Println(err)
}
