# ERD — Greeting

## `greetings`

One row holds globally shared current greeting.

| Column | PostgreSQL type | Constraints | Purpose |
|---|---|---|---|
| `id` | `boolean` | primary key, check (`id`) | Enforces singleton row using value `true`. |
| `text` | `text` | not null, check (`btrim(text) <> ''`) | Current saved greeting. |

No relationships. No personal data, history, or timestamps are in scope.

## Seed and lifecycle

Migration inserts `(true, 'Hello, World!')` with conflict ignored. Read endpoint treats absent row as initial greeting; save upserts singleton row. Last successful write wins.
