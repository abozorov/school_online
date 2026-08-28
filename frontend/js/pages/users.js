import { api } from '../api.js';
import { getRoleLabel } from '../utils.js';
import { showToast } from '../components/toast.js';
import { openModal, closeModal } from '../components/modal.js';

let usersList = [];
let classroomsList = [];
let classroomsMap = {};
let subjectsList = [];
let subjectsMap = {};

function formatExperience(days) {
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


export async function renderUsersPage() {
    const mainContent = document.getElementById('main-content');
    if (!mainContent) return;

    mainContent.innerHTML = `
        <div class="page-header">
            <div>
                <h1 class="page-title">Управление пользователями</h1>
                <p class="page-subtitle">Просмотр, создание, привязка ролевых профилей (классы, дети, предметы, стаж)</p>
            </div>
            <button id="add-user-btn" class="btn btn-primary">
                <span>➕</span> Добавить пользователя
            </button>
        </div>

        <div class="glass-card mb-4" style="margin-bottom: 1.5rem; padding: 1rem 1.25rem;">
            <div style="display: flex; gap: 1rem; flex-wrap: wrap; align-items: center;">
                <div style="flex: 1; min-width: 240px;">
                    <input type="text" id="user-search-input" class="form-input" placeholder="Поиск по имени, email, логину..." style="padding: 0.55rem 0.85rem;">
                </div>
                <div>
                    <select id="user-role-filter" class="form-select" style="padding: 0.55rem 0.85rem;">
                        <option value="">Все роли</option>
                        <option value="staff">Администратор</option>
                        <option value="teacher">Преподаватель</option>
                        <option value="student">Ученик</option>
                        <option value="parent">Родитель</option>
                        <option value="user">Обычный пользователь</option>
                    </select>
                </div>
            </div>
        </div>

        <div id="users-table-container" class="table-container">
            <div class="loading-spinner-wrapper" style="min-height: 200px;">
                <div class="spinner"></div>
                <p>Загрузка списка пользователей...</p>
            </div>
        </div>
    `;

    document.getElementById('add-user-btn').addEventListener('click', () => {
        openUserModal();
    });

    document.getElementById('user-search-input').addEventListener('input', filterAndRenderUsers);
    document.getElementById('user-role-filter').addEventListener('change', filterAndRenderUsers);

    await loadData();
}

async function loadData() {
    try {
        const [usersData, classroomsData, subjectsData] = await Promise.all([
            api.getUsersList(),
            api.getClassroomsList().catch(() => []),
            api.getSubjectsList().catch(() => [])
        ]);
        usersList = usersData || [];
        classroomsList = classroomsData || [];
        subjectsList = subjectsData || [];

        classroomsMap = {};
        classroomsList.forEach(c => {
            classroomsMap[c.id] = `${c.grade_number}${c.letter || ''}`;
        });

        subjectsMap = {};
        subjectsList.forEach(s => {
            subjectsMap[s.id] = s.name;
        });

        filterAndRenderUsers();
    } catch (error) {
        console.error('Failed to load users:', error);
        document.getElementById('users-table-container').innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">⚠️</div>
                <p>Не удалось загрузить пользователей (${error.message})</p>
                <button class="btn btn-secondary btn-sm" style="margin-top: 1rem;" onclick="location.reload()">Повторить</button>
            </div>
        `;
    }
}

function getChildrenForParent(parentUser) {
    const studentIds = parentUser.parent_role?.students_id || parentUser.students_id || [];
    if (!Array.isArray(studentIds) || studentIds.length === 0) return [];
    
    return usersList.filter(u => studentIds.includes(Number(u.id)));
}

function filterAndRenderUsers() {
    const search = document.getElementById('user-search-input')?.value.toLowerCase() || '';
    const roleFilter = document.getElementById('user-role-filter')?.value.toLowerCase() || '';

    const filtered = usersList.filter(u => {
        const matchesSearch = (
            (u.name && u.name.toLowerCase().includes(search)) ||
            (u.email && u.email.toLowerCase().includes(search)) ||
            (u.username && u.username.toLowerCase().includes(search))
        );
        const matchesRole = !roleFilter || (u.role && u.role.toLowerCase() === roleFilter);
        return matchesSearch && matchesRole;
    });

    const container = document.getElementById('users-table-container');
    if (!container) return;

    if (filtered.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">👤</div>
                <p>Пользователи не найдены</p>
            </div>
        `;
        return;
    }

    let rowsHTML = filtered.map(u => {
        const role = u.role?.toLowerCase() || 'user';
        let extraInfoHTML = '';

        if (role === 'student') {
            const classId = u.student_role?.classroom_id || u.classroom_id;
            const className = classroomsMap[classId] ? `Класс ${classroomsMap[classId]}` : (classId ? `Класс #${classId}` : 'Без класса');
            extraInfoHTML = `<div style="font-size: 0.78rem; color: var(--primary); font-weight: 500;">🏫 ${className}</div>`;
        } else if (role === 'parent') {
            const children = getChildrenForParent(u);
            if (children.length > 0) {
                const namesStr = children.map(c => c.name || `@${c.username}`).join(', ');
                extraInfoHTML = `<div style="font-size: 0.78rem; color: var(--accent); font-weight: 500;">👶 Дети (${children.length}): ${namesStr}</div>`;
            } else {
                extraInfoHTML = `<div style="font-size: 0.78rem; color: var(--text-subtle);">👶 Дети не привязаны</div>`;
            }
        } else if (role === 'staff') {
            const pos = u.staff_role?.position ? `Должность: ${u.staff_role.position}` : '';
            const exp = u.staff_role?.experience ? `Стаж: ${formatExperience(u.staff_role.experience)}` : '';
            const details = [pos, exp].filter(Boolean).join(' • ');
            if (details) {
                extraInfoHTML = `<div style="font-size: 0.78rem; color: var(--text-muted);">💼 ${details}</div>`;
            }
        } else if (role === 'teacher') {
            const subjIds = u.teacher_role?.subjects_id || [];
            const subNames = subjIds.map(id => subjectsMap[id] || `#${id}`).join(', ');
            const subjStr = subNames ? `Предметы: ${subNames}` : '';
            const expStr = u.teacher_role?.experience ? `Стаж: ${formatExperience(u.teacher_role.experience)}` : '';
            const details = [subjStr, expStr].filter(Boolean).join(' • ');
            if (details) {
                extraInfoHTML = `<div style="font-size: 0.78rem; color: var(--text-muted);">📚 ${details}</div>`;
            }
        }

        return `
            <tr>
                <td><strong>#${u.id}</strong></td>
                <td>
                    <div style="font-weight: 600;">${u.name || '—'}</div>
                    <div style="font-size: 0.8rem; color: var(--text-subtle);">@${u.username || 'user'}</div>
                    ${extraInfoHTML}
                </td>
                <td>${u.email}</td>
                <td>
                    <span class="role-pill ${role}">${getRoleLabel(u.role)}</span>
                </td>
                <td>${u.phone_number || '—'}</td>
                <td>${u.birthday || '—'}</td>
                <td style="text-align: right;">
                    ${role === 'parent' ? `
                    <button class="btn btn-secondary btn-sm view-parent-children-btn" data-id="${u.id}" title="Просмотр детей">👶</button>
                    ` : ''}
                    <button class="btn btn-secondary btn-sm edit-user-btn" data-id="${u.id}" title="Редактировать">✏️</button>
                    <button class="btn btn-danger btn-sm delete-user-btn" data-id="${u.id}" title="Удалить">🗑️</button>
                </td>
            </tr>
        `;
    }).join('');

    container.innerHTML = `
        <table class="data-table">
            <thead>
                <tr>
                    <th>ID</th>
                    <th>ФИО / Логин</th>
                    <th>Email</th>
                    <th>Роль</th>
                    <th>Телефон</th>
                    <th>Дата рождения</th>
                    <th style="text-align: right;">Действия</th>
                </tr>
            </thead>
            <tbody>
                ${rowsHTML}
            </tbody>
        </table>
    `;

    // Attach row events
    container.querySelectorAll('.view-parent-children-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const id = Number(e.currentTarget.getAttribute('data-id'));
            const parentUser = usersList.find(u => u.id === id);
            if (parentUser) openParentChildrenModal(parentUser);
        });
    });

    container.querySelectorAll('.edit-user-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const id = Number(e.currentTarget.getAttribute('data-id'));
            const user = usersList.find(u => u.id === id);
            if (user) openUserModal(user);
        });
    });

    container.querySelectorAll('.delete-user-btn').forEach(btn => {
        btn.addEventListener('click', async (e) => {
            const id = Number(e.currentTarget.getAttribute('data-id'));
            if (confirm(`Вы уверены, что хотите удалить пользователя #${id}?`)) {
                try {
                    await api.deleteUser(id);
                    showToast('Пользователь успешно удален', 'success');
                    await loadData();
                } catch (err) {
                    showToast(err.message, 'error');
                }
            }
        });
    });
}

function openParentChildrenModal(parentUser) {
    const children = getChildrenForParent(parentUser);

    let bodyHTML = '';
    if (children.length === 0) {
        bodyHTML = `
            <div class="empty-state" style="padding: 2rem 1rem;">
                <div class="empty-state-icon">👶</div>
                <p>У родителя <strong>${parentUser.name || `@${parentUser.username}`}</strong> нет привязанных детей.</p>
            </div>
        `;
    } else {
        const rowsHTML = children.map(c => {
            const classId = c.student_role?.classroom_id || c.classroom_id;
            const className = classroomsMap[classId] ? `Класс ${classroomsMap[classId]}` : (classId ? `Класс #${classId}` : '—');
            return `
                <tr>
                    <td><strong>#${c.id}</strong></td>
                    <td>
                        <div style="font-weight: 600;">${c.name || '—'}</div>
                        <div style="font-size: 0.8rem; color: var(--text-subtle);">@${c.username || 'student'}</div>
                    </td>
                    <td><span class="role-pill student">${className}</span></td>
                    <td>${c.email || '—'}</td>
                    <td>${c.phone_number || '—'}</td>
                </tr>
            `;
        }).join('');

        bodyHTML = `
            <div style="margin-bottom: 1rem; font-size: 0.95rem; color: var(--text-muted);">
                Дети родителя <strong>${parentUser.name || `@${parentUser.username}`}</strong> (всего: ${children.length}):
            </div>
            <div class="table-container">
                <table class="data-table">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>ФИО / Логин</th>
                            <th>Класс</th>
                            <th>Email</th>
                            <th>Телефон</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${rowsHTML}
                    </tbody>
                </table>
            </div>
        `;
    }

    const footerHTML = `
        <button type="button" class="btn btn-primary" onclick="document.getElementById('modal-close-btn').click()">Закрыть</button>
    `;

    openModal({
        title: `Дети родителя ${parentUser.name || `@${parentUser.username}`}`,
        bodyHTML,
        footerHTML
    });
}

function openUserModal(user = null) {
    const isEdit = !!user;
    const currentRole = user?.role?.toLowerCase() || 'user';

    const bodyHTML = `
        <form id="user-form">
            <div class="form-group">
                <label class="form-label">ФИО</label>
                <input type="text" id="user-name" class="form-input" required value="${user?.name || ''}" placeholder="Иванов Иван Иванович">
            </div>
            
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                <div class="form-group">
                    <label class="form-label">Username (Логин)</label>
                    <input type="text" id="user-username" class="form-input" required value="${user?.username || ''}" placeholder="ivanov_i">
                </div>
                <div class="form-group">
                    <label class="form-label">Email</label>
                    <input type="email" id="user-email" class="form-input" required value="${user?.email || ''}" placeholder="ivanov@school.com">
                </div>
            </div>

            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                <div class="form-group">
                    <label class="form-label">Роль</label>
                    <select id="user-role" class="form-select">
                        <option value="user" ${currentRole === 'user' ? 'selected' : ''}>Пользователь</option>
                        <option value="staff" ${currentRole === 'staff' ? 'selected' : ''}>Администратор (Staff)</option>
                        <option value="teacher" ${currentRole === 'teacher' ? 'selected' : ''}>Учитель (Teacher)</option>
                        <option value="student" ${currentRole === 'student' ? 'selected' : ''}>Ученик (Student)</option>
                        <option value="parent" ${currentRole === 'parent' ? 'selected' : ''}>Родитель (Parent)</option>
                    </select>
                </div>
                <div class="form-group">
                    <label class="form-label">Телефон</label>
                    <input type="text" id="user-phone" class="form-input" value="${user?.phone_number || ''}" placeholder="79991234567">
                </div>
            </div>

            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                <div class="form-group">
                    <label class="form-label">Дата рождения (ДД.ММ.ГГГГ)</label>
                    <input type="text" id="user-birthday" class="form-input" value="${user?.birthday || ''}" placeholder="15.05.2008">
                </div>
                ${!isEdit ? `
                <div class="form-group">
                    <label class="form-label">Пароль (мин. 8 символов)</label>
                    <input type="password" id="user-password" class="form-input" required value="Password123" placeholder="••••••••">
                </div>
                ` : ''}
            </div>

            <!-- Role-Specific Fields Container -->
            
            <!-- 1. Student Role Fields -->
            <div id="student-role-fields" class="role-fields-panel" style="display: ${currentRole === 'student' ? 'block' : 'none'}; border-top: 1px dashed var(--border-color); padding-top: 1rem; margin-top: 0.75rem;">
                <div class="form-group">
                    <label class="form-label">Школьный класс (для Ученика)</label>
                    <select id="student-classroom-id" class="form-select">
                        <option value="">-- Выберите класс --</option>
                        ${classroomsList.map(c => `
                            <option value="${c.id}" ${user?.student_role?.classroom_id === c.id ? 'selected' : ''}>${c.grade_number}${c.letter || ''} (${c.academic_year || ''})</option>
                        `).join('')}
                    </select>
                </div>
            </div>

            <!-- 2. Parent Role Fields -->
            <div id="parent-role-fields" class="role-fields-panel" style="display: ${currentRole === 'parent' ? 'block' : 'none'}; border-top: 1px dashed var(--border-color); padding-top: 1rem; margin-top: 0.75rem;">
                <div class="form-group">
                    <label class="form-label">ID Детей (для Родителя, через запятую)</label>
                    <input type="text" id="parent-students-id" class="form-input" value="${(user?.parent_role?.students_id || []).join(', ')}" placeholder="например: 5, 8">
                    <div style="font-size: 0.78rem; color: var(--text-subtle); margin-top: 0.25rem;">
                        Укажите списком ID учеников через запятую (например: 5, 8)
                    </div>
                </div>
            </div>

            <!-- 3. Staff Role Fields -->
            <div id="staff-role-fields" class="role-fields-panel" style="display: ${currentRole === 'staff' ? 'block' : 'none'}; border-top: 1px dashed var(--border-color); padding-top: 1rem; margin-top: 0.75rem;">
                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                    <div class="form-group">
                        <label class="form-label">Должность (Staff)</label>
                        <input type="text" id="staff-position" class="form-input" value="${user?.staff_role?.position || ''}" placeholder="Завуч, Директор, Секретарь">
                    </div>
                    <div class="form-group">
                        <label class="form-label">Стаж работы (в днях)</label>
                        <input type="number" id="staff-experience" class="form-input" min="0" value="${user?.staff_role?.experience || 0}" placeholder="365">
                    </div>
                </div>
            </div>

            <!-- 4. Teacher Role Fields -->
            <div id="teacher-role-fields" class="role-fields-panel" style="display: ${currentRole === 'teacher' ? 'block' : 'none'}; border-top: 1px dashed var(--border-color); padding-top: 1rem; margin-top: 0.75rem;">
                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                    <div class="form-group">
                        <label class="form-label">ID Предметов (через запятую)</label>
                        <input type="text" id="teacher-subjects-id" class="form-input" value="${(user?.teacher_role?.subjects_id || []).join(', ')}" placeholder="например: 1, 3">
                        <div style="font-size: 0.78rem; color: var(--text-subtle); margin-top: 0.25rem;">
                            Доступные предметы: ${subjectsList.map(s => `${s.name} (ID: ${s.id})`).join(', ') || 'нет'}
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Педагогический стаж (в днях)</label>
                        <input type="number" id="teacher-experience" class="form-input" min="0" value="${user?.teacher_role?.experience || 0}" placeholder="365">
                    </div>
                </div>
            </div>
        </form>
    `;

    const footerHTML = `
        <button type="button" class="btn btn-secondary" onclick="document.getElementById('modal-close-btn').click()">Отмена</button>
        <button type="button" id="save-user-submit" class="btn btn-primary">${isEdit ? 'Сохранить изменения' : 'Создать пользователя'}</button>
    `;

    openModal({
        title: isEdit ? `Редактирование пользователя #${user.id}` : 'Создание нового пользователя',
        bodyHTML,
        footerHTML
    });

    // Toggle role fields dynamically
    const roleSelect = document.getElementById('user-role');
    roleSelect.addEventListener('change', (e) => {
        const val = e.target.value;
        document.querySelectorAll('.role-fields-panel').forEach(el => el.style.display = 'none');
        if (val === 'student') document.getElementById('student-role-fields').style.display = 'block';
        if (val === 'parent') document.getElementById('parent-role-fields').style.display = 'block';
        if (val === 'staff') document.getElementById('staff-role-fields').style.display = 'block';
        if (val === 'teacher') document.getElementById('teacher-role-fields').style.display = 'block';
    });

    document.getElementById('save-user-submit').addEventListener('click', async () => {
        const name = document.getElementById('user-name').value;
        const username = document.getElementById('user-username').value;
        const email = document.getElementById('user-email').value;
        const role = document.getElementById('user-role').value;
        const phone_number = document.getElementById('user-phone').value || undefined;
        const birthday = document.getElementById('user-birthday').value || undefined;

        let student_role = undefined;
        let parent_role = undefined;
        let staff_role = undefined;
        let teacher_role = undefined;

        if (role === 'student') {
            const classVal = document.getElementById('student-classroom-id').value;
            if (classVal) {
                student_role = { classroom_id: Number(classVal) };
            }
        } else if (role === 'parent') {
            const kidsStr = document.getElementById('parent-students-id').value;
            if (kidsStr) {
                const ids = kidsStr.split(',').map(s => Number(s.trim())).filter(n => !isNaN(n) && n > 0);
                if (ids.length > 0) {
                    parent_role = { students_id: ids };
                }
            }
        } else if (role === 'staff') {
            const position = document.getElementById('staff-position').value;
            const experience = Number(document.getElementById('staff-experience').value) || 0;
            staff_role = { position, experience };
        } else if (role === 'teacher') {
            const subjStr = document.getElementById('teacher-subjects-id').value;
            const experience = Number(document.getElementById('teacher-experience').value) || 0;
            let subjects_id = [];
            if (subjStr) {
                subjects_id = subjStr.split(',').map(s => Number(s.trim())).filter(n => !isNaN(n) && n > 0);
            }
            teacher_role = { subjects_id, experience };
        }

        try {
            if (isEdit) {
                await api.updateUser({
                    id: user.id,
                    name,
                    username,
                    email,
                    role,
                    phone_number,
                    birthday,
                    student_role,
                    parent_role,
                    staff_role,
                    teacher_role
                });
                showToast('Данные пользователя обновлены', 'success');
            } else {
                const password = document.getElementById('user-password').value;
                await api.createUser({
                    name,
                    username,
                    email,
                    role,
                    password,
                    phone_number,
                    birthday,
                    student_role,
                    parent_role,
                    staff_role,
                    teacher_role
                });
                showToast('Новый пользователь успешно создан', 'success');
            }
            closeModal();
            await loadData();
        } catch (err) {
            showToast(err.message, 'error');
        }
    });
}
