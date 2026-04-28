package utils

type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	ErrorCode int         `json:"error_code,omitempty"`
}

func Success(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func Error(message string, code int) APIResponse {
	return APIResponse{
		Success:   false,
		Message:   message,
		ErrorCode: code,
		Data:      map[string]interface{}{},
	}
}
