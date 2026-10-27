# UI Documentation

## Architecture
- Static files in `web/`, served by Go Fiber
- htmx 2.0.11 vendored locally
- Adapter pattern: `api-ext.js` sends JSON, `render.js` builds HTML
- All dynamic values escaped via `UI.esc`

## Adding a New Screen
Follow the pattern in renderers and add routes/elements.

## Development
Run `make run` and open http://localhost:8080/
