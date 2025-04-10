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
	OK:       "HTTP/1.1 200 OK\r\n",
	NotFound: "HTTP/1.1 404 Not Found\r\n\r\n",
}

type contentType int

const (
	TextPlain contentType = iota
)

var contentTypeMap = map[contentType]string{
	TextPlain: "text/plain\\r\\n",
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

	var buff = make([]byte, 1024)

	_, err := conn.Read(buff)

	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}

	urlParts := strings.Split(getUrl(string(buff)), "/")

	// for _, part := range urlParts {
	// 	fmt.Printf("Part %s\n", part)
	// }

	switch urlParts[1] {
	case "":
		msg := "Hello World"
		produceResponse(conn, msg, OK, TextPlain, len(msg))
	case "echo":
		msg := urlParts[2]
		// fmt.Printf("Echo %s\n", msg)
		produceResponse(conn, msg, OK, TextPlain, len(msg))
	default:
		msg := ""
		produceResponse(conn, msg, NotFound, TextPlain, len(msg))
	}

	defer conn.Close()

	// if strings.Compare(reqPart, "/") == 0 {
	// 	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	// } else {
	// 	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	// }

}

func getUrl(url string) string {
	parts := strings.Split(url, " ")

	// for _, part := range parts {
	// 	fmt.Printf("Part %s\n", part)
	// }

	return parts[1]
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
