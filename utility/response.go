package utility

type PaginationMeta struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

func BuildPaginationResponseSuccess(message string, data any, meta PaginationMeta) Response {
	res := Response{
		Status:  true,
		Message: message,
		Error:   nil,
		Data:    data,
		Meta:    meta,
	}
	return res
}

func BuildResponseSuccess(message string, data any) Response {
	res := Response{
		Status:  true,
		Message: message,
		Error:   nil,
		Data:    data,
	}
	return res
}

func BuildResponseFailed(message string, err any) Response {
	res := Response{
		Status:  false,
		Message: message,
		Error:   err,
		Data:    nil,
	}
	return res
}
