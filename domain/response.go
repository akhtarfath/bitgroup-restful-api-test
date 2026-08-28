package domain

type ResponseWithData struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ResponseMessage struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewSuccessResponse(message string, data any) ResponseWithData {
	return ResponseWithData{
		Code:    200,
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func NewCreatedResponse(message string, data any) ResponseWithData {
	return ResponseWithData{
		Code:    201,
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func NewSuccessMessage(message string) ResponseMessage {
	return ResponseMessage{
		Code:    200,
		Status:  "success",
		Message: message,
	}
}

func NewErrorMessage(code int, message string) ResponseMessage {
	return ResponseMessage{
		Code:    code,
		Status:  "error",
		Message: message,
	}
}
