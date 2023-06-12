package arc

import "github.com/libsv/go-bt/v2"

// Policy is the unmarshalled version of the payload envelope
type Policy struct {
	MaxScriptSizePolicy uint32  `json:"maxscriptsizepolicy"`
	MaxTxSigOpsCount    uint32  `json:"maxtxsigopscount"`
	MaxTxSizePolicy     uint32  `json:"maxtxsizepolicy"`
	MiningFee           *bt.Fee `json:"miningFee"`
}

// PolicyQuoteModel is the unmarshalled version of the payload envelope
type PolicyQuoteModel struct {
	Policy    Policy `json:"policy"`
	Timestamp string `json:"timestamp"`
}
