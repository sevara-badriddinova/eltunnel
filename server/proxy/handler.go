package proxy

import (
	"eltunnel/server/protocol"
	"eltunnel/server/tunnel"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// handles incoming HTTP requests from internet and forwards to the right tunnel
func HandleRequest(manager *tunnel.Manager, w http.ResponseWriter, r *http.Request) {
	// xyz.localhost:8081
	host := r.Host
	// split into ["xyz", "localhost:8081"]
	parts := strings.Split(host, ".")
	// get subdomain
	subdomain := parts[0]
	fmt.Println(subdomain)

	conn, exists := manager.GetClient(subdomain)
	if !exists {
		http.Error(w, "tunnel not found", http.StatusNotFound)
		return
	} else {
		fmt.Println("found tunnel for: " + subdomain)
	}

	// body - content of message, headers - notes outside of envelope
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	headers := make(map[string]string)
	for key, values := range r.Header {
		headers[key] = values[0]
	}

	// pack into standard format
	requestPayload := protocol.HttpRequestPayload{
		Method:  r.Method,
		Path:    r.URL.Path,
		Body:    string(body),
		Headers: headers,
	}

	// unique id with timestamp
	requestID := fmt.Sprintf("%d", time.Now().UnixNano())

	// pack requestPayload into bytes, wrap in message with type and ticket number
	// pack the whole message into bytes ready to send
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		http.Error(w, "failed to marshal payload", http.StatusInternalServerError)
		return
	}
	message := protocol.Message{
		Type:      "http_request",
		RequestID: requestID,
		Payload:   payloadBytes,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "failed to marshal message", http.StatusInternalServerError)
		return
	}

	// creating a waiting channel
	responseChan := make(chan protocol.HttpResponsePayload)
	manager.AddPendingRequest(requestID, responseChan)

	// send down the tunnel
	err = conn.WriteMessage(websocket.TextMessage, messageBytes)
	fmt.Println("waiting for response...")
	if err != nil {
		http.Error(w, "failed to send message", http.StatusInternalServerError)
		return
	}

	// wait for ticket number to be called
	response := <-responseChan
	// response arrives, sent to browser
	for key, value := range response.Headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(response.StatusCode)
	w.Write([]byte(response.Body))

}
