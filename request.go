package minercraft

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Retryable can be implemented to identify a struct as retryable, in this case an error can be deemed retryable.
type Retryable interface {
	// IsRetryable if a method has this method then it's a retryable error.
	IsRetryable()
}

// IsRetryable can be passed an error to check if it is retryable
func IsRetryable(err error) bool {
	var e Retryable
	return errors.As(err, &e)
}

// ErrRetryable indicates a retryable error.
//
// To check an error is a retryable error do:
//
//	errors.Is(err, minercraft.Retryable)
type ErrRetryable struct{ err error }

func (e ErrRetryable) Error() string {
	return e.err.Error()
}

// IsRetryable returns true denoting this is retryable.
func (e ErrRetryable) IsRetryable() {}

// Is allows the underlying error to be checked that it is a certain error type.
func (e ErrRetryable) Is(err error) bool { return errors.Is(e.err, err) }

// As will return true if the error can be cast to the target.
func (e ErrRetryable) As(target any) bool { return errors.As(e.err, target) }

// ErrorResponse is the response returned from mAPI on error.
type ErrorResponse struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Status  int    `json:"status"`
	Detail  string `json:"detail"`
	TraceID string `json:"traceId"`
	// Errors will return a list of formatting errors in the case of a bad request
	// being sent to mAPI.
	Errors map[string][]string `json:"errors"`
}

// Error defines the ErrorResponse as an error, an error can be converted
// to it using the below:
//
//	 var errResp ErrorResponse
//	 if errors.As(testErr, &errResp) {
//		 // handle error
//		 fmt.Println(errResp.Title)
//	 }
func (e ErrorResponse) Error() string {
	sb := strings.Builder{}
	for field, warnings := range e.Errors {
		sb.WriteString("[" + field + ": ")
		sb.WriteString(strings.Join(warnings, ", "))
		sb.WriteString("]")
	}
	return fmt.Sprintf("title: %s \n detail: %s \n traceID: %s \n validation errors: %s",
		e.Title, e.Detail, e.TraceID, sb.String())
}

// RequestResponse is the response from a request
type RequestResponse struct {
	BodyContents []byte `json:"body_contents"` // Raw body response
	Error        error  `json:"error"`         // If an error occurs
	Method       string `json:"method"`        // Method is the HTTP method used
	PostData     string `json:"post_data"`     // PostData is the post data submitted if POST/PUT request
	StatusCode   int    `json:"status_code"`   // StatusCode is the last code from the request
	URL          string `json:"url"`           // URL is used for the request
}

// httpPayload is used for a httpRequest
type httpPayload struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	Token  string `json:"token"`
	Data   []byte `json:"data"`
}

// httpRequest is a generic request wrapper that can be used without constraints.
//
// If response.Error isn't nil it can be checked for being retryable by calling errors.Is(err, minercraft.Retryable),
// this means the request returned an intermittent / transient error and can be retried depending on client
// requirements.
//
// It can also be converted to the ErrorResponse type to get the error detail as shown:
//
//		var errResp ErrorResponse
//		if errors.As(testErr, &errResp) {
//	    // handle error
//	    fmt.Println(errResp.Title)
//	 }
func httpRequest(ctx context.Context, client *Client,
	payload *httpPayload) (response *RequestResponse) {

	// Set reader
	var bodyReader io.Reader

	// Start the response
	response = new(RequestResponse)

	// Add post data if applicable
	if payload.Method == http.MethodPost || payload.Method == http.MethodPut {
		bodyReader = bytes.NewBuffer(payload.Data)
		response.PostData = string(payload.Data)
	}

	// Store for debugging purposes
	response.Method = payload.Method
	response.URL = payload.URL

	// Start the request
	var request *http.Request
	if request, response.Error = http.NewRequestWithContext(
		ctx, payload.Method, payload.URL, bodyReader,
	); response.Error != nil {
		return
	}

	// Change the header (user agent is in case they block default Go user agents)
	request.Header.Set("User-Agent", client.Options.UserAgent)

	// Set the content type on Method
	if payload.Method == http.MethodPost || payload.Method == http.MethodPut {
		request.Header.Set("Content-Type", "application/json")
	}

	// Set a token if supplied
	if len(payload.Token) > 0 {
		request.Header.Set("Authorization", payload.Token)
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

	if resp.Body != nil {
		// Read the body
		response.BodyContents, response.Error = io.ReadAll(resp.Body)
	}
	// Check status code
	if http.StatusOK == resp.StatusCode {
		return
	}

	// indicates that resubmitting this request could be successful when mAPI
	// is available again.
	retryable := response.StatusCode >= 500 && response.StatusCode <= 599
	// unexpected status, write an error.
	if response.BodyContents == nil {
		// There's no "body" present, so just echo status code.
		statusErr := fmt.Errorf("status code: %d does not match %d", resp.StatusCode, http.StatusOK)
		if !retryable {
			response.Error = statusErr
			return
		}
		response.Error = ErrRetryable{err: statusErr}
		return
	}
	// Have a "body" so map to an error type and add to the error message.
	var errBody ErrorResponse
	if err := json.Unmarshal(response.BodyContents, &errBody); err != nil {
		response.Error = fmt.Errorf("failed to unmarshal mapi error response: %w", err)
		return
	}
	if retryable {
		response.Error = ErrRetryable{err: errBody}
		return
	}
	response.Error = errBody
	return
}
