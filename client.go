package minercraft

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gojektech/heimdall/v6"
	"github.com/gojektech/heimdall/v6/httpclient"
)

// HTTPInterface is used for the http client (mocking heimdall)
type HTTPInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the parent struct that contains the miner clients and list of miners to use
type Client struct {
	apiType    APIType        // The API type to use
	httpClient HTTPInterface  // Interface for all HTTP requests
	miners     []*Miner       // List of loaded miners
	minerAPIs  []*MinerAPIs   // List of loaded miners APIs
	Options    *ClientOptions // Client options config
}

// AddMiner will add a new miner to the list of miners
func (c *Client) AddMiner(miner Miner, apis []API) error {
	// Check if miner name is empty
	if len(miner.Name) == 0 {
		return errors.New("missing miner name")
	}

	// Check if apis is empty or nil
	if len(apis) == 0 || apis == nil {
		return errors.New("at least one API must be provided")
	}

	// Check if a miner with that name already exists
	existingMiner := c.MinerByName(miner.Name)
	if existingMiner != nil {
		return fmt.Errorf("miner %s already exists", miner.Name)
	}

	// Check if a miner with the minerID already exists
	if len(miner.MinerID) > 0 {
		if existingMiner = c.MinerByID(miner.MinerID); existingMiner != nil {
			return fmt.Errorf("miner %s already exists", miner.MinerID)
		}
	}

	// Check if the MinerAPIs already exist for the given MinerID
	existingMinerAPIs := c.MinerAPIsByMinerID(miner.MinerID)
	if existingMinerAPIs != nil {
		return fmt.Errorf("miner APIs for MinerID %s already exist", miner.MinerID)
	}

	// Check if the API types are valid
	for _, api := range apis {
		if !isValidAPIType(api.Type) {
			return fmt.Errorf("invalid API type: %s", api.Type)
		}
	}

	// Check if the API types are unique within the provided APIs
	apiTypes := make(map[APIType]bool)
	for _, api := range apis {
		if apiTypes[api.Type] {
			return fmt.Errorf("duplicate API type found: %s", api.Type)
		}
		apiTypes[api.Type] = true
	}

	// Check if the MinerID is unique or generate a new one
	if len(miner.MinerID) == 0 || !c.isUniqueMinerID(miner.MinerID) {
		miner.MinerID = generateUniqueMinerID()
	}

	// Append the new miner
	c.miners = append(c.miners, &miner)

	// Append the new miner APIs
	c.minerAPIs = append(c.minerAPIs, &MinerAPIs{
		MinerID: miner.MinerID,
		APIs:    apis,
	})

	return nil
}

// RemoveMiner will remove a miner from the list
func (c *Client) RemoveMiner(miner *Miner) bool {
	for i, m := range c.miners {
		if m.Name == miner.Name || m.MinerID == miner.MinerID {
			c.miners[i] = c.miners[len(c.miners)-1]
			c.miners = c.miners[:len(c.miners)-1]
			return true
		}
	}
	// Miner not found
	return false
}

// MinerByName will return a miner given a name
func (c *Client) MinerByName(name string) *Miner {
	return MinerByName(c.miners, name)
}

// MinerByName will return a miner from a given set of miners
func MinerByName(miners []*Miner, minerName string) *Miner {
	for index, miner := range miners {
		if strings.EqualFold(minerName, miner.Name) {
			return miners[index]
		}
	}
	return nil
}

// MinerByID will return a miner given a miner id
func (c *Client) MinerByID(minerID string) *Miner {
	return MinerByID(c.miners, minerID)
}

// MinerByID will return a miner from a given set of miners
func MinerByID(miners []*Miner, minerID string) *Miner {
	for index, miner := range miners {
		if strings.EqualFold(minerID, miner.MinerID) {
			return miners[index]
		}
	}
	return nil
}

// MinerAPIByMinerID will return a miner's API given a miner id and API type
func (c *Client) MinerAPIByMinerID(minerID string, apiType APIType) (*API, error) {
	for _, minerAPI := range c.minerAPIs {
		if minerAPI.MinerID == minerID {
			for i := range minerAPI.APIs {
				if minerAPI.APIs[i].Type == apiType {
					return &minerAPI.APIs[i], nil
				}
			}
		}
	}
	return nil, &APINotFoundError{MinerID: minerID, APIType: apiType}
}

// MinerAPIsByMinerID will return a miner's APIs given a miner id
func (c *Client) MinerAPIsByMinerID(minerID string) *MinerAPIs {
	for _, minerAPIs := range c.minerAPIs {
		if minerAPIs.MinerID == minerID {
			return minerAPIs
		}
	}
	return nil
}

// ActionRouteByAPIType will return the route for a given action and API type
func ActionRouteByAPIType(actionName APIActionName, apiType APIType) (string, error) {
	for _, apiRoute := range Routes {
		if apiRoute.Name == actionName {
			for _, route := range apiRoute.Routes {
				if route.APIType == apiType {
					return route.Route, nil
				}
			}
		}
	}
	return "", &ActionRouteNotFoundError{ActionName: actionName, APIType: apiType}
}

// MinerUpdateToken will find a miner by name and update the token
func (c *Client) MinerUpdateToken(name, token string, apiType APIType) {
	if miner := c.MinerByName(name); miner != nil {
		api, _ := c.MinerAPIByMinerID(miner.MinerID, apiType)
		api.Token = token
	}
}

// Miners will return the list of miners
func (c *Client) Miners() []*Miner {
	return c.miners
}

// UserAgent will return the user agent
func (c *Client) UserAgent() string {
	return c.Options.UserAgent
}

// ClientOptions holds all the configuration for connection, dialer and transport
type ClientOptions struct {
	BackOffExponentFactor          float64       `json:"back_off_exponent_factor"`
	BackOffInitialTimeout          time.Duration `json:"back_off_initial_timeout"`
	BackOffMaximumJitterInterval   time.Duration `json:"back_off_maximum_jitter_interval"`
	BackOffMaxTimeout              time.Duration `json:"back_off_max_timeout"`
	DialerKeepAlive                time.Duration `json:"dialer_keep_alive"`
	DialerTimeout                  time.Duration `json:"dialer_timeout"`
	RequestRetryCount              int           `json:"request_retry_count"`
	RequestTimeout                 time.Duration `json:"request_timeout"`
	TransportExpectContinueTimeout time.Duration `json:"transport_expect_continue_timeout"`
	TransportIdleTimeout           time.Duration `json:"transport_idle_timeout"`
	TransportMaxIdleConnections    int           `json:"transport_max_idle_connections"`
	TransportTLSHandshakeTimeout   time.Duration `json:"transport_tls_handshake_timeout"`
	UserAgent                      string        `json:"user_agent"`
}

// DefaultClientOptions will return a ClientOptions struct with the default settings.
// Useful for starting with the default and then modifying as needed
func DefaultClientOptions() (clientOptions *ClientOptions) {
	return &ClientOptions{
		BackOffExponentFactor:          2.0,
		BackOffInitialTimeout:          2 * time.Millisecond,
		BackOffMaximumJitterInterval:   2 * time.Millisecond,
		BackOffMaxTimeout:              10 * time.Millisecond,
		DialerKeepAlive:                20 * time.Second,
		DialerTimeout:                  5 * time.Second,
		RequestRetryCount:              2,
		RequestTimeout:                 30 * time.Second,
		TransportExpectContinueTimeout: 3 * time.Second,
		TransportIdleTimeout:           20 * time.Second,
		TransportMaxIdleConnections:    10,
		TransportTLSHandshakeTimeout:   5 * time.Second,
		UserAgent:                      defaultUserAgent,
	}
}

// NewClient creates a new client for requests
//
// clientOptions: inject custom client options on load
// customHTTPClient: use your own custom HTTP client
// customMiners: use your own custom list of miners
func NewClient(clientOptions *ClientOptions, customHTTPClient HTTPInterface,
	apiType APIType, customMiners []*Miner, customMinersAPIDef []*MinerAPIs) (client ClientInterface, err error) {

	// Create the new client
	return createClient(clientOptions, apiType, customHTTPClient, customMiners, customMinersAPIDef)
}

// createClient will make a new http client based on the options provided
func createClient(options *ClientOptions, apiType APIType, customHTTPClient HTTPInterface,
	customMiners []*Miner, customMinersAPIDef []*MinerAPIs) (c *Client, err error) {

	// Create a client
	c = new(Client)

	// For now set MerchantAPI as the default if not set
	if apiType == "" {
		apiType = MAPI
	}

	// Set the client API type
	c.apiType = apiType

	// Set options (either default or user modified)
	if options == nil {
		options = DefaultClientOptions()
	}

	// Set the options
	c.Options = options

	// Load custom vs pre-defined
	if len(customMiners) > 0 && len(customMinersAPIDef) > 0 {
		c.miners = customMiners
		c.minerAPIs = customMinersAPIDef
	} else {
		c.miners, err = DefaultMiners()
		if err != nil {
			return nil, err
		}

		c.minerAPIs, err = DefaultMinersAPIs()
		if err != nil {
			return nil, err
		}
	}

	// Is there a custom HTTP client to use?
	if customHTTPClient != nil {
		c.httpClient = customHTTPClient
		return
	}

	// dial is the net dialer for clientDefaultTransport
	dial := &net.Dialer{KeepAlive: options.DialerKeepAlive, Timeout: options.DialerTimeout}

	// clientDefaultTransport is the default transport struct for the HTTP client
	clientDefaultTransport := &http.Transport{
		DialContext:           dial.DialContext,
		ExpectContinueTimeout: options.TransportExpectContinueTimeout,
		IdleConnTimeout:       options.TransportIdleTimeout,
		MaxIdleConns:          options.TransportMaxIdleConnections,
		Proxy:                 http.ProxyFromEnvironment,
		TLSHandshakeTimeout:   options.TransportTLSHandshakeTimeout,
	}

	// Determine the strategy for the http client
	if options.RequestRetryCount <= 0 {

		// no retry enabled
		c.httpClient = httpclient.NewClient(
			httpclient.WithHTTPTimeout(options.RequestTimeout),
			httpclient.WithHTTPClient(&http.Client{
				Transport: clientDefaultTransport,
				Timeout:   options.RequestTimeout,
			}),
		)
		return
	}

	// Retry enabled - create exponential back-off
	c.httpClient = httpclient.NewClient(
		httpclient.WithHTTPTimeout(options.RequestTimeout),
		httpclient.WithRetrier(heimdall.NewRetrier(
			heimdall.NewExponentialBackoff(
				options.BackOffInitialTimeout,
				options.BackOffMaxTimeout,
				options.BackOffExponentFactor,
				options.BackOffMaximumJitterInterval,
			))),
		httpclient.WithRetryCount(options.RequestRetryCount),
		httpclient.WithHTTPClient(&http.Client{
			Transport: clientDefaultTransport,
			Timeout:   options.RequestTimeout,
		}),
	)

	return
}

// DefaultMiners will parse the config JSON and return a list of miners
func DefaultMiners() (miners []*Miner, err error) {
	err = json.Unmarshal([]byte(KnownMiners), &miners)
	return
}

// DefaultMinersAPIs will parse the config JSON and return a list of miner APIs
func DefaultMinersAPIs() (minerAPIs []*MinerAPIs, err error) {
	err = json.Unmarshal([]byte(KnownMinersAPIs), &minerAPIs)
	return
}

// APIType will return the API type
func (c *Client) APIType() APIType {
	return c.apiType
}

// isValidAPIType will return true if the API type is valid and part of our predefined list
func isValidAPIType(apiType APIType) bool {
	switch apiType {
	case MAPI, Arc:
		return true
	default:
		return false
	}
}

// isUniqueMinerID will return true if the miner ID is unique
func (c *Client) isUniqueMinerID(minerID string) bool {
	for _, miner := range c.miners {
		if miner.MinerID == minerID {
			return false
		}
	}
	return true
}

// generateUniqueMinerID will generate a unique miner ID
func generateUniqueMinerID() string {
	const idLength = 8
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	id := make([]byte, idLength)
	for i := range id {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ""
		}
		id[i] = letters[num.Int64()]
	}

	return string(id)
}
