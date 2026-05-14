# API Documentation

This document describes the available API endpoints and their corresponding payloads.

## General Endpoints

### Root
- **GET** `/`
- **Description:** Root endpoint.
- **Response:** `{}`

### Ping
- **GET** `/ping`
- **Description:** Ping endpoint to check server health.
- **Response:** Object

---

## Notes Endpoints

### List Notes
- **GET** `/notes/`
- **Description:** Retrieves a list of all notes.
- **Response:** Array of `NoteRead` objects.
```json
[
  {
    "id": 1,
    "title": "Meeting Notes",
    "content": "Discussed the new feature."
  }
]
```

### Create Note
- **POST** `/notes/`
- **Description:** Creates a new note.
- **Payload:** `NoteCreate`
```json
{
  "title": "Meeting Notes",
  "content": "Discussed the new feature."
}
```
- **Response:** The created `NoteRead` object.

### Search Notes
- **GET** `/notes/search/`
- **Description:** Searches notes by query.
- **Query Parameters:**
  - `q` (string, optional): The search query.
- **Response:** Array of matching `NoteRead` objects.

### Get Note
- **GET** `/notes/{note_id}`
- **Description:** Retrieves a specific note by ID.
- **Response:** `NoteRead` object.

### Update Note
- **PUT** `/notes/{note_id}`
- **Description:** Updates an existing note.
- **Payload:** `NoteUpdate`
```json
{
  "title": "Updated Meeting Notes",
  "content": "Discussed the new feature and timeline."
}
```
- **Response:** The updated `NoteRead` object.

### Delete Note
- **DELETE** `/notes/{note_id}`
- **Description:** Deletes a specific note.
- **Response:** `204 No Content`

### Extract Note Data
- **POST** `/notes/{note_id}/extract`
- **Description:** Extracts tags and action items from a note.
- **Response:** `ExtractResult` object.
```json
{
  "tags": ["meeting", "feature"],
  "action_items": [
    {
      "id": 1,
      "description": "Follow up with design team",
      "completed": false
    }
  ]
}
```

---

## Action Items Endpoints

### List Action Items
- **GET** `/action-items/`
- **Description:** Retrieves a list of all action items.
- **Response:** Array of `ActionItemRead` objects.
```json
[
  {
    "id": 1,
    "description": "Follow up with design team",
    "completed": false
  }
]
```

### Create Action Item
- **POST** `/action-items/`
- **Description:** Creates a new action item.
- **Payload:** `ActionItemCreate`
```json
{
  "description": "Follow up with design team"
}
```
- **Response:** The created `ActionItemRead` object.

### Complete Action Item
- **PUT** `/action-items/{item_id}/complete`
- **Description:** Marks an action item as completed.
- **Response:** The updated `ActionItemRead` object.
