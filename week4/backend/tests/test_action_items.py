def test_complete_sets_completed_true(client):
    r = client.post("/action-items/", json={"description": "Write docs"})
    assert r.status_code == 201
    item_id = r.json()["id"]

    r = client.put(f"/action-items/{item_id}/complete")
    assert r.status_code == 200
    data = r.json()
    assert data["completed"] is True
    assert data["id"] == item_id


def test_complete_nonexistent_item_returns_404(client):
    r = client.put("/action-items/99999/complete")
    assert r.status_code == 404


def test_create_and_complete_action_item(client):
    payload = {"description": "Ship it"}
    r = client.post("/action-items/", json=payload)
    assert r.status_code == 201, r.text
    item = r.json()
    assert item["completed"] is False

    r = client.put(f"/action-items/{item['id']}/complete")
    assert r.status_code == 200
    done = r.json()
    assert done["completed"] is True

    r = client.get("/action-items/")
    assert r.status_code == 200
    items = r.json()
    assert len(items) == 1
