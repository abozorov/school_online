import { api } from '../api.js';
import { showToast } from '../components/toast.js';
import { openModal, closeModal } from '../components/modal.js';

let classroomsList = [];
let usersList = [];

export async function renderClassroomsPage() {
    const mainContent = document.getElementById('main-content');
    if (!mainContent) return;

    mainContent.innerHTML = `
        <div class="page-header">
            <div>
                <h1 class="page-title">Школьные классы</h1>
                <p class="page-subtitle">Формирование классов, списки учеников и привязка классных руководителей</p>
            </div>
            <button id="add-classroom-btn" class="btn btn-primary">
                <span>➕</span> Создать класс
            </button>
        </div>

        <div id="classrooms-grid-container">
            <div class="loading-spinner-wrapper" style="min-height: 200px;">
                <div class="spinner"></div>
                <p>Загрузка списка классов...</p>
            </div>
        </div>
    `;

    document.getElementById('add-classroom-btn').addEventListener('click', () => {
        openClassroomModal();
    });

    await loadData();
}

async function loadData() {
    try {
        const [classroomsData, usersData] = await Promise.all([
            api.getClassroomsList(),
            api.getUsersList().catch(() => [])
        ]);
        classroomsList = classroomsData || [];
        usersList = usersData || [];
        renderClassroomsGrid();
    } catch (error) {
        console.error('Failed to load classrooms data:', error);
        document.getElementById('classrooms-grid-container').innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">🏫</div>
                <p>Не удалось загрузить классы (${error.message})</p>
            </div>
        `;
    }
}

function renderClassroomsGrid() {
    const container = document.getElementById('classrooms-grid-container');
    if (!container) return;

    if (!classroomsList || classroomsList.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">🏫</div>
                <p>Список классов пуст. Создайте первый класс!</p>
            </div>
        `;
        return;
    }

    let cardsHTML = classroomsList.map(c => {
        const className = `${c.grade_number}${c.letter || ''}`;
        
        // Find students assigned to this classroom
        const classroomStudents = usersList.filter(u => {
            const isStudent = u.role?.toLowerCase() === 'student' || !!u.student_role;
            const classId = u.student_role?.classroom_id || u.classroom_id;
            return isStudent && Number(classId) === Number(c.id);
        });

        // Find Hometown Teacher Name
        const teacherUser = usersList.find(u => Number(u.id) === Number(c.hometown_teacher_id));
        const teacherName = teacherUser ? (teacherUser.name || `@${teacherUser.username}`) : (c.hometown_teacher_id ? `#${c.hometown_teacher_id}` : 'Не назначен');

        return `
            <div class="glass-card" style="position: relative; display: flex; flex-direction: column; justify-content: space-between;">
                <div>
                    <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 1rem;">
                        <div>
                            <div style="font-family: var(--font-heading); font-size: 1.75rem; font-weight: 700; color: var(--primary);">
                                ${className}
                            </div>
                            <div style="font-size: 0.85rem; color: var(--text-muted);">
                                Учебный год: <strong>${c.academic_year || '2024.1'}</strong>
                            </div>
                        </div>
                        <div style="display: flex; gap: 0.35rem;">
                            <button class="icon-btn edit-classroom-btn" data-id="${c.id}" title="Редактировать">✏️</button>
                            <button class="icon-btn delete-classroom-btn" data-id="${c.id}" title="Удалить">🗑️</button>
                        </div>
                    </div>

                    <div style="font-size: 0.9rem; color: var(--text-subtle); padding-top: 0.75rem; border-top: 1px solid var(--border-color);">
                        <div style="margin-bottom: 0.4rem;">👨‍🏫 Руководитель: <strong>${teacherName}</strong></div>
                        <div>👥 Состав: <strong>${classroomStudents.length} учеников</strong></div>
                    </div>
                </div>

                <div style="margin-top: 1.25rem; padding-top: 0.75rem; border-top: 1px dashed var(--border-color);">
                    <button class="btn btn-secondary btn-sm view-students-btn" data-id="${c.id}" style="width: 100%;">
                        👥 Ученики класса (${classroomStudents.length})
                    </button>
                </div>
            </div>
        `;
    }).join('');

    container.innerHTML = `
        <div class="grid-cols-3">
            ${cardsHTML}
        </div>
    `;

    // Event listeners
    container.querySelectorAll('.view-students-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const id = Number(e.currentTarget.getAttribute('data-id'));
            const cls = classroomsList.find(item => item.id === id);
            if (cls) openStudentsListModal(cls);
        });
    });

    container.querySelectorAll('.edit-classroom-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const id = Number(e.currentTarget.getAttribute('data-id'));
            const cls = classroomsList.find(item => item.id === id);
            if (cls) openClassroomModal(cls);
        });
    });

    container.querySelectorAll('.delete-classroom-btn').forEach(btn => {
        btn.addEventListener('click', async (e) => {
            const id = Number(e.currentTarget.getAttribute('data-id'));
            if (confirm(`Удалить класс #${id}?`)) {
                try {
                    await api.deleteClassroom(id);
                    showToast('Класс успешно удален', 'success');
                    await loadData();
                } catch (err) {
                    showToast(err.message, 'error');
                }
            }
        });
    });
}

function openStudentsListModal(classroom) {
    const className = `${classroom.grade_number}${classroom.letter || ''}`;
    
    const students = usersList.filter(u => {
        const isStudent = u.role?.toLowerCase() === 'student' || !!u.student_role;
        const classId = u.student_role?.classroom_id || u.classroom_id;
        return isStudent && Number(classId) === Number(classroom.id);
    });

    let bodyHTML = '';
    if (students.length === 0) {
        bodyHTML = `
            <div class="empty-state" style="padding: 2rem 1rem;">
                <div class="empty-state-icon">🎓</div>
                <p>В классе <strong>${className}</strong> пока нет зарегистрированных учеников.</p>
            </div>
        `;
    } else {
        const rowsHTML = students.map(s => `
            <tr>
                <td><strong>#${s.id}</strong></td>
                <td>
                    <div style="font-weight: 600;">${s.name || '—'}</div>
                    <div style="font-size: 0.8rem; color: var(--text-subtle);">@${s.username || 'student'}</div>
                </td>
                <td>${s.email || '—'}</td>
                <td>${s.phone_number || '—'}</td>
                <td>${s.birthday || '—'}</td>
            </tr>
        `).join('');

        bodyHTML = `
            <div style="margin-bottom: 1rem; font-size: 0.95rem; color: var(--text-muted);">
                Всего учеников в классе <strong>${className}</strong>: <strong>${students.length}</strong>
            </div>
            <div class="table-container">
                <table class="data-table">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>ФИО / Логин</th>
                            <th>Email</th>
                            <th>Телефон</th>
                            <th>Дата рождения</th>
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
        title: `Ученики класса ${className}`,
        bodyHTML,
        footerHTML
    });
}

function openClassroomModal(classroom = null) {
    const isEdit = !!classroom;

    // Filter teachers list
    const teachersList = usersList.filter(u => !u.role || u.role.toLowerCase() === 'teacher' || u.role.toLowerCase() === 'staff');
    let teacherSelectHTML = '';
    if (teachersList.length > 0) {
        teacherSelectHTML = `
            <select id="classroom-teacher-id" class="form-select">
                <option value="">-- Не назначен --</option>
                ${teachersList.map(u => `
                    <option value="${u.id}" ${classroom?.hometown_teacher_id === u.id ? 'selected' : ''}>${u.name || u.username} (ID: ${u.id})</option>
                `).join('')}
            </select>
        `;
    } else {
        teacherSelectHTML = `<input type="number" id="classroom-teacher-id" class="form-input" value="${classroom?.hometown_teacher_id || ''}" placeholder="ID Учителя">`;
    }

    const bodyHTML = `
        <form id="classroom-form">
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                <div class="form-group">
                    <label class="form-label">Параллель (1-11)</label>
                    <input type="number" id="classroom-grade" class="form-input" min="1" max="11" required value="${classroom?.grade_number || 9}">
                </div>
                <div class="form-group">
                    <label class="form-label">Буква класса</label>
                    <input type="text" id="classroom-letter" class="form-input" maxlength="1" required value="${classroom?.letter || 'А'}" placeholder="А">
                </div>
            </div>

            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                <div class="form-group">
                    <label class="form-label">Классный руководитель</label>
                    ${teacherSelectHTML}
                </div>
                <div class="form-group">
                    <label class="form-label">Учебный год (формат YYYY.N)</label>
                    <input type="text" id="classroom-year" class="form-input" required value="${classroom?.academic_year || '2024.1'}" placeholder="2024.1">
                </div>
            </div>
        </form>
    `;

    const footerHTML = `
        <button type="button" class="btn btn-secondary" onclick="document.getElementById('modal-close-btn').click()">Отмена</button>
        <button type="button" id="save-classroom-btn" class="btn btn-primary">${isEdit ? 'Обновить класс' : 'Создать класс'}</button>
    `;

    openModal({
        title: isEdit ? `Редактирование класса #${classroom.id}` : 'Создание нового класса',
        bodyHTML,
        footerHTML
    });

    document.getElementById('save-classroom-btn').addEventListener('click', async () => {
        const grade_number = Number(document.getElementById('classroom-grade').value);
        const letter = document.getElementById('classroom-letter').value.toUpperCase();
        const teacherVal = document.getElementById('classroom-teacher-id').value;
        const hometown_teacher_id = teacherVal ? Number(teacherVal) : null;
        const academic_year = document.getElementById('classroom-year').value;

        try {
            if (isEdit) {
                await api.updateClassroom({
                    id: classroom.id,
                    grade_number,
                    letter,
                    hometown_teacher_id,
                    academic_year
                });
                showToast('Класс успешно обновлен', 'success');
            } else {
                await api.createClassroom({
                    grade_number,
                    letter,
                    hometown_teacher_id,
                    academic_year
                });
                showToast('Новый класс успешно создан', 'success');
            }
            closeModal();
            await loadData();
        } catch (err) {
            showToast(err.message, 'error');
        }
    });
}
