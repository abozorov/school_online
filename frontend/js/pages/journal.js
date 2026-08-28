import { api } from '../api.js';
import { state } from '../state.js';
import { getDefaultDateRange, formatDate } from '../utils.js';
import { showToast } from '../components/toast.js';
import { openModal, closeModal } from '../components/modal.js';

let journalEntries = [];
let classrooms = [];

export async function renderJournalPage() {
    const mainContent = document.getElementById('main-content');
    if (!mainContent) return;

    const userRole = state.user?.role?.toLowerCase() || 'user';
    const isTeacher = userRole === 'teacher';
    const isStaff = userRole === 'staff';
    const isStudentOrParent = userRole === 'student' || userRole === 'parent';

    mainContent.innerHTML = `
        <div class="page-header">
            <div>
                <h1 class="page-title">Электронный журнал и успеваемость</h1>
                <p class="page-subtitle">Просмотр оценок, выданных домашних заданий и отметки посещаемости</p>
            </div>
            ${(isTeacher || isStaff) ? `
            <button id="add-grade-btn" class="btn btn-primary">
                <span>✏️</span> Поставить оценку / ДЗ
            </button>
            ` : ''}
        </div>

        <div class="glass-card" style="margin-bottom: 1.5rem; padding: 1rem 1.25rem;">
            <div style="display: flex; gap: 1rem; flex-wrap: wrap; align-items: center; justify-content: space-between;">
                
                <div style="display: flex; gap: 1rem; flex-wrap: wrap; align-items: center;">
                    ${!isStudentOrParent ? `
                    <div style="display: flex; align-items: center; gap: 0.5rem;">
                        <label class="form-label" style="margin: 0;">Класс:</label>
                        <select id="journal-classroom-select" class="form-select" style="padding: 0.45rem 0.85rem; min-width: 150px;">
                            <option value="">Выберите класс</option>
                        </select>
                    </div>
                    ` : ''}

                    <div style="display: flex; align-items: center; gap: 0.5rem;">
                        <label class="form-label" style="margin: 0;">ID Ученика:</label>
                        <input type="number" id="journal-student-id-input" class="form-input" placeholder="ID Ученика" style="padding: 0.45rem 0.85rem; width: 130px;" value="${isStudentOrParent ? (state.user?.id || 1) : ''}">
                    </div>

                    <div style="display: flex; align-items: center; gap: 0.5rem;">
                        <label class="form-label" style="margin: 0;">Период:</label>
                        <input type="text" id="journal-date-range-input" class="form-input" placeholder="01.01.2024-26.08.2026" style="padding: 0.45rem 0.85rem; width: 210px;" value="${getDefaultDateRange()}">
                    </div>
                </div>

                <button id="load-journal-btn" class="btn btn-primary btn-sm">
                    🔍 Показать журнал
                </button>
            </div>
        </div>

        <div id="journal-table-container" class="table-container">
            <div class="empty-state">
                <div class="empty-state-icon">📖</div>
                <p>Нажмите "Показать журнал" для загрузки ведомости</p>
            </div>
        </div>
    `;

    if (isTeacher || isStaff) {
        document.getElementById('add-grade-btn')?.addEventListener('click', () => {
            openGradeModal();
        });

        // Load classrooms list
        try {
            classrooms = await api.getClassroomsList();
            const select = document.getElementById('journal-classroom-select');
            if (select && classrooms) {
                classrooms.forEach(c => {
                    const opt = document.createElement('option');
                    opt.value = c.id;
                    opt.textContent = `${c.grade_number}${c.letter || ''}`;
                    select.appendChild(opt);
                });
            }
        } catch (e) {
            console.warn('Failed to load classrooms for journal:', e);
        }
    }

    document.getElementById('load-journal-btn').addEventListener('click', fetchJournalData);

    // Auto-load for student/parent
    if (isStudentOrParent) {
        fetchJournalData();
    }
}

async function fetchJournalData() {
    const container = document.getElementById('journal-table-container');
    if (!container) return;

    const classroomSelect = document.getElementById('journal-classroom-select');
    const studentInput = document.getElementById('journal-student-id-input');
    const dateRangeInput = document.getElementById('journal-date-range-input');

    const classroomId = classroomSelect ? classroomSelect.value : null;
    const studentId = studentInput ? studentInput.value : null;
    const dateRange = dateRangeInput ? dateRangeInput.value.trim() : getDefaultDateRange();

    if (!classroomId && !studentId) {
        showToast('Укажите класс или ID ученика для загрузки журнала', 'warning');
        return;
    }

    container.innerHTML = `
        <div class="loading-spinner-wrapper" style="min-height: 200px;">
            <div class="spinner"></div>
            <p>Загрузка записей журнала...</p>
        </div>
    `;

    try {
        if (studentId) {
            journalEntries = await api.getJournalByStudent(studentId, dateRange);
        } else if (classroomId) {
            journalEntries = await api.getJournalByClassroom(classroomId, dateRange);
        }

        renderJournalTable();
    } catch (err) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">⚠️</div>
                <p>Не удалось загрузить журнал (${err.message})</p>
                <div style="font-size: 0.82rem; color: var(--text-subtle); margin-top: 0.5rem;">
                    Проверьте формат периода дат (должен быть: DD.MM.YYYY-DD.MM.YYYY)
                </div>
            </div>
        `;
    }
}

function renderJournalTable() {
    const container = document.getElementById('journal-table-container');
    if (!container) return;

    // Handle payload if wrapped in list or object
    const list = Array.isArray(journalEntries) ? journalEntries : (journalEntries?.items || []);

    if (list.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">📖</div>
                <p>Записи за указанный период отсутствуют</p>
            </div>
        `;
        return;
    }

    const userRole = state.user?.role?.toLowerCase() || 'user';
    const canEdit = userRole === 'teacher' || userRole === 'staff';

    const rowsHTML = list.map(item => {
        const dateStr = item.date ? formatDate(new Date(item.date)) : '—';
        
        let gradeBadge = '—';
        if (item.grade !== undefined && item.grade !== null) {
            gradeBadge = `<span class="badge badge-grade-${item.grade}">${item.grade}</span>`;
        }

        let attendanceStr = '—';
        if (item.attendance === true) attendanceStr = '<span style="color: var(--success);">Присутствовал</span>';
        if (item.attendance === false) attendanceStr = '<span style="color: var(--danger); font-weight:600;">Отсутствовал (Н)</span>';

        return `
            <tr>
                <td><strong>${dateStr}</strong></td>
                <td>Урок №${item.lesson_number || 1}</td>
                <td>Ученик #${item.student_id}</td>
                <td>Предмет #${item.subject_id}</td>
                <td>${gradeBadge}</td>
                <td>${attendanceStr}</td>
                <td>${item.homework || '—'}</td>
                ${canEdit ? `
                <td style="text-align: right;">
                    <button class="btn btn-secondary btn-sm edit-grade-btn" 
                        data-student="${item.student_id}" 
                        data-classroom="${item.classroom_id}"
                        data-subject="${item.subject_id}"
                        data-teacher="${item.teacher_id}"
                        data-lesson="${item.lesson_number}"
                        data-date="${item.date}">✏️ Изменить</button>
                </td>
                ` : ''}
            </tr>
        `;
    }).join('');

    container.innerHTML = `
        <table class="data-table">
            <thead>
                <tr>
                    <th>Дата</th>
                    <th>Урок</th>
                    <th>Ученик</th>
                    <th>Предмет</th>
                    <th>Оценка</th>
                    <th>Посещаемость</th>
                    <th>Домашнее задание</th>
                    ${canEdit ? '<th style="text-align: right;">Действия</th>' : ''}
                </tr>
            </thead>
            <tbody>
                ${rowsHTML}
            </tbody>
        </table>
    `;

    if (canEdit) {
        container.querySelectorAll('.edit-grade-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const target = e.currentTarget;
                openGradeModal({
                    student_id: Number(target.getAttribute('data-student')),
                    classroom_id: Number(target.getAttribute('data-classroom')),
                    subject_id: Number(target.getAttribute('data-subject')),
                    teacher_id: Number(target.getAttribute('data-teacher')),
                    lesson_number: Number(target.getAttribute('data-lesson')),
                    date: target.getAttribute('data-date')
                });
            });
        });
    }
}

function openGradeModal(existing = null) {
    const isEdit = !!existing;

    const bodyHTML = `
        <form id="journal-update-form">
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                <div class="form-group">
                    <label class="form-label">ID Ученика</label>
                    <input type="number" id="j-student-id" class="form-input" required value="${existing?.student_id || 1}">
                </div>
                <div class="form-group">
                    <label class="form-label">ID Класса</label>
                    <input type="number" id="j-classroom-id" class="form-input" required value="${existing?.classroom_id || 1}">
                </div>
            </div>

            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                <div class="form-group">
                    <label class="form-label">ID Предмета</label>
                    <input type="number" id="j-subject-id" class="form-input" required value="${existing?.subject_id || 1}">
                </div>
                <div class="form-group">
                    <label class="form-label">ID Учителя</label>
                    <input type="number" id="j-teacher-id" class="form-input" required value="${existing?.teacher_id || state.user?.id || 1}">
                </div>
            </div>

            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                <div class="form-group">
                    <label class="form-label">Номер урока (1-8)</label>
                    <input type="number" id="j-lesson-num" class="form-input" min="1" max="8" required value="${existing?.lesson_number || 1}">
                </div>
                <div class="form-group">
                    <label class="form-label">Оценка (1-5)</label>
                    <select id="j-grade" class="form-select">
                        <option value="">Без оценки</option>
                        <option value="5">5 (Отлично)</option>
                        <option value="4">4 (Хорошо)</option>
                        <option value="3">3 (Удовлетворительно)</option>
                        <option value="2">2 (Неудовлетворительно)</option>
                    </select>
                </div>
            </div>

            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                <div class="form-group">
                    <label class="form-label">Посещаемость</label>
                    <select id="j-attendance" class="form-select">
                        <option value="true">Присутствовал</option>
                        <option value="false">Отсутствовал (Н)</option>
                    </select>
                </div>
                <div class="form-group">
                    <label class="form-label">Дата урока</label>
                    <input type="date" id="j-date" class="form-input" required value="${existing?.date ? new Date(existing.date).toISOString().split('T')[0] : new Date().toISOString().split('T')[0]}">
                </div>
            </div>

            <div class="form-group">
                <label class="form-label">Домашнее задание</label>
                <textarea id="j-homework" class="form-textarea" rows="3" placeholder="Задания на следующий урок...">${existing?.homework || ''}</textarea>
            </div>
        </form>
    `;

    const footerHTML = `
        <button type="button" class="btn btn-secondary" onclick="document.getElementById('modal-close-btn').click()">Отмена</button>
        <button type="button" id="save-journal-entry-btn" class="btn btn-primary">Сохранить запись</button>
    `;

    openModal({
        title: isEdit ? 'Обновление оценки в журнале' : 'Выставление оценки / ДЗ',
        bodyHTML,
        footerHTML
    });

    document.getElementById('save-journal-entry-btn').addEventListener('click', async () => {
        const student_id = Number(document.getElementById('j-student-id').value);
        const classroom_id = Number(document.getElementById('j-classroom-id').value);
        const subject_id = Number(document.getElementById('j-subject-id').value);
        const teacher_id = Number(document.getElementById('j-teacher-id').value);
        const lesson_number = Number(document.getElementById('j-lesson-num').value);
        
        const gradeVal = document.getElementById('j-grade').value;
        const grade = gradeVal ? Number(gradeVal) : null;
        
        const attVal = document.getElementById('j-attendance').value;
        const attendance = attVal === 'true';
        
        const homework = document.getElementById('j-homework').value || null;
        const dateInputVal = document.getElementById('j-date').value;
        const date = new Date(dateInputVal).toISOString();

        try {
            await api.updateJournal({
                student_id,
                classroom_id,
                subject_id,
                teacher_id,
                lesson_number,
                grade,
                attendance,
                homework,
                date
            });

            showToast('Запись успешно сохранена в журнале', 'success');
            closeModal();
            fetchJournalData();
        } catch (err) {
            showToast(err.message, 'error');
        }
    });
}
