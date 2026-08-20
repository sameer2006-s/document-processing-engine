# Document Processing Engine

A small, simple project for uploading documents, running OCR, generating thumbnails and tags, and chatting with the results.

Intuitive end to end: upload a file, Temporal runs the pipeline, then search or chat with what came out.

```
  upload ──▶ Temporal workflow
                  │
       ┌──────────┼──────────┐
       ▼          ▼          ▼
      OCR     thumbnail    tags
       │          │          │
       └──────────┼──────────┘
                  ▼
             search / chat
```

## Features

- Auth (register / login)
- Document upload to MinIO
- OCR, thumbnail, and tagging via Temporal
- Search your files
- Chat with a document (GitHub Models or OpenRouter)

## Stack

| Piece        | Role                       |
|--------------|----------------------------|
| Go + Gin     | API (`cmd/api`)            |
| Temporal     | OCR → thumbnail → tags     |
| Postgres     | users, docs, chat history  |
| MinIO        | object storage             |
| Python OCR   | OCR service                |
| React + Vite | frontend (`web`)           |

## Quick start

**1. Start infra**

```bash
docker compose up -d
```

Starts Postgres, Temporal (+ UI on `:8080`), MinIO (`:9000` / console `:9001`), and OCR (`:8090`).

**2. Env**

Fill in `.env` (DB, MinIO, chat provider keys). Local compose defaults work out of the box for most values.

**3. Run API + worker**

```bash
go run ./cmd/api
go run ./cmd/worker
```

**4. Frontend (dev)**

```bash
cd web
npm install
npm run dev
```

For a production UI build (served by the API from `web/dist`):

```bash
cd web && npm run build
```

## Ports

| Service       | Port |
|---------------|------|
| API           | 7070 |
| Temporal UI   | 8080 |
| MinIO API     | 9000 |
| MinIO console | 9001 |
| OCR           | 8090 |
| Postgres      | 5432 |

## API

```
POST   /register
POST   /login
POST   /upload
GET    /documents
DELETE /documents/:id
GET    /get-file/:id
GET    /documents/:id/thumbnail
GET    /search-my-files
POST   /chat
POST   /chat/stream
```
