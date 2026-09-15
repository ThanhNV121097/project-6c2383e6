# Architecture Overview — Hello World Acceptance 6

## Stack

| Part | Choice | Reason |
|---|---|---|
| Frontend | Next.js 15 App Router, TypeScript, Tailwind v3, ESLint | Required UI stack; App Router matches container build. |
| Backend | Go 1.22 HTTP server | Small API with standard `net/http`; no router dependency. |
| Database | PostgreSQL 16 | Required durable shared greeting. |
| Database driver | `pgx/v5` stdlib adapter | PostgreSQL connectivity without ORM. |

## Layout

```text
code/backend/cmd/api/main.go     server, migration boot, health check
code/backend/migrations/         ordered SQL migration pairs and embedded files
code/frontend/app/               composition root, layout, shared tokens
code/frontend/components/        one story component per feature
code/frontend/lib/               API clients; mock files removed when API lands
```

## Contracts and conventions

- Backend starts by requiring `DATABASE_URL`, applies SQL migrations in filename order, then listens on `PORT`, `APP_PORT`, or `8080`.
- `GET /healthz` succeeds only after migrations and database `SELECT 1` succeed.
- Public API paths use `/v1/...`; deploy proxy owns `/api` prefix removal.
- Migration history lives in `schema_migrations`; each migration has `.up.sql` and `.down.sql` files. Current migration seeds one singleton greeting row.
- Frontend `app/page.tsx` remains Server Component composition root. Story interaction belongs in default-exported client component.
- Shared visual values belong in `app/globals.css` tokens. CSS modules must use those tokens without fallback values.
- Input validation occurs at API boundary; queries use parameters. No authentication, external services, or user data.

## Design decisions

| Decision | Rejected | Tradeoff |
|---|---|---|
| One singleton database row | Greeting history or multi-row model | Meets single shared greeting scope; migration needed for future history. |
| Standard library HTTP | Third-party router | Fewer dependencies; route count remains tiny. |
| SQL migrations on server boot | Manual migration command | Empty runtime DB becomes usable automatically; startup waits for DB. |
| CSS tokens and modules | Inline styles or Tailwind utilities | Token gate remains enforceable; more CSS declarations. |

## Environment

| Service | Keys |
|---|---|
| Backend | `DATABASE_URL`, `PORT`, optional `APP_PORT` fallback |
| Frontend | `NEXT_PUBLIC_API_URL`, `API_ORIGIN` for server-side requests |
| Compose | `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, optional port and memory settings |

Copy relevant `.env.example` file before local overrides. Never commit `.env` files.

## Run and verify

Run `docker compose --profile local up --build` from repository root. Frontend serves `http://localhost:3000`; backend health check is `http://localhost:8080/healthz`. Compose waits for PostgreSQL before backend and backend health before frontend.

## Compatibility and rollout

Migrations are append-only and filename ordered. `CREATE INDEX CONCURRENTLY` must run outside transaction when later added. Current singleton schema has no external compatibility concern. Unknown: deployment proxy routing is managed outside this repository; service paths intentionally exclude `/api`.
