package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
    "strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

    var conn net.Conn

    conn, err = l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

    buffer := make([]byte, 1024)
    _, err = conn.Read(buffer)

    if err != nil{
		fmt.Println("Error reading data from connection: ", err.Error())
		os.Exit(1)
    }

    strBuf := string(buffer)

    startGet := strings.Index(strBuf, "GET")
    endGet := startGet + len("GET") - 1
    startHttp := strings.Index(strBuf, "HTTP")

    reqStr := strBuf[endGet+1:startHttp]

    reqStr = strings.TrimSpace(reqStr)


    var randStart int
    var randStr string
    var echoMsg string

    response := "HTTP/1.1 404 Not Found\r\n\r\n"

    if reqStr == "/" {
        response = "HTTP/1.1 200 OK\r\n\r\n"
    }else if strings.HasPrefix(reqStr, "/echo"){
        randStart = strings.Index(reqStr[1:], "/")
        randStr = reqStr[randStart+2:]
        echoMsg = reqStr[1:randStart+1]
    } else if strings.HasPrefix(reqStr, "/user-agent"){
        echoMsg = "user-agent"
    }

    if reqStr == "/" {
        response = "HTTP/1.1 200 OK\r\n\r\n"
    }else if echoMsg == "echo"{
        response = "HTTP/1.1 200 OK\r\n"
        response += "Content-Type: text/plain\r\n" + "Content-Length: " + fmt.Sprintf("%d",len(randStr)) + "\r\n\r\n"  + randStr
    }else if echoMsg == "user-agent"{
        uaPos := strings.Index(strBuf, "User-Agent:")
        uaEnd := strings.Index(strBuf[uaPos:], "\r")
        uaContent := strings.TrimSpace(strBuf[uaPos+len("User-Agent:"):uaPos+uaEnd])
        response = "HTTP/1.1 200 OK\r\n"
        response += "Content-Type: text/plain\r\n" + "Content-Length: " + fmt.Sprintf("%d",uaEnd) + "\r\n\r\n"  + uaContent
    }

    _, err = conn.Write([]byte(response))

    if err != nil{
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
    }
}
