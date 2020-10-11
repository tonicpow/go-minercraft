package minercraft

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// RequestResponse is the response from a request
type RequestResponse struct {
	BodyContents []byte `json:"body_contents"` // Raw body response
	Error        error  `json:"error"`         // If an error occurs
	Method       string `json:"method"`        // Method is the HTTP method used
	PostData     string `json:"post_data"`     // PostData is the post data submitted if POST/PUT request
	StatusCode   int    `json:"status_code"`   // StatusCode is the last code from the request
	URL          string `json:"url"`           // URL is used for the request
}

// httpRequest is a generic request wrapper that can be used without constraints
func httpRequest(client *Client, method, url, token string, payload []byte) (response *RequestResponse) {

	// Set reader
	var bodyReader io.Reader

	// Start the response
	response = new(RequestResponse)

	// Add post data if applicable
	if method == http.MethodPost || method == http.MethodPut {
		bodyReader = bytes.NewBuffer(payload)
		response.PostData = string(payload)
	}

	// Store for debugging purposes
	response.Method = method
	response.URL = url

	// Start the request
	var request *http.Request
	if request, response.Error = http.NewRequestWithContext(context.Background(), method, url, bodyReader); response.Error != nil {
		return
	}

	// Change the header (user agent is in case they block default Go user agents)
	request.Header.Set("User-Agent", client.Options.UserAgent)

	// Set the content type on Method
	if method == http.MethodPost || method == http.MethodPut {
		request.Header.Set("Content-Type", "application/json")
	}

	// Set a token if supplied
	if len(token) > 0 {
		request.Header.Set("token", token)
	}

	// Fire the http request
	var resp *http.Response
	if resp, response.Error = client.httpClient.Do(request); response.Error != nil {
		if resp != nil {
			response.StatusCode = resp.StatusCode
		}
		return
	}

	// Close the response body
	defer func() {
		_ = resp.Body.Close()
	}()

	// Set the status
	response.StatusCode = resp.StatusCode

	// Check status code
	if http.StatusOK != resp.StatusCode {
		response.Error = fmt.Errorf("status code: %d does not match %d", resp.StatusCode, http.StatusOK)
		return
	}

	// Read the body
	response.BodyContents, response.Error = ioutil.ReadAll(resp.Body)

	return
}