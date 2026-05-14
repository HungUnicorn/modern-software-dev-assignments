def test_create_and_list_notes(client):
    payload = {"title": "Test", "content": "Hello world"}
    r = client.post("/notes/", json=payload)
    assert r.status_code == 201, r.text
    data = r.json()
    assert data["title"] == "Test"

    r = client.get("/notes/")
    assert r.status_code == 200
    items = r.json()
    assert len(items) >= 1

    r = client.get("/notes/search/")
    assert r.status_code == 200

    r = client.get("/notes/search/", params={"q": "Hello"})
    assert r.status_code == 200
    items = r.json()
    assert len(items) >= 1


def test_update_note(client):
    r = client.post("/notes/", json={"title": "Old", "content": "Old content"})
    note_id = r.json()["id"]

    r = client.put(f"/notes/{note_id}", json={"title": "New", "content": "New content"})
    assert r.status_code == 200
    data = r.json()
    assert data["title"] == "New"
    assert data["content"] == "New content"
    assert data["id"] == note_id


def test_update_note_404(client):
    r = client.put("/notes/99999", json={"title": "X", "content": "Y"})
    assert r.status_code == 404


def test_delete_note(client):
    r = client.post("/notes/", json={"title": "Bye", "content": "Gone"})
    note_id = r.json()["id"]

    r = client.delete(f"/notes/{note_id}")
    assert r.status_code == 204

    r = client.get(f"/notes/{note_id}")
    assert r.status_code == 404


def test_delete_note_404(client):
    r = client.delete("/notes/99999")
    assert r.status_code == 404
