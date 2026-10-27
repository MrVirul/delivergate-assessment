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
    const body = document.getElementById('modal-body');
    body.innerHTML = UI.render ? UI.render('item-form', null) : RENDERERS['item-form'](null);
    if (typeof htmx !== 'undefined' && htmx.process) {
      htmx.process(body);
    }
    modal.showModal();
  }
}

function openCustomerAdd() {
  const modal = document.getElementById('modal');
  if (modal) {
    document.getElementById('modal-title').textContent = 'Add customer';
    const body = document.getElementById('modal-body');
    body.innerHTML = UI.render ? UI.render('customer-form', null) : RENDERERS['customer-form'](null);
    if (typeof htmx !== 'undefined' && htmx.process) {
      htmx.process(body);
    }
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
  const expires = sessionStorage.getItem('expires_at');
  if (!user || !token) {
    window.location.replace('/login.html');
    return;
  }
  if (expires) {
    try {
      if (Date.now() > parseInt(expires)) {
        sessionStorage.removeItem('token');
        sessionStorage.removeItem('user');
        sessionStorage.removeItem('expires_at');
        window.location.replace('/login.html');
        return;
      }
    } catch(e) {}
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
      try { modal.showModal(); } catch (e) {}
      const first = target.querySelector('input, select, textarea');
      if (first) try { first.focus(); } catch (e) {}
      // snapshot forms for edit change detection
      try {
        const form = target.querySelector('form[data-edit-id]');
        if (form && !form._confirmSnapshot) {
          const data = {};
          const fd = new FormData(form);
          for (const [k,v] of fd.entries()) { data[k] = v; }
          form._confirmSnapshot = data;
        }
      } catch (e) {}
    }
  }
});

document.body.addEventListener('click', (e) => {
  const delItemBtn = e.target.closest('[data-delete-item]');
  if (delItemBtn) {
    e.preventDefault();
    const id = delItemBtn.getAttribute('data-delete-item');
    const name = delItemBtn.getAttribute('data-delete-name') || 'this item';
    confirmAction('Delete item', 'Delete ' + name + '? This cannot be undone.', () => {
      const token = UI.token ? UI.token() : sessionStorage.getItem('token');
      fetch('/api/items/' + id, {
        method: 'DELETE',
        headers: token ? { 'Authorization': 'Bearer ' + token } : {}
      }).then((resp) => {
        if (resp.status === 401) {
          sessionStorage.removeItem('token');
          sessionStorage.removeItem('user');
          sessionStorage.removeItem('expires_at');
          UI.toast && UI.toast('Your session has ended. Please sign in again.', 'error');
          setTimeout(() => window.location.replace('/login.html'), 500);
          return;
        }
        if (resp.status === 204 || resp.ok) {
          UI.toast && UI.toast('Item deleted', 'success');
          const tr = delItemBtn.closest('tr');
          if (tr) tr.remove();
          htmx.trigger(document.body, 'items-changed');
        } else {
          UI.toast && UI.toast('Failed to delete item', 'error');
        }
      }).catch(() => {
        UI.toast && UI.toast('Failed to delete item', 'error');
      });
    });
    return;
  }
  const delCustBtn = e.target.closest('[data-delete-customer]');
  if (delCustBtn) {
    e.preventDefault();
    const id = delCustBtn.getAttribute('data-delete-customer');
    const name = delCustBtn.getAttribute('data-delete-name') || 'this customer';
    confirmAction('Delete customer', 'Delete ' + name + '? This cannot be undone.', () => {
      const token = UI.token ? UI.token() : sessionStorage.getItem('token');
      fetch('/api/customers/' + id, {
        method: 'DELETE',
        headers: token ? { 'Authorization': 'Bearer ' + token } : {}
      }).then((resp) => {
        if (resp.status === 401) {
          sessionStorage.removeItem('token');
          sessionStorage.removeItem('user');
          sessionStorage.removeItem('expires_at');
          UI.toast && UI.toast('Your session has ended. Please sign in again.', 'error');
          setTimeout(() => window.location.replace('/login.html'), 500);
          return;
        }
        if (resp.status === 204 || resp.ok) {
          UI.toast && UI.toast('Customer deleted', 'success');
          const tr = delCustBtn.closest('tr');
          if (tr) tr.remove();
          htmx.trigger(document.body, 'customers-changed');
        } else {
          UI.toast && UI.toast('Failed to delete customer', 'error');
        }
      }).catch(() => {
        UI.toast && UI.toast('Failed to delete customer', 'error');
      });
    });
    return;
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
function confirmAction(title, message, onConfirm) {
  const modal = document.getElementById('confirm-modal');
  const titleEl = document.getElementById('confirm-title');
  const msgEl = document.getElementById('confirm-message');
  const okBtn = document.getElementById('confirm-ok');
  const cancelBtn = document.getElementById('confirm-cancel');
  const cancelX = document.getElementById('confirm-cancel-x');
  if (!modal) {
    if (confirm(message)) onConfirm && onConfirm();
    return;
  }
  if (titleEl) titleEl.textContent = title || 'Confirm';
  if (msgEl) msgEl.textContent = message || 'Are you sure?';
  const cleanup = () => {
    okBtn.onclick = null;
    cancelBtn.onclick = null;
    if (cancelX) cancelX.onclick = null;
    modal.onclose = null;
  };
  okBtn.onclick = () => {
    cleanup();
    modal.close();
    try { onConfirm && onConfirm(); } catch (e) {}
  };
  cancelBtn.onclick = () => {
    cleanup();
    modal.close();
  };
  if (cancelX) cancelX.onclick = () => {
    cleanup();
    modal.close();
  };
  modal.onclose = cleanup;
  try { modal.showModal(); } catch (e) { if (confirm(message)) onConfirm && onConfirm(); }
}

function performLogout() {
  const token = UI.token ? UI.token() : sessionStorage.getItem('token');
  const doRedirect = () => {
    try {
      sessionStorage.removeItem('token');
      sessionStorage.removeItem('user');
      sessionStorage.removeItem('expires_at');
    } catch (e) {}
    window.location.replace('/login.html');
  };
  if (token) {
    try {
      fetch('/api/auth/logout', {
        method: 'POST',
        headers: { 'Authorization': 'Bearer ' + token },
        keepalive: true
      }).then(doRedirect).catch(doRedirect);
      setTimeout(doRedirect, 1500);
      return;
    } catch (err) {
      doRedirect();
      return;
    }
  }
  doRedirect();
}

document.addEventListener('click', (e) => {
  const btn = e.target.closest ? e.target.closest('#logout-btn') : (e.target && e.target.id === 'logout-btn' ? e.target : null);
  if (btn) {
    e.preventDefault();
    confirmAction('Sign out', 'Are you sure you want to sign out?', performLogout);
    return;
  }
});

document.addEventListener('submit', (e) => {
  const form = e.target;
  if (form && form.id === 'logout-form') {
    e.preventDefault();
    performLogout();
  }
});

document.body.addEventListener('htmx:afterRequest', (evt) => {
  const path = evt.detail.pathInfo?.requestPath || evt.detail.requestConfig?.path || '';
  const verb = (evt.detail.requestConfig?.verb || '').toLowerCase();
  if (evt.detail.successful) {
    if (verb === 'delete') {
      if (path.includes('/api/items')) {
        try { UI.toast && UI.toast('Item deleted', 'success'); } catch(e) {}
        // if row removed and might be last on page, reload to adjust pagination
        setTimeout(function(){ htmx.trigger(document.body, 'items-changed'); }, 50);
      } else if (path.includes('/api/customers')) {
        try { UI.toast && UI.toast('Customer deleted', 'success'); } catch(e) {}
        setTimeout(function(){ htmx.trigger(document.body, 'customers-changed'); }, 50);
      }
    }
  }
}, true);
