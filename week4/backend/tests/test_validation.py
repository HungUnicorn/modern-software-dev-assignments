import pytest


@pytest.mark.parametrize(
    "payload",
    [
        {"title": "", "content": "valid"},
        {"title": "valid", "content": ""},
        {"title": "", "content": ""},
        {"title": "x" * 201, "content": "valid"},
    ],
)
def test_create_note_invalid(client, payload):
    r = client.post("/notes/", json=payload)
    assert r.status_code == 422


@pytest.mark.parametrize(
    "payload",
    [
        {"title": "", "content": "valid"},
        {"title": "valid", "content": ""},
        {"title": "x" * 201, "content": "valid"},
    ],
)
def test_update_note_invalid(client, payload):
    r = client.post("/notes/", json={"title": "Seed", "content": "Seed"})
    note_id = r.json()["id"]
    r = client.put(f"/notes/{note_id}", json=payload)
    assert r.status_code == 422


def test_create_action_item_empty_description(client):
    r = client.post("/action-items/", json={"description": ""})
    assert r.status_code == 422


def test_create_note_missing_fields(client):
    r = client.post("/notes/", json={"title": "only title"})
    assert r.status_code == 422

    r = client.post("/notes/", json={"content": "only content"})
    assert r.status_code == 422
