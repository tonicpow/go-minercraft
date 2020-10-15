package minercraft

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPValidFeeQuote for mocking requests
type mockHTTPValidFeeQuote struct{}

// Do is a mock http request
func (m *mockHTTPValidFeeQuote) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/feeQuote") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T21:26:17.410Z\",\"expiryTime\":\"2020-10-09T21:36:17.410Z\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565\",\"currentHighestBlockHeight\":656169,\"minerReputation\":null,\"fees\":[{\"id\":1,\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"id\":2,\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}}]}",
   	 	"signature": "3045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc4",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPError for mocking requests
type mockHTTPError struct{}

// Do is a mock http request
func (m *mockHTTPError) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	if req.URL.String() != "" {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf(`http timeout`)
	}

	// Default is valid
	return resp, nil
}

// mockHTTPBadRequest for mocking requests
type mockHTTPBadRequest struct{}

// Do is a mock http request
func (m *mockHTTPBadRequest) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	if req.URL.String() != "" {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPInvalidJSON for mocking requests
type mockHTTPInvalidJSON struct{}

// Do is a mock http request
func (m *mockHTTPInvalidJSON) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	if req.URL.String() != "" {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{invalid:json}`)))
		resp.StatusCode = http.StatusOK
	}

	// Default is valid
	return resp, nil
}

// mockHTTPMissingFees for mocking requests
type mockHTTPMissingFees struct{}

// Do is a mock http request
func (m *mockHTTPMissingFees) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/feeQuote") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T21:26:17.410Z\",\"expiryTime\":\"2020-10-09T21:36:17.410Z\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565\",\"currentHighestBlockHeight\":656169,\"minerReputation\":null,\"fees\":[]}",
   	 	"signature": "3045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc4",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPInvalidSignature for mocking requests
type mockHTTPInvalidSignature struct{}

// Do is a mock http request
func (m *mockHTTPInvalidSignature) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Invalid sig response
	if strings.Contains(req.URL.String(), "/mapi/feeQuote") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T21:26:17.410Z\",\"expiryTime\":\"2020-10-09T21:36:17.410Z\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565\",\"currentHighestBlockHeight\":656169,\"minerReputation\":null,\"fees\":[{\"id\":1,\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"id\":2,\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}}]}",
   	 	"signature": "03045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc40",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Invalid sig response
	if strings.Contains(req.URL.String(), "/mapi/tx/"+testTx) {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-11T15:41:29.269Z\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"blockHash\":\"0000000000000000050a09fe90b0e8542bba9e712edb8cc9349e61888fe45ac5\",\"blockHeight\":612530,\"confirmations\":43923,\"minerId\":\"0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087\",\"txSecondMempoolExpiry\":0}",
   	 	"signature": "03045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc40",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Invalid sig response
	if strings.Contains(req.URL.String(), "/mapi/tx") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-01-15T11:40:29.826Z\",\"txid\":\"6bdbcfab0526d30e8d68279f79dff61fb4026ace8b7b32789af016336e54f2f0\",\"returnResult\":\"success\",\"resultDescription\":\"\",\"minerId\":\"03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031\",\"currentHighestBlockHash\":\"71a7374389afaec80fcabbbf08dcd82d392cf68c9a13fe29da1a0c853facef01\",\"currentHighestBlockHeight\":207,\"txSecondMempoolExpiry\":0}",
    	"signature": "03045022100f65ae83b20bc60e7a5f0e9c1bd9aceb2b26962ad0ee35472264e83e059f4b9be022010ca2334ff088d6e085eb3c2118306e61ec97781e8e1544e75224533dcc323790",
    	"publicKey": "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPBetterRate for mocking requests
type mockHTTPBetterRate struct{}

// Do is a mock http request
func (m *mockHTTPBetterRate) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if req.URL.String() == defaultProtocol+"merchantapi.taal.com/mapi/feeQuote" {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T21:26:17.410Z\",\"expiryTime\":\"2020-10-09T21:36:17.410Z\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565\",\"currentHighestBlockHeight\":656169,\"minerReputation\":null,\"fees\":[{\"id\":1,\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":475,\"bytes\":1000},\"relayFee\":{\"satoshis\":150,\"bytes\":1000}},{\"id\":2,\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}}]}",
   	 	"signature": "3045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc4",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	if req.URL.String() == defaultProtocol+"merchantapi.matterpool.io/mapi/feeQuote" {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T22:08:26.236Z\",\"expiryTime\":\"2020-10-09T22:18:26.236Z\",\"minerId\":\"0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087\",\"currentHighestBlockHash\":\"0000000000000000028285a9168c95457521a743765f499de389c094e883f42a\",\"currentHighestBlockHeight\":656171,\"minerReputation\":null,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":405,\"bytes\":1000},\"relayFee\":{\"satoshis\":100,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":100,\"bytes\":1000}}]}",
    	"signature": "3044022011f90db2661726eb2659c3447ccaa9fd3368194f87d5d86a23e673c45d5d714502200c51eb600e3370b49d759aa4d441000286937b0803037a1d6de4c5a5c559d74c",
    	"publicKey": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	if req.URL.String() == defaultProtocol+"www.ddpurse.com/openapi/mapi/feeQuote" {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T22:09:04.433Z\",\"expiryTime\":\"2020-10-09T22:19:04.433Z\",\"minerId\":null,\"currentHighestBlockHash\":\"0000000000000000028285a9168c95457521a743765f499de389c094e883f42a\",\"currentHighestBlockHeight\":656171,\"minerReputation\":null,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":350,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":430,\"bytes\":1000},\"relayFee\":{\"satoshis\":175,\"bytes\":1000}}]}",
    	"signature": null,"publicKey": null,"encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPMissingFeeType for mocking requests
type mockHTTPMissingFeeType struct{}

// Do is a mock http request
func (m *mockHTTPMissingFeeType) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid response
	if strings.Contains(req.URL.String(), "/mapi/feeQuote") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T21:26:17.410Z\",\"expiryTime\":\"2020-10-09T21:36:17.410Z\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565\",\"currentHighestBlockHeight\":656169,\"minerReputation\":null,\"fees\":[{\"id\":2,\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}}]}",
   	 	"signature": "3045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc4",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// TestClient_FeeQuote tests the method FeeQuote()
func TestClient_FeeQuote(t *testing.T) {
	t.Parallel()

	testSignature := "3045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc4"
	testPublicKey := "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270"

	// Create a client
	client := newTestClient(&mockHTTPValidFeeQuote{})

	// Create a req
	response, err := client.FeeQuote(client.MinerByName(MinerTaal))
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

// TestClient_FeeQuoteParsedValues tests the method FeeQuote()
func TestClient_FeeQuoteParsedValues(t *testing.T) {
	t.Parallel()

	testID := "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270"

	// Create a client
	client := newTestClient(&mockHTTPValidFeeQuote{})

	// Create a req
	response, err := client.FeeQuote(client.MinerByName(MinerTaal))
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if response == nil {
		t.Fatalf("expected response to not be nil")
	}

	// Test parsed values
	if response.Miner.Name != MinerTaal {
		t.Fatalf("expected response.Miner.Name to be %s, got %s", MinerTaal, response.Miner.Name)
	}
	if response.Quote.MinerID != testID {
		t.Fatalf("expected response.Quote.MinerID to be %s, got %s", testID, response.Quote.MinerID)
	}
	if response.Quote.ExpirationTime != "2020-10-09T21:36:17.410Z" {
		t.Fatalf("expected response.Quote.ExpirationTime to be %s, got %s", "2020-10-09T21:36:17.410Z", response.Quote.ExpirationTime)
	}
	if response.Quote.Timestamp != "2020-10-09T21:26:17.410Z" {
		t.Fatalf("expected response.Quote.Timestamp to be %s, got %s", "2020-10-09T21:26:17.410Z", response.Quote.Timestamp)
	}
	if response.Quote.APIVersion != testAPIVersion {
		t.Fatalf("expected response.Quote.APIVersion to be %s, got %s", testAPIVersion, response.Quote.APIVersion)
	}
	if response.Quote.CurrentHighestBlockHash != "0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565" {
		t.Fatalf("expected response.Quote.CurrentHighestBlockHash to be %s, got %s", "0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565", response.Quote.CurrentHighestBlockHash)
	}
	if response.Quote.CurrentHighestBlockHeight != 656169 {
		t.Fatalf("expected response.Quote.CurrentHighestBlockHeight to be %d, got %d", 656169, response.Quote.CurrentHighestBlockHeight)
	}
	if len(response.Quote.Fees) != 2 {
		t.Fatalf("expected response.Quote.Fees to be length of %d, got: %d", 2, len(response.Quote.Fees))
	}
}

// TestClient_FeeQuoteGetRate tests the method FeeQuote()
func TestClient_FeeQuoteGetRate(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPValidFeeQuote{})

	// Create a req
	response, err := client.FeeQuote(client.MinerByName(MinerTaal))
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if response == nil {
		t.Fatalf("expected response to not be nil")
	}

	// Test getting rate from request
	var rate uint64
	rate, err = response.Quote.CalculateFee(FeeCategoryMining, FeeTypeData, 1000)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if rate != 500 {
		t.Fatalf("rate was %d but expected: %d", rate, 500)
	}

	// Test relay rate
	rate, err = response.Quote.CalculateFee(FeeCategoryRelay, FeeTypeData, 1000)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if rate != 250 {
		t.Fatalf("rate was %d but expected: %d", rate, 250)
	}

}

// TestClient_FeeQuoteInvalidMiner tests the method FeeQuote()
func TestClient_FeeQuoteInvalidMiner(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPValidFeeQuote{})

	// Create a req
	response, err := client.FeeQuote(nil)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_FeeQuoteHTTPError tests the method FeeQuote()
func TestClient_FeeQuoteHTTPError(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPError{})

	// Create a req
	response, err := client.FeeQuote(client.MinerByName(MinerTaal))
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_FeeQuoteBadRequest tests the method FeeQuote()
func TestClient_FeeQuoteBadRequest(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPBadRequest{})

	// Create a req
	response, err := client.FeeQuote(client.MinerByName(MinerTaal))
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_FeeQuoteInvalidJSON tests the method FeeQuote()
func TestClient_FeeQuoteInvalidJSON(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPInvalidJSON{})

	// Create a req
	response, err := client.FeeQuote(client.MinerByName(MinerTaal))
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_FeeQuoteInvalidSignature tests the method FeeQuote()
func TestClient_FeeQuoteInvalidSignature(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPInvalidSignature{})

	// Create a req
	response, err := client.FeeQuote(client.MinerByName(MinerTaal))
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_FeeQuoteMissingFees tests the method FeeQuote()
func TestClient_FeeQuoteMissingFees(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPMissingFees{})

	// Create a req
	response, err := client.FeeQuote(client.MinerByName(MinerTaal))
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// ExampleClient_FeeQuote example using FeeQuote()
func ExampleClient_FeeQuote() {
	// Create a client (using a test client vs NewClient())
	client := newTestClient(&mockHTTPValidFeeQuote{})

	// Create a req
	response, err := client.FeeQuote(client.MinerByName(MinerTaal))
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("got quote from: %s", response.Miner.Name)
	// Output:got quote from: Taal
}

// BenchmarkClient_FeeQuote benchmarks the method FeeQuote()
func BenchmarkClient_FeeQuote(b *testing.B) {
	client := newTestClient(&mockHTTPValidFeeQuote{})
	for i := 0; i < b.N; i++ {
		_, _ = client.FeeQuote(client.MinerByName(MinerTaal))
	}
}

// TestFeePayload_CalculateFee tests the method CalculateFee()
func TestFeePayload_CalculateFee(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPValidFeeQuote{})

	// Create a req
	response, err := client.FeeQuote(client.MinerByName(MinerTaal))
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if response == nil {
		t.Fatalf("expected response to not be nil")
	}

	// Mining & Data
	var fee uint64
	fee, err = response.Quote.CalculateFee(FeeCategoryMining, FeeTypeData, 1000)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if fee != 500 {
		t.Fatalf("fee was: %d but expected: %d", fee, 500)
	}

	// Mining and standard
	fee, err = response.Quote.CalculateFee(FeeCategoryMining, FeeTypeStandard, 1000)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if fee != 500 {
		t.Fatalf("fee was: %d but expected: %d", fee, 500)
	}

	// Relay & Data
	fee, err = response.Quote.CalculateFee(FeeCategoryRelay, FeeTypeData, 1000)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if fee != 250 {
		t.Fatalf("fee was: %d but expected: %d", fee, 250)
	}

	// Relay and standard
	fee, err = response.Quote.CalculateFee(FeeCategoryRelay, FeeTypeStandard, 1000)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if fee != 250 {
		t.Fatalf("fee was: %d but expected: %d", fee, 250)
	}
}

// ExampleFeePayload_CalculateFee example using CalculateFee()
func ExampleFeePayload_CalculateFee() {
	// Create a client (using a test client vs NewClient())
	client := newTestClient(&mockHTTPValidBestQuote{})

	// Create a req
	response, err := client.BestQuote(FeeCategoryMining, FeeTypeData)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Calculate fee for tx
	var fee uint64
	fee, err = response.Quote.CalculateFee(FeeCategoryMining, FeeTypeData, 1000)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Note: cannot show response since the miner might be different each time
	fmt.Printf("got best quote and fee for 1000 byte tx is: %d", fee)
	// Output:got best quote and fee for 1000 byte tx is: 420
}

// BenchmarkFeePayload_CalculateFee benchmarks the method CalculateFee()
func BenchmarkFeePayload_CalculateFee(b *testing.B) {
	client := newTestClient(&mockHTTPValidBestQuote{})
	response, _ := client.BestQuote(FeeCategoryMining, FeeTypeData)
	for i := 0; i < b.N; i++ {
		_, _ = response.Quote.CalculateFee(FeeCategoryMining, FeeTypeData, 1000)
	}
}

// TestFeePayload_CalculateFeeZero tests the method CalculateFee()
func TestFeePayload_CalculateFeeZero(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPValidFeeQuote{})

	// Create a req
	response, err := client.FeeQuote(client.MinerByName(MinerTaal))
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if response == nil {
		t.Fatalf("expected response to not be nil")
	}

	// Zero tx size produces 0 fee and error
	var fee uint64
	fee, err = response.Quote.CalculateFee(FeeCategoryMining, FeeTypeData, 0)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if fee != 1 {
		t.Fatalf("fee was: %d but expected: %d", fee, 1)
	}
}

// TestFeePayload_CalculateFeeMissingFeeType tests the method CalculateFee()
func TestFeePayload_CalculateFeeMissingFeeType(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPMissingFeeType{})

	// Create a req
	response, err := client.FeeQuote(client.MinerByName(MinerTaal))
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if response == nil {
		t.Fatalf("expected response to not be nil")
	}

	// Zero tx size produces 0 fee and error
	var fee uint64
	fee, err = response.Quote.CalculateFee(FeeCategoryRelay, FeeTypeStandard, 1000)
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if fee != 1 {
		t.Fatalf("fee was: %d but expected: %d", fee, 1)
	}
}
