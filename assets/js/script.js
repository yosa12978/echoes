function toDateString_(dateISO8601, timezone = 'GMT') {
    const d = new Date(dateISO8601);
    return d.toLocaleDateString('en-US', {
        timeZone: 'UTC',
        hour: 'numeric',
        minute: 'numeric',
        year: "numeric",
        month: "short",
        day: "numeric",
        hour12: true
    })
}

function renderMarkdown(md) {
    return DOMPurify.sanitize(marked.parse(md), { ADD_TAGS: ["iframe"], ADD_ATTR: ['allow', 'allowfullscreen', 'frameborder', 'scrolling'] });
}
