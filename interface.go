package minercraft

import (
	"context"
	"time"
)

// ClientInterface is the MinerCraft client interface
type ClientInterface interface {
	AddMiner(miner Miner) error
	BestQuote(ctx context.Context, feeCategory, feeType string) (*FeeQuoteResponse, error)
	FastestQuote(ctx context.Context, timeout time.Duration) (*FeeQuoteResponse, error)
	FeeQuote(ctx context.Context, miner *Miner) (*FeeQuoteResponse, error)
	MinerByID(minerID string) *Miner
	MinerByName(name string) *Miner
	Miners() []*Miner
	MinerUpdateToken(name, token string)
	PolicyQuote(ctx context.Context, miner *Miner) (*PolicyQuoteResponse, error)
	QueryTransaction(ctx context.Context, miner *Miner, txID string) (*QueryTransactionResponse, error)
	RemoveMiner(miner *Miner) bool
	SubmitTransaction(ctx context.Context, miner *Miner, tx *Transaction) (*SubmitTransactionResponse, error)
	SubmitTransactions(ctx context.Context, miner *Miner, txs []Transaction) (*SubmitTransactionsResponse, error)
	UserAgent() string
}
