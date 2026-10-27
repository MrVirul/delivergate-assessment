try {
  const u = JSON.parse(sessionStorage.getItem('user')||'null');
  const t = sessionStorage.getItem('token');
  if (u && t) { window.location.replace('/'); }
} catch(e) {}

document.body.addEventListener('htmx:afterRequest', (evt) => {
  const path = evt.detail.pathInfo?.requestPath || evt.detail.requestConfig?.path || '';
  if (path.includes('/api/auth/login')) {
    const errEl = document.getElementById('login-error');
    if (evt.detail.successful && evt.detail.xhr) {
      try {
        const res = JSON.parse(evt.detail.xhr.responseText);
        if (res.data && res.data.access_token) {
          sessionStorage.setItem('token', res.data.access_token);
          sessionStorage.setItem('user', JSON.stringify(res.data.user));
          sessionStorage.setItem('expires_at', String(Date.now() + (res.data.expires_in || 3600)*1000));
          window.location.replace('/');
          return;
        }
      } catch(e) {
        if (errEl) errEl.textContent = 'Login failed';
        return;
      }
    }
    if (errEl) errEl.textContent = 'Invalid email or password';
  }
});
