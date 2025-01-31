function handleErrors() {
    toastr['error']('Failed to call the API', 'Error', { "toastClass": "toast-dark" });
}
var currentPage = 1;
var usersPerPage = 9; // This should limit the number of users displayed per page
var usersData = []; // This will hold all the users' data

(function() {
    // Load users for the first page
    loadUsers(currentPage);
})();

// Render user option for selection
function renderUser(user) {
    return `<option value="${user.username}" selected>${user.username}</option>`;
}

// Render user in table format
function renderUserTable(user) {
    return `<tr>
        <td>${user.id}</td>
        <td>${user.username}</td>
        <td>${user.conns}</td>
        <td>${user.servers}</td>
        <td>${user.duration}</td>
        <td>${user.permissions}</td>
        <td>
            <div class="dropdown">
                <button type="button" class="btn p-0 dropdown-toggle hide-arrow" data-bs-toggle="dropdown"><i class="fa fa-ellipsis-v fa-2x"></i></button>
                <div class="dropdown-menu">
                    <a class="dropdown-item" href="javascript:void(0);" onclick="editUser('${user.username}')"><i class="ti ti-pencil me-1"></i>Edit</a>
                    <a class="dropdown-item" href="javascript:void(0);" onclick="deleteUser('${user.username}')"><i class="ti ti-trash me-1"></i>Delete</a>
                </div>
            </div>
        </td>
    </tr>`;
}

// Load users and handle pagination
function loadUsers(page) {
    $.post('/api/admin/user-list', { page: page, perPage: usersPerPage }, function(data) {
        if (data.status === "success") {
            usersData = data.users; // Store the users data

            // Render the user list
            const UsersList = $('#userOption');
            UsersList.empty(); // Clear previous options
            data.users.forEach(function(user) {
                const UserItem = renderUser(user);
                UsersList.append(UserItem);
            });

            // Render pagination and handle pagination links
            const totalPages = Math.ceil(usersData.length / usersPerPage);
            displayUsers(currentPage); 
            renderPagination(totalPages);
        } else {
            handleErrors();
        }
    });
}

function renderPagination() {
    const totalPages = Math.ceil(usersData.length / usersPerPage);
    const paginationContainer = $(".pagination");
    paginationContainer.empty(); 

    paginationContainer.append(`<li class="page-item prev ${currentPage === 1 ? 'disabled' : ''}">
                        <a class="page-link" href="#" onclick="changePage(${currentPage - 1})"><i class="ti ti-chevrons-left ti-xs"></i></a>
                    </li>`);

    for (let i = 1; i <= totalPages; i++) {
        paginationContainer.append(`<li class="page-item ${i === currentPage ? 'active' : ''}">
                            <a class="page-link" href="#" onclick="changePage(${i})">${i}</a>
                        </li>`);
    }

    paginationContainer.append(`<li class="page-item next ${currentPage === totalPages ? 'disabled' : ''}">
                        <a class="page-link" href="#" onclick="changePage(${currentPage + 1})"><i class="ti ti-chevrons-right ti-xs"></i></a>
                    </li>`);
}

// Update pagination based on current page
function updatePagination(totalPages, currentPage) {
    const paginationList = $('.pagination');
    paginationList.empty(); // Clear previous pagination links

    // Previous page link
    paginationList.append(`
        <li class="page-item ${currentPage === 1 ? 'disabled' : ''}">
            <a class="page-link" href="#" onclick="changePage(${currentPage - 1})">
                <i class="ti ti-chevrons-left ti-xs"></i>
            </a>
        </li>
    `);

    // Page number links
    for (let i = 1; i <= totalPages; i++) {
        paginationList.append(`
            <li class="page-item ${i === currentPage ? 'active' : ''}">
                <a class="page-link" href="#" onclick="changePage(${i})">${i}</a>
            </li>
        `);
    }

    // Next page link
    paginationList.append(`
        <li class="page-item ${currentPage === totalPages ? 'disabled' : ''}">
            <a class="page-link" href="#" onclick="changePage(${currentPage + 1})">
                <i class="ti ti-chevrons-right ti-xs"></i>
            </a>
        </li>
    `);
}
function displayUsers(page) {
    const UsersList = $('#userTable');
    UsersList.empty(); 

    const startIndex = (page - 1) * usersPerPage;
    const endIndex = startIndex + usersPerPage;
    const usersToDisplay = usersData.slice(startIndex, endIndex);

    usersToDisplay.forEach(function (user) {
        const UserItem = renderUserTable(user);
        UsersList.append(UserItem);
    });
}

// Change page based on pagination click
function changePage(page) {
    const totalPages = Math.ceil(usersData.length / usersPerPage); // Ensure correct page calculation
    if (page < 1 || page > totalPages) return; // Prevent invalid page numbers

    currentPage = page;
    displayUsers(currentPage); 
    renderPagination(); 
}


function editUser(username) {
    $.post('/api/admin/user-list', { page: currentPage }, function(data) {
        if (data.status == "success") {
            var user = data.users.find(u => u.username === username);
            if (user) {
                // Populate the form with user data
                document.getElementById('userOption').value = user.username;
                document.getElementById('concurrents').value = user.conns;
                document.getElementById('servers').value = user.servers;
                document.getElementById("testbalance").value = user.balance;
                document.getElementById('mbt').value = user.duration;
                
                // Set the expiration date field
                var expireDate = new Date(user.expiry * 1000); // Convert from Unix timestamp to JavaScript date
                document.getElementById('expire').value = expireDate.toISOString().split('T')[0]; // Set date in yyyy-mm-dd format

                // Populate the roles (if the roles are stored in user.ranks)
                var roleOptions = document.getElementById('roleOption');
                // Clear previous selections
                for (var option of roleOptions.options) {
                    option.selected = false;
                }
            } else {
                toastr['error']('User not found', 'Error', { "toastClass": "toast-dark" });
            }
        } else {
            handleErrors();
        }
    });
}
function UpdateUser() {
    // Retrieve values from form fields
    var username = document.getElementById('userOption').value;
    var concurrents = parseInt(document.getElementById('concurrents').value);
    var servers = parseInt(document.getElementById('servers').value);
    var balance = parseInt(document.getElementById("testbalance").value);
    var mbt = parseInt(document.getElementById('mbt').value); 
    // Get the value of the date input field
    var expireDateValue = document.getElementById("expire").value;
    console.log("username:", username);
console.log("concurrents:", concurrents);
console.log("servers:", servers);
console.log("balance:", balance);
console.log("mbt:", mbt);
    // Convert the date string to a Date object
    var expireDate = new Date(expireDateValue);

    // Get the Unix timestamp for the whole date
    var unixTimestamp = expireDate.getTime() / 1000;

    // Check if any field is empty or invalid
    if (!username || isNaN(concurrents) || isNaN(servers) || isNaN(balance) || isNaN(mbt)) {
        toastr['error']('Please fill in all fields with valid values!', 'Error', { "toastClass": "toast-dark" });
        return; // Exit the function if validation fails
    }
    
    // Retrieve selected roles from the multi-select dropdown
    var roleOptions = document.getElementById('roleOption').selectedOptions;
    var roles = Array.from(roleOptions).map(option => option.value);

    // Prepare the data to be sent in the AJAX request
    var userData = {
        username: username,
        concurrents: concurrents,
        servers: servers,
        duration: mbt,
        ranks: roles.map(role => ({ name: role, has: true })),
        expiry: unixTimestamp,
        balance: balance
    };
    console.log('User Data:', userData);
    // Make a POST request to the server
    fetch('/api/admin/update-user', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(userData)
    })
    .then(response => {
        if (response.ok) {
            toastr['success']('User updated successfully', 'Success', { "toastClass": "toast-dark" });
        } else {
            toastr['error']('Update Error', 'Error', { "toastClass": "toast-dark" });
        }
    })
    .catch(error => {
        toastr['error']('Failed to update user: ' + error.message, 'Error', { "toastClass": "toast-dark" });
    });
}

function deleteUser(username) {
    var deleteUser = {
        username: username
    };
    console.log('User Data:', deleteUser);
    // Make a POST request to the server
    fetch('/api/admin/delete-user', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(deleteUser)
    })
    .then(response => {
        if (response.ok) {
            window.location.reload();
            toastr['success']('User deleted successfully', 'Error', { "toastClass": "toast-dark" });
        } else {
            toastr['error']('Update Error', 'Error', { "toastClass": "toast-dark" });
        }
    })
    .catch(error => {
        console.error('Failed to delete user:', error.message);
        toastr['error']('Failed to delete user: ' + error.message, 'Error', { "toastClass": "toast-dark" });
    });
}
