// Custom confirmation card using htmx:confirm
(function() {
  const esc = function(v) {
    return String(v == null ? '' : v).replace(/[&<>"']/g, function(c) {
      return ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'})[c];
    });
  };

  function getCard() {
    let card = document.getElementById('confirm-card');
    if (card) return card;
    card = document.createElement('dialog');
    card.id = 'confirm-card';
    card.className = 'confirm-dialog';
    card.setAttribute('role', 'alertdialog');
    card.innerHTML = '<div class="confirm-card">' +
      '<div class="confirm-header"><strong id="confirm-card-title"></strong></div>' +
      '<div class="confirm-body">' +
      '<p id="confirm-card-record" class="confirm-record"></p>' +
      '<p id="confirm-card-sentence" class="confirm-sentence"></p>' +
      '<div id="confirm-card-changes" class="confirm-changes"></div>' +
      '</div>' +
      '<div class="confirm-actions">' +
      '<button type="button" class="btn" id="confirm-card-safe">Go back</button>' +
      '<button type="button" class="btn btn-danger" id="confirm-card-ok">Confirm</button>' +
      '</div>' +
      '</div>';
    document.body.appendChild(card);
    return card;
  }

  function closeCard(card) {
    if (card && card.open) {
      try { card.close(); } catch(e) {}
    }
  }

  function getFormDataSnapshot(form) {
    const data = {};
    try {
      const fd = new FormData(form);
      for (const [k,v] of fd.entries()) {
        data[k] = v;
      }
    } catch(e) {}
    return data;
  }

  document.body.addEventListener('click', function(e) {
    const safe = e.target.closest('#confirm-card-safe');
    if (safe) {
      const card = document.getElementById('confirm-card');
      closeCard(card);
      return;
    }
  });

  document.body.addEventListener('keydown', function(e) {
    if (e.key === 'Enter') {
      const card = document.getElementById('confirm-card');
      if (card && card.open && card.dataset.danger === 'true') {
        e.preventDefault();
        e.stopPropagation();
      }
    }
  }, true);

  document.body.addEventListener('htmx:confirm', function(evt) {
    const question = evt.detail.question || '';
    if (!question) {
      return;
    }
    evt.preventDefault();

    const elt = evt.detail.elt;
    if (!elt) { evt.detail.issueRequest(true); return; }

    // Check for edit forms - detect changes
    if (elt.tagName === 'FORM' && elt.hasAttribute('data-edit-id')) {
      // store snapshot if not present
      if (!elt._confirmSnapshot) {
        elt._confirmSnapshot = getFormDataSnapshot(elt);
      }
      const current = getFormDataSnapshot(elt);
      const changes = [];
      const snap = elt._confirmSnapshot || {};
      for (const k in current) {
        if (Object.prototype.hasOwnProperty.call(snap, k)) {
          if (String(snap[k]) !== String(current[k])) {
            // format nicely for common fields
            if (k === 'price_cents') {
              const oldv = (parseFloat(snap[k]) || 0);
              const newv = (parseFloat(current[k]) || 0);
              changes.push('Price: ' + oldv.toFixed(2) + ' -> ' + newv.toFixed(2));
            } else {
              changes.push(k + ': ' + snap[k] + ' -> ' + current[k]);
            }
          }
        } else {
          changes.push(k + ': (empty) -> ' + current[k]);
        }
      }
      if (changes.length === 0) {
        // No changes - show toast and don't send request
        try {
          if (window.UI && UI.toast) UI.toast('No changes to save.', 'info');
        } catch(e) {}
        return;
      }
      elt.setAttribute('data-confirm-changes', JSON.stringify(changes));
    }

    const titleAttr = elt.getAttribute ? elt.getAttribute('data-confirm-title') : '';
    const actionAttr = elt.getAttribute ? elt.getAttribute('data-confirm-action') : '';
    const kind = elt.getAttribute ? elt.getAttribute('data-confirm-kind') : 'danger';

    const card = getCard();
    card.dataset.kind = kind;
    card.dataset.danger = kind === 'danger' ? 'true' : 'false';

    const titleEl = document.getElementById('confirm-card-title');
    const recordEl = document.getElementById('confirm-card-record');
    const sentenceEl = document.getElementById('confirm-card-sentence');
    const changesEl = document.getElementById('confirm-card-changes');
    const okBtn = document.getElementById('confirm-card-ok');
    const safeBtn = document.getElementById('confirm-card-safe');

    if (titleEl) titleEl.textContent = titleAttr || 'Confirm';
    if (recordEl) {
      const rec = elt.getAttribute ? elt.getAttribute('data-confirm-record') : '';
      if (rec) recordEl.innerHTML = '<strong>' + esc(rec) + '</strong>';
      else recordEl.innerHTML = '';
    }
    if (sentenceEl) sentenceEl.textContent = actionAttr || (kind === 'danger' ? 'This cannot be undone.' : '');

    if (changesEl) {
      changesEl.innerHTML = '';
      try {
        const changes = elt.getAttribute('data-confirm-changes');
        if (changes) {
          const arr = JSON.parse(changes);
          if (Array.isArray(arr) && arr.length > 0) {
            for (let i = 0; i < arr.length; i++) {
              const p = document.createElement('p');
              p.className = 'confirm-change-line';
              p.textContent = arr[i];
              changesEl.appendChild(p);
            }
          }
        }
      } catch(e) {}
    }

    if (okBtn) {
      okBtn.textContent = (elt.getAttribute ? elt.getAttribute('data-confirm-ok') : '') || (kind === 'danger' ? 'Delete item' : 'Save changes');
      okBtn.disabled = false;
    }
    if (safeBtn) {
      safeBtn.textContent = (elt.getAttribute ? elt.getAttribute('data-confirm-safe') : '') || (kind === 'danger' ? 'Keep item' : 'Go back');
      safeBtn.focus();
    }

    card.onclick = function(e) {
      if (e.target === card) {
        closeCard(card);
      }
    };

    card.onclose = function() {
      try {
        if (elt && typeof elt.focus === 'function') elt.focus();
      } catch(e) {}
    };

    try {
      card.showModal();
    } catch(e) {
      card.setAttribute('open', '');
    }

    const handler = function() {
      okBtn.onclick = null;
      safeBtn.onclick = null;
      card.onclose = null;
      card.onclick = null;
      closeCard(card);
      okBtn.disabled = true;
      evt.detail.issueRequest(true);
      document.removeEventListener('keydown', keyHandler, true);
    };
    const keyHandler = function(e) {
      if (!card.open) return;
      if (e.key === 'Escape') {
        closeCard(card);
        document.removeEventListener('keydown', keyHandler, true);
        return;
      }
    };
    okBtn.onclick = handler;
    safeBtn.onclick = function() {
      okBtn.onclick = null;
      safeBtn.onclick = null;
      card.onclose = null;
      card.onclick = null;
      closeCard(card);
      document.removeEventListener('keydown', keyHandler, true);
    };
    document.addEventListener('keydown', keyHandler, true);
  }, true);
})();
