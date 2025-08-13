# Go Real-Time Chat App

A learning project to explore **Go** concepts through building a simple WebSocket-based real-time chat.  
The goal is to practice **HTTP servers**, **goroutines**, **channels**, and **concurrency patterns** in Go while also touching databases, persistence, and deployment.

> **Note:** This is not intended as a production-grade chat app.  

---

## Milestone Plan

### **Milestone 1 – Basic HTTP Server**
- Set up a simple HTTP server using **Chi** (or Fiber/Gin).
- Serve a static HTML file (`index.html`) from a `static/` folder.
- Test by opening the page via the Go server.

---

### **Milestone 2 – Simple WebSocket Endpoint**
- Add a `/ws` route that upgrades HTTP to WebSocket using `gorilla/websocket`.
- On connection:
  - Goroutine to read messages from the client.
  - Goroutine to send test messages from server to client.
- Connect from the frontend using native JavaScript `WebSocket` API.

---

### **Milestone 3 – Broadcast Hub**
- Create a `Hub` struct to manage:
  - `register` channel
  - `unregister` channel
  - `broadcast` channel
  - `clients` map
- Hub runs in its own goroutine, processing messages from these channels.
- Each client runs dedicated **read** and **write** goroutines.
- When a message is received → push to `broadcast` → Hub sends to all clients.

---

### **Milestone 4 – Basic Frontend UI**
- Simple HTML + JS:
  - Text input + send button.
  - Scrollable message list.
- Append incoming messages to the list.
- Send text from the input box to `/ws`.

---

### **Milestone 5 – Handle Disconnects & Errors**
- Remove disconnected clients from `clients` map.
- Ensure goroutines exit cleanly to avoid memory leaks.
- Add server logs for connection lifecycle events.

---

### **Milestone 6 – Add Usernames**
- On connect, let the frontend send a username (via first message or query param).
- Broadcast messages with the username prefix.

---

### **Milestone 7 – Persistence (Optional but Recommended)**
- Integrate SQLite/PostgreSQL with `sqlx` or `gorm`.
- Save messages with username + timestamp.
- Send recent chat history to new clients on connect.

---

### **Milestone 8 – Async Features**
- Add background goroutines for:
  - Periodic DB writes in batches.
  - Cleaning up inactive clients.
  - Scheduled announcements.

---

### **Milestone 9 – Deployment**
- Deploy to **Render** or **Fly.io**.
- Test with real users.
- Monitor logs and handle basic scaling.

---

## How to Use This Plan
- Each milestone should be a separate commit (e.g., `feat: milestone 3 – added broadcast hub`).
- Keep a `LEARNINGS.md` file with bullet points of what was learned in each milestone.
- Avoid adding unnecessary features too early — focus on mastering the core concepts first.
