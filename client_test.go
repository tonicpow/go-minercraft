package minercraft

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAPIVersion = "1.4.0" // Old version was 0.1.0
	testEncoding   = "UTF-8"
	testMimeType   = "application/json"
	testMinerID    = "1234567"
	testMinerName  = "TestMiner"
	testMinerToken = "0987654321"
	testMinerURL   = "https://testminer.com"
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
func newTestClient(httpClient HTTPInterface) ClientInterface {
	client, _ := NewClient(nil, httpClient, "", nil, nil)
	// client.httpClient = httpClient
	return client
}

// TestNewClient tests the method NewClient()
func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("valid new client", func(t *testing.T) {
		client, err := NewClient(nil, nil, "", nil, nil)
		assert.NotNil(t, client)
		assert.NoError(t, err)

		// Test default miners
		assert.Equal(t, 4, len(client.Miners()))
	})

	t.Run("custom http client", func(t *testing.T) {
		client, err := NewClient(nil, http.DefaultClient, "", nil, nil)
		assert.NotNil(t, client)
		assert.NoError(t, err)
	})

	t.Run("default miners", func(t *testing.T) {
		client, err := NewClient(nil, nil, "", nil, nil)
		assert.NotNil(t, client)
		assert.NoError(t, err)

		// Get Taal
		miner := client.MinerByName(MinerTaal)
		assert.Equal(t, MinerTaal, miner.Name)

		// Get Mempool
		miner = client.MinerByName(MinerMempool)
		assert.Equal(t, MinerMempool, miner.Name)

		// Get Matterpool
		miner = client.MinerByName(MinerMatterpool)
		assert.Equal(t, MinerMatterpool, miner.Name)

		// Get GorillaPool
		miner = client.MinerByName(MinerGorillaPool)
		assert.Equal(t, MinerGorillaPool, miner.Name)
	})

	t.Run("custom miners", func(t *testing.T) {
		miners := []*Miner{{
			MinerID: testMinerID,
			Name:    testMinerName,
		}}

		minerAPIs := []*MinerAPIs{{MinerID: testMinerID, APIs: []API{{Token: testMinerToken, URL: testMinerURL, Type: MAPI}}}}

		client, err := NewClient(nil, nil, "", miners, minerAPIs)
		assert.NotNil(t, client)
		assert.NoError(t, err)

		// Get test miner
		miner := client.MinerByName(testMinerName)
		assert.Equal(t, testMinerName, miner.Name)

		assert.Equal(t, 1, len(client.Miners()))
	})
}

// ExampleNewClient example using NewClient()
func ExampleNewClient() {
	client, err := NewClient(nil, nil, "", nil, nil)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("created new client with %d default miners", len(client.Miners()))
	// Output:created new client with 4 default miners
}

// BenchmarkNewClient benchmarks the method NewClient()
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewClient(nil, nil, "", nil, nil)
	}
}

// TestDefaultClientOptions tests setting DefaultClientOptions()
func TestDefaultClientOptions(t *testing.T) {
	t.Parallel()

	t.Run("default client options", func(t *testing.T) {
		options := DefaultClientOptions()

		assert.Equal(t, defaultUserAgent, options.UserAgent)
		assert.Equal(t, 2.0, options.BackOffExponentFactor)
		assert.Equal(t, 2*time.Millisecond, options.BackOffInitialTimeout)
		assert.Equal(t, 2*time.Millisecond, options.BackOffMaximumJitterInterval)
		assert.Equal(t, 10*time.Millisecond, options.BackOffMaxTimeout)
		assert.Equal(t, 20*time.Second, options.DialerKeepAlive)
		assert.Equal(t, 5*time.Second, options.DialerTimeout)
		assert.Equal(t, 2, options.RequestRetryCount)
		assert.Equal(t, 30*time.Second, options.RequestTimeout)
		assert.Equal(t, 3*time.Second, options.TransportExpectContinueTimeout)
		assert.Equal(t, 20*time.Second, options.TransportIdleTimeout)
		assert.Equal(t, 10, options.TransportMaxIdleConnections)
		assert.Equal(t, 5*time.Second, options.TransportTLSHandshakeTimeout)
	})

	t.Run("no retry", func(t *testing.T) {
		options := DefaultClientOptions()
		options.RequestRetryCount = 0
		client, err := NewClient(options, nil, "", nil, nil)
		assert.NotNil(t, client)
		assert.NoError(t, err)
	})
}

// ExampleDefaultClientOptions example using DefaultClientOptions()
func ExampleDefaultClientOptions() {
	options := DefaultClientOptions()
	options.UserAgent = "Custom UserAgent v1.0"
	client, err := NewClient(options, nil, "", nil, nil)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("created new client with user agent: %s", client.UserAgent())
	// Output:created new client with user agent: Custom UserAgent v1.0
}

// BenchmarkDefaultClientOptions benchmarks the method DefaultClientOptions()
func BenchmarkDefaultClientOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = DefaultClientOptions()
	}
}

// TestClient_AddMiner tests the method AddMiner()
func TestClient_AddMiner(t *testing.T) {
	t.Parallel()

	t.Run("valid cases", func(t *testing.T) {

		// Create the list of tests
		var tests = []struct {
			testCase     string
			inputMiner   Miner
			inputAPIs    []API
			expectedName string
			expectedURL  string
		}{
			{
				"valid miner",
				Miner{
					MinerID: testMinerID,
					Name:    "Test",
				},
				[]API{
					API{
						Token: testMinerToken,
						URL:   "https://testminer.com",
						Type:  MAPI,
					},
				},
				"Test",
				"https://testminer.com",
			},
		}

		// Run tests
		client := newTestClient(&mockHTTPDefaultClient{})
		for _, test := range tests {
			t.Run(test.testCase, func(t *testing.T) {
				err := client.AddMiner(test.inputMiner, test.inputAPIs)
				assert.NoError(t, err)

				// Get the miner
				miner := client.MinerByName(test.inputMiner.Name)
				assert.Equal(t, test.expectedName, miner.Name)
				assert.Equal(t, test.expectedURL, test.inputAPIs[0].URL)
			})
		}
	})

	t.Run("invalid cases", func(t *testing.T) {

		// Create the list of tests
		var tests = []struct {
			testCase   string
			inputMiner Miner
			inputAPIs  []API
		}{
			{
				"duplicate miner - by name",
				Miner{
					MinerID: testMinerID + "123",
					Name:    "Test",
				},
				[]API{
					API{
						Token: testMinerToken,
						URL:   testMinerURL,
						Type:  MAPI,
					},
				},
			},
			{
				"duplicate miner - by id",
				Miner{
					MinerID: testMinerID,
					Name:    "Test123",
				},
				[]API{
					API{
						Token: testMinerToken,
						URL:   testMinerURL,
						Type:  MAPI,
					},
				},
			},
			{
				"missing miner name",
				Miner{
					MinerID: testMinerID,
					Name:    "",
				},
				[]API{
					API{
						Token: testMinerToken,
						URL:   testMinerURL,
						Type:  MAPI,
					},
				},
			},
			{
				"missing miner url",
				Miner{
					MinerID: testMinerID,
					Name:    "TestURL",
				},
				[]API{
					API{
						Token: testMinerToken,
						URL:   testMinerURL,
						Type:  MAPI,
					},
				},
			},
			{
				"invalid miner url - http",
				Miner{
					MinerID: testMinerID,
					Name:    "TestURL",
				},
				[]API{
					API{
						Token: testMinerToken,
						URL:   "www.domain.com",
						Type:  MAPI,
					},
				},
			},
			{
				"invalid miner url - trigger parse error",
				Miner{
					MinerID: testMinerID,
					Name:    "TestURL",
				},
				[]API{
					API{
						Token: testMinerToken,
						URL:   "postgres://user:abc{DEf1=ghi@example.com:5432/db?sslmode=require",
						Type:  MAPI,
					},
				},
			},
		}

		// Run tests
		client := newTestClient(&mockHTTPDefaultClient{})

		// Add a miner to start
		err := client.AddMiner(Miner{MinerID: testMinerID, Name: "Test"}, []API{API{URL: testMinerURL, Type: Arc}})
		assert.NoError(t, err)

		for _, test := range tests {
			t.Run(test.testCase, func(t *testing.T) {
				err = client.AddMiner(test.inputMiner, test.inputAPIs)
				assert.Error(t, err)
			})
		}
	})
}

// ExampleClient_AddMiner example using AddMiner()
func ExampleClient_AddMiner() {
	client, err := NewClient(nil, nil, "", nil, nil)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Add a miner
	if err = client.AddMiner(Miner{Name: testMinerName}, []API{API{URL: testMinerURL, Type: Arc}}); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get miner by name
	fmt.Printf("created new miner named: %s", client.MinerByName(testMinerName).Name)
	// Output:created new miner named: TestMiner
}

// BenchmarkClient_AddMiner benchmarks the method AddMiner()
func BenchmarkClient_AddMiner(b *testing.B) {
	client, _ := NewClient(nil, nil, "", nil, nil)
	for i := 0; i < b.N; i++ {
		_ = client.AddMiner(Miner{Name: testMinerName}, []API{API{URL: testMinerURL, Type: Arc}})
	}
}

// TestClient_MinerByName tests the method MinerByName()
func TestClient_MinerByName(t *testing.T) {
	t.Parallel()

	t.Run("get valid miner", func(t *testing.T) {
		client := newTestClient(&mockHTTPDefaultClient{})

		// Add a valid miner
		err := client.AddMiner(Miner{
			Name: testMinerName,
		}, []API{API{URL: testMinerURL, Type: Arc}})
		assert.NoError(t, err)

		// Get valid miner
		miner := client.MinerByName(testMinerName)
		assert.NotNil(t, miner)
	})

	t.Run("get invalid miner", func(t *testing.T) {
		client := newTestClient(&mockHTTPDefaultClient{})

		// Add a valid miner
		err := client.AddMiner(Miner{
			Name: testMinerName,
		}, []API{API{URL: testMinerURL, Type: Arc}},
		)
		assert.NoError(t, err)

		// Get invalid miner
		miner := client.MinerByName("Unknown")
		assert.Nil(t, miner)
	})
}

// ExampleClient_MinerByName example using MinerByName()
func ExampleClient_MinerByName() {
	client, err := NewClient(nil, nil, "", nil, nil)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Add a miner
	if err = client.AddMiner(Miner{Name: testMinerName}, []API{API{URL: testMinerURL, Type: Arc}}); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get miner by name
	fmt.Printf("created new miner named: %s", client.MinerByName(testMinerName).Name)
	// Output:created new miner named: TestMiner
}

// BenchmarkClient_MinerByName benchmarks the method MinerByName()
func BenchmarkClient_MinerByName(b *testing.B) {
	client, _ := NewClient(nil, nil, "", nil, nil)
	_ = client.AddMiner(Miner{Name: testMinerName}, []API{API{URL: testMinerURL, Type: MAPI}})
	for i := 0; i < b.N; i++ {
		_ = client.MinerByName(testMinerName)
	}
}

// TestClient_MinerByID tests the method MinerByID()
func TestClient_MinerByID(t *testing.T) {
	t.Parallel()

	t.Run("get valid miner", func(t *testing.T) {
		client := newTestClient(&mockHTTPDefaultClient{})

		// Add a valid miner
		err := client.AddMiner(Miner{
			Name:    testMinerName,
			MinerID: testMinerID,
		}, []API{API{URL: testMinerURL, Type: MAPI}})
		assert.NoError(t, err)

		// Get valid miner
		miner := client.MinerByID(testMinerID)
		assert.NotNil(t, miner)
	})

	t.Run("get invalid miner", func(t *testing.T) {
		client := newTestClient(&mockHTTPDefaultClient{})

		// Add a valid miner
		err := client.AddMiner(Miner{
			Name:    testMinerName,
			MinerID: testMinerID,
		}, []API{API{URL: testMinerURL, Type: MAPI}})
		assert.NoError(t, err)

		// Get invalid miner
		miner := client.MinerByID("00000")
		assert.Nil(t, miner)
	})
}

// ExampleClient_MinerByID example using MinerByID()
func ExampleClient_MinerByID() {
	client, err := NewClient(nil, nil, "", nil, nil)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Add a miner
	if err = client.AddMiner(Miner{Name: testMinerName, MinerID: testMinerID}, []API{API{URL: testMinerURL, Type: MAPI}}); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get miner by id
	fmt.Printf("created new miner named: %s", client.MinerByID(testMinerID).Name)
	// Output:created new miner named: TestMiner
}

// BenchmarkClient_MinerByID benchmarks the method MinerByID()
func BenchmarkClient_MinerByID(b *testing.B) {
	client, _ := NewClient(nil, nil, "", nil, nil)
	_ = client.AddMiner(Miner{Name: testMinerName, MinerID: testMinerID}, []API{API{URL: testMinerURL, Type: MAPI}})
	for i := 0; i < b.N; i++ {
		_ = client.MinerByID(testMinerID)
	}
}

// TestClient_MinerUpdateToken tests the method MinerUpdateToken()
func TestClient_MinerUpdateToken(t *testing.T) {
	t.Parallel()

	t.Run("update valid miner", func(t *testing.T) {
		client := newTestClient(&mockHTTPDefaultClient{})

		// Add a valid miner
		err := client.AddMiner(Miner{
			Name:    testMinerName,
			MinerID: testMinerID,
		}, []API{{URL: testMinerURL, Type: MAPI, Token: testMinerToken}})
		assert.NoError(t, err)

		// Update a valid miner token
		client.MinerUpdateToken(testMinerName, "99999", MAPI)

		// Get valid miner
		miner := client.MinerByID(testMinerID)
		api, err := client.MinerAPIByMinerID(testMinerID, MAPI)
		assert.NoError(t, err)
		assert.NotNil(t, miner)
		assert.Equal(t, "99999", api.Token)
	})

	t.Run("update unknown miner", func(t *testing.T) {
		client := newTestClient(&mockHTTPDefaultClient{})

		// Add a valid miner
		err := client.AddMiner(Miner{
			Name:    testMinerName,
			MinerID: testMinerID,
		}, []API{API{URL: testMinerURL, Type: MAPI, Token: testMinerToken}})
		assert.NoError(t, err)

		// Update a invalid miner token
		client.MinerUpdateToken("Unknown", "99999", MAPI)
	})
}

// ExampleClient_MinerUpdateToken example using MinerUpdateToken()
func ExampleClient_MinerUpdateToken() {
	client, err := NewClient(nil, nil, "", nil, nil)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Update existing miner token
	client.MinerUpdateToken(MinerTaal, "9999", MAPI)

	taal := client.MinerByName(MinerTaal)
	api, err := client.MinerAPIByMinerID(taal.MinerID, client.APIType())
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Get miner by id
	fmt.Printf("miner token found: %s", api.Token)
	// Output:miner token found: 9999
}

// BenchmarkClient_MinerUpdateToken benchmarks the method MinerUpdateToken()
func BenchmarkClient_MinerUpdateToken(b *testing.B) {
	client, _ := NewClient(nil, nil, "", nil, nil)
	for i := 0; i < b.N; i++ {
		_ = client.MinerByName(MinerTaal)
	}
}

// TestClient_RemoveMiner will remove a miner by name or ID
func TestClient_RemoveMiner(t *testing.T) {
	t.Parallel()

	t.Run("remove a valid miner", func(t *testing.T) {
		client := newTestClient(&mockHTTPDefaultClient{})

		// Remove miner
		removed := client.RemoveMiner(client.MinerByName(MinerTaal))
		assert.Equal(t, true, removed)

		// Try to get the miner
		miner := client.MinerByName(MinerTaal)
		assert.Nil(t, miner)
	})

	t.Run("remove an invalid miner", func(t *testing.T) {
		client := newTestClient(&mockHTTPDefaultClient{})

		// Unknown miner
		dummyMiner := &Miner{
			MinerID: "dummy",
			Name:    "dummy",
		}

		// Remove miner
		removed := client.RemoveMiner(dummyMiner)
		assert.Equal(t, false, removed)
	})

	t.Run("remove a nil miner", func(t *testing.T) {
		client := newTestClient(&mockHTTPDefaultClient{})

		// Remove miner
		assert.Panics(t, func() {
			removed := client.RemoveMiner(nil)
			assert.Equal(t, false, removed)
		})
	})
}

// ExampleClient_MinerUpdateToken example using RemoveMiner()
func ExampleClient_RemoveMiner() {
	client, err := NewClient(nil, nil, "", nil, nil)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	// Update existing miner token
	client.RemoveMiner(client.MinerByName(MinerTaal))

	// Show response
	fmt.Printf("total miners: %d", len(client.Miners()))
	// Output:total miners: 3
}

// BenchmarkClient_RemoveMiner benchmarks the method RemoveMiner()
func BenchmarkClient_RemoveMiner(b *testing.B) {
	client, _ := NewClient(nil, nil, "", nil, nil)
	miner := client.MinerByName(MinerTaal)
	for i := 0; i < b.N; i++ {
		_ = client.RemoveMiner(miner)
	}
}

// TestDefaultMiners will test the method DefaultMiners()
func TestDefaultMiners(t *testing.T) {
	t.Run("default json", func(t *testing.T) {
		miners, err := DefaultMiners()
		require.NoError(t, err)
		require.NotNil(t, miners)
		assert.Equal(t, 4, len(miners))
		assert.Equal(t, MinerTaal, miners[0].Name)
		assert.Equal(t, MinerMempool, miners[1].Name)
		assert.Equal(t, MinerGorillaPool, miners[3].Name)
		assert.Equal(t, MinerMatterpool, miners[2].Name)
	})
}

// ExampleDefaultMiners example using DefaultMiners()
func ExampleDefaultMiners() {
	miners, _ := DefaultMiners()
	fmt.Printf("total miners: %d", len(miners))
	// Output:total miners: 4
}

// BenchmarkDefaultMiners benchmarks the method DefaultMiners()
func BenchmarkDefaultMiners(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DefaultMiners()
	}
}
