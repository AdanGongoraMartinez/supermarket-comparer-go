package core

import (
	"encoding/json"

	"supermarket-comparer-go/internal/errors"
)

type APIResponse struct {
	StatusCode int  `json:"-"`
	Success    bool `json:"success"`
	Data       any  `json:"data,omitempty"`
	Error      any  `json:"error,omitempty"`
}

func (r *APIResponse) ToJSON() []byte {
	data, _ := json.Marshal(r)
	return data
}

func HandleResult[T any](value T, err error, successStatus int) *APIResponse {
	if err != nil {
		return &APIResponse{
			StatusCode: getErrorStatus(err),
			Success:    false,
			Error:      err.Error(),
		}
	}
	return &APIResponse{
		StatusCode: successStatus,
		Success:    true,
		Data:       value,
	}
}

func HandleEmptyResult(err error, successStatus int) *APIResponse {
	if err != nil {
		return &APIResponse{
			StatusCode: getErrorStatus(err),
			Success:    false,
			Error:      err.Error(),
		}
	}
	return &APIResponse{
		StatusCode: successStatus,
		Success:    true,
	}
}

func getErrorStatus(err error) int {
	if err == nil {
		return 500
	}
	if sc, ok := err.(errors.StatusCoder); ok {
		return sc.GetStatusCode()
	}
	return 400
}

