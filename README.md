# DB Parser

A full-stack application with a Go backend server and a static HTML/CSS/JS frontend.

## Project Structure

```
db-parser/
├── backend/          # Go HTTP server
│   ├── main.go       # Main server file with API endpoints
│   ├── go.mod        # Go module file
│   └── README.md     # Backend documentation
├── frontend/         # Static HTML/CSS/JS site
│   ├── index.html    # Main HTML page
│   ├── styles.css    # Styles and layout
│   ├── app.js        # JavaScript for API calls
│   └── README.md     # Frontend documentation
└── README.md         # This file
```

## Getting Started

### Prerequisites

- Go 1.16 or higher
- A modern web browser
- (Optional) Python 3 or Node.js for serving the frontend

### Running the Backend

1. Navigate to the backend directory:
```bash
cd backend
```

2. Run the Go server:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

### Running the Frontend

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Serve the static files using Python:
```bash
python3 -m http.server 3000
```

Or using Node.js:
```bash
npx http-server -p 3000
```

3. Open your browser and go to `http://localhost:3000`

## Features

### Backend
- RESTful API endpoints
- CORS enabled for frontend integration
- JSON response format
- Health check endpoint
- Sample data endpoints

### Frontend
- Modern, responsive design
- API integration with backend
- Interactive UI to test endpoints
- Error handling and loading states

## API Endpoints

- `GET /health` - Check server health status
- `GET /api` - Get API information and available endpoints
- `GET /api/data` - Retrieve sample data

## Development

See individual README files in the `backend/` and `frontend/` directories for more detailed information about each component.

## License

MIT