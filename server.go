// server.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	Conn     net.Conn
	Username string
	Ch       chan string
}

type Server struct {
	Clients map[string]*Client
	Mutex   sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		Clients: make(map[string]*Client),
	}
}

func (s *Server) Broadcast(sender string, msg string) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	for uname, c := range s.Clients {
		// don't echo to sender (optional)
		if uname == sender {
			continue
		}
		select {
		case c.Ch <- msg:
		default:
			// if client's channel is full, drop message (avoid blocking)
		}
	}
}

func (s *Server) AddClient(c *Client) {
	s.Mutex.Lock()
	s.Clients[c.Username] = c
	s.Mutex.Unlock()
	s.Broadcast("", fmt.Sprintf("*** %s has joined the chat ***", c.Username))
}

func (s *Server) RemoveClient(c *Client) {
	s.Mutex.Lock()
	delete(s.Clients, c.Username)
	s.Mutex.Unlock()
	s.Broadcast("", fmt.Sprintf("*** %s has left the chat ***", c.Username))
	close(c.Ch)
	c.Conn.Close()
}

func (s *Server) ListUsers() []string {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	users := make([]string, 0, len(s.Clients))
	for u := range s.Clients {
		users = append(users, u)
	}
	return users
}

func handleConnection(conn net.Conn, s *Server) {
	defer conn.Close()
	// Prompt for username
	conn.Write([]byte("Enter username: "))
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}
	rawName := strings.TrimSpace(scanner.Text())
	if rawName == "" {
		rawName = conn.RemoteAddr().String()
	}
	// Ensure unique username
	username := rawName
	i := 1
	for {
		s.Mutex.RLock()
		_, exists := s.Clients[username]
		s.Mutex.RUnlock()
		if !exists {
			break
		}
		username = fmt.Sprintf("%s%d", rawName, i)
		i++
	}

	client := &Client{
		Conn:     conn,
		Username: username,
		Ch:       make(chan string, 10),
	}
	s.AddClient(client)

	// Goroutine: write messages from server to client
	go func() {
		writer := bufio.NewWriter(conn)
		for msg := range client.Ch {
			writer.WriteString(msg + "\n")
			writer.Flush()
		}
	}()

	// Tell the connecting client welcome text
	conn.Write([]byte(fmt.Sprintf("Welcome, %s! Commands: /users, /quit\n", client.Username)))
	// Read incoming messages from client
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "/") {
			handleCommand(line, client, s)
			continue
		}
		// Broadcast message
		msg := fmt.Sprintf("%s: %s", client.Username, line)
		s.Broadcast(client.Username, msg)
	}
	// Connection closed / client disconnected
	s.RemoveClient(client)
}

func handleCommand(cmd string, c *Client, s *Server) {
	cmd = strings.TrimSpace(cmd)
	switch {
	case cmd == "/quit":
		c.Conn.Write([]byte("Goodbye!\n"))
		// Remove will close channel and conn
		s.RemoveClient(c)
	case cmd == "/users":
		users := s.ListUsers()
		c.Ch <- fmt.Sprintf("Online (%d): %s", len(users), strings.Join(users, ", "))
	default:
		c.Ch <- "Unknown command. Available: /users, /quit"
	}
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Chat server started on :9000")

	server := NewServer()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleConnection(conn, server)
	}
}
