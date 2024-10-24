package main

import (
    "bytes"
    "fmt"
    "net"
    "strings"
)

// Function to handle incoming connections
func main() {
    // Listen on port 6379
    listener, err := net.Listen("tcp", ":6379")
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    defer listener.Close()

    fmt.Println("Listening on port :6379")

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        go handleConnection(conn)
    }
}

// Handle incoming connections
func handleConnection(conn net.Conn) {
    defer conn.Close()
    fmt.Println("Client connected:", conn.RemoteAddr())

    var buffer bytes.Buffer
    readBuffer := make([]byte, 1024)

    for {
        n, err := conn.Read(readBuffer)
        if err != nil {
            fmt.Println("Error reading:", err)
            return
        }

        if n == 0 {
            fmt.Println("Client disconnected:", conn.RemoteAddr())
            return
        }

        // Accumulate the buffer
        buffer.Write(readBuffer[:n])

        // Process the buffer if it contains a complete RESP command
        processCommand(&buffer, conn)
    }
}

// Process incoming RESP commands
func processCommand(buffer *bytes.Buffer, conn net.Conn) {
    command := buffer.String()

    // Check if the command is complete
    if strings.Contains(command, "\r\n") {
        // Split the command based on the RESP format
        parts := strings.Split(command, "\r\n")
        if len(parts) > 0 {
            // The first part is the type of command
            if parts[0] == "*1" && len(parts) > 2 && parts[2] == "PING" {
                response := "+PONG\r\n"
                conn.Write([]byte(response))
                buffer.Reset() // Clear the buffer after processing
                fmt.Println("Responded with PONG")
            } else {
                response := "-ERR unknown command\r\n"
                conn.Write([]byte(response))
                buffer.Reset() // Clear the buffer after processing
                fmt.Println("Responded with error for unknown command")
            }
        }
    }
}
