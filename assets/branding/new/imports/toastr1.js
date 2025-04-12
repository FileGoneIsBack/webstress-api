const injectCSS = () => {
    if (!document.getElementById('toastify-styles')) {
        const style = document.createElement('style');
        style.id = 'toastify-styles';
        style.innerHTML = `
            .toast-container {
                position: fixed;
                top: 1rem;
                right: 1rem;
                z-index: 9999;
                perspective: 1000px;
            }

            .toast-content {
                display: flex;
                align-items: center;
                gap: 12px;
                color: white;
                font-family: system-ui, -apple-system, sans-serif;
                background: rgba(26, 26, 26, 0.95);
                box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
                backdrop-filter: blur(10px);
                min-width: 300px;
                max-width: 400px;
                padding: 12px 24px;
                border-radius: 8px;
                border: 1px solid rgba(255, 51, 51, 0.2);
            }

            .toast-icon {
                display: flex;
                align-items: center;
                justify-content: center;
                padding: 8px;
                border-radius: 50%;
                background: rgba(255, 51, 51, 0.1);
            }

            .toast-icon svg {
                width: 24px;
                height: 24px;
                stroke: var(--primary);
                animation: pulse 2s infinite;
            }

            .toast-message {
                font-size: 0.95rem;
                font-weight: 500;
                line-height: 1.4;
                flex: 1;
                word-break: break-word;
            }

            .toastify {
                opacity: 0;
                transform: translateX(100%) rotateY(-45deg);
                animation: slideInRotate 0.5s cubic-bezier(0.68, -0.55, 0.265, 1.55) forwards;
                border-radius: 8px;
                padding: 0;
                background: transparent;
            }

            .toastify.on {
                opacity: 1;
                transform: translateX(0) rotateY(0);
            }

            .toastify:hover {
                transform: translateY(-3px) scale(1.02);
                transition: all 0.3s ease;
                cursor: pointer;
            }

            .toastify-progress {
                background: linear-gradient(to right, var(--primary), rgba(255, 51, 51, 0.5));
                height: 3px;
                bottom: 0;
                left: 0;
                opacity: 0.7;
                border-radius: 0 0 4px 4px;
                box-shadow: 0 0 8px var(--primary);
            }

            @keyframes slideInRotate {
                from {
                    opacity: 0;
                    transform: translateX(100%) rotateY(-45deg);
                }
                to {
                    opacity: 1;
                    transform: translateX(0) rotateY(0);
                }
            }

            @keyframes pulse {
                0% {
                    transform: scale(1);
                    opacity: 1;
                }
                50% {
                    transform: scale(1.1);
                    opacity: 0.7;
                }
                100% {
                    transform: scale(1);
                    opacity: 1;
                }
            }

            .toast-exit {
                animation: slideOutRotate 0.5s cubic-bezier(0.68, -0.55, 0.265, 1.55) forwards;
            }

            @keyframes slideOutRotate {
                from {
                    opacity: 1;
                    transform: translateX(0) rotateY(0);
                }
                to {
                    opacity: 0;
                    transform: translateX(-100%) rotateY(45deg);
                }
            }
        `;
        document.head.appendChild(style);
    }
};

const showToast = (message, type = 'info') => {
    injectCSS();
    
    const toastConfig = {
        success: {
            background: 'linear-gradient(135deg, rgba(26, 26, 26, 0.95), rgba(26, 26, 26, 0.98))',
            icon: `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                    <polyline points="22 4 12 14.01 9 11.01"></polyline>
                   </svg>`
        },
        error: {
            background: 'linear-gradient(135deg, rgba(26, 26, 26, 0.95), rgba(26, 26, 26, 0.98))',
            icon: `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="15" y1="9" x2="9" y2="15"></line>
                    <line x1="9" y1="9" x2="15" y2="15"></line>
                   </svg>`
        },
        warning: {
            background: 'linear-gradient(135deg, rgba(26, 26, 26, 0.95), rgba(26, 26, 26, 0.98))',
            icon: `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
                    <line x1="12" y1="9" x2="12" y2="13"></line>
                    <line x1="12" y1="17" x2="12.01" y2="17"></line>
                   </svg>`
        },
        info: {
            background: 'linear-gradient(135deg, rgba(26, 26, 26, 0.95), rgba(26, 26, 26, 0.98))',
            icon: `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="12" y1="16" x2="12" y2="12"></line>
                    <line x1="12" y1="8" x2="12.01" y2="8"></line>
                   </svg>`
        }
    };

    let toastContainer = document.querySelector('.toast-container');
    if (!toastContainer) {
        toastContainer = document.createElement('div');
        toastContainer.className = 'toast-container';
        document.body.appendChild(toastContainer);
    }

    Toastify({
        text: `<div class="toast-content">
                <div class="toast-icon">${toastConfig[type].icon}</div>
                <div class="toast-message">${message}</div>
               </div>`,
        duration: 3000,
        gravity: "top",
        position: "right",
        escapeMarkup: false,
        style: {
            background: toastConfig[type].background,
            boxShadow: '0 8px 24px rgba(0, 0, 0, 0.15)'
        },
        onClick: function() {
            const toast = this.toastElement;
            if (toast) {
                toast.classList.add('toast-exit');
                setTimeout(() => this.hideToast(), 500);
            }
        }
    }).showToast();
};

window.showSuccessToast = (message) => showToast(message, 'success');
window.showErrorToast = (message) => showToast(message, 'error');
window.showWarningToast = (message) => showToast(message, 'warning');
window.showInfoToast = (message) => showToast(message, 'info');

window.showToast = showToast;