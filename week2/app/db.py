from __future__ import annotations

import sqlite3
from pathlib import Path
from typing import Optional, Generator


BASE_DIR = Path(__file__).resolve().parents[1]
DATA_DIR = BASE_DIR / "data"
DB_PATH = DATA_DIR / "app.db"


def ensure_data_directory_exists() -> None:
    DATA_DIR.mkdir(parents=True, exist_ok=True)


def get_db_connection() -> Generator[sqlite3.Connection, None, None]:
    ensure_data_directory_exists()
    connection = sqlite3.connect(DB_PATH, check_same_thread=False)
    connection.row_factory = sqlite3.Row
    try:
        yield connection
    finally:
        connection.close()


def init_db() -> None:
    ensure_data_directory_exists()
    with sqlite3.connect(DB_PATH) as connection:
        cursor = connection.cursor()
        cursor.execute(
            """
            CREATE TABLE IF NOT EXISTS notes (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                content TEXT NOT NULL,
                created_at TEXT DEFAULT (datetime('now'))
            );
            """
        )
        cursor.execute(
            """
            CREATE TABLE IF NOT EXISTS action_items (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                note_id INTEGER,
                text TEXT NOT NULL,
                done INTEGER DEFAULT 0,
                created_at TEXT DEFAULT (datetime('now')),
                FOREIGN KEY (note_id) REFERENCES notes(id)
            );
            """
        )
        connection.commit()


def insert_note(connection: sqlite3.Connection, content: str) -> int:
    cursor = connection.cursor()
    cursor.execute("INSERT INTO notes (content) VALUES (?)", (content,))
    connection.commit()
    return int(cursor.lastrowid)


def list_notes(connection: sqlite3.Connection) -> list[sqlite3.Row]:
    cursor = connection.cursor()
    cursor.execute("SELECT id, content, created_at FROM notes ORDER BY id DESC")
    return list(cursor.fetchall())


def get_note(connection: sqlite3.Connection, note_id: int) -> Optional[sqlite3.Row]:
    cursor = connection.cursor()
    cursor.execute(
        "SELECT id, content, created_at FROM notes WHERE id = ?",
        (note_id,),
    )
    return cursor.fetchone()


def insert_action_items(connection: sqlite3.Connection, items: list[str], note_id: Optional[int] = None) -> list[int]:
    cursor = connection.cursor()
    ids: list[int] = []
    for item in items:
        cursor.execute(
            "INSERT INTO action_items (note_id, text) VALUES (?, ?)",
            (note_id, item),
        )
        ids.append(int(cursor.lastrowid))
    connection.commit()
    return ids


def list_action_items(connection: sqlite3.Connection, note_id: Optional[int] = None) -> list[sqlite3.Row]:
    cursor = connection.cursor()
    if note_id is None:
        cursor.execute(
            "SELECT id, note_id, text, done, created_at FROM action_items ORDER BY id DESC"
        )
    else:
        cursor.execute(
            "SELECT id, note_id, text, done, created_at FROM action_items WHERE note_id = ? ORDER BY id DESC",
            (note_id,),
        )
    return list(cursor.fetchall())


def mark_action_item_done(connection: sqlite3.Connection, action_item_id: int, done: bool) -> None:
    cursor = connection.cursor()
    cursor.execute(
        "UPDATE action_items SET done = ? WHERE id = ?",
        (1 if done else 0, action_item_id),
    )
    connection.commit()
