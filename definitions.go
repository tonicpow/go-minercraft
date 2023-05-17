package minercraft

import (
	"encoding/json"

	"github.com/libsv/go-bk/envelope"
)

// Miner is a configuration per miner, including connection url, auth token, etc
type Miner struct {
	MinerID string `json:"miner_id,omitempty"`
	Name    string `json:"name,omitempty"`
	APIs    []API  `json:"apis,omitempty"`
}

// API is a configuration per miner, including connection url, auth token, etc
type API struct {
	Type  APIType `json:"type,omitempty"`
	Token string  `json:"token,omitempty"`
	URL   string  `json:"url,omitempty"`
}

// APIType is the type of API
type APIType string

const (
	// mAPI stands for Merchant API
	mAPI APIType = "mAPI"
	// Arc stands for Arc API
	Arc APIType = "Arc"
)

// JSONEnvelope is a standard response from the Merchant API requests
//
// This type wraps the go-bk JSONEnvelope which performs validation of the
// signatures (if we have any) and will return true / false if valid.
//
// We wrap this, so we can append some additional miner info and a validated
// helper property to indicate if the envelope is or isn't valid.
// Consumers can also independently validate the envelope.
type JSONEnvelope struct {
	Miner     *Miner `json:"miner"`     // Custom field for our internal Miner configuration
	Validated bool   `json:"validated"` // Custom field if the signature has been validated
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

	// verify JSONEnvelope
	p.Validated, err = p.IsValid()
	return
}
