# Service Contracts — Greeting

Base paths omit deploy proxy `/api` prefix.

## Read greeting

`GET /v1/greeting`

Auth: none.

Success `200`:

```json
{"text":"Hello, World!"}
```

Returns singleton greeting. If row is absent, returns initial value `Hello, World!`.

## Save greeting

`PUT /v1/greeting`

Auth: none.

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

## Story extension — Persist editable greeting

No endpoint changes needed. `GET /v1/greeting` and `PUT /v1/greeting` are complete contract for UI mock shape `{ "text": string }`.

| Endpoint | Status codes | Error codes |
|---|---|---|
| `GET /v1/greeting` | `200`, `500`, `503` | `INTERNAL`, `UNAVAILABLE` |
| `PUT /v1/greeting` | `200`, `400`, `422`, `500`, `503` | `MALFORMED_REQUEST`, `VALIDATION_FAILED`, `INTERNAL`, `UNAVAILABLE` |

`INTERNAL` covers a query that fails after database dependency is reached. `UNAVAILABLE` covers connection refusal or unavailable database dependency only. Backend replaces localStorage mock with this contract; no UI body-shape change.
