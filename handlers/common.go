package handlers

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Errors []string `json:"errors"`
}
