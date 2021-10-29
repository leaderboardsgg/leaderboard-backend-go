package handlers

type ErrorResponse struct {
	Errors []error `json:"error"`
}
