package dto

// JSONRequest JSON处理请求
type JSONRequest struct {
	Text   string `json:"text" binding:"required"`
	Indent int    `json:"indent,omitempty"`
}

// JSONResponse JSON处理响应
type JSONResponse struct {
	Result  string `json:"result,omitempty"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// ValidateResponse JSON验证响应
type ValidateResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
