package minercraft

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPValidQuery for mocking requests
type mockHTTPValidQuery struct{}

// Do is a mock http request
func (m *mockHTTPValidQuery) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/tx/"+testTx) {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-10T13:07:26.014Z\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"blockHash\":\"0000000000000000050a09fe90b0e8542bba9e712edb8cc9349e61888fe45ac5\",\"blockHeight\":612530,\"confirmations\":43733,\"minerId\":\"0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087\",\"txSecondMempoolExpiry\":0}",
    	"signature": "3044022066a8a39ff5f5eae818636aa03fdfc386ea4f33f41993cf41d4fb6d4745ae032102206a8895a6f742d809647ad1a1df12230e9b480275853ed28bc178f4b48afd802a",
    	"publicKey": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPBadQuery for mocking requests
type mockHTTPBadQuery struct{}

// Do is a mock http request
func (m *mockHTTPBadQuery) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/tx/"+testTx) {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{}",
    	"signature": "3044022066a8a39ff5f5eae818636aa03fdfc386ea4f33f41993cf41d4fb6d4745ae032102206a8895a6f742d809647ad1a1df12230e9b480275853ed28bc178f4b48afd802a",
    	"publicKey": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// TestClient_QueryTransaction tests the method QueryTransaction()
func TestClient_QueryTransaction(t *testing.T) {
	t.Parallel()

	testSignature := "3044022066a8a39ff5f5eae818636aa03fdfc386ea4f33f41993cf41d4fb6d4745ae032102206a8895a6f742d809647ad1a1df12230e9b480275853ed28bc178f4b48afd802a"
	testPublicKey := "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087"

	// Create a client
	client := newTestClient(&mockHTTPValidQuery{})

	// Create a req
	response, err := client.QueryTransaction(client.MinerByName(MinerMatterpool), testTx)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if response == nil {
		t.Fatalf("expected response to not be nil")
	}

	// Check returned values
	if !response.Validated {
		t.Fatalf("expected response.Validated to be true, got false")
	}
	if response.Signature != testSignature {
		t.Fatalf("expected response.Signature to be %s, got %s", testSignature, response.Signature)
	}
	if response.PublicKey != testPublicKey {
		t.Fatalf("expected response.PublicKey to be %s, got %s", testPublicKey, response.PublicKey)
	}
	if response.Encoding != testEncoding {
		t.Fatalf("expected response.Encoding to be %s, got %s", testEncoding, response.Encoding)
	}
	if response.MimeType != testMimeType {
		t.Fatalf("expected response.MimeType to be %s, got %s", testMimeType, response.MimeType)
	}
}

// ExampleClient_QueryTransaction example using QueryTransaction()
func ExampleClient_QueryTransaction() {
	// Create a client (using a test client vs NewClient())
	client := newTestClient(&mockHTTPValidQuery{})

	// Create a req
	response, err := client.QueryTransaction(client.MinerByName(MinerTaal), testTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("got tx status %s from: %s", response.Query.ReturnResult, response.Miner.Name)
	// Output:got tx status success from: Taal
}

// BenchmarkClient_QueryTransaction benchmarks the method QueryTransaction()
func BenchmarkClient_QueryTransaction(b *testing.B) {
	client := newTestClient(&mockHTTPValidQuery{})
	miner := client.MinerByName(MinerTaal)
	for i := 0; i < b.N; i++ {
		_, _ = client.QueryTransaction(miner, testTx)
	}
}

// TestClient_QueryTransactionParsedValues tests the method QueryTransaction()
func TestClient_QueryTransactionParsedValues(t *testing.T) {
	t.Parallel()

	testID := "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087"

	// Create a client
	client := newTestClient(&mockHTTPValidQuery{})

	// Create a req
	response, err := client.QueryTransaction(client.MinerByName(MinerMatterpool), testTx)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if response == nil {
		t.Fatalf("expected response to not be nil")
	}

	// Test parsed values
	if response.Miner.Name != MinerMatterpool {
		t.Fatalf("expected response.Miner.Name to be %s, got %s", MinerTaal, response.Miner.Name)
	}
	if response.Query.MinerID != testID {
		t.Fatalf("expected response.Query.MinerID to be %s, got %s", testID, response.Query.MinerID)
	}
	if response.Query.Timestamp != "2020-10-10T13:07:26.014Z" {
		t.Fatalf("expected response.Query.Timestamp to be %s, got %s", "2020-10-10T13:07:26.014Z", response.Query.Timestamp)
	}
	if response.Query.APIVersion != testAPIVersion {
		t.Fatalf("expected response.Query.APIVersion to be %s, got %s", testAPIVersion, response.Query.APIVersion)
	}
	if response.Query.BlockHash != "0000000000000000050a09fe90b0e8542bba9e712edb8cc9349e61888fe45ac5" {
		t.Fatalf("expected response.Query.BlockHash to be %s, got %s", "0000000000000000050a09fe90b0e8542bba9e712edb8cc9349e61888fe45ac5", response.Query.BlockHash)
	}
	if response.Query.BlockHeight != 612530 {
		t.Fatalf("expected response.Query.BlockHeight to be %d, got %d", 612530, response.Query.BlockHeight)
	}
	if response.Query.ReturnResult != QueryTransactionSuccess {
		t.Fatalf("expected response.Query.ReturnResult to be %s, got %s", "success", response.Query.ReturnResult)
	}
}

// TestClient_QueryTransactionInvalidMiner tests the method QueryTransaction()
func TestClient_QueryTransactionInvalidMiner(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPValidFeeQuote{})

	// Create a req
	response, err := client.QueryTransaction(nil, testTx)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_QueryTransactionHTTPError tests the method QueryTransaction()
func TestClient_QueryTransactionHTTPError(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPError{})

	// Create a req
	response, err := client.QueryTransaction(client.MinerByName(MinerMatterpool), testTx)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_QueryTransactionBadRequest tests the method QueryTransaction()
func TestClient_QueryTransactionBadRequest(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPBadRequest{})

	// Create a req
	response, err := client.QueryTransaction(client.MinerByName(MinerMatterpool), testTx)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_QueryTransactionInvalidJSON tests the method QueryTransaction()
func TestClient_QueryTransactionInvalidJSON(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPInvalidJSON{})

	// Create a req
	response, err := client.QueryTransaction(client.MinerByName(MinerMatterpool), testTx)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_QueryTransactionInvalidSignature tests the method QueryTransaction()
func TestClient_QueryTransactionInvalidSignature(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPInvalidSignature{})

	// Create a req
	response, err := client.QueryTransaction(client.MinerByName(MinerMatterpool), testTx)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_QueryTransactionBadQuery tests the method QueryTransaction()
func TestClient_QueryTransactionBadQuery(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPBadQuery{})

	// Create a req
	response, err := client.QueryTransaction(client.MinerByName(MinerMatterpool), testTx)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}
