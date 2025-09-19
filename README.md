# 💬 Chat Application (Go)

A simple **TCP-based chat application** built with Go.  
It allows multiple clients to connect to a server, send messages, and interact using basic commands.  

---

## 🚀 Features
- TCP server that handles multiple clients concurrently  
- Clients can send and receive broadcast messages  
- Commands:  
  - `/users` → List all connected users  
  - `/quit`  → Disconnect from the server  
- Sends a **welcome message** to each client  
- Graceful client disconnection handling  

---

## 📂 Project Structure
```
chat-app/
├── server.go   # Chat server
├── client.go   # Chat client
├── go.mod
└── README.md
```

---

## ⚡ Usage

### Start the Server
```bash
go run server.go
```

### Start Clients
Open multiple terminals and run:
```bash
go run client.go
```

---

## 💻 Example
**Terminal 1 (Server):**
```
Server started on :9000
```

**Terminal 2 (Client 1):**
```
Connected to chat server.
Welcome 127.0.0.1:54321!
127.0.0.1:54321 joined the chat
/users
Online users: 127.0.0.1:54321
Hello everyone!
```

**Terminal 3 (Client 2):**
```
Connected to chat server.
Welcome 127.0.0.1:54322!
127.0.0.1:54322 joined the chat
127.0.0.1:54321: Hello everyone!
/quit
Goodbye!
```

---

## 🧑‍💻 Future Improvements
- Add authentication (username/password)  
- Support private messaging (`/msg user message`)  
- Save chat history to a file or database  
- Add WebSocket version for browsers  

---

## 📜 License
MIT License – free to use and modify.
