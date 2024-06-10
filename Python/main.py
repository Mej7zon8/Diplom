import asyncio
import threading
from concurrent.futures import ThreadPoolExecutor
from fastapi import FastAPI
from pydantic import BaseModel
from starlette.requests import Request
from starlette.responses import JSONResponse

from logic import processor

app = FastAPI()
executor = ThreadPoolExecutor(max_workers=2)
thread_local = threading.local()


def run_in_executor(func, *args):
    loop = asyncio.get_event_loop()
    return loop.run_in_executor(executor, func, *args)


@app.exception_handler(Exception)
async def exception_handler(request, exc):
    return JSONResponse(
        status_code=500,
        content={"detail": str(exc)},
    )


@app.get("/")
async def root():
    return {"message": "Still Alive"}


"""
    Text categorization endpoint
"""


class CategorizeRequest(BaseModel):
    text: str = ""
    candidate_labels: list[str] = []


@app.post("/categorize")
async def categorize(request: Request):
    content = CategorizeRequest.model_validate_json(await request.body())
    return {
        "result": await run_in_executor(run_process_task, content)
    }


def run_process_task(request: CategorizeRequest):
    if not hasattr(thread_local, 'processor'):
        thread_local.processor = processor.Processor()
    return thread_local.processor.predict(request.text, request.candidate_labels)
