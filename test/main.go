package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// 定义全局变量
var clients = make(map[*websocket.Conn]bool) // 已连接的客户端
var broadcast = make(chan Message)           // 消息广播通道
var mutex = &sync.Mutex{}                    // 互斥锁，用于同步操作

// 消息结构
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

// WebSocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 处理 WebSocket 连接
// 处理连接
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 将HTTP连接升级为WebSocket连接
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// 将新连接的客户端添加到客户端列表
	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	for {
		// 读取客户端发送的消息
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			// 从客户端列表中删除该客户端
			delete(clients, ws)
			break
		}

		// 将消息发送到广播通道
		broadcast <- msg
	}
}

// 处理消息广播
// 处理消息
func handleMessages() {
	for {
		// 从广播通道中接收消息
		msg := <-broadcast
		// 将消息发送给所有连接的客户端
		for client := range clients {
			// 将消息发送给客户端
			err := client.WriteJSON(msg)
			// 如果发送失败，则关闭客户端连接，并从客户端列表中删除
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// 处理WebSocket连接
	http.HandleFunc("/ws", handleConnections)

	// 处理消息
	go handleMessages()

	// 提供静态文件服务
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Server started at :9090")
	// 监听端口
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
