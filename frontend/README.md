# Frontend Static Site

A simple HTML/CSS/JavaScript static site that connects to the backend Go server.

## Features

- Clean, modern UI design
- Responsive layout
- API integration with backend server
- Interactive buttons to test different endpoints
- Error handling and loading states
- CORS-enabled for local development

## Files

- `index.html` - Main HTML structure
- `styles.css` - Styling and layout
- `app.js` - JavaScript for API calls and interactivity

## Running the Frontend

### Option 1: Using a Simple HTTP Server

#### Python 3
```bash
cd frontend
python3 -m http.server 3000
```

#### Node.js (http-server)
```bash
cd frontend
npx http-server -p 3000
```

#### Go
```bash
cd frontend
go run -m http.server 3000
```

### Option 2: Open Directly in Browser

You can also open `index.html` directly in your browser, though some features may not work due to CORS restrictions.

## Usage

1. Make sure the backend server is running on `http://localhost:8080`
2. Open the frontend in your browser
3. Click the buttons to interact with the backend:
   - **Check Server Status**: Verify the backend is running
   - **Fetch Data**: Retrieve sample data from the backend
   - **Get API Info**: View available API endpoints

## Configuration

To change the backend API URL, edit the `API_BASE_URL` constant in `app.js`:

```javascript
const API_BASE_URL = 'http://localhost:8080';
```

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## Development

The frontend uses vanilla JavaScript with no build step required. Simply edit the files and refresh your browser to see changes.

### File Structure
```
frontend/
├── index.html    # Main HTML page
├── styles.css    # Styles and layout
├── app.js        # JavaScript logic
└── README.md     # This file
```
