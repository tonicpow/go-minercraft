package minercraft

import "time"

const (

	// version is the current package version
	version = "v0.9.1"

	// defaultUserAgent is the default user agent for all requests
	defaultUserAgent string = "go-minercraft: " + version

	// defaultFastQuoteTimeout is used for the FastestQuote timeout
	defaultFastQuoteTimeout = 20 * time.Second
)

const (
	// mAPI stands for Merchant API
	mAPI APIType = "mAPI"
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

var Routes = []APIRoute{
	{
		Name: PolicyQuote,
		Routes: []APISpecificRoute{
			{Route: mAPIRoutePolicyQuote, APIType: mAPI},
		},
	},
	{
		Name: FeeQuote,
		Routes: []APISpecificRoute{
			{Route: mAPIRouteFeeQuote, APIType: mAPI},
		},
	},
	{
		Name: QueryTx,
		Routes: []APISpecificRoute{
			{Route: mAPIRouteQueryTx, APIType: mAPI},
		},
	},
	{
		Name: SubmitTx,
		Routes: []APISpecificRoute{
			{Route: mAPIRouteSubmitTx, APIType: mAPI},
		},
	},
	{
		Name: SubmitTxs,
		Routes: []APISpecificRoute{
			{Route: mAPIRouteSubmitTxs, APIType: mAPI},
		},
	},
}

const (
	// MinerTaal is the name of the known miner for "Taal"
	MinerTaal = "Taal"

	// MinerMempool is the name of the known miner for "Mempool"
	MinerMempool = "Mempool"

	// MinerMatterpool is the name of the known miner for "Matterpool"
	MinerMatterpool = "Matterpool"

	// MinerGorillaPool is the name of the known miner for "GorillaPool"
	MinerGorillaPool = "GorillaPool"
)

// KnownMiners is a pre-filled list of known miners
// Any pre-filled tokens are for free use only
// update your custom token with client.MinerUpdateToken("name", "token")
const KnownMiners = `
[
   {
      "name":"Taal",
      "miner_id":"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270",
      "apis":[
         {
            "token":"",
            "url":"https://merchantapi.taal.com",
            "type":"mAPI"
         },
         {
            "token":"",
            "url":"https://tapi.taal.com/arc/v1",
            "type":"Arc"
         }
      ]
   },
   {
      "name":"Mempool",
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
      "name":"Matterpool",
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
      "name":"GorillaPool",
      "miner_id":"03ad780153c47df915b3d2e23af727c68facaca4facd5f155bf5018b979b9aeb83",
      "apis":[
         {
            "token":"",
            "url":"https://merchantapi.gorillapool.io",
            "type":"mAPI"
         },
         {
            "token":"",
            "url":"https://arc.gorillapool.io/v1",
            "type":"Arc"
         }
      ]
   }
]
`
