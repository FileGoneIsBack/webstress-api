function handleErrors() {
    console.log['error']('Failed to call the API', 'Error', { "toastClass": "toast-dark" });
}
document.addEventListener('DOMContentLoaded', () => {
    $.post('/api/dashboard/data', function (data) {
        if (typeof data === 'string') {
            data = JSON.parse(data);
        }
        loadTimeline(data.News)
        populateData(data);
    }).fail(handleErrors);
    $.post('/api/admin/user-list', function (data) {
        if (typeof data === 'string') {
            data = JSON.parse(data);
        }
        console.log(data)
    }).fail(handleErrors);
    // Initialize attacks chart
    const ctx = document.getElementById('attacksChart').getContext('2d');
    
    const gradient = ctx.createLinearGradient(0, 0, 0, 200);
    gradient.addColorStop(0, 'rgba(255, 51, 51, 0.3)');
    gradient.addColorStop(1, 'rgba(255, 51, 51, 0)');

    const attacksChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
            datasets: [{
                label: 'Total Attacks',
                data: [65000, 70000, 85000, 75000, 90000, 100000, 95000, 110000, 120000, 130000, 125000, 140000],
                borderColor: '#ff3333',
                backgroundColor: gradient,
                fill: true,
                tension: 0.4,
                borderWidth: 2,
                pointRadius: 0,
                pointHoverRadius: 4
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: false
                },
                tooltip: {
                    mode: 'index',
                    intersect: false,
                    backgroundColor: 'rgba(26, 26, 26, 0.9)',
                    titleColor: '#ff3333',
                    bodyColor: '#cccccc',
                    borderColor: 'rgba(255, 51, 51, 0.1)',
                    borderWidth: 1,
                    padding: 10,
                    displayColors: false
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    grid: {
                        color: 'rgba(255, 255, 255, 0.1)',
                        drawBorder: false
                    },
                    ticks: {
                        color: '#cccccc',
                        padding: 10,
                        callback: function(value) {
                            return value.toLocaleString();
                        }
                    }
                },
                x: {
                    grid: {
                        display: false
                    },
                    ticks: {
                        color: '#cccccc',
                        padding: 5
                    }
                }
            }
        }
    });

    let servers = [];
    function populateData(data) {
        const all = [...data.serversLayer4, ...data.serversLayer7];
        servers = all.map(s => ({
          name: s.name,
          type: s.Responsetime === -1 ? 'api' : 'attack',
          load: s.load || 0,
          status: s.load > 90 ? 'error' : s.load > 75 ? 'warning' : 'operational',
          slots: s.slots,
          connections: s.runningAttacks || 0,
          method: s.responsetime
        }));
    
        displayServers(document.querySelector('.toggle-btn.active')?.dataset.type || 'all');
    };

    // Function to create server cards
    const createServerCard = (server) => {
        const card = document.createElement('div');
        card.className = 'server-card';
        
        const statusClass = server.status === 'warning' ? 'warning' : 
                          server.status === 'error' ? 'error' : '';
        
        card.innerHTML = `
            <div class="server-name">
                ${server.name}
                ${server.type === 'api' ? 
                    `<span class="server-type">API</span>` : 
                    `<span class="server-type">${server.method}</span>`}
            </div>
            <div class="server-status">
                <span class="status-indicator ${statusClass}"></span>
                ${server.status.charAt(0).toUpperCase() + server.status.slice(1)}
            </div>
            <div class="server-stats">
                <div class="stat-row">
                    <div>
                        <span>Region</span>
                        <span>${server.slots}</span>
                    </div>
                </div>
                <div class="stat-row">
                    <div>
                        <span>Running Attacks</span>
                        <span>${server.connections.toLocaleString()}</span>
                    </div>
                </div>
                <div class="stat-row">
                    <div>
                        <span>Load</span>
                        <span>${server.load}%</span>
                    </div>
                    <div class="load-bar">
                        <div class="load-fill" style="width: ${server.load}%"></div>
                    </div>
                </div>
            </div>
        `;
        
        return card;
    };

    // Function to filter and display servers
    const displayServers = (type = 'all') => {
        const serverGrid = document.querySelector('.server-grid');
        serverGrid.style.opacity = '0';
        
        setTimeout(() => {
            serverGrid.innerHTML = '';
            const filteredServers = type === 'all' ? 
                servers : 
                servers.filter(server => server.type === type);
            
            filteredServers.forEach(server => {
                serverGrid.appendChild(createServerCard(server));
            });
            
            serverGrid.style.opacity = '1';
        }, 300);
    };

    // Initialize server display
    displayServers();

    // Handle server type toggle
    document.querySelector('.server-toggle').addEventListener('click', (e) => {
        if (e.target.classList.contains('toggle-btn')) {
            document.querySelectorAll('.toggle-btn').forEach(btn => {
                btn.classList.remove('active');
            });
            e.target.classList.add('active');
            displayServers(e.target.dataset.type);
        }
    });

    // Handle news form submission
    const newsForm = document.getElementById('newsForm');
    newsForm.addEventListener('submit', (e) => {
        e.preventDefault();
        
        const titleInput = newsForm.querySelector('input');
        const contentInput = newsForm.querySelector('textarea');
        const now = Math.floor(Date.now() / 1000); 
        var userData = {
            title: titleInput.value,
            from: "admin",
            content: contentInput.value,
            date: now
        };
        fetch('/api/admin/news', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        })
        /* const newsItem = document.createElement('div');
        newsItem.className = 'news-item';
        newsItem.innerHTML = `
            <div class="news-content-wrapper">
                <div class="news-title">${titleInput.value}</div>
                <div class="news-content">${contentInput.value}</div>
                <div class="news-date">Just now</div>
            </div>
            <button class="news-delete">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
            </button>
        `;
        newsList.insertBefore(newsItem, newsList.firstChild);
        */
        showToast('News update posted successfully', 'success');
        
        newsForm.reset();
    });

    /* Handle news deletion
    document.addEventListener('click', (e) => {
        if (e.target.closest('.news-delete')) {
            const newsItem = e.target.closest('.news-item');
            newsItem.style.opacity = '0';
            newsItem.style.transform = 'translateX(-20px)';
            setTimeout(() => {
                newsItem.remove();
                showToast('News update deleted', 'success');
            }, 300);
        }
    });
    */

    function loadTimeline(news) {
        const timelineList = $('#news-timeline');
        timelineList.empty();
    
        news.forEach(function (item) {
            const timelineItem = renderTimeline(item);
            timelineList.append(timelineItem);
        });
    }

    function renderTimeline(item) {
        const formattedDate = formatDate(item.Date); 
        return `
            <div class="news-item">
                <div class="news-content-wrapper">
                    <div class="news-title">${item.Title}</div>
                    <div class="news-content">${item.Content}</div>
                    <div class="news-date">${formattedDate}</div>
                </div>
                <button class="news-delete">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <line x1="18" y1="6" x2="6" y2="18"></line>
                        <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                </button>
            </div>
        `;
    }
});

function formatDate(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric',
        hour: '2-digit', minute: '2-digit'
    });
}