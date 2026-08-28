import { api } from '../api.js';
import { showToast } from '../components/toast.js';
import { openModal, closeModal } from '../components/modal.js';

let subjectsList = [];

export async function renderSubjectsPage() {
    const mainContent = document.getElementById('main-content');
    if (!mainContent) return;

    mainContent.innerHTML = `
        <div class="page-header">
            <div>
                <h1 class="page-title">Учебные предметы</h1>
                <p class="page-subtitle">Просмотр, создание и редактирование предметов школьной программы</p>
            </div>
            <button id="add-subject-btn" class="btn btn-primary">
                <span>➕</span> Создать предмет
            </button>
        </div>

        <div id="subjects-container">
            <div class="loading-spinner-wrapper" style="min-height: 200px;">
                <div class="spinner"></div>
                <p>Загрузка списка предметов...</p>
            </div>
        </div>
    `;

    document.getElementById('add-subject-btn')?.addEventListener('click', () => {
        openSubjectModal();
    });

    await loadSubjects();
}

async function loadSubjects() {
    try {
        subjectsList = await api.getSubjectsList();
        renderSubjectsGrid();
    } catch (error) {
        console.error('Failed to load subjects:', error);
        document.getElementById('subjects-container').innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">📚</div>
                <p>Не удалось загрузить предметы (${error.message})</p>
            </div>
        `;
    }
}

function renderSubjectsGrid() {
    const container = document.getElementById('subjects-container');
    if (!container) return;

    if (!subjectsList || subjectsList.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">📚</div>
                <p>Список предметов пуст. Создайте первый предмет!</p>
            </div>
        `;
        return;
    }

    const cardsHTML = subjectsList.map(s => `
        <div class="glass-card" style="display: flex; flex-direction: column; justify-content: space-between;">
            <div>
                <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 0.75rem;">
                    <div>
                        <span style="font-size: 0.75rem; background: rgba(99, 102, 241, 0.15); color: var(--primary); padding: 0.2rem 0.5rem; border-radius: 4px; font-weight: 600;">
                            ID: #${s.id}
                        </span>
                        <h3 style="font-family: var(--font-heading); font-size: 1.2rem; font-weight: 700; color: var(--text-main); margin-top: 0.4rem;">
                            ${s.name}
                        </h3>
                    </div>
                    <button class="icon-btn edit-subject-btn" data-id="${s.id}" title="Редактировать">✏️</button>
                </div>
                <p style="font-size: 0.88rem; color: var(--text-muted); line-height: 1.5; margin-bottom: 1rem;">
                    ${s.description || 'Описание отсутствует'}
                </p>
            </div>
        </div>
    `).join('');

    container.innerHTML = `
        <div class="grid-cols-3">
            ${cardsHTML}
        </div>
    `;

    container.querySelectorAll('.edit-subject-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const id = Number(e.currentTarget.getAttribute('data-id'));
            const sub = subjectsList.find(item => item.id === id);
            if (sub) openSubjectModal(sub);
        });
    });
}

function openSubjectModal(subject = null) {
    const isEdit = !!subject;

    const bodyHTML = `
        <form id="subject-modal-form">
            <div class="form-group">
                <label class="form-label" for="subject-name">Название предмета</label>
                <input type="text" id="subject-name" class="form-input" required value="${subject?.name || ''}" placeholder="например: Алгебра и начала анализа">
            </div>

            <div class="form-group">
                <label class="form-label" for="subject-desc">Описание предмета</label>
                <textarea id="subject-desc" class="form-textarea" rows="4" placeholder="Краткая аннотация курса...">${subject?.description || ''}</textarea>
            </div>
        </form>
    `;

    const footerHTML = `
        <button type="button" class="btn btn-secondary" onclick="document.getElementById('modal-close-btn').click()">Отмена</button>
        <button type="button" id="save-subject-btn" class="btn btn-primary">${isEdit ? 'Сохранить изменения' : 'Создать предмет'}</button>
    `;

    openModal({
        title: isEdit ? `Редактирование предмета #${subject.id}` : 'Создание нового предмета',
        bodyHTML,
        footerHTML
    });

    document.getElementById('save-subject-btn').addEventListener('click', async () => {
        const name = document.getElementById('subject-name').value;
        const description = document.getElementById('subject-desc').value;

        try {
            if (isEdit) {
                await api.updateSubject({ id: subject.id, name, description });
                showToast(`Предмет "${name}" успешно обновлен`, 'success');
            } else {
                const res = await api.createSubject({ name, description });
                showToast(`Предмет "${name}" успешно создан! (ID: ${res.id})`, 'success');
            }
            closeModal();
            await loadSubjects();
        } catch (err) {
            showToast(err.message, 'error');
        }
    });
}
