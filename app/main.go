package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var _ = net.Listen
var _ = os.Exit

type HttpResponseStatusCode int

const (
	OK HttpResponseStatusCode = iota
	NotFound
)

var httpResponseStatusCodeMap = map[HttpResponseStatusCode]string{
	OK:       "HTTP/1.1 200 OK\r\n\r\n",
	NotFound: "HTTP/1.1 404 Not Found\r\n\r\n",
}

type contentType int

const (
	TextPlain contentType = iota
)

var contentTypeMap = map[contentType]string{
	TextPlain: "text/plain\r\n",
}

type Header struct {
	Name, Value string
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var buff = make([]byte, 1024)

	_, err := conn.Read(buff)

	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}

	urlParts := getUrl(buff)
	headers := getHeaders(buff)

	switch urlParts[1] {
	case "":
		msg := ""
		produceResponse(conn, msg, OK, TextPlain, len(msg))
	case "echo":
		msg := urlParts[2]
		produceResponse(conn, msg, OK, TextPlain, len(msg))
	case "user-agent":
		msg, err := getHeaderValue(headers, "User-Agent")
		if err != nil {
			fmt.Printf("Error getting header value: %s\n", err.Error())
			produceResponse(conn, "", NotFound, TextPlain, 0)
			return
		}

		produceResponse(conn, msg, OK, TextPlain, len(msg))
	default:
		msg := ""
		produceResponse(conn, msg, NotFound, TextPlain, len(msg))
	}
}

func getUrl(buff []byte) []string {
	parts := strings.Split(string(buff), " ")
	return strings.Split(parts[1], "/")
}

func getHeaders(buff []byte) []Header {
	lines := strings.Split(string(buff), "\r\n")

	headersValues := lines[3:]

	headers := make([]Header, 0, len(headersValues))

	for _, header := range headersValues {
		headerParts := strings.Split(header, ": ")

		if len(headerParts) != 2 {
			continue
		}

		headers = append(headers, Header{
			Name:  headerParts[0],
			Value: headerParts[1],
		})
	}

	return headers
}

func getHeaderValue(headers []Header, name string) (string, error) {
	name = strings.ToLower(name)
	for _, header := range headers {
		if name == strings.ToLower(header.Name) {
			return header.Value, nil
		}
	}

	return "", fmt.Errorf("header %s not found", name)
}

func produceResponse(
	conn net.Conn,
	msg string,
	statusCode HttpResponseStatusCode,
	contentType contentType,
	contentLength int,
) {
	resp := httpResponseStatusCodeMap[statusCode] + "Content-Type: " + contentTypeMap[contentType] + "Content-Length: " + strconv.Itoa(contentLength) + "\r\n\r\n" + msg
	fmt.Printf("%s\n", resp)
	conn.Write([]byte(resp))
}
