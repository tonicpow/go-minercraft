package minercraft

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonicpow/go-minercraft/v2/apis/mapi"
)

// MockClient mocks the http client.
type MockClient struct {
	MockDo func(req *http.Request) (*http.Response, error)
}

// Do implement the mock version of Do method for http client interface
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

const rawTx = "01000000010136836d73f29cbe648bc2aeea20286502a3c2f2d3cff54522d0cc76bb755e9f000000006a4730440220533430d6f29d9437f94c60f0a59c857d254108fb9c375415fa53e9248d8bec5d0220606810830a6175dbee71da54fff6378a7aadfdb4b4714dc500cb3466ad4500004121027ae06a5b3fe1de495fa9d4e738e48810b8b06fa6c959a5305426f78f42b48f8cffffffff018c949800000000001976a91482932cf55b847ffa52832d2bbec2838f658f226788ac00000000"
const submitResponse = `{
	"payload": "{\"apiVersion\":\"1.3.0\",\"timestamp\":\"2020-11-13T08:31:56.5722511Z\",\"minerId\":\"030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e\",\"currentHighestBlockHash\":\"08dc4bb006fc7e7186544343c3ccbb5a773d0a19cd2ccff1fa52f51eb6faf2ab\",\"currentHighestBlockHeight\":151,\"txSecondMempoolExpiry\":0,\"txs\":[{\"txid\":\"3145011f34a00d0666ea265b87c8e44108f87d3b53b853976906519ee8e1475f\",\"returnResult\":\"failure\",\"resultDescription\":\"Missing inputs\",\"conflictedWith\":[{\"txid\":\"86e1b384d3d169fd6aa4d34cf2d6f487436da54154befaab5a1fb25f844d65a8\",\"size\":191,\"hex\":\"01000000010136836d73f29cbe648bc2aeea20286502a3c2f2d3cff54522d0cc76bb755e9f000000006a4730440220761fb63128d4184fc142f2e854c499c52422db0136191f29f0bbe0969b6021770220536d72606d49dbbd244d2633b8b19031234138f045c530cc773e6e72bb34c62c4121027ae06a5b3fe1de495fa9d4e738e48810b8b06fa6c959a5305426f78f42b48f8cffffffff0198929800000000001976a91482932cf55b847ffa52832d2bbec2838f658f226788ac00000000\"}]},{\"txid\":\"c8a087b1ee775fa29697511ecd64e800941c8a22db6ed0989fb27a1d2d6798da\",\"returnResult\":\"success\",\"resultDescription\":\"\"}],\"failureCount\":1}",
	"signature": "304402200c4b0dc179906581eb32953abeddbef5799d302d82367aa9a469d79c15f932f3022029e827af6122290d5e1b80c50676b0336f4c7658ca67ba4819396dff9c6239a6",
	"publicKey": "030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e",
	"encoding": "UTF-8",
	"mimetype": "application/json"
}`

func TestClient_SubmitTransactions(t *testing.T) {
	tests := map[string]struct {
		txs        []Transaction
		miner      string
		mockHTTPFn func(*http.Request) (*http.Response, error)
		code       int
		body       string
		exp        *SubmitTransactionsResponse
		err        error
	}{
		"invalid miner should return error": {
			err: errors.New("miner was nil"),
		},
		"empty txs should return error": {
			miner: MinerGorillaPool,
			err:   errors.New("no transactions"),
			txs:   []Transaction{},
		},
		"one valid and one invalid tx should return specific response": {
			miner: MinerGorillaPool,
			txs:   []Transaction{{RawTx: rawTx}},
			exp: &SubmitTransactionsResponse{
				Payload: UnifiedTxsPayload{
					APIVersion:                "1.3.0",
					Timestamp:                 time.Date(2020, time.November, 13, 8, 31, 56, 572251100, time.UTC),
					MinerID:                   "030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e",
					CurrentHighestBlockHash:   "08dc4bb006fc7e7186544343c3ccbb5a773d0a19cd2ccff1fa52f51eb6faf2ab",
					CurrentHighestBlockHeight: 151,
					TxSecondMempoolExpiry:     0,
					Txs: []UnifiedTx{{
						TxID:              "3145011f34a00d0666ea265b87c8e44108f87d3b53b853976906519ee8e1475f",
						ReturnResult:      "failure",
						ResultDescription: "Missing inputs",
						ConflictedWith: []mapi.ConflictedWith{{
							TxID: "86e1b384d3d169fd6aa4d34cf2d6f487436da54154befaab5a1fb25f844d65a8",
							Size: 191,
							Hex:  "01000000010136836d73f29cbe648bc2aeea20286502a3c2f2d3cff54522d0cc76bb755e9f000000006a4730440220761fb63128d4184fc142f2e854c499c52422db0136191f29f0bbe0969b6021770220536d72606d49dbbd244d2633b8b19031234138f045c530cc773e6e72bb34c62c4121027ae06a5b3fe1de495fa9d4e738e48810b8b06fa6c959a5305426f78f42b48f8cffffffff0198929800000000001976a91482932cf55b847ffa52832d2bbec2838f658f226788ac00000000",
						}},
					}, {
						TxID:              "c8a087b1ee775fa29697511ecd64e800941c8a22db6ed0989fb27a1d2d6798da",
						ReturnResult:      "success",
						ResultDescription: "",
						ConflictedWith:    nil,
					}},
					FailureCount: 1,
				},
				Signature: "304402200c4b0dc179906581eb32953abeddbef5799d302d82367aa9a469d79c15f932f3022029e827af6122290d5e1b80c50676b0336f4c7658ca67ba4819396dff9c6239a6",
				PublicKey: "030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e",
				Encoding:  "UTF-8",
				MimeType:  "application/json",
			},
			mockHTTPFn: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(submitResponse))),
				}, nil
			},
		},
	}

	for name, test := range tests {
		client := newTestClient(&MockClient{MockDo: test.mockHTTPFn})

		t.Run(name, func(t *testing.T) {
			result, err := client.SubmitTransactions(context.Background(), client.MinerByName(test.miner), test.txs)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.exp, result)
		})
	}
}
