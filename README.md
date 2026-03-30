# Eltunnel

A lightweight tunneling tool (similar to ngrok) built in Go and Python that exposes local applications through a public URL using persistent WebSocket connections.

---

## Overview

Eltunnel lets you run any app locally and access it through a generated subdomain.

No deployment. No cloud setup. Just run and share.

---

## How It Works

1. Python client connects to the Go server via WebSocket  
2. Client registers a random subdomain  
3. Go server maps subdomain → client connection  
4. Incoming HTTP requests are routed through the tunnel  
5. Client forwards request to localhost  
6. Response is sent back through the tunnel to the browser  

---

## Project Structure

### Server (Go)

- `server/protocol/messages.go`  
  Defines the protocol (register, http_request, http_response, error)

- `server/tunnel/manager.go`  
  Manages active tunnels (subdomain → connection), thread-safe with mutex

- `server/tunnel/handler.go`  
  Handles WebSocket connections and response routing

- `server/proxy/handler.go`  
  Routes incoming HTTP requests through the correct tunnel

- `server/main.go`  
  Entry point  
  - Port `8081` → WebSocket clients  
  - Port `8080` → HTTP traffic  

---

### Client (Python)

- `client/protocol/messages.py`  
  Python version of the shared protocol

- `client/tunnel/client.py`  
  Connects to server, registers subdomain, listens for requests

- `client/proxy/forwarder.py`  
  Forwards requests to localhost and returns responses  

---

## Full Flow

1. Run client → registers subdomain  
2. Open `<subdomain>.localhost:8080`  
3. Server receives request  
4. Routes through WebSocket tunnel  
5. Client forwards to localhost  
6. Response returns to browser  

---

## Running the Project

### 1. Start local app
```bash
python3 -m http.server 3000
```

### 2. Start Go server
```bash
cd server
go run main.go
```

### 3. Start Python client
```bash
cd client
PYTHONPATH=. python3 tunnel/client.py
```

### 4. Test
```bash
curl -H "Host: <subdomain>.localhost" http://localhost:8080
```

---

## Performance

- ~7ms average latency (local testing)
- Efficient concurrency via Go goroutines + channels
- Identified limitation: Python client currently processes requests sequentially

---

## Design Decisions

- WebSockets for persistent bidirectional communication  
- Go for concurrent, high-performance server  
- Python for simple client implementation  
- Custom protocol for request/response handling  
- Request IDs to prevent response mismatches  

---

## Use Cases

- Share local projects instantly  
- Test webhooks  
- Demo apps without deployment  
- Access local apps from other devices  

---

## Future Improvements

- Async request handling in Python (fix concurrency bottleneck)  
- Auto-reconnect support  
- Logging + metrics  
- Timeouts and failure handling  
- Deploy with real domain + wildcard DNS  
