# Week 6 Write-up
Tip: To preview this markdown file
- On Mac, press `Command (⌘) + Shift + V`
- On Windows/Linux, press `Ctrl + Shift + V`

## Instructions

Fill out all of the `TODO`s in this file.

## Submission Details

Name: **TODO** \
SUNet ID: **TODO** \
Citations: **TODO**

This assignment took me about **TODO** hours to do. 


## Brief findings overview 
### Scan Categories Summary
*   **SAST (Static Application Security Testing):** Semgrep OSS identified **5 blocking findings** within the Python backend codebase (`week6/backend/`). These findings spanned critical risk categories including Cross-Origin Resource Sharing (CORS) misconfigurations, SQL Injection (SQLi), Remote Code Execution (RCE), Command Injection, and Server-Side Request Forgery (SSRF) / Arbitrary File Read risks.
*   **Secrets:** No hardcoded secrets or exposed environment credentials were recovered during this local repository scan.
*   **SCA (Software Composition Analysis):** Because the local scan was restricted to the open-source engine without a platform login, formal reachable dependency tracking was deferred; tracking focused heavily on first-party code analysis.

## Fix #1 Insecure CORS Policy
*   **File and Line(s):** `week6/backend/app/main.py` (Line 24)
*   **Rule/Category:** `python.fastapi.security.wildcard-cors.wildcard-cors`
*   **Risk Description:** Enforcing `allow_origins=["*"]` alongside `allow_credentials=True` permits any external domain to make authenticated, cross-origin requests to the backend API using the victim browser's active session state.
*   **Mitigation Change:** 

```python
# BEFORE
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# AFTER
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:8000", "[http://127.0.0.1:8000](http://127.0.0.1:8000)"], 
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

## Fix 2: SQL Injection (SQLi)
File and Line(s): week6/backend/app/routers/notes.py (Lines 71–79)

Rule/Category: python.sqlalchemy.security.audit.avoid-sqlalchemy-text.avoid-sqlalchemy-text

Risk Description: Utilizing Python f-strings to pass input parameter q straight into a raw text() SQL statement allows malicious string breakouts, giving attackers a path to run arbitrary backend database queries.

Mitigation Change:
```python
# BEFORE
@router.get("/unsafe-search", response_model=list[NoteRead])
def unsafe_search(q: str, db: Session = Depends(get_db)) -> list[NoteRead]:
    sql = text(
        f"""
        SELECT id, title, content, created_at, updated_at
        FROM notes
        WHERE title LIKE '%{q}%' OR content LIKE '%{q}%'
        ORDER BY created_at DESC
        LIMIT 50
        """
    )
    rows = db.execute(sql).all()

# AFTER
@router.get("/unsafe-search", response_model=list[NoteRead])
def unsafe_search(q: str, db: Session = Depends(get_db)) -> list[NoteRead]:
    sql = text(
        """
        SELECT id, title, content, created_at, updated_at
        FROM notes
        WHERE title LIKE :q OR content LIKE :q
        ORDER BY created_at DESC
        LIMIT 50
        """
    )
    rows = db.execute(sql, {"q": f"%{q}%"}).all()
```
## Fix 3: OS Command Injection via Subprocess
File and Line(s): week6/backend/app/routers/notes.py (Line 112)

Rule/Category: python.lang.security.audit.subprocess-shell-true.subprocess-shell-true

Risk Description: Running subprocess.run() with shell=True forces arguments to run through an intermediary shell interpreter (like /bin/sh). This makes it easy for meta-characters (like ;, |, or &&) to stack unauthorized shell command executions.

Mitigation Change:
```python
# BEFORE
@router.get("/debug/run")
def debug_run(cmd: str) -> dict[str, str]:
    import subprocess

    completed = subprocess.run(cmd, shell=True, capture_output=True, text=True)  # noqa: S602,S603
    return {"returncode": str(completed.returncode), "stdout": completed.stdout, "stderr": completed.stderr}

# AFTER
@router.get("/debug/run")
def debug_run(cmd: str) -> dict[str, str]:
    import subprocess
    import shlex

    safe_cmd = shlex.split(cmd)
    completed = subprocess.run(safe_cmd, shell=False, capture_output=True, text=True) 
    return {"returncode": str(completed.returncode), "stdout": completed.stdout, "stderr": completed.stderr}
```