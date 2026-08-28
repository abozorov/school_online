import { api } from '../api.js';
import { state } from '../state.js';
import { showToast } from '../components/toast.js';

export function renderAuthPage() {
    const mainContent = document.getElementById('main-content');
    if (!mainContent) return;

    mainContent.innerHTML = `
        <div style="max-width: 440px; margin: 3rem auto 0 auto;">
            <div class="glass-card" style="padding: 2.5rem 2rem;">
                <div style="text-align: center; margin-bottom: 2rem;">
                    <div class="brand-logo-icon" style="width: 56px; height: 56px; margin: 0 auto 1rem auto; font-size: 1.6rem;">SO</div>
                    <h2 style="font-family: var(--font-heading); font-size: 1.75rem; font-weight: 700;">Вход в систему</h2>
                    <p style="color: var(--text-muted); font-size: 0.9rem; margin-top: 0.35rem;">Электронная школа School Online</p>
                </div>

                <form id="login-form">
                    <div class="form-group">
                        <label class="form-label" for="login-email">Email адрес</label>
                        <input type="email" id="login-email" class="form-input" placeholder="teacher@school.edu" required value="staff@school.com">
                    </div>

                    <div class="form-group">
                        <label class="form-label" for="login-password">Пароль</label>
                        <input type="password" id="login-password" class="form-input" placeholder="••••••••" required value="password123">
                    </div>

                    <button type="submit" class="btn btn-primary" style="width: 100%; margin-top: 1rem; padding: 0.85rem;">
                        Войти в кабинет
                    </button>
                </form>

                <div style="margin-top: 2rem; padding-top: 1.5rem; border-top: 1px solid var(--border-color);">
                    <p style="font-size: 0.82rem; color: var(--text-muted); text-align: center; margin-bottom: 0.75rem;">
                        Тестовый быстрый вход:
                    </p>
                    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.5rem;">
                        <button class="btn btn-secondary btn-sm demo-role-btn" data-email="staff@school.com">Администратор</button>
                        <button class="btn btn-secondary btn-sm demo-role-btn" data-email="teacher@school.com">Преподаватель</button>
                        <button class="btn btn-secondary btn-sm demo-role-btn" data-email="student@school.com">Ученик</button>
                        <button class="btn btn-secondary btn-sm demo-role-btn" data-email="parent@school.com">Родитель</button>
                    </div>
                </div>
            </div>
        </div>
    `;

    // Demo role buttons
    mainContent.querySelectorAll('.demo-role-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const email = e.currentTarget.getAttribute('data-email');
            document.getElementById('login-email').value = email;
            document.getElementById('login-password').value = 'password123';
        });
    });

    // Form submit
    const form = document.getElementById('login-form');
    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const email = document.getElementById('login-email').value;
        const password = document.getElementById('login-password').value;

        const submitBtn = form.querySelector('button[type="submit"]');
        submitBtn.disabled = true;
        submitBtn.textContent = 'Авторизация...';

        try {
            const res = await api.login(email, password);
            showToast('Успешный вход в систему!', 'success');
            state.setAuth(res);
            
            // Set initial active page depending on role
            const role = state.user?.role?.toLowerCase() || 'user';
            if (role === 'staff') state.setActivePage('users');
            else if (role === 'teacher') state.setActivePage('journal');
            else state.setActivePage('journal');

        } catch (error) {
            // If server fails (e.g. backend offline), allow dev fallback auth for testing UI
            console.warn('Backend login failed, using auth fallback if needed:', error);
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = 'Войти в кабинет';
        }
    });
}
