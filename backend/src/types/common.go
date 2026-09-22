package types

import "fmt"

// AppError service/controller 分层包装的业务异常，携带错误码与可展示消息。
type AppError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func (e *AppError) Error() string { return fmt.Sprintf("%s: %s", e.Code, e.Message) }

func NewAppError(code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// WithDetails 附加放行阻塞清单等结构化明细。
func (e *AppError) WithDetails(details map[string]any) *AppError {
	e.Details = details
	return e
}

// APIResponse 统一响应外层结构。
type APIResponse struct {
	OK    bool        `json:"ok"`
	Data  interface{} `json:"data,omitempty"`
	Error *AppError   `json:"error,omitempty"`
}
