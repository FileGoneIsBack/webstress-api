// Tab switching
const tabs = document.querySelectorAll('.tab');
const forms = document.querySelectorAll('.form');

function activateTab(tabName) {
    // Remove active class from all tabs and forms
    tabs.forEach(t => t.classList.remove('active'));
    forms.forEach(f => f.classList.remove('active'));
    
    // Add active class to selected tab and form
    const selectedTab = document.querySelector(`[data-tab="${tabName}"]`);
    const selectedForm = document.querySelector(`#${tabName}Form`);
    
    if (selectedTab && selectedForm) {
        selectedTab.classList.add('active');
        selectedForm.classList.add('active');
    }
}

// Handle URL fragments
function handleUrlFragment() {
    const hash = window.location.hash.slice(1).toLowerCase();
    if (hash === 'signup') {
        activateTab('signup');
    } else {
        activateTab('login');
    }
}

// Listen for hash changes
window.addEventListener('hashchange', handleUrlFragment);

// Handle initial load
handleUrlFragment();

// Tab click handling
tabs.forEach(tab => {
    tab.addEventListener('click', () => {
        const tabName = tab.dataset.tab;
        window.location.hash = tabName.toLowerCase();
    });
});

// Captcha
function generateCaptcha() {
    const num1 = Math.floor(Math.random() * 10);
    const num2 = Math.floor(Math.random() * 10);
    const operators = ["+", "-", "*"];
    const operator = operators[Math.floor(Math.random() * operators.length)];
    
    let answer;
    switch (operator) {
        case "+": answer = num1 + num2; break;
        case "-": answer = num1 - num2; break;
        case "*": answer = num1 * num2; break;
    }
    
    const captchaText = `${num1} ${operator} ${num2}`;
    document.getElementById("signup-captcha").value = captchaText;
    
    return answer.toString();
}

let currentCaptcha = '';

function refreshCaptcha() {
    currentCaptcha = generateCaptcha();
}

document.getElementById('refreshCaptcha').addEventListener('click', refreshCaptcha);
refreshCaptcha();

// Terms of Service
const tosText = `
Terms of Service

1. Acceptance of Terms
By accessing and using this service, you accept and agree to be bound by the terms and provision of this agreement.

2. User Account
You are responsible for maintaining the confidentiality of your account and password.

3. Privacy Policy
Your privacy is important to us. Please review our Privacy Policy to understand how we collect and use your information.

4. Termination
We reserve the right to terminate or suspend your account at any time without notice.

5. Changes to Terms
We reserve the right to modify these terms at any time. Please review these terms periodically for changes.
`;

document.getElementById('tosButton').addEventListener('click', () => {
    Swal.fire({
        title: 'Terms of Service',
        html: tosText.replace(/\n/g, '<br>'),
        width: 600,
        padding: '3em',
        color: '#1a1a1a',
        confirmButtonColor: '#ff3333'
    });
});

// Form Validation
document.getElementById("signupForm").addEventListener("submit", (e) => {
    e.preventDefault();
    
    const password = document.getElementById("signup-password").value;
    const confirmPassword = document.getElementById("signup-confirm-password").value;
    const captchaInput = document.getElementById("captchaInput").value;
    const tosCheckbox = document.getElementById("tos");

    if (captchaInput !== currentCaptcha) {
        showToast("Invalid CAPTCHA", "error");
        refreshCaptcha();
        document.getElementById("captchaInput").value = "";
        return;
    }

    if (!tosCheckbox.checked) {
        showToast("Please accept the Terms of Service", "error");
        return;
    }
    e.target.submit(); 
});

document.getElementById('loginForm').addEventListener('submit', (e) => {
    e.target.submit();
});