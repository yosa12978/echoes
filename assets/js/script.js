function toDateString_(dateISO8601, timezone = 'GMT') {
    const d = new Date(dateISO8601);
    let diffInTime = (new Date()).getTime() - d.getTime();
    let diffInDays = Math.round(diffInTime / (1000 * 3600 * 24));
    if (diffInDays > 7) {
        return "on " + d.toLocaleDateString('en-US', {
            timeZone: 'UTC',
            year: "numeric",
            month: "long",
            day: "numeric",
            hour12: true
        })
    } 
    return moment(dateISO8601).utc().fromNow()
}

function renderMarkdown(md) {
    return DOMPurify.sanitize(marked.parse(md), { ADD_TAGS: ["iframe"], ADD_ATTR: ['allow', 'allowfullscreen', 'frameborder', 'scrolling'] });
}
