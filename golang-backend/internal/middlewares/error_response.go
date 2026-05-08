package middlewares

type errorResponse struct {
	ErrorDetails string `json:"error"`
	UserMessage  string `json:"message"`
	APIError     bool   `json:"APIError"`
}
