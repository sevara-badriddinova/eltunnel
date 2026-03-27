from dataclasses import dataclass

@dataclass
class Message:
    type: str
    request_id: str
    payload: dict

@dataclass
class RegisterPayload:
    subdomain: str

@dataclass
class RegisteredPayload:
    url: str

@dataclass
class HttpRequestPayload:
    method: str
    path: str
    headers: dict
    body: str

@dataclass
class HttpResponsePayload:
    status_code: int
    headers: dict
    body: str

@dataclass
class ErrorPayload:
    err: str
