package core

import "encoding/json"

type APIResponse struct {
	StatusCode int         `json:"-"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
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
			Error:      err,
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
			Error:      err,
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
	return 400
}