package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	ID         int         `json:"id"`
	Body       string      `json:"body"`
	OwnerCarID int         `json:"owner_car_id"`
	PreOrderID int         `json:"pre_order_id"`
	ListTruck  []ListTruck `json:"list_truck"`
}

type ListTruck struct {
	TruckCarID int `json:"truck_car_id"`
}

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{ID: messageType, Body: string(p)}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{ID: 1, Body: "New User Joined..."})
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{ID: 1, Body: "User Disconnected..."})
			}
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

func websocketRandPing(conn *websocket.Conn) {
	for {
		err := conn.WriteMessage(websocket.TextMessage, []byte("randping"))
		if err != nil {
			log.Println(err)
			return
		}
		time.Sleep(time.Duration(rand.Intn(int(time.Second * 3))))
	}
}

func serveWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	go websocketRandPing(conn)
	client.Read()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

var (
	pool            *Pool
	mockMessageConn *websocket.Conn
)

func setupRoutes() {
	pool := NewPool()
	go pool.Start()

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		homePage(w, r)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})

	//http.HandleFunc("/sendMockMessage", func(w http.ResponseWriter, r *http.Request) {
	//	sendMockMessage(w, r, pool)
	//})
}

//func createMockMessage() Message {
//	// Tạo và trả về một tin nhắn giả mạo
//	return Message{
//		ID:         1,
//		Body:       "Mock Message",
//		OwnerCarID: 1,
//		PreOrderID: 1,
//		TruckID:    1,
//	}
//}

//func sendMockMessage(w http.ResponseWriter, r *http.Request, pool *Pool) {
//	// Gửi tin nhắn mock vào kết nối WebSocket
//	mockMessage := createMockMessage()
//	pool.Broadcast <- mockMessage
//
//	fmt.Fprintf(w, "Sent Mock Message")
//}

func main() {
	setupRoutes()
	http.ListenAndServe("178.128.57.128:8081", nil)
}
