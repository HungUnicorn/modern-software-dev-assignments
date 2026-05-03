from __future__ import annotations

import sqlite3
from typing import List

from fastapi import APIRouter, HTTPException, Depends

from ..db import get_db_connection, insert_note, get_note, list_notes
from ..schemas import NoteCreate, NoteRead

router = APIRouter(prefix="/notes", tags=["notes"])


@router.post("", response_model=NoteRead)
def create_note(
    payload: NoteCreate,
    conn: sqlite3.Connection = Depends(get_db_connection)
) -> NoteRead:
    note_id = insert_note(conn, payload.content)
    row = get_note(conn, note_id)
    if not row:
        raise HTTPException(status_code=500, detail="Failed to create note")
    return NoteRead(id=row["id"], content=row["content"], created_at=row["created_at"])


@router.get("/{note_id}", response_model=NoteRead)
def get_single_note(
    note_id: int,
    conn: sqlite3.Connection = Depends(get_db_connection)
) -> NoteRead:
    row = get_note(conn, note_id)
    if row is None:
        raise HTTPException(status_code=404, detail="Note not found")
    return NoteRead(id=row["id"], content=row["content"], created_at=row["created_at"])


@router.get("", response_model=List[NoteRead])
def get_all_notes(conn: sqlite3.Connection = Depends(get_db_connection)) -> List[NoteRead]:
    rows = list_notes(conn)
    return [NoteRead(id=row["id"], content=row["content"], created_at=row["created_at"]) for row in rows]
