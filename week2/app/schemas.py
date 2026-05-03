from typing import List, Optional
from pydantic import BaseModel, Field

class NoteCreate(BaseModel):
    content: str = Field(..., min_length=1, description="The content of the note")

class NoteRead(BaseModel):
    id: int
    content: str
    created_at: str

class ActionItemRead(BaseModel):
    id: int
    note_id: Optional[int]
    text: str
    done: bool
    created_at: str

class ExtractRequest(BaseModel):
    text: str = Field(..., min_length=1, description="Text to extract action items from")
    save_note: bool = False

class ExtractedItem(BaseModel):
    id: int
    text: str

class ExtractResponse(BaseModel):
    note_id: Optional[int]
    items: List[ExtractedItem]

class MarkDoneRequest(BaseModel):
    done: bool = True

class MarkDoneResponse(BaseModel):
    id: int
    done: bool
