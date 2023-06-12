package main

import (
	"log"

	"github.com/tonicpow/go-minercraft/v2"
)

func main() {
	apiType := minercraft.MAPI
	// Create a new client
	client, err := minercraft.NewClient(nil, nil, apiType, nil, nil)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Add a custom miner!
	if err = client.AddMiner(minercraft.Miner{
		Name: "Custom",
	}, []minercraft.API{{URL: "https://mapi.customminer.com", Type: apiType}}); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Show all miners loaded
	for _, miner := range client.Miners() {
		var api *minercraft.API
		api, err = client.MinerAPIByMinerID(miner.MinerID, apiType)
		if err != nil {
			log.Fatalf("error occurred: %s", err.Error())
		}
		log.Printf("miner: %s (%s)", miner.Name, api.URL)
	}
}
