from backend.app.services.extract import extract_action_items, extract_tags


def test_extract_action_items():
    text = """
    This is a note
    - TODO: write tests
    - Ship it!
    Not actionable
    """.strip()
    items = extract_action_items(text)
    assert "TODO: write tests" in items
    assert "Ship it!" in items


def test_extract_tags_basic():
    tags = extract_tags("Working on #backend and #api today")
    assert tags == ["backend", "api"]


def test_extract_tags_empty():
    assert extract_tags("no tags here") == []


def test_extract_tags_deduplication_not_applied():
    tags = extract_tags("#foo and #foo again")
    assert tags.count("foo") == 2


def test_extract_note_endpoint(client):
    r = client.post(
        "/notes/", json={"title": "Sprint", "content": "- TODO: write docs\n#backend #api"}
    )
    note_id = r.json()["id"]

    r = client.post(f"/notes/{note_id}/extract")
    assert r.status_code == 200
    data = r.json()
    assert "backend" in data["tags"]
    assert "api" in data["tags"]
    assert any("write docs" in item["description"] for item in data["action_items"])


def test_extract_note_endpoint_404(client):
    r = client.post("/notes/99999/extract")
    assert r.status_code == 404
