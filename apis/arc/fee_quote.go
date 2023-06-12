// Package arc provides the API structures for the ARC service
package arc

import "github.com/libsv/go-bt/v2"

// FeePayload is the unmarshalled version of the payload envelope
type FeePayload struct {
	Fees *bt.Fee `json:"fees"`
}
