package minercraft

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// mockHTTPValidFastestQuote for mocking requests
type mockHTTPValidFastestQuote struct{}

// Do is a mock http request
func (m *mockHTTPValidFastestQuote) Do(req *http.Request) (*http.Response, error) {
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
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T21:26:17.410Z\",\"expiryTime\":\"2020-10-09T21:36:17.410Z\",\"minerId\":\"03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270\",\"currentHighestBlockHash\":\"0000000000000000035c5f8c0294802a01e500fa7b95337963bb3640da3bd565\",\"currentHighestBlockHeight\":656169,\"minerReputation\":null,\"fees\":[{\"id\":1,\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":400,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"id\":2,\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":225,\"bytes\":1000}}]}",
   	 	"signature": "3045022100eed49f6bf75d8f975f581271e3df658fbe8ec67e6301ea8fc25a72d18c92e30e022056af253f0d24db6a8fde4e2c1ee95e7a5ecf2c7cdc93246f8328c9e0ca582fc4",
    	"publicKey": "03e92d3e5c3f7bd945dfbf48e7a99393b1bfb3f11f380ae30d286e7ff2aec5a270","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	if req.URL.String() == defaultProtocol+"merchantapi.matterpool.io/mapi/feeQuote" {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T22:08:26.236Z\",\"expiryTime\":\"2020-10-09T22:18:26.236Z\",\"minerId\":\"0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087\",\"currentHighestBlockHash\":\"0000000000000000028285a9168c95457521a743765f499de389c094e883f42a\",\"currentHighestBlockHeight\":656171,\"minerReputation\":null,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":400,\"bytes\":1000},\"relayFee\":{\"satoshis\":100,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":430,\"bytes\":1000},\"relayFee\":{\"satoshis\":110,\"bytes\":1000}}]}",
    	"signature": "3044022011f90db2661726eb2659c3447ccaa9fd3368194f87d5d86a23e673c45d5d714502200c51eb600e3370b49d759aa4d441000286937b0803037a1d6de4c5a5c559d74c",
    	"publicKey": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087","encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	if req.URL.String() == defaultProtocol+"www.ddpurse.com/openapi/mapi/feeQuote" {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    	"payload": "{\"apiVersion\":\"` + testAPIVersion + `\",\"timestamp\":\"2020-10-09T22:09:04.433Z\",\"expiryTime\":\"2020-10-09T22:19:04.433Z\",\"minerId\":null,\"currentHighestBlockHash\":\"0000000000000000028285a9168c95457521a743765f499de389c094e883f42a\",\"currentHighestBlockHeight\":656171,\"minerReputation\":null,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":420,\"bytes\":1000},\"relayFee\":{\"satoshis\":150,\"bytes\":1000}}]}",
    	"signature": null,"publicKey": null,"encoding": "` + testEncoding + `","mimetype": "` + testMimeType + `"}`)))
	}

	// Default is valid
	return resp, nil
}

// TestClient_FastestQuote tests the method FastestQuote()
func TestClient_FastestQuote(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPValidFastestQuote{})

	// Create a req
	response, err := client.FastestQuote()
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	} else if response == nil {
		t.Fatalf("expected response to not be nil")
	}

	// Check returned values
	if response.Encoding != testEncoding {
		t.Fatalf("expected response.Encoding to be %s, got %s", testEncoding, response.Encoding)
	}
	if response.MimeType != testMimeType {
		t.Fatalf("expected response.MimeType to be %s, got %s", testMimeType, response.MimeType)
	}

	// Check that we got fees
	if len(response.Quote.Fees) != 2 {
		t.Fatalf("expected response.Quote.Fees to be a length of %d, got %d", 2, len(response.Quote.Fees))
	}
}

// TestClient_FastestQuoteHTTPError tests the method FastestQuote()
func TestClient_FastestQuoteHTTPError(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPError{})

	// Create a req
	response, err := client.FastestQuote()
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_FastestQuoteBadRequest tests the method FastestQuote()
func TestClient_FastestQuoteBadRequest(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPBadRequest{})

	// Create a req
	response, err := client.FastestQuote()
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// TestClient_FastestQuoteInvalidJSON tests the method FastestQuote()
func TestClient_FastestQuoteInvalidJSON(t *testing.T) {
	t.Parallel()

	// Create a client
	client := newTestClient(&mockHTTPInvalidJSON{})

	// Create a req
	response, err := client.FastestQuote()
	if err == nil {
		t.Fatalf("error should have occurred")
	} else if response != nil {
		t.Fatalf("expected response to be nil")
	}
}

// ExampleClient_FastestQuote example using FastestQuote()
func ExampleClient_FastestQuote() {
	// Create a client (using a test client vs NewClient())
	client := newTestClient(&mockHTTPValidFastestQuote{})

	// Create a req
	_, err := client.FastestQuote()
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Note: cannot show response since the miner might be different each time
	fmt.Printf("got fastest quote!")
	// Output:got fastest quote!
}

// BenchmarkClient_FastestQuote benchmarks the method FastestQuote()
func BenchmarkClient_FastestQuote(b *testing.B) {
	client := newTestClient(&mockHTTPValidFastestQuote{})
	for i := 0; i < b.N; i++ {
		_, _ = client.FastestQuote()
	}
}
