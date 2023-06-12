package arc

import "time"

// QueryTxModel is the unmarshalled version of the payload envelope
type QueryTxModel struct {
	BlockHash   string `json:"blockHash,omitempty"`
	BlockHeight int64  `json:"blockHeight,omitempty"`
	// TODO: Specify the type - currently no information on this in the docs
	ExtraInfo struct{}  `json:"extraInfo,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	TxStatus  TxStatus  `json:"txStatus,omitempty"`
	TxID      string    `json:"txid,omitempty"`
}
