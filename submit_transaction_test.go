package minercraft

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPValidSubmission for mocking requests
type mockHTTPValidSubmission struct{}

// Do is a mock http request
func (m *mockHTTPValidSubmission) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/tx") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-01-15T11:40:29.826Z\",\"txid\":\"6bdbcfab0526d30e8d68279f79dff61fb4026ace8b7b32789af016336e54f2f0\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"minerId\":\"03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031\",\"currentHighestBlockHash\":\"71a7374389afaec80fcabbbf08dcd82d392cf68c9a13fe29da1a0c853facef01\",\"currentHighestBlockHeight\":207,\"txSecondMempoolExpiry\":0}",
    	"signature": "3045022100f65ae83b20bc60e7a5f0e9c1bd9aceb2b26962ad0ee35472264e83e059f4b9be022010ca2334ff088d6e085eb3c2118306e61ec97781e8e1544e75224533dcc32379",
    	"publicKey": "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPBadSubmission for mocking requests
type mockHTTPBadSubmission struct{}

// Do is a mock http request
func (m *mockHTTPBadSubmission) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/tx") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{}",
    	"signature": "3044022066a8a39ff5f5eae818636aa03fdfc386ea4f33f41993cf41d4fb6d4745ae032102206a8895a6f742d809647ad1a1df12230e9b480275853ed28bc178f4b48afd802a",
    	"publicKey": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// TestClient_SubmitTransaction tests the method SubmitTransaction()
func TestClient_SubmitTransaction(t *testing.T) {
	t.Parallel()

	testSignature := "3045022100f65ae83b20bc60e7a5f0e9c1bd9aceb2b26962ad0ee35472264e83e059f4b9be022010ca2334ff088d6e085eb3c2118306e61ec97781e8e1544e75224533dcc32379"
	testPublicKey := "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031"

	// Create a client
	client := newTestClient(&mockHTTPValidSubmission{})

	tx := &Transaction{
		RawTx: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000",
	}

	// Create a req
	response, err := client.SubmitTransaction(client.MinerByName(MinerMatterpool), tx)
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

// ExampleClient_SubmitTransaction example using SubmitTransaction()
func ExampleClient_SubmitTransaction() {
	// Create a client (using a test client vs NewClient())
	client := newTestClient(&mockHTTPValidSubmission{})

	tx := &Transaction{
		RawTx: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000",
	}

	// Create a req
	response, err := client.SubmitTransaction(client.MinerByName(MinerTaal), tx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("submitted tx to: %s", response.Miner.Name)
	// Output:submitted tx to: Taal
}

// BenchmarkClient_SubmitTransaction benchmarks the method SubmitTransaction()
func BenchmarkClient_SubmitTransaction(b *testing.B) {
	client := newTestClient(&mockHTTPValidSubmission{})
	miner := client.MinerByName(MinerTaal)
	tx := &Transaction{
		RawTx: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000",
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.SubmitTransaction(miner, tx)
	}
}

// TestClient_SubmitTransactionParsedValues tests the method SubmitTransaction()
func TestClient_SubmitTransactionParsedValues(t *testing.T) {
	t.Parallel()

	testID := "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031"

	tx := &Transaction{
		RawTx: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000",
	}

	// Create a client
	client := newTestClient(&mockHTTPValidSubmission{})

	// Create a req
	response, err := client.SubmitTransaction(client.MinerByName(MinerMatterpool), tx)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if response == nil {
		t.Fatalf("expected response to not be nil")
	}

	// Test parsed values
	if response.Miner.Name != MinerMatterpool {
		t.Fatalf("expected response.Miner.Name to be %s, got %s", MinerTaal, response.Miner.Name)
	}
	if response.Results.MinerID != testID {
		t.Fatalf("expected response.Results.MinerID to be %s, got %s", testID, response.Results.MinerID)
	}
	if response.Results.Timestamp != "2020-01-15T11:40:29.826Z" {
		t.Fatalf("expected response.Results.Timestamp to be %s, got %s", "2020-01-15T11:40:29.826Z", response.Results.Timestamp)
	}
	if response.Results.APIVersion != testAPIVersion {
		t.Fatalf("expected response.Results.APIVersion to be %s, got %s", testAPIVersion, response.Results.APIVersion)
	}
	if response.Results.ReturnResult != QueryTransactionSuccess {
		t.Fatalf("expected response.Results.ReturnResult to be %s, got %s", "success", response.Results.ReturnResult)
	}
	if response.Results.TxID != "6bdbcfab0526d30e8d68279f79dff61fb4026ace8b7b32789af016336e54f2f0" {
		t.Fatalf("expected response.Results.TxID to be %s, got %s", "6bdbcfab0526d30e8d68279f79dff61fb4026ace8b7b32789af016336e54f2f0", response.Results.TxID)
	}
}

// TestClient_SubmitTransactionInvalidMiner tests the method SubmitTransaction()
func TestClient_SubmitTransactionInvalidMiner(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPValidSubmission{})

	tx := &Transaction{
		RawTx: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000",
	}

	// Create a req
	response, err := client.SubmitTransaction(nil, tx)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_SubmitTransactionHTTPError tests the method SubmitTransaction()
func TestClient_SubmitTransactionHTTPError(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPError{})

	tx := &Transaction{
		RawTx: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000",
	}

	// Create a req
	response, err := client.SubmitTransaction(client.MinerByName(MinerMatterpool), tx)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_SubmitTransactionBadRequest tests the method SubmitTransaction()
func TestClient_SubmitTransactionBadRequest(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPBadRequest{})

	tx := &Transaction{
		RawTx: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000",
	}

	// Create a req
	response, err := client.SubmitTransaction(client.MinerByName(MinerMatterpool), tx)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_SubmitTransactionInvalidJSON tests the method SubmitTransaction()
func TestClient_SubmitTransactionInvalidJSON(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPInvalidJSON{})

	tx := &Transaction{
		RawTx: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000",
	}

	// Create a req
	response, err := client.SubmitTransaction(client.MinerByName(MinerMatterpool), tx)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_SubmitTransactionInvalidSignature tests the method SubmitTransaction()
func TestClient_SubmitTransactionInvalidSignature(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPInvalidSignature{})

	tx := &Transaction{
		RawTx: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000",
	}

	// Create a req
	response, err := client.SubmitTransaction(client.MinerByName(MinerMatterpool), tx)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_SubmitTransactionBadSubmit tests the method SubmitTransaction()
func TestClient_SubmitTransactionBadSubmission(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPBadSubmission{})

	tx := &Transaction{
		RawTx: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1c03d7c6082f7376706f6f6c2e636f6d2f3edff034600055b8467f0040ffffffff01247e814a000000001976a914492558fb8ca71a3591316d095afc0f20ef7d42f788ac00000000",
	}

	// Create a req
	response, err := client.SubmitTransaction(client.MinerByName(MinerMatterpool), tx)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}
