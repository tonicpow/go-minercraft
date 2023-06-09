package mapi

/*
Example policyQuote response from Merchant API:

{
    "payload": "{\"apiVersion\":\"1.4.0\",\"timestamp\":\"2021-11-12T13:17:47.7498672Z\",\"expiryTime\":\"2021-11-12T13:27:47.7498672Z\",\"minerId\":\"030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e\",\"currentHighestBlockHash\":\"45628be2fe616167b7da399ab63455e60ffcf84147730f4af4affca90c7d437e\",\"currentHighestBlockHeight\":234,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}}],\"callbacks\":[{\"ipAddress\":\"123.456.789.123\"}],\"policies\":{\"skipscriptflags\":[\"MINIMALDATA\",\"DERSIG\",\"NULLDUMMY\",\"DISCOURAGE_UPGRADABLE_NOPS\",\"CLEANSTACK\"],\"maxtxsizepolicy\":99999,\"datacarriersize\":100000,\"maxscriptsizepolicy\":100000,\"maxscriptnumlengthpolicy\":100000,\"maxstackmemoryusagepolicy\":10000000,\"limitancestorcount\":1000,\"limitcpfpgroupmemberscount\":10,\"acceptnonstdoutputs\":true,\"datacarrier\":true,\"dustrelayfee\":150,\"maxstdtxvalidationduration\":99,\"maxnonstdtxvalidationduration\":100,\"dustlimitfactor\":10}}",
    "signature": "30440220708e2e62a393f53c43d172bc1459b4daccf9cf23ff77cff923f09b2b49b94e0a022033792bee7bc3952f4b1bfbe9df6407086b5dbfc161df34fdee684dc97be72731",
    "publicKey": "030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e",
    "encoding": "UTF-8",
    "mimetype": "application/json"
}
*/

/*
Example PolicyQuoteResponse.Payload (unmarshalled):

{
    "apiVersion": "1.4.0",
    "timestamp": "2021-11-12T13:17:47.7498672Z",
    "expiryTime": "2021-11-12T13:27:47.7498672Z",
    "minerId": "030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e",
    "currentHighestBlockHash": "45628be2fe616167b7da399ab63455e60ffcf84147730f4af4affca90c7d437e",
    "currentHighestBlockHeight": 234,
    "fees": [
        {
            "feeType": "standard",
            "miningFee": {
                "satoshis": 500,
                "bytes": 1000
            },
            "relayFee": {
                "satoshis": 250,
                "bytes": 1000
            }
        },
        {
            "feeType": "data",
            "miningFee": {
                "satoshis": 500,
                "bytes": 1000
            },
            "relayFee": {
                "satoshis": 250,
                "bytes": 1000
            }
        }
    ],
    "callbacks": [
        {
            "ipAddress": "123.456.789.123"
        }
    ],
	  "policies": {
		"skipscriptflags": [
		  "MINIMALDATA",
		  "DERSIG",
		  "NULLDUMMY",
		  "DISCOURAGE_UPGRADABLE_NOPS",
		  "CLEANSTACK"
		],
		"maxtxsizepolicy": 99999,
		"datacarriersize": 100000,
		"maxscriptsizepolicy": 100000,
		"maxscriptnumlengthpolicy": 100000,
		"maxstackmemoryusagepolicy": 10000000,
		"limitancestorcount": 1000,
		"limitcpfpgroupmemberscount": 10,
		"acceptnonstdoutputs": true,
		"datacarrier": true,
		"maxstdtxvalidationduration": 99,
		"maxnonstdtxvalidationduration": 100,
		"minconsolidationfactor": 10,
		"maxconsolidationinputscriptsize": 100,
		"minconfconsolidationinput": 10,
		"acceptnonstdconsolidationinput": false
	  }
}
*/

// Policy is the unmarshalled version of the payload envelope
type Policy struct {
	AcceptNonStdOutputs             bool         `json:"acceptnonstdoutputs"`
	DataCarrier                     bool         `json:"datacarrier"`
	DataCarrierSize                 uint32       `json:"datacarriersize"`
	LimitAncestorCount              uint32       `json:"limitancestorcount"`
	LimitCpfpGroupMembersCount      uint32       `json:"limitcpfpgroupmemberscount"`
	MaxNonStdTxValidationDuration   uint32       `json:"maxnonstdtxvalidationduration"`
	MaxScriptNumLengthPolicy        uint32       `json:"maxscriptnumlengthpolicy"`
	MaxScriptSizePolicy             uint32       `json:"maxscriptsizepolicy"`
	MaxStackMemoryUsagePolicy       uint64       `json:"maxstackmemoryusagepolicy"`
	MaxStdTxValidationDuration      uint32       `json:"maxstdtxvalidationduration"`
	MaxTxSizePolicy                 uint32       `json:"maxtxsizepolicy"`
	SkipScriptFlags                 []ScriptFlag `json:"skipscriptflags"`
	MaxConsolidationFactor          uint32       `json:"minconsolidationfactor"`
	MaxConsolidationInputScriptSize uint32       `json:"maxconsolidationinputscriptsize"`
	MinConfConsolidationInput       uint32       `json:"minconfconsolidationinput"`
	AcceptNonStdConsolidationInput  bool         `json:"acceptnonstdconsolidationinput"`
}

// PolicyQuoteModel is the unmarshalled version of the payload envelope
type PolicyQuoteModel struct {
	APIVersion                string           `json:"apiVersion"`
	Timestamp                 string           `json:"timestamp"`
	ExpiryTime                string           `json:"expiryTime"`
	MinerID                   string           `json:"minerId"`
	CurrentHighestBlockHash   string           `json:"currentHighestBlockHash"`
	CurrentHighestBlockHeight uint64           `json:"currentHighestBlockHeight"`
	Fees                      []FeeObj         `json:"fees"`
	Callbacks                 []PolicyCallback `json:"callbacks"`
	Policies                  Policy           `json:"policies"`
}

// ScriptFlag is a flag used in the policy quote
type ScriptFlag string

// All known script flags
const (
	// FlagCleanStack is the CLEANSTACK flag
	FlagCleanStack ScriptFlag = "CLEANSTACK"
	// FlagDerSig is the DERSIG flag
	FlagDerSig ScriptFlag = "DERSIG"
	// FlagDiscourageUpgradableNops is the DISCOURAGE_UPGRADABLE_NOPS flag
	FlagDiscourageUpgradableNops ScriptFlag = "DISCOURAGE_UPGRADABLE_NOPS"
	// FlagMinimalData is the MINIMALDATA flag
	FlagMinimalData ScriptFlag = "MINIMALDATA"
	// FlagNullDummy is the NULLDUMMY flag
	FlagNullDummy ScriptFlag = "NULLDUMMY"
)

// PolicyCallback is the callback address
type PolicyCallback struct {
	IPAddress string `json:"ipAddress"`
}
