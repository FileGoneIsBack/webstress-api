let loader = null;

(function() {
    loader = new DashboardLoader();
    setInterval(function () {
           $.post('/api/dashboard/data', function (data) {
               var json = $.parseJSON(data);
               loadTimeline(json.News)
           }).fail(handleErrors);
       }, 2000);
       loader.advance('Fetching dashboard data...', 10);
   })();

function loadTimeline(news) {
    const timelineList = $('#news-timeline');
    timelineList.empty(); // clear existing items
    news.forEach(function (item) {
        const timelineItem = renderTimeline(item);
        timelineList.append(timelineItem);
    });
    loader.advance('Loading stats...', 40);
}

function renderTimeline(item) {
    const formattedDate = formatDate(item.Date); 
    loader.advance('Ready, finalizing dashboard...', 80);
    loader.hide();
    return `
        <div class="news-item">
            <div class="news-title">${item.Title}</div>
            <div class="news-content">${item.Content}</div>
            <div class="news-date">${formattedDate}</div>
        </div>
    `;
    
}

//helper func
function formatDate(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric',
        hour: '2-digit', minute: '2-digit'
    });
}

