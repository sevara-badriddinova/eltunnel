import asyncio
import json
import websockets
import string
import random

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

def random_subdomain():
    res = random.choices(string.ascii_lowercase, k = 5)
    return ''.join(res)
    
if __name__ == "__main__":
    asyncio.run(connect("ws://localhost:8081/tunnel", 3000))