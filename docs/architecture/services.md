# Service Contracts — Greeting

Base paths omit deploy proxy `/api` prefix.

## Read greeting

`GET /v1/greeting`

Success `200`:

```json
{"text":"Hello, World!"}
```

Returns singleton greeting. If row is absent, returns initial value `Hello, World!`.

## Save greeting

`PUT /v1/greeting`

Request:

```json
{"text":"New greeting"}
```

`text` must be string. Server trims leading and trailing whitespace; trimmed value must be non-empty. Unknown fields, wrong types, and invalid JSON fail as malformed request.

Success `200`:

```json
{"text":"New greeting"}
```

Save uses singleton upsert. Last successful save wins.

## Health

`GET /healthz` returns `200` only after migration succeeds and database accepts `SELECT 1`. It has no response contract for product UI.

## Error envelope

All API failures use JSON:

```json
{"error":{"code":"MALFORMED_REQUEST","message":"Request body is malformed."}}
```

| HTTP | Code | Message | Failure kind |
|---|---|---|---|
| 400 | `MALFORMED_REQUEST` | `Request body is malformed.` | Bad JSON, wrong type, unknown field. |
| 422 | `VALIDATION_FAILED` | `Greeting must not be empty.` | Well-formed text violates non-empty-after-trim rule. |
| 500 | `INTERNAL` | `Internal server error.` | Failed database query or unexpected server error. |
| 503 | `UNAVAILABLE` | `Service unavailable.` | Refused or unavailable database dependency. |
