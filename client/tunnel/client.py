import asyncio
import json
import websockets
import string
import random
from proxy.forwarder import forward_request
from protocol.messages import HttpRequestPayload

async def connect(server_url, port):
    message = {
        "type": "register",
        "request_id": "",
        "payload": {
            "subdomain": random_subdomain()
        }
    }
    message_json = json.dumps(message)
    async with websockets.connect(server_url) as ws:
        await ws.send(message_json)
        response = await ws.recv()
        response = json.loads(response)
        print("Tunnel ready at response " + response["Payload"]["URL"])
        while True:
            try:
                msg = await ws.recv()
                msg = json.loads(msg)
                print("full message:", msg) 
                payload = {k.lower(): v for k, v in msg["Payload"].items()}
                request = HttpRequestPayload(**payload)
                response = await forward_request(port, request)
                print("sending response back...")
                await ws.send(json.dumps({
                    "Type": "http_response",
                    "RequestID": msg["RequestID"],
                    "Payload": {
                        "StatusCode": response.status_code,
                        "Headers": response.headers,
                        "Body": response.body
                    }
                }))
                print("response sent!")
            except Exception as e:
                print("error:", e)

def random_subdomain():
    res = random.choices(string.ascii_lowercase, k = 5)
    return ''.join(res)
    
if __name__ == "__main__":
    asyncio.run(connect("ws://localhost:8081/tunnel", 3000))