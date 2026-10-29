// EDA-Game static reader. Hash routes: #/  #/tracks  #/glossary  #/about  #/d/<domain>[/<concept>]
const app = document.getElementById("app");
const cache = {};
let index = null;

const store = {
  get() { try { return JSON.parse(localStorage.getItem("eda-game") || "{}"); } catch { return {}; } },
  set(v) { try { localStorage.setItem("eda-game", JSON.stringify(v)); } catch {} },
};
// Where the learner last was, for "Continue" on the home page.
const last = {
  get() { try { return JSON.parse(localStorage.getItem("eda-game-last") || "null"); } catch { return null; } },
  set(v) { try { localStorage.setItem("eda-game-last", JSON.stringify(v)); } catch {} },
};
const STREAK = 3; // correct answers in a row to mark a concept learned here

async function load(path) {
  if (!cache[path]) cache[path] = fetch(path).then((r) => r.json());
  return cache[path];
}
const esc = (s) => String(s).replace(/[&<>"]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" }[c]));
const unwrapP = (html) => html.replace(/^\s*<p>|<\/p>\s*$/g, "");
const humanize = (slug) => { const s = String(slug || "").replace(/_/g, " "); return s.charAt(0).toUpperCase() + s.slice(1); };
const pct = (n, d) => (d ? Math.round((100 * n) / d) : 0);
const meter = (n, d, label) => `<span class="meter" role="img" aria-label="${esc(label || `${n} of ${d} learned`)}"><i style="width:${pct(n, d)}%"></i></span>`;

// Markdown with $…$ and $$…$$ math: lift the math out, render the markdown, put KaTeX back.
window.edaRender = (md) => render(md);
function render(md) {
  const math = [];
  const keep = (tex, display) => { math.push([tex, display]); return `@@M${math.length - 1}@@`; };
  md = md.replace(/\$\$([\s\S]+?)\$\$/g, (_, t) => keep(t, true))
         .replace(/\\\[([\s\S]+?)\\\]/g, (_, t) => keep(t, true))
         .replace(/\\\(([\s\S]+?)\\\)/g, (_, t) => keep(t, false))
         .replace(/(^|[^\\$])\$([^\s$](?:[^$\n]*?[^\s\\$])?)\$/g, (_, pre, t) => pre + keep(t, false));
  let html = marked.parse(md);
  return html.replace(/@@M(\d+)@@/g, (_, i) => {
    const [tex, display] = math[+i];
    try { return katex.renderToString(tex, { displayMode: display, throwOnError: false }); }
    catch { return esc(tex); }
  });
}

function progress(domain, p = store.get()) {
  return Object.keys(p).filter((k) => p[k].done && p[k].domain === domain).length;
}
const domainOf = (pid) => index.domains.find((d) => pid.startsWith(d.prefix + "."));

// Readable subtopic titles come from the glossary; until it loads, the slug stands in.
let subTitleMap = null;
function subTitles() {
  if (!subTitleMap) subTitleMap = load("data/glossary.json").then((g) => {
    const m = {};
    g.domains.forEach((d) => d.subdomains.forEach((s) => s.concepts.forEach((c) => {
      const parts = c.id.split(".");
      m[`${d.id}:${parts[parts.length - 2]}`] = s.title;
    })));
    return m;
  }).catch(() => ({}));
  return subTitleMap;
}
let titles = {};
const subTitle = (domainId, sub) => titles[`${domainId}:${sub}`] || humanize(sub);
function upgradeSubTitles(domainId) {
  subTitles().then((m) => {
    titles = m;
    document.querySelectorAll("[data-sub]").forEach((el) => { el.textContent = subTitle(domainId, el.dataset.sub); });
  });
}

// Inside claude.ai the page is a frame, and keys go to whichever frame was clicked last, so a
// reader who scrolls here from the chat presses keys into claude.ai. Ask for focus when the
// pointer moves over the page, and mark the body while the page cannot hear keys so the
// practice hint can say to click first.
let focusAsked = 0;
addEventListener("pointermove", () => {
  if (document.hasFocus() || Date.now() - focusAsked < 400) return;
  focusAsked = Date.now();
  try { window.focus(); } catch {}
}, { passive: true });
const markFocus = () => document.body.classList.toggle("unfocused", !document.hasFocus());
addEventListener("focus", markFocus);
addEventListener("blur", markFocus);
markFocus();

// One keyboard handler at a time, owned by the practice box on screen.
let keyHandler = null;
addEventListener("keydown", (e) => {
  if (e.key === "Escape" && document.body.classList.contains("drawer-open")) return closeDrawer(true);
  if (keyHandler && !document.body.classList.contains("drawer-open")) keyHandler(e);
});

/* ---------- Home ---------- */

function home() {
  const p = store.get();
  const layers = {};
  index.domains.forEach((d) => (layers[d.layer] ||= []).push(d));
  const learned = Object.values(p).filter((v) => v && v.done).length;
  const here = last.get();
  const first = index.domains[0];
  const cta = here && index.domains.some((d) => d.id === here.d)
    ? `<a class="btn primary lg" href="#/d/${here.d}/${here.c}">Continue</a>
       <span class="cta-note">${esc(here.label)}</span>`
    : `<a class="btn primary lg" href="#/d/${first.id}">Start learning</a>
       <span class="cta-note">Begins with ${esc(first.name)}</span>`;
  app.innerHTML = `
  <section class="hero">
    <h1>Learn how chip design tools work.</h1>
    <p class="lede">Electronic design automation, from first principles to agent-first tools.</p>
    <div class="cta">${cta}</div>
    ${learned ? `<p class="overall">${meter(learned, index.concepts)}<span>${learned} of ${index.concepts} concepts learned</span></p>` : ""}
  </section>
  <section class="map" aria-labelledby="map-h">
    <div class="map-head"><h2 id="map-h">Domains</h2><p class="muted">Each layer builds on the ones above it.</p></div>
    ${Object.keys(layers).sort((a, b) => a - b).map((l) => `
    <section class="layer" aria-label="Layer ${l}">
      <h3 class="layer-n">Layer ${l}</h3>
      <ul class="tiles">${layers[l].map((d) => {
        const n = progress(d.id, p);
        return `<li><a class="tile${n ? " started" : ""}${n === d.count ? " complete" : ""}" href="#/d/${d.id}" title="${esc(d.description)}">
          <span class="t-name">${esc(d.name)}</span>
          <span class="t-meta">${n ? `${meter(n, d.count)}<span>${n}/${d.count}</span>` : `<span>${d.count} concepts</span>`}</span>
        </a></li>`;
      }).join("")}</ul>
    </section>`).join("")}
  </section>`;
}

/* ---------- Tracks ---------- */

function tracks() {
  const p = store.get();
  const byId = Object.fromEntries(index.domains.map((d) => [d.id, d]));
  app.innerHTML = `<div class="page">
  <header class="page-head"><h1>Tracks</h1><p class="lede">Paths through the domains, in dependency order.</p></header>
  <ul class="tracks">${index.tracks.map((t) => {
    const ds = t.courses.map((c) => byId[c]).filter(Boolean);
    const total = ds.reduce((s, d) => s + d.count, 0);
    const done = ds.reduce((s, d) => s + progress(d.id, p), 0);
    const nextD = ds.find((d) => progress(d.id, p) < d.count) || ds[0];
    return `<li class="track">
      <div class="tr-main">
        <h2>${esc(t.name)}</h2>
        <p>${esc(t.description)}</p>
        <p class="tr-meta">${ds.length} domains · ${total} concepts${done ? ` · ${pct(done, total)}% learned` : ""}</p>
      </div>
      <a class="btn" href="#/d/${nextD.id}">${done ? "Continue" : "Start"}</a>
      <details class="tr-path"><summary>Show the path</summary>
        <ol>${ds.map((d) => { const n = progress(d.id, p);
          return `<li><a href="#/d/${d.id}">${esc(d.name)}</a>${n ? `<span class="muted">${n}/${d.count}</span>` : ""}</li>`; }).join("")}</ol>
      </details>
    </li>`;
  }).join("")}</ul></div>`;
}

/* ---------- About ---------- */

function about() {
  app.innerHTML = `<div class="page"><article class="prose">${render(`# About

EDA-Game teaches electronic design automation as a directed acyclic graph of concepts. Every domain sits on a
layer, and a concept depends only on concepts in its own domain, in lower layers, or in the mathematics substrate,
so the order you meet things in is always one where the prerequisites come first.

## This site

This is the **static reader**: every lesson, plus a fixed set of practice problems per concept. A concept
counts as learned after ${STREAK} correct answers in a row. Progress is stored only in this browser.

## The full engine

The adaptive engine (fresh problems on every draw, spaced review, a placement diagnostic, mastery gating
over the whole graph) runs locally from the repository:

\`\`\`
go build -o bin/mathua ./cmd/mathua
DATABASE_URL=sqlite:eda.db ./bin/mathua --serve --port 8080
\`\`\`

Built on the open-source [mathua](https://github.com/chuma-beep/mathua) engine (MIT).`)}</article></div>`;
}

/* ---------- Domain: landing and concept ---------- */

function sidebar(meta, concepts, cur, p) {
  const n = progress(meta.id, p);
  let sub = null;
  const groups = [];
  concepts.forEach((x) => {
    if (x.sub !== sub) { groups.push({ sub: x.sub, items: [] }); sub = x.sub; }
    groups[groups.length - 1].items.push(x);
  });
  return `<aside class="side" id="side" aria-label="Concepts in ${esc(meta.name)}">
    <div class="side-head">
      <a class="side-dom${cur ? "" : " on"}" href="#/d/${meta.id}"${cur ? "" : ` aria-current="page"`}>${esc(meta.name)}</a>
      <div class="side-prog">${meter(n, meta.count)}<span>${n}/${meta.count}</span></div>
    </div>
    ${groups.map((g) => `<h3 class="side-h" data-sub="${esc(g.sub)}">${esc(subTitle(meta.id, g.sub))}</h3>
      <ul>${g.items.map((x) => {
        const s = p[x.id];
        const cls = [x.id === cur ? "on" : "", s?.done ? "done" : s?.n ? "started" : ""].join(" ").trim();
        return `<li><a class="${cls}" href="#/d/${meta.id}/${x.id}"${x.id === cur ? ` aria-current="page"` : ""}>
          <span class="dot" aria-hidden="true"></span><span>${esc(x.label)}</span>${s?.done ? `<span class="sr-only"> (learned)</span>` : ""}</a></li>`;
      }).join("")}</ul>`).join("")}
  </aside><div class="scrim" hidden></div>`;
}

function openDrawer() {
  document.body.classList.add("drawer-open");
  document.querySelector(".scrim")?.removeAttribute("hidden");
  document.querySelector(".contents-btn")?.setAttribute("aria-expanded", "true");
  const side = document.getElementById("side");
  (side.querySelector("a.on") || side.querySelector("a"))?.focus({ preventScroll: true });
  centerCurrent();
}
function closeDrawer(refocus) {
  document.body.classList.remove("drawer-open");
  document.querySelector(".scrim")?.setAttribute("hidden", "");
  const btn = document.querySelector(".contents-btn");
  btn?.setAttribute("aria-expanded", "false");
  if (refocus) btn?.focus();
}
function centerCurrent() {
  const side = document.getElementById("side"), on = side?.querySelector("li a.on");
  if (side && on) side.scrollTop = on.offsetTop - side.clientHeight / 2 + on.offsetHeight / 2;
}

function studyShell(meta, concepts, cur, p, body) {
  app.innerHTML = `<div class="study">
    <button class="contents-btn" aria-expanded="false" aria-controls="side">
      <svg width="16" height="16" viewBox="0 0 16 16" aria-hidden="true"><path d="M2 4h12M2 8h12M2 12h12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
      <span>${esc(meta.name)}</span><span class="muted">${progress(meta.id, p)}/${meta.count}</span>
    </button>
    ${sidebar(meta, concepts, cur, p)}
    <div class="reader">${body}</div>
  </div>`;
  document.querySelector(".contents-btn").onclick = () =>
    document.body.classList.contains("drawer-open") ? closeDrawer(true) : openDrawer();
  document.querySelector(".scrim").onclick = () => closeDrawer(true);
  centerCurrent();
  upgradeSubTitles(meta.id);
}

async function domain(id, conceptId) {
  const meta = index.domains.find((d) => d.id === id);
  if (!meta) return home();
  const concepts = await load(`data/${id}.json`);
  const c = concepts.find((x) => x.id === conceptId);
  const p = store.get();
  if (!c) return domainLanding(meta, concepts, p);

  last.set({ d: id, c: c.id, label: c.label });
  const pre = c.prereqs.map(([pid, label]) => {
    const owner = domainOf(pid);
    return owner ? `<a href="#/d/${owner.id}/${pid}">${esc(label)}</a>` : `<span title="Mathematics prerequisite">${esc(label)}</span>`;
  });
  const i = concepts.indexOf(c);
  const prev = concepts[i - 1], next = concepts[i + 1];
  studyShell(meta, concepts, c.id, p, `
    <nav class="crumbs" aria-label="Breadcrumb">
      <a href="#/d/${id}">${esc(meta.name)}</a><span aria-hidden="true">›</span><span data-sub="${esc(c.sub)}">${esc(subTitle(id, c.sub))}</span>
    </nav>
    <article class="prose lesson" id="lesson">${render(c.lesson || `# ${c.label}`)}</article>
    <section class="practice" id="practice" aria-labelledby="ph"></section>
    <nav class="pager" aria-label="Previous and next concept">
      ${prev ? `<a class="pg prev" href="#/d/${id}/${prev.id}" rel="prev"><span class="pg-k">Previous</span><span class="pg-l">${esc(prev.label)}</span></a>` : `<span></span>`}
      ${next ? `<a class="pg next" href="#/d/${id}/${next.id}" rel="next"><span class="pg-k">Next</span><span class="pg-l">${esc(next.label)}</span></a>` : `<a class="pg next" href="#/d/${id}"><span class="pg-k">Done</span><span class="pg-l">Back to ${esc(meta.name)}</span></a>`}
    </nav>`);
  if (pre.length) {
    const h1 = document.querySelector("#lesson h1");
    const html = `<p class="builds"><span class="builds-k">Builds on</span> ${pre.join('<span class="sep" aria-hidden="true">·</span>')}</p>`;
    h1 ? h1.insertAdjacentHTML("afterend", html) : document.getElementById("lesson").insertAdjacentHTML("afterbegin", html);
  }
  practice(id, c, next);
  window.__edaCtx = { page: "concept", concept: c.id, domain: id, label: c.label };
}

function domainLanding(meta, concepts, p) {
  const n = progress(meta.id, p);
  const subs = [];
  concepts.forEach((x) => {
    let s = subs.find((g) => g.sub === x.sub);
    if (!s) subs.push((s = { sub: x.sub, items: [] }));
    s.items.push(x);
  });
  const here = last.get();
  const resume = here && here.d === meta.id && concepts.find((x) => x.id === here.c && !p[x.id]?.done);
  const nextUp = resume || concepts.find((x) => !p[x.id]?.done);
  const started = !!resume || concepts.some((x) => p[x.id]?.n);
  const up = (meta.upstream || []).map((u) => {
    const d = index.domains.find((x) => x.id === u || x.prefix === u);
    return d ? `<a href="#/d/${d.id}">${esc(d.name)}</a>` : `<span>${esc(humanize(u))}</span>`;
  });
  studyShell(meta, concepts, null, p, `
    <header class="dom-head">
      <p class="eyebrow">Layer ${meta.layer} · ${meta.count} concepts</p>
      <h1>${esc(meta.name)}</h1>
      <p class="lede">${esc(meta.description)}</p>
      <div class="dom-prog">${meter(n, meta.count)}<span>${n || started ? `${n} of ${meta.count} learned` : "Not started"}</span></div>
      <div class="cta">
        ${nextUp ? `<a class="btn primary lg" href="#/d/${meta.id}/${nextUp.id}">${started ? "Continue" : "Start"}</a>
          <span class="cta-note">${esc(nextUp.label)}</span>` : `<span class="cta-note">Every concept here is learned.</span>`}
      </div>
      ${up.length ? `<p class="builds"><span class="builds-k">Builds on</span> ${up.join('<span class="sep" aria-hidden="true">·</span>')}</p>` : ""}
    </header>
    <section aria-labelledby="subs-h">
      <h2 id="subs-h" class="sec-h">Subtopics</h2>
      <ol class="subs">${subs.map((s) => {
        const k = s.items.filter((x) => p[x.id]?.done).length;
        const go = s.items.find((x) => !p[x.id]?.done) || s.items[0];
        return `<li><a href="#/d/${meta.id}/${go.id}">
          <span class="s-title" data-sub="${esc(s.sub)}">${esc(subTitle(meta.id, s.sub))}</span>
          <span class="s-first">${esc(s.items[0].label)}${s.items.length > 1 ? ` and ${s.items.length - 1} more` : ""}</span>
          <span class="s-meta">${k ? `${meter(k, s.items.length)}<span>${k}/${s.items.length}</span>` : `<span>${s.items.length}</span>`}</span>
        </a></li>`;
      }).join("")}</ol>
    </section>`);
}

/* ---------- Practice ---------- */

// Options are written into the stem as "A) …" paragraphs, as the engine's client parses them.
function splitChoices(q) {
  const parts = q.split(/\n\n(?=[A-E]\) )/);
  if (parts.length < 3) return { stem: q, options: [] };
  return { stem: parts[0], options: parts.slice(1).map((o) => o.replace(/^[A-E]\) /, "").trim()) };
}

// Numeric answers are accepted within half a unit of the last decimal place shown.
function numericMatch(expected, given) {
  const num = (s) => parseFloat(String(s).replace(/,/g, "").replace(/[^\d.eE+-].*$/, ""));
  const e = num(expected), g = num(given);
  if (Number.isNaN(e) || Number.isNaN(g)) return String(expected).trim().toLowerCase() === String(given).trim().toLowerCase();
  const dec = (String(expected).match(/\.(\d+)/) || ["", ""])[1].length;
  return Math.abs(e - g) <= 0.5 * 10 ** -dec + 1e-9;
}

function streakDots(st) {
  const on = st.done ? STREAK : st.streak;
  const label = st.done ? "Learned" : `Streak ${st.streak} of ${STREAK}`;
  return `<span class="streak${st.done ? " learned" : ""}" aria-label="${label}">
    ${Array.from({ length: STREAK }, (_, k) => `<i class="${k < on ? "on" : ""}"></i>`).join("")}
    <span aria-hidden="true">${st.done ? "Learned" : `${st.streak}/${STREAK}`}</span></span>`;
}

function practice(domainId, c, nextConcept, focusBox) {
  const box = document.getElementById("practice");
  keyHandler = null;
  if (!c.problems.length) {
    box.innerHTML = `<div class="p-head"><h2 id="ph">Practice</h2></div><p class="muted">No practice problems for this concept.</p>`;
    return;
  }
  const st = store.get()[c.id] || { domain: domainId, streak: 0, done: false, n: 0 };
  const prob = c.problems[st.n % c.problems.length];
  const { stem, options } = splitChoices(prob.q);
  const keys = options.map((_, k) => String.fromCharCode(65 + k));
  box.innerHTML = `
    <div class="p-head"><h2 id="ph">Practice</h2>${streakDots(st)}</div>
    <div class="q prose">${render(stem)}</div>
    ${options.length ? `
      <div class="opts" role="radiogroup" aria-labelledby="ph">${options.map((o, k) => `
        <button type="button" role="radio" aria-checked="false" data-o="${k}">
          <kbd aria-hidden="true">${keys[k]}</kbd><span class="o-t">${unwrapP(render(o))}</span></button>`).join("")}
      </div>
      <div class="p-act"><button class="btn primary" id="check" disabled>Check</button>
        <span class="hint" data-keys>Press ${keys[0]}–${keys[keys.length - 1]} or 1–${keys.length} to choose, Enter to check</span></div>`
    : `
      <form class="num" autocomplete="off">
        <label for="ans" class="sr-only">Your answer</label>
        <input id="ans" class="field" placeholder="Your answer" inputmode="decimal" spellcheck="false">
        <button class="btn primary">Check</button>
      </form>
      <p class="hint" data-keys>Enter to check. Numbers are accepted to the precision shown in the answer.</p>`}
    <div id="fb" aria-live="polite"></div>`;

  let chosen = null, graded = false;
  const optBtns = [...box.querySelectorAll("button[data-o]")];
  const checkBtn = box.querySelector("#check");
  const select = (k) => {
    if (graded || k < 0 || k >= options.length) return;
    chosen = k;
    optBtns.forEach((b, j) => b.setAttribute("aria-checked", String(j === k)));
    checkBtn.disabled = false;
  };
  const nextProblem = () => practice(domainId, c, nextConcept, true);

  const grade = (ok) => {
    graded = true;
    st.n += 1;
    st.streak = ok ? st.streak + 1 : 0;
    const justLearned = !st.done && st.streak >= STREAK;
    if (st.streak >= STREAK) st.done = true;
    const all = store.get(); all[c.id] = st; store.set(all);
    box.classList.add(ok ? "is-good" : "is-bad");
    box.querySelector(".streak").outerHTML = streakDots(st);
    optBtns.forEach((b, j) => {
      b.disabled = true;
      if (options[j].trim() === prob.a.trim()) b.classList.add("right");
      else if (j === chosen) b.classList.add("wrong");
    });
    box.querySelectorAll("form input, form button, #check").forEach((el) => (el.disabled = true));
    box.querySelector(".hint")?.remove();
    if (checkBtn) checkBtn.closest(".p-act").remove();
    box.querySelector("form button")?.remove();
    document.getElementById("fb").innerHTML = `<div class="fb ${ok ? "good" : "bad"}">
        <p class="fb-t"><span class="fb-i" aria-hidden="true">${ok ? "✓" : "✕"}</span>
        ${ok ? (justLearned ? "Correct. Concept learned." : "Correct.") : `Not quite. The answer is ${unwrapP(render(prob.a))}.`}</p>
        ${prob.e ? `<div class="ex prose">${render(prob.e)}</div>` : ""}
      </div>
      <div class="p-act">
        <button class="btn primary" id="again">Next problem</button>
        ${st.done && nextConcept ? `<a class="btn" href="#/d/${domainId}/${nextConcept.id}">Next concept</a>` : ""}
        <span class="hint" data-keys>Enter for the next problem</span>
      </div>`;
    const again = document.getElementById("again");
    again.onclick = nextProblem;
    again.focus({ preventScroll: true });
    document.getElementById("fb").scrollIntoView({ block: "nearest", behavior: matchMedia("(prefers-reduced-motion: reduce)").matches ? "auto" : "smooth" });
    // Refresh the sidebar's state for this concept.
    const link = document.querySelector(`#side a[href="#/d/${domainId}/${c.id}"]`);
    if (link && st.done) link.classList.add("done");
  };

  optBtns.forEach((b) => (b.onclick = () => select(+b.dataset.o)));
  if (checkBtn) checkBtn.onclick = () => chosen !== null && grade(options[chosen].trim() === prob.a.trim());
  const form = box.querySelector("form");
  if (form) form.onsubmit = (ev) => {
    ev.preventDefault();
    const v = form.querySelector("input").value;
    if (v.trim() && !graded) grade(numericMatch(prob.a, v));
  };

  keyHandler = (e) => {
    if (e.metaKey || e.ctrlKey || e.altKey || !document.body.contains(box)) return;
    const t = e.target, typing = t.tagName === "INPUT" || t.tagName === "TEXTAREA" || t.tagName === "SELECT";
    if (graded) {
      if (e.key === "Enter" && t.tagName !== "A") { e.preventDefault(); nextProblem(); }
      return;
    }
    if (!options.length || typing) return;
    let k = -1;
    if (/^[1-9]$/.test(e.key)) k = +e.key - 1;
    else if (/^[a-e]$/i.test(e.key)) k = e.key.toUpperCase().charCodeAt(0) - 65;
    if (k >= 0 && k < options.length) { e.preventDefault(); select(k); optBtns[k].focus({ preventScroll: true }); return; }
    if (e.key === "Enter") {
      if (t.dataset && t.dataset.o !== undefined) select(+t.dataset.o);
      if (chosen !== null) { e.preventDefault(); grade(options[chosen].trim() === prob.a.trim()); }
    }
  };
  if (focusBox) {
    (form ? form.querySelector("input") : optBtns[0])?.focus({ preventScroll: true });
    box.scrollIntoView({ block: "nearest" });
  }
}

/* ---------- Glossary ---------- */

// Every domain, subdomain and concept with its definition and key ideas, filtered as you type.
async function glossary() {
  app.innerHTML = `<div class="page"><p class="muted loading">Loading glossary…</p></div>`;
  const g = await load("data/glossary.json");
  const total = g.domains.reduce((n, d) => n + d.subdomains.reduce((m, s) => m + s.concepts.length, 0), 0);
  app.innerHTML = `<div class="page glossary">
    <header class="page-head"><h1>Glossary</h1>
      <p class="lede">Every concept with its definition, in layer order.</p></header>
    <div class="gbar" role="search">
      <input type="search" class="field" id="gq" placeholder="Search ${total} concepts" aria-label="Search the glossary" autocomplete="off" spellcheck="false">
      <select class="field" id="gj" aria-label="Jump to a domain"><option value="">Jump to…</option>
        ${g.domains.map((d) => `<option value="g-${d.id}">${esc(d.name)}</option>`).join("")}</select>
      <p class="gn" id="gn" aria-live="polite"></p>
    </div>
    <div id="gl">${g.domains.map((d) => `<section class="gdom" id="g-${d.id}" aria-labelledby="gh-${d.id}">
      <h2 id="gh-${d.id}">${esc(d.name)} <span class="eyebrow">Layer ${d.layer}</span></h2>
      ${d.covers.length ? `<details class="covers"><summary>Topics covered</summary><ul>${d.covers.map((c) => `<li>${esc(c)}</li>`).join("")}</ul></details>` : ""}
      ${d.subdomains.map((s) => `<div class="gsub"><h3>${esc(s.title)}</h3><dl>
        ${s.concepts.map((c) => `<div class="gc" data-text="${esc([c.label, c.definition, ...c.ideas, ...c.terms].join(" ").toLowerCase())}">
          <dt><a href="#/d/${d.id}/${c.id}">${esc(c.label)}</a></dt>
          <dd><p class="g-def">${unwrapP(render(c.definition))}</p>
          ${c.ideas.length ? `<p class="g-ideas">${c.ideas.map(esc).join('<span aria-hidden="true"> · </span>')}</p>` : ""}</dd>
        </div>`).join("")}</dl></div>`).join("")}
    </section>`).join("")}</div>
    <p class="empty muted" id="gempty" hidden>No concepts match.</p></div>`;
  const q = document.getElementById("gq"), count = document.getElementById("gn");
  const filter = () => {
    const v = q.value.trim().toLowerCase();
    let shown = 0;
    document.querySelectorAll(".gc").forEach((el) => {
      const on = !v || el.dataset.text.includes(v);
      el.hidden = !on; shown += on;
    });
    document.querySelectorAll(".gsub, .gdom").forEach((el) => { el.hidden = !el.querySelector(".gc:not([hidden])"); });
    count.textContent = v ? `${shown} of ${total}` : "";
    document.getElementById("gempty").hidden = shown > 0;
  };
  q.oninput = filter;
  document.getElementById("gj").onchange = (ev) => {
    const v = ev.target.value;
    if (!v) return;
    if (q.value) { q.value = ""; filter(); }
    document.getElementById(v).scrollIntoView();
    ev.target.value = "";
  };
}

/* ---------- Router ---------- */

async function route() {
  keyHandler = null;
  closeDrawer(false);
  index ||= await load("data/index.json");
  const [, page, a, b] = location.hash.split("/");
  const navKey = page === "d" || !page ? "home" : page;
  document.querySelectorAll("[data-nav]").forEach((el) =>
    el.dataset.nav === navKey ? el.setAttribute("aria-current", "page") : el.removeAttribute("aria-current"));
  document.body.dataset.page = page === "d" ? "study" : navKey;
  if (page === "d") await domain(a, b);
  else if (page === "glossary") await glossary();
  else if (page === "tracks") tracks();
  else if (page === "about") about();
  else if (page === "notes") app.innerHTML = "";
  else home();
  document.dispatchEvent(new CustomEvent("eda:render", { detail: window.__edaCtx || { page: navKey } }));
  window.__edaCtx = null;
  const h1 = app.querySelector("h1");
  document.title = page && h1 ? `${h1.textContent.trim()} · EDA-Game` : "EDA-Game";
  window.scrollTo(0, 0);
  app.focus({ preventScroll: true });
}
window.addEventListener("hashchange", route);
route();
