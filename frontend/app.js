// API Base URL - use relative path to work with any domain
const API_BASE_URL = '';

// State management
let allFiles = [];
let selectedFiles = new Set();
let eventSource = null;

// DOM elements
const directoryPath = document.getElementById('directoryPath');
const loadFilesBtn = document.getElementById('loadFilesBtn');
const fileSelectionSection = document.getElementById('fileSelectionSection');
const fileList = document.getElementById('fileList');
const selectAllBtn = document.getElementById('selectAllBtn');
const selectNoneBtn = document.getElementById('selectNoneBtn');
const selectedCount = document.getElementById('selectedCount');
const configSection = document.getElementById('configSection');
const separator = document.getElementById('separator');
const submitSection = document.getElementById('submitSection');
const submitBtn = document.getElementById('submitBtn');
const progressSection = document.getElementById('progressSection');
const progressLog = document.getElementById('progressLog');

// Event listeners
loadFilesBtn.addEventListener('click', loadFiles);
selectAllBtn.addEventListener('click', selectAll);
selectNoneBtn.addEventListener('click', selectNone);
submitBtn.addEventListener('click', processFiles);

// Load files from directory
async function loadFiles() {
    const path = directoryPath.value.trim();
    if (!path) {
        alert('Please enter a directory path');
        return;
    }

    loadFilesBtn.disabled = true;
    loadFilesBtn.innerHTML = '<span class="loading"></span> Loading...';

    try {
        const response = await fetch(`${API_BASE_URL}/api/list-files`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ path }),
        });

        const data = await response.json();

        if (data.status === 'success' && data.data && data.data.length > 0) {
            allFiles = data.data;
            displayFiles();
            fileSelectionSection.style.display = 'block';
            configSection.style.display = 'block';
            submitSection.style.display = 'block';
        } else {
            alert(data.message || 'No .txt files found in directory');
            fileSelectionSection.style.display = 'none';
            configSection.style.display = 'none';
            submitSection.style.display = 'none';
        }
    } catch (error) {
        alert(`Error loading files: ${error.message}`);
    } finally {
        loadFilesBtn.disabled = false;
        loadFilesBtn.textContent = 'Load Files';
    }
}

// Display files in list
function displayFiles() {
    fileList.innerHTML = '';
    selectedFiles.clear();

    allFiles.forEach((file, index) => {
        const fileItem = document.createElement('div');
        fileItem.className = 'file-item';

        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.id = `file-${index}`;
        checkbox.checked = true;
        checkbox.addEventListener('change', updateSelectedCount);
        selectedFiles.add(file.name);

        const label = document.createElement('label');
        label.htmlFor = `file-${index}`;
        label.textContent = file.name;

        const sizeSpan = document.createElement('span');
        sizeSpan.className = 'file-size';
        sizeSpan.textContent = formatFileSize(file.size);

        fileItem.appendChild(checkbox);
        fileItem.appendChild(label);
        fileItem.appendChild(sizeSpan);
        fileList.appendChild(fileItem);
    });

    updateSelectedCount();
}

// Format file size
function formatFileSize(bytes) {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
}

// Select all files
function selectAll() {
    const checkboxes = fileList.querySelectorAll('input[type="checkbox"]');
    checkboxes.forEach((cb, index) => {
        cb.checked = true;
        selectedFiles.add(allFiles[index].name);
    });
    updateSelectedCount();
}

// Select none
function selectNone() {
    const checkboxes = fileList.querySelectorAll('input[type="checkbox"]');
    checkboxes.forEach(cb => {
        cb.checked = false;
    });
    selectedFiles.clear();
    updateSelectedCount();
}

// Update selected count
function updateSelectedCount() {
    selectedFiles.clear();
    const checkboxes = fileList.querySelectorAll('input[type="checkbox"]');
    checkboxes.forEach((cb, index) => {
        if (cb.checked) {
            selectedFiles.add(allFiles[index].name);
        }
    });
    selectedCount.textContent = `${selectedFiles.size} files selected`;
}

// Process files
async function processFiles() {
    if (selectedFiles.size === 0) {
        alert('Please select at least one file to process');
        return;
    }

    const path = directoryPath.value.trim();
    const sep = separator.value || ':';
    const files = Array.from(selectedFiles);

    // Confirm processing
    if (!confirm(`Process ${files.length} files with separator "${sep}"?`)) {
        return;
    }

    // Disable UI
    submitBtn.disabled = true;
    submitBtn.innerHTML = '<span class="loading"></span> Processing...';
    loadFilesBtn.disabled = true;
    selectAllBtn.disabled = true;
    selectNoneBtn.disabled = true;

    // Show progress section
    progressSection.style.display = 'block';
    progressLog.innerHTML = '<div class="log-entry info">Connecting to server...</div>';

    // Connect to SSE
    connectSSE();

    // Start processing
    try {
        const response = await fetch(`${API_BASE_URL}/api/process`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                path,
                files,
                separator: sep,
            }),
        });

        const data = await response.json();

        if (data.status !== 'success') {
            addLogEntry(`Error: ${data.message}`, 'error');
            resetUI();
        } else {
            addLogEntry(data.message, 'info');
        }
    } catch (error) {
        addLogEntry(`Error: ${error.message}`, 'error');
        resetUI();
    }
}

// Connect to SSE
function connectSSE() {
    if (eventSource) {
        eventSource.close();
    }

    eventSource = new EventSource(`${API_BASE_URL}/api/events`);

    eventSource.onmessage = (event) => {
        const message = event.data;
        
        if (message === 'DONE') {
            addLogEntry('Processing completed!', 'done');
            resetUI();
            eventSource.close();
            eventSource = null;
        } else {
            const type = message.includes('Error') || message.includes('error') ? 'error' : 
                        message.includes('Completed') || message.includes('complete') ? 'success' : 'info';
            addLogEntry(message, type);
        }
    };

    eventSource.onerror = (error) => {
        console.error('SSE error:', error);
        if (eventSource.readyState === EventSource.CLOSED) {
            addLogEntry('Connection to server lost', 'error');
        }
    };
}

// Add log entry
function addLogEntry(message, type = 'info') {
    const entry = document.createElement('div');
    entry.className = `log-entry ${type}`;
    const timestamp = new Date().toLocaleTimeString();
    entry.textContent = `[${timestamp}] ${message}`;
    progressLog.appendChild(entry);
    progressLog.scrollTop = progressLog.scrollHeight;
}

// Reset UI after processing
function resetUI() {
    submitBtn.disabled = false;
    submitBtn.textContent = 'Process Files';
    loadFilesBtn.disabled = false;
    selectAllBtn.disabled = false;
    selectNoneBtn.disabled = false;
}

// Cleanup on page unload
window.addEventListener('beforeunload', () => {
    if (eventSource) {
        eventSource.close();
    }
});
