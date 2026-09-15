# Story — Persist editable greeting

Module: `greeting`
Plan item: Persist editable greeting

## User story

As a Visitor, I want to view and save the shared greeting, so that the page reflects the greeting persisted in PostgreSQL through the Go API after reload.

## In scope

- Read the current singleton greeting from the backend API backed by PostgreSQL.
- Seed and show `Hello, World!` when no greeting has been changed since first product run.
- Render one centered greeting screen matching the approved design: large heading, visually hidden `Greeting` label, text input, `Save` button, and status line.
- Keep input value in sync with current greeting on load and after save.
- Save non-empty greetings after trimming leading and trailing whitespace.
- Persist saved greeting so reload shows the saved value in heading and input.
- Validate empty or whitespace-only input on the client-visible path: do not save, show `Enter a greeting.`, and return focus to input.
- Preserve saved greeting capitalization and punctuation exactly, except trimming at save boundary.
- Allow all visitors to read and save with no sign-in.
- Use last successful save as the current shared greeting when multiple visitors save.

## Out of scope

- Sign-in, roles, per-visitor greetings, and permissions beyond anonymous Visitor access.
- Navigation, extra sections, secondary pages, greeting history, deletion, timestamps, and author metadata.
- External services.
- Dedicated loading screens and API error screens; approved design has no such states, and API error envelope belongs to TL service contract.
- Multiple greetings or list management; this story owns one global singleton greeting only.
- Styling beyond approved minimal design tokens and component states in `design/design-system.md`.

## UI scope

- Screen: Greeting page only.
- Section: `main > section.greeting-section`, labelled by `h1#greeting`.
- Default state: heading shows current greeting; form below contains visually hidden label text `Greeting`, text input prefilled with current greeting, primary `Save` button, and reserved `p#status` live region.
- Validation state: whitespace-only submit keeps stored greeting unchanged, focuses greeting input, and shows status text `Enter a greeting.`
- Success state: successful save updates heading and input to saved greeting and shows status text `Saved.`
- Responsive scope: page works from 320px wide with no horizontal page scroll; form stacks into one column up to 520px.
- Accessibility scope: associated input label, keyboard-reachable controls, visible focus outline, and polite live status updates.

## Acceptance criteria

- SC-1 [GREETING-001 AC-1]: Given no greeting has been changed since first product run, when Visitor opens the page, heading text is `Hello, World!`.
- SC-2 [GREETING-001 AC-2]: Given stored greeting is `Pipeline OK`, when Visitor opens the page, heading text is `Pipeline OK`.
- SC-3 [GREETING-001 AC-3]: Given stored greeting contains capitalization and punctuation `Hello, Database!`, when Visitor opens the page, heading text is exactly `Hello, Database!`.
- SC-4 [GREETING-002 AC-1]: Given stored greeting is `Hello, World!`, when Visitor opens the page, text input value is `Hello, World!`.
- SC-5 [GREETING-002 AC-2]: Given page is open, when Visitor inspects form controls, text input has accessible label `Greeting`.
- SC-6 [GREETING-002 AC-2]: Given page is open, when Visitor inspects form controls, visible button text is `Save`.
- SC-7 [GREETING-002 AC-3]: Given page is open and input value is `New greeting`, when Visitor saves, heading text becomes `New greeting`.
- SC-8 [GREETING-002 AC-4]: Given page is open and input value is `  Trim me  `, when Visitor saves, heading text becomes `Trim me`.
- SC-9 [GREETING-002 AC-5]: Given page is open and input value is empty or whitespace only, when Visitor saves, status text becomes `Enter a greeting.`
- SC-10 [GREETING-002 AC-6]: Given page is open and input value is empty or whitespace only, when Visitor saves, focus is on the greeting input.
- SC-11 [GREETING-002 AC-7]: Given page is open and input value is empty or whitespace only, when Visitor saves, stored greeting remains unchanged.
- SC-12 [GREETING-002 AC-8]: Given page is open and input value is `Saved value`, when Visitor saves, status text becomes `Saved.`
- SC-13 [GREETING-002 AC-9]: Given greeting `Saved value` was saved, when Visitor reloads the page, heading text is `Saved value`.
- SC-14 [GREETING-002 AC-10]: Given greeting `Saved value` was saved, when Visitor reloads the page, text input value is `Saved value`.
- SC-15 [GREETING-001 boundary]: Given stored greeting is long text, when Visitor opens the page at 320px width or wider, heading wraps within the centered section and page has no horizontal scroll.
- SC-16 [GREETING-002 boundary]: Given page is open and input value is long non-empty text, when Visitor saves, save is accepted and heading wraps within the centered section with no horizontal scroll at 320px width or wider.
- SC-17 [GREETING-002 permission]: Given Visitor is not signed in, when Visitor saves a non-empty greeting, save succeeds with no authentication prompt.
- SC-18 [GREETING-002 conflict]: Given two visitors save different greetings, when both saves complete, last successful save is current greeting after reload.

## Dependencies

- Backend API from `docs/architecture/overview.md`: public `/v1/...` endpoint serving current greeting and accepting updates.
- PostgreSQL singleton greeting row seeded through migrations with initial value `Hello, World!`.
- Frontend scaffold using Next.js App Router, Server Component `app/page.tsx`, and one client story component for form interaction.
- Design tokens in `app/globals.css` matching `design/design-system.md`.
- No external accounts, credentials, sign-in providers, or third-party services.

## Notes for later stages

- Test cases must cover every `SC-n` above in `docs/greeting/test-cases/persist-editable-greeting.md`.
- UI implementation must use mock data first, then backend story replaces it with API client.
- Backend validation must trim input and reject empty saved value at API boundary.
