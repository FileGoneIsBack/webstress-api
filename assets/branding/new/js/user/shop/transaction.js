function handleErrors() {
    console.log['error']('Failed to call the API', 'Error', { "toastClass": "toast-dark" });
}
startCountdown();
var paymentIdCookie = document.cookie.split('; ').find(row => row.startsWith('payment_id='))?.split('=')[1];


$.post('/api/payments/status', { payment_id: paymentIdCookie }, function (response) {
    var data = $.parseJSON(response);
    console.log(data);
    if (data.status == "success") {
        console.log(data.message)
        updatePaymentStatus(data);
    } else if (data.status == "error") {
        console.log(data.message)
    } else {
        //toastr['error']('Failed to load transaction', 'Error', { "toastClass": "toast-dark" });
    }
}).fail(handleErrors);

function updatePaymentStatus(data) {
    const payment = data.payment_info;
    if (payment.status === 'finished') {
        window.location.href = '/dash'; //redirect to send noti later
        return; 
    }
    // Set status badge
    const statusBadge = document.querySelector('.status-badge');
    const statusMap = {
        waiting: 'Pending Payment',
        confirming: 'Waiting for confirmations',
        finished: 'Confirmed',
        expired: 'Payment Expired',
        partially_paid: 'Partially Paid'
    };
    statusBadge.textContent = statusMap[payment.status] || 'Unknown Status';
    statusBadge.className = `status-badge ${payment.status}`; 

    // Set BTC Address
    document.getElementById('btcAddress').textContent = payment.crypto_address;

    // Set USD amount
    document.querySelector('.detail-value.highlight').textContent = `$${payment.amount.toFixed(2)}`;

    // Set crypto amount + symbol
    const cryptoNames = {
        btc: 'BTC',
        eth: 'ETH',
        ltc: 'LTC',
        xmr: 'XMR',
        usdttrc20: 'USDT TRC20',
        usdterc20: 'USDT ERC20',
        bnbmainnet: 'BNB',
        trx: 'TRX'
    };
    const cryptoSymbol = cryptoNames[payment.crypto_coin.toLowerCase()] || payment.crypto_coin;
    document.querySelectorAll('.detail-value')[2].textContent = `${payment.crypto_amount} ${cryptoSymbol}`;

    // Set confirmation status
    const confirmationsText = `Waiting for payment (${payment.recieved || 0}/3 confirmations)`;
    document.querySelectorAll('.detail-value')[3].textContent = confirmationsText;

    // Format timestamps (UNIX to readable)
    const createdDate = new Date(payment.date * 1000); 
    const expiresDate = new Date((payment.date + 900) * 1000);

    document.querySelectorAll('.detail-value')[4].textContent = createdDate.toUTCString();
    document.querySelectorAll('.detail-value')[5].textContent = expiresDate.toUTCString();

    // Update QR code
    document.querySelector('.qr-code').src = payment.qr_code;
}

function startCountdown() {
    const countdownSpan = document.querySelector('.expire-countdown');
    let totalSeconds = 15 * 60; 
    let timer;

    function updateCountdown() {
        if (totalSeconds <= 0) {
            countdownSpan.textContent = '00:00';
            clearInterval(timer);
            return;
        }

        const minutes = Math.floor(totalSeconds / 60).toString().padStart(2, '0');
        const seconds = (totalSeconds % 60).toString().padStart(2, '0');
        countdownSpan.textContent = `${minutes}:${seconds}`;
        totalSeconds--;
    }

    updateCountdown(); 
    timer = setInterval(updateCountdown, 1000);
}

setInterval(function () {
    $.post('/api/payments/status', { payment_id: paymentIdCookie }, function (response) {
        var data = $.parseJSON(response);
        if (data.status == "success") {
            console.log(data.message)
            updatePaymentStatus(data);
        } else if (data.status == "error") {
            console.log(data.message)
        } else {
            console.log(data.message)
        }
    }).fail(handleErrors);
}, 1000);