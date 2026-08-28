import { state } from './state.js';
import { renderNavbar } from './components/navbar.js';
import { renderAuthPage } from './pages/auth.js';
import { renderUsersPage } from './pages/users.js';
import { renderClassroomsPage } from './pages/classrooms.js';
import { renderSchedulePage } from './pages/schedule.js';
import { renderJournalPage } from './pages/journal.js';
import { renderSubjectsPage } from './pages/subjects.js';

function router() {
    renderNavbar();

    const mainContent = document.getElementById('main-content');
    if (!mainContent) return;

    // Auth guard: if no user logged in, force auth page
    if (!state.token || !state.user) {
        renderAuthPage();
        return;
    }

    const role = state.user?.role?.toLowerCase() || 'user';

    switch (state.activePage) {
        case 'users':
            if (role === 'staff') {
                renderUsersPage();
            } else {
                state.setActivePage('journal');
            }
            break;

        case 'classrooms':
            renderClassroomsPage();
            break;

        case 'schedule':
            renderSchedulePage();
            break;

        case 'journal':
            renderJournalPage();
            break;

        case 'subjects':
            if (role === 'staff') {
                renderSubjectsPage();
            } else {
                state.setActivePage('journal');
            }
            break;

        case 'dashboard':
        default:
            renderJournalPage();
            break;
    }
}

// Initialize Application
document.addEventListener('DOMContentLoaded', () => {
    // Hide global loader if visible
    const loader = document.getElementById('global-loader');
    if (loader) loader.style.display = 'none';

    // Subscribe router to state changes
    state.subscribe(router);

    // Initial render
    router();
});
