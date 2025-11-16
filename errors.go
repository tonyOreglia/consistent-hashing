package consistent_hash

import "encoding/json"

// ApiError represents the body of an error response.
type ApiError struct {
	// Desc is a human readable description of the error that was encountered.
	Desc string

	// Err is the raw error that was encountered.
	Err error
}

// MarshalJSON helper to convert the ApiError into JSON friendly response.
func (e ApiError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Description string `json:"description"`
		Error       string `json:"error"`
	}{
		Description: e.Desc,
		Error:       e.Err.Error(),
	})
}
