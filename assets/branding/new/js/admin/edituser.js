let currentPage = 1;
const usersPerPage = 5;
let usersData = [];  // Make sure this is declared in a scope accessible by all functions

document.addEventListener('DOMContentLoaded', () => {
    // Load users from the server and render list/pagination.
    $.post('/api/admin/user-list', function (data) {
    if(data.status === "success"){
        // Save the users data into a global variable so the user list can be rendered later.
        usersData = data.users;
    }
    
    // Populate the user select dropdown using renderUserOption
    const UsersList = $('#userSelect');
    data.users.forEach(function(user) {
        const UserItem = renderUserOption(user); // Use the new function name here.
        UsersList.append(UserItem);
    });

    // Continue with your code to render pagination and display the user list...
    loadUsers(data);
});

    function loadUsers(data) {
        if (data.status == "success") {
            renderPagination();
            displayUsers(currentPage); 
        }
    }

    function renderUserOption(user) {
     return `<option value="${user.username}">${user.username}</option>`;
    }
    function renderPagination() {
        const totalPages = Math.ceil(usersData.length / usersPerPage);
        const paginationContainer = $(".pagination");
        paginationContainer.empty();

        // Previous button
        paginationContainer.append(`
            <button class="page-button prev ${currentPage === 1 ? 'disabled' : ''}" 
                ${currentPage === 1 ? 'disabled' : ''} 
                onclick="changePage(${currentPage - 1})">
                &laquo; Prev
            </button>
        `);

        // Page numbers
        for (let i = 1; i <= totalPages; i++) {
            paginationContainer.append(`
                <button class="page-button ${i === currentPage ? 'active' : ''}" 
                    onclick="changePage(${i})">
                    ${i}
                </button>
            `);
        }

        // Next button
        paginationContainer.append(`
            <button class="page-button next ${currentPage === totalPages ? 'disabled' : ''}" 
                ${currentPage === totalPages ? 'disabled' : ''} 
                onclick="changePage(${currentPage + 1})">
                Next &raquo;
            </button>
        `);
    }

    function displayUsers(page) {
        // Use the correct container ID for the updated HTML
        const usersList = $('#userList');
        usersList.empty();

        const startIndex = (page - 1) * usersPerPage;
        const endIndex = startIndex + usersPerPage;
        const usersToDisplay = usersData.slice(startIndex, endIndex);

        usersToDisplay.forEach(function (user) {
            const userItem = renderUser(user);
            usersList.append(userItem);
        });
    }

    function renderUser(user) {
    const rolesHTML = user.permissions && Array.isArray(user.permissions)
        ? user.permissions.map(role => `<span class="role-badge">${role.name}</span>`).join('')
        : '';
    return `
        <div class="user-item" data-user-id="${user.id}">
            <div class="user-info">
                <span class="user-id">${user.id}</span>
                <span class="user-name">${user.username}</span>
                <div class="user-roles">
                    ${rolesHTML}
                </div>
            </div>
        </div>
    `;
}

    // Called from the pagination buttons (via inline onclick)
    window.changePage = function(page) {
        const totalPages = Math.ceil(usersData.length / usersPerPage);
        if (page < 1 || page > totalPages) return;

        currentPage = page;
        displayUsers(currentPage);
        renderPagination();
    };

    // ======================
    // *** Copy User to Form ***
    // ======================
     $('#userList').on('click', '.user-item', function() {
        const userId = $(this).data('user-id').toString();
        const user = usersData.find(u => String(u.id) === userId);
        if (user) {
            populateUserForm(user);
        }
    });
    $('#userSelect').on('change', function() {
        const selectedUserId = $(this).val();
        if (!selectedUserId) return; // No user selected
        const user = usersData.find(u => String(u.id) === selectedUserId);
        if (user) {
            populateUserForm(user);
        }
    });

    function populateUserForm(user) {
        console.log(user)
    $('#userSelect').val(user.username);
    $('#mbt').val(user.duration);

    $('#cons').val(user.conns);
    $('#balance').val(user.balance);
    selectedCategories.clear();
    $('.select-option').removeClass('selected');

    if (user.permissions && Array.isArray(user.permissions)) {
        user.permissions.forEach(role => {
            selectedCategories.add(role.name);
            $(`.select-option[data-value="${role.name}"]`).addClass('selected');
        });
    }

    updateSelectedCategories();
}

    // ===============================
    // *** Multi-select Functionality ***
    // ===============================
    const select = document.querySelector('.form2-select');
    const input = document.querySelector('.select-input');
    const dropdown = document.querySelector('.select-dropdown');
    const selectedCategories = new Set();

    const updateSelectedCategories = () => {
        input.innerHTML = selectedCategories.size === 0 
            ? '<span class="select-placeholder">Select categories</span>'
            : Array.from(selectedCategories).map(category => `
                <span class="select-item">
                    ${category}
                    <span class="select-item-remove" data-value="${category}">×</span>
                </span>
            `).join('');
    };

    input.addEventListener('click', () => {
        dropdown.classList.toggle('active');
    });

    document.addEventListener('click', (e) => {
        if (!select.contains(e.target)) {
            dropdown.classList.remove('active');
        }
    });

    dropdown.addEventListener('click', (e) => {
        const option = e.target.closest('.select-option');
        if (option) {
            const value = option.dataset.value;
            if (selectedCategories.has(value)) {
                selectedCategories.delete(value);
                option.classList.remove('selected');
            } else {
                selectedCategories.add(value);
                option.classList.add('selected');
            }
            updateSelectedCategories();
        }
    });

    input.addEventListener('click', (e) => {
        if (e.target.classList.contains('select-item-remove')) {
            const value = e.target.dataset.value;
            selectedCategories.delete(value);
            document.querySelector(`.select-option[data-value="${value}"]`)
                .classList.remove('selected');
            updateSelectedCategories();
            e.stopPropagation();
        }
    });

    // ============================
    // *** Form Submission Handling ***
    // ============================
    // Ensure that the form element has the id "advancedForm"
    document.getElementById('advancedForm').addEventListener('submit', (e) => {
        e.preventDefault();
        const username = document.getElementById('userSelect').value;
        const concurrents = Number(document.getElementById('cons')?.value || 0);
        const servers = Number(0);
        const duration = Number(document.getElementById('mbt')?.value || 0);
        const roles = Array.from(selectedCategories);
        //const unixTimestamp = Math.floor(Date.now() / 1000);
        const balance = Number(document.getElementById('balance')?.value || 0);


        var userData = {
            username: username,
            concurrents: concurrents,
            servers: servers,
            duration: duration,
            ranks: roles.map(role => ({ name: role, has: true })),
            //expiry: unixTimestamp,
            balance: balance
        };
        console.log(userData);

        $.ajax({
    url: '/api/admin/update-user',
    type: 'POST',
    data: JSON.stringify(userData), // userData is your JavaScript object
    contentType: 'application/json', // Tell the server that the data is JSON
    success: function(response) {
        console.log(response);
    },
    error: function(xhr, status, error) {
        console.error("Error updating user:", error);
    }
});
    });

    // ================================
    // *** Cancel Button Handling ***
    // ================================
    document.querySelector('.secondary-button').addEventListener('click', () => {
        if (confirm('Are you sure you want to cancel? All changes will be lost.')) {
            window.location.href = '/dashboard';
        }
    });
});

    