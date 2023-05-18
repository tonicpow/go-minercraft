package main

import (
	"context"
	"log"

	"github.com/tonicpow/go-minercraft"
)

func main() {

	// Create a new client
	client, err := minercraft.NewClient(nil, nil, minercraft.Arc, nil, nil)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Select the miner
	miner := client.MinerByName(minercraft.MinerGorillaPool)

	// Query the transaction status
	var response *minercraft.QueryTransactionResponse
	if response, err = client.QueryTransaction(context.Background(), miner, "9c5f5244ee45e8c3213521c1d1d5df265d6c74fb108961a876917073d65fef14"); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Display the results
	log.Printf("miner: %s", response.Miner.Name)
	log.Printf("status: %s [%s]", response.Query.ReturnResult, response.Query.ResultDescription)
	log.Printf("payload validated: %v", response.Validated)
}
