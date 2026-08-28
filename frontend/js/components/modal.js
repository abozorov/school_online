export function openModal({ title, bodyHTML, footerHTML, onClose }) {
    const overlay = document.getElementById('modal-container');
    const titleEl = document.getElementById('modal-title');
    const bodyEl = document.getElementById('modal-body');
    const footerEl = document.getElementById('modal-footer');
    const closeBtn = document.getElementById('modal-close-btn');

    if (!overlay) return;

    titleEl.textContent = title || '';
    bodyEl.innerHTML = bodyHTML || '';
    footerEl.innerHTML = footerHTML || '';

    overlay.classList.remove('hidden');

    const handleClose = () => {
        overlay.classList.add('hidden');
        if (onClose) onClose();
    };

    closeBtn.onclick = handleClose;
    overlay.onclick = (e) => {
        if (e.target === overlay) handleClose();
    };
}

export function closeModal() {
    const overlay = document.getElementById('modal-container');
    if (overlay) {
        overlay.classList.add('hidden');
    }
}
