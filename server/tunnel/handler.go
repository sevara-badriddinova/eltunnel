package tunnel

import (
	"eltunnel/server/protocol"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleTunnel(manager *Manager, w http.ResponseWriter, r *http.Request) {
	// moving connection to Websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("new client connected")

	// msg raw bytes
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(msg)

	// outer content of Message struct, message.RequestID, message.Payload
	var message protocol.Message
	// json string -> go struct
	err = json.Unmarshal(msg, &message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(message.Type)

	// open the inner box - payload to find subdomain
	var registerPayload protocol.RegisterPayload
	err = json.Unmarshal(message.Payload, &registerPayload)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(registerPayload.Subdomain)

	// write subdomain into manager
	manager.Register(registerPayload.Subdomain, conn)

	// assign the url
	registeredPayload := protocol.RegisteredPayload{
		URL: "http://" + registerPayload.Subdomain + ".localhost:8081",
	}
	// converting to raw bytes
	payloadBytes, err := json.Marshal(registeredPayload)
	if err != nil {
		fmt.Println(err)
		return
	}

	// putting the url into Message protocol
	response := protocol.Message{
		Type:      "registered",
		RequestID: "",
		Payload:   payloadBytes,
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}

	// sending down the tunnel to python client
	err = conn.WriteMessage(websocket.TextMessage, responseBytes)
	if err != nil {
		fmt.Println(err)
	}

	for {
		_, msg, err := conn.ReadMessage()
		// client disconnected, remove from manager list, break out of the loop
		if err != nil {
			manager.Remove(registerPayload.Subdomain)
			break
		}
		err = json.Unmarshal(msg, &message)
		if err != nil {
			fmt.Println(err)
			break
		}
		// unmarshal message, payload and send request to waiting proxy goroutine
		var responsePayload protocol.HttpResponsePayload
		err = json.Unmarshal(message.Payload, &responsePayload)
		if err != nil {
			fmt.Println(err)
			fmt.Println("received message type:", message.Type)
			break
		}
		fmt.Println("received response for requestID:", message.RequestID)
		fmt.Println("resolving pending request...")
		manager.ResolvePendingRequest(message.RequestID, responsePayload)
		fmt.Println("resolved!")
	}
}
