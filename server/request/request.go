package request

import (
	"encoding/json"
	"errors"
)

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

var ErrInvalidSuccessResponse = errors.New("expected response to be a valid SuccessResponse")
var ErrDataInvalidJson = errors.New("the data key was not valid json")

// This function will unmarshal the response into a handlers.SuccessResponse,
// re-marshal the contents of `data`, and unmarshal it into the expected destination.
// This allows for a testing flow that can still assert by unmarshalling
// SuccessResponse.Data into concrete types.
func UnmarshalSuccessResponseData(
	successRepsonseBytes []byte,
	dest interface{},
) ([]byte, error) {
	var response SuccessResponse
	err := json.Unmarshal(successRepsonseBytes, &response)
	if err != nil {
		return nil, ErrInvalidSuccessResponse
	}
	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		return dataBytes, ErrDataInvalidJson
	}
	err = json.Unmarshal(dataBytes, dest)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type ErrorResponse struct {
	Errors []error `json:"errors"`
}
