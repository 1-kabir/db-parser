# Email File Processor

A high-performance email processing system with a Go backend and minimal black & white frontend. Efficiently processes thousands of email files with duplicate detection and validation.

## Features

### Core Functionality
- ✅ **Batch File Processing**: Handle 6,000+ files efficiently with worker pools
- ✅ **Email Validation**: Regex-based email format validation
- ✅ **Duplicate Detection**: Global duplicate detection across all files
- ✅ **Configurable Separator**: Change email/data separator (default: colon)
- ✅ **Real-time Progress**: SSE-based live progress updates
- ✅ **Memory Efficient**: Line-by-line streaming processing
- ✅ **Smart Output**: Files renamed with line counts, invalid emails logged

### Backend (Go)
- Concurrent file processing with worker pools
- Line-by-line streaming for memory efficiency
- Email validation with comprehensive regex
- Global duplicate tracking across files
- SSE (Server-Sent Events) for real-time progress
- Automatic output directory creation
- Invalid email logging to `[INVALID].txt` files

### Frontend (HTML/CSS/JS)
- Minimal black & white monospace theme
- Directory path input with file discovery
- File selection with Select All/None
- Configurable email separator
- Real-time SSE progress display
- Processing status with timestamps

## Project Structure

```
db-parser/
├── backend/          # Go HTTP server
│   ├── main.go       # Main server with API endpoints & SSE
│   ├── processor.go  # Email processing logic
│   ├── go.mod        # Go module file
│   └── README.md     # Backend documentation
├── frontend/         # Static HTML/CSS/JS site
│   ├── index.html    # Main HTML page
│   ├── styles.css    # Minimal B&W styles
│   ├── app.js        # JavaScript for API calls & SSE
│   └── README.md     # Frontend documentation
└── README.md         # This file
```

## Getting Started

### Prerequisites

- Go 1.16 or higher
- A modern web browser
- Python 3 (for serving frontend)

### Quick Start

**Option 1: Using PM2 (Recommended)**
```bash
cd /path/to/db-parser
go build -o backend/email-processor backend/main.go backend/processor.go
pm2 start ecosystem.config.js
```

**Option 2: Manual Start**
```bash
# Terminal 1 - Backend (default port 8080)
cd backend
go run main.go

# Or with custom port:
PORT=9000 go run main.go
```

The frontend files are automatically served by the backend.

**Accessing the Application:**
- Open `http://localhost:8080` in your browser (or your custom port)

### Port Configuration

**Default:** Port `8080`

**Change the port:**
```bash
PORT=9000 go run main.go
# or with compiled binary:
PORT=9000 ./backend/email-processor
```

**Using PM2:** Edit `ecosystem.config.js` and update the `env` section:
```javascript
env: {
  PORT: '9000'  // Change this to your desired port
}
```

Open `http://localhost:9000` in your browser.

## Usage

1. **Select Directory**: Enter the path to your folder containing .txt files
2. **Load Files**: Click "Load Files" to discover all .txt files
3. **Select Files**: Choose which files to process (Select All/None available)
4. **Configure**: Set email separator (default is `:`)
5. **Process**: Click "Process Files" and watch real-time progress

## How It Works

### Email Processing
Each line in your .txt files is processed as follows:
```
email@example.com:additional_data
```

The processor:
1. Extracts the email (before separator)
2. Validates email format
3. Checks for duplicates across ALL files
4. Writes valid emails to output file
5. Logs invalid/duplicate emails to `[INVALID].txt`

### Output Files
All processed files go to an `outputs/` directory in your source folder:
- `filename {N}.txt` - Cleaned emails (N = line count)
- `filename [INVALID].txt` - Invalid/duplicate emails with reasons

### Example
**Input:** `sample {100}.txt` with 100 lines, 80 valid, 20 invalid  
**Output:**
- `outputs/sample {80}.txt` - 80 valid emails
- `outputs/sample [INVALID].txt` - 20 invalid entries with tags:
  - `[DUPLICATE]` - Email already exists
  - `[INVALID_FORMAT]` - Invalid email format
  - `[NO_EMAIL]` - No email found in line

## Performance

- **Concurrent Processing**: Uses worker pools (CPU cores, max 8)
- **Memory Efficient**: Streams files line-by-line
- **Fast**: Processes 10,000+ line files in seconds
- **Scalable**: Handles 6,000+ files in batches

## API Endpoints

- `GET /health` - Server health check
- `POST /api/list-files` - List .txt files in directory
- `POST /api/process` - Start processing job
- `GET /api/events` - SSE stream for progress updates

## Development

See individual README files in the `backend/` and `frontend/` directories for detailed information.

## License

MIT