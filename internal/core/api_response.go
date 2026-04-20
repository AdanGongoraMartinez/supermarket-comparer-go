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

func HandleResult[T any](result *Result[T], successStatus int) *APIResponse {
	if result.IsSuccess() {
		return &APIResponse{
			StatusCode: successStatus,
			Success:   true,
			Data:      result.GetValue(),
		}
	}
	return &APIResponse{
		StatusCode: getErrorStatus(result.GetError()),
		Success:   false,
		Error:     result.GetError(),
	}
}

func HandleEmptyResult[T any](result *Result[T], successStatus int) *APIResponse {
	if result.IsSuccess() {
		return &APIResponse{
			StatusCode: successStatus,
			Success:   true,
		}
	}
	return &APIResponse{
		StatusCode: getErrorStatus(result.GetError()),
		Success:   false,
		Error:     result.GetError(),
	}
}

func getErrorStatus(err error) int {
	if err == nil {
		return 500
	}
	return 400
}