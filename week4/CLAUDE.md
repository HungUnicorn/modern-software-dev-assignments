# Week 4: Developer Command Center - Guidance

## Build/Run Commands
- Install: `pip install -r backend/requirements.txt`
- Run App: `make run` (from week4 directory)
- Run Tests: `make test`
- Formatting: `make format`
- Linting: `make lint`

## Architecture & Code Style
- **Clean Code:** Use meaningful naming and small functions.
- **Dependency Inversion:** Routers should depend on abstractions/dependencies, not direct DB instances.
- **FastAPI:** Use Pydantic models for request/response validation. 
- **Database:** SQLite is used. Logic should remain in `backend/app/`.

## Testing Gating
- Before finishing any task, run `make test` and `make lint`. 
- If tests fail, diagnose and fix the implementation immediately.
