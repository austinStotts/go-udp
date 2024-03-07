package main

import (
    "fmt"
    "net"
)

func main() {
    // Listen for UDP packets on port 8080
    udpAddr, err := net.ResolveUDPAddr("udp", ":8080")
    if err != nil {
        fmt.Println("Error resolving UDP address:", err)
        return
    }

    // Create UDP connection
    udpConn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        fmt.Println("Error listening:", err)
        return
    }
    defer udpConn.Close()

    fmt.Println("UDP server is listening on port 8080...")

    // Buffer for receiving data
    buffer := make([]byte, 1024)

    for {
        // Read data from UDP connection
        n, addr, err := udpConn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Println("Error reading from UDP:", err)
            continue
        }

        // Print received data
        fmt.Printf("Received %d bytes from %s: %s\n", n, addr, string(buffer[:n]))

        // Echo received data back to the client
        _, err = udpConn.WriteToUDP(buffer[:n], addr)
        if err != nil {
            fmt.Println("Error writing to UDP:", err)
            continue
        }
    }
}
