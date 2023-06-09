package minercraft

import (
	"encoding/json"

	"github.com/libsv/go-bk/envelope"
)

// Miner is a configuration per miner, including connection url, auth token, etc
type Miner struct {
	MinerID string `json:"miner_id,omitempty"`
	Name    string `json:"name,omitempty"`
}

// MinerAPIs is a configuration per miner, including connection url, auth token, etc
type MinerAPIs struct {
	MinerID string `json:"miner_id,omitempty"`
	APIs    []API  `json:"apis,omitempty"`
}

// APIType is the type of available APIs
type APIType string

// APIActionName is the name of the action for the API
type APIActionName string

// API is a configuration per miner, including connection url, auth token, etc
type API struct {
	Type  APIType `json:"type,omitempty"`
	Token string  `json:"token,omitempty"`
	URL   string  `json:"url,omitempty"`
}

// APIRoute contains the routes for a specific API related to a specific action
type APIRoute struct {
	Name   APIActionName      `json:"name,omitempty"`
	Routes []APISpecificRoute `json:"routes,omitempty"`
}

// APISpecificRoute contains route definition for a specific API type
type APISpecificRoute struct {
	Route   string  `json:"route,omitempty"`
	APIType APIType `json:"apitype,omitempty"`
}

// JSONEnvelope is a standard response from the Merchant API requests
//
// This type wraps the go-bk JSONEnvelope which performs validation of the
// signatures (if we have any) and will return true / false if valid.
//
// We wrap this, so we can append some additional miner info and a validated
// helper property to indicate if the envelope is or isn't valid.
// Consumers can also independently validate the envelope.
type JSONEnvelope struct {
	Miner     *Miner  `json:"miner"`     // Custom field for our internal Miner configuration
	Validated bool    `json:"validated"` // Custom field if the signature has been validated
	ApiType   APIType `json:"apiType"`   // Custom field for the API type
	envelope.JSONEnvelope
}

// process will take the raw payload bytes, unmarshall into a JSONEnvelope
// and validate the signature vs payload.
func (p *JSONEnvelope) process(miner *Miner, bodyContents []byte) (err error) {
	// Set the miner on the response
	p.Miner = miner

	// Unmarshal the response
	if err = json.Unmarshal(bodyContents, &p); err != nil {
		return
	}

	return
}
