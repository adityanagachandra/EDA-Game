// EDA-Game static reader. Hash routes: #/  #/tracks  #/about  #/d/<domain>[/<concept>]
const app = document.getElementById("app");
const cache = {};
let index = null;

const store = {
  get() { try { return JSON.parse(localStorage.getItem("eda-game") || "{}"); } catch { return {}; } },
  set(v) { try { localStorage.setItem("eda-game", JSON.stringify(v)); } catch {} },
};
const STREAK = 3; // correct answers in a row to mark a concept learned here

async function load(path) {
  if (!cache[path]) cache[path] = fetch(path).then((r) => r.json());
  return cache[path];
}
const esc = (s) => String(s).replace(/[&<>"]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" }[c]));

// Markdown with $…$ and $$…$$ math: lift the math out, render the markdown, put KaTeX back.
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

function progress(domain) {
  const p = store.get();
  return Object.keys(p).filter((k) => p[k].done && p[k].domain === domain).length;
}

function home() {
  const layers = {};
  index.domains.forEach((d) => (layers[d.layer] ||= []).push(d));
  const done = Object.values(store.get()).filter((v) => v.done).length;
  app.innerHTML = `
  <section class="hero">
    <h1>From first principles to agent-first EDA tools.</h1>
    <p>${index.domains.length} domains and ${index.concepts} concepts, from the RTL-to-GDS flow and its mathematics to
    timing, placement, routing, signoff and the agent-first engine architectures that come next.
    Each concept has a short lesson and practice problems.</p>
    <div class="stats"><div><b>${index.concepts}</b><span>concepts</span></div>
    <div><b>${index.domains.length}</b><span>domains</span></div>
    <div><b>13</b><span>dependency layers</span></div>
    <div><b>${done}</b><span>learned on this device</span></div></div>
  </section>
  ${Object.keys(layers).sort((a, b) => a - b).map((l) => `
    <div class="layer"><h3>Layer ${l}</h3><div class="grid">
    ${layers[l].map((d) => {
      const n = progress(d.id);
      return `<a class="card" href="#/d/${d.id}"><h4>${esc(d.name)}</h4><p>${esc(d.description)}</p>
        <div class="bar"><i style="width:${(100 * n) / d.count}%"></i></div>
        <div class="meta">${n} / ${d.count} concepts</div></a>`;
    }).join("")}</div></div>`).join("")}`;
}

function tracks() {
  const name = Object.fromEntries(index.domains.map((d) => [d.id, d.name]));
  app.innerHTML = `<h1>Tracks</h1><p class="muted">Each track takes its domains in dependency order.</p>
  ${index.tracks.map((t) => `<div class="track"><h2>${esc(t.name)}</h2><p class="muted">${esc(t.description)}</p>
    <ol>${t.courses.filter((c) => name[c]).map((c) => `<li><a href="#/d/${c}">${esc(name[c])}</a></li>`).join("")}</ol></div>`).join("")}`;
}

function about() {
  app.innerHTML = `<div class="lesson">${render(`# About EDA-Game

EDA-Game teaches electronic design automation as a directed acyclic graph of concepts. Every domain sits on a
layer, and a concept depends only on concepts in its own domain, in lower layers, or in the mathematics substrate,
so the order you meet things in is always one where the prerequisites come first.

This page is the **static reader**: every lesson, plus a fixed set of practice problems per concept. A concept
counts as learned here after ${STREAK} correct answers in a row, stored only in this browser.

The **full adaptive engine** (fresh problems on every draw, spaced review, a placement diagnostic, mastery gating
over the whole graph) runs locally from the repository:

\`\`\`
go build -o bin/mathua ./cmd/mathua
DATABASE_URL=sqlite:eda.db ./bin/mathua --serve --port 8080
\`\`\`

Built on the open-source [mathua](https://github.com/chuma-beep/mathua) engine (MIT).`)}</div>`;
}

async function domain(id, conceptId) {
  const meta = index.domains.find((d) => d.id === id);
  if (!meta) return home();
  const concepts = await load(`data/${id}.json`);
  const c = concepts.find((x) => x.id === conceptId) || concepts[0];
  const p = store.get();
  let sub = "";
  const side = concepts.map((x) => {
    const head = x.sub !== sub ? `<h5>${esc((sub = x.sub).replace(/_/g, " "))}</h5>` : "";
    const cls = [x.id === c.id ? "on" : "", p[x.id]?.done ? "done" : ""].join(" ");
    return `${head}<a class="${cls}" href="#/d/${id}/${x.id}">${esc(x.label)}</a>`;
  }).join("");
  const domainOf = (pid) => index.domains.find((d) => pid.startsWith(d.prefix + "."));
  const pre = c.prereqs.map(([pid, label]) => {
    const owner = domainOf(pid);
    return owner ? `<a href="#/d/${owner.id}/${pid}">${esc(label)}</a>` : `<span title="mathematics prerequisite">${esc(label)}</span>`;
  }).join("");
  const i = concepts.indexOf(c);
  const next = concepts[i + 1];
  app.innerHTML = `<div class="split"><aside class="side"><a href="#/">← All domains</a>
    <h5 style="margin-top:12px">${esc(meta.name)}</h5>${side}</aside>
    <section><div class="lesson">${render(c.lesson || `# ${c.label}`)}
      ${pre ? `<p class="prereq">Builds on: ${pre}</p>` : ""}</div>
      <div class="practice" id="practice"></div>
      ${next ? `<p style="margin-top:18px"><a href="#/d/${id}/${next.id}">Next: ${esc(next.label)} →</a></p>` : ""}
    </section></div>`;
  practice(id, c);
  window.scrollTo(0, 0);
}

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

function practice(domainId, c) {
  const box = document.getElementById("practice");
  if (!c.problems.length) { box.innerHTML = `<p class="muted">No practice problems for this concept.</p>`; return; }
  const p = store.get();
  const st = p[c.id] || { domain: domainId, streak: 0, done: false, n: 0 };
  const prob = c.problems[st.n % c.problems.length];
  const { stem, options } = splitChoices(prob.q);
  box.innerHTML = `<div class="row" style="justify-content:space-between">
      <b>Practice</b><span class="pill">${st.done ? "Learned" : `Streak ${st.streak} / ${STREAK}`}</span></div>
    <div class="q">${render(stem)}</div>
    ${options.length ? `<div class="opts">${options.map((o, k) => `<button data-o="${k}">${String.fromCharCode(65 + k)}) ${render(o).replace(/^<p>|<\/p>\s*$/g, "")}</button>`).join("")}</div>`
      : `<form class="row"><input class="ans" placeholder="Your answer" autocomplete="off"><button class="btn primary">Check</button></form>`}
    <div id="fb"></div>`;
  const grade = (ok, given) => {
    st.n += 1;
    st.streak = ok ? st.streak + 1 : 0;
    if (st.streak >= STREAK) st.done = true;
    const all = store.get(); all[c.id] = st; store.set(all);
    document.getElementById("fb").innerHTML = `<div class="fb ${ok ? "good" : "bad"}">
      <b>${ok ? "Correct." : `Not quite. The answer is ${render(prob.a).replace(/^<p>|<\/p>\s*$/g, "")}.`}</b>
      ${prob.e ? `<div class="ex">${render(prob.e)}</div>` : ""}</div>
      <div class="row"><button class="btn primary" id="again">Next problem</button></div>`;
    box.querySelectorAll("button[data-o], form button").forEach((b) => (b.disabled = true));
    document.getElementById("again").onclick = () => practice(domainId, c);
  };
  box.querySelectorAll("button[data-o]").forEach((b) => b.onclick = () => grade(options[+b.dataset.o].trim() === prob.a.trim(), options[+b.dataset.o]));
  const form = box.querySelector("form");
  if (form) form.onsubmit = (ev) => { ev.preventDefault(); const v = form.querySelector("input").value; if (v.trim()) grade(numericMatch(prob.a, v), v); };
}

// Every domain, subdomain and concept with its definition and key ideas, filtered as you type.
async function glossary() {
  const g = await load("data/glossary.json");
  const total = g.domains.reduce((n, d) => n + d.subdomains.reduce((m, s) => m + s.concepts.length, 0), 0);
  app.innerHTML = `<h1>Glossary</h1>
    <p class="muted">${g.domains.length} domains, ${g.domains.reduce((n, d) => n + d.subdomains.length, 0)} subtopics and
    ${total} concepts in layer order, each with its definition and the ideas its lesson covers.</p>
    <input class="ans" id="gq" placeholder="Filter: slack, Elmore, BDD, Steiner…" style="width:100%;margin:8px 0 4px">
    <p class="meta" id="gn"></p>
    <nav class="toc">${g.domains.map((d) => `<a href="#/glossary" data-jump="g-${d.id}">${esc(d.name)}</a>`).join("")}</nav>
    <div id="gl">${g.domains.map((d) => `<section class="gdom" id="g-${d.id}">
      <h2>${esc(d.name)} <span class="pill">Layer ${d.layer}</span></h2>
      ${d.covers.length ? `<details><summary class="meta">Topics covered</summary><ul class="meta">${d.covers.map((c) => `<li>${esc(c)}</li>`).join("")}</ul></details>` : ""}
      ${d.subdomains.map((s) => `<div class="gsub"><h3>${esc(s.title)}</h3>
        ${s.concepts.map((c) => `<div class="gc" data-text="${esc([c.label, c.definition, ...c.ideas, ...c.terms].join(" ").toLowerCase())}">
          <a href="#/d/${d.id}/${c.id}"><b>${esc(c.label)}</b></a>
          <div>${render(c.definition).replace(/^<p>|<\/p>\s*$/g, "")}</div>
          ${c.ideas.length ? `<div class="meta">${c.ideas.map(esc).join(" · ")}</div>` : ""}
        </div>`).join("")}</div>`).join("")}
    </section>`).join("")}</div>`;
  const q = document.getElementById("gq"), count = document.getElementById("gn");
  const filter = () => {
    const v = q.value.trim().toLowerCase();
    let shown = 0;
    document.querySelectorAll(".gc").forEach((el) => {
      const on = !v || el.dataset.text.includes(v);
      el.hidden = !on; shown += on;
    });
    document.querySelectorAll(".gsub, .gdom").forEach((el) => { el.hidden = !el.querySelector(".gc:not([hidden])"); });
    count.textContent = v ? `${shown} matching concepts` : "";
  };
  q.oninput = filter;
  document.querySelectorAll("[data-jump]").forEach((a) => a.onclick = (ev) => {
    ev.preventDefault(); document.getElementById(a.dataset.jump).scrollIntoView();
  });
}

async function route() {
  index ||= await load("data/index.json");
  const [, page, a, b] = location.hash.split("/");
  if (page === "d") return domain(a, b);
  if (page === "glossary") return glossary();
  if (page === "tracks") return tracks();
  if (page === "about") return about();
  home();
}
window.addEventListener("hashchange", route);
route();
