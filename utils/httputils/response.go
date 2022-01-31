package httputils

// ValidResponse struct can be used to formulate a standard response
type ValidResponse struct {
	Message string `json:"message"`
}

// NewResponse is a constructor
func NewResponse(message string) ValidResponse {
	return ValidResponse{Message: message}
}
