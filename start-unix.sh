#!/bin/bash
echo "========================================"
echo "  Email File Processor"
echo "========================================"
echo ""
echo "Starting server on http://localhost:8080"
echo ""
echo "Open your browser and go to:"
echo "  http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop the server"
echo "========================================"
echo ""

cd "$(dirname "$0")/backend"
./email-processor
