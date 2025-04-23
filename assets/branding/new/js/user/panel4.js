document.addEventListener('DOMContentLoaded', () => {
    // Load attack methods on page load
    loadMethods();
    updateAttacks();

    // Handle sliders
    document.querySelectorAll('.slider').forEach(slider => {
        const valueDisplay = slider.nextElementSibling;
        
        slider.addEventListener('input', () => {
            const value = slider.value;
            const unit = slider.classList.contains('threads') ? ' threads' : ' PPS';
            valueDisplay.textContent = value + unit;
        });
    });

    // Handle form submission
    document.querySelector('.panel-form').addEventListener('submit', async (e) => {
        e.preventDefault();

        // Get form values
        const host = document.querySelector('input[placeholder="Enter IP address"]').value;
        const port = document.querySelector('input[placeholder="Enter port"]').value;
        const duration = document.querySelector('input[placeholder="Attack duration"]').value;
        const method = document.querySelector('.form-select').value;
        const threads = document.querySelector('input[placeholder="Enter threads"]').value;
        const concurrents = document.querySelector('input[placeholder="Enter conncurent attaks"]').value;
        //const threads = document.querySelector('.slider.threads').value;
        const pps = document.querySelector('.slider:not(.threads)').value;

        // Disable button while sending request
        const submitButton = e.target.querySelector('.submit-button');
        submitButton.disabled = true;

        try {
            const response = await fetch('/api/start', {
                method: 'POST',
                headers: { "Content-Type": "application/x-www-form-urlencoded" },
                body: new URLSearchParams({
                    host: host,
                    port: port,
                    duration: duration,
                    method: method,
                    threads:threads,
                    pps: pps,
                    concurrents: concurrents
                  })
            });

            const data = await response.json();

            if (data.status === "success") {
                showToast(data.message, 'success');
                addTargetToList(host, port, duration);
            } else {
                showToast(data.message || 'Failed to launch attack', 'error');
            }
        } catch (error) {
            showToast('Error connecting to server', 'error');
        }

        submitButton.disabled = false;
    });

    // Handle target removal
    document.querySelector('.target-list').addEventListener('click', (e) => {
        if (e.target.closest('.target-action')) {
            const targetItem = e.target.closest('.target-item');
            targetItem.remove();
            showToast('Target removed', 'success');
        }
    });
});

/**
 * Loads attack methods dynamically from the API and populates the dropdown.
 */
async function loadMethods() {
    try {
        const response = await fetch('/api/attacks/methods', { method: 'POST' });
        const data = await response.json();

        if (data.status === "success") {
            const methodSelect = document.querySelector('.form-select');
            methodSelect.innerHTML = '<option value="">Select method</option>';

            data.methods.forEach(method => {
                const option = document.createElement("option");
                option.value = method.method;
                option.textContent = method.panel_method;
                methodSelect.appendChild(option);
            });
        } else {
            showToast('Failed to load attack methods', 'warning');
        }
    } catch (error) {
        showToast('Error fetching methods', 'error');
    }
}

/**
 * Adds a new attack target to the list.
 */
function addTargetToList(host, port, duration) {
    const targetList = document.querySelector('.target-list');
    const targetCount = targetList.children.length + 1;

    const targetItem = document.createElement('div');
    targetItem.className = 'target-item';

    const id = `target-${targetCount}`;
    let remaining = parseInt(duration);

    targetItem.innerHTML = `
        <span class="target-number">#${targetCount}</span>
        <span class="target-info">${host}:${port}</span>
        <span class="target-time" id="${id}">${remaining}s</span>
        <span class="target-action">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
        </span>
    `;

    targetList.appendChild(targetItem);

    // Start countdown
    const interval = setInterval(() => {
        remaining--;
        const timeElement = document.getElementById(id);
        if (timeElement) {
            timeElement.textContent = `${remaining}s`;
        }

        if (remaining <= 0) {
            clearInterval(interval);
            targetItem.remove(); // remove the item when timer ends
        }
    }, 1000);
}

/**
 * Fetches and updates running attacks dynamically.
 */
async function updateAttacks() {
    try {
        const response = await fetch('/api/attacks/running', { method: 'POST' });
        const data = await response.json();

        if (data.status === "success") {
            const targetList = document.querySelector('.target-list');
            targetList.innerHTML = ''; // Clear the list before repopulating

            data.data.forEach(attack => {
                addTargetToList(attack.target, attack.port, attack.expires);
            });
        } else {
            showToast('Failed to load active attacks', 'warning');
        }
    } catch (error) {
        showToast('Error fetching attacks', 'error');
    }
}

/**
 * Stops an attack by ID.
 */
async function stopAttack(id) {
    try {
        const response = await fetch(`/api/stop/${id}`, { method: 'POST' });
        const data = await response.json();

        if (data.status === "success") {
            showToast('Attack stopped successfully', 'success');
            updateAttacks();
        } else {
            showToast(data.message || 'Failed to stop attack', 'warning');
        }
    } catch (error) {
        showToast('Error stopping attack', 'error');
    }
}

