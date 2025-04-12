function handleErrors() {
    showToast('Failed to call the API', 'info');
}


document.addEventListener('DOMContentLoaded', () => {
    // Handle form submission
    document.querySelector('.panel-form').addEventListener('submit', (e) => {
        e.preventDefault();

        var amount= e.target.querySelector('input[type="number"]').value;	
        var coin= e.target.querySelector('select').value;

        // Show success message
        showToast(`Transaction is being generated...`, 'info');
        
        $.post('/api/payments/create', {amount, coin}, function(data) {
            var response = $.parseJSON(data);
            if (response.status == "success") {
                document.cookie = `payment_id=${response.id}; path=/`;
    
                window.location.href = 'transaction';
            } else if (response.status == "error") {
                showToast(response.message, 'error');
            } else {
                showToast('Failed to create transaction', 'error');
            }
        }).fail(handleErrors);

        targetList.insertBefore(targetItem, targetList.firstChild);
        
        // Clear form
        e.target.reset();
    });

    // Mobile menu toggle
    const menuToggle = document.querySelector('.menu-toggle');
    const sideNav = document.querySelector('.side-nav');
    
    menuToggle?.addEventListener('click', () => {
        sideNav.classList.toggle('active');
    });
});


let previousPayments = '';

function updatePayments() {
    $.post('/api/payments/history', function (data) {
        var response = $.parseJSON(data);
        if (response.status === "success" && Array.isArray(response.payments)) {
            const currentPayments = JSON.stringify(response.payments);

            // Check if anything has changed
            if (currentPayments === previousPayments) {
                return;
            }

            previousPayments = currentPayments; 
            const paymentInfoArray = response.payments;
            const transactionsTableBody = document.getElementById("target-list");

            transactionsTableBody.innerHTML = '';

            paymentInfoArray.forEach((payment) => {
                const newRow = document.createElement('div');

                const cryptoNames = {
                    BTC: 'BTC',
                    eth: 'ethereum',
                    ltc: 'litecoin',
                    xmr: 'monero',
                    usdttrc20: 'USDT TRC20',
                    usdterc20: 'USDT ERC20',
                    bnbmainnet: 'Binance Coin Mainnet',
                    trx: 'TRON'
                };

                const cryptoCoin = payment.coin;
                const cryptoName = cryptoNames[cryptoCoin] || 'Unknown Crypto';

                const status = {
                    waiting: `Waiting for transaction`,
                    confirming: `Waiting for confirmations`,
                    finished: `Confirmed`,
                    expired: `Payment Expired`,
                    partially_paid: `Partially Paid`,
                }[payment.status] || 'Unknown';

                newRow.className = 'target-item';
                newRow.innerHTML = `
                    <span class="target-number">#${payment.id}</span>
                    <span class="target-info">$${payment.amount} ${cryptoName}</span>
                    <span class="target-time">${status}</span>
                    <span class="target-action">
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                            <polyline points="7 10 12 15 17 10"></polyline>
                            <line x1="12" y1="15" x2="12" y2="3"></line>
                        </svg>
                    </span>
                `;

                newRow.addEventListener("click", () => {
                    document.cookie = `payment_id=${payment.id}; path=/`;
                    window.location.href = 'transaction';
                });

                transactionsTableBody.prepend(newRow);
            });
        } else if (response.status === "error") {
            showToast(response.message || 'Something went wrong', 'error');
        } else {
            showToast('Failed to load transaction', 'error');
        }
    }).fail(handleErrors);
}

updatePayments();

setInterval(function () {
    updatePayments();
}, 3000);
