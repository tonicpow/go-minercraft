package main

import (
	"context"
	"log"

	"github.com/tonicpow/go-minercraft"
	"github.com/tonicpow/go-minercraft/apis/mapi"
)

func main() {

	// Create a new client
	client, err := minercraft.NewClient(nil, nil, "", nil, nil)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	log.Printf("querying %d miners for the best rate...", len(client.Miners()))

	// Fetch quotes from all miners
	var response *minercraft.FeeQuoteResponse
	response, err = client.BestQuote(context.Background(), mapi.FeeCategoryMining, mapi.FeeTypeData)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	log.Printf("found best quote: %s", response.Miner.Name)
}
