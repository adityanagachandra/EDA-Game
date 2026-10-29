#!/usr/bin/env python3
"""Build the static EDA-Game reader into site/ for GitHub Pages.

GitHub Pages serves files, not the Go engine, so the reader carries every EDA
lesson and a fixed set of practice problems per concept (drawn by
cmd/sitedump). Mastery is tracked in the browser only. The full adaptive
engine still runs with ./bin/mathua --serve.

    python3 scripts/build_site.py          # writes site/
"""
import json
import pathlib
import shutil
import subprocess

ROOT = pathlib.Path(__file__).resolve().parent.parent
SITE = ROOT / "site"
TEMPLATE = ROOT / "scripts" / "site_template"


def load_concepts():
    labels, by_domain = {}, {}
    for f in sorted((ROOT / "data" / "concepts").glob("*.json")):
        data = json.loads(f.read_text())
        concepts = data if isinstance(data, list) else data["concepts"]
        for c in concepts:
            if "id" not in c or "domain" not in c:
                continue  # enrichment.json annotates concepts, it does not define them
            labels[c["id"]] = c.get("label", c["id"])
            by_domain.setdefault(c["domain"], []).append(c)
    return labels, by_domain


def topo(concepts):
    """Order a domain's concepts so in-domain prerequisites come first."""
    ids = {c["id"] for c in concepts}
    done, order = set(), []
    index = {c["id"]: c for c in concepts}

    def visit(cid, stack=()):
        if cid in done or cid in stack:
            return
        for p in index[cid].get("prerequisites", []):
            if p in ids:
                visit(p, stack + (cid,))
        done.add(cid)
        order.append(index[cid])

    for c in concepts:
        visit(c["id"])
    return order


def main():
    curriculum = json.loads((ROOT / "docs/eda/curriculum.json").read_text())
    courses = {c["id"]: c for c in json.loads((ROOT / "data/courses.json").read_text())}
    tracks = [t for t in json.loads((ROOT / "data/destinations.json").read_text())
              if t["id"].startswith("eda-")]
    labels, by_domain = load_concepts()

    domains = curriculum["domains"]
    ids = "\n".join(c["id"] for d in domains for c in by_domain[d["domain"]])
    dump = subprocess.run(["go", "run", "./cmd/sitedump"], input=ids, cwd=ROOT,
                          capture_output=True, text=True, check=True)
    problems = json.loads(dump.stdout)

    if SITE.exists():
        shutil.rmtree(SITE)
    shutil.copytree(TEMPLATE, SITE)
    (SITE / "data").mkdir()
    (SITE / ".nojekyll").write_text("")

    index_domains = []
    for d in domains:
        name = d["domain"]
        course = courses.get(name, {})
        out = []
        for c in topo(by_domain[name]):
            lesson = ROOT / "data/lessons/teaching" / f"{c['id']}.md"
            out.append({
                "id": c["id"],
                "label": c.get("label", c["id"]),
                "sub": c.get("subdomain", "").split(".")[-1],
                "type": c.get("grading_type", ""),
                "prereqs": [[p, labels.get(p, p)] for p in c.get("prerequisites", [])],
                "lesson": lesson.read_text() if lesson.exists() else "",
                "problems": problems.get(c["id"], []),
            })
        (SITE / "data" / f"{name}.json").write_text(json.dumps(out, separators=(",", ":")))
        index_domains.append({
            "id": name,
            "name": course.get("name", name),
            "description": course.get("description", ""),
            "prefix": d["prefix"],
            "layer": d["layer"],
            "upstream": d["upstream"],
            "count": len(out),
        })

    index = {
        "domains": index_domains,
        "tracks": [{"id": t["id"], "name": t["name"], "description": t.get("description", ""),
                    "courses": t["courses"]} for t in tracks],
        "concepts": sum(d["count"] for d in index_domains),
    }
    (SITE / "data" / "index.json").write_text(json.dumps(index, indent=1))
    write_glossary()
    bust_cache()
    print(f"site/: {len(index_domains)} domains, {index['concepts']} concepts, "
          f"{sum(len(v) for v in problems.values())} problems")


def bust_cache():
    """Version app.js and style.css by content hash, so returning visitors never mix old and new."""
    import hashlib
    page = SITE / "index.html"
    html = page.read_text()
    for name in ("app.js", "style.css"):
        tag = hashlib.sha1((SITE / name).read_bytes()).hexdigest()[:8]
        html = html.replace(f'"{name}"', f'"{name}?v={tag}"')
    page.write_text(html)


def write_glossary():
    """Slim docs/eda/site_glossary.json (from scripts/build_glossary.py) for the Glossary page.

    Sample problems are dropped, since each concept page already carries its practice set, and
    index terms without a letter ("0", "10n") are dropped as noise.
    """
    src = ROOT / "docs/eda/site_glossary.json"
    if not src.exists():
        return
    g = json.loads(src.read_text())
    domains = [{
        "id": d["id"], "name": d["name"], "layer": d["layer"], "covers": d.get("covers", []),
        "subdomains": [{
            "title": s["title"],
            "concepts": [{"id": c["id"], "label": c["label"], "definition": c.get("definition", ""),
                          "ideas": c.get("key_ideas", []), "terms": c.get("terms", [])}
                         for c in s["concepts"]],
        } for s in d["subdomains"]],
    } for d in g["domains"]]
    terms = [t for t in g["index"] if any(ch.isalpha() for ch in t["term"]) and len(t["term"]) > 2]
    (SITE / "data" / "glossary.json").write_text(
        json.dumps({"domains": domains, "index": terms}, separators=(",", ":")))


if __name__ == "__main__":
    main()
