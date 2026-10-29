#!/usr/bin/env python3
"""
eda_check.py, the cross-domain gate for one or more EDA domains.

The Go contract test (edakit.CheckDomain) proves a domain's generators, lessons and shards in
isolation. This script proves the domain's place in the whole graph, which the Go test cannot
see while other domains are still being written:

  - every id carries the domain prefix and every subdomain the domain name
  - every prerequisite, key prerequisite and encompassed id resolves to the domain itself, the
    mathematics substrate, an existing concept of an upstream domain, or an anchor of an upstream
    domain (docs/eda/curriculum.json); frontier concepts may also use agent_first_eda anchors
  - agent_first_eda never depends on a frontier concept
  - every anchor the curriculum assigns to the domain exists
  - the domain is acyclic, every encompassed id is an ancestor, and the density floors hold
  - the frontier subdomain has 4 to 6 concepts
  - each concept has a lesson and a shard

Usage: python3 scripts/eda_check.py timing_analysis [more domains ...] | --all
"""
import json
import os
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))


def load(path):
    with open(path, encoding="utf-8") as f:
        return json.load(f)


def main(argv):
    cur = load(os.path.join(ROOT, "docs", "eda", "curriculum.json"))
    doms = {d["domain"]: d for d in cur["domains"]}
    by_key = {d["key"]: d for d in cur["domains"]}
    math_ids = {}
    eda_ids = {}
    cdir = os.path.join(ROOT, "data", "concepts")
    for name in os.listdir(cdir):
        if not name.endswith(".json") or name == "enrichment.json":
            continue
        rows = load(os.path.join(cdir, name))
        for c in rows:
            if c["domain"] in cur["math_domains"]:
                math_ids[c["id"]] = c["domain"]
            elif c["domain"] in doms:
                eda_ids[c["id"]] = c["domain"]
    targets = list(doms) if argv == ["--all"] else argv
    if not targets:
        print(__doc__)
        return 2
    agentic_anchors = set(by_key["agentic"]["anchors"])
    failed = 0
    for name in targets:
        if name not in doms:
            print(f"FAIL {name}: not a curriculum domain")
            failed += 1
            continue
        errs = check(name, doms, by_key, math_ids, eda_ids, agentic_anchors)
        if errs:
            failed += 1
            print(f"FAIL {name}: {len(errs)} issue(s)")
            for e in errs[:60]:
                print("   ", e)
            if len(errs) > 60:
                print(f"    ... and {len(errs) - 60} more")
        else:
            print(f"OK   {name}")
    return 1 if failed else 0


def check(name, doms, by_key, math_ids, eda_ids, agentic_anchors):
    spec = doms[name]
    errs = []
    path = os.path.join(ROOT, "data", "concepts", name + ".json")
    if not os.path.exists(path):
        return [f"missing {os.path.relpath(path, ROOT)}"]
    rows = load(path)
    own = {c["id"]: c for c in rows}
    if len(own) != len(rows):
        errs.append("duplicate ids inside the domain file")
    upstream = set(spec["upstream"])
    upstream_doms = {by_key[u]["domain"] for u in upstream if u in by_key}
    upstream_anchors = set()
    for u in upstream:
        if u in by_key:
            upstream_anchors |= set(by_key[u]["anchors"])

    def resolves(ref, frontier):
        if ref in own or ref in math_ids:
            return True
        if ref in upstream_anchors:
            return True
        if eda_ids.get(ref) in upstream_doms and ".frontier." not in ref:
            return True
        if frontier and ref in agentic_anchors:
            return True
        return False

    frontier_count = 0
    for c in rows:
        cid = c["id"]
        if not cid.startswith(spec["prefix"] + "."):
            errs.append(f"{cid}: id must start with {spec['prefix']}.")
        if c.get("domain") != name:
            errs.append(f"{cid}: domain {c.get('domain')!r}")
        sub = c.get("subdomain", "")
        if not sub.startswith(name + "."):
            errs.append(f"{cid}: subdomain {sub!r} must start with {name}.")
        frontier = sub == name + ".frontier"
        if frontier:
            frontier_count += 1
            if ".frontier." not in cid:
                errs.append(f"{cid}: frontier concepts use ids {spec['prefix']}.frontier.<name>")
        pre = c.get("prerequisites") or []
        enc = c.get("encompasses") or []
        key = c.get("key_prerequisites") or []
        if not pre:
            errs.append(f"{cid}: no prerequisites")
        if not key:
            errs.append(f"{cid}: no key_prerequisites")
        for k in key:
            if k not in pre:
                errs.append(f"{cid}: key_prerequisite {k} is not a prerequisite")
        for ref in pre + enc:
            if not resolves(ref, frontier):
                errs.append(f"{cid}: reference {ref} does not resolve (own domain, math substrate, upstream anchor"
                            f"{', agentic anchor' if frontier else ''})")
            if name == "agent_first_eda" and ".frontier." in ref:
                errs.append(f"{cid}: agent_first_eda must not depend on frontier concept {ref}")
        for f_ in ("teaching", "kp"):
            ext = ".md" if f_ == "teaching" else ".json"
            if not os.path.exists(os.path.join(ROOT, "data", "lessons", f_, cid + ext)):
                errs.append(f"{cid}: missing data/lessons/{f_}/{cid}{ext}")

    for a in spec["anchors"]:
        if a not in own:
            errs.append(f"anchor {a} ({spec['anchors'][a]}) is not defined")
    if spec["key"] != "agentic" and not (4 <= frontier_count <= 6):
        errs.append(f"frontier subdomain {name}.frontier has {frontier_count} concepts, want 4 to 6")
    lo, hi = int(spec["target"] * 0.85), int(spec["target"] * 1.2) + 1
    if not (lo <= len(rows) <= hi):
        errs.append(f"{len(rows)} concepts, target {spec['target']} (accept {lo}-{hi})")

    # Acyclic inside the domain, and encompassed ids must be ancestors.
    edges = {cid: [p for p in (c.get("prerequisites") or []) + (c.get("encompasses") or []) if p in own]
             for cid, c in own.items()}
    state = {}

    def visit(n, stack):
        state[n] = 1
        for m in edges[n]:
            if state.get(m) == 1:
                errs.append(f"cycle through {n} -> {m}")
            elif m not in state:
                visit(m, stack)
        state[n] = 2

    for n in own:
        if n not in state:
            visit(n, [])

    memo = {}

    def ancestors(n):
        if n in memo:
            return memo[n]
        memo[n] = set()
        out = set()
        for p in own[n].get("prerequisites") or []:
            out.add(p)
            if p in own:
                out |= ancestors(p)
        memo[n] = out
        return out

    with_enc = 0
    for cid, c in own.items():
        enc = c.get("encompasses") or []
        if enc:
            with_enc += 1
        anc = ancestors(cid)
        for e in enc:
            if e not in anc:
                errs.append(f"{cid}: encompasses {e}, which is not one of its prerequisites or their ancestors")
    if rows and with_enc / len(rows) < 0.55:
        errs.append(f"only {with_enc}/{len(rows)} concepts declare encompasses, want at least 55%")
    indeg = sum(len(c.get("prerequisites") or []) + len(c.get("encompasses") or []) for c in rows)
    if rows and indeg / len(rows) < 2.3:
        errs.append(f"mean in-degree {indeg / len(rows):.2f}, want at least 2.3 (prerequisites + encompasses)")
    return errs


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
