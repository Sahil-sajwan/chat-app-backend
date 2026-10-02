# Chat Application Backend

A lightweight, concurrent chat application backend built with Go and MongoDB. This service facilitates room-based real-time communication using WebSockets with persistent storage for rooms and message history.

---

## 🚀 Features

- **REST API:** Powered by Gin for room creation, password validation, and join-token authorization.
- **Real-time Messaging:** Powered by Gorilla WebSocket for low-latency, bi-directional communication.
- **Room Management:** Create password-protected chat rooms saved directly to MongoDB.
- **Persistent Storage:** Integrated MongoDB to persist room metadata and chat history.
- **Chat History Replay:** Automatically streams past message history to newly connected users upon joining a room.
- **Token Authorization & Graceful Reconnects:** Temporary join tokens for WebSocket security and disconnect grace periods for seamless reconnects.

---

## 🛠️ Tech Stack

- **Language:** Go 1.22.5
- **Web Framework:** [Gin Gonic](https://github.com/gin-gonic/gin)
- **WebSocket Engine:** [Gorilla WebSocket](https://github.com/gorilla/websocket)
- **Database:** [MongoDB](https://www.mongodb.com/) (using `go.mongodb.org/mongo-driver/v2`)

---

## 📋 API Endpoints

| Method   | Endpoint             | Description                                                   |
| -------- | -------------------- | ------------------------------------------------------------- |
| **POST** | `/create-room/:name` | Create a new chat room (Body: `rpass`)                        |
| **POST** | `/join-room/auth`    | Authenticate room password & receive join token (`rname`, `rpass`) |
| **GET**  | `/join-room/:name`   | Join a chat room via WebSocket query (`?rname=...&token=...`) |

---

## 🏃 Getting Started

### Prerequisites

- **Go** 1.22.5 or higher installed.
- **MongoDB** running locally (`mongodb://localhost:27017`) or accessible via remote URI (e.g. MongoDB Atlas).

### Environment Setup (Optional)

Specify your MongoDB connection string using the `MONGO_URI` environment variable:

```bash
export MONGO_URI="mongodb://localhost:27017"
```

*(Defaults to `mongodb://localhost:27017` if omitted)*

### Clone & Run

```bash
git clone https://github.com/Sahil-sajwan/chat-app-backend.git
cd chat-app-backend
go mod download
go run main.go
```

The backend server runs at `http://localhost:8080`.

---

## 🔮 Future Scope

The following milestones are planned for future development:

- **JWT Authentication:** Implement JWT-based user authentication and user accounts.
- **Pagination for History:** Support cursor/offset-based pagination for loading historical messages dynamically.
- **Media & File Attachments:** Support sharing media, images, and file attachments in chat rooms.

