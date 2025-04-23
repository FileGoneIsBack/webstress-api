class DashboardLoader {
    constructor() {
        this.createLoader();
        this.progress = 0;
        this.loadingSteps = [
            'Initializing dashboard',
            'Loading components',
            'Fetching data',
            'Preparing interface'
        ];
        this.currentStep = 0;
    }
    

    createLoader() {
        const loader = document.createElement('div');
        loader.className = 'loader-container';
        loader.innerHTML = `
            <div class="loader-logo">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                </svg>
            </div>
            <div class="loader-bar">
                <div class="loader-progress"></div>
            </div>
            <div class="loader-percentage">0%</div>
            <div class="loader-text">Initializing dashboard</div>
        `;
        document.body.appendChild(loader);
        this.loader = loader;
        this.progressBar = loader.querySelector('.loader-progress');
        this.percentageText = loader.querySelector('.loader-percentage');
        this.loadingText = loader.querySelector('.loader-text');
    }

    updateProgress(progress) {
        this.progress = Math.min(progress, 100);
        this.progressBar.style.width = `${this.progress}%`;
        this.percentageText.textContent = `${Math.round(this.progress)}%`;
        
        // Update loading step text
        const stepIndex = Math.floor((this.progress / 100) * this.loadingSteps.length);
        if (stepIndex < this.loadingSteps.length && stepIndex !== this.currentStep) {
            this.currentStep = stepIndex;
            this.loadingText.textContent = this.loadingSteps[stepIndex];
        }
    }
    advance(stepText, percent) {
        this.progress = percent;
        this.progressBar.style.width = `${percent}%`;
        this.percentageText.textContent = `${percent}%`;
        this.loadingText.textContent = stepText;
    }
    simulateLoading() {
        return new Promise((resolve) => {
            let progress = 0;
            const interval = setInterval(() => {
                progress += Math.random() * 15;
                this.updateProgress(progress);
                
                if (progress >= 100) {
                    clearInterval(interval);
                    setTimeout(() => {
                        this.hide();
                        resolve();
                    }, 500);
                }
            }, 200);
        });
    }

    hide() {
        this.loader.classList.add('fade-out');
        setTimeout(() => {
            this.loader.remove();
        }, 500);
    }
}

