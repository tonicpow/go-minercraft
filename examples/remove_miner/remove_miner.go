package main

import (
	"log"

	"github.com/tonicpow/go-minercraft/v2"
)

func main() {

	// Create a new client
	client, err := minercraft.NewClient(nil, nil, "", nil, nil)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Remove a miner
	client.RemoveMiner(client.MinerByName(minercraft.MinerTaal))

	// Show all miners loaded
	// TODO: Align with new structure
	// for _, miner := range client.Miners() {
	// 	log.Printf("miner: %s (%s)", miner.Name, miner.URL)
	// }
}
