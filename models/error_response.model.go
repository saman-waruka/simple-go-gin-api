package models

// ErrorResponse represents a standard error response
type ErrorUnauthorizedResponse struct {
    Error string `json:"error" example:"Unauthorized"`
}

type ErrorInternalServerResponse struct {
	Error string `json:"error" example:"Internal Server Error"`
}

type ErrorBadRequestCreateUserResponse struct {
	Error string `json:"error" example:"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'required' tag"`
}