package minercraft

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ErrorResponse_AsError(t *testing.T) {
	// Ensure error maps correctly to the type when converted to a normal type
	err := ErrorResponse{
		Type:    "Type123",
		Title:   "Title123",
		Status:  http.StatusConflict,
		Detail:  "Detail123",
		TraceID: "TraceID123",
		Errors:  nil,
	}
	var testErr error = err
	var errResp ErrorResponse
	// Use assert.ErrorAs to verify the type of the error
	require.ErrorAs(t, testErr, &errResp)

	assert.EqualValues(t, err, errResp)
}

func Test_ErrRetryable(t *testing.T) {
	// Ensure ErrRetryable can be checked for Retryable
	err := ErrRetryable{err: ErrorResponse{}}
	var testErr error = err

	// test using our helper method
	assert.True(t, IsRetryable(testErr))
	require.Error(t, err)

	// test using the As method for interface Retryable
	var r Retryable
	require.ErrorAs(t, testErr, &r)

	// test using the As method for ErrorResponse
	var e ErrorResponse
	require.ErrorAs(t, testErr, &e)

	err = ErrRetryable{err: fmt.Errorf("new err")}
	assert.False(t, errors.As(err, &e))
}

func Test_ErrorResponse_Error(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		// TODO: add test properties
		err ErrorResponse
		exp string
	}{
		"error response should print string when error called": {
			err: ErrorResponse{
				Type:    "Type123",
				Title:   "Title123",
				Status:  http.StatusConflict,
				Detail:  "Detail123",
				TraceID: "TraceID123",
				Errors:  nil,
			},
			exp: "title: Title123 \n detail: Detail123 \n traceID: TraceID123 \n validation errors: ",
		}, "error response should print string including validation errors when they are present": {
			err: ErrorResponse{
				Type:    "Type123",
				Title:   "Title123",
				Status:  http.StatusConflict,
				Detail:  "Detail123",
				TraceID: "TraceID123",
				Errors: map[string][]string{
					"field1": {
						"failed1", "failed2",
					},
				},
			},
			exp: "title: Title123 \n detail: Detail123 \n traceID: TraceID123 \n validation errors: [field1: failed1, failed2]",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.exp, test.err.Error())
		})
	}
}
