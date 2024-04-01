package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
    "strings"
    "flag"
    "io/ioutil"
)

func handleConnection(conn net.Conn, dirPath string) {

    defer conn.Close()

    buffer := make([]byte, 1024)
    _, err := conn.Read(buffer)

    if err != nil{
		fmt.Println("Error reading data from connection: ", err.Error())
		os.Exit(1)
    }

    strBuf := string(buffer)
    var reqType string
    var startReq int

    if strings.HasPrefix(strBuf, "GET") {
        startReq = strings.Index(strBuf, "GET")
        reqType = "GET"
    } else if strings.HasPrefix(strBuf, "POST") {
        startReq = strings.Index(strBuf, "POST")
        reqType = "POST"
    }

    endReq := startReq + len(reqType) - 1
    startHttp := strings.Index(strBuf, "HTTP")

    reqStr := strBuf[endReq+1:startHttp]

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
    } else if strings.HasPrefix(reqStr, "/files"){
        echoMsg = "files"
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
        response += "Content-Type: text/plain\r\n" + "Content-Length: " + fmt.Sprintf("%d",len(uaContent)) + "\r\n\r\n"  + uaContent
    } else if echoMsg == "files" && reqType == "GET"{

        filePos := strings.Index(reqStr, "/files") + len("/files/")
        fileName := reqStr[filePos:]

        files, err := os.ReadDir(dirPath)
        if err != nil{
            fmt.Printf("Error reading the path: ", err.Error())
            os.Exit(1)
        }

        for _, file := range files{
            if file.Name() == fileName{
                fileContent, err := ioutil.ReadFile(dirPath + "/" + fileName)
                if err != nil{
                    fmt.Println("Error reading the file: ", err.Error())
                    os.Exit(1)
                }
                fileContentStr := string(fileContent)
                response = "HTTP/1.1 200 OK\r\n"
                response += "Content-Type: application/octet-stream\r\n" + "Content-Length: " + fmt.Sprintf("%d",len(fileContentStr)) + "\r\n\r\n"  + fileContentStr
                break
            }
        }
    } else if echoMsg == "files" && reqType == "POST"{
        filePos := strings.Index(reqStr, "/files") + len("/files/")
        fileName := reqStr[filePos:]

        file, err := os.Create(dirPath + "/" + fileName)
        if err != nil{
            fmt.Printf("Could not create the file: %s\n", err.Error())
            os.Exit(1)
        }
        defer file.Close()

        bodyStart := strings.Index(strBuf, "\r\n\r\n")
        bodyEnd := strings.Index(strBuf, "\x00")

        data := strBuf[bodyStart+len("\r\n\r\n"):bodyEnd]

        // fmt.Println(bodyEnd)
        // fmt.Println(data[bodyEnd:])
        //
        // fmt.Println("DATA LEN: ", len(data))
        //
        _, err = file.Write([]byte(data))

        if err != nil {
            fmt.Println("Error writing to file:", err)
            os.Exit(1)
        }

        response = "HTTP/1.1 201 OK\r\n\r\n"
    }

    _, err = conn.Write([]byte(response))

    if err != nil{
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
    }
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
    dirPath := flag.String("directory", "", "directory for the file to be sent")
	flag.Parse()

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
    defer l.Close()
    var conn net.Conn

    for {

        conn, err = l.Accept()
        if err != nil {
            fmt.Println("Error accepting connection: ", err.Error())
            os.Exit(1)
        }

        go handleConnection(conn, *dirPath)
    }
}
