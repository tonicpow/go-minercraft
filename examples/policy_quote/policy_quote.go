package main

import (
	"context"
	"log"

	"github.com/tonicpow/go-minercraft"
)

func main() {

	// Create a new client
	client, err := minercraft.NewClient(nil, nil, "", nil, nil)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Select the miner
	miner := client.MinerByName(minercraft.MinerTaal)

	// Get a policy quote from a miner
	var response *minercraft.PolicyQuoteResponse
	if response, err = client.PolicyQuote(context.Background(), miner); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Display the results
	// todo: At this time (1/10/21) there is no policy response from any miner
	log.Printf("miner: %s", response.Miner.Name)
	log.Printf("callbacks: %+v", response.Quote.Callbacks)
	log.Printf("policy quote: %+v", response.Quote.Policies)
}
