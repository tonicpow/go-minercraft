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
	miner := client.MinerByName(minercraft.MinerTaal)

	// Transactions
	var transactions []*minercraft.Transaction
	transactions = append(transactions, &minercraft.Transaction{RawTx: "0100000001d6d1607b208b30c0a3fe21d563569c4d2a0f913604b4c5054fe267da6be324ab220000006b4830450221009a965dcd5d42983090a63cfd761038ff8adcea621c46a68a205f326292a95383022061b8d858f366c69f3ebd30a60ccafe36faca4e242ac3d2edd3bf63b669bcf23b4121034e871e147aa4a3e2f1665eaf76cf9264d089b6a91702af92bd6ce33bac84a765ffffffff0123020000000000001976a914d8819a7197d3e221e15f4348203fdecfd29fa2b888ac00000000"})

	// Submit transaction
	var response *minercraft.SubmitTransactionsResponse
	if response, err = client.SubmitTransactions(
		miner,
		transactions,
	); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// todo: issue with request - returning 404

	// Display the results
	log.Printf("miner: %s", response.Miner.Name)
	log.Printf("payload validated: %v", response.Validated)
}
