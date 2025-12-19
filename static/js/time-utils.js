/**
 * Formats a date string into a localized string with a smart timezone acronym.
 * @param {string} isoDate - ISO 8601 date string (e.g. "2025-12-19T14:00:00Z")
 * @param {string} locale - Locale string (default "en-US")
 * @returns {string|null} Formatted string "YYYY-MM-DD HH:MM:SS ACRO" or null if invalid
 */
export function formatLocalizedDate(isoDate, locale = 'en-US') {
    if (!isoDate) return null;
    const date = new Date(isoDate);
    if (isNaN(date.getTime())) return null;

    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');

    const timeZoneName = getTimeZoneName(date, locale);

    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds} ${timeZoneName}`;
}

/**
 * Extracts a short timezone name (acronym) from a date.
 * Falls back to abbreviated long name if standard API returns an offset.
 * @param {Date} date
 * @param {string} locale
 * @returns {string} Timezone acronym (e.g. "JST", "EST") or short name
 */
export function getTimeZoneName(date, locale = 'en-US') {
    // Try to get the short name first (e.g. JST, EST, GMT+9)
    const short = new Intl.DateTimeFormat(locale, { timeZoneName: 'short' })
        .formatToParts(date)
        .find(part => part.type === 'timeZoneName')?.value || '';

    // If it's not an offset style (doesn't contain GMT, UTC, or +s/-), return it
    if (!short.match(/GMT|UTC|[+-]\d/)) {
        return short;
    }

    // Fallback: Try to abbreviate the long name from date.toString()
    // e.g. "Sat Dec 20 ... (Japan Standard Time)" -> "Japan Standard Time"
    const longNameMatch = /\(([^)]+)\)/.exec(date.toString());
    if (longNameMatch) {
        const longName = longNameMatch[1];
        // Create acronym: "Japan Standard Time" -> "JST"
        const acronym = longName.match(/\b[A-Z]/g)?.join('');
        if (acronym && acronym.length >= 2) {
            return acronym;
        }
    }

    return short;
}
