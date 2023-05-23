package arc

import "time"

type TxStatus int

const (
	// List of statuses available here: https://github.com/bitcoin-sv/arc
	UNKNOWN              TxStatus = iota // 0
	QUEUED                               // 1
	RECEIVED                             // 2
	STORED                               // 3
	ANNOUNCED_TO_NETWORK                 // 4
	REQUESTED_BY_NETWORK                 // 5
	SENT_TO_NETWORK                      // 6
	ACCEPTED_BY_NETWORK                  // 7
	SEEN_ON_NETWORK                      // 8
	MINED                                // 9
	CONFIRMED            TxStatus = 108  // 108
	REJECTED             TxStatus = 109  // 109
)

func (s TxStatus) String() string {
	statuses := [...]string{
		"UNKOWN",
		"QUEUED",
		"RECEIVED",
		"SOTRED",
		"ANNOUNCED_TO_NETWORK",
		"REQUESTED_BY_NETWORK",
		"SENT_TO_NETWORK",
		"ACCEPTED_BY_NETWORK",
		"SEEN_ON_NETWORK",
		"MINED",
		"CONFIRMED",
		"REJECTED",
	}

	if s < UNKNOWN || s > REJECTED {
		return "Can't parse status"
	}

	return statuses[s]
}

type QueryTxModel struct {
	BlockHash   string `json:"blockHash,omitempty"`
	BlockHeight int64  `json:"blockHeight,omitempty"`
	// TODO: Specify the type - currently no information on this in the docs
	ExtraInfo struct{}  `json:"extraInfo,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	TxStatus  TxStatus  `json:"txStatus,omitempty"`
	TxID      string    `json:"txid,omitempty"`
}
