// Utility Functions

/**
 * Format Date to DD.MM.YYYY
 */
export function formatDate(dateObj = new Date()) {
    const d = new Date(dateObj);
    if (isNaN(d.getTime())) return '';
    const day = String(d.getDate()).padStart(2, '0');
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const year = d.getFullYear();
    return `${day}.${month}.${year}`;
}

/**
 * Generate default date_range string (e.g. 01.01.2024 - today)
 */
export function getDefaultDateRange() {
    const now = new Date();
    const startOfYear = new Date(now.getFullYear(), 0, 1);
    return `${formatDate(startOfYear)}-${formatDate(now)}`;
}

/**
 * Decode JWT token payload without external libraries
 */
export function parseJwt(token) {
    try {
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(window.atob(base64).split('').map(c => {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));

        return JSON.parse(jsonPayload);
    } catch (e) {
        return null;
    }
}

/**
 * Translate backend role code to human-readable Russian label
 */
export function getRoleLabel(role) {
    switch (role?.toLowerCase()) {
        case 'staff': return 'Администратор';
        case 'teacher': return 'Преподаватель';
        case 'student': return 'Ученик';
        case 'parent': return 'Родитель';
        case 'user': return 'Пользователь';
        default: return role || 'Гость';
    }
}

/**
 * Day of week names (1-7: Mon-Sun)
 */
export const DAYS_OF_WEEK = [
    { id: 1, name: 'Понедельник' },
    { id: 2, name: 'Вторник' },
    { id: 3, name: 'Среда' },
    { id: 4, name: 'Четверг' },
    { id: 5, name: 'Пятница' },
    { id: 6, name: 'Суббота' },
    { id: 7, name: 'Воскресенье' }
];

export function getDayName(dayNumber) {
    const found = DAYS_OF_WEEK.find(d => d.id === Number(dayNumber));
    return found ? found.name : `День ${dayNumber}`;
}

/**
 * Format experience in days to years and months (omitting leftover days)
 */
export function formatExperience(days) {
    const totalDays = Number(days);
    if (!totalDays || totalDays <= 0) return '0 месяцев';

    const years = Math.floor(totalDays / 365);
    const months = Math.floor((totalDays % 365) / 30);

    const parts = [];
    if (years > 0) {
        let yearWord = 'лет';
        const mod10 = years % 10;
        const mod100 = years % 100;
        if (mod10 === 1 && mod100 !== 11) {
            yearWord = 'год';
        } else if (mod10 >= 2 && mod10 <= 4 && (mod100 < 10 || mod100 >= 20)) {
            yearWord = 'года';
        }
        parts.push(`${years} ${yearWord}`);
    }

    if (months > 0) {
        let monthWord = 'месяцев';
        const mod10 = months % 10;
        const mod100 = months % 100;
        if (mod10 === 1 && mod100 !== 11) {
            monthWord = 'месяц';
        } else if (mod10 >= 2 && mod10 <= 4 && (mod100 < 10 || mod100 >= 20)) {
            monthWord = 'месяца';
        }
        parts.push(`${months} ${monthWord}`);
    }

    return parts.length > 0 ? parts.join(' ') : '0 месяцев';
}

