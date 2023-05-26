package arc

import "github.com/libsv/go-bt/v2"

type Policy struct {
	MaxScriptSizePolicy uint32  `json:"maxscriptsizepolicy"`
	MaxTxSigOpsCount    uint32  `json:"maxtxsigopscount"`
	MaxTxSizePolicy     uint32  `json:"maxtxsizepolicy"`
	MiningFee           *bt.Fee `json:"miningFee"`
}

type PolicyQuoteModel struct {
	Policy    Policy `json:"policy"`
	Timestamp string `json:"timestamp"`
}
