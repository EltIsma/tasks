package common

type ErrorResponse struct {
	Error     string `json:"error" example:"error meessage"`
	ErrorCode int    `json:"error_code" example:"0"`
}

func NewErrorResponse(message string, code int) *ErrorResponse {
	return &ErrorResponse{
		Error:     message,
		ErrorCode: code,
	}
}