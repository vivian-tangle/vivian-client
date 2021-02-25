package main

import (
	"fmt"
	"log"
	"math/rand"

	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/vivian-tangle/vivian-client/config"
	"github.com/vivian-tangle/vivian-client/network"
)

func main() {
	fmt.Println("Hello world!")
	c := config.Config{}
	c.LoadConfig()
	// n := network.Network{Config: &c}
	// ac := account.Account{Seed: "", Config: &c}
	// ac.GetSeed()
	// ac.HelloWorldTx()
	// ac.ZeroValueTx("Hello world", domain.TagPreorder)
	// ac.ReadTxTagMsg("LGKZQJGPLRGRQQAQTVIWSRNBBUWNQBHGGCHQJNRPVPNBWXQXGFPSFMJKKFTIQCARNDEJI9FGDGSWVA999")
	// dm := domain.Domain{Config: &c, Account: &ac}
	// err := dm.PreorderName("vivian.vi")
	// if err != nil {
	// 	panic(err)
	// }
	// err := dm.RegisterName("vivian.vi", "Hello Vivian!")
	// if err != nil {
	// 	panic(err)
	// }
	nw := network.Network{Config: &c}
	rand.Seed(666)
	port1 := rand.Intn(100) + 10000
	port2 := port1 + 1

	done := make(chan bool, 1)

	h1 := nw.MakeRandomNode(port1, done)
	h2 := nw.MakeRandomNode(port2, done)

	h1.Peerstore().AddAddrs(h2.ID(), h2.Addrs(), peerstore.PermanentAddrTTL)
	h2.Peerstore().AddAddrs(h1.ID(), h1.Addrs(), peerstore.PermanentAddrTTL)

	log.Printf("This is a conversation between %s and %s\n", h1.ID(), h2.ID())

	// send messages using the protocols
	h1.Ping(h2.Host)
	h2.Ping(h1.Host)

	// block until all responses have been processed
	for i := 0; i < 4; i++ {
		<-done
	}
}
