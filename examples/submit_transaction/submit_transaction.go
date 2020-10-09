package main

import (
	"log"

	"github.com/tonicpow/go-minercraft"
)

func main() {

	// Create a new client
	client, err := minercraft.NewClient(nil, nil)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Select the miner
	miner := client.MinerByName("taal")

	// Submit transaction
	var response *minercraft.SubmitTransactionResponse
	if response, err = client.SubmitTransaction(
		miner,
		&minercraft.Transaction{RawTx: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000"},
	); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Display the results
	log.Printf("miner: %s", response.Miner.Name)
	log.Printf("status: %s [%s]", response.Submission.ReturnResult, response.Submission.ResultDescription)
	log.Printf("payload validated: %v", response.Validated)
}
