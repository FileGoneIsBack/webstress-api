            // Function to update the blacklist (POST request)
            function UpdateBlacklist() {
                var host = $('#host').val();  // Get the value from the input field
        
                // Check if the host field is empty
                if (host === '') {
                    toastr['error']('Please fill in all fields!', 'Error', { "toastClass": "toast-dark" });
                    return;  // Exit if validation fails
                }
        
                // Prepare the data to be sent in the AJAX request
                var BlacklistData = { host: host };
        
                // Make a POST request to add the host to the blacklist
                $.ajax({
                    url: '/api/admin/blacklist',  // Your backend endpoint
                    type: 'POST',
                    contentType: 'application/json',
                    data: JSON.stringify(BlacklistData),
                    success: function(response) {
                        // If successful, show success message
                        toastr['success']('Host added to the blacklist successfully!', 'Success', { "toastClass": "toast-dark" });
        
                        // Reload the blacklist list to show the newly added host
                        loadBlacklist();
                    },
                    error: function(xhr, status, error) {
                        // If error occurs, show error message
                        toastr['error']('Failed to add host to the blacklist: ' + error, 'Error', { "toastClass": "toast-dark" });
                    }
                });
            }

            function UpdateAuth() {
                var token = $('#token').val();  
                var exp = $('#date').val();
                if (token === '' || exp === '') {
                    toastr['error']('Please fill in all fields!', 'Error', { "toastClass": "toast-dark" });
                    return;
                }
                var dateObj = new Date(exp);
                var timestamp = dateObj.getTime();
                var BlacklistData = { token: token, exp: timestamp };
                $.ajax({
                    url: '/api/admin/AddToken', 
                    type: 'POST',
                    contentType: 'application/json',
                    data: JSON.stringify(BlacklistData),
                    success: function (response) {
                        toastr['success']('token added to the invite manager successfully!', 'Success', { "toastClass": "toast-dark" });

                        loadBlacklist();
                    },
                    error: function (xhr, status, error) {
                        toastr['error']('Failed to add token to the invite manager: ' + error, 'Error', { "toastClass": "toast-dark" });
                    }
                });
            }
                    
            // Function to load the blacklist (GET request)
            function loadBlacklist() {
                $.ajax({
                    url: '/api/admin/blacklist',
                    type: 'GET',
                    contentType: 'application/json',
                    success: function (response) {
                        var apiList = $('#apiList'); 
                        apiList.empty(); 

                        response.forEach(function (host, index) {
                            var rowID = 'row-' + (index + 1); 
                            var deleteButton = '<button onclick="del(\'' + rowID + '\')" id="action-btn-' + rowID + '" class="btn btn-danger btn-md" type="button">Delete</button>';

                            apiList.append('<tr id="' + rowID + '">' +
                                '<td>' + (index + 1) + '</td>' +
                                '<td>' + host + '</td>' + 
                                '<td>' + deleteButton + '</td>' +
                                '</tr>');
                        });
                    },
                    error: function (xhr, status, error) {
                        toastr['error']('Failed to load blacklist: ' + error, 'Error', { "toastClass": "toast-dark" });
                    }
                });
            }

            // Function to load the invites (GET request)
            function loadInvites() {
                $.ajax({
                    url: '/api/admin/invites',
                    type: 'GET',
                    contentType: 'application/json',
                    success: function (response) {
                        var apiList = $('#blacklist');  // Ensure this targets the correct table body
                        apiList.empty();  // Clear the existing table contents

                        // Assuming response is an array of invite objects
                        response.forEach(function (invite, index) {
                            var rowID = 'row-' + (index + 1);  // Generate a unique ID for each row
                            var deleteButton = '<button onclick="delInvite(\'' + invite.Token + '\')" id="action-btn-' + rowID + '" class="btn btn-danger btn-md" type="button">Delete</button>';
                            var userAvatars = '';

                            // Check if 'Uses' is an array before calling forEach
                            if (Array.isArray(invite.Uses)) {
                                invite.Uses.forEach(function (username) {
                                    userAvatars += `
                                        <li data-bs-toggle="tooltip" data-bs-placement="top" class="avatar avatar-xs pull-up" title="${username}">
                                            <img src="https://demos.pixinvent.com/vuexy-html-admin-template/documentation/assets/img/avatars/7.png" alt="Avatar" class="rounded-circle avatar-img-small">
                                        </li>
                                    `;
                                });
                            } else {
                                console.error("Invite 'Uses' is not an array:", invite.Uses);
                            }

                            // Append a new row to the table with the delete button
                            apiList.append('<tr id="' + rowID + '">' +
                                '<td>' + (index + 1) + '</td>' +
                                '<td>' + invite.Token + '</td>' +
                                '<td><ul class="list-unstyled users-list m-0 avatar-group d-flex align-items-center">' + userAvatars + '</ul></td>' +  // Display the list of avatars for users
                                '<td>' + deleteButton + '</td>' +
                                '</tr>');
                        });
                    },
                    error: function (xhr, status, error) {
                        toastr['error']('Failed to load invites: ' + error, 'Error', { "toastClass": "toast-dark" });
                    }
                });
            }


            // Function to handle the DELETE button click
            function del(rowID) {
                var host = $('#' + rowID).find('td').eq(1).text();

                $.ajax({
                    url: '/api/admin/removeBlacklist', 
                    type: 'POST',
                    contentType: 'application/json',
                    data: JSON.stringify({ host: host }), 
                    success: function(response) {
                        $('#' + rowID).remove();

                        toastr['success']('Host successfully removed from blacklist', 'Success', { "toastClass": "toast-dark" });
                    },
                    error: function(xhr, status, error) {
                        console.error('Error removing host:', error);
                        toastr['error']('Failed to remove host: ' + error, 'Error', { "toastClass": "toast-dark" });
                    }
                });
            }
            // Function to handle the DELETE button click
            function delInvite(rowID) {
                var token = $('#' + rowID).find('td').eq(1).text(); 
                $.ajax({
                    url: '/api/admin/deleteInvite',  
                    type: 'POST',
                    contentType: 'application/json',
                    data: JSON.stringify({ token: token }),  
                    success: function (response) {
                        $('#' + rowID).remove();
                        toastr['success']('Invite successfully removed', 'Success', { "toastClass": "toast-dark" });
                    },
                    error: function (xhr, status, error) {
                        console.error('Error removing invite:', error);
                        toastr['error']('Failed to remove invite: ' + error, 'Error', { "toastClass": "toast-dark" });
                    }
                });
            }