package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Địa chỉ WebSocket server của bạn
	serverAddr := "ws://localhost:8081/ws"

	// Kết nối tới server WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		fmt.Println("Error connecting to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Gửi tin nhắn đến server mỗi 5 giây
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				message := []byte("Hello, From Xuan Vu Server!")
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					fmt.Println("Error sending message:", err)
					return
				}
				fmt.Println("Message sent to server:", string(message))

			// Ngắt kết nối khi nhận được tín hiệu ngắt từ người dùng (Ctrl+C)
			case <-getInterruptSignal():
				fmt.Println("Closing connection...")
				return
			}
		}
	}()

	// Chờ người dùng nhấn Ctrl+C để kết thúc chương trình
	select {}
}

// Hàm trợ giúp để bắt tín hiệu ngắt từ người dùng (Ctrl+C)
func getInterruptSignal() <-chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	return interrupt
}
