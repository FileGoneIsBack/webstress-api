let plansFromServer = [];
let addonsFromServer = [];
let currentCategory = 'upgrades';
document.addEventListener('DOMContentLoaded', () => {
    $.post('/api/dashboard/plans', function (data) {
        try {
            const parsed = $.parseJSON(data);
    
            // ✅ Corrected assignment
            plansFromServer = Array.isArray(parsed.plans) ? parsed.plans : [];
            
            if (currentCategory === 'upgrades') renderPlans(currentCategory);
        } catch (err) {
            console.error("Failed to parse plans data:", err);
        }
    });
    
    $.post('/api/dashboard/addons', function (data) {
        try {
            const parsed = $.parseJSON(data);
    
            // ✅ Corrected assignment
            addonsFromServer = Array.isArray(parsed.addons) ? parsed.addons : [];
    
            if (currentCategory === 'addons') renderPlans(currentCategory);
        } catch (err) {
            console.error("Failed to parse addons data:", err);
        }
    });


    // Create toggle buttons
    const header = document.querySelector('h1');
    const toggleContainer = document.createElement('div');
    toggleContainer.className = 'toggle-container';
    toggleContainer.innerHTML = `
        <button class="toggle-button active" data-category="upgrades">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="display: inline-block; vertical-align: middle; margin-right: 8px;">
                <path d="M23 6l-9.5 9.5-5-5L1 18"></path>
                <path d="M17 6h6v6"></path>
            </svg>
            Upgrades
        </button>
        <button class="toggle-button" data-category="addons">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="display: inline-block; vertical-align: middle; margin-right: 8px;">
                <path d="M12 2v20M2 12h20"></path>
            </svg>
            Addons
        </button>
    `;
    header.after(toggleContainer);

    // Handle toggle clicks
    document.querySelectorAll('.toggle-button').forEach(button => {
        button.addEventListener('click', () => {
            const category = button.dataset.category;
            if (category === currentCategory) return;

            // Update active state
            document.querySelectorAll('.toggle-button').forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');

            // Update content
            currentCategory = category;
            renderPlans(category);
        });
    });


    // Add hover effect to feature items
    document.addEventListener('mouseover', (e) => {
        if (e.target.classList.contains('feature-item')) {
            e.target.style.transform = 'translateX(10px)';
            e.target.style.color = 'var(--primary)';
        }
    });

    document.addEventListener('mouseout', (e) => {
        if (e.target.classList.contains('feature-item')) {
            e.target.style.transform = 'translateX(0)';
            e.target.style.color = 'var(--text)';
        }
    });
});

const renderPlans = (category) => {
    const plansGrid = document.querySelector('.plans-grid');
    const plans = category === 'upgrades' ? plansFromServer : addonsFromServer;

    plansGrid.style.opacity = '0';

    setTimeout(() => {
        plansGrid.innerHTML = '';

        plans.forEach((plan, index) => {
            const planCard = document.createElement('div');
            planCard.className = 'plan-card';
            planCard.style.animationDelay = `${index * 0.1}s`;

            const features = [];
            if (category === 'upgrades') {
                features.push(`Duration: ${plan.duration} seconds`);
                features.push(`${plan.concurrents} Concurrent Attacks`);
                if (plan.api) features.push('API Access');
                if (plan.vip) features.push('VIP Support');
                expiryText = `${plan.expiry} days from purchase`;
            } else {
                features.push(plan.text || 'Addon Feature');
                expiryText = 'Never Expires'; 
            }

            planCard.innerHTML = `
                <div class="plan-header">
                    <h2>${plan.name}</h2>
                    <div class="plan-price">
                        <span class="currency">$</span>
                        <span class="amount">${plan.price}</span>
                        <span class="period">/once</span>
                    </div>
                </div>
                <div class="plan-features">
                    ${features.map(f => `
                        <div class="feature-item">
                            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <polyline points="20 6 9 17 4 12"></polyline>
                            </svg>
                            <span>${f}</span>
                        </div>
                    `).join('')}
                </div>
                <div class="plan-expire">
                    Expires: <span>${expiryText}</span>
                </div>
                <button class="plan-button">Select ${category === 'upgrades' ? 'Plan' : 'Addon'}</button>
            `;
            planCard.querySelector('.plan-button').addEventListener('click', () => {
                buyPlan(plan.name, category);  // Pass the plan name and category (either 'upgrades' or 'addons')
            });

            plansGrid.appendChild(planCard);
        });

        plansGrid.style.opacity = '1';
    }, 300);
};


function buyPlan(type, category) {
    if (category === 'upgrades') {
        // Confirm the action for the plan (upgrades)
        Swal.fire({
            icon: 'question',
            iconColor: '#FF3333',
            text: 'Are you sure you want to buy this membership?',
            showCancelButton: true,
            denyButtonText: `Buy Now`,
            background: '#1A1A1A',
            color: '#fff'
        }).then((result) => {
            if (result.isConfirmed) {
                // If confirmed, post to the plan endpoint
                $.post('/api/payments/buy', { plan_name: type }, function (data) {
                    var response = $.parseJSON(data);
                    if (response.status === "success") {
                        showToast('Payment successful, your membership was updated', 'success');
                    } else if (response.status === "error") {
                        showToast(response.message, 'warning');
                    } else {
                        showToast('Failed to process payment!', 'error');
                    }
                }).fail(handleErrors);
            }
        });
    } else if (category === 'addons') {
        // Confirm the action for the addon
        Swal.fire({
            icon: 'question',
            iconColor: '#FF3333',
            text: 'Are you sure you want to buy this addon?',
            showCancelButton: true,
            denyButtonText: `Buy Now`,
            background: '#1A1A1A',
            color: '#fff'
        }).then((result) => {
            if (result.isConfirmed) {
                // If confirmed, post to the addon endpoint
                $.post('/api/payments/buyaddon', { addon_name: type }, function (data) {
                    var response = $.parseJSON(data);
                    if (response.status === "success") {
                        showToast('Payment successful, your addon was updated', 'success');
                    } else if (response.status === "error") {
                        showToast(response.message, 'warning');
                    } else {
                        showToast('Failed to process payment!', 'error');
                    }
                }).fail(handleErrors);
            }
        });
    }
}
