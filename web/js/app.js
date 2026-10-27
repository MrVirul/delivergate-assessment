function closeModal() {
  const modal = document.getElementById('modal');
  if (modal) modal.close();
}

function openItemEdit(id) {
  const modal = document.getElementById('modal');
  if (modal) {
    document.getElementById('modal-title').textContent = 'Edit item';
    modal.showModal();
  }
}

function openCustomerEdit(id) {
  const modal = document.getElementById('modal');
  if (modal) {
    document.getElementById('modal-title').textContent = 'Edit customer';
    modal.showModal();
  }
}

function openItemAdd() {
  const modal = document.getElementById('modal');
  if (modal) {
    document.getElementById('modal-title').textContent = 'Add item';
    document.getElementById('modal-body').innerHTML = UI.render ? UI.render('item-form', null) : RENDERERS['item-form'](null);
    modal.showModal();
  }
}

function openCustomerAdd() {
  const modal = document.getElementById('modal');
  if (modal) {
    document.getElementById('modal-title').textContent = 'Add customer';
    document.getElementById('modal-body').innerHTML = UI.render ? UI.render('customer-form', null) : RENDERERS['customer-form'](null);
    modal.showModal();
  }
}

if (typeof UI !== 'undefined' && !UI.render) {
  UI.render = function(key, body) {
    if (RENDERERS && RENDERERS[key]) {
      return RENDERERS[key](body);
    }
    return '';
  };
}

function initApp() {
  const user = UI.user();
  const token = UI.token();
  if (!user || !token) {
    window.location.replace('/login.html');
    return;
  }
  const nameEl = document.getElementById('user-name');
  const roleEl = document.getElementById('user-role');
  if (nameEl) nameEl.textContent = user.name || user.email || '';
  if (roleEl) roleEl.textContent = user.role || '';

  // Show/hide add buttons
  const addItemBtn = document.getElementById('add-item-btn');
  if (addItemBtn && !UI.can('items:create')) {
    addItemBtn.style.display = 'none';
  }
  const addCustBtn = document.getElementById('add-customer-btn');
  if (addCustBtn && !UI.can('customers:create')) {
    addCustBtn.style.display = 'none';
  }

  // Tab switching
  document.querySelectorAll('.tab').forEach(tab => {
    tab.addEventListener('click', () => {
      document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
      tab.classList.add('active');
      const target = tab.dataset.target;
      document.querySelectorAll('.tab-panel').forEach(p => p.style.display = 'none');
      const panel = document.getElementById(target);
      if (panel) panel.style.display = 'block';
      if (target === 'customers-panel') {
        htmx.trigger('#customers-results', 'load');
      }
    });
  });
}

document.addEventListener('DOMContentLoaded', () => {
  if (window.location.pathname.includes('login')) {
    return;
  }
  initApp();
});

// Global htmx listeners
document.body.addEventListener('htmx:afterSwap', (evt) => {
  const target = evt.detail.target;
  if (target && target.id === 'modal-body') {
    const modal = document.getElementById('modal');
    if (modal) {
      document.getElementById('modal-title').textContent = target.querySelector('h3')?.textContent || 'Form';
      modal.showModal();
      const first = target.querySelector('input, select, textarea');
      if (first) first.focus();
    }
  }
});

document.body.addEventListener('htmx:afterRequest', (evt) => {
  const path = evt.detail.pathInfo?.requestPath || evt.detail.requestConfig?.path || '';
  if (evt.detail.successful) {
    // check if it's a form submit
    if (path.includes('/api/items') && (evt.detail.requestConfig?.verb === 'post' || evt.detail.requestConfig?.verb === 'put')) {
      closeModal();
      UI.toast('Item saved', 'success');
      htmx.trigger(document.body, 'items-changed');
    }
    if (path.includes('/api/customers') && (evt.detail.requestConfig?.verb === 'post' || evt.detail.requestConfig?.verb === 'put')) {
      closeModal();
      UI.toast('Customer saved', 'success');
      htmx.trigger(document.body, 'customers-changed');
    }
  } else {
    // handle errors
    if (evt.detail.xhr) {
      try {
        const res = JSON.parse(evt.detail.xhr.responseText);
        if (evt.detail.xhr.status === 401) {
          sessionStorage.removeItem('token');
          sessionStorage.removeItem('user');
          sessionStorage.removeItem('expires_at');
          UI.toast('Your session has ended. Please sign in again.', 'error');
          setTimeout(() => window.location.replace('/login.html'), 500);
          return;
        }
        if (evt.detail.xhr.status === 403) {
          UI.toast("You don't have permission to do that.", 'error');
          return;
        }
        if (evt.detail.xhr.status === 422 && res.error?.details) {
          const details = res.error.details;
          // try to show near fields
          details.forEach(d => {
            const el = document.querySelector('[name="' + d.field.replace(/"/g,'\\"') + '"]');
            if (el) {
              el.setAttribute('aria-invalid', 'true');
              let msg = el.parentElement.querySelector('.field-error');
              if (!msg) {
                msg = document.createElement('p');
                msg.className = 'field-error';
                el.parentElement.appendChild(msg);
              }
              msg.textContent = d.message;
            }
          });
          return;
        }
        if (res.error?.message) {
          UI.toast(res.error.message, 'error');
        }
      } catch (e) {
        UI.toast('Request failed', 'error');
      }
    }
  }
});

document.body.addEventListener('htmx:responseError', (evt) => {
  if (evt.detail.xhr?.status === 401) {
    sessionStorage.removeItem('token');
    sessionStorage.removeItem('user');
    sessionStorage.removeItem('expires_at');
    UI.toast('Your session has ended. Please sign in again.', 'error');
    setTimeout(() => window.location.replace('/login.html'), 500);
    return;
  }
  if (evt.detail.xhr?.status === 403) {
    UI.toast("You don't have permission to do that.", 'error');
    return;
  }
  UI.toast('Something went wrong', 'error');
});

document.body.addEventListener('htmx:sendError', () => {
  UI.toast('Cannot reach the server', 'error');
});

// Logout
document.addEventListener('click', (e) => {
  if (e.target && e.target.id === 'logout-btn') {
    const token = UI.token();
    if (token) {
      fetch('/api/auth/logout', {
        method: 'POST',
        headers: { 'Authorization': 'Bearer ' + token }
      }).finally(() => {
        sessionStorage.removeItem('token');
        sessionStorage.removeItem('user');
        sessionStorage.removeItem('expires_at');
        window.location.replace('/login.html');
      });
    } else {
      sessionStorage.removeItem('token');
      sessionStorage.removeItem('user');
      sessionStorage.removeItem('expires_at');
      window.location.replace('/login.html');
    }
  }
});
