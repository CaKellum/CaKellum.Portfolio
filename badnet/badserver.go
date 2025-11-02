package badnet

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

var serverLogger io.Writer

type ServerConfiguration struct {
	MaximumRequestSize int
	Logger             io.Writer
	Port               string
	Network            string
}

func (config ServerConfiguration) fillEmpty() {
	defaultConfig := defaultConfig()
	if config.Port == "" {
		config.Port = defaultConfig.Port
	}
	if config.Network == "" {
		config.Network = defaultConfig.Network
	}
	if config.MaximumRequestSize == 0 {
		config.MaximumRequestSize = defaultConfig.MaximumRequestSize
	}
}

func defaultConfig() ServerConfiguration {
	return ServerConfiguration{
		MaximumRequestSize: 1024,
		Logger:             nil,
		Port:               ":8080",
		Network:            "tcp",
	}
}

func parseRequest(buff []byte) (Request, error) {
	req := emptyRequest()
	temp := strings.Split(string(buff), fmt.Sprintf("%s%s", registerNewLine, registerNewLine))
	strHeaders := temp[0]
	var body []byte
	if len(temp) == 2 {
		body = []byte(temp[1])
	}
	headerGroups := strings.Split(strHeaders, registerNewLine)
	statusLine := strings.Split(headerGroups[0], " ")
	req.Method = HTTPRequestMethod(statusLine[0])
	req.Path = statusLine[1]
	req.Version = HTTPVersion(statusLine[2])
	for _, headerGroup := range headerGroups[1:] {
		keyVal := strings.Split(headerGroup, ":")
		req.Headers[keyVal[0]] = strings.TrimSpace(keyVal[1])
	}
	if sizeStr := req.Headers[ContentLength]; sizeStr != "" {
		size, sizeErr := strconv.Atoi(sizeStr)
		if sizeErr == nil {
			body = body[:size]
		}
	}
	req.Data = body
	return req, nil
}

func writeHeaders(headers HTTPHeaders) string {
	var headerStr string
	for key, value := range headers {
		headerStr += fmt.Sprintf("%s: %s %s", key, value, registerNewLine)
	}
	return headerStr + registerNewLine
}

func writeResponseToConnection(res Response, conn *net.Conn) {
	headerString := writeHeaders(res.Headers)
	statusLine := fmt.Sprintf("%s %d %s", res.Version, res.ResponseCode, res.ResponseMsg)
	fmt.Fprint(*conn, statusLine)
	fmt.Fprint(*conn, headerString)
	(*conn).Write(res.Data)
}

func handleConnection(conn net.Conn) error {
	defer conn.Close()
	buff := make([]byte, 1460)
	conn.Read(buff)
	req, reqErr := parseRequest(buff)
	if reqErr != nil {
		return reqErr
	}

	handler := handlerFor(req)
	if handler == nil {
		msg := fmt.Sprintf("Error No Handler for %s", req.Path)
		slog(msg)
		return nil
	}

	res := handler(req)
	res.Headers[ContentLength] = fmt.Sprintf("%d", len(res.Data))
	writeResponseToConnection(res, &conn)
	return nil
}

// short for server log
func slog(msg string) {
	if serverLogger != nil {
		serverLogger.Write([]byte(msg))
	} else {
		fmt.Println(msg)
	}
}

func StartServer(config ServerConfiguration) {
	config.fillEmpty()
	l, lErr := net.Listen(config.Network, config.Port)
	serverLogger = config.Logger
	if lErr != nil {
		msg := fmt.Sprintf("%v", lErr)
		slog(msg)
		return
	}
	for {
		conn, connErr := l.Accept()
		if connErr != nil {
			msg := fmt.Sprintf("%v", lErr)
			slog(msg)
			continue
		}
		handleErr := handleConnection(conn)
		if handleErr != nil {
			msg := fmt.Sprintf("%v", lErr)
			slog(msg)
			continue
		}
	}
}
