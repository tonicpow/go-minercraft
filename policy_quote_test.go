package minercraft

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

// mockHTTPValidPolicyQuote for mocking requests
type mockHTTPValidPolicyQuote struct{}

// Do is a mock http request
func (m *mockHTTPValidPolicyQuote) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/policyQuote") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    "payload": "{\"apiVersion\":\"1.4.0\",\"timestamp\":\"2021-11-12T13:17:47.7498672Z\",\"expiryTime\":\"2021-11-12T13:27:47.7498672Z\",\"minerId\":\"030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e\",\"currentHighestBlockHash\":\"45628be2fe616167b7da399ab63455e60ffcf84147730f4af4affca90c7d437e\",\"currentHighestBlockHeight\":234,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}}],\"callbacks\":[{\"ipAddress\":\"123.456.789.123\"}],\"policies\":{\"skipscriptflags\":[\"MINIMALDATA\",\"DERSIG\",\"NULLDUMMY\",\"DISCOURAGE_UPGRADABLE_NOPS\",\"CLEANSTACK\"],\"maxtxsizepolicy\":99999,\"datacarriersize\":100000,\"maxscriptsizepolicy\":100000,\"maxscriptnumlengthpolicy\":100000,\"maxstackmemoryusagepolicy\":10000000,\"limitancestorcount\":1000,\"limitcpfpgroupmemberscount\":10,\"acceptnonstdoutputs\":true,\"datacarrier\":true,\"dustrelayfee\":150,\"maxstdtxvalidationduration\":99,\"maxnonstdtxvalidationduration\":100,\"dustlimitfactor\":10}}",
    "signature": "30440220708e2e62a393f53c43d172bc1459b4daccf9cf23ff77cff923f09b2b49b94e0a022033792bee7bc3952f4b1bfbe9df6407086b5dbfc161df34fdee684dc97be72731",
    "publicKey": "030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e",
    "encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// TestClient_PolicyQuote tests the method PolicyQuote()
func TestClient_PolicyQuote(t *testing.T) {

	t.Run("get a valid policy quote", func(t *testing.T) {

		defer goleak.VerifyNone(t)

		// Create a client
		client := newTestClient(&mockHTTPValidPolicyQuote{})

		// Create a req
		response, err := client.PolicyQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Check returned values
		assert.Equal(t, "30440220708e2e62a393f53c43d172bc1459b4daccf9cf23ff77cff923f09b2b49b94e0a022033792bee7bc3952f4b1bfbe9df6407086b5dbfc161df34fdee684dc97be72731", *response.Signature)
		assert.Equal(t, "030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e", *response.PublicKey)
		assert.Equal(t, testEncoding, response.Encoding)
		assert.Equal(t, testMimeType, response.MimeType)

		// Test flags
		flags := []ScriptFlag{FlagMinimalData, FlagDerSig, FlagNullDummy, FlagDiscourageUpgradableNops, FlagCleanStack}

		// Confirm all policy fields are the right type
		assert.Equal(t, true, response.Quote.Policies.AcceptNonStdOutputs)
		assert.Equal(t, true, response.Quote.Policies.DataCarrier)
		assert.Equal(t, uint32(10), response.Quote.Policies.DustLimitFactor)
		assert.Equal(t, uint32(10), response.Quote.Policies.LimitCpfpGroupMembersCount)
		assert.Equal(t, uint32(100), response.Quote.Policies.MaxNonStdTxValidationDuration)
		assert.Equal(t, uint32(1000), response.Quote.Policies.LimitAncestorCount)
		assert.Equal(t, uint32(100000), response.Quote.Policies.DataCarrierSize)
		assert.Equal(t, uint32(100000), response.Quote.Policies.MaxScriptNumLengthPolicy)
		assert.Equal(t, uint32(100000), response.Quote.Policies.MaxScriptSizePolicy)
		assert.Equal(t, uint32(150), response.Quote.Policies.DustRelayFee)
		assert.Equal(t, uint32(99), response.Quote.Policies.MaxStdTxValidationDuration)
		assert.Equal(t, uint32(99999), response.Quote.Policies.MaxTxSizePolicy)
		assert.Equal(t, uint64(10000000), response.Quote.Policies.MaxStackMemoryUsagePolicy)
		assert.Equal(t, flags, response.Quote.Policies.SkipScriptFlags)

		assert.Equal(t, 500, response.Quote.FeePayload.GetFee(FeeTypeStandard).MiningFee.Satoshis)
		assert.Equal(t, 1000, response.Quote.FeePayload.GetFee(FeeTypeStandard).MiningFee.Bytes)
	})

	t.Run("invalid miner", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPValidPolicyQuote{})
		response, err := client.PolicyQuote(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("http error", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPError{})
		response, err := client.PolicyQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("bad request", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPBadRequest{})
		response, err := client.PolicyQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPInvalidJSON{})
		response, err := client.PolicyQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("invalid signature", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPInvalidSignature{})
		response, err := client.PolicyQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("missing fees", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := newTestClient(&mockHTTPMissingFees{})
		response, err := client.PolicyQuote(context.Background(), client.MinerByName(MinerTaal))
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

// ExampleClient_FeeQuote example using PolicyQuote()
func ExampleClient_PolicyQuote() {
	// Create a client (using a test client vs NewClient())
	client := newTestClient(&mockHTTPValidPolicyQuote{})

	// Create a req
	response, err := client.PolicyQuote(context.Background(), client.MinerByName(MinerTaal))
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("got quote from: %s", response.Miner.Name)
	// Output:got quote from: Taal
}

// BenchmarkClient_FeeQuote benchmarks the method PolicyQuote()
func BenchmarkClient_PolicyQuote(b *testing.B) {
	client := newTestClient(&mockHTTPValidPolicyQuote{})
	for i := 0; i < b.N; i++ {
		_, _ = client.PolicyQuote(context.Background(), client.MinerByName(MinerTaal))
	}
}
