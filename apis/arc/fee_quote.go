package arc

import "github.com/libsv/go-bt/v2"

// FeePayload is the unmarshalled version of the payload envelope
type FeePayload struct {
	Fees *bt.Fee `json:"fees"`
}
