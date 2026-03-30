package protocol

import "encoding/json"

type Message struct {
	Type      string
	RequestID string
	Payload   json.RawMessage
}

type RegisterPayload struct {
	Subdomain string
}

type RegisteredPayload struct {
	URL string
}

type HttpRequestPayload struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    string
}

type HttpResponsePayload struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

type ErrorPayload struct {
	Error string
}
