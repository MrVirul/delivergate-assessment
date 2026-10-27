const CURRENCY = 'LKR';
const SHOP_NAME = 'Scoop & Co.';

const PERMISSIONS = {
  Owner: ['items:read','items:create','items:update','items:delete','customers:read','customers:create','customers:update','customers:delete'],
  Manager: ['items:read','items:create','items:update','customers:read','customers:create','customers:update'],
  Cashier: ['items:read','items:update','customers:read','customers:create','customers:update'],
};

const UI = {
  esc: (v) => String(v ?? '').replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c])),
  money: (cents) => {
    const val = (typeof cents === 'number' ? cents : 0) / 100;
    try {
      return new Intl.NumberFormat(undefined, { style: 'currency', currency: CURRENCY }).format(val);
    } catch (e) {
      return CURRENCY + ' ' + val.toFixed(2);
    }
  },
  user: () => {
    try {
      return JSON.parse(sessionStorage.getItem('user') || 'null');
    } catch (e) {
      return null;
    }
  },
  token: () => sessionStorage.getItem('token'),
  can: (perm) => {
    const u = UI.user();
    const perms = PERMISSIONS[u?.role] || [];
    return perms.includes(perm);
  },
  coerceForm: (formData, form) => {
    const obj = {};
    for (const [key, value] of formData.entries()) {
      if (value === '' || value === null || value === undefined) {
        continue;
      }
      const input = form?.querySelector('[name="' + key.replace(/"/g, '\\"') + '"]');
      if (input && input.type === 'number') {
        if (value === '') continue;
        const num = Number(value);
        if (!isNaN(num)) obj[key] = num;
        continue;
      }
      if (input && input.dataset.money === 'true') {
        if (value === '') continue;
        const num = parseFloat(value);
        if (!isNaN(num)) {
          obj[key] = Math.round(num * 100);
        }
        continue;
      }
      if (typeof value === 'string') {
        obj[key] = value.trim();
      } else {
        obj[key] = value;
      }
    }
    return obj;
  },
  toast: (message, kind) => {
    const container = document.getElementById('toast-container');
    if (!container) return;
    const toast = document.createElement('div');
    toast.className = 'toast' + (kind ? ' ' + kind : '');
    toast.setAttribute('role', 'status');
    toast.setAttribute('aria-live', 'polite');
    toast.textContent = message;
    const btn = document.createElement('button');
    btn.className = 'btn btn-ghost';
    btn.textContent = '✕';
    btn.style.marginLeft = '0.75rem';
    btn.onclick = () => toast.remove();
    toast.appendChild(btn);
    container.appendChild(toast);
    const timeout = kind === 'error' ? 6000 : 4000;
    setTimeout(() => {
      if (toast.parentElement) toast.remove();
    }, timeout);
  }
};

function gaugeHtml(qty, reorder) {
  const q = typeof qty === 'number' ? qty : 0;
  const r = typeof reorder === 'number' ? reorder : 0;
  const low = q <= r;
  const pct = Math.max(5, Math.min(100, r === 0 ? 50 : (q / (r * 3)) * 100));
  return '<div class="gauge"><div class="gauge-bar"><div class="gauge-fill' + (low ? ' low' : '') + '" style="width:' + pct + '%"></div></div><span class="gauge-text' + (low ? ' low' : '') + '">' + UI.esc(String(q)) + (low ? ' · Low stock' : '') + '</span></div>';
}

const RENDERERS = {
  'items-list': (body) => {
    const items = body.data || [];
    const meta = body.meta || {};
    let html = '<div class="panel">';
    html += '<div class="table-wrap"><table><thead><tr><th>Name</th><th>Category</th><th>Stock</th><th>Price</th><th>Actions</th></tr></thead><tbody>';
    if (items.length === 0) {
      html += '<tr><td colspan="5">No items match your search. Clear filters.</td></tr>';
    } else {
      for (const it of items) {
        html += '<tr>';
        html += '<td data-label="Name">' + UI.esc(it.name) + '</td>';
        html += '<td data-label="Category">' + UI.esc(it.category) + '</td>';
        html += '<td data-label="Stock">' + gaugeHtml(it.quantity, it.reorder_level) + '</td>';
        html += '<td data-label="Price">' + UI.money(it.price_cents) + '</td>';
        html += '<td class="actions" data-label="Actions">';
        if (UI.can('items:update')) {
          html += '<button class="btn" hx-get="/api/items/' + it.id + '" hx-target="#modal-body" hx-swap="innerHTML" data-render="item-form" onclick="openItemEdit(' + it.id + ')">Edit</button>';
        }
        if (UI.can('items:delete')) {
          html += '<button class="btn btn-danger" hx-delete="/api/items/' + it.id + '" hx-target="closest tr" hx-swap="outerHTML" hx-confirm="Delete ' + UI.esc(it.name).replace(/"/g, '&quot;') + '? This cannot be undone.">Delete</button>';
        }
        html += '</td></tr>';
      }
    }
    html += '</tbody></table></div>';
    html += '<div class="pagination"><div>Showing page ' + (meta.page || 1) + ' · Total ' + (meta.total || 0) + '</div><div style="display:flex;gap:0.5rem;">';
    if ((meta.page || 1) > 1) {
      html += '<button class="btn" hx-get="/api/items" hx-include="#items-filters" hx-vals=\'{"page":' + ((meta.page||1)-1) + '}\' hx-target="#items-results" data-render="items-list">Previous</button>';
    }
    if ((meta.page || 1) < (meta.total_pages || 1)) {
      html += '<button class="btn" hx-get="/api/items" hx-include="#items-filters" hx-vals=\'{"page":' + ((meta.page||1)+1) + '}\' hx-target="#items-results" data-render="items-list">Next</button>';
    }
    html += '</div></div></div>';
    return html;
  },
  'item-form': (body) => {
    const it = body && body.data ? body.data : null;
    const isEdit = !!it;
    let html = '<form ' + (isEdit ? 'hx-put="/api/items/' + it.id + '"' : 'hx-post="/api/items"') + ' hx-swap="none" hx-disabled-elt="find button[type=submit]"><h3>' + (isEdit ? 'Edit item' : 'Add item') + '</h3>';
    html += '<div class="form-group"><label>Name</label><input name="name" required value="' + UI.esc(isEdit ? it.name : '') + '"></div>';
    html += '<div class="form-group"><label>Category</label><select name="category">' +
      ['flavor','topping','cone','cup','packaging','other'].map(c => '<option value="' + c + '"' + (isEdit && it.category === c ? ' selected' : '') + '>' + c + '</option>').join('') +
      '</select></div>';
    html += '<div class="form-group"><label>Description</label><textarea name="description">' + UI.esc(isEdit ? (it.description||'') : '') + '</textarea></div>';
    html += '<div class="form-group"><label>Quantity</label><input type="number" name="quantity" min="0" value="' + (isEdit ? it.quantity : 0) + '"></div>';
    html += '<div class="form-group"><label>Unit</label><select name="unit">' +
      ['tub','scoop','kg','litre','pcs','box'].map(c => '<option value="' + c + '"' + (isEdit && it.unit === c ? ' selected' : '') + '>' + c + '</option>').join('') +
      '</select></div>';
    html += '<div class="form-group"><label>Price</label><input data-money="true" name="price_cents" value="' + (isEdit ? (it.price_cents/100).toFixed(2) : '0.00') + '"></div>';
    html += '<div class="form-group"><label>Reorder level</label><input type="number" name="reorder_level" min="0" value="' + (isEdit ? it.reorder_level : 0) + '"></div>';
    html += '<div style="display:flex;gap:0.5rem;justify-content:flex-end;"><button type="button" class="btn" onclick="closeModal()">Cancel</button><button type="submit" class="btn btn-primary">Save item</button></div>';
    html += '<div id="item-form-errors"></div></form>';
    return html;
  },
  'customers-list': (body) => {
    const items = body.data || [];
    const meta = body.meta || {};
    let html = '<div class="panel">';
    html += '<div class="table-wrap"><table><thead><tr><th>Name</th><th>Email</th><th>Phone</th><th>Actions</th></tr></thead><tbody>';
    if (items.length === 0) {
      html += '<tr><td colspan="4">No customers yet. Add the first one.</td></tr>';
    } else {
      for (const it of items) {
        html += '<tr>';
        html += '<td data-label="Name">' + UI.esc(it.name) + '</td>';
        html += '<td data-label="Email">' + UI.esc(it.email || '') + '</td>';
        html += '<td data-label="Phone">' + UI.esc(it.phone || '') + '</td>';
        html += '<td class="actions" data-label="Actions">';
        if (UI.can('customers:update')) {
          html += '<button class="btn" hx-get="/api/customers/' + it.id + '" hx-target="#modal-body" hx-swap="innerHTML" data-render="customer-form" onclick="openCustomerEdit(' + it.id + ')">Edit</button>';
        }
        if (UI.can('customers:delete')) {
          html += '<button class="btn btn-danger" hx-delete="/api/customers/' + it.id + '" hx-target="closest tr" hx-swap="outerHTML" hx-confirm="Delete ' + UI.esc(it.name).replace(/"/g, '&quot;') + '? This cannot be undone.">Delete</button>';
        }
        html += '</td></tr>';
      }
    }
    html += '</tbody></table></div>';
    html += '<div class="pagination"><div>Showing page ' + (meta.page || 1) + ' · Total ' + (meta.total || 0) + '</div><div style="display:flex;gap:0.5rem;">';
    if ((meta.page || 1) > 1) {
      html += '<button class="btn" hx-get="/api/customers" hx-include="#customers-filters" hx-vals=\'{"page":' + ((meta.page||1)-1) + '}\' hx-target="#customers-results" data-render="customers-list">Previous</button>';
    }
    if ((meta.page || 1) < (meta.total_pages || 1)) {
      html += '<button class="btn" hx-get="/api/customers" hx-include="#customers-filters" hx-vals=\'{"page":' + ((meta.page||1)+1) + '}\' hx-target="#customers-results" data-render="customers-list">Next</button>';
    }
    html += '</div></div></div>';
    return html;
  },
  'customer-form': (body) => {
    const it = body && body.data ? body.data : null;
    const isEdit = !!it;
    let html = '<form ' + (isEdit ? 'hx-put="/api/customers/' + it.id + '"' : 'hx-post="/api/customers"') + ' hx-swap="none" hx-disabled-elt="find button[type=submit]"><h3>' + (isEdit ? 'Edit customer' : 'Add customer') + '</h3>';
    html += '<div class="form-group"><label>Name</label><input name="name" required value="' + UI.esc(isEdit ? it.name : '') + '"></div>';
    html += '<div class="form-group"><label>Email</label><input name="email" value="' + UI.esc(isEdit && it.email ? it.email : '') + '"></div>';
    html += '<div class="form-group"><label>Phone</label><input name="phone" required value="' + UI.esc(isEdit ? it.phone : '') + '"></div>';
    html += '<div class="form-group"><label>Address</label><input name="address" value="' + UI.esc(isEdit ? (it.address||'') : '') + '"></div>';
    html += '<div class="form-group"><label>Notes</label><textarea name="notes">' + UI.esc(isEdit ? (it.notes||'') : '') + '</textarea></div>';
    html += '<div style="display:flex;gap:0.5rem;justify-content:flex-end;"><button type="button" class="btn" onclick="closeModal()">Cancel</button><button type="submit" class="btn btn-primary">Save customer</button></div>';
    html += '<div id="customer-form-errors"></div></form>';
    return html;
  },
  'login': (body) => {
    return '';
  }
};
