/**
 * ShambaChain — Shared Application Utilities
 * Handles: API calls, auth state, toast notifications, loading states, navigation
 */

// API BASE
const API_BASE = '/api';

//AUTH HELPERS
const Auth = {
  getToken() {
    return localStorage.getItem('sc_token');
  },
  setToken(token) {
    localStorage.setItem('sc_token', token);
  },
  setUser(user) {
    localStorage.setItem('sc_user', JSON.stringify(user));
  },
  getUser() {
    try {
      return JSON.parse(localStorage.getItem('sc_user') || 'null');
    } catch {
      return null;
    }
  },
  isLoggedIn() {
    return !!this.getToken();
  },
  logout() {
    localStorage.removeItem('sc_token');
    localStorage.removeItem('sc_user');
    localStorage.removeItem('selectedRole');
    localStorage.removeItem('cachedUsername');
    window.location.href = 'login.html';
  },
};

// HTTP CLIENT
const Http = {
  async request(method, path, body = null, requiresAuth = false) {
    const headers = { 'Content-Type': 'application/json' };
    if (requiresAuth) {
      const token = Auth.getToken();
      if (!token) {
        Auth.logout();
        return;
      }
      headers['Authorization'] = `Bearer ${token}`;
    }

    const options = { method, headers };
    if (body) options.body = JSON.stringify(body);

    const res = await fetch(`${API_BASE}${path}`, options);
    const data = await res.json().catch(() => ({}));

    if (!res.ok) {
      throw new Error(data.error || data.message || `Request failed (${res.status})`);
    }
    return data;
  },

  get(path, requiresAuth = false) {
    return this.request('GET', path, null, requiresAuth);
  },
  post(path, body, requiresAuth = false) {
    return this.request('POST', path, body, requiresAuth);
  },
};

//TOAST NOTIFICATIONS
const Toast = {
  _container: null,

  _getContainer() {
    if (!this._container) {
      this._container = document.createElement('div');
      this._container.id = 'sc-toast-container';
      this._container.style.cssText = `
        position: fixed; top: 20px; right: 20px; z-index: 9999;
        display: flex; flex-direction: column; gap: 10px;
        pointer-events: none;
      `;
      document.body.appendChild(this._container);
    }
    return this._container;
  },

  show(message, type = 'info', duration = 3500) {
    const colors = {
      success: { bg: '#022c22', border: '#22c55e', icon: '✓' },
      error:   { bg: '#7f1d1d', border: '#ef4444', icon: '✕' },
      info:    { bg: '#1e3a5f', border: '#3b82f6', icon: 'ℹ' },
      warning: { bg: '#78350f', border: '#f59e0b', icon: '⚠' },
    };
    const c = colors[type] || colors.info;

    const toast = document.createElement('div');
    toast.style.cssText = `
      background: ${c.bg}; border-left: 4px solid ${c.border};
      color: #fff; padding: 12px 16px; border-radius: 8px;
      font-size: 14px; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      box-shadow: 0 4px 20px rgba(0,0,0,0.3); pointer-events: all;
      display: flex; align-items: center; gap: 10px; min-width: 260px; max-width: 360px;
      opacity: 0; transform: translateX(20px);
      transition: opacity 0.25s ease, transform 0.25s ease;
    `;
    toast.innerHTML = `<span style="font-weight:700;font-size:16px;">${c.icon}</span><span>${message}</span>`;

    this._getContainer().appendChild(toast);

    // Animate in
    requestAnimationFrame(() => {
      toast.style.opacity = '1';
      toast.style.transform = 'translateX(0)';
    });

    // Animate out and remove
    setTimeout(() => {
      toast.style.opacity = '0';
      toast.style.transform = 'translateX(20px)';
      setTimeout(() => toast.remove(), 300);
    }, duration);
  },

  success(msg) { this.show(msg, 'success'); },
  error(msg)   { this.show(msg, 'error'); },
  info(msg)    { this.show(msg, 'info'); },
  warning(msg) { this.show(msg, 'warning'); },
};

// BUTTON LOADING STATE 
const Btn = {
  setLoading(btn, loading, originalText = null) {
    if (loading) {
      btn.dataset.originalText = btn.innerHTML;
      btn.innerHTML = `<span style="display:inline-flex;align-items:center;gap:8px;">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"
          style="animation:sc-spin 0.8s linear infinite;">
          <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/>
        </svg>
        Loading...
      </span>`;
      btn.disabled = true;
      btn.style.opacity = '0.75';
    } else {
      btn.innerHTML = originalText || btn.dataset.originalText || btn.innerHTML;
      btn.disabled = false;
      btn.style.opacity = '1';
    }
  },
};

// Inject spin keyframe once
(function injectSpinStyle() {
  if (document.getElementById('sc-spin-style')) return;
  const style = document.createElement('style');
  style.id = 'sc-spin-style';
  style.textContent = `@keyframes sc-spin { to { transform: rotate(360deg); } }`;
  document.head.appendChild(style);
})();

//  FORM HELPERS
const Form = {
  // Show inline error under a field
  setError(inputEl, message) {
    this.clearError(inputEl);
    inputEl.style.borderColor = '#ef4444';
    const err = document.createElement('p');
    err.className = 'sc-field-error';
    err.style.cssText = 'color:#ef4444;font-size:12px;margin-top:4px;';
    err.textContent = message;
    inputEl.parentNode.insertBefore(err, inputEl.nextSibling);
  },
  clearError(inputEl) {
    inputEl.style.borderColor = '';
    const existing = inputEl.parentNode.querySelector('.sc-field-error');
    if (existing) existing.remove();
  },
  clearAllErrors(formEl) {
    formEl.querySelectorAll('.sc-field-error').forEach(e => e.remove());
    formEl.querySelectorAll('input, select, textarea').forEach(el => {
      el.style.borderColor = '';
    });
  },
};

// NAVIGATION GUARD
// Call on protected pages to redirect unauthenticated users
function requireAuth() {
  if (!Auth.isLoggedIn()) {
    Toast.warning('Please log in to continue.');
    setTimeout(() => { window.location.href = 'login.html'; }, 800);
    return false;
  }
  return true;
}

//  SIGN OUT WIRING 
// Attach to any element with data-action="signout"
document.addEventListener('DOMContentLoaded', () => {
  document.querySelectorAll('[data-action="signout"]').forEach(el => {
    el.addEventListener('click', (e) => {
      e.preventDefault();
      Auth.logout();
    });
  });
});
