# Action Item Extractor

A full-stack web application designed to help you quickly extract actionable items from your notes. It features a modern FastAPI backend with SQLite for data persistence and a minimal Vanilla JS frontend.

The project offers two extraction mechanisms:
1. **Heuristic Extraction:** Extracts items based on standard prefix markers (like `- [ ]`, `TODO:`, `*`, etc.).
2. **LLM Extraction:** Uses a local Large Language Model via Ollama to intelligently parse and extract action items from freeform text.

## Features
- Save text notes and retrieve them.
- Deterministic extraction using regex heuristics.
- AI-powered extraction utilizing structured JSON outputs via local LLMs (`llama3.1`).
- Check off and manage your completed action items.
- Modern FastAPI backend utilizing clean architecture, Pydantic schemas, and dependency injection.

---

## 🚀 Setup and Installation

### Prerequisites
- **Python 3.9+** (Conda environment recommended)
- **Ollama** installed locally (if you plan to use the LLM extraction feature)

### 1. Environment Setup
Create and activate your Python environment (e.g., using `conda`):
```bash
conda create -n cs146s python=3.12
conda activate cs146s
```

Install the required dependencies:
```bash
pip install fastapi uvicorn pydantic ollama pytest python-dotenv
```

### 2. LLM Setup (Optional but recommended)
To use the LLM-powered extraction (`/extract-llm`), ensure the Ollama daemon is running, and pull the required model:
```bash
ollama run llama3.1
```

### 3. Running the Server
Start the development server using Uvicorn from the root directory:
```bash
uvicorn app.main:app --reload
```
The application will automatically initialize the local SQLite database (`data/app.db`) on startup.

Navigate to **`http://127.0.0.1:8000`** in your browser to access the frontend application.

---

## 📡 API Endpoints

The API follows RESTful principles, utilizing Pydantic schemas for data validation.

### Frontend
- **`GET /`**
  - Serves the main `index.html` UI.

### Notes
- **`POST /notes`**
  - Creates and saves a new note.
  - **Body:** `{"content": "string"}`
- **`GET /notes`**
  - Retrieves a list of all saved notes.
- **`GET /notes/{note_id}`**
  - Retrieves a specific note by its ID.

### Action Items
- **`POST /action-items/extract`**
  - Extracts action items using built-in programmatic heuristics.
  - **Body:** `{"text": "string", "save_note": boolean}`
- **`POST /action-items/extract-llm`**
  - Extracts action items intelligently using the local Ollama LLM.
  - **Body:** `{"text": "string", "save_note": boolean}`
- **`GET /action-items`**
  - Retrieves all extracted action items. Accepts an optional `?note_id={id}` query parameter to filter by a specific note.
- **`POST /action-items/{action_item_id}/done`**
  - Toggles the completion status of a specific action item.
  - **Body:** `{"done": boolean}`

---

## 🧪 Running the Test Suite

The project uses `pytest` for unit testing. Tests include mocking the Ollama responses to ensure reliable execution regardless of your local LLM state.

To run the test suite, ensure your environment is activated, then run:

```bash
pytest tests/
```

Alternatively, using conda run:
```bash
conda run -n cs146s pytest tests/
```
