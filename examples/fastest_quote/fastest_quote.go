package main

import (
	"log"

	"github.com/tonicpow/go-minercraft"
)

func main() {

	// Create a new client
	client, err := minercraft.NewClient(nil, nil, nil)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	log.Printf("querying %d miners for the fastest response...", len(client.Miners))

	// Fetch fastest quote from all miners
	var response *minercraft.FeeQuoteResponse
	response, err = client.FastestQuote()
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	log.Printf("found quote: %s", response.Miner.Name)
}
