// EDA-Game annotations: highlights, margin notes and Ask Claude on any selected passage.
//
// Where it runs decides where things are kept. Inside the claude.ai artifact, notes and
// conversations live in the artifact's database (every device, and Claude can read them),
// and Ask Claude answers on the reader's own Claude account. On GitHub Pages there is no
// runtime: notes stay in this browser and Ask Claude is hidden.
(() => {
  const LOCAL = "eda-game-notes";
  const state = { notes: [], threads: [], db: null, sample: null, ready: false, ctx: null };
  const listeners = new Set();
  const changed = () => listeners.forEach((fn) => fn());
  const now = () => new Date().toISOString();
  const uid = () => Date.now().toString(36) + Math.random().toString(36).slice(2, 8);
  const escH = (s) => String(s ?? "").replace(/[&<>"]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" }[c]));

  /* ---------- Storage: artifact db when present, else this browser ---------- */

  const local = {
    read() { try { return JSON.parse(localStorage.getItem(LOCAL) || "{}"); } catch { return {}; } },
    write(v) { try { localStorage.setItem(LOCAL, JSON.stringify(v)); } catch {} },
  };

  async function start() {
    const use = window.claude?.use;
    if (use) {
      const [db, sample] = await Promise.all([use("db").catch(() => null), use("sample").catch(() => null)]);
      state.db = db; state.sample = sample;
    }
    if (state.db) {
      const sub = (name, key) => state.db.collection(name).onSnapshot((snap) => {
        state[key] = snap.docs.map((d) => ({ id: d.id, ...d.data() }));
        changed();
      }, () => {});
      sub("notes", "notes");
      sub("threads", "threads");
    } else {
      const v = local.read();
      state.notes = v.notes || [];
      state.threads = [];
    }
    state.ready = true;
    document.body.classList.toggle("can-ask", !!state.sample);
    changed();
  }

  async function save(kind, item) {
    const key = kind === "notes" ? "notes" : "threads";
    item.updated = now();
    const i = state[key].findIndex((x) => x.id === item.id);
    if (i >= 0) state[key][i] = item; else state[key].push(item);
    changed();
    if (state.db) {
      const { id, ...body } = item;
      await state.db.collection(kind).doc(id).set(JSON.parse(JSON.stringify(body)));
    } else if (kind === "notes") {
      local.write({ notes: state.notes });
    }
  }
  async function remove(kind, id) {
    const key = kind === "notes" ? "notes" : "threads";
    state[key] = state[key].filter((x) => x.id !== id);
    changed();
    if (state.db) await state.db.collection(kind).doc(id).delete();
    else if (kind === "notes") local.write({ notes: state.notes });
  }

  /* ---------- Text anchoring inside the lesson ---------- */

  // Lesson text as one string, skipping KaTeX's hidden MathML so a quote reads like the page.
  function textNodes(root) {
    const out = [];
    const walk = document.createTreeWalker(root, NodeFilter.SHOW_TEXT, {
      acceptNode: (n) => (n.parentElement.closest(".katex-mathml, .ann-ui") ? NodeFilter.FILTER_REJECT : NodeFilter.FILTER_ACCEPT),
    });
    for (let n; (n = walk.nextNode());) out.push(n);
    return out;
  }
  function flatten(root) {
    const nodes = textNodes(root);
    let text = "";
    const spans = nodes.map((n) => { const s = { node: n, start: text.length }; text += n.data; s.end = text.length; return s; });
    return { text, spans };
  }
  function rangeOffsets(root, range) {
    const { text, spans } = flatten(root);
    const pos = (container, offset) => {
      if (container.nodeType === Node.TEXT_NODE) {
        const s = spans.find((x) => x.node === container);
        if (s) return s.start + offset;
      }
      // An element boundary: the first text node at or after it.
      const probe = document.createRange();
      probe.setStart(container, offset);
      const s = spans.find((x) => probe.comparePoint(x.node, 0) >= 0);
      return s ? s.start : text.length;
    };
    const a = pos(range.startContainer, range.startOffset), b = pos(range.endContainer, range.endOffset);
    return { text, a: Math.min(a, b), b: Math.max(a, b) };
  }
  function locate(root, quote, prefix) {
    const { text, spans } = flatten(root);
    let best = -1;
    for (let i = text.indexOf(quote); i >= 0; i = text.indexOf(quote, i + 1)) {
      if (best < 0) best = i;
      if (prefix && text.slice(Math.max(0, i - prefix.length), i) === prefix) { best = i; break; }
    }
    return best < 0 ? null : { spans, a: best, b: best + quote.length };
  }
  function wrap(spans, a, b, make) {
    for (const s of spans) {
      if (s.end <= a || s.start >= b) continue;
      const from = Math.max(a, s.start) - s.start, to = Math.min(b, s.end) - s.start;
      if (to <= from || !s.node.parentNode) continue;
      let node = s.node;
      if (from > 0) node = node.splitText(from);
      if (to - from < node.data.length) node.splitText(to - from);
      const m = make();
      node.parentNode.insertBefore(m, node);
      m.appendChild(node);
    }
  }
  function unpaint(root) {
    root.querySelectorAll("mark.ann").forEach((m) => m.replaceWith(...m.childNodes));
    root.normalize();
  }
  function paint() {
    const root = document.getElementById("lesson");
    if (!root || !state.ctx) return;
    unpaint(root);
    const mine = state.notes.filter((n) => n.concept === state.ctx.concept);
    for (const n of mine) {
      const hit = locate(root, n.quote, n.prefix);
      if (!hit) continue;
      wrap(hit.spans, hit.a, hit.b, () => {
        const m = document.createElement("mark");
        m.className = "ann" + (n.text ? " has-note" : "");
        m.dataset.id = n.id;
        m.tabIndex = 0;
        if (n.text) m.title = n.text;
        return m;
      });
    }
    renderMargin();
  }

  /* ---------- Selection toolbar ---------- */

  const bar = document.createElement("div");
  bar.className = "ann-ui ann-bar";
  bar.setAttribute("role", "toolbar");
  bar.setAttribute("aria-label", "Selection");
  bar.hidden = true;
  bar.innerHTML = `
    <button data-act="hl" title="Highlight (H)">Highlight</button>
    <button data-act="note" title="Add a note (N)">Note</button>
    <button data-act="ask" class="ask-only" title="Start a new question about this passage (A)">Ask Claude</button>
    <button data-act="follow" class="ask-only" title="Quote this into the open conversation (F)">Follow up</button>`;
  let pending = null; // {quote, prefix} for the current selection

  function selectionInLesson() {
    const sel = getSelection();
    const root = document.getElementById("lesson");
    if (!sel || sel.isCollapsed || !root || !sel.rangeCount) return null;
    const range = sel.getRangeAt(0);
    if (!root.contains(range.commonAncestorContainer)) return null;
    const { text, a, b } = rangeOffsets(root, range);
    const quote = text.slice(a, b).replace(/\s+$/, "");
    if (quote.trim().length < 2) return null;
    return { quote, prefix: text.slice(Math.max(0, a - 40), a), rect: range.getBoundingClientRect() };
  }
  function showBar() {
    const s = selectionInLesson();
    if (!s) { bar.hidden = true; pending = null; return; }
    pending = s;
    bar.querySelector('[data-act="follow"]').hidden = !openThread();
    bar.hidden = false;
    const w = bar.offsetWidth, h = bar.offsetHeight;
    const x = Math.min(Math.max(8, s.rect.left + s.rect.width / 2 - w / 2), innerWidth - w - 8);
    const above = s.rect.top - h - 10;
    bar.style.left = x + "px";
    bar.style.top = (above > 60 ? above : s.rect.bottom + 10) + "px";
  }
  let selTimer = 0;
  document.addEventListener("selectionchange", () => { clearTimeout(selTimer); selTimer = setTimeout(showBar, 180); });
  addEventListener("scroll", () => { if (!bar.hidden) showBar(); }, { passive: true });
  bar.addEventListener("mousedown", (e) => e.preventDefault()); // keep the selection alive
  bar.addEventListener("click", (e) => {
    const act = e.target.closest("button")?.dataset.act;
    if (act && pending) act === "hl" ? highlight(pending) : act === "note" ? noteFor(pending) : act === "ask" ? ask(pending) : followUp(pending);
  });
  document.addEventListener("keydown", (e) => {
    if (bar.hidden || !pending || e.metaKey || e.ctrlKey || e.altKey) return;
    if (e.target.closest?.("input, textarea")) return;
    const k = e.key.toLowerCase();
    const map = { h: () => highlight(pending), n: () => noteFor(pending), a: () => state.sample && ask(pending), f: () => state.sample && openThread() && followUp(pending) };
    if (map[k]) { e.preventDefault(); e.stopPropagation(); map[k](); }
    if (e.key === "Escape") { bar.hidden = true; getSelection().removeAllRanges(); }
  }, true);

  function base(sel) {
    const c = state.ctx;
    return { concept: c.concept, domain: c.domain, label: c.label, quote: sel.quote, prefix: sel.prefix, created: now() };
  }
  async function highlight(sel) {
    bar.hidden = true;
    getSelection().removeAllRanges();
    const n = { id: uid(), kind: "highlight", text: "", ...base(sel) };
    await save("notes", n);
    paint();
    return n;
  }
  async function noteFor(sel) {
    const n = await highlight(sel);
    editNote(n.id);
  }

  /* ---------- Note editor ---------- */

  const pop = document.createElement("div");
  pop.className = "ann-ui ann-pop";
  pop.hidden = true;
  pop.setAttribute("role", "dialog");
  pop.setAttribute("aria-label", "Note");
  function editNote(id) {
    const n = state.notes.find((x) => x.id === id);
    const mark = document.querySelector(`mark.ann[data-id="${id}"]`);
    if (!n) return;
    pop.innerHTML = `<blockquote>${escH(n.quote.length > 160 ? n.quote.slice(0, 160) + "…" : n.quote)}</blockquote>
      <textarea id="ann-text" rows="4" placeholder="Write a note to yourself">${escH(n.text)}</textarea>
      <div class="ann-row">
        <button class="btn primary" data-p="save">Save</button>
        ${state.sample ? `<button class="btn" data-p="ask">Ask Claude</button>` : ""}
        <span class="ann-flex"></span>
        <button class="btn quiet" data-p="del">Remove</button>
      </div>`;
    pop.hidden = false;
    const r = (mark || document.getElementById("lesson")).getBoundingClientRect();
    pop.style.left = Math.min(Math.max(8, r.left), innerWidth - pop.offsetWidth - 8) + "px";
    pop.style.top = Math.min(r.bottom + 8, innerHeight - pop.offsetHeight - 8) + "px";
    const ta = pop.querySelector("textarea");
    ta.focus();
    ta.setSelectionRange(ta.value.length, ta.value.length);
    const done = async (keep) => {
      pop.hidden = true;
      if (keep) { n.text = ta.value.trim(); n.kind = n.text ? "note" : "highlight"; await save("notes", n); paint(); }
      mark?.focus();
    };
    pop.onclick = async (e) => {
      const p = e.target.closest("button")?.dataset.p;
      if (p === "save") done(true);
      if (p === "del") { pop.hidden = true; await remove("notes", n.id); paint(); }
      if (p === "ask") { await done(true); ask({ quote: n.quote, prefix: n.prefix }); }
    };
    ta.onkeydown = (e) => {
      if (e.key === "Enter" && (e.metaKey || e.ctrlKey)) { e.preventDefault(); done(true); }
      if (e.key === "Escape") { e.preventDefault(); done(false); }
    };
  }
  document.addEventListener("click", (e) => {
    const m = e.target.closest?.("mark.ann");
    if (m && getSelection().isCollapsed) { editNote(m.dataset.id); return; }
    if (!pop.hidden && !e.target.closest(".ann-pop")) pop.hidden = true;
  });
  document.addEventListener("keydown", (e) => {
    const m = e.target.closest?.("mark.ann");
    if (m && e.key === "Enter") { e.preventDefault(); editNote(m.dataset.id); }
  });

  /* ---------- Ask Claude ---------- */

  const panel = document.createElement("aside");
  panel.className = "ann-ui ann-chat";
  panel.hidden = true;
  panel.setAttribute("aria-label", "Ask Claude");
  let thread = null; // the open conversation
  let ctl = null;

  const openThread = () => (thread && !panel.hidden && thread.concept === state.ctx?.concept ? thread : null);

  function lessonText() {
    const root = document.getElementById("lesson");
    return root ? flatten(root).text.replace(/\n{3,}/g, "\n\n").slice(0, 14000) : "";
  }
  function rules(t) {
    return `You are a tutor inside EDA-Game, a course on electronic design automation for someone becoming a tool designer.
The reader is studying the concept "${t.label}". Answer their questions about it precisely and concisely, the way a senior EDA engineer would explain it to a strong colleague. Use worked numbers where they help. Write Markdown; put mathematics in $...$ or $$...$$ (KaTeX). Keep answers under about 300 words unless the reader asks for more.

The lesson they are reading:
"""
${lessonText()}
"""`;
  }

  function ask(sel) {
    bar.hidden = true;
    getSelection().removeAllRanges();
    thread = { id: uid(), ...base(sel), turns: [] };
    openPanel(sel.quote);
  }
  function followUp(sel) {
    bar.hidden = true;
    getSelection().removeAllRanges();
    const box = panel.querySelector("#ask-text");
    box.value = `> ${sel.quote.replace(/\n/g, " ")}\n\n` + box.value;
    box.focus();
  }
  function openThreadById(id) {
    thread = state.threads.find((t) => t.id === id) || null;
    if (thread) openPanel();
  }

  function openPanel(freshQuote) {
    panel.hidden = false;
    document.body.classList.add("chat-open");
    drawThread();
    const box = panel.querySelector("#ask-text");
    if (freshQuote) box.placeholder = "Ask about this passage. Enter sends; leave it empty to get an explanation.";
    box.focus();
  }
  function closePanel() {
    ctl?.abort();
    panel.hidden = true;
    document.body.classList.remove("chat-open");
  }

  function drawThread() {
    const t = thread;
    const render = window.edaRender || escH;
    panel.innerHTML = `
      <header class="ann-chat-h">
        <b>Ask Claude</b><span class="muted">${escH(t?.label || "")}</span>
        <span class="ann-flex"></span>
        <button class="btn quiet" data-c="new" title="Start a new question">New</button>
        <button class="btn quiet" data-c="close" aria-label="Close">Close</button>
      </header>
      <div class="ann-chat-body" id="ask-log" aria-live="polite">
        ${t?.quote ? `<blockquote class="ann-quote">${escH(t.quote.length > 400 ? t.quote.slice(0, 400) + "…" : t.quote)}</blockquote>` : ""}
        ${(t?.turns || []).map((m) => `<div class="msg ${m.role}">${m.role === "assistant" ? render(m.content) : escH(m.shown || m.content).replace(/\n/g, "<br>")}</div>`).join("")}
        ${!t?.turns?.length ? `<p class="muted small">Your question and Claude's answer are saved with this lesson.</p>` : ""}
      </div>
      <form class="ann-chat-f" id="ask-form">
        <textarea id="ask-text" rows="3" placeholder="Ask a follow-up. Enter sends, Shift+Enter adds a line."></textarea>
        <div class="ann-row"><span class="muted small" id="ask-status"></span><span class="ann-flex"></span>
          <button type="button" class="btn quiet" id="ask-stop" hidden>Stop</button>
          <button class="btn primary" id="ask-send">Send</button></div>
      </form>`;
    const log = panel.querySelector("#ask-log");
    log.scrollTop = log.scrollHeight;
    panel.querySelector('[data-c="close"]').onclick = closePanel;
    panel.querySelector('[data-c="new"]').onclick = () => { thread = { id: uid(), ...base({ quote: "", prefix: "" }), turns: [] }; drawThread(); panel.querySelector("#ask-text").focus(); };
    const box = panel.querySelector("#ask-text");
    box.onkeydown = (e) => {
      if (e.key === "Enter" && !e.shiftKey && !e.isComposing) { e.preventDefault(); send(); }
      if (e.key === "Escape") { e.preventDefault(); closePanel(); }
    };
    panel.querySelector("#ask-form").onsubmit = (e) => { e.preventDefault(); send(); };
    panel.querySelector("#ask-stop").onclick = () => ctl?.abort();
  }

  async function send() {
    if (!state.sample || !thread) return;
    const box = panel.querySelector("#ask-text");
    let q = box.value.trim();
    if (!q) q = thread.quote && !thread.turns.length ? "Explain this passage." : "";
    if (!q) return;
    const t = thread;
    const first = !t.turns.length;
    const lead = first && t.quote ? `The reader highlighted this passage:\n"""\n${t.quote}\n"""\n\n${q}` : q;
    t.turns.push({ role: "user", content: lead, shown: q });
    drawThread();
    const log = panel.querySelector("#ask-log");
    const bubble = document.createElement("div");
    bubble.className = "msg assistant thinking";
    bubble.textContent = "Thinking…";
    log.appendChild(bubble);
    log.scrollTop = log.scrollHeight;
    const status = panel.querySelector("#ask-status"), stop = panel.querySelector("#ask-stop"), sendBtn = panel.querySelector("#ask-send");
    stop.hidden = false; sendBtn.disabled = true;
    ctl = new AbortController();
    const render = window.edaRender || escH;
    const input = [{ role: "user", content: rules(t) }, ...t.turns.slice(-12).map(({ role, content }) => ({ role, content }))];
    try {
      const { text, truncated } = await state.sample(input, {
        cache: false, signal: ctl.signal,
        onText: ({ text }) => { bubble.classList.remove("thinking"); bubble.innerHTML = render(text); log.scrollTop = log.scrollHeight; },
      });
      t.turns.push({ role: "assistant", content: text });
      if (truncated) status.textContent = "The answer was cut short. Ask for less at a time.";
      await save("threads", t);
    } catch (e) {
      bubble.classList.remove("thinking");
      if (e?.text) { bubble.innerHTML = render(e.text); t.turns.push({ role: "assistant", content: e.text }); await save("threads", t); }
      else bubble.remove();
      const copy = {
        cancelled: "", not_granted: "Ask Claude is turned off for this page. You can allow it from the artifact's Permissions menu.",
        sampling_disabled: "Claude is not available on this account.", rate_limited: "You have hit a usage limit. Try again in a little while.",
        session_expired: "Sign in to claude.ai again, then retry.", refused: "Claude declined that question. Try rephrasing it.",
        prompt_too_large: "That conversation is too long. Start a new question.",
      };
      status.textContent = e?.code in copy ? copy[e.code] : "Something went wrong reaching Claude. Try again.";
      if (e?.code === "not_granted" || e?.code === "sampling_disabled") { state.sample = null; document.body.classList.remove("can-ask"); }
    } finally {
      stop.hidden = true; sendBtn.disabled = false;
      const b = panel.querySelector("#ask-text"); if (b) { b.value = ""; b.focus(); }
    }
  }

  /* ---------- Margin list for the current concept ---------- */

  function renderMargin() {
    const host = document.getElementById("ann-margin");
    if (!host || !state.ctx) return;
    const notes = state.notes.filter((n) => n.concept === state.ctx.concept).sort((a, b) => a.created.localeCompare(b.created));
    const threads = state.threads.filter((t) => t.concept === state.ctx.concept).sort((a, b) => b.updated.localeCompare(a.updated));
    if (!notes.length && !threads.length) {
      host.innerHTML = `<p class="muted small ann-hint">Select any passage to highlight it, add a note${state.sample ? " or ask Claude about it" : ""}.</p>`;
      return;
    }
    host.innerHTML = `<h2 class="ann-h">Your notes on this page</h2>
      <ul class="ann-list">
        ${notes.map((n) => `<li><button class="ann-item" data-jump="${n.id}"><span class="ann-q">${escH(n.quote.slice(0, 140))}${n.quote.length > 140 ? "…" : ""}</span>${n.text ? `<span class="ann-t">${escH(n.text)}</span>` : ""}</button></li>`).join("")}
        ${threads.map((t) => `<li><button class="ann-item" data-thread="${t.id}"><span class="ann-k">Conversation</span><span class="ann-t">${escH(t.turns.find((m) => m.role === "user")?.shown || "Question")}</span></button></li>`).join("")}
      </ul>`;
    host.querySelectorAll("[data-jump]").forEach((b) => b.onclick = () => {
      const m = document.querySelector(`mark.ann[data-id="${b.dataset.jump}"]`);
      if (m) { m.scrollIntoView({ block: "center", behavior: matchMedia("(prefers-reduced-motion: reduce)").matches ? "auto" : "smooth" }); m.classList.add("flash"); setTimeout(() => m.classList.remove("flash"), 1200); }
      else editNote(b.dataset.jump);
    });
    host.querySelectorAll("[data-thread]").forEach((b) => b.onclick = () => openThreadById(b.dataset.thread));
  }

  /* ---------- All notes page ---------- */

  function notesPage(app) {
    const groups = {};
    for (const n of state.notes) (groups[n.concept] ||= { domain: n.domain, label: n.label, notes: [], threads: [] }).notes.push(n);
    for (const t of state.threads) (groups[t.concept] ||= { domain: t.domain, label: t.label, notes: [], threads: [] }).threads.push(t);
    const keys = Object.keys(groups).sort((a, b) => {
      const last = (g) => [...g.notes, ...g.threads].map((x) => x.updated || x.created).sort().pop();
      return last(groups[b]).localeCompare(last(groups[a]));
    });
    app.innerHTML = `<div class="page narrow"><h1>Notes</h1>
      <p class="lede">${state.db ? "Everything you have highlighted, written or asked, on every device." : "Highlights and notes are kept in this browser."}</p>
      ${!state.ready ? `<p class="muted">Loading…</p>` : !keys.length ? `<p class="muted">Nothing yet. Open any lesson and select a passage to highlight it or add a note.</p>` :
        keys.map((k) => { const g = groups[k]; return `<section class="ann-group">
          <h2><a href="#/d/${g.domain}/${k}">${escH(g.label)}</a></h2>
          ${g.notes.map((n) => `<div class="ann-entry"><blockquote>${escH(n.quote)}</blockquote>${n.text ? `<p>${escH(n.text)}</p>` : ""}</div>`).join("")}
          ${g.threads.map((t) => `<details class="ann-entry"><summary>${escH(t.turns.find((m) => m.role === "user")?.shown || "Conversation")}</summary>
            ${t.turns.map((m) => `<div class="msg ${m.role}">${m.role === "assistant" ? (window.edaRender || escH)(m.content) : escH(m.shown || m.content)}</div>`).join("")}</details>`).join("")}
        </section>`; }).join("")}</div>`;
  }

  /* ---------- Wiring ---------- */

  document.addEventListener("eda:render", (e) => {
    const { page, concept, domain, label } = e.detail || {};
    bar.hidden = true; pop.hidden = true;
    state.ctx = page === "concept" ? { concept, domain, label } : null;
    if (thread && thread.concept !== state.ctx?.concept) closePanel();
    if (page === "notes") notesPage(document.getElementById("app"));
    if (state.ctx) {
      document.getElementById("practice")?.insertAdjacentHTML("beforebegin", `<section class="ann-margin" id="ann-margin" aria-label="Your notes"></section>`);
      paint();
    }
  });
  listeners.add(() => {
    if (state.ctx) paint();
    if (document.body.dataset.page === "notes") notesPage(document.getElementById("app"));
  });

  document.addEventListener("DOMContentLoaded", () => document.body.append(bar, pop, panel));
  if (document.readyState !== "loading") document.body.append(bar, pop, panel);
  start();
})();
