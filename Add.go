package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func main() {
    // Handle WebSocket connections
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        // Upgrade the HTTP connection to a WebSocket connection
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Println("Upgrade error:", err)
            return
        }
        defer conn.Close()

        // Continuously read messages from the WebSocket connection
        for {
            // Read message from the client
            messageType, p, err := conn.ReadMessage()
            if err != nil {
                log.Println("Read error:", err)
                return
            }

            // Print received message
            fmt.Printf("Received message: %s\n", p)

            // Echo the message back to the client
            err = conn.WriteMessage(messageType, p)
            if err != nil {
                log.Println("Write error:", err)
                return
            }
        }
    })

    // Start the HTTP server on port 8080
    fmt.Println("Server is listening on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
