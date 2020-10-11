package minercraft

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const (
	testEncoding   = "UTF-8"
	testMimeType   = "application/json"
	testMinerID    = "1234567"
	testMinerName  = "TestMiner"
	testMinerToken = "0987654321"
	testMinerURL   = defaultProtocol + "testminer.com"
	testTx         = "7e0c4651fc256c0433bd704d7e13d24c8d10235f4b28ba192849c5d318de974b"
)

// mockHTTPDefaultClient for mocking requests
type mockHTTPDefaultClient struct{}

// Do is a mock http request
func (m *mockHTTPDefaultClient) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	if req.URL.String() == "/test" {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"message":"test"}`)))
	}

	// Default is valid
	return resp, nil
}

// newTestClient returns a client for mocking (using a custom HTTP interface)
func newTestClient(httpClient httpInterface) *Client {
	client, _ := NewClient(nil, nil)
	client.httpClient = httpClient
	return client
}

// TestNewClient tests the method NewClient()
func TestNewClient(t *testing.T) {
	t.Parallel()

	client, err := NewClient(nil, nil)

	if client == nil {
		t.Fatal("failed to load client")
	} else if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Test default miners
	if len(client.Miners) != 3 {
		t.Fatalf("expected %d default miners, got %d", 3, len(client.Miners))
	}
}

// ExampleNewClient example using NewClient()
func ExampleNewClient() {
	client, err := NewClient(nil, nil)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("created new client with %d default miners", len(client.Miners))
	// Output:created new client with 3 default miners
}

// BenchmarkNewClient benchmarks the method NewClient()
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewClient(nil, nil)
	}
}

// TestNewClient_CustomHttpClient tests new client with custom HTTP client
func TestNewClient_CustomHttpClient(t *testing.T) {
	t.Parallel()

	client, err := NewClient(nil, http.DefaultClient)

	if client == nil {
		t.Fatal("failed to load client")
	} else if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
}

// TestNewClient_DefaultMiners tests NewClient with default miners
func TestNewClient_DefaultMiners(t *testing.T) {
	t.Parallel()

	client, err := NewClient(nil, http.DefaultClient)

	if client == nil {
		t.Fatal("failed to load client")
	} else if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Get Taal
	miner := client.MinerByName(MinerTaal)
	if miner.Name != MinerTaal {
		t.Fatalf("expected miner named %s", MinerTaal)
	}

	// Get Mempool
	miner = client.MinerByName(MinerMempool)
	if miner.Name != MinerMempool {
		t.Fatalf("expected miner named %s", MinerMempool)
	}

	// Get Matterpool
	miner = client.MinerByName(MinerMatterpool)
	if miner.Name != MinerMatterpool {
		t.Fatalf("expected miner named %s", MinerMatterpool)
	}
}

// TestDefaultClientOptions tests setting DefaultClientOptions()
func TestDefaultClientOptions(t *testing.T) {
	t.Parallel()

	options := DefaultClientOptions()

	if options.UserAgent != defaultUserAgent {
		t.Fatalf("expected value: %s got: %s", defaultUserAgent, options.UserAgent)
	}

	if options.BackOffExponentFactor != 2.0 {
		t.Fatalf("expected value: %f got: %f", 2.0, options.BackOffExponentFactor)
	}

	if options.BackOffInitialTimeout != 2*time.Millisecond {
		t.Fatalf("expected value: %v got: %v", 2*time.Millisecond, options.BackOffInitialTimeout)
	}

	if options.BackOffMaximumJitterInterval != 2*time.Millisecond {
		t.Fatalf("expected value: %v got: %v", 2*time.Millisecond, options.BackOffMaximumJitterInterval)
	}

	if options.BackOffMaxTimeout != 10*time.Millisecond {
		t.Fatalf("expected value: %v got: %v", 10*time.Millisecond, options.BackOffMaxTimeout)
	}

	if options.DialerKeepAlive != 20*time.Second {
		t.Fatalf("expected value: %v got: %v", 20*time.Second, options.DialerKeepAlive)
	}

	if options.DialerTimeout != 5*time.Second {
		t.Fatalf("expected value: %v got: %v", 5*time.Second, options.DialerTimeout)
	}

	if options.RequestRetryCount != 2 {
		t.Fatalf("expected value: %v got: %v", 2, options.RequestRetryCount)
	}

	if options.RequestTimeout != 10*time.Second {
		t.Fatalf("expected value: %v got: %v", 10*time.Second, options.RequestTimeout)
	}

	if options.TransportExpectContinueTimeout != 3*time.Second {
		t.Fatalf("expected value: %v got: %v", 3*time.Second, options.TransportExpectContinueTimeout)
	}

	if options.TransportIdleTimeout != 20*time.Second {
		t.Fatalf("expected value: %v got: %v", 20*time.Second, options.TransportIdleTimeout)
	}

	if options.TransportMaxIdleConnections != 10 {
		t.Fatalf("expected value: %v got: %v", 10, options.TransportMaxIdleConnections)
	}

	if options.TransportTLSHandshakeTimeout != 5*time.Second {
		t.Fatalf("expected value: %v got: %v", 5*time.Second, options.TransportTLSHandshakeTimeout)
	}
}

// ExampleDefaultClientOptions example using DefaultClientOptions()
func ExampleDefaultClientOptions() {
	options := DefaultClientOptions()
	options.UserAgent = "Custom UserAgent v1.0"
	client, err := NewClient(options, nil)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("created new client with user agent: %s", client.Options.UserAgent)
	// Output:created new client with user agent: Custom UserAgent v1.0
}

// BenchmarkDefaultClientOptions benchmarks the method DefaultClientOptions()
func BenchmarkDefaultClientOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = DefaultClientOptions()
	}
}

// TestDefaultClientOptions_NoRetry will set 0 retry counts
func TestDefaultClientOptions_NoRetry(t *testing.T) {
	t.Parallel()

	options := DefaultClientOptions()
	options.RequestRetryCount = 0
	client, err := NewClient(options, nil)

	if client == nil {
		t.Fatal("failed to load client")
	} else if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
}

// TestClient_AddMiner tests the method AddMiner()
func TestClient_AddMiner(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		inputMiner    Miner
		expectedName  string
		expectedURL   string
		expectedNil   bool
		expectedError bool
	}{
		{Miner{
			MinerID: testMinerID,
			Name:    "Test",
			Token:   testMinerToken,
			URL:     testMinerURL,
		}, "Test", "testminer.com", false, false},
		{Miner{
			MinerID: testMinerID,
			Name:    "Test",
			Token:   testMinerToken,
			URL:     testMinerURL,
		}, "Test", "testminer.com", false, true},
		{Miner{
			MinerID: testMinerID,
			Name:    "Test2",
			Token:   testMinerToken,
			URL:     testMinerURL,
		}, "Test", "testminer.com", true, true},
		{Miner{
			MinerID: testMinerID,
			Name:    "",
			Token:   testMinerToken,
			URL:     testMinerURL,
		}, "Test", "testminer.com", true, true},
		{Miner{
			MinerID: testMinerID,
			Name:    "Test2",
			Token:   testMinerToken,
			URL:     "",
		}, "Test", "testminer.com", true, true},
	}

	// Create a client
	client := newTestClient(&mockHTTPDefaultClient{})

	// Run tests
	for _, test := range tests {
		if err := client.AddMiner(test.inputMiner); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.inputMiner, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error was expected", t.Name(), test.inputMiner)
		}
		// Get the miner
		miner := client.MinerByName(test.inputMiner.Name)
		if miner == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was not expected", t.Name(), test.inputMiner)
		} else if miner != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was expected", t.Name(), test.inputMiner)
		} else if miner != nil && miner.Name != test.expectedName {
			t.Errorf("%s Failed: [%v] inputted and [%s] expected but got: %s", t.Name(), test.inputMiner, test.expectedName, miner.Name)
		} else if miner != nil && miner.URL != test.expectedURL {
			t.Errorf("%s Failed: [%v] inputted and [%s] expected but got: %s", t.Name(), test.inputMiner, test.expectedURL, miner.URL)
		}
	}
}

// ExampleClient_AddMiner example using AddMiner()
func ExampleClient_AddMiner() {
	client, err := NewClient(nil, nil)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Add a miner
	if err = client.AddMiner(Miner{Name: testMinerName, URL: testMinerURL}); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get miner by name
	fmt.Printf("created new miner named: %s", client.MinerByName(testMinerName).Name)
	// Output:created new miner named: TestMiner
}

// BenchmarkClient_AddMiner benchmarks the method AddMiner()
func BenchmarkClient_AddMiner(b *testing.B) {
	client, _ := NewClient(nil, nil)
	for i := 0; i < b.N; i++ {
		_ = client.AddMiner(Miner{Name: testMinerName, URL: testMinerURL})
	}
}

// TestClient_MinerByName tests the method MinerByName()
func TestClient_MinerByName(t *testing.T) {
	t.Parallel()

	client := newTestClient(&mockHTTPDefaultClient{})

	// Add a valid miner
	err := client.AddMiner(Miner{
		Name: testMinerName,
		URL:  testMinerURL,
	})
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Get valid miner
	miner := client.MinerByName(testMinerName)
	if miner == nil {
		t.Fatalf("expected miner to not be nil using: %s", testMinerName)
	}

	// Get invalid miner
	miner = client.MinerByName("Unknown")
	if miner != nil {
		t.Fatalf("expected miner to be nil but got: %v", miner)
	}

}

// ExampleClient_MinerByName example using MinerByName()
func ExampleClient_MinerByName() {
	client, err := NewClient(nil, nil)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Add a miner
	if err = client.AddMiner(Miner{Name: testMinerName, URL: testMinerURL}); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get miner by name
	fmt.Printf("created new miner named: %s", client.MinerByName(testMinerName).Name)
	// Output:created new miner named: TestMiner
}

// BenchmarkClient_MinerByName benchmarks the method MinerByName()
func BenchmarkClient_MinerByName(b *testing.B) {
	client, _ := NewClient(nil, nil)
	_ = client.AddMiner(Miner{Name: testMinerName, URL: testMinerURL})
	for i := 0; i < b.N; i++ {
		_ = client.MinerByName(testMinerName)
	}
}

// TestClient_MinerByID tests the method MinerByID()
func TestClient_MinerByID(t *testing.T) {
	t.Parallel()

	client := newTestClient(&mockHTTPDefaultClient{})

	// Add a valid miner
	err := client.AddMiner(Miner{
		Name:    testMinerName,
		MinerID: testMinerID,
		URL:     testMinerURL,
	})
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Get valid miner
	miner := client.MinerByID(testMinerID)
	if miner == nil {
		t.Fatalf("expected miner to not be nil using: %s", testMinerID)
	}

	// Get invalid miner
	miner = client.MinerByID("00000")
	if miner != nil {
		t.Fatalf("expected miner to be nil but got: %v", miner)
	}
}

// ExampleClient_MinerByID example using MinerByID()
func ExampleClient_MinerByID() {
	client, err := NewClient(nil, nil)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Add a miner
	if err = client.AddMiner(Miner{Name: testMinerName, MinerID: testMinerID, URL: testMinerURL}); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get miner by id
	fmt.Printf("created new miner named: %s", client.MinerByID(testMinerID).Name)
	// Output:created new miner named: TestMiner
}

// BenchmarkClient_MinerByID benchmarks the method MinerByID()
func BenchmarkClient_MinerByID(b *testing.B) {
	client, _ := NewClient(nil, nil)
	_ = client.AddMiner(Miner{Name: testMinerName, MinerID: testMinerID, URL: testMinerURL})
	for i := 0; i < b.N; i++ {
		_ = client.MinerByID(testMinerID)
	}
}

// TestClient_MinerUpdateToken tests the method MinerUpdateToken()
func TestClient_MinerUpdateToken(t *testing.T) {
	t.Parallel()

	client := newTestClient(&mockHTTPDefaultClient{})

	// Add a valid miner
	err := client.AddMiner(Miner{
		Name:    testMinerName,
		MinerID: testMinerID,
		Token:   testMinerToken,
		URL:     testMinerURL,
	})
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	// Update a valid miner token
	client.MinerUpdateToken(testMinerName, "99999")

	// Get valid miner
	miner := client.MinerByID(testMinerID)
	if miner == nil {
		t.Fatalf("expected miner to not be nil using: %s", testMinerID)
	} else if miner.Token != "99999" {
		t.Fatalf("failed to update token to %s got: %s", "99999", miner.Token)
	}

	// Update a invalid miner token
	client.MinerUpdateToken("Unknown", "99999")
}

// ExampleClient_MinerUpdateToken example using MinerUpdateToken()
func ExampleClient_MinerUpdateToken() {
	client, err := NewClient(nil, nil)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Update existing miner token
	client.MinerUpdateToken(MinerTaal, "9999")

	// Get miner by id
	fmt.Printf("miner token found: %s", client.MinerByName(MinerTaal).Token)
	// Output:miner token found: 9999
}

// BenchmarkClient_MinerUpdateToken benchmarks the method MinerUpdateToken()
func BenchmarkClient_MinerUpdateToken(b *testing.B) {
	client, _ := NewClient(nil, nil)
	for i := 0; i < b.N; i++ {
		_ = client.MinerByName(MinerTaal)
	}
}
