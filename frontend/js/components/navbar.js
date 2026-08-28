import { state } from '../state.js';
import { getRoleLabel } from '../utils.js';

export function renderNavbar() {
    const container = document.getElementById('navbar-container');
    if (!container) return;

    if (!state.user) {
        container.innerHTML = `
            <div class="navbar-inner">
                <a href="#" class="brand-logo">
                    <div class="brand-logo-icon">SO</div>
                    <span>School Online</span>
                </a>
            </div>
        `;
        return;
    }

    const role = state.user.role?.toLowerCase() || 'user';

    // Navigation items based on role
    let navItemsHTML = '';
    
    if (role === 'staff') {
        navItemsHTML = `
            <li class="nav-item"><button data-page="users" class="${state.activePage === 'users' ? 'active' : ''}">👥 Пользователи</button></li>
            <li class="nav-item"><button data-page="classrooms" class="${state.activePage === 'classrooms' ? 'active' : ''}">🏫 Классы</button></li>
            <li class="nav-item"><button data-page="schedule" class="${state.activePage === 'schedule' ? 'active' : ''}">📅 Расписание</button></li>
            <li class="nav-item"><button data-page="journal" class="${state.activePage === 'journal' ? 'active' : ''}">📖 Журнал</button></li>
            <li class="nav-item"><button data-page="subjects" class="${state.activePage === 'subjects' ? 'active' : ''}">📚 Предметы</button></li>
        `;
    } else if (role === 'teacher') {
        navItemsHTML = `
            <li class="nav-item"><button data-page="schedule" class="${state.activePage === 'schedule' ? 'active' : ''}">📅 Мое расписание</button></li>
            <li class="nav-item"><button data-page="journal" class="${state.activePage === 'journal' ? 'active' : ''}">📖 Выставление оценок</button></li>
            <li class="nav-item"><button data-page="classrooms" class="${state.activePage === 'classrooms' ? 'active' : ''}">🏫 Классы</button></li>
        `;
    } else if (role === 'student' || role === 'parent') {
        navItemsHTML = `
            <li class="nav-item"><button data-page="journal" class="${state.activePage === 'journal' ? 'active' : ''}">📖 Дневник и Оценки</button></li>
            <li class="nav-item"><button data-page="schedule" class="${state.activePage === 'schedule' ? 'active' : ''}">📅 Расписание уроков</button></li>
        `;
    } else {
        navItemsHTML = `
            <li class="nav-item"><button data-page="dashboard" class="${state.activePage === 'dashboard' ? 'active' : ''}">🏠 Главная</button></li>
        `;
    }

    container.innerHTML = `
        <div class="navbar-inner">
            <a href="#" class="brand-logo">
                <div class="brand-logo-icon">SO</div>
                <span>School Online</span>
            </a>

            <ul class="nav-links">
                ${navItemsHTML}
            </ul>

            <div class="user-profile-badge">
                <div style="display:flex; flex-direction:column; align-items:flex-end;">
                    <span style="font-size:0.88rem; font-weight:600;">${state.user.email || 'Пользователь'}</span>
                    <span class="role-pill ${role}">${getRoleLabel(role)}</span>
                </div>
                <button id="logout-btn" class="btn btn-secondary btn-sm" title="Выйтииз системы">Выйти</button>
            </div>
        </div>
    `;

    // Attach click listeners for nav buttons
    container.querySelectorAll('.nav-item button').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const page = e.currentTarget.getAttribute('data-page');
            if (page) state.setActivePage(page);
        });
    });

    const logoutBtn = container.querySelector('#logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', () => {
            state.logout();
        });
    }
}
