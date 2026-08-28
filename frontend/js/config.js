// Application Configuration
export const CONFIG = {
    // API Gateway URL (defaults to current origin)
    API_BASE_URL: window.location.origin || 'http://localhost:8080',
    STORAGE_KEY_TOKEN: 'school_online_jwt',
    STORAGE_KEY_REFRESH: 'school_online_refresh',
    STORAGE_KEY_USER: 'school_online_user_info'
};
