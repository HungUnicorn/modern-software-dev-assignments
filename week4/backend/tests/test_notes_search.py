def test_search_matches_title(client):
    client.post("/notes/", json={"title": "Meeting Notes", "content": "Discuss roadmap"})
    r = client.get("/notes/search/", params={"q": "meeting"})
    assert r.status_code == 200
    results = r.json()
    assert any("Meeting" in n["title"] for n in results)


def test_search_matches_content(client):
    client.post("/notes/", json={"title": "Sprint", "content": "Finish the Dashboard feature"})
    r = client.get("/notes/search/", params={"q": "dashboard"})
    assert r.status_code == 200
    results = r.json()
    assert any("Dashboard" in n["content"] for n in results)


def test_search_case_insensitive(client):
    client.post("/notes/", json={"title": "CamelCase Title", "content": "some content"})
    r = client.get("/notes/search/", params={"q": "camelcase"})
    assert r.status_code == 200
    results = r.json()
    assert any("CamelCase" in n["title"] for n in results)


def test_search_no_match_returns_empty(client):
    r = client.get("/notes/search/", params={"q": "xyznonexistent"})
    assert r.status_code == 200
    assert r.json() == []


def test_search_no_query_returns_all(client):
    client.post("/notes/", json={"title": "Alpha", "content": "first"})
    client.post("/notes/", json={"title": "Beta", "content": "second"})
    r = client.get("/notes/search/")
    assert r.status_code == 200
    assert len(r.json()) >= 2
