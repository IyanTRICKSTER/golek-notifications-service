package response

import (
	status "golek_notifications_service/pkg/operation_status"
)

type Response struct {
	status     bool
	statusCode status.OperationStatus
	message    string
	data       any
}

func (r Response) GetStatus() bool {
	return r.status
}

func (r Response) GetStatusCode() status.OperationStatus {
	return r.statusCode
}

func (r Response) GetMessage() string {
	return r.message
}

func (r Response) GetData() any {
	return r.data
}

func (r *Response) SetMessage(msg string) {
	r.message = msg
}

func (r Response) ErrorIs(code status.OperationStatus) bool {
	if r.statusCode == code {
		return true
	}
	return false
}

func (r Response) IsFailed() bool {
	if r.status == true {
		return false
	}
	return true
}

func (r Response) ToMapStringInterface() map[string]interface{} {
	return map[string]interface{}{
		"status":     r.status,
		"statusCode": r.statusCode,
		"message":    r.message,
		"data":       r.data,
	}
}

func New(status bool, statusCode status.OperationStatus, message string, data any) Response {
	return Response{status: status, statusCode: statusCode, message: message, data: data}
}
