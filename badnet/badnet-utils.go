package badnet

import (
	"fmt"
	"maps"
	"regexp"
)

const registerNewLine = "\r\n"

// MARK: HTTP Header type
type HTTPHeaders map[string]string

// MARK: HTTP Request Headers Keys
const (
	Host                    = "Host"
	UserAgent               = "User-Agent"
	Accept                  = "Accept"
	AcceptLang              = "Accept-Language"
	Referer                 = "Referer"
	Connection              = "Connection"
	UpgradeInsecureRequests = "Upgrade-Insecure-Requests"
	IfModSince              = "If-Modified-Since"
	IfNoneMatch             = "If-None-Match"
	CacheCtrl               = "Cache-Control"
	ContentLength           = "Contnet-Length"
	ContentType             = "Contnet-Type"
)

type HTTPMimeType string

const (
	// MARK: image types
	gif  HTTPMimeType = "image/gif"
	jpg  HTTPMimeType = "image/jpeg"
	webp HTTPMimeType = "image/webp"
	png  HTTPMimeType = "image/png"
	heic HTTPMimeType = "image/heic"

	// MARK: text types
	text       HTTPMimeType = "text/plain"
	html       HTTPMimeType = "text/html"
	css        HTTPMimeType = "text/css"
	javaScript HTTPMimeType = "text/javascript"
	md         HTTPMimeType = "text/markdown"
)

// MARK: HTTP Response codes
// TODO: Make some constants for these
type HTTPResponseCode int

const (
	// MARK: Informational codes
	Continue   HTTPResponseCode = 100
	Switching  HTTPResponseCode = 101
	Processing HTTPResponseCode = 102
	Earlyhints HTTPResponseCode = 103

	// MARK: Successful codes
	OK                          HTTPResponseCode = 200
	Created                     HTTPResponseCode = 201
	Accepted                    HTTPResponseCode = 202
	NonauthoritativeInformation HTTPResponseCode = 203
	NoContent                   HTTPResponseCode = 204
	ResetContent                HTTPResponseCode = 205
	PartialContent              HTTPResponseCode = 206
	MultiStatus                 HTTPResponseCode = 207
	AlreadyReported             HTTPResponseCode = 208
	IMUsed                      HTTPResponseCode = 226

	// MARK: Redirect messages
	MultipleChoices   HTTPResponseCode = 300
	MovedPermanently  HTTPResponseCode = 301
	Found             HTTPResponseCode = 302
	SeeOther          HTTPResponseCode = 303
	NotModified       HTTPResponseCode = 304
	TemporaryRedirect HTTPResponseCode = 307
	PermanentRedirect HTTPResponseCode = 308

	// MARK: Client Error
	BadRequest                  HTTPResponseCode = 400
	Unauthorized                HTTPResponseCode = 401
	PaymentRequired             HTTPResponseCode = 402
	Forbidden                   HTTPResponseCode = 403
	NotFound                    HTTPResponseCode = 404
	MethodNotAllowed            HTTPResponseCode = 405
	NotAcceptable               HTTPResponseCode = 406
	ProxyAuthenticationRequired HTTPResponseCode = 407
	RequestTimeOut              HTTPResponseCode = 408
	Conflict                    HTTPResponseCode = 409
	Gone                        HTTPResponseCode = 410
	LengthRequired              HTTPResponseCode = 411
	PreconditionFailed          HTTPResponseCode = 412
	ContentTooLarge             HTTPResponseCode = 413
	URITooLong                  HTTPResponseCode = 414
	UnsupportedMediaType        HTTPResponseCode = 415
	RangeNotSatisfiable         HTTPResponseCode = 416
	ExpectationFailed           HTTPResponseCode = 417
	ImATeapot                   HTTPResponseCode = 418
	MisdirectedRequest          HTTPResponseCode = 421
	UnprocessableContent        HTTPResponseCode = 422
	Locked                      HTTPResponseCode = 423
	FailedDependency            HTTPResponseCode = 424
	TooEarly                    HTTPResponseCode = 425
	UpgradeRequired             HTTPResponseCode = 426
	PreconditionRequired        HTTPResponseCode = 428
	TooManyRequest              HTTPResponseCode = 429
	RequestHeaderFieldsTooLarge HTTPResponseCode = 431
	UnavailableForLegalReasons  HTTPResponseCode = 451

	// MARK: Server Error
	InternalServerError          HTTPResponseCode = 500
	NotImplemented               HTTPResponseCode = 501
	BadgateWay                   HTTPResponseCode = 502
	ServiceUnavailable           HTTPResponseCode = 503
	GatewayTimeout               HTTPResponseCode = 504
	HTTPVersionNotSupported      HTTPResponseCode = 505
	VariantAlsoNegotiates        HTTPResponseCode = 506
	InsufficientStorage          HTTPResponseCode = 507
	LoopDetected                 HTTPResponseCode = 508
	NotExtended                  HTTPResponseCode = 510
	NetworAuthenticationRequired HTTPResponseCode = 511
)

func (code HTTPResponseCode) responseMessage() string {
	msg := ""

	switch code {
	case OK:
		msg = "Ok"
	case InternalServerError:
		msg = "Internal Server Error"
	case NotFound:
		msg = "Not Found"
	}

	return msg
}

// MARK: HTTP Request Methods
type HTTPRequestMethod string

const (
	GET     HTTPRequestMethod = "GET"
	HEAD    HTTPRequestMethod = "HEAD"
	POST    HTTPRequestMethod = "POST"
	PUT     HTTPRequestMethod = "PUT"
	DELETE  HTTPRequestMethod = "DELETE"
	CONNECT HTTPRequestMethod = "CONNECT"
	OPTIONS HTTPRequestMethod = "OPTIONS"
	TRACE   HTTPRequestMethod = "TRACE"
	PATCH   HTTPRequestMethod = "PATCH"
)

// MARK: HTTP Versions ~for now only support 1.1
type HTTPVersion string

const (
	V0_9 HTTPVersion = "HTTP/0.9"
	V1_1 HTTPVersion = "HTTP/1.1"
	V2   HTTPVersion = "HTTP/2"
)

// MARK: HTTP Core types
type Request struct {
	Path    string
	Version HTTPVersion
	Method  HTTPRequestMethod
	Headers HTTPHeaders
	Data    []byte
}

func emptyRequest() Request {
	return Request{
		Path:    "",
		Version: V1_1,
		Method:  GET,
		Headers: make(HTTPHeaders),
		Data:    nil,
	}
}

type Response struct {
	Version      HTTPVersion
	ResponseCode HTTPResponseCode
	ResponseMsg  string
	Headers      HTTPHeaders
	Data         []byte
}

// MARK: Server Utils

type RequestHandler func(Request) Response
type pathSectionMapping map[string]RequestHandler

var pathMap map[HTTPRequestMethod]pathSectionMapping = make(map[HTTPRequestMethod]pathSectionMapping)

func (method HTTPRequestMethod) RegisterPath(path string, handler RequestHandler) {
	if pathMap[method] == nil {
		pathMap[method] = make(pathSectionMapping)
	}
	pathMap[method][path] = handler
}

func invalidPath(req Request) Response {

	return Response{
		ResponseCode: 404,
		ResponseMsg:  "File Not Found",
		Version:      V1_1,
		Headers: HTTPHeaders{
			ContentType:   "text/plain",
			ContentLength: fmt.Sprintf("%d", len("File Not Found")),
		},
		Data: []byte("File Not Found"),
	}
}

func handlerFor(req Request) RequestHandler {
	paths := pathMap[req.Method]
	for path := range maps.Keys(paths) {
		matched, _ := regexp.Match(path, []byte(req.Path))
		if matched {
			return paths[path]
		}
	}
	return invalidPath
}
