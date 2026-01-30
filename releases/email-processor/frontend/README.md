# Frontend - Email File Processor UI

Minimal black & white interface for email file processing with real-time progress updates.

## Design Philosophy

- **Minimal**: Clean black & white theme, monospace font
- **Functional**: Every element serves a purpose
- **Real-time**: SSE-based live progress updates
- **Responsive**: Works on desktop and tablet screens

## Features

### UI Components

1. **Directory Input**
   - Text input for file path
   - "Load Files" button
   - Validates directory existence

2. **File Selection**
   - Checkbox list of all .txt files
   - File sizes displayed
   - "Select All" / "Select None" buttons
   - Selected count display

3. **Configuration**
   - Email separator input (default: `:`)
   - Supports any character(s)

4. **Processing**
   - Large "Process Files" button
   - Confirmation dialog
   - Loading state with animation

5. **Progress Log**
   - Real-time SSE updates
   - Timestamped entries
   - Auto-scroll to bottom
   - Color-coded messages:
     - Info: Standard messages
     - Success: Completion messages
     - Error: Error messages
     - Done: Final completion

## Files

### `index.html`
Clean semantic HTML structure:
- No unnecessary divs
- Proper heading hierarchy
- Accessible form elements
- Progressive disclosure (sections show as needed)

### `styles.css`
Minimal black & white styling:
- Monospace font family
- Pure black (#000) and white (#fff)
- 1px/2px borders
- No gradients, shadows, or transitions
- Simple hover states (invert colors)
- Custom scrollbar styling

### `app.js`
Vanilla JavaScript (no frameworks):
- ES6+ features
- Async/await for API calls
- EventSource for SSE
- Efficient DOM manipulation
- No external dependencies

## Usage Flow

1. User enters directory path
2. Click "Load Files" → API call to `/api/list-files`
3. Files displayed with checkboxes (all selected by default)
4. User adjusts selection and separator
5. Click "Process Files" → Confirmation dialog
6. SSE connection established to `/api/events`
7. Progress updates stream in real-time
8. "DONE" message triggers UI reset

## API Integration

### Endpoints Used

**`POST /api/list-files`**
```javascript
fetch('http://localhost:8080/api/list-files', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ path: '/path/to/files' })
})
```

**`POST /api/process`**
```javascript
fetch('http://localhost:8080/api/process', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    path: '/path/to/files',
    files: ['file1.txt', 'file2.txt'],
    separator: ':'
  })
})
```

**`GET /api/events`** (SSE)
```javascript
const eventSource = new EventSource('http://localhost:8080/api/events');
eventSource.onmessage = (event) => {
  console.log(event.data);
};
```

## Running

### Option 1: Python HTTP Server
```bash
cd frontend
python3 -m http.server 3000
```

### Option 2: Node.js HTTP Server
```bash
cd frontend
npx http-server -p 3000
```

### Option 3: Go HTTP Server
Create `serve.go`:
```go
package main
import (
    "net/http"
    "log"
)
func main() {
    log.Fatal(http.ListenAndServe(":3000", http.FileServer(http.Dir("."))))
}
```
Run: `go run serve.go`

Open `http://localhost:3000` in your browser.

## Configuration

### Backend URL
Edit `app.js`:
```javascript
const API_BASE_URL = 'http://localhost:8080'; // Change here
```

### Port
Change when starting HTTP server:
```bash
python3 -m http.server 8000  # Use port 8000
```

## Browser Compatibility

Tested on:
- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

**Requirements:**
- ES6+ support
- EventSource API (SSE)
- Fetch API
- CSS Grid

## Development

### File Structure
```
frontend/
├── index.html    # Main HTML page
├── styles.css    # Minimal B&W styles
├── app.js        # JavaScript logic
└── README.md     # This file
```

### Making Changes

**HTML:**
- Keep structure semantic
- Use progressive disclosure
- Maintain accessibility

**CSS:**
- Keep it minimal (black & white only)
- Use monospace font
- No animations except loading spinner
- Maintain high contrast

**JavaScript:**
- Keep vanilla (no frameworks)
- Use modern ES6+ features
- Handle errors gracefully
- Clean up SSE connections

## State Management

Simple state object:
```javascript
let allFiles = [];           // All discovered files
let selectedFiles = new Set(); // Currently selected
let eventSource = null;        // SSE connection
```

## Event Handling

### Load Files
1. Validate path input
2. Call `/api/list-files`
3. Display files with checkboxes
4. Show config and submit sections

### Select All/None
1. Toggle all checkboxes
2. Update `selectedFiles` Set
3. Update count display

### Process Files
1. Validate selection
2. Show confirmation dialog
3. Disable UI elements
4. Connect SSE
5. Call `/api/process`
6. Stream progress updates

### SSE Messages
- Parse message type
- Add to log with timestamp
- Auto-scroll to bottom
- On "DONE": reset UI

## Error Handling

- Network errors → Alert dialog
- Empty path → Alert before API call
- No files selected → Alert before processing
- SSE errors → Console log + status message
- Invalid directory → Display error from API

## Performance

- Efficient DOM updates (DocumentFragment when needed)
- Debounced checkbox updates
- Minimal reflows/repaints
- Lazy image loading (N/A)
- No unnecessary libraries

## Accessibility

- Semantic HTML
- Proper label associations
- Keyboard navigation
- Focus visible states
- Screen reader friendly

## Future Enhancements

Possible additions (not implemented):
- Dark mode toggle (B&W remains)
- Keyboard shortcuts
- Download progress as CSV
- Retry failed files
- Pause/resume processing
