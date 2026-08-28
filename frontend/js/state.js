import { CONFIG } from './config.js';
import { parseJwt } from './utils.js';

class AppState {
    constructor() {
        this.token = localStorage.getItem(CONFIG.STORAGE_KEY_TOKEN) || null;
        this.refreshToken = localStorage.getItem(CONFIG.STORAGE_KEY_REFRESH) || null;
        this.user = this.initUser();
        this.activePage = 'dashboard';
        this.listeners = new Set();
    }

    initUser() {
        const cachedUser = localStorage.getItem(CONFIG.STORAGE_KEY_USER);
        if (cachedUser) {
            try { return JSON.parse(cachedUser); } catch(e) {}
        }
        if (this.token) {
            const claims = parseJwt(this.token);
            if (claims) {
                return {
                    id: claims.user_id || claims.sub || claims.id,
                    email: claims.email || '',
                    role: claims.role || 'user'
                };
            }
        }
        return null;
    }

    setAuth(tokens, userInfo = null) {
        this.token = tokens.jwt_token || tokens.jwt || null;
        this.refreshToken = tokens.refresh_token || tokens.refresh || null;
        
        if (this.token) {
            localStorage.setItem(CONFIG.STORAGE_KEY_TOKEN, this.token);
        }
        if (this.refreshToken) {
            localStorage.setItem(CONFIG.STORAGE_KEY_REFRESH, this.refreshToken);
        }

        if (userInfo) {
            this.user = userInfo;
            localStorage.setItem(CONFIG.STORAGE_KEY_USER, JSON.stringify(userInfo));
        } else if (this.token) {
            const claims = parseJwt(this.token);
            this.user = {
                id: claims?.user_id || claims?.sub || 1,
                email: claims?.email || '',
                role: claims?.role || 'user'
            };
            localStorage.setItem(CONFIG.STORAGE_KEY_USER, JSON.stringify(this.user));
        }

        this.notify();
    }

    logout() {
        this.token = null;
        this.refreshToken = null;
        this.user = null;
        localStorage.removeItem(CONFIG.STORAGE_KEY_TOKEN);
        localStorage.removeItem(CONFIG.STORAGE_KEY_REFRESH);
        localStorage.removeItem(CONFIG.STORAGE_KEY_USER);
        this.activePage = 'auth';
        this.notify();
    }

    setActivePage(page) {
        this.activePage = page;
        this.notify();
    }

    subscribe(listener) {
        this.listeners.add(listener);
        return () => this.listeners.delete(listener);
    }

    notify() {
        for (const listener of this.listeners) {
            listener(this);
        }
    }
}

export const state = new AppState();
