from __future__ import annotations

import sqlite3
from typing import List, Optional

from fastapi import APIRouter, HTTPException, Depends

from ..db import (
    get_db_connection,
    insert_note,
    insert_action_items,
    list_action_items,
    mark_action_item_done,
)
from ..schemas import (
    ExtractRequest,
    ExtractResponse,
    ActionItemRead,
    MarkDoneRequest,
    MarkDoneResponse,
    ExtractedItem,
)
from ..services.extract import extract_action_items, extract_action_items_llm

router = APIRouter(prefix="/action-items", tags=["action-items"])


@router.post("/extract", response_model=ExtractResponse)
def extract(
    payload: ExtractRequest,
    conn: sqlite3.Connection = Depends(get_db_connection)
) -> ExtractResponse:
    note_id: Optional[int] = None
    if payload.save_note:
        note_id = insert_note(conn, payload.text)

    items = extract_action_items(payload.text)
    ids = insert_action_items(conn, items, note_id=note_id)
    
    extracted_items = [
        ExtractedItem(id=i, text=t)
        for i, t in zip(ids, items)
    ]
    
    return ExtractResponse(note_id=note_id, items=extracted_items)


@router.post("/extract-llm", response_model=ExtractResponse)
def extract_llm(
    payload: ExtractRequest,
    conn: sqlite3.Connection = Depends(get_db_connection)
) -> ExtractResponse:
    note_id: Optional[int] = None
    if payload.save_note:
        note_id = insert_note(conn, payload.text)

    items = extract_action_items_llm(payload.text)
    ids = insert_action_items(conn, items, note_id=note_id)
    
    extracted_items = [
        ExtractedItem(id=i, text=t)
        for i, t in zip(ids, items)
    ]
    
    return ExtractResponse(note_id=note_id, items=extracted_items)


@router.get("", response_model=List[ActionItemRead])
def list_all(
    note_id: Optional[int] = None,
    conn: sqlite3.Connection = Depends(get_db_connection)
) -> List[ActionItemRead]:
    rows = list_action_items(conn, note_id=note_id)
    return [
        ActionItemRead(
            id=r["id"],
            note_id=r["note_id"],
            text=r["text"],
            done=bool(r["done"]),
            created_at=r["created_at"],
        )
        for r in rows
    ]


@router.post("/{action_item_id}/done", response_model=MarkDoneResponse)
def mark_done(
    action_item_id: int,
    payload: MarkDoneRequest,
    conn: sqlite3.Connection = Depends(get_db_connection)
) -> MarkDoneResponse:
    mark_action_item_done(conn, action_item_id, payload.done)
    return MarkDoneResponse(id=action_item_id, done=payload.done)
