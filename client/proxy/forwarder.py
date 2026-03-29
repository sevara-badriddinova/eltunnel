import aiohttp, asyncio
from protocol.messages import HttpRequestPayload, HttpResponsePayload

async def forward_request(port, request):
    async with aiohttp.ClientSession() as session:
        async with session.request(request.method, f"http://localhost:{port}{request.path}") as resp:
            body = await resp.text()
            status = resp.status
            headers = dict(resp.headers)
            print("got response:", status, body[:50])  # add this
            return HttpResponsePayload(status, headers, body)