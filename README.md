# ai-doc-platform 

## Tech Stack

Go 1.25 (chi) · Python 3.13 (FastAPI, uvicorn) · PostgreSQL 16 · Docker Compose

## Prerequisites

- Docker and Docker Compose
- A `.env` file in the repo root. Compose substitutes these values into `docker-compose.yml` at startup. Required keys:
  `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `DATABASE_URL_DOCKER`, `INFERENCE_URL_DOCKER`.

## Running the stack

```bash
docker compose up -d --build
docker compose ps        # postgres + inference-service show (healthy), gateway is Up
```

```bash
echo "placeholder" > invoice.txt          # filename "invoice" → classified as finance
```

```bash
curl -s -F "file=@invoice.txt" http://localhost:8080/upload
#   -> {"id":"<uuid>","filename":"invoice.txt"}

curl -s http://localhost:8080/documents/<uuid>
#   -> {"document":{...,"status":"done"},
#       "prediction":{"label":"finance","confidence":0.8}}
```

## API

| Method | Path              | Purpose                               |
|--------|-------------------|---------------------------------------|
| GET    | `/health`         | Gateway liveness                      |
| POST   | `/upload`         | Upload a document (multipart `file`)  |
| GET    | `/documents/{id}` | Fetch a document + its prediction     |

## Local development (without Docker)

**Gateway:**
```bash
cd gateway-service
go run ./cmd/server/main.go        # needs a reachable Postgres + inference
go test ./...                      # unit tests (fast, no Docker)
go test -tags=integration ./...    # repository integration tests (needs Docker)
```
**Inference:**
```bash
cd inference-service
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8001
pytest
```

## Teardown

```bash
docker compose down       # stop and remove containers
docker compose down -v    # also wipe postgres + uploads volumes (clean slate)
```
