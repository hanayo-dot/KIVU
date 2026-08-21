document.addEventListener('DOMContentLoaded', () => {
    // 1. Handle Login Form Submission
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', (e) => {
            e.preventDefault();
            localStorage.setItem('auth', 'true');
            window.location.href = 'index.html';
        });
    }

    // 2. Enforce Session Security Guard across pages
    if (!window.location.pathname.includes('login.html')) {
        if (localStorage.getItem('auth') !== 'true') {
            window.location.href = 'login.html';
        }
    }

    // 3. Handle Observation Form Submission & Toast Animation
    const obsForm = document.getElementById('observation-form');
    if (obsForm) {
        obsForm.addEventListener('submit', (e) => {
            e.preventDefault();
            const toast = document.getElementById('success-toast');
            if (toast) {
                toast.classList.add('show');
                setTimeout(() => {
                    toast.classList.remove('show');
                }, 4000);
            }
        });
    }
});

// 4. Category Selector Toggle for Reports Page
function categoryToggle(element) {
    document.querySelectorAll('.cat-btn').forEach(btn => {
        btn.classList.remove('active');
        btn.style.background = 'var(--bg-color)';
        btn.style.borderColor = 'var(--card-border)';
        btn.style.color = 'var(--text-muted)';
    });
    element.classList.add('active');
    element.style.background = 'rgba(239, 68, 68, 0.1)';
    element.style.borderColor = 'var(--accent-red)';
    element.style.color = 'var(--text-main)';
}