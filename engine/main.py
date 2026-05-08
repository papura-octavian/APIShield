from fastapi import FastAPI
from pydantic import BaseModel
from typing import Literal

app = FastAPI()

class CheckRequest(BaseModel):
    method : str
    path : str
    headers : dict[str, str]
    body : str

class CheckResponse(BaseModel):
    action : Literal["allow", "block"]
    reason : str | None = None

@app.get("/")
async def root():
    return {"message": "hello from engine"}

@app.post("/check")
async def check(req: CheckRequest) -> CheckResponse :
    if "/admin" in req.path:
        return CheckResponse(action="block", reason="admin path access denied")

    return CheckResponse(action="allow")