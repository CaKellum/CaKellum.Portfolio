package badnet

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

// MARK: HTTP Response codes
// TODO: Make some constants for these
type HTTPResponseCode int

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

func handlerFor(req Request) RequestHandler { return pathMap[req.Method][req.Path] }
