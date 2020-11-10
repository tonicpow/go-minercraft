package minercraft

const (

	// version is the current package version
	version = "v0.0.7"

	// defaultUserAgent is the default user agent for all requests
	defaultUserAgent string = "go-minercraft: " + version

	// defaultProtocol is used for url endpoints in requests
	defaultProtocol = "https://"
)

const (
	// routeFeeQuote is the route for getting a fee quote
	routeFeeQuote = "/mapi/feeQuote"

	// routeQueryTx is the route for querying a transaction
	routeQueryTx = "/mapi/tx/"

	// routeSubmitTx is the route for submit a transaction
	routeSubmitTx = "/mapi/tx"
)

const (
	// MinerTaal is the name of the known miner for "Taal"
	MinerTaal = "Taal"

	// MinerMempool is the name of the known miner for "Mempool"
	MinerMempool = "Mempool"

	// MinerMatterpool is the name of the known miner for "Matterpool"
	MinerMatterpool = "Matterpool"
)