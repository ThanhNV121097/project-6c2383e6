# SRS — Greeting

Module: `greeting`
Design: [View the approved design](http://localhost:8080/design/6c2383e6-89b9-4ad7-9d0f-8dccdaa6309b)
Design system: `design/design-system.md`

> One file per module, at `docs/{module}/SRS.md`. It covers only the functions
> that belong to this module. Never write `docs/SRS.md`.

## 1. Purpose

Greeting module lets any visitor read and change the single greeting for "Hello World Acceptance 6". It proves the end-to-end product path: persisted data, backend API, and frontend UI all reflect the same greeting after save and reload. Without this module, the project would not satisfy its acceptance goal.

## 2. Actors

| Actor | Who they are | What they may do in this module |
|---|---|---|
| Visitor | Anyone who opens the page, with no sign-in | View current greeting, edit the greeting input, and save a new greeting |

## 3. Scope

**In scope** — the functions specified below, by their plan titles:

- Persist editable greeting

**Out of scope** — name what a reader would reasonably expect here and say where it lives instead. This section prevents the same argument twice.

- Sign-in and roles — deliberately not built; brief says no sign-in.
- Navigation and extra sections — deliberately not built; approved design has one centered greeting section only.
- External services — deliberately not built; brief says no external services.
- Multiple greetings, history, timestamps, author metadata, and deletion — deliberately not built; plan has one stored greeting and one save action.
- Loading and API error screens — deliberately not built; approved design shows only the default screen. API error envelope belongs in TL service contract.

## 4. Functional requirements

### 4.1 Persist editable greeting

**Requirement GREETING-001 — Show stored greeting**

*As a* Visitor, *I want to* see the current stored greeting as the page heading, *so that* the page reflects persisted data instead of fixed copy.

Behaviour:

1. When Visitor opens the page, the page obtains the current stored greeting.
2. Page shows that greeting in the single `h1` heading.
3. First product run starts with greeting text `Hello, World!`.
4. Heading preserves saved greeting text exactly, including capitalization and punctuation.

**Acceptance criteria** — each is proved by at least one test case in `docs/greeting/test-cases/persist-editable-greeting.md`, through the story plan that cites it (`SC-1 [GREETING-001 AC-1]`). Given/When/Then, no compound conditions: one behaviour per criterion.

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | No greeting has been changed since first product run | Visitor opens the page | Heading text is `Hello, World!` |
| AC-2 | Stored greeting is `Pipeline OK` | Visitor opens the page | Heading text is `Pipeline OK` |
| AC-3 | Stored greeting contains capitalization and punctuation `Hello, Database!` | Visitor opens the page | Heading text is exactly `Hello, Database!` |

**Failure, boundary and permission behaviour** — the part most often skipped and most often the source of bugs. Every case this function actually has needs a defined outcome; "should not happen" is not an outcome.

| Case | Condition | Expected behaviour |
|---|---|---|
| Invalid input | Not applicable to reading greeting | No read input exists |
| Boundary | Stored greeting has long text | Heading wraps within the centered section; no horizontal page scroll at 320px and wider |
| Not found | Stored greeting is absent on first product run | Initial greeting `Hello, World!` is used as current greeting |
| Not permitted | Visitor is not signed in | Visitor may read greeting; no sign-in required |
| Conflict | Multiple visitors read at same time | Each page view shows current stored greeting available for that request |
| Upstream failure | API or database is unavailable while reading | Not applicable: approved design has no loading or error screen; service contract specifies API error envelope |

**Data touched** — the fields this function reads and writes, in product terms. The physical schema is TL's job in `docs/architecture/erd.md`; this is the list that document has to satisfy.

| Field | Type | Required | Rule |
|---|---|---|---|
| Greeting text | text | yes | Initial value is `Hello, World!`; displayed text preserves saved value exactly; empty saved value is not allowed |

**Requirement GREETING-002 — Edit and save greeting**

*As a* Visitor, *I want to* enter a new greeting and save it, *so that* the page shows my new greeting.

Behaviour:

1. Page shows one form below the heading.
2. Form contains an accessible label `Greeting` for the text input; label is visually hidden in the approved design.
3. Text input value matches current greeting when the page loads.
4. Form contains one primary button labelled `Save`.
5. When Visitor submits a non-empty greeting, surrounding whitespace is removed before saving.
6. After a successful save, heading text changes to the saved greeting.
7. After a successful save, input value changes to the saved greeting.
8. After a successful save, status message text is `Saved.` in the reserved status line.
9. If Visitor submits empty text after trimming whitespace, nothing is saved, focus returns to input, and status message text is `Enter a greeting.`

**Acceptance criteria** — each is proved by at least one test case in `docs/greeting/test-cases/persist-editable-greeting.md`, through the story plan that cites it (`SC-1 [GREETING-002 AC-1]`). Given/When/Then, no compound conditions: one behaviour per criterion.

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | Stored greeting is `Hello, World!` | Visitor opens the page | Text input value is `Hello, World!` |
| AC-2 | Page is open | Visitor inspects form controls | Text input has accessible label `Greeting` and visible button text is `Save` |
| AC-3 | Page is open and input value is `New greeting` | Visitor saves | Heading text becomes `New greeting` |
| AC-4 | Page is open and input value is `  Trim me  ` | Visitor saves | Heading text becomes `Trim me` |
| AC-5 | Page is open and input value is empty or whitespace only | Visitor saves | Status text becomes `Enter a greeting.` |
| AC-6 | Page is open and input value is empty or whitespace only | Visitor saves | Focus is on the greeting input |
| AC-7 | Page is open and input value is empty or whitespace only | Visitor saves | Stored greeting remains unchanged |
| AC-8 | Page is open and input value is `Saved value` | Visitor saves | Status text becomes `Saved.` |
| AC-9 | Greeting `Saved value` was saved | Visitor reloads the page | Heading text is `Saved value` |
| AC-10 | Greeting `Saved value` was saved | Visitor reloads the page | Text input value is `Saved value` |

**Failure, boundary and permission behaviour** — the part most often skipped and most often the source of bugs. Every case this function actually has needs a defined outcome; "should not happen" is not an outcome.

| Case | Condition | Expected behaviour |
|---|---|---|
| Invalid input | Visitor submits empty or whitespace-only text | No save occurs; status line shows `Enter a greeting.`; focus returns to input |
| Boundary | Visitor submits text with leading or trailing whitespace | Whitespace is trimmed before save and display |
| Boundary | Visitor submits long text | Save is accepted if text remains non-empty after trimming; heading wraps within the centered section; no horizontal page scroll at 320px and wider |
| Not found | Stored greeting is absent before save | Current value is `Hello, World!`; Visitor can save over it |
| Not permitted | Visitor is not signed in | Visitor may save greeting; no sign-in required |
| Conflict | Two visitors save different greetings | Last successful save is current greeting after reload |
| Upstream failure | API or database is unavailable while saving | Not applicable: approved design has no loading or error screen; service contract specifies API error envelope |

**Data touched** — the fields this function reads and writes, in product terms. The physical schema is TL's job in `docs/architecture/erd.md`; this is the list that document has to satisfy.

| Field | Type | Required | Rule |
|---|---|---|---|
| Greeting text | text | yes | Saved value is trimmed; after trimming it must contain at least 1 character; saved value persists after reload |

## 5. Screens

The design is the source of truth for appearance; this section maps functions onto it so nothing in the design is unaccounted for and nothing specified here is missing from the design.

List only the states the approved design actually shows. The approved design shows one screen state: `default`.

| Screen | Section in the design | Functions it serves | States that must exist |
|---|---|---|---|
| Greeting page | `main > section.greeting-section` with `h1#greeting`, form `#greeting-form`, visually hidden `label` text `Greeting`, `input#greeting-input`, `button` text `Save`, and `p#status` live region | GREETING-001, GREETING-002 | default |

## 6. Non-functional requirements

Only what is real for this module. Delete rows that do not apply rather than inventing a number nobody will check.

| Area | Requirement |
|---|---|
| Accessibility | Text input has associated label `Greeting`; status line uses live region semantics; input and button are keyboard reachable; focus outline is visible on input and button; text contrast is at least 4.5:1 |
| Responsive | Page works at 320px and wider with no horizontal page scroll; form stacks into one column at widths up to 520px |
| Persistence | A saved greeting remains current after browser reload |
| Privacy | No personal data is stored; only shared greeting text is stored |

## 7. Dependencies and assumptions

- **Depends on:** Backend API, for reading current greeting and saving changed greeting.
- **Depends on:** PostgreSQL database, for persisting greeting across reloads.
- **Assumption:** Greeting is shared globally for all visitors. If false, sign-in or visitor identity would become new scope.

| Open question | Proposed default | Who decides |
|---|---|---|
| — | No open questions | — |

## 8. Traceability

Every plan item in this module appears exactly once, and every requirement id traces to a test case. A gap in this table is a gap in the build.

| Plan item | Requirement ids | Test cases |
|---|---|---|
| Persist editable greeting | GREETING-001, GREETING-002 | `test-cases/persist-editable-greeting.md` |
