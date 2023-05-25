package mapi

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

// PolicyPayload is the unmarshalled version of the payload envelope
type PolicyPayload struct {
	FeePayload                   // Inherit the same structure as the fee payload
	Callbacks  []*PolicyCallback `json:"callbacks"` // IP addresses of double-spend notification servers such as mAPI reference implementation
	Policies   *Policy           `json:"policies"`  // values of miner policies as configured by the mAPI reference implementation administrator
}

// ScriptFlag is a flag used in the policy quote
type ScriptFlag string

// All known script flags
const (
	FlagCleanStack               ScriptFlag = "CLEANSTACK"
	FlagDerSig                   ScriptFlag = "DERSIG"
	FlagDiscourageUpgradableNops ScriptFlag = "DISCOURAGE_UPGRADABLE_NOPS"
	FlagMinimalData              ScriptFlag = "MINIMALDATA"
	FlagNullDummy                ScriptFlag = "NULLDUMMY"
)

// Policy is the struct of a policy (from policy quote response)
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

// PolicyCallback is the callback address
type PolicyCallback struct {
	IPAddress string `json:"ipAddress"`
}
