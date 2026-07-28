# Chat Application Backend

A lightweight, concurrent chat application backend built with Go. This service facilitates room-based communication using WebSockets.

---

## 🚀 Features

- **REST API:** Powered by Gin for room creation and management.
- **Real-time Messaging:** Powered by Gorilla WebSocket for low-latency, bi-directional communication.
- **Room Orchestration:** Create rooms with a name and secure them with a password.

---

## 🛠️ Tech Stack

- **Language:** Go 1.22.5
- **Web Framework:** [Gin Gonic](https://github.com/gin-gonic/gin)
- **WebSocket Engine:** [Gorilla WebSocket](https://github.com/gorilla/websocket)

---

## 📋 API Endpoints

| Method   | Endpoint             | Description                     |
| -------- | -------------------- | ------------------------------- |
| **POST** | `/create-room/:name` | Create a new chat room          |
| **POST** | `/join-room/auth`    | Authenticate/Verify room access |
| **GET**  | `/join-room/:name`   | Join a chat room via WebSocket  |

---

## 🏃 Getting Started

### Prerequisites

- Go 1.22.5 or higher installed.

### 1. Clone & Run

```bash
git clone https://github.com/Sahil-sajwan/chat-app-backend.git
cd chat-app-backend
go mod download
go run main.go

```

---

## 🔮 Future Scope

The following milestones are planned for future development:

- **Database Integration:** Integrate **PostgreSQL** using **GORM** to replace the current in-memory storage. This will enable persistent room storage and user metadata.
- **Message Persistence:** Transition from volatile in-memory messaging to database-backed storage, allowing users to retrieve historical messages when joining a room.
- **Authentication & Security:** Enhance room security with JWT-based sessions and refined user management.
