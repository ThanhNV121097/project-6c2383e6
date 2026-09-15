# Design System — Hello World Acceptance 6

> Source of truth: approved `index.html`.
> Every value below is extracted from it. Changing a value here without changing approved design is defect.

Last updated: 2026-09-15

## 1. Foundations

### 1.1 Color

Semantic tokens. Name by job, never by hue.

| Token | Value | Used for |
|---|---|---|
| `--color-bg` | `#FFFFFF` | Page background, input background, button text |
| `--color-text` | `#000000` | Heading, input text, status text, input border |
| `--color-primary-action` | `#2563EB` | Save button background and border, focus outline |

#### Contrast audit

Every text-on-background pair actually used. Body text ≥ 4.5:1, large text (≥ 18.66px bold or ≥ 24px) ≥ 3:1, UI borders ≥ 3:1.

| Foreground | Background | Ratio | Passes |
|---|---|---|---|
| `--color-text` | `--color-bg` | `21:1` | AA |
| `--color-bg` | `--color-primary-action` | `5.17:1` | AA |
| `--color-primary-action` | `--color-bg` | `5.17:1` | AA for focus indicator and UI border |

### 1.2 Spacing

Base unit: `1px`, because approved design includes accessibility clipping and focus offsets at 1px and 3px. Main layout follows a sparse 12/24/32 rhythm.

| Token | Value |
|---|---|
| `--space-0` | `0` |
| `--space-1px` | `1px` |
| `--space-focus-offset` | `3px` |
| `--space-form-gap` | `12px` |
| `--space-control-y` | `14px` |
| `--space-control-x` | `16px` |
| `--space-button-x` | `22px` |
| `--space-page-padding` | `24px` |
| `--space-heading-gap` | `32px` |

### 1.3 Typography

Font families, loaded by system fallback only:

- Body: `Arial, Helvetica, sans-serif`
- Headings: `Arial, Helvetica, sans-serif`

| Token | Size | Line height | Weight | Used for |
|---|---|---|---|---|
| `--text-status` | `16px` | Browser default normal | Browser default normal | Status message |
| `--text-control` | `18px` | Browser default normal | `400` input, `700` button | Text field and Save button |
| `--text-display` | `clamp(40px, 8vw, 72px)` | `1.05` | `700` | h1 greeting |

Heading levels are used in order: page uses one `h1` only.

| Token | Value | Used for |
|---|---|---|
| `--font-weight-body` | Browser default normal | Input and status text |
| `--font-weight-action` | `700` | Save button |
| `--font-weight-heading` | `700` | h1 greeting |

### 1.4 Radius, border, shadow, motion

| Token | Value | Used for |
|---|---|---|
| `--radius-control` | `0` | Text field and button |
| `--border-control` | `2px solid` | Text field and button border |
| `--outline-focus` | `3px solid` | Text field and button focus outline |

Motion: no animation or transition is used.

### 1.5 Layout and breakpoints

| Name | Min width | Container | Columns | Gutter |
|---|---|---|---|---|
| `default` | `0` | `min(100%, 560px)` greeting section | 1 centred section | `12px` form gap |
| `compact` | `max-width: 520px` media query | `min(100%, 560px)` greeting section | Form stacks into one column | `12px` form gap |

## 2. Components

### 2.1 GreetingSection

**Purpose** — Single page region for reading and changing persisted greeting. Do not use for navigation or extra content.

**Anatomy** — `[h1 greeting] [form: visually hidden label, text input, Save button] [status message]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Default | `--color-bg`, `--color-text`, `--space-page-padding`, `--space-heading-gap` | Only page section |

**Sizes**

| Size | Width | Padding | Text token |
|---|---|---|---|
| Default | `min(100%, 560px)` | `24px` via parent `main` | `--text-display` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Centered large greeting above form; status line reserves 24px height | `--color-bg`, `--color-text`, `--text-display`, `--space-heading-gap` |

**Accessibility** — Section labelled by h1 through `aria-labelledby="greeting"`. Status uses `role="status"` and `aria-live="polite"`.

### 2.2 TextInput

**Purpose** — Edits greeting text. Use when visitor changes current greeting.

**Anatomy** — `[visually hidden label] [text input]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Default | `--color-bg`, `--color-text`, `--border-control`, `--radius-control` | Greeting input |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Default | Content-defined; 18px text plus 14px vertical padding and 2px border | `14px 16px` | `--text-control` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | White background, black text, 2px black border | `--color-bg`, `--color-text`, `--border-control` |
| Focus | 3px blue outline with 3px offset | `--color-primary-action`, `--outline-focus`, `--space-focus-offset` |

**Accessibility** — Input has associated label. Label is visually hidden but accessible. Input uses `required` and `autocomplete="off"`. Keyboard focus must keep visible outline. Minimum hit target exceeds 44×44px.

### 2.3 PrimaryButton

**Purpose** — Saves current greeting. Use for primary submit action only.

**Anatomy** — `[label]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Primary | `--color-primary-action`, `--color-bg`, `--border-control`, `--radius-control` | Save action |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Default | Content-defined; 18px text plus 14px vertical padding and 2px border | `14px 22px` | `--text-control` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Blue background and border, white bold text, pointer cursor | `--color-primary-action`, `--color-bg`, `--font-weight-action` |
| Hover | Same colors as default | `--color-primary-action`, `--color-bg` |
| Focus | 3px blue outline with 3px offset | `--color-primary-action`, `--outline-focus`, `--space-focus-offset` |

**Accessibility** — Native `button type="submit"`. Keyboard focus must keep visible outline. Minimum hit target exceeds 44×44px.

### 2.4 StatusMessage

**Purpose** — Announces save success or validation error below form.

**Anatomy** — `[message text]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Default | `--color-text`, `--text-status`, `--space-form-gap` | Save and validation feedback |

**Sizes**

| Size | Min height | Margin | Text token |
|---|---|---|---|
| Default | `24px` | `12px 0 0` | `--text-status` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Empty | Reserved 24px line, no text | `--color-text` |
| Validation error | Text reads `Enter a greeting.` | `--color-text`, `--text-status` |
| Success | Text reads `Saved.` | `--color-text`, `--text-status` |

**Accessibility** — Use `role="status"` and `aria-live="polite"` so updates announce without moving focus. Validation error returns focus to input.

## 3. Content and formatting

- Voice and tone: plain, direct, functional.
- Date, time, number, and currency formats: not used in approved design.
- Capitalization: heading preserves saved greeting text exactly; button uses title case `Save`; label uses title case `Greeting`; status messages use sentence case.
- Empty-state and error-message wording pattern: short sentence with direct action, shown in status line. Approved validation wording: `Enter a greeting.`

## 4. Known deviations

| Where | Deviation | Why it stands | Follow-up |
|---|---|---|---|
| Spacing scale | Uses 1px, 3px, 14px, 22px alongside 12/24/32 layout rhythm | Approved mockup includes precise accessibility clipping, focus offset, and control padding | Keep unless stakeholder asks for stricter spacing scale |
| PrimaryButton | Hover state repeats default colors with no visible change | Approved mockup draws this exact hover rule | Change only if stakeholder requests stronger hover feedback |

Approved design avoids AI defaults: no purple/indigo palette, no gradients, no maximum rounding, no heavy shadows, no generic multi-section layout, no emoji iconography, no filler copy, and focus states remain visible.

## 5. Change log

| Date | Change | Design PR |
|---|---|---|
| 2026-09-15 | Initial design system extracted from approved `index.html` | This PR |
