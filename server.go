package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	clients   = make(map[net.Conn]string) // connection -> username
	broadcast = make(chan string)
	mutex     = sync.Mutex{}
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Assign a default username
	username := conn.RemoteAddr().String()

	mutex.Lock()
	clients[conn] = username
	mutex.Unlock()

	// ✅ Send welcome message only to this client
	conn.Write([]byte(fmt.Sprintf("Welcome %s!\n", username)))

	// Broadcast join event to others
	broadcast <- fmt.Sprintf("%s joined the chat", username)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()

		// Commands
		if strings.HasPrefix(msg, "/") {
			switch msg {
			case "/users":
				mutex.Lock()
				var userList []string
				for _, name := range clients {
					userList = append(userList, name)
				}
				mutex.Unlock()
				conn.Write([]byte("Online users: " + strings.Join(userList, ", ") + "\n"))

			case "/quit":
				conn.Write([]byte("Goodbye!\n"))
				mutex.Lock()
				delete(clients, conn)
				mutex.Unlock()
				broadcast <- fmt.Sprintf("%s left the chat", username)
				return

			default:
				conn.Write([]byte("Unknown command\n"))
			}
			continue
		}

		// Broadcast normal message
		broadcast <- fmt.Sprintf("%s: %s", username, msg)
	}

	// Handle disconnect (e.g., Ctrl+C)
	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
	broadcast <- fmt.Sprintf("%s disconnected", username)
}

func broadcaster() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for conn := range clients {
			fmt.Fprintln(conn, msg)
		}
		mutex.Unlock()
	}
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Server started on :9000")

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}
