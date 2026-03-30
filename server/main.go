package main

import (
	"eltunnel/server/proxy"
	"eltunnel/server/tunnel"
	"fmt"
	"net/http"
)

// proxy - with internet users (browser, curl)
// tunnel - client (python)
func main() {
	manager := tunnel.NewManager()
	http.HandleFunc("/tunnel", func(w http.ResponseWriter, r *http.Request) {
		tunnel.HandleTunnel(manager, w, r)
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.HandleRequest(manager, w, r)
	})

	fmt.Println("proxy started at 8080")
	go http.ListenAndServe(":8080", mux)

	fmt.Println("tunnel started at 8081")
	http.ListenAndServe(":8081", nil)

}
