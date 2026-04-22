# Chirpy

Chirpy is a small Go HTTP API for a Twitter-like microblogging service. It supports user registration, login, JWT-based auth, refresh tokens, chirp creation and deletion, and a webhook that upgrades a user to `Chirpy Red`.

## Stack

- Go
- PostgreSQL
- `sqlc`-generated database layer
- JWT access tokens
- Argon2 password hashing

## Features

- Create user accounts
- Log in with email and password
- Issue access tokens and refresh tokens
- Create, list, fetch, and delete chirps
- Filter chirps by author and sort by creation time
- Update a user's email and password
- Revoke refresh tokens
- Upgrade users through a Polka webhook
- Basic admin metrics and reset endpoints

## Project Structure

```text
.
├── main.go                     # HTTP server and route registration
├── handler_*.go                # Route handlers
├── internal/auth/              # JWT, API key, password, and token helpers
├── internal/database/          # sqlc-generated queries and models
├── sql/schema/                 # Database migrations
├── sql/queries/                # SQL used by sqlc
├── assets/                     # Static assets
└── index.html                  # Static app entry served from /app/
```

## Requirements

- Go `1.25+`
- PostgreSQL

## Environment Variables

Create a `.env` file in the project root:

```env
DB_URL="postgres://username:password@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
TOKEN_SECRET="replace-with-a-long-random-secret"
POLKA_KEY="replace-with-your-polka-api-key"
```

Notes:

- `DB_URL` is required or the server exits on startup.
- `PLATFORM=dev` enables `POST /admin/reset`.
- `TOKEN_SECRET` is used to sign JWT access tokens.
- `POLKA_KEY` is used to authenticate `POST /api/polka/webhooks`.

## Database Setup

Create a PostgreSQL database named `chirpy`, then apply the SQL files in `sql/schema/` in order:

1. `001_users.sql`
2. `002_chirps.sql`
3. `003_users.sql`
4. `004_tokens.sql`
5. `005_users.sql`

Example:

```sh
createdb chirpy
psql "$DB_URL" -f sql/schema/001_users.sql
psql "$DB_URL" -f sql/schema/002_chirps.sql
psql "$DB_URL" -f sql/schema/003_users.sql
psql "$DB_URL" -f sql/schema/004_tokens.sql
psql "$DB_URL" -f sql/schema/005_users.sql
```

## Run Locally

```sh
go run .
```

The server listens on `http://localhost:8080`.

Static files are served from:

- `GET /app/`

## API Overview

### Health and Admin

- `GET /api/healthz`
- `GET /admin/metrics`
- `POST /admin/reset`

### Users and Auth

- `POST /api/users` - create a user
- `POST /api/login` - log in and receive access/refresh tokens
- `POST /api/refresh` - exchange a refresh token for a new access token
- `POST /api/revoke` - revoke a refresh token
- `PUT /api/users` - update email and password for the authenticated user

### Chirps

- `POST /api/chirps` - create a chirp
- `GET /api/chirps` - list chirps
- `GET /api/chirps/{chirpID}` - fetch one chirp
- `DELETE /api/chirps/{chirpID}` - delete a chirp owned by the authenticated user

Supported query params for `GET /api/chirps`:

- `author_id=<uuid>` filters chirps by user
- `sort=asc|desc` sorts by `created_at`

### Webhooks

- `POST /api/polka/webhooks` - accepts Polka events authenticated with `Authorization: ApiKey <key>`

Only the `user.upgraded` event changes state, setting `is_chirpy_red` to `true` for the user in the payload.

## Auth Notes

- Access tokens are sent as `Authorization: Bearer <token>`.
- Refresh tokens are also sent in the `Authorization` header for refresh and revoke requests.
- Access tokens expire after 1 hour.

## Example Requests

Create a user:

```sh
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}'
```

Log in:

```sh
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}'
```

Create a chirp:

```sh
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access-token>" \
  -d '{"body":"hello world"}'
```

List chirps for one author in descending order:

```sh
curl "http://localhost:8080/api/chirps?author_id=<user-uuid>&sort=desc"
```

## Testing

Run the test suite with:

```sh
go test ./...
```
