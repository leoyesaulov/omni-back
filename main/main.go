package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func main() {
	http.Handle("/ws", websocket.Handler(handleConnection))
	log.Println("WebSocket server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnection(ws *websocket.Conn) {
	defer func() {
		ws.Close()
		log.Printf("[%s] Connection closed", getClientInfo(ws))
	}()

	clientIP := getClientIP(ws.Request())
	log.Printf("[%s] New connection attempt | Origin: %s | Protocol: %s",
		clientIP,
		ws.Request().Header.Get("Origin"),
		ws.Request().Proto,
	)

	var msg string
	err := websocket.Message.Receive(ws, &msg)
	if err != nil {
		log.Printf("[%s] Connection error: %v", clientIP, err)
		return
	}

	log.Printf("[%s] Successfully upgraded to WebSocket", clientIP)
	log.Printf("[%s] Received message: %s", clientIP, msg)

	// Обработка сообщения
	if msg == "Hello, backend" {
		websocket.Message.Send(ws, "Hello from Go server!")
		log.Printf("[%s] Sent response", clientIP)
	} else {
		websocket.Message.Send(ws, "Unknown request")
		log.Printf("[%s] Invalid request received: %s", clientIP, msg)
	}
}

func getClientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "unknown"
	}
	return ip
}

func getClientInfo(ws *websocket.Conn) string {
	return fmt.Sprintf("%s | %s",
		getClientIP(ws.Request()),
		time.Now().Format("2006-01-02 15:04:05"),
	)
}
