import { CONFIG } from './config.js';
import { state } from './state.js';
import { showToast } from './components/toast.js';

class ApiClient {
    async request(endpoint, options = {}) {
        const url = `${CONFIG.API_BASE_URL}${endpoint}`;
        
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };

        if (state.token) {
            headers['Authorization'] = `Bearer ${state.token}`;
        }

        const config = {
            ...options,
            headers
        };

        try {
            const response = await fetch(url, config);
            
            // Handle 401 Unauthorized
            if (response.status === 401) {
                showToast('Сессия истекла. Пожалуйста, войдите снова.', 'error');
                state.logout();
                throw new Error('Unauthorized');
            }

            // Handle 403 Permission Denied
            if (response.status === 403) {
                showToast('Доступ запрещен (недостаточно прав).', 'error');
                throw new Error('Forbidden');
            }

            const data = await response.json().catch(() => ({}));

            if (!response.ok) {
                const errMsg = data.message || data.error || `Ошибка сервера (${response.status})`;
                throw new Error(errMsg);
            }

            return data;
        } catch (error) {
            if (error.message !== 'Unauthorized' && error.message !== 'Forbidden') {
                showToast(error.message || 'Ошибка подключения к серверу', 'error');
            }
            throw error;
        }
    }

    // Auth
    login(email, password) {
        return this.request('/api/auth/login', {
            method: 'POST',
            body: JSON.stringify({ email, password })
        });
    }

    // Users
    getUserById(id) {
        return this.request(`/api/user/${id}`, { method: 'GET' });
    }

    getUsersList() {
        return this.request('/api/user/list', { method: 'GET' });
    }

    createUser(userData) {
        return this.request('/api/user', {
            method: 'POST',
            body: JSON.stringify(userData)
        });
    }

    updateUser(userData) {
        return this.request('/api/user', {
            method: 'PATCH',
            body: JSON.stringify(userData)
        });
    }

    deleteUser(id) {
        return this.request(`/api/user/${id}`, { method: 'DELETE' });
    }

    getSubjectById(id) {
        return this.request(`/api/user/subject/${id}`, { method: 'GET' });
    }

    getSubjectsList() {
        return this.request('/api/user/subject/list', { method: 'GET' });
    }

    createSubject(subjectData) {
        return this.request('/api/user/subject', {
            method: 'POST',
            body: JSON.stringify(subjectData)
        });
    }

    // Classrooms
    getClassroomById(id) {
        return this.request(`/api/classroom/${id}`, { method: 'GET' });
    }

    getClassroomsList() {
        return this.request('/api/classroom/list', { method: 'GET' });
    }

    createClassroom(classroomData) {
        return this.request('/api/classroom', {
            method: 'POST',
            body: JSON.stringify(classroomData)
        });
    }

    updateClassroom(classroomData) {
        return this.request('/api/classroom', {
            method: 'PATCH',
            body: JSON.stringify(classroomData)
        });
    }

    deleteClassroom(id) {
        return this.request(`/api/classroom/${id}`, { method: 'DELETE' });
    }

    // Schedule
    getScheduleByClassroom(classroomId) {
        return this.request(`/api/schedule/classroom/${classroomId}`, { method: 'GET' });
    }

    getScheduleByTeacher(teacherId) {
        return this.request(`/api/schedule/teacher/${teacherId}`, { method: 'GET' });
    }

    createSchedule(scheduleData) {
        return this.request('/api/schedule', {
            method: 'POST',
            body: JSON.stringify(scheduleData)
        });
    }

    updateSchedule(scheduleData) {
        return this.request('/api/schedule', {
            method: 'PATCH',
            body: JSON.stringify(scheduleData)
        });
    }

    deleteSchedule(id) {
        return this.request(`/api/schedule/${id}`, { method: 'DELETE' });
    }

    // Journal / Rating
    getJournalByStudent(studentId, dateRangeStr) {
        return this.request(`/api/journal/student/${studentId}`, {
            method: 'POST',
            body: JSON.stringify({ date_range: dateRangeStr })
        });
    }

    getJournalByClassroom(classroomId, dateRangeStr) {
        return this.request(`/api/journal/classroom/${classroomId}`, {
            method: 'POST',
            body: JSON.stringify({ date_range: dateRangeStr })
        });
    }

    updateJournal(journalData) {
        return this.request('/api/journal', {
            method: 'PATCH',
            body: JSON.stringify(journalData)
        });
    }
}

export const api = new ApiClient();
