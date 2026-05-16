package badnet

import "fmt"

func SuccessResponse(code HTTPResponseCode, mimeType HTTPMimeType, data []byte) Response {
	return Response{
		Version:      V1_1,
		ResponseCode: code,
		ResponseMsg:  code.responseMessage(),
		Headers: HTTPHeaders{
			ContentLength: fmt.Sprintf("%d", len(data)),
			ContentType:   string(mimeType),
		},
		Data: data,
	}
}

func FailureResponse(code HTTPResponseCode) Response {
	msg := code.responseMessage()
	return Response{
		Version:      V1_1,
		ResponseCode: code,
		ResponseMsg:  msg,
		Headers: HTTPHeaders{
			ContentLength: fmt.Sprintf("%d", len(msg)),
			ContentType:   string(text),
		},
		Data: []byte(msg),
	}
}
