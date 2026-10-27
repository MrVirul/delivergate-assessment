---
name: htmx
description: How to build interactive web UIs with htmx 2.x, plain HTML, CSS and vanilla JavaScript, without a framework or build step. Use this skill whenever the user mentions htmx, hx-get / hx-post / hx-target / hx-swap attributes, "HTML over the wire", server-driven UI, a lightweight frontend without React, or wants a simple frontend for a REST or JSON API using only HTML, CSS and JavaScript. Use it even if the user only says "simple frontend" and the stack is htmx. It pins the version to htmx 2.x (NOT 4), lists the attributes and events that actually exist, and gives a tested pattern for using htmx with a JSON API and a strict Content-Security-Policy.
---

# htmx 2.x skill

Use this when writing or reviewing any htmx code. Follow it strictly: htmx has a version split and many tutorials online mix versions.

## 0. Version rule (read first)

- **Use htmx 2.x. Pin `htmx.org@2.0.11`** (the latest 2.x at the time of writing; verify with `npm view htmx.org@2 version` if you have network access and use the newest 2.x patch).
- **Do NOT use htmx 4** (`htmx.org@next`, `four.htmx.org` docs, `htmx:before:request` style event names, `hx-partial`, `:inherited` modifiers). htmx 4 was released in August 2026 and swaps fetch for XMLHttpRequest, renames every event and changes defaults. It is not the `latest` npm tag, and most examples in your training data are htmx 1.x/2.x.
- Vendor the file locally so the app works offline and is not affected by CDN changes:

```bash
mkdir -p web/vendor
curl -L -o web/vendor/htmx.min.js https://cdn.jsdelivr.net/npm/htmx.org@2.0.11/dist/htmx.min.js
head -c 200 web/vendor/htmx.min.js   # sanity check: must start with the htmx 2.x banner/code, not an HTML error page
```

If the version above does not exist, use the newest `2.0.x` and note it in the decision log. Never use an unpinned CDN URL.

## 1. Mental model

htmx lets any element issue an HTTP request and swap the response into the page, declared with attributes:

```
trigger  ->  request  ->  response  ->  swap into target
```

Normally the server returns **HTML fragments**. If the backend is a JSON API, add a small adapter that converts JSON to HTML in the browser (section 8). The rest of the page stays plain HTML + CSS + a little vanilla JS.

## 2. Core attributes (all valid in 2.x)

| Attribute | Meaning | Example |
|---|---|---|
| `hx-get` `hx-post` `hx-put` `hx-patch` `hx-delete` | Issue the request on trigger | `hx-get="/api/items"` |
| `hx-trigger` | What fires it. Default: `click` (buttons/links), `change` (inputs), `submit` (forms) | `hx-trigger="load"` |
| `hx-target` | Where the response goes. CSS selector, `this`, `closest <sel>`, `find <sel>`, `next`, `previous` | `hx-target="#results"` |
| `hx-swap` | How it goes in (section 4) | `hx-swap="outerHTML"` |
| `hx-include` | Add other elements' values to the request | `hx-include="#filters"` |
| `hx-vals` | Extra JSON values (static JSON only) | `hx-vals='{"page": 2}'` |
| `hx-headers` | Extra headers (static JSON) | `hx-headers='{"X-Foo":"1"}'` |
| `hx-indicator` | Element that gets class `htmx-request` while loading | `hx-indicator="#spinner"` |
| `hx-disabled-elt` | Disable elements during the request | `hx-disabled-elt="find button"` |
| `hx-confirm` | Browser `confirm()` before sending | `hx-confirm="Delete this item?"` |
| `hx-push-url` | Update browser URL/history | `hx-push-url="true"` |
| `hx-boost` | Turn normal links/forms into AJAX | `<body hx-boost="true">` |
| `hx-select` | Pick part of the response to swap | `hx-select="#list"` |
| `hx-swap-oob` | Swap extra elements out-of-band (by id) | `<div id="x" hx-swap-oob="true">` |
| `hx-sync` | Coordinate overlapping requests | `hx-sync="this:replace"` |
| `hx-ext` | Enable an extension (inherited by children) | `<body hx-ext="api">` |
| `hx-params` | Filter which params are sent | `hx-params="*"` |

Inheritance: in 2.x most `hx-*` attributes inherit from parents (so `hx-target` on a form applies to its buttons). Use `hx-target="unset"` or `hx-disinherit` to stop that.

## 3. Triggers

```html
<input type="search" name="search"
       hx-get="/api/items" hx-target="#results"
       hx-trigger="input changed delay:300ms, search"
       hx-sync="this:replace">                 <!-- debounce + cancel stale requests -->

<div hx-get="/api/items" hx-trigger="load, items-changed from:body"> <!-- load once + reload on custom event -->
<div hx-get="/more" hx-trigger="revealed">      <!-- lazy load when scrolled into view -->
<div hx-get="/clock" hx-trigger="every 30s">    <!-- polling -->
```

Modifiers: `changed`, `delay:`, `throttle:`, `once`, `from:<selector>`, `target:<selector>`, `consume`, `queue:`. Filters like `click[ctrlKey]` need `allowEval` and must NOT be used under a strict CSP.

Fire your own events from JS: `htmx.trigger(document.body, 'items-changed')`.

## 4. Swaps

`hx-swap` values: `innerHTML` (default), `outerHTML`, `beforebegin`, `afterbegin`, `beforeend`, `afterend`, `delete`, `none`.
Modifiers: `swap:100ms`, `settle:100ms`, `scroll:top`, `show:top`, `focus-scroll:true`, `transition:true`.

- Removing a row after DELETE: `hx-target="closest tr" hx-swap="outerHTML"` and respond with empty content (needs `responseHandling` to allow swapping 204, see section 6).
- Replacing a form after save: `hx-swap="outerHTML"` on the form.
- Do nothing visually (you will handle it in JS): `hx-swap="none"`.

## 5. Events (htmx 2.x names — camelCase)

| Event | Use |
|---|---|
| `htmx:configRequest` | Last chance to change `evt.detail.headers`, `.parameters`, `.verb`, `.path`. Add the `Authorization` header here |
| `htmx:beforeRequest` / `htmx:afterRequest` | `afterRequest` has `evt.detail.successful`, `.failed`, `.xhr`, `.elt` |
| `htmx:beforeSwap` | Change swap behaviour: `evt.detail.shouldSwap`, `.target`, `.serverResponse` |
| `htmx:afterSwap` / `htmx:afterSettle` | Run code after new DOM is in place (open a dialog, focus an input) |
| `htmx:responseError` | HTTP 4xx/5xx (when `error:true` in `responseHandling`) |
| `htmx:sendError` | Network failure (server down) |
| `htmx:timeout` | Request timeout |
| `htmx:load` | Fired on each newly added element |

Listen with plain `document.body.addEventListener('htmx:afterRequest', handler)`. **Do not use `hx-on:*` attributes** under a strict CSP: they are evaluated with `eval`.

## 6. Configuration (put in a meta tag — no inline script needed)

```html
<meta name="htmx-config" content='{
  "allowEval": false,
  "selfRequestsOnly": true,
  "defaultSwapStyle": "innerHTML",
  "responseHandling": [
    {"code": "204", "swap": true},
    {"code": "[23]..", "swap": true},
    {"code": "[45]..", "swap": false, "error": true},
    {"code": "...", "swap": false}
  ]
}'>
```

- `allowEval:false` disables `hx-on`, `js:` expressions in `hx-vals`, and trigger filters. Good for CSP and safer.
- `selfRequestsOnly:true` blocks requests to other origins.
- `responseHandling` decides per status code what to do. By default, 204 and 4xx/5xx are not swapped. The config above swaps 204 (so a DELETE can remove a row with an empty response) and flags 4xx/5xx as errors so `htmx:responseError` fires.

## 7. Forms

```html
<form hx-post="/api/customers" hx-swap="none" hx-disabled-elt="find button[type=submit]">
  <label>Name <input name="name" required></label>
  <button type="submit">Save customer</button>
</form>
```

- Values are collected from named inputs. Everything is a **string** unless you convert it (section 8).
- Checkbox: sent only when checked.
- A form inside a `<dialog>`: open with `dialog.showModal()` in an `htmx:afterSwap` listener; close in an `htmx:afterRequest` listener when `evt.detail.successful`.
- `hx-delete` and `hx-get` send parameters in the URL; `hx-post/put/patch` send them in the body.

## 8. Pattern: htmx with a JSON REST API

htmx swaps text into the DOM; it does not understand JSON. Two things must be adapted:

1. **Sending JSON** (and converting number inputs to numbers).
2. **Rendering JSON responses as HTML.**

Do both in one extension using the extension hooks `onEvent`, `encodeParameters`, `transformResponse` (all exist in 2.x). Reference implementation (verify in the browser network tab and adapt; `UI` is your own module with `render`, `coerceForm`):

```js
// web/js/api-ext.js  (load after htmx and after render.js)
htmx.defineExtension('api', {
  onEvent(name, evt) {
    if (name === 'htmx:configRequest') {
      const d = evt.detail;
      const token = sessionStorage.getItem('token');
      if (token) d.headers['Authorization'] = 'Bearer ' + token;
      if (['POST', 'PUT', 'PATCH'].includes(d.verb.toUpperCase())) {
        d.headers['Content-Type'] = 'application/json';
      }
    }
    return true;
  },

  // Called for non-GET requests that send a body. `parameters` is FormData-like in 2.x.
  encodeParameters(xhr, parameters, elt) {
    xhr.overrideMimeType('text/json');
    return JSON.stringify(UI.coerceForm(parameters, elt));
  },

  // Convert a JSON response into an HTML string before htmx swaps it.
  transformResponse(text, xhr, elt) {
    if (xhr.status >= 400) return text;                       // errors handled in htmx:responseError
    const type = xhr.getResponseHeader('Content-Type') || '';
    if (!type.includes('json') || !text) return text;          // HTML partials and 204 pass through
    const key = elt.closest('[data-render]')?.dataset.render;  // which renderer to use
    if (!key) return '';
    return UI.render(key, JSON.parse(text));
  },
});
```

`UI.coerceForm(formData, form)` rules (put them in your code, they matter):

- `<input type="number">` -> `Number(value)`; empty -> omit the key.
- Money inputs (`data-money`) -> `Math.round(parseFloat(value) * 100)` (cents).
- Optional text left blank -> `null` (or omit) according to what the API validates; check the backend DTO, do not guess.
- Trim all strings.

`UI.render(key, body)` is a map of small pure functions `(apiBody) => htmlString`, one per `data-render` key. **Every interpolated value must go through `esc()`** (section 9).

Page wiring:

```html
<body hx-ext="api">
  <form id="filters" hx-get="/api/items" hx-target="#results" data-render="items-list"
        hx-trigger="input changed delay:300ms, submit, change from:#category" hx-sync="this:replace">
    <input type="search" name="search"> <select id="category" name="category">...</select>
  </form>

  <div id="results" hx-get="/api/items" hx-include="#filters"
       hx-trigger="load, items-changed from:body" data-render="items-list"></div>
</body>
```

Pagination buttons produced by the renderer:

```html
<button hx-get="/api/items" hx-include="#filters" hx-vals='{"page": 2}'
        hx-target="#results" data-render="items-list">Next</button>
```

Global handlers in `app.js` (plain `addEventListener`):

- `htmx:responseError`: parse `evt.detail.xhr.responseText`, read `error.message` and `error.details` (field errors), show a toast. If `status === 401` clear the session and redirect to `/login.html`. If `403` show "You don't have permission".
- `htmx:sendError`: toast "Cannot reach the server".
- Mutations succeed -> `htmx.trigger(document.body, 'items-changed')` to refresh the list and close the dialog.

## 9. Security (non-negotiable)

1. **Escape all dynamic text** when building HTML strings:

```js
const esc = (v) => String(v ?? '').replace(/[&<>"']/g, (c) =>
  ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));
```

   Never put API data into `innerHTML` unescaped, and never into inline `on*=` attributes or `<script>`.
2. **Strict CSP friendly:** no inline `<script>`, no inline event handlers, no `hx-on:*`, `allowEval:false`. Put JS in files, config in the meta tag. If the server sends a CSP header (for example from a security middleware), open the browser console and fix any "violates Content Security Policy" error by adjusting the server policy narrowly, not by removing it.
3. **Tokens:** a JWT in `sessionStorage` is acceptable for a simple app but is readable by any XSS, which is why rule 1 matters. Clear it on logout and on any 401.
4. The UI hides buttons the user cannot use, but **the server is the real authority**. Always handle 401/403 gracefully.
5. `selfRequestsOnly:true`; no third-party scripts except self-hosted ones.

## 10. Accessibility and UX defaults

- Every input has a `<label>`; buttons are `<button>`, links are `<a>`.
- Toast container has `role="status" aria-live="polite"`.
- Use native `<dialog>` for modals (focus trap and Escape for free). Return focus to the opener on close.
- Show progress with `hx-indicator` / `htmx-request` classes and disable submit buttons with `hx-disabled-elt`.
- Empty and error states are written for the user (what happened, what to do next).
- Respect `prefers-reduced-motion`; keep visible `:focus-visible` outlines.

## 11. Common mistakes to avoid

| Mistake | Fix |
|---|---|
| Using htmx 4 docs/events (`htmx:before:request`) | Use 2.x camelCase names (section 5) |
| Number fields sent as strings -> API returns 400/422 | Convert in `coerceForm` (section 8) |
| 204 response does not remove the row | Add `{"code":"204","swap":true}` to `responseHandling` |
| 4xx/5xx silently ignored | Mark `[45]..` with `error:true` and handle `htmx:responseError` |
| `hx-vals` with `js:` or `hx-on` under CSP | Use static JSON or `htmx:configRequest` |
| Duplicate `id`s after swaps | Ids must be unique; prefer classes or `closest` targets |
| Search fires on every keystroke and results arrive out of order | `delay:300ms` and `hx-sync="this:replace"` |
| Forgetting to re-bind JS to new elements | Use event delegation on `document` / `body`, or listen to `htmx:load` |
| Putting secrets or tokens in URLs | Use the `Authorization` header |
| `hx-target` inherited by an unrelated child | Set an explicit `hx-target` on the child |
| Swapping a `<tr>` into a non-table parent | Target the `<tbody>` or the row itself with `outerHTML` |

## 12. Debugging checklist

1. Browser DevTools -> Network: is the request sent, with the right method, headers (`Authorization`, `Content-Type`) and body?
2. Is the response status what you expect? Which `responseHandling` rule applies?
3. Console: any CSP violation or JS error from the renderer?
4. `htmx.logAll()` in the console prints every htmx event while debugging (remove afterwards).
5. Confirm the loaded file is 2.x: `htmx.version` in the console.
