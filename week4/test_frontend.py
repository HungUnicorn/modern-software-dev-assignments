import json
import urllib.request

req = urllib.request.Request("http://localhost:8000/notes/")
with urllib.request.urlopen(req) as response:
    notes = json.loads(response.read().decode())
    print("Notes:", notes)

if notes:
    n = notes[0]
    print(f"Testing PUT /notes/{n['id']}")
    data = json.dumps({"title": "Updated", "content": "Updated content"}).encode()
    req = urllib.request.Request(f"http://localhost:8000/notes/{n['id']}", data=data, method="PUT")
    req.add_header("Content-Type", "application/json")
    try:
        with urllib.request.urlopen(req) as response:
            print("PUT response:", response.status, response.read().decode())
    except Exception as e:
        print("PUT error:", e)

    print(f"Testing DELETE /notes/{n['id']}")
    req = urllib.request.Request(f"http://localhost:8000/notes/{n['id']}", method="DELETE")
    try:
        with urllib.request.urlopen(req) as response:
            print("DELETE response:", response.status, response.read().decode())
    except Exception as e:
        print("DELETE error:", e)
