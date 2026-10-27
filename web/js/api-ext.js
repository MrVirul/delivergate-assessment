htmx.defineExtension('api', {
  onEvent: function(name, evt) {
    if (name === 'htmx:configRequest') {
      const d = evt.detail;
      const token = sessionStorage.getItem('token');
      if (token) {
        d.headers['Authorization'] = 'Bearer ' + token;
      }
      if (['POST', 'PUT', 'PATCH'].includes(d.verb.toUpperCase())) {
        d.headers['Content-Type'] = 'application/json';
      }
    }
    return true;
  },
  encodeParameters: function(xhr, parameters, elt) {
    xhr.overrideMimeType('text/json');
    // build form
    const form = elt.tagName === 'FORM' ? elt : (elt.querySelector('form') || elt.closest('form'));
    const obj = UI.coerceForm ? UI.coerceForm(parameters, form) : {};
    // if no entries, try to read form
    if (Object.keys(obj).length === 0 && form) {
      const fd = new FormData(form);
      obj = UI.coerceForm(fd, form);
    }
    return JSON.stringify(obj);
  },
  transformResponse: function(text, xhr, elt) {
    if (xhr.status >= 400) {
      return text;
    }
    const type = xhr.getResponseHeader('Content-Type') || '';
    if (!type.includes('json') || !text) {
      return text;
    }
    if (xhr.status === 204) {
      return '';
    }
    let body;
    try {
      body = JSON.parse(text);
    } catch (e) {
      return text;
    }
    const key = elt.closest('[data-render]')?.dataset.render || elt.dataset.render;
    if (!key) {
      return text;
    }
    if (RENDERERS && RENDERERS[key]) {
      try {
        return RENDERERS[key](body);
      } catch (e) {
        return text;
      }
    }
    return text;
  }
});
