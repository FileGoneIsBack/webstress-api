document.addEventListener('DOMContentLoaded', function () {
    document.querySelector('.send-button')?.addEventListener('click', updateTicket);
    document.querySelector('.status-select')?.addEventListener('change', function () {
        updateStatus(this.value);
    });
    document.querySelector('.chat-input textarea')?.addEventListener('input', function () {
        this.style.height = 'auto';
        this.style.height = (this.scrollHeight) + 'px';
    });
});
function getAllCookies() {
    return Object.fromEntries(
        document.cookie.split('; ').map(cookie => {
            const [name, value] = cookie.split('=');
            return [name, decodeURIComponent(value)];
        })
    );
}

function renderMessages(messages) {
    const messagesContainer = document.querySelector('.chat-messages');
    getCurrentUserID().then(currentUserID => {
        messagesContainer.innerHTML = ''; // Clear existing

        messages.forEach(msg => {
            const messageEl = document.createElement('div');
            messageEl.classList.add('message');
            messageEl.classList.add(msg.user_id === currentUserID ? 'user' : 'support');

            const avatar = document.createElement('div');
            avatar.className = 'message-avatar';
            avatar.textContent = msg.user_id === currentUserID ? 'You' : 'SA';

            const content = document.createElement('div');
            content.className = 'message-content';
            content.innerHTML = `
                <p>${msg.message}</p>
                <span class="message-time">${new Date(msg.created_at).toLocaleTimeString([], {hour: '2-digit', minute: '2-digit'})}</span>
            `;

            if (msg.user_id === currentUserID) {
                messageEl.appendChild(content);
                messageEl.appendChild(avatar);
            } else {
                messageEl.appendChild(avatar);
                messageEl.appendChild(content);
            }

            messagesContainer.appendChild(messageEl);
        });

        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    });
}

function getCurrentUserID() {
    return $.get('/api/dashboard/user_id').then(function(data) {
        return data.user_id;
    }).fail(function(xhr, status, error) {
        console.error('Failed to fetch current user ID:', error);
    });
}

function updateTicket() {
    const cookies = getAllCookies();
    const ticketId = cookies.ticket_id;
    const textarea = document.querySelector('.chat-input textarea');
    const message = textarea.value.trim();
    if (!message) return;
    $.ajax({
        url: '/api/tickets/update',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
    ticketid: parseInt(ticketId), 
    message: message
}),
        success: function (response) {
            const data = $.parseJSON(response);
            if (data.status === 'success') {
                textarea.value = "";
                fetchAndRenderNewMessages();
            } else {
                console.log(data.message);
            }
        },
    });
}

var lastMessageId = null;

//on load
(function() {
    fetchAndRenderNewMessages();
})();


function updateMessagesPeriodically() {
    setInterval(function() {
        fetchAndRenderNewMessages();
    }, 1000);
}

function fetchAndRenderNewMessages() {
    const cookies = getAllCookies();
    const ticketId = cookies.ticket_id;
    $.post('/api/tickets/messages', { ticketid: ticketId, lastMessageId })
        .done(data => {
            console.log(data)
            const jsonData = $.parseJSON(data);
            if (jsonData.messages.length > 0) {
                lastMessageId = jsonData.messages[jsonData.messages.length - 1].id;
            }
            renderMessages(jsonData.messages);
        })
}


updateMessagesPeriodically();
fetchUserInfo();

function updateStatus(status) {
    var ticketIdStr = getCookie('ticket_id');
    var ticketId = parseInt(ticketIdStr, 10);


    var jsonData = {
        ticketid: ticketId,
        status: status
    };

    console.log(jsonData)

    $.ajax({
        url: '/api/admin/updateTicket',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(jsonData),
        success: function(response) {
            var data = $.parseJSON(response);
            if (data.status === "success") {
                console.log('Ticket status updated successfully');
            } else {
                console.log('Failed to update ticket status');
            }
        },
        error: function(xhr, status, error) {
            console.log(error)
        }
    });
}
function fetchUserInfo() {
    $.get("/api/admin/GetUser")
        .done(function (user) {
            $(".user-details h2").text(user.Username);
            $(".user-email").text("TG ID: " + user.Tele);
            $(".stat-value").eq(0).text(user.Balance);
            // Format expiry timestamp
            if (user.Expiry) {
                const expiryDate = new Date(user.Expiry * 1000); // convert UNIX to ms
                const formatted = expiryDate.toLocaleDateString("en-US", {
                    year: "numeric",
                    month: "long",
                    day: "numeric"
                });
                $(".stat-value").eq(4).text(formatted); // "Member Since"
            }

            $(".stat-value").eq(1).text(user.Membership);

            // Optionally, update ranks
            if (Array.isArray(user.Ranks)) {
                const roleList = user.Ranks.map(r => r.name).join(", ");
                $(".stat-value").eq(2).text(roleList);
            }

        })
        .fail(function (xhr) {
            console.error("Failed to load user:", xhr.responseText);
            showToast("User not found", "error");
        });
}
    