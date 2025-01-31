function handleErrors() {
    toastr['error']('Failed to call the API', 'Error', { "toastClass": "toast-dark" });
}
let currentPage = 1; 
const usersPerPage = 9; 

(function () {

    function loadUsers(data) {
        if (data.status == "success") {
            usersData = data.users; 
            renderPagination(); 
            displayUsers(currentPage); 
        }
    }

    $.post('/api/admin/user-list', function (data) {
        loadUsers(data);
    });

})();

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
function renderUser(user) {
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
                            <a class="dropdown-item" href="javascript:void(0);"><i class="ti ti-pencil me-1"></i>Edit</a>
                            <a class="dropdown-item" href="javascript:void(0);" onclick="deleteUser2('${user.username}')"><i class="ti ti-trash me-1"></i>Delete</a>
                        </div>
                    </div>
                </td>
            </tr>`;
}


function displayUsers(page) {
    const UsersList = $('#user-list');
    UsersList.empty(); 

    const startIndex = (page - 1) * usersPerPage;
    const endIndex = startIndex + usersPerPage;
    const usersToDisplay = usersData.slice(startIndex, endIndex);

    usersToDisplay.forEach(function (user) {
        const UserItem = renderUser(user);
        UsersList.append(UserItem);
    });
}

function changePage(page) {
    const totalPages = Math.ceil(usersData.length / usersPerPage);
    if (page < 1 || page > totalPages) return; 

    currentPage = page;
    displayUsers(currentPage); 
    renderPagination(); 
}

function deleteUser2(username) {
    var deleteUser = {
        username: username
    };
    console.log('User Data:', deleteUser);
    fetch('/api/admin/delete-user', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(deleteUser)
    })
        .then(response => {
            if (response.ok) {
                console.log('User deleted successfully');
                location.reload();
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
