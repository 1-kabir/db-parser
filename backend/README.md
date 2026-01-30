# Backend Go Server

A simple HTTP server built with Go that exposes RESTful API endpoints.

## Features

- Basic HTTP server with multiple endpoints
- CORS enabled for frontend integration
- JSON response format
- Health check endpoint
- Sample data endpoints

## Prerequisites

- Go 1.16 or higher

## Installation

```bash
cd backend
go mod download
```

## Running the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check
- **URL**: `/health`
- **Method**: `GET`
- **Description**: Check if the server is running
- **Response**: JSON object with server status and timestamp

### API Info
- **URL**: `/api`
- **Method**: `GET`
- **Description**: Get API information and available endpoints
- **Response**: JSON object with API version and endpoints list

### Sample Data
- **URL**: `/api/data`
- **Method**: `GET`
- **Description**: Retrieve sample data items
- **Response**: JSON array of sample data objects

## Response Format

All endpoints return JSON responses in the following format:

```json
{
  "message": "Description of the response",
  "status": "ok|success|error",
  "data": {} // Optional data field
}
```

## Development

To build the server:

```bash
go build -o server main.go
```

To run the built binary:

```bash
./server
```
