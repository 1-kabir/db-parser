// API Base URL - Update this to match your backend server
const API_BASE_URL = 'http://localhost:8080';

// Helper function to make API calls
async function fetchAPI(endpoint) {
    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('API call failed:', error);
        throw error;
    }
}

// Helper function to display data in a formatted way
function displayJSON(data) {
    return `<pre>${JSON.stringify(data, null, 2)}</pre>`;
}

// Check server status
async function checkServerStatus() {
    const statusBox = document.getElementById('status');
    const btn = document.getElementById('checkStatusBtn');
    
    btn.disabled = true;
    btn.innerHTML = '<span class="loading"></span> Checking...';
    statusBox.className = 'status-box';
    statusBox.innerHTML = '<p>Connecting to server...</p>';

    try {
        const data = await fetchAPI('/health');
        statusBox.className = 'status-box success';
        statusBox.innerHTML = `
            <p class="success-message">${data.message}</p>
            <p><strong>Status:</strong> ${data.status}</p>
            <p><strong>Timestamp:</strong> ${data.data.timestamp}</p>
        `;
    } catch (error) {
        statusBox.className = 'status-box error';
        statusBox.innerHTML = `
            <p class="error-message">Failed to connect to server</p>
            <p>Error: ${error.message}</p>
            <p>Make sure the backend server is running on <code>${API_BASE_URL}</code></p>
        `;
    } finally {
        btn.disabled = false;
        btn.innerHTML = 'Check Server Status';
    }
}

// Fetch sample data
async function fetchData() {
    const dataBox = document.getElementById('data');
    const btn = document.getElementById('fetchDataBtn');
    
    btn.disabled = true;
    btn.innerHTML = '<span class="loading"></span> Fetching...';
    dataBox.innerHTML = '<p>Loading data...</p>';

    try {
        const response = await fetchAPI('/api/data');
        
        if (response.data && Array.isArray(response.data)) {
            let html = `<p class="success-message">${response.message}</p>`;
            response.data.forEach(item => {
                html += `
                    <div class="data-item">
                        <h3>${item.name}</h3>
                        <p><strong>ID:</strong> ${item.id}</p>
                        <p><strong>Description:</strong> ${item.description}</p>
                    </div>
                `;
            });
            dataBox.innerHTML = html;
        } else {
            dataBox.innerHTML = `<p>No data available</p>`;
        }
    } catch (error) {
        dataBox.className = 'data-box';
        dataBox.innerHTML = `
            <p class="error-message">Failed to fetch data</p>
            <p>Error: ${error.message}</p>
            <p>Make sure the backend server is running on <code>${API_BASE_URL}</code></p>
        `;
    } finally {
        btn.disabled = false;
        btn.innerHTML = 'Fetch Data';
    }
}

// Get API information
async function getApiInfo() {
    const infoBox = document.getElementById('api-info');
    const btn = document.getElementById('getApiInfoBtn');
    
    btn.disabled = true;
    btn.innerHTML = '<span class="loading"></span> Loading...';
    infoBox.innerHTML = '<p>Fetching API information...</p>';

    try {
        const data = await fetchAPI('/api');
        
        let html = `<p class="success-message">${data.message}</p>`;
        html += `<p><strong>Status:</strong> ${data.status}</p>`;
        html += `<p><strong>Version:</strong> ${data.data.version}</p>`;
        html += `<p><strong>Available Endpoints:</strong></p>`;
        html += '<ul style="margin-left: 20px; margin-top: 10px;">';
        data.data.endpoints.forEach(endpoint => {
            html += `<li><code>${API_BASE_URL}${endpoint}</code></li>`;
        });
        html += '</ul>';
        
        infoBox.innerHTML = html;
    } catch (error) {
        infoBox.className = 'info-box';
        infoBox.innerHTML = `
            <p class="error-message">Failed to fetch API information</p>
            <p>Error: ${error.message}</p>
            <p>Make sure the backend server is running on <code>${API_BASE_URL}</code></p>
        `;
    } finally {
        btn.disabled = false;
        btn.innerHTML = 'Get API Info';
    }
}

// Event listeners
document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('checkStatusBtn').addEventListener('click', checkServerStatus);
    document.getElementById('fetchDataBtn').addEventListener('click', fetchData);
    document.getElementById('getApiInfoBtn').addEventListener('click', getApiInfo);

    // Display welcome message
    console.log('Frontend loaded successfully!');
    console.log(`Backend API: ${API_BASE_URL}`);
});
