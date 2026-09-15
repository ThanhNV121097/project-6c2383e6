# Test Cases — Persist editable greeting

Module: `greeting`
Function: Persist editable greeting
Risk level: Medium. This story writes shared persisted data through frontend, API, and database; failures can lose user changes or break acceptance path. Scope remains small: one public singleton greeting, no accounts, no external services.

## Page and interaction cases

**Scenario**: Initial greeting appears in heading
**Given**: No greeting has been changed since first product run, or singleton greeting row is absent
**When**: Visitor opens the page
**Then**: `h1#greeting` text is exactly `Hello, World!`
Traces: SC-1 (GREETING-001 AC-1)
Check: render_url

**Scenario**: Stored greeting appears in heading
**Given**: Stored greeting is `Pipeline OK`
**When**: Visitor opens the page
**Then**: `h1#greeting` text is exactly `Pipeline OK`
Traces: SC-2 (GREETING-001 AC-2)
Check: render_url

**Scenario**: Capitalization and punctuation are preserved in heading
**Given**: Stored greeting is `Hello, Database!`
**When**: Visitor opens the page
**Then**: `h1#greeting` text is exactly `Hello, Database!`
Traces: SC-3 (GREETING-001 AC-3)
Check: render_url

**Scenario**: Input is prefilled from stored greeting
**Given**: Stored greeting is `Hello, World!`
**When**: Visitor opens the page
**Then**: `input#greeting-input` value is exactly `Hello, World!`
Traces: SC-4 (GREETING-002 AC-1)
Check: render_url

**Scenario**: Text input has accessible label
**Given**: Page is open
**When**: Visitor inspects form controls
**Then**: `input#greeting-input` has accessible label text exactly `Greeting`
Traces: SC-5 (GREETING-002 AC-2)
Check: render_url

**Scenario**: Save button text is visible
**Given**: Page is open
**When**: Visitor inspects form controls
**Then**: Visible button text is exactly `Save`
Traces: SC-6 (GREETING-002 AC-2)
Check: render_url

**Scenario**: Save updates heading
**Given**: Page is open and `input#greeting-input` value is `New greeting`
**When**: Visitor activates `button` with text `Save`
**Then**: `h1#greeting` text becomes exactly `New greeting`
Traces: SC-7 (GREETING-002 AC-3)
Check: interact_page

**Scenario**: Save trims surrounding whitespace
**Given**: Page is open and `input#greeting-input` value is `  Trim me  `
**When**: Visitor activates `button` with text `Save`
**Then**: `h1#greeting` text becomes exactly `Trim me`, and `input#greeting-input` value becomes exactly `Trim me`
Traces: SC-8 (GREETING-002 AC-4)
Check: interact_page

**Scenario**: Empty submit shows validation status
**Given**: Page is open with stored greeting `Keep me`, and `input#greeting-input` value is empty
**When**: Visitor activates `button` with text `Save`
**Then**: `p#status` text becomes exactly `Enter a greeting.`
Traces: SC-9 (GREETING-002 AC-5)
Check: interact_page

**Scenario**: Whitespace-only submit focuses input
**Given**: Page is open with stored greeting `Keep me`, and `input#greeting-input` value is `   `
**When**: Visitor activates `button` with text `Save`
**Then**: Focus is on `input#greeting-input`
Traces: SC-10 (GREETING-002 AC-6)
Check: interact_page

**Scenario**: Invalid empty submit does not change stored greeting
**Given**: Page is open with stored greeting `Keep me`, and `input#greeting-input` value is empty
**When**: Visitor activates `button` with text `Save`, then reloads the page
**Then**: `h1#greeting` text is exactly `Keep me`, and `input#greeting-input` value is exactly `Keep me`
Traces: SC-11 (GREETING-002 AC-7)
Check: interact_page

**Scenario**: Successful save shows saved status
**Given**: Page is open and `input#greeting-input` value is `Saved value`
**When**: Visitor activates `button` with text `Save`
**Then**: `p#status` text becomes exactly `Saved.`
Traces: SC-12 (GREETING-002 AC-8)
Check: interact_page

**Scenario**: Reload keeps saved greeting in heading
**Given**: Visitor saved greeting `Saved value`
**When**: Visitor reloads the page
**Then**: `h1#greeting` text is exactly `Saved value`
Traces: SC-13 (GREETING-002 AC-9)
Check: interact_page

**Scenario**: Reload keeps saved greeting in input
**Given**: Visitor saved greeting `Saved value`
**When**: Visitor reloads the page
**Then**: `input#greeting-input` value is exactly `Saved value`
Traces: SC-14 (GREETING-002 AC-10)
Check: interact_page

**Scenario**: Long stored greeting wraps without horizontal scroll
**Given**: Stored greeting is a long non-empty text of at least 200 characters
**When**: Visitor opens the page at 320px viewport width
**Then**: `h1#greeting` wraps within `section.greeting-section`, document horizontal scroll width is not greater than viewport width, and no horizontal page scroll exists
Traces: SC-15 (GREETING-001 boundary)
Check: measure_styles

**Scenario**: Long saved greeting is accepted and wraps without horizontal scroll
**Given**: Page is open at 320px viewport width and `input#greeting-input` value is a long non-empty text of at least 200 characters
**When**: Visitor activates `button` with text `Save`
**Then**: `h1#greeting` text becomes the submitted long text, `h1#greeting` wraps within `section.greeting-section`, and document horizontal scroll width is not greater than viewport width
Traces: SC-16 (GREETING-002 boundary)
Check: interact_page

**Scenario**: Anonymous visitor can save greeting
**Given**: Visitor has no sign-in session and page is open with `input#greeting-input` value `Anonymous save`
**When**: Visitor activates `button` with text `Save`
**Then**: No authentication prompt appears, `h1#greeting` text becomes exactly `Anonymous save`, and `p#status` text becomes exactly `Saved.`
Traces: SC-17 (GREETING-002 permission)
Check: interact_page

**Scenario**: Last successful save wins after concurrent visitors
**Given**: Two visitor pages are open from the same stored greeting
**When**: First visitor saves `First save`, then second visitor saves `Second save`, then either visitor reloads the page
**Then**: `h1#greeting` text is exactly `Second save`, and `input#greeting-input` value is exactly `Second save`
Traces: SC-18 (GREETING-002 conflict)
Check: interact_page

## Service contract cases

**Scenario**: GET greeting returns current text
**Given**: Stored greeting is `API greeting`
**When**: Client requests `GET /v1/greeting`
**Then**: Response status is `200` and JSON body is exactly `{"text":"API greeting"}`
Traces: contract (GET /v1/greeting)
Check: fetch_url

**Scenario**: GET greeting returns initial value when row is absent
**Given**: Singleton greeting row is absent
**When**: Client requests `GET /v1/greeting`
**Then**: Response status is `200` and JSON body is exactly `{"text":"Hello, World!"}`
Traces: contract (GET /v1/greeting)
Check: fetch_url

**Scenario**: PUT greeting saves trimmed text
**Given**: API is available
**When**: Client sends `PUT /v1/greeting` with JSON body `{"text":"  New greeting  "}`
**Then**: Response status is `200` and JSON body is exactly `{"text":"New greeting"}`
Traces: contract (PUT /v1/greeting)
Check: fetch_url

**Scenario**: PUT greeting rejects invalid JSON
**Given**: API is available
**When**: Client sends `PUT /v1/greeting` with malformed JSON body `{`
**Then**: Response status is `400` and JSON body is exactly `{"error":{"code":"MALFORMED_REQUEST","message":"Request body is malformed."}}`
Traces: contract (PUT /v1/greeting)
Check: fetch_url

**Scenario**: PUT greeting rejects wrong type
**Given**: API is available
**When**: Client sends `PUT /v1/greeting` with JSON body `{"text":123}`
**Then**: Response status is `400` and JSON body is exactly `{"error":{"code":"MALFORMED_REQUEST","message":"Request body is malformed."}}`
Traces: contract (PUT /v1/greeting)
Check: fetch_url

**Scenario**: PUT greeting rejects unknown field
**Given**: API is available
**When**: Client sends `PUT /v1/greeting` with JSON body `{"text":"Hello","extra":"ignored?"}`
**Then**: Response status is `400` and JSON body is exactly `{"error":{"code":"MALFORMED_REQUEST","message":"Request body is malformed."}}`
Traces: contract (PUT /v1/greeting)
Check: fetch_url

**Scenario**: PUT greeting rejects empty after trim
**Given**: API is available
**When**: Client sends `PUT /v1/greeting` with JSON body `{"text":"   "}`
**Then**: Response status is `422` and JSON body is exactly `{"error":{"code":"VALIDATION_FAILED","message":"Greeting must not be empty."}}`
Traces: contract (PUT /v1/greeting)
Check: fetch_url

**Scenario**: API database query failure returns internal error
**Given**: Backend is running and a database query for greeting fails
**When**: Client requests `GET /v1/greeting`
**Then**: Response status is `500` and JSON body is exactly `{"error":{"code":"INTERNAL","message":"Internal server error."}}`
Traces: contract (GET /v1/greeting)
Check: manual

**Scenario**: Health returns OK after migration and database check
**Given**: Migration has succeeded and database accepts `SELECT 1`
**When**: Client requests `GET /healthz`
**Then**: Response status is `200`
Traces: contract (GET /healthz)
Check: fetch_url
