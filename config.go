package minercraft

const (

	// version is the current package version
	version = "v0.0.1"

	// defaultUserAgent is the default user agent for all requests
	defaultUserAgent string = "go-minercraft: " + version
)

// KnownMiners is a pre-filled list of known miners
const KnownMiners = `
[
  {
   "name": "Taal",
   "miner_id": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270",
   "token": "",
   "url": "merchantapi.taal.com"
  },
  {
   "name": "Mempool",
   "miner_id": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270",
   "token": "561b756d12572020ea9a104c3441b71790acbbce95a6ddbf7e0630971af9424b",
   "url": "www.ddpurse.com/openapi"
  },
  {
   "name": "Matterpool",
   "miner_id": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087",
   "token": "",
   "url": "merchantapi.matterpool.io"
  }
]
`
