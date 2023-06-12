package minercraft

import "fmt"

// APINotFoundError is returned when an API definition is not found for a miner
type APINotFoundError struct {
	MinerID string
	APIType APIType
}

// ActionRouteNotFoundError is returned when an action route is not found for a miner
type ActionRouteNotFoundError struct {
	ActionName APIActionName
	APIType    APIType
}

// Error returns the error message related to the APINotFoundError
func (e *APINotFoundError) Error() string {
	return fmt.Sprintf("API definition not found for MinerID: %s and APIType: %s", e.MinerID, e.APIType)
}

// Error returns the error message related to the ActionRouteNotFoundError
func (e *ActionRouteNotFoundError) Error() string {
	return fmt.Sprintf("Action route not found for ActionName: %s and APIType: %s", e.ActionName, e.APIType)
}
