function GetKey() {
    $.ajax({
        url: '/api/manager/api',  // The endpoint from which you're fetching data
        type: 'GET',              // HTTP method: GET
        dataType: 'json',         // Expected response type
        success: function(response) {
            // This function will run on successful response
            console.log('Data fetched:', response);
            $('#apiKeyField').text(response.key); // Example: inserting the API key into a field
        },
        error: function(xhr, status, error) {
            // This function will run if there's an error with the request
            console.error('Error fetching data:', error);
        }
    });
}

function RegenKey() {
    $.post('/api/manager/api', function(response) {
        // Success callback
        console.log('API key regenerated:', response);
        
        // Update the UI with the new API key
        $('#apiKeyField').text(response.new_key);

        // Optionally: Update the copy button data attribute with the new key
        const copyButton = document.querySelector('.copy-button[data-value]');
        if (copyButton) {
            copyButton.dataset.value = response.new_key;
        }

        // Show success message (e.g., a toast)
        showToast('API key regenerated successfully', 'success');

        // Add the regeneration event to the activity list
        const activityList = document.querySelector('.activity-list');
        if (activityList) {
            const newActivity = document.createElement('div');
            newActivity.className = 'activity-item';
            newActivity.innerHTML = `
                <div class="activity-icon success">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="20 6 9 17 4 12"></polyline>
                    </svg>
                </div>
                <div class="activity-content">
                    <div class="activity-title">API Key Regenerated</div>
                    <div class="activity-time">Just now</div>
                </div>
                <div class="activity-status success">Success</div>
            `;
            activityList.insertBefore(newActivity, activityList.firstChild);
        }
    }).fail(function(xhr, status, error) {
        // Error callback
        console.error('Error regenerating API key:', error);
        showToast('Failed to regenerate API key', 'error');
    });
}

document.addEventListener('DOMContentLoaded', () => {
    key = GetKey();
    const apiKey = key;
    let isKeyVisible = false;

    // Handle API key visibility toggle
    const visibilityToggle = document.querySelector('.visibility-toggle');
    const apiKeyField = document.querySelector('.api-key-field');

    if (visibilityToggle && apiKeyField) {
        visibilityToggle.addEventListener('click', () => {
            isKeyVisible = !isKeyVisible;
            apiKeyField.textContent = isKeyVisible ? apiKey : '••••••••••••••••••••••••••••••';
            
            // Update icon
            visibilityToggle.innerHTML = isKeyVisible 
                ? `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path>
                    <line x1="1" y1="1" x2="23" y2="23"></line>
                   </svg>`
                : `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                    <circle cx="12" cy="12" r="3"></circle>
                   </svg>`;
        });
    }

    // Handle copy functionality for individual API key or other data
    const copyButtons = document.querySelectorAll('.copy-button');
    copyButtons.forEach(button => {
        button.addEventListener('click', async () => {
            const textToCopy = button.dataset.value;
            try {
                await navigator.clipboard.writeText(textToCopy);
                showToast('Copied to clipboard', 'success');
            } catch (err) {
                showToast('Failed to copy text', 'error');
            }
        });
    });

    // Handle copy all credentials functionality
    const copyAllButton = document.getElementById('copyAll');
    if (copyAllButton) {
        copyAllButton.addEventListener('click', async () => {
            const endpoint = 'https://api.secureauth.pro/v1';
            const textToCopy = `API Endpoint: ${endpoint}\nAPI Key: ${apiKey}`;
            
            try {
                await navigator.clipboard.writeText(textToCopy);
                showToast('All credentials copied to clipboard', 'success');
            } catch (err) {
                showToast('Failed to copy credentials', 'error');
            }
        });
    }

    // Handle API key regeneration (with confirmation dialog)
    const regenerateButton = document.querySelector('.regenerate-button');
    if (regenerateButton) {
        regenerateButton.addEventListener('click', () => {
            const confirmRegenerate = confirm(
                'Are you sure you want to regenerate your API key? ' +
                'This will invalidate your current key immediately.'
            );
    
            if (confirmRegenerate) {
                RegenKey();
            }
        });
    }
    // Animate stats (bars) on load
    const chartBars = document.querySelectorAll('.chart-bar');
    chartBars.forEach(bar => {
        const width = bar.style.width;
        bar.style.width = '0';
        setTimeout(() => {
            bar.style.width = width;
        }, 100);
    });
});

// Fetch API Key from the server (via GET request)
document.addEventListener('DOMContentLoaded', async () => {
    const apiKeyField = document.querySelector('.api-key-field');

    if (!apiKeyField) {
        console.warn('.api-key-field not found in DOM.');
        return;
    }

    try {
        const response = await fetch('/api/manager/api', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (response.ok) {
            const data = await response.json();
            apiKeyField.textContent = data.key;
        } else {
            console.error('Failed to fetch API key');
        }
    } catch (error) {
        console.error('Error fetching API key:', error);
    }
});
// Regenerate API Key and update on the frontend (via POST request)
document.addEventListener('DOMContentLoaded', () => {
    const regenerateButton = document.querySelector('.regenerate-button');
    const apiKeyField = document.querySelector('.api-key-field');

    regenerateButton.addEventListener('click', async () => {
        const confirmRegenerate = confirm(
            'Are you sure you want to regenerate your API key? This will invalidate your current key immediately.'
        );

        if (confirmRegenerate) {
            try {
                const response = await fetch('/api/manager/api', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });

                if (response.ok) {
                    const data = await response.json();
                    const newKey = data.new_key; // Get the newly generated API key
                    apiKeyField.textContent = newKey; // Update the UI with the new key
                    showToast('API key regenerated successfully', 'success'); // Show success toast
                } else {
                    console.error('Failed to regenerate API key');
                }
            } catch (error) {
                console.error('Error regenerating API key:', error);
            }
        }
    });
});

