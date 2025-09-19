// client.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr := "127.0.0.1:9000"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Unable to connect:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to", addr)
	// Start goroutine to read from server
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		os.Exit(0) // server closed or connection lost
	}()

	// Read stdin and send lines to server
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		_, err := fmt.Fprintln(conn, line)
		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
		if line == "/quit" {
			// Give server a moment to reply/close
			return
		}
	}
}
