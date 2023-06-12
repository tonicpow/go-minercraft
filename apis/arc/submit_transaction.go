package arc

import "time"

// SubmitTxModel is the unmarshalled version of the payload envelope
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
