package main

import (
    "fmt"
    "net"
    "time"
)

func main() {
    // Server address
    serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
    if err != nil {
        fmt.Println("Error resolving server address:", err)
        return
    }

    // Establish UDP connection
    conn, err := net.DialUDP("udp", nil, serverAddr)
    if err != nil {
        fmt.Println("Error connecting to server:", err)
        return
    }
    defer conn.Close()

    // Data to be sent
    data := []byte("Hello, UDP server!")

    // Interval for sending data (in milliseconds)
    interval := 1000 // 1 second

    fmt.Println("UDP client is sending data to server...")

    // Send data at set intervals
    for {
        _, err := conn.Write(data)
        if err != nil {
            fmt.Println("Error sending data:", err)
            return
        }

        fmt.Printf("Sent: %s\n", data)

        time.Sleep(time.Duration(interval) * time.Millisecond)
    }
}