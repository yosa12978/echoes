function toDateString_(dateISO8601, timezone = 'GMT') {
    const d = new Date(dateISO8601);
    return d.getUTCFullYear().toString() +
        "/" + (d.getUTCMonth() + 1).toString() +
        "/" + d.getUTCDate().toString() +
        " - " + d.toLocaleString('en-US', { timeZone: 'UTC', hour: 'numeric', minute: 'numeric', hour12: true });
}

function renderMarkdown(md) {
    return DOMPurify.sanitize(marked.parse(md), { ADD_TAGS: ["iframe"], ADD_ATTR: ['allow', 'allowfullscreen', 'frameborder', 'scrolling'] });
}
