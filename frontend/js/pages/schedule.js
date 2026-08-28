import { api } from '../api.js';
import { state } from '../state.js';
import { DAYS_OF_WEEK } from '../utils.js';
import { showToast } from '../components/toast.js';
import { openModal, closeModal } from '../components/modal.js';

let currentSchedule = [];
let classrooms = [];
let classroomsMap = {};
let usersList = [];
let usersMap = {};
let subjectsList = [];
let subjectsMap = {};

async function ensureUsersLoaded() {
    if (usersList.length > 0) return;
    try {
        usersList = await api.getUsersList();
        if (Array.isArray(usersList)) {
            usersMap = {};
            usersList.forEach(u => {
                usersMap[u.id] = u.name || u.username || `Пользователь #${u.id}`;
            });
        }
    } catch (e) {
        console.warn('Failed to load users for schedule names mapping:', e);
    }
}

async function ensureSubjectsLoaded() {
    try {
        subjectsList = await api.getSubjectsList();
        if (Array.isArray(subjectsList)) {
            subjectsMap = {};
            subjectsList.forEach(s => {
                subjectsMap[s.id] = s.name;
            });
        }
    } catch (e) {
        console.warn('Failed to load subjects list for schedule:', e);
    }
}

function updateClassroomsMap() {
    classroomsMap = {};
    if (Array.isArray(classrooms)) {
        classrooms.forEach(c => {
            classroomsMap[c.id] = `${c.grade_number}${c.letter || ''}`;
        });
    }
}

export async function renderSchedulePage() {
    const mainContent = document.getElementById('main-content');
    if (!mainContent) return;

    const userRole = state.user?.role?.toLowerCase() || 'user';
    const isStaff = userRole === 'staff';

    mainContent.innerHTML = `
        <div class="page-header">
            <div>
                <h1 class="page-title">Расписание занятий</h1>
                <p class="page-subtitle">Сетка уроков по дням недели, классам и предметам</p>
            </div>
            ${isStaff ? `
            <button id="add-schedule-btn" class="btn btn-primary">
                <span>➕</span> Добавить урок
            </button>
            ` : ''}
        </div>

        <div class="glass-card" style="margin-bottom: 1.5rem; padding: 1rem 1.25rem;">
            <div style="display: flex; gap: 1rem; flex-wrap: wrap; align-items: center;">
                <div style="display: flex; align-items: center; gap: 0.5rem;">
                    <label class="form-label" style="margin: 0;">Класс:</label>
                    <select id="schedule-classroom-select" class="form-select" style="padding: 0.45rem 0.85rem; min-width: 150px;">
                        <option value="">Выберите класс</option>
                    </select>
                </div>

                <div style="display: flex; align-items: center; gap: 0.5rem;">
                    <label class="form-label" style="margin: 0;">или ID Учителя:</label>
                    <input type="number" id="schedule-teacher-input" class="form-input" placeholder="ID Учителя" style="padding: 0.45rem 0.85rem; width: 140px;">
                    <button id="search-teacher-schedule-btn" class="btn btn-secondary btn-sm">Найти</button>
                </div>
            </div>
        </div>

        <div id="schedule-view-container">
            <div class="empty-state">
                <div class="empty-state-icon">📅</div>
                <p>Выберите класс или введите ID учителя для просмотра расписания</p>
            </div>
        </div>
    `;

    if (isStaff) {
        document.getElementById('add-schedule-btn')?.addEventListener('click', () => {
            openScheduleModal();
        });
    }

    // Preload users & subjects maps
    await Promise.all([
        ensureUsersLoaded(),
        ensureSubjectsLoaded()
    ]);

    // Load classrooms list for dropdown and classrooms map
    try {
        classrooms = await api.getClassroomsList();
        updateClassroomsMap();

        const select = document.getElementById('schedule-classroom-select');
        if (select && classrooms) {
            classrooms.forEach(c => {
                const opt = document.createElement('option');
                opt.value = c.id;
                opt.textContent = `${c.grade_number}${c.letter || ''} (${c.academic_year || ''})`;
                select.appendChild(opt);
            });

            select.addEventListener('change', (e) => {
                const val = e.target.value;
                if (val) loadClassroomSchedule(val);
            });
            
            // Auto select first classroom if available
            if (classrooms.length > 0) {
                select.value = classrooms[0].id;
                loadClassroomSchedule(classrooms[0].id);
            }
        }
    } catch (e) {
        console.warn('Failed to load classrooms for schedule dropdown:', e);
    }

    document.getElementById('search-teacher-schedule-btn')?.addEventListener('click', () => {
        const teacherId = document.getElementById('schedule-teacher-input').value;
        if (teacherId) loadTeacherSchedule(teacherId);
    });
}

async function loadClassroomSchedule(classroomId) {
    const container = document.getElementById('schedule-view-container');
    if (!container) return;

    container.innerHTML = `
        <div class="loading-spinner-wrapper" style="min-height: 200px;">
            <div class="spinner"></div>
            <p>Загрузка расписания класса...</p>
        </div>
    `;

    try {
        await Promise.all([ensureUsersLoaded(), ensureSubjectsLoaded()]);
        currentSchedule = await api.getScheduleByClassroom(classroomId);
        renderScheduleGrid();
    } catch (err) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">⚠️</div>
                <p>Не удалось загрузить расписание (${err.message})</p>
            </div>
        `;
    }
}

async function loadTeacherSchedule(teacherId) {
    const container = document.getElementById('schedule-view-container');
    if (!container) return;

    container.innerHTML = `
        <div class="loading-spinner-wrapper" style="min-height: 200px;">
            <div class="spinner"></div>
            <p>Загрузка расписания учителя...</p>
        </div>
    `;

    try {
        await Promise.all([ensureUsersLoaded(), ensureSubjectsLoaded()]);
        currentSchedule = await api.getScheduleByTeacher(teacherId);
        renderScheduleGrid();
    } catch (err) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">⚠️</div>
                <p>Не удалось загрузить расписание учителя (${err.message})</p>
            </div>
        `;
    }
}

function renderScheduleGrid() {
    const container = document.getElementById('schedule-view-container');
    if (!container) return;

    if (!currentSchedule || currentSchedule.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">📅</div>
                <p>Уроков в расписании пока нет</p>
            </div>
        `;
        return;
    }

    const isStaff = state.user?.role?.toLowerCase() === 'staff';

    // Group by day of week (1 to 7)
    let daysHTML = DAYS_OF_WEEK.map(day => {
        const lessons = currentSchedule
            .filter(item => Number(item.day_of_week) === day.id)
            .sort((a, b) => a.lesson_number - b.lesson_number);

        if (lessons.length === 0) return '';

        const lessonsListHTML = lessons.map(l => {
            const subjectName = subjectsMap[l.subject_id] || `Предмет #${l.subject_id}`;
            const classroomLabel = classroomsMap[l.classroom_id] ? `Класс ${classroomsMap[l.classroom_id]}` : `Класс #${l.classroom_id}`;
            const teacherName = usersMap[l.teacher_id] || `#${l.teacher_id}`;
            
            return `
            <div class="glass-card" style="padding: 0.85rem 1rem; margin-bottom: 0.75rem; border-left: 3px solid var(--primary);">
                <div style="display: flex; justify-content: space-between; align-items: center;">
                    <div style="font-weight: 600; font-size: 0.95rem; color: var(--text-main);">
                        Урок №${l.lesson_number} — ${subjectName}
                    </div>
                    ${isStaff ? `
                    <div style="display:flex; gap:0.25rem;">
                        <button class="icon-btn edit-sched-btn" data-id="${l.id}" title="Редактировать">✏️</button>
                        <button class="icon-btn delete-sched-btn" data-id="${l.id}" title="Удалить">🗑️</button>
                    </div>
                    ` : ''}
                </div>
                <div style="font-size: 0.82rem; color: var(--text-subtle); margin-top: 0.35rem; display: flex; flex-wrap: wrap; gap: 0.5rem; align-items: center;">
                    <span>👥 <strong>${classroomLabel}</strong></span>
                    <span>•</span>
                    <span>👨‍🏫 Учитель: <strong>${teacherName}</strong></span>
                    ${l.room ? `<span>•</span> <span>🚪 Кабинет: <strong>${l.room}</strong></span>` : ''}
                </div>
            </div>
            `;
        }).join('');

        return `
            <div class="glass-card" style="padding: 1.25rem;">
                <h3 style="font-family: var(--font-heading); font-size: 1.1rem; color: var(--accent); margin-bottom: 1rem; border-bottom: 1px solid var(--border-color); padding-bottom: 0.5rem;">
                    ${day.name}
                </h3>
                ${lessonsListHTML}
            </div>
        `;
    }).join('');

    container.innerHTML = `
        <div class="grid-cols-3">
            ${daysHTML}
        </div>
    `;

    if (isStaff) {
        container.querySelectorAll('.edit-sched-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const id = Number(e.currentTarget.getAttribute('data-id'));
                const sched = currentSchedule.find(s => s.id === id);
                if (sched) openScheduleModal(sched);
            });
        });

        container.querySelectorAll('.delete-sched-btn').forEach(btn => {
            btn.addEventListener('click', async (e) => {
                const id = Number(e.currentTarget.getAttribute('data-id'));
                if (confirm(`Удалить этот урок из расписания?`)) {
                    try {
                        await api.deleteSchedule(id);
                        showToast('Урок удален из расписания', 'success');
                        const classroomId = document.getElementById('schedule-classroom-select')?.value;
                        if (classroomId) await loadClassroomSchedule(classroomId);
                    } catch (err) {
                        showToast(err.message, 'error');
                    }
                }
            });
        });
    }
}

function openScheduleModal(sched = null) {
    const isEdit = !!sched;

    // Filter teachers if usersList is populated
    const teachersList = usersList.filter(u => !u.role || u.role.toLowerCase() === 'teacher' || u.role.toLowerCase() === 'staff');
    const teacherOptionsList = teachersList.length > 0 ? teachersList : usersList;

    let teacherSelectHTML = '';
    if (teacherOptionsList.length > 0) {
        teacherSelectHTML = `
            <select id="sched-teacher-id" class="form-select">
                ${teacherOptionsList.map(u => `
                    <option value="${u.id}" ${sched?.teacher_id === u.id ? 'selected' : ''}>${u.name || u.username} (ID: ${u.id})</option>
                `).join('')}
            </select>
        `;
    } else {
        teacherSelectHTML = `<input type="number" id="sched-teacher-id" class="form-input" required value="${sched?.teacher_id || 1}">`;
    }

    let classroomSelectHTML = '';
    if (Array.isArray(classrooms) && classrooms.length > 0) {
        classroomSelectHTML = `
            <select id="sched-classroom-id" class="form-select">
                ${classrooms.map(c => `
                    <option value="${c.id}" ${sched?.classroom_id === c.id ? 'selected' : ''}>${c.grade_number}${c.letter || ''} (${c.academic_year || ''})</option>
                `).join('')}
            </select>
        `;
    } else {
        classroomSelectHTML = `<input type="number" id="sched-classroom-id" class="form-input" required value="${sched?.classroom_id || 1}">`;
    }

    let subjectSelectHTML = '';
    if (Array.isArray(subjectsList) && subjectsList.length > 0) {
        subjectSelectHTML = `
            <select id="sched-subject-id" class="form-select">
                ${subjectsList.map(s => `
                    <option value="${s.id}" ${sched?.subject_id === s.id ? 'selected' : ''}>${s.name} (ID: ${s.id})</option>
                `).join('')}
            </select>
        `;
    } else {
        subjectSelectHTML = `<input type="number" id="sched-subject-id" class="form-input" required value="${sched?.subject_id || 1}">`;
    }

    const bodyHTML = `
        <form id="schedule-form">
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                <div class="form-group">
                    <label class="form-label">Класс</label>
                    ${classroomSelectHTML}
                </div>
                <div class="form-group">
                    <label class="form-label">Предмет</label>
                    ${subjectSelectHTML}
                </div>
            </div>

            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                <div class="form-group">
                    <label class="form-label">Преподаватель</label>
                    ${teacherSelectHTML}
                </div>
                <div class="form-group">
                    <label class="form-label">День недели (1-7)</label>
                    <select id="sched-day" class="form-select">
                        ${DAYS_OF_WEEK.map(d => `
                            <option value="${d.id}" ${sched?.day_of_week === d.id ? 'selected' : ''}>${d.name}</option>
                        `).join('')}
                    </select>
                </div>
            </div>

            <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 1rem;">
                <div class="form-group">
                    <label class="form-label">Номер урока (1-8)</label>
                    <input type="number" id="sched-lesson-num" class="form-input" min="1" max="8" required value="${sched?.lesson_number || 1}">
                </div>
                <div class="form-group">
                    <label class="form-label">Кабинет</label>
                    <input type="number" id="sched-room" class="form-input" value="${sched?.room || ''}" placeholder="101">
                </div>
                <div class="form-group">
                    <label class="form-label">Учебный год</label>
                    <input type="text" id="sched-year" class="form-input" required value="${sched?.academic_year || '2024.1'}">
                </div>
            </div>
        </form>
    `;

    const footerHTML = `
        <button type="button" class="btn btn-secondary" onclick="document.getElementById('modal-close-btn').click()">Отмена</button>
        <button type="button" id="save-sched-btn" class="btn btn-primary">${isEdit ? 'Сохранить урок' : 'Добавить урок'}</button>
    `;

    openModal({
        title: isEdit ? `Редактирование урока #${sched.id}` : 'Добавление урока в расписание',
        bodyHTML,
        footerHTML
    });

    document.getElementById('save-sched-btn').addEventListener('click', async () => {
        const classroom_id = Number(document.getElementById('sched-classroom-id').value);
        const subject_id = Number(document.getElementById('sched-subject-id').value);
        const teacher_id = Number(document.getElementById('sched-teacher-id').value);
        const day_of_week = Number(document.getElementById('sched-day').value);
        const lesson_number = Number(document.getElementById('sched-lesson-num').value);
        const roomVal = document.getElementById('sched-room').value;
        const room = roomVal ? Number(roomVal) : null;
        const academic_year = document.getElementById('sched-year').value;

        try {
            if (isEdit) {
                await api.updateSchedule({
                    id: sched.id,
                    classroom_id,
                    subject_id,
                    teacher_id,
                    day_of_week,
                    lesson_number,
                    room,
                    academic_year
                });
                showToast('Урок успешно обновлен', 'success');
            } else {
                await api.createSchedule({
                    classroom_id,
                    subject_id,
                    teacher_id,
                    day_of_week,
                    lesson_number,
                    room,
                    academic_year
                });
                showToast('Урок успешно добавлен в расписание', 'success');
            }
            closeModal();
            const classroomId = document.getElementById('schedule-classroom-select')?.value;
            if (classroomId) await loadClassroomSchedule(classroomId);
        } catch (err) {
            showToast(err.message, 'error');
        }
    });
}
