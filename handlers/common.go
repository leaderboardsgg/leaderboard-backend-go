package handlers

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Errors []error `json:"error"`
}
