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

## Story extension — Persist editable greeting

`greetings` satisfies this story without new entities, columns, foreign keys, or indexes. Singleton key lookup and upsert use primary key `id`; no extra index serves a story query.

### Mock shape review

UI mock at `code/frontend/lib/mock/persist-editable-greeting.ts` exports `Greeting` as `{ "text": string }` and reads/saves one value. This is sound and matches both endpoint success bodies. Its localStorage persistence is UI-stage-only and must be removed when backend client replaces mock; no frontend response-shape change is needed.

### Migration plan

**Forward:** create `greetings` with `id boolean primary key check (id)`, `text text not null check (btrim(text) <> '')`; insert `(true, 'Hello, World!')` with `ON CONFLICT (id) DO NOTHING`. Safe on a populated database: table is new, and conflict-safe seed never overwrites an existing greeting.

**Backward:** drop `greetings`. This deletes current greeting and is only safe before production data needs retention; restore requires backup or reapplying forward migration, which seeds initial value.
