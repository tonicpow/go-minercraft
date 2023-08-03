package minercraft

import "time"

const (

	// version is the current package version
	version = "v2.0.3"

	// defaultUserAgent is the default user agent for all requests
	defaultUserAgent string = "go-minercraft: " + version

	// defaultFastQuoteTimeout is used for the FastestQuote timeout
	defaultFastQuoteTimeout = 20 * time.Second
)

const (
	// MAPI stands for Merchant API
	MAPI APIType = "mAPI"
	// Arc stands for Arc API
	Arc APIType = "Arc"
)

const (
	// PolicyQuote is the name of the PolicyQuote API action
	PolicyQuote APIActionName = "PolicyQuote"
	// FeeQuote is the name of the FeeQuote API action
	FeeQuote APIActionName = "FeeQuote"
	// QueryTx is the name of the Query Transaction API action
	QueryTx APIActionName = "QueryTx"
	// SubmitTx is the name of the Submit Transaction API action
	SubmitTx APIActionName = "SubmitTx"
	// SubmitTxs is the name of the Submit multiple Transactions API action
	SubmitTxs APIActionName = "SubmitTxs"
)

// mAPI routes
const (
	// mAPIRoutePolicyQuote is the route for getting a policy quote
	mAPIRoutePolicyQuote = "/mapi/policyQuote"

	// mAPIRouteFeeQuote is the route for getting a fee quote
	mAPIRouteFeeQuote = "/mapi/feeQuote"

	// mAPIRouteQueryTx is the route for querying a transaction
	mAPIRouteQueryTx = "/mapi/tx/"

	// mAPIRouteSubmitTx is the route for submit a transaction
	mAPIRouteSubmitTx = "/mapi/tx"

	// mAPIRouteSubmitTxs is the route for submit batched transactions
	mAPIRouteSubmitTxs = "/mapi/txs"
)

// Arc routes
const (
	// arcRoutePolicyQuote is the route for getting a policy quote
	arcRoutePolicyQuote = "/v1/policy"
	// arcRouteQueryTx is the route for querying a transaction
	arcRouteQueryTx = "/v1/tx/"
	// arcRouteSubmitTx is the route for submit a transaction
	arcRouteSubmitTx = "/v1/tx"
	// arcRouteSubmitTxs is the route for submit batched transactions
	arcRouteSubmitTxs = "/v1/txs"
)

// Routes is a list of known actions with it's routes for the different APIs
var Routes = []APIRoute{
	{
		Name: PolicyQuote,
		Routes: []APISpecificRoute{
			{Route: mAPIRoutePolicyQuote, APIType: MAPI},
			{Route: arcRoutePolicyQuote, APIType: Arc},
		},
	},
	{
		Name: FeeQuote,
		Routes: []APISpecificRoute{
			{Route: mAPIRouteFeeQuote, APIType: MAPI},
		},
	},
	{
		Name: QueryTx,
		Routes: []APISpecificRoute{
			{Route: mAPIRouteQueryTx, APIType: MAPI},
			{Route: arcRouteQueryTx, APIType: Arc},
		},
	},
	{
		Name: SubmitTx,
		Routes: []APISpecificRoute{
			{Route: mAPIRouteSubmitTx, APIType: MAPI},
			{Route: arcRouteSubmitTx, APIType: Arc},
		},
	},
	{
		Name: SubmitTxs,
		Routes: []APISpecificRoute{
			{Route: mAPIRouteSubmitTxs, APIType: MAPI},
			{Route: arcRouteSubmitTxs, APIType: Arc},
		},
	},
}

const (
	// MinerTaal is the name of the known miner for "Taal"
	MinerTaal = "Taal"

	// MinerGorillaPool is the name of the known miner for "GorillaPool"
	MinerGorillaPool = "GorillaPool"

	// MinerMempool is the name of the known miner for "Mempool"
	MinerMempool = "Mempool"

	// MinerMatterpool is the name of the known miner for "Matterpool"
	MinerMatterpool = "Matterpool"
)

// KnownMiners is a pre-filled list of known miners
const KnownMiners = `
[
   {
      "name":"Taal",
      "miner_id":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270"
   },
   {
      "name":"GorillaPool",
      "miner_id":"03ad780153c47df915b3d2e23af727c68facaca4facd5f155bf5018b979b9aeb83"
   }
]
`

// KnownMinersAll is a pre-filled list of known miners
// deprecated: use KnownMiners instead
const KnownMinersAll = `
[
   {
      "name":"Taal",
      "miner_id":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270"
   },
   {
      "name":"Mempool",
      "miner_id":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270"
   },
   {
      "name":"Matterpool",
      "miner_id":"0253a9b2d017254b91704ba52aad0df5ca32b4fb5cb6b267ada6aefa2bc5833a93"
   },
   {
      "name":"GorillaPool",
      "miner_id":"03ad780153c47df915b3d2e23af727c68facaca4facd5f155bf5018b979b9aeb83"
   }
]
`

// KnownMinersAPIs is a pre-filled list of known miners with their APIs
// Any pre-filled tokens are for free use only
// update your custom token with client.MinerUpdateToken("name", "token")
const KnownMinersAPIs = `
[
	{
		"miner_id":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270",
		"apis":[
		   {
			  "token":"",
			  "url":"https://merchantapi.taal.com",
			  "type":"mAPI"
		   },
		   {
			  "token":"",
			  "url":"https://tapi.taal.com/arc",
			  "type":"Arc"
		   }
		]
	},
	{
		"miner_id":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270",
		"apis":[
		   {
			  "token":"561b756d12572020ea9a104c3441b71790acbbce95a6ddbf7e0630971af9424b",
			  "url":"https://www.ddpurse.com/openapi",
			  "type":"mAPI"
		   }
		]
	},
	{
		"miner_id":"0253a9b2d017254b91704ba52aad0df5ca32b4fb5cb6b267ada6aefa2bc5833a93",
		"apis":[
		   {
			  "token":"",
			  "url":"https://merchantapi.matterpool.io",
			  "type":"mAPI"
		   }
		]
	},
	{
		"miner_id":"03ad780153c47df915b3d2e23af727c68facaca4facd5f155bf5018b979b9aeb83",
      	"apis":[
         	{
            	"token":"",
            	"url":"https://merchantapi.gorillapool.io",
            	"type":"mAPI"
         	},
         	{
            	"token":"",
            	"url":"https://arc.gorillapool.io",
            	"type":"Arc"
         	}
      	]
	}
]
`

// KnownMinersAPIsAll is a pre-filled list of known miners with their APIs
// Any pre-filled tokens are for free use only
// update your custom token with client.MinerUpdateToken("name", "token")
// deprecated: use KnownMinersAPIsAll instead
const KnownMinersAPIsAll = `
[
	{
		"miner_id":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270",
		"apis":[
		   {
			  "token":"",
			  "url":"https://merchantapi.taal.com",
			  "type":"mAPI"
		   },
		   {
			  "token":"",
			  "url":"https://tapi.taal.com/arc",
			  "type":"Arc"
		   }
		]
	},
	{
		"miner_id":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270",
		"apis":[
		   {
			  "token":"561b756d12572020ea9a104c3441b71790acbbce95a6ddbf7e0630971af9424b",
			  "url":"https://www.ddpurse.com/openapi",
			  "type":"mAPI"
		   }
		]
	},
	{
		"miner_id":"0253a9b2d017254b91704ba52aad0df5ca32b4fb5cb6b267ada6aefa2bc5833a93",
		"apis":[
		   {
			  "token":"",
			  "url":"https://merchantapi.matterpool.io",
			  "type":"mAPI"
		   }
		]
	},
	{
		"miner_id":"03ad780153c47df915b3d2e23af727c68facaca4facd5f155bf5018b979b9aeb83",
      	"apis":[
         	{
            	"token":"",
            	"url":"https://merchantapi.gorillapool.io",
            	"type":"mAPI"
         	},
         	{
            	"token":"",
            	"url":"https://arc.gorillapool.io",
            	"type":"Arc"
         	}
      	]
	}
]
`
