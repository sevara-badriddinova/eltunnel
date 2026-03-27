import asyncio
import json
import websockets

async def connect(server_url, port):
    message = {
        "type": "register",
        "request_id": "",
        "payload": {
            "subdomain": "alice"
        }
    }
    message_json = json.dumps(message)
    async with websockets.connect(server_url) as ws:
        await ws.send(message_json)
        response = await ws.recv()
        response = json.loads(response)
        print(response)

if __name__ == "__main__":
    asyncio.run(connect("ws://localhost:8081/tunnel", 3000))