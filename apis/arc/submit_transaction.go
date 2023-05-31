package arc

import "time"

type SubmitTxModel struct {
	BlockHash   string    `json:"blockHash,omitempty"`
	BlockHeight int64     `json:"blockHeight,omitempty"`
	ExtraInfo   string    `json:"extraInfo,omitempty"`
	Status      int       `json:"status,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
	Title       string    `json:"title,omitempty"`
	TxStatus    TxStatus  `json:"txStatus,omitempty"`
	TxID        string    `json:"txid,omitempty"`
}
