(function() {
    $.post('/api/dashboard/data', function (data) {
        var json = $.parseJSON(data);
        console.log(json);
        updateCharts(json.networkInfo);
        populateData(json);
        populateL4Servers(json.serversLayer4);
        populateL7Servers(json.serversLayer7);
    }).fail(handleErrors);
})();

// Handle AJAX errors
function handleErrors(xhr, status, error) {
    console.error("AJAX Request Failed:", status, error);
    showToast("Failed to load dashboard data");
}

// Function to populate Layer 4 servers
function populateL4Servers(serversLayer4) {
    const serverContainerL4 = $('#serverContainerL4'); 
    serverContainerL4.empty();
    serversLayer4.forEach(server => {
        const serverStatusClass = server.status === 'Online' ? 'running' : 'stopped';

        let usedSlots = server.runningAttacks || 0; 
        let totalSlots = server.slots || 1; 

        const loadPercentage = (usedSlots / totalSlots) * 100;

        let methodsList = server.methods.map(method => `
            <li class="method-item">
                <span class="method-name">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M20 12V8H6a2 2 0 0 1-2-2V4"></path>
                        <path d="M4 12v4h14a2 2 0 0 1 2 2v2"></path>
                    </svg>
                    ${method}
                </span>
            </li>
        `).join(''); 

        const card = `
            <div class="server-item" data-id="server-${server.id}">
                <div class="server-header">
                    <div class="server-info">
                        <div class="server-name">${server.name}</div>
                        <div class="server-status">
                            <span class="status-indicator ${serverStatusClass}"></span>
                            <span>${server.status}</span>
                        </div>
                    </div>
                    <div class="server-load" data-slots="${usedSlots}/${totalSlots}">${usedSlots}/${totalSlots} slots</div>
                </div>
                <div class="server-load-bar">
                    <div class="load-bar-fill" style="width: ${loadPercentage}%"></div>
                </div>
                <div class="server-methods">
                    <ul class="method-list">
                        ${methodsList}
                    </ul>
                </div>
            </div>
        `;

        serverContainerL4.append(card); 
    });

    animateServerItems(); 
}

// Function to populate Layer 7 servers
function populateL7Servers(serversLayer7) {
    const serverContainerL7 = $('#serverContainerL7'); // Ensure this container exists
    serverContainerL7.empty(); // Clear previous servers

    serversLayer7.forEach(server => {
        const serverStatusClass = server.status === 'Online' ? 'running' : 'stopped';

        // Ensure `server.slots` is in the correct format
        let usedSlots = server.runningAttacks || 0; 
        let totalSlots = server.slots || 1; 

        const loadPercentage = (usedSlots / totalSlots) * 100;  // Calculate load percentage

        let methodsList = server.methods.map(method => `
            <li class="method-item">
                <span class="method-name">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M20 12V8H6a2 2 0 0 1-2-2V4"></path>
                        <path d="M4 12v4h14a2 2 0 0 1 2 2v2"></path>
                    </svg>
                    ${method}
                </span>
            </li>
        `).join(''); // Create the list of methods

        // Create the server card HTML
        const card = `
            <div class="server-item" data-id="server-${server.id}">
                <div class="server-header">
                    <div class="server-info">
                        <div class="server-name">${server.name}</div>
                        <div class="server-status">
                            <span class="status-indicator ${serverStatusClass}"></span>
                            <span>${server.status}</span>
                        </div>
                    </div>
                    <div class="server-load" data-slots="${usedSlots}/${totalSlots}">${usedSlots}/${totalSlots} slots</div>
                </div>
                <div class="server-load-bar">
                    <div class="load-bar-fill" style="width: ${loadPercentage}%"></div>
                </div>
                <div class="server-methods">
                    <ul class="method-list">
                        ${methodsList}
                    </ul>
                </div>
            </div>
        `;

        serverContainerL7.append(card);  // Append the card to the Layer 7 container
    });

    animateServerItems();  // Ensure animations apply to newly added items
}


// Populate page data (balance, membership expiry, etc.)
function populateData(data) {
    //$('#bal').text(data.userInfo.balance);
    //const expiry = formatDate(data.userInfo.membership_expire);
    //$('#expiry').text(expiry);
    $('#layer7_network').text(`${data.networkInfo.Layer7}/${data.networkInfo.Layer7Total} slots in use`);
    $('#layer4_network').text(`${data.networkInfo.Layer4}/${data.networkInfo.Layer4Total} slots in use`);
}

// Function to update charts with the data from your API
function updateCharts(networkInfo) {
    // Extract Layer4 and Layer7 values
    const layer4Data = networkInfo.Layer4 || 0;
    const layer7Data = networkInfo.Layer7 || 0;

    // Update Layer4 chart smoothly
    if (layer4Chart) {
        // Add the new data point
        layer4Chart.data.datasets[0].data.push(layer4Data);
        if (layer4Chart.data.datasets[0].data.length > 12) {  // Limit to the last 12 data points
            layer4Chart.data.datasets[0].data.shift();
        }

        // Smooth update with animation
        layer4Chart.update({
            duration: 750,  // Animation duration in ms
            easing: 'easeInOutQuart',  // Easing function for smooth transition
            lazy: true  // Lazy update for smooth transition
        });
    }

    // Update Layer7 chart smoothly
    if (layer7Chart) {
        // Add the new data point
        layer7Chart.data.datasets[0].data.push(layer7Data);
        if (layer7Chart.data.datasets[0].data.length > 12) {  // Limit to the last 12 data points
            layer7Chart.data.datasets[0].data.shift();
        }

        // Smooth update with animation
        layer7Chart.update({
            duration: 750,  // Animation duration in ms
            easing: 'easeInOutQuart',  // Easing function for smooth transition
            lazy: true  // Lazy update for smooth transition
        });
    }
}

// Chart setup code for Layer4 and Layer7 network usage with min and max Y-axis values
const sharedOptions = { 
    responsive: true,
    maintainAspectRatio: false,  // Make sure chart takes full space
    scales: {
        y: { 
            beginAtZero: true,
            min: 0,  // Set min value for Y-axis
            max: 100,  // Set max value for Y-axis
        }
    },
    animation: {
        duration: 750,  // Default duration for chart updates
        easing: 'easeInOutQuart'  // Default easing function for smooth transitions
    }
};

const layer4Ctx = document.getElementById('L4Chart').getContext('2d');
const layer4Gradient = layer4Ctx.createLinearGradient(0, 0, 0, 200);
layer4Gradient.addColorStop(0, 'rgba(51, 153, 255, 0.3)');
layer4Gradient.addColorStop(1, 'rgba(51, 153, 255, 0)');

const layer4Chart = new Chart(layer4Ctx, {
    type: 'line',
    data: {
        labels: Array.from({ length: 12 }, (_, i) => `${(11 - i) * 5}m`),
        datasets: [{
            data: [],  // This will be updated dynamically
            borderColor: '#3399ff',
            backgroundColor: layer4Gradient,
            fill: true,
            borderWidth: 2,
        }]
    },
    options: sharedOptions
});

// Layer7 chart
const layer7Ctx = document.getElementById('L7Chart').getContext('2d');
const layer7Gradient = layer7Ctx.createLinearGradient(0, 0, 0, 200);
layer7Gradient.addColorStop(0, 'rgba(255, 51, 51, 0.3)');
layer7Gradient.addColorStop(1, 'rgba(255, 51, 51, 0)');

const layer7Chart = new Chart(layer7Ctx, {
    type: 'line',
    data: {
        labels: Array.from({ length: 12 }, (_, i) => `${(11 - i) * 5}m`),
        datasets: [{
            data: [],  // This will be updated dynamically
            borderColor: '#ff3333',
            backgroundColor: layer7Gradient,
            fill: true,
        }]
    },
    options: sharedOptions
});

// Periodically update charts with new data (every 5 seconds in this example)
setInterval(() => {
    $.post('/api/dashboard/data', function(data) {
        var json = $.parseJSON(data);
        updateCharts(json.networkInfo);  // Pass network info to update the charts
    }).fail(handleErrors);
}, 5000);

const animateServerItems = () => {
    // Only target items that haven't been animated yet
    document.querySelectorAll('.server-item:not(.animated)').forEach((item, index) => {
        item.style.animation = `fadeInUp 0.3s ease-out ${index * 0.1}s forwards`;
        item.classList.add('animated'); // prevent reanimation

        // Attach click handler only once
        item.addEventListener('click', () => {
            item.classList.toggle('expanded');
        });
    });
};
