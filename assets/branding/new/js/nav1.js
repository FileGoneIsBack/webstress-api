// Handle user dropdown
const initUserDropdown = () => {
    const userButton = document.querySelector('.user-button');
    const userDropdown = document.querySelector('.user-dropdown');
    
    if (userButton && userDropdown) {
        userButton.addEventListener('click', (e) => {
            e.stopPropagation();
            userDropdown.classList.toggle('active');
            // Close notifications if open
            const notificationsDropdown = document.querySelector('.notifications-dropdown');
            if (notificationsDropdown) {
                notificationsDropdown.classList.remove('active');
            }
        });

        // Prevent dropdown close when clicking inside
        userDropdown.addEventListener('click', (e) => e.stopPropagation());
    }
};

// Handle notifications dropdown
const initNotificationsDropdown = () => {
    const notificationsButton = document.querySelector('.notifications');
    const notificationsDropdown = document.querySelector('.notifications-dropdown');
    
    if (notificationsButton && notificationsDropdown) {
        notificationsButton.addEventListener('click', (e) => {
            e.stopPropagation();
            notificationsDropdown.classList.toggle('active');
            // Close user dropdown if open
            const userDropdown = document.querySelector('.user-dropdown');
            if (userDropdown) {
                userDropdown.classList.remove('active');
            }
        });

        // Prevent dropdown close when clicking inside
        notificationsDropdown.addEventListener('click', (e) => e.stopPropagation());
    }
};

// Close dropdowns when clicking outside
const initOutsideClickHandler = () => {
    document.addEventListener('click', () => {
        const userDropdown = document.querySelector('.user-dropdown');
        const notificationsDropdown = document.querySelector('.notifications-dropdown');
        
        if (userDropdown) {
            userDropdown.classList.remove('active');
        }
        if (notificationsDropdown) {
            notificationsDropdown.classList.remove('active');
        }
    });
};

// Mobile menu toggle
const initMobileMenu = () => {
    const menuToggle = document.querySelector('.menu-toggle');
    const sideNav = document.querySelector('.side-nav');
    
    if (menuToggle && sideNav) {
        menuToggle.addEventListener('click', () => {
            sideNav.classList.toggle('active');
        });
    }
};

// Handle navigation dropdowns
const initNavDropdowns = () => {
    const navItems = document.querySelectorAll('.nav-item');
    
    navItems.forEach(item => {
        const dropdown = item.nextElementSibling;
        if (dropdown && dropdown.classList.contains('nav-dropdown')) {
            item.addEventListener('click', (e) => {
                e.preventDefault();
                
                // Close other dropdowns
                navItems.forEach(otherItem => {
                    if (otherItem !== item) {
                        otherItem.classList.remove('expanded');
                        const otherDropdown = otherItem.nextElementSibling;
                        if (otherDropdown && otherDropdown.classList.contains('nav-dropdown')) {
                            otherDropdown.style.maxHeight = '0px';
                        }
                    }
                });
                
                // Toggle current dropdown
                item.classList.toggle('expanded');
                if (item.classList.contains('expanded')) {
                    dropdown.style.maxHeight = `${dropdown.scrollHeight}px`;
                } else {
                    dropdown.style.maxHeight = '0px';
                }
            });
        }
    });

    // Initialize dropdowns to closed state
    document.querySelectorAll('.nav-dropdown').forEach(dropdown => {
        dropdown.style.maxHeight = '0px';
    });
};

// Set active states based on current page
const initActiveStates = () => {
    const currentPath = window.location.pathname;
    document.querySelectorAll('.nav-dropdown-item').forEach(item => {
        if (item.getAttribute('href') === currentPath) {
            item.classList.add('active');
            const parentDropdown = item.closest('.nav-dropdown');
            const parentNavItem = parentDropdown?.previousElementSibling;
            if (parentNavItem) {
                parentNavItem.classList.add('expanded');
                parentDropdown.style.maxHeight = `${parentDropdown.scrollHeight}px`;
            }
        }
    });
};


// Initialize all navigation functionality
const initNavigation = () => {
    initUserDropdown();
    initNotificationsDropdown();
    initOutsideClickHandler();
    initMobileMenu();
    initNavDropdowns();
    initActiveStates();
};

// Export the initialization function
window.initNavigation = initNavigation;

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', initNavigation);
