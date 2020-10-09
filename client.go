package minercraft

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gojektech/heimdall/v6"
	"github.com/gojektech/heimdall/v6/httpclient"
)

// httpInterface is used for the http client (mocking heimdall)
type httpInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the parent struct that contains the miner clients and list of miners to use
type Client struct {
	httpClient httpInterface  // Interface for all HTTP requests
	Miners     []*Miner       // List of loaded miners
	Options    *ClientOptions // Client options config
}

// AddMiner will add a new miner to the list of miners
func (c *Client) AddMiner(miner Miner) error {

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

	// Remove any protocol(s)
	miner.URL = strings.Replace(miner.URL, "https://", "", -1)
	miner.URL = strings.Replace(miner.URL, "http://", "", -1)

	// Append the new miner
	c.Miners = append(c.Miners, &miner)
	return nil
}

// MinerByName will return a miner given a name
func (c *Client) MinerByName(name string) *Miner {
	for index, miner := range c.Miners {
		if strings.EqualFold(name, miner.Name) {
			return c.Miners[index]
		}
	}
	return nil
}

// MinerByID will return a miner given a miner id
func (c *Client) MinerByID(minerID string) *Miner {
	for index, miner := range c.Miners {
		if strings.EqualFold(minerID, miner.MinerID) {
			return c.Miners[index]
		}
	}
	return nil
}

// MinerUpdateToken will find a miner by name and update the token
func (c *Client) MinerUpdateToken(name, token string) {
	if miner := c.MinerByName(name); miner != nil {
		miner.UpdateToken(token)
	}
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

// DefaultClientOptions will return an Options struct with the default settings.
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
		RequestTimeout:                 10 * time.Second,
		TransportExpectContinueTimeout: 3 * time.Second,
		TransportIdleTimeout:           20 * time.Second,
		TransportMaxIdleConnections:    10,
		TransportTLSHandshakeTimeout:   5 * time.Second,
		UserAgent:                      defaultUserAgent,
	}
}

// NewClient creates a new client for requests
func NewClient(clientOptions *ClientOptions, customHTTPClient *http.Client) (client *Client, err error) {

	// Create the new client
	client = createClient(clientOptions, customHTTPClient)

	// Load all known miners
	err = json.Unmarshal([]byte(KnownMiners), &client.Miners)

	return
}

// createClient will make a new http client based on the options provided
func createClient(options *ClientOptions, customHTTPClient *http.Client) (c *Client) {

	// Create a client
	c = new(Client)

	// Is there a custom HTTP client to use?
	if customHTTPClient != nil {
		c.httpClient = customHTTPClient
		return
	}

	// Set options (either default or user modified)
	if options == nil {
		options = DefaultClientOptions()
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

	// Determine the strategy for the http client (no retry enabled)
	if options.RequestRetryCount <= 0 {
		c.httpClient = httpclient.NewClient(
			httpclient.WithHTTPTimeout(options.RequestTimeout),
			httpclient.WithHTTPClient(&http.Client{
				Transport: clientDefaultTransport,
				Timeout:   options.RequestTimeout,
			}),
		)
	} else { // Retry enabled
		// Create exponential back-off
		backOff := heimdall.NewExponentialBackoff(
			options.BackOffInitialTimeout,
			options.BackOffMaxTimeout,
			options.BackOffExponentFactor,
			options.BackOffMaximumJitterInterval,
		)

		c.httpClient = httpclient.NewClient(
			httpclient.WithHTTPTimeout(options.RequestTimeout),
			httpclient.WithRetrier(heimdall.NewRetrier(backOff)),
			httpclient.WithRetryCount(options.RequestRetryCount),
			httpclient.WithHTTPClient(&http.Client{
				Transport: clientDefaultTransport,
				Timeout:   options.RequestTimeout,
			}),
		)
	}

	c.Options = options

	return
}

// RequestResponse is the response from a request
type RequestResponse struct {
	BodyContents []byte `json:"body_contents"` // Raw body response
	Error        error  `json:"error"`         // If an error occurs
	Method       string `json:"method"`        // Method is the HTTP method used
	PostData     string `json:"post_data"`     // PostData is the post data submitted if POST/PUT request
	StatusCode   int    `json:"status_code"`   // StatusCode is the last code from the request
	URL          string `json:"url"`           // URL is used for the request
}

// httpRequest is a generic request wrapper that can be used without constraints
func httpRequest(client *Client, method, url, token string, payload []byte) (response *RequestResponse) {

	// Set reader
	var bodyReader io.Reader

	// Start the response
	response = new(RequestResponse)

	// Add post data if applicable
	if method == http.MethodPost || method == http.MethodPut {
		bodyReader = bytes.NewBuffer(payload)
		response.PostData = string(payload)
	}

	// Store for debugging purposes
	response.Method = method
	response.URL = url

	// Start the request
	var request *http.Request
	if request, response.Error = http.NewRequestWithContext(context.Background(), method, url, bodyReader); response.Error != nil {
		return
	}

	// Change the header (user agent is in case they block default Go user agents)
	request.Header.Set("User-Agent", client.Options.UserAgent)

	// Set the content type on Method
	if method == http.MethodPost || method == http.MethodPut {
		request.Header.Set("Content-Type", "application/json")
	}

	// Set a token if supplied
	if len(token) > 0 {
		request.Header.Set("token", token)
	}

	// Fire the http request
	var resp *http.Response
	if resp, response.Error = client.httpClient.Do(request); response.Error != nil {
		if resp != nil {
			response.StatusCode = resp.StatusCode
		}
		return
	}

	// Close the response body
	defer func() {
		_ = resp.Body.Close()
	}()

	// Set the status
	response.StatusCode = resp.StatusCode

	// Check status code
	if http.StatusOK != resp.StatusCode {
		response.Error = fmt.Errorf("status code: %d does not match %d", resp.StatusCode, http.StatusOK)
		return
	}

	// Read the body
	response.BodyContents, response.Error = ioutil.ReadAll(resp.Body)

	return
}
