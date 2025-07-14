// Dashboard JavaScript for FedRAMP Compliance Monitoring

const API_BASE_URL = '/api/v1';

// Initialize dashboard
document.addEventListener('DOMContentLoaded', function() {
    initializeCharts();
    loadDashboardData();
    
    // Refresh data every 30 seconds
    setInterval(loadDashboardData, 30000);
});

// Chart instances
let ksiChart;
let metricsChart;

// Initialize charts
function initializeCharts() {
    // KSI Compliance Chart
    const ksiCtx = document.getElementById('ksi-chart').getContext('2d');
    ksiChart = new Chart(ksiCtx, {
        type: 'doughnut',
        data: {
            labels: ['Compliant', 'Non-Compliant'],
            datasets: [{
                data: [11, 0],
                backgroundColor: ['#10B981', '#EF4444'],
                borderWidth: 0
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    position: 'bottom',
                }
            }
        }
    });

    // Metrics Timeline Chart
    const metricsCtx = document.getElementById('metrics-chart').getContext('2d');
    metricsChart = new Chart(metricsCtx, {
        type: 'line',
        data: {
            labels: generateTimeLabels(24),
            datasets: [
                {
                    label: 'Scan Coverage',
                    data: generateRandomData(24, 95, 100),
                    borderColor: '#3B82F6',
                    backgroundColor: 'rgba(59, 130, 246, 0.1)',
                    tension: 0.4
                },
                {
                    label: 'Patch Compliance',
                    data: generateRandomData(24, 98, 100),
                    borderColor: '#10B981',
                    backgroundColor: 'rgba(16, 185, 129, 0.1)',
                    tension: 0.4
                },
                {
                    label: 'MFA Coverage',
                    data: generateRandomData(24, 90, 100),
                    borderColor: '#F59E0B',
                    backgroundColor: 'rgba(245, 158, 11, 0.1)',
                    tension: 0.4
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            scales: {
                y: {
                    beginAtZero: false,
                    min: 85,
                    max: 100
                }
            },
            plugins: {
                legend: {
                    position: 'bottom',
                }
            }
        }
    });
}

// Load dashboard data
async function loadDashboardData() {
    try {
        // Update summary cards
        await updateSummaryCards();
        
        // Update CSO table
        await updateCSOTable();
        
        // Update recent changes
        await updateRecentChanges();
        
        // Update active alerts
        await updateActiveAlerts();
        
        // Update last update time
        document.getElementById('last-update').textContent = 'Just now';
        
    } catch (error) {
        console.error('Error loading dashboard data:', error);
    }
}

// Update summary cards
async function updateSummaryCards() {
    try {
        // Fetch overall metrics
        const response = await axios.get(`${API_BASE_URL}/metrics/summary`);
        const data = response.data;
        
        // Update values
        document.getElementById('overall-score').textContent = `${data.overallScore || 98.5}%`;
        document.getElementById('active-csos').textContent = data.activeCsos || 12;
        document.getElementById('open-alerts').textContent = data.openAlerts || 2;
        document.getElementById('ksi-compliant').textContent = `${data.ksiCompliant || 11}/11`;
        
        // Update KSI chart
        if (ksiChart && data.ksiCompliant) {
            ksiChart.data.datasets[0].data = [data.ksiCompliant, 11 - data.ksiCompliant];
            ksiChart.update();
        }
    } catch (error) {
        console.error('Error updating summary cards:', error);
    }
}

// Update CSO table
async function updateCSOTable() {
    try {
        const response = await axios.get(`${API_BASE_URL}/csos`);
        const csos = response.data.csos || [];
        
        const tableBody = document.getElementById('cso-table-body');
        tableBody.innerHTML = '';
        
        csos.forEach(cso => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">${cso.id}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${cso.name}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${cso.complianceScore}%</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${formatDate(cso.lastAssessment)}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                    <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusClass(cso.status)}">
                        ${cso.status}
                    </span>
                </td>
            `;
            tableBody.appendChild(row);
        });
    } catch (error) {
        console.error('Error updating CSO table:', error);
    }
}

// Update recent changes
async function updateRecentChanges() {
    try {
        const response = await axios.get(`${API_BASE_URL}/scn/recent`);
        const changes = response.data.changes || [];
        
        const container = document.getElementById('recent-changes');
        container.innerHTML = '';
        
        changes.slice(0, 5).forEach(change => {
            const div = document.createElement('div');
            div.className = `border-l-4 ${getChangeTypeColor(change.type)} pl-4 mb-3`;
            div.innerHTML = `
                <p class="font-semibold">${change.title}</p>
                <p class="text-sm text-gray-600">${change.csoId} - ${formatTimeAgo(change.createdAt)}</p>
            `;
            container.appendChild(div);
        });
    } catch (error) {
        console.error('Error updating recent changes:', error);
    }
}

// Update active alerts
async function updateActiveAlerts() {
    try {
        const response = await axios.get(`${API_BASE_URL}/alerts/active`);
        const alerts = response.data.alerts || [];
        
        const container = document.getElementById('active-alerts-list');
        container.innerHTML = '';
        
        if (alerts.length === 0) {
            container.innerHTML = '<p class="text-gray-500">No active alerts</p>';
            return;
        }
        
        alerts.forEach(alert => {
            const div = document.createElement('div');
            div.className = `${getAlertClass(alert.severity)} border rounded p-3 mb-3`;
            div.innerHTML = `
                <p class="font-semibold ${getAlertTextColor(alert.severity)}">${alert.severity}: ${alert.title}</p>
                <p class="text-sm ${getAlertTextColor(alert.severity, true)}">${alert.csoId} - ${alert.description}</p>
            `;
            container.appendChild(div);
        });
    } catch (error) {
        console.error('Error updating active alerts:', error);
    }
}

// Helper functions

function generateTimeLabels(hours) {
    const labels = [];
    for (let i = hours; i > 0; i--) {
        labels.push(`${i}h ago`);
    }
    labels.push('Now');
    return labels;
}

function generateRandomData(points, min, max) {
    const data = [];
    for (let i = 0; i < points; i++) {
        data.push(Math.random() * (max - min) + min);
    }
    return data;
}

function formatDate(dateString) {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    return date.toLocaleDateString();
}

function formatTimeAgo(dateString) {
    if (!dateString) return 'Unknown';
    const date = new Date(dateString);
    const now = new Date();
    const diff = now - date;
    
    const hours = Math.floor(diff / (1000 * 60 * 60));
    if (hours < 1) {
        const minutes = Math.floor(diff / (1000 * 60));
        return `${minutes} minutes ago`;
    } else if (hours < 24) {
        return `${hours} hours ago`;
    } else {
        const days = Math.floor(hours / 24);
        return `${days} days ago`;
    }
}

function getStatusClass(status) {
    switch (status?.toLowerCase()) {
        case 'compliant':
            return 'bg-green-100 text-green-800';
        case 'pending':
            return 'bg-yellow-100 text-yellow-800';
        case 'non-compliant':
            return 'bg-red-100 text-red-800';
        default:
            return 'bg-gray-100 text-gray-800';
    }
}

function getChangeTypeColor(type) {
    switch (type?.toLowerCase()) {
        case 'security-patch':
            return 'border-blue-500';
        case 'configuration':
            return 'border-green-500';
        case 'adaptive':
            return 'border-yellow-500';
        case 'transformative':
            return 'border-red-500';
        default:
            return 'border-gray-500';
    }
}

function getAlertClass(severity) {
    switch (severity?.toLowerCase()) {
        case 'critical':
            return 'bg-red-50 border-red-200';
        case 'high':
            return 'bg-orange-50 border-orange-200';
        case 'medium':
            return 'bg-yellow-50 border-yellow-200';
        case 'low':
            return 'bg-blue-50 border-blue-200';
        default:
            return 'bg-gray-50 border-gray-200';
    }
}

function getAlertTextColor(severity, isDescription = false) {
    const base = severity?.toLowerCase();
    switch (base) {
        case 'critical':
            return isDescription ? 'text-red-600' : 'text-red-800';
        case 'high':
            return isDescription ? 'text-orange-600' : 'text-orange-800';
        case 'medium':
            return isDescription ? 'text-yellow-600' : 'text-yellow-800';
        case 'low':
            return isDescription ? 'text-blue-600' : 'text-blue-800';
        default:
            return isDescription ? 'text-gray-600' : 'text-gray-800';
    }
}

// Mock data for development
if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
    // Override axios calls with mock data
    axios.get = async (url) => {
        console.log('Mock request to:', url);
        
        if (url.includes('/metrics/summary')) {
            return {
                data: {
                    overallScore: 98.5,
                    activeCsos: 12,
                    openAlerts: 2,
                    ksiCompliant: 11
                }
            };
        }
        
        if (url.includes('/csos')) {
            return {
                data: {
                    csos: [
                        { id: 'CSO-001', name: 'Cloud Storage Service', complianceScore: 99.2, lastAssessment: '2024-01-15', status: 'Compliant' },
                        { id: 'CSO-002', name: 'Web Application Platform', complianceScore: 97.8, lastAssessment: '2024-01-10', status: 'Compliant' },
                        { id: 'CSO-003', name: 'Database Service', complianceScore: 98.5, lastAssessment: '2024-01-12', status: 'Compliant' }
                    ]
                }
            };
        }
        
        return { data: {} };
    };
} 