#!/usr/bin/env python3
"""Build the EDA curriculum glossary.

Writes docs/eda/GLOSSARY.md and docs/eda/site_glossary.json from:

  docs/eda/curriculum.json          domains, layers, upstream domains, covers
  data/courses.json                 domain display names and descriptions
  data/concepts/<domain>.json       concept labels and subdomains
  data/lessons/teaching/<id>.md     lesson text (definition, key ideas, terms)
  go run ./cmd/sitedump -per 1      one sample problem per concept

Usage:
  python3 scripts/build_glossary.py                  # runs sitedump itself
  python3 scripts/build_glossary.py --problems f.json  # reuse a saved dump

The glossary must not name people, companies or customers. Lesson sentences
that attribute work to a person or name a vendor are rewritten through
REWRITES or dropped, and the build fails if a DENY pattern survives anywhere
in the output. Eponymous technique names (Lee's algorithm, Elmore delay) are
terms, not attributions, and are kept.
"""

import argparse
import json
import re
import subprocess
import sys
from collections import defaultdict
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
OUT_MD = ROOT / "docs/eda/GLOSSARY.md"
OUT_JSON = ROOT / "docs/eda/site_glossary.json"

# Readable names for each subdomain key, per domain.
SUBDOMAINS = {
    "eda_landscape": {"overview": "What EDA is", "stages": "Stages of the flow", "repr": "Representations and formats",
                      "metrics": "Quality metrics", "branches": "Branches of EDA", "industry": "The industry and its standards",
                      "history": "History"},
    "boolean_reasoning": {"func": "Boolean functions", "decomp": "Decomposition", "cover": "Covering and minimization",
                          "repr": "Representations (BDDs, AIGs)", "order": "Variable ordering", "cnf": "CNF and encodings",
                          "algo": "Algorithms", "complexity": "Complexity"},
    "eda_graph_algorithms": {"model": "Graph models of circuits", "dag": "DAG algorithms", "path": "Shortest paths",
                             "partition": "Partitioning", "spectral": "Spectral methods", "flow": "Network flow",
                             "steiner": "Steiner trees", "coloring": "Graph coloring"},
    "computational_geometry": {"rect": "Rectangles and rectilinear shapes", "poly": "Polygons", "sweep": "Sweep-line algorithms",
                               "index": "Spatial indexing", "prox": "Proximity queries", "decomp": "Polygon decomposition",
                               "layout": "Layout geometry"},
    "optimization_for_eda": {"formulate": "Formulating problems", "lp": "Linear programming", "ip": "Integer programming",
                             "convex": "Convex optimization", "gradient": "Gradient methods", "lagrange": "Lagrangian methods",
                             "meta": "Metaheuristics", "dp": "Dynamic programming", "multi": "Multi-objective optimization",
                             "bayes": "Bayesian optimization"},
    "numerical_methods_for_eda": {"sparse": "Sparse matrices", "direct": "Direct solvers", "nonlinear": "Nonlinear equations",
                                  "ode": "ODE integration", "iter": "Iterative solvers", "pde": "PDE discretization",
                                  "mor": "Model order reduction", "hpc": "High-performance computing"},
    "variability_statistics": {"var": "Sources of variation", "gauss": "Gaussian models", "mc": "Monte Carlo",
                               "sens": "Sensitivity analysis", "fit": "Model fitting", "yield": "Yield", "rel": "Reliability"},
    "circuits_and_electromagnetics": {"circuit": "Circuit laws", "ac": "AC analysis", "signal": "Signals",
                                      "field": "Fields", "tline": "Transmission lines", "emc": "EMC and coupling",
                                      "solver": "Field solvers"},
    "semiconductor_devices": {"physics": "Semiconductor physics", "junction": "PN junctions", "mos": "MOS transistors",
                              "sce": "Short-channel effects", "advanced": "Advanced devices", "var": "Device variability",
                              "rel": "Device reliability", "scaling": "Scaling", "model": "Compact models"},
    "digital_logic_design": {"number": "Number systems", "comb": "Combinational logic", "seq": "Sequential logic",
                             "fsm": "Finite state machines", "arith": "Arithmetic circuits", "mem": "Memories",
                             "arch": "Digital architecture"},
    "fabrication_and_dtco": {"process": "Process steps", "litho": "Lithography", "beol": "Back end of line",
                             "pdk": "Process design kits", "yield": "Manufacturing yield", "dtco": "Design-technology co-optimization"},
    "cmos_circuits_and_cells": {"gate": "CMOS gates", "delay": "Gate delay", "power": "Power", "seq": "Sequential cells",
                                "cell": "Standard-cell design", "lib": "Cell libraries", "special": "Special cells"},
    "hdl_and_compilers": {"rtl": "RTL concepts", "verilog": "Verilog", "sv": "SystemVerilog", "lang": "Other hardware languages",
                          "compile": "Compilation", "elab": "Elaboration", "opt": "Compiler optimization",
                          "ir": "Intermediate representations"},
    "logic_synthesis": {"overview": "Overview", "twolevel": "Two-level minimization", "multilevel": "Multi-level synthesis",
                        "dc": "Don't-cares", "aig": "AIG optimization", "fsm": "FSM synthesis", "seq": "Sequential optimization",
                        "datapath": "Datapath synthesis", "tool": "Synthesis tools"},
    "formal_verification": {"sat": "SAT solving", "bdd": "BDD-based verification", "equiv": "Equivalence checking",
                            "mc": "Model checking", "bmc": "Bounded model checking", "induction": "Induction",
                            "ic3": "IC3 / PDR", "smt": "SMT", "sva": "Assertions (SVA)"},
    "functional_verification": {"event": "Event-driven simulation", "cycle": "Cycle-based simulation", "tb": "Testbenches",
                                "random": "Constrained-random stimulus", "cov": "Coverage", "assert": "Assertions",
                                "uvm": "UVM methodology", "cdc": "Clock-domain crossing", "accel": "Acceleration and emulation"},
    "circuit_simulation": {"tool": "Simulators and netlists", "mna": "Modified nodal analysis", "dc": "DC analysis",
                           "tran": "Transient analysis", "ac": "AC and noise analysis", "pss": "Periodic steady state",
                           "var": "Variation analysis", "mor": "Model order reduction", "fast": "Fast simulation"},
    "technology_mapping": {"subject": "Subject graphs", "tree": "Tree covering", "cut": "Cut enumeration",
                           "lut": "LUT mapping", "match": "Matching", "load": "Load-dependent mapping",
                           "size": "Gate sizing", "buffer": "Buffering", "tools": "Mapping tools"},
    "high_level_synthesis": {"ir": "HLS intermediate representations", "sched": "Scheduling", "bind": "Binding and allocation",
                             "pipe": "Pipelining", "mem": "Memory architecture", "sys": "HLS systems and tools"},
    "timing_analysis": {"graph": "The timing graph", "delay": "Delay calculation", "check": "Timing checks",
                        "clock": "Clocks", "constraint": "Constraints", "path": "Path analysis",
                        "var": "Variation-aware timing", "incremental": "Incremental timing", "closure": "Timing closure"},
    "test_and_dft": {"fault": "Fault models", "fsim": "Fault simulation", "atpg": "Test pattern generation",
                     "scan": "Scan design", "compress": "Test compression", "bist": "Logic BIST", "mbist": "Memory BIST",
                     "jtag": "Boundary scan (JTAG)", "ate": "Test equipment"},
    "floorplanning_and_placement": {"overview": "Overview", "format": "Placement formats", "wl": "Wirelength models",
                                    "floorplan": "Floorplanning", "macro": "Macro placement", "global": "Global placement",
                                    "legal": "Legalization", "detail": "Detailed placement", "driven": "Timing- and congestion-driven placement"},
    "analog_mixed_signal_rf": {"ss": "Small-signal analysis", "amp": "Amplifiers", "layout": "Analog layout",
                               "size": "Analog sizing", "conv": "Data converters", "pll": "PLLs", "rf": "RF design",
                               "ms": "Mixed-signal flows", "special": "Special topics"},
    "fpga_eda": {"arch": "FPGA architecture", "map": "Technology mapping to LUTs", "pack": "Packing", "place": "FPGA placement",
                 "route": "FPGA routing", "impl": "Implementation", "flow": "Toolchains and bitstreams"},
    "clock_and_power": {"objective": "Clock-tree objectives", "topology": "Clock-tree topologies", "zeroskew": "Zero-skew trees",
                        "buffered": "Buffered clock trees", "skew": "Useful skew", "gating": "Clock gating",
                        "pdn": "Power delivery networks", "lowpower": "Low-power design", "analysis": "Power and IR analysis"},
    "routing": {"grid": "Routing grids", "maze": "Maze routing", "pattern": "Pattern routing", "global": "Global routing",
                "layer": "Layer assignment", "detail": "Detailed routing", "channel": "Channel routing",
                "driven": "Timing-driven routing"},
    "extraction_and_signoff": {"resistance": "Resistance extraction", "capacitance": "Capacitance extraction",
                               "methods": "Extraction methods", "inductance": "Inductance", "format": "Parasitic formats",
                               "reduction": "Parasitic reduction", "si": "Signal integrity", "signoff": "Signoff"},
    "physical_verification_dfm": {"format": "Layout formats", "drc": "Design rule checking", "lvs": "Layout versus schematic",
                                  "antenna": "Antenna rules", "fill": "Metal fill", "dfm": "Design for manufacturing",
                                  "tapeout": "Tape-out"},
    "pcb_design": {"flow": "PCB flow and tools", "stackup": "Stackups", "via": "Vias", "layout": "PCB layout",
                   "diff": "Differential pairs", "si": "Signal integrity", "assembly": "Assembly",
                   "verify": "PCB verification", "output": "Manufacturing outputs"},
    "packaging_and_integration": {"basic": "Package basics", "wlp": "Wafer-level packaging", "adv": "Advanced packaging",
                                  "chiplet": "Chiplets", "thermal": "Thermal analysis", "mech": "Mechanical stress",
                                  "rel": "Package reliability", "eda3d": "EDA for 3D integration"},
    "eda_software_infrastructure": {"db": "Design databases", "api": "APIs and scripting", "parse": "Parsers",
                                    "arch": "Tool architecture", "incr": "Incremental computation", "par": "Parallelism",
                                    "qa": "Testing and QA", "pdk": "Open PDKs", "flow": "Flow orchestration"},
    "ml_for_eda": {"found": "Foundations", "gnn": "Graph neural networks", "predict": "Prediction models",
                   "rl": "Reinforcement learning", "surrogate": "Surrogate models", "diff": "Differentiable EDA",
                   "gen": "Generative models", "llm": "Language models for EDA", "eval": "Evaluation"},
    "agent_first_eda": {"architecture": "Architecture", "harness": "The harness", "tiers": "Tiers of autonomy",
                        "verification": "Verification of agent output", "safety": "Safety", "loops": "Agent loops",
                        "models": "Models", "learning": "Learning from runs", "orchestration": "Orchestration"},
}
FRONTIER = "Frontier: agent-first futures"

# Sentence rewrites that remove attributions to people and vendor names.
REWRITES = [
    ("Tools include Synopsys Design Compiler and Cadence Genus commercially, and Yosys with ABC in open source.",
     "Tools include commercial synthesizers and, in open source, Yosys with ABC."),
    ("Synopsys Design Constraints (SDC)", "SDC"),
    ("The Synopsys Design Constraints format, SDC,", "The SDC format"),
    ("Lee's maze-routing algorithm, from 1961, is", "Lee's maze-routing algorithm is"),
    ("The Kernighan-Lin algorithm of 1970 was", "The Kernighan-Lin algorithm was"),
    ("Hanan's theorem, from 1966, removes", "Hanan's theorem removes"),
    ("Iterated 1-Steiner, due to Kahng and Robins, improves", "Iterated 1-Steiner improves"),
    ("Pelgrom, Duinmaijer and Welbers showed in 1989 that the mismatch", "Pelgrom's law states that the mismatch"),
    ("In 1974 Dennard and colleagues showed that shrinking", "Dennard scaling: shrinking"),
    ("BSIM, the Berkeley Short-channel IGFET Model family, is", "BSIM, the short-channel IGFET model family, is"),
    ("Espresso, developed at IBM and Berkeley in the early 1980s, is", "Espresso is"),
    ("DAG-aware rewriting, introduced in ABC in 2006, counts", "DAG-aware rewriting, as implemented in ABC, counts"),
    ("SPICE, written at Berkeley in the early 1970s, set", "SPICE set"),
    ("ngspice descends from Berkeley SPICE3 with many extensions; Xyce was written at Sandia National Laboratories for large parallel simulations.",
     "ngspice descends from SPICE3 with many extensions; Xyce was written for large parallel simulations."),
    ("called pnjlim in Berkeley SPICE,", "called pnjlim in SPICE,"),
    ("ABC, from UC Berkeley, is", "ABC is"),
    ("Commercial tools such as AMD Vitis HLS use", "Commercial HLS tools use"),
    ("Clark's 1961 formulas give", "Clark's formulas give"),
    ("SCOAP, the Sandia Controllability/Observability Analysis Program, computes",
     "SCOAP (controllability/observability analysis) computes"),
    ("Open PDKs such as SkyWater SKY130 made", "Open PDKs such as SKY130 made"),
    ("Projects such as IceStorm (Lattice iCE40), Trellis (Lattice ECP5) and X-Ray (AMD Xilinx 7-series) recovered it",
     "Open documentation projects recovered it"),
    ("Each FPGA vendor ships its own implementation suite: AMD Vivado (with Vitis for software and HLS), Intel Quartus Prime, Lattice Radiant (and Diamond for older families) and Microchip Libero SoC.",
     "Each FPGA vendor ships its own implementation suite, covering synthesis, placement, routing, timing and bitstream generation."),
    ("Lee's algorithm (1961) is", "Lee's algorithm is"),
    ("Hadlock's router (1977) keeps", "Hadlock's router keeps"),
    ("Soukup's router (1978) takes", "Soukup's router takes"),
    ("Line-probe routers, by Mikami and Tabuchi (1968) and Hightower (1969), search", "Line-probe routers search"),
    ("The left-edge algorithm of Hashimoto and Stevens (1971) is", "The left-edge algorithm is"),
    ("Deutsch's dogleg router (1976) applies", "The dogleg router applies"),
    ("The Prim-Dijkstra construction of Alpert, Hu, Huang and Kahng blends", "The Prim-Dijkstra construction blends"),
    ("This lesson maps the main families: KiCad, LibrePCB and gEDA/Lepton in open source; Altium Designer, Cadence Allegro and OrCAD, Siemens Xpedition and PADS, Zuken, and Autodesk's EAGLE lineage commercially.",
     "This lesson maps the main families of PCB tools, open-source and commercial."),
    ("SKY130 is the SkyWater 130 nm process, released as an open process design kit with Google's backing.",
     "SKY130 is a 130 nm process released as an open process design kit."),
    ("GF180MCU from GlobalFoundries and SG13G2 from IHP widen", "GF180MCU and SG13G2 widen"),
    ("The RL macro placer published in Nature in 2021 became", "A published RL macro placer became"),
    ("By late 2026 the major EDA vendors and several start-ups sell or preview agentic products, and model providers have entered partnerships with EDA companies.",
     "EDA vendors and start-ups now sell or preview agentic products, often built on general-purpose language models."),
    # Domain "covers" lists in curriculum.json.
    ("Industry: Cadence, Synopsys, Siemens EDA (Mentor), Ansys, Keysight, Silvaco, Zuken, Altium; foundry model",
     "Industry: the major EDA vendors and their product lines; foundry model"),
    ("Major tools by role (Virtuoso, Innovus, Genus, Xcelium, Spectre, Design Compiler, Fusion Compiler, PrimeTime, VCS, Calibre, Questa, Tessent, Vivado, Quartus) - what class of problem each solves",
     "Major tools by role (custom layout, place and route, synthesis, logic and circuit simulation, timing signoff, physical verification, test, FPGA implementation) - what class of problem each solves"),
    ("landmark history: Mead-Conway, SPICE at Berkeley, Berkeley CAD tools, Brayton, Sangiovanni-Vincentelli, Newton, Pederson, Mead, Conway",
     "landmark history: structured VLSI design, SPICE, the university CAD tools behind today's logic synthesis and layout"),
    ("Standards bodies: IEEE, Accellera, Si2, IPC, SEMI", "Standards bodies for languages, formats, boards and manufacturing"),
    ("vendor flows: AMD Vivado/Vitis, Intel Quartus Prime, Lattice Radiant, Microchip Libero", "vendor flows"),
    ("Tools: KiCad, LibrePCB, gEDA/Lepton, Altium Designer, Cadence Allegro/OrCAD, Siemens Xpedition/PADS, Zuken, Autodesk EAGLE",
     "Tools: open-source (KiCad, LibrePCB, gEDA/Lepton) and commercial PCB suites"),
    ("Open PDKs: SkyWater SKY130, GF180MCU, IHP SG13G2", "Open PDKs: SKY130, GF180MCU, SG13G2"),
    # Concept labels.
    ("Mead-Conway, SPICE and the Berkeley CAD tools", "Structured VLSI design, SPICE and the university CAD tools"),
    ("GF180MCU and IHP SG13G2", "GF180MCU and SG13G2"),
    # Sample problems.
    ("continues the Berkeley MIS and SIS logic synthesis line", "continues the MIS and SIS logic synthesis line"),
    ("ABC, from Mishchenko and Brayton, is the successor to SIS.", "ABC is the successor to SIS."),
    ("Why do Vivado users save design checkpoints", "Why do FPGA tool users save design checkpoints"),
]

# Anything matching here must not reach the output. The build fails if it does.
DENY = re.compile(
    r"~/|/Users/"
    r"|\b(?:Synopsys|Cadence|Xilinx|Intel|Microchip|Siemens|Zuken|Autodesk|Altium|Google|GlobalFoundries|SkyWater"
    r"|Sandia|Berkeley|IBM|Accellera|Nvidia|NVIDIA|TSMC|Samsung|Apple|Mentor|Ansys|Keysight|OpenAI|Anthropic|Microsoft"
    r"|Vivado|Quartus|Libero|Virtuoso|Spectre|OrCAD|Xpedition|IHP)\b"
    r"|\bAMD (?:Vi|Xi)|\bLattice (?:iCE|ECP|Radiant|Diamond)|\bNature\b"
    r"|\b(?:Mead|Conway|Brayton|Sangiovanni|Vincentelli|Kahng|Robins|Duinmaijer|Welbers|Hightower|Hashimoto|Alpert)\b"
    r"|\bet al\b|\band colleagues\b|\(\s*(?:19|20)\d\d\s*\)"
)

# Bolded words that are list labels rather than terms (single capitalized words).
ACRONYM = re.compile(r"^[A-Z0-9][A-Z0-9/-]+s?$")


def scrub(text):
    for old, new in REWRITES:
        text = text.replace(old, new)
    return text


def clean_inline(text):
    """Strip markdown emphasis and links, keep math as written."""
    text = re.sub(r"\*\*(.+?)\*\*", r"\1", text)
    text = re.sub(r"(?<!\*)\*(?!\s)(.+?)(?<!\s)\*(?!\*)", r"\1", text)
    text = re.sub(r"\[([^\]]+)\]\([^)]+\)", r"\1", text)
    text = re.sub(r"!\[[^\]]*\]\([^)]+\)", "", text)
    return re.sub(r"\s+", " ", text).strip()


def sentences(par):
    # Split on sentence ends, but not inside $...$ math or after "e.g." / "i.e.".
    parts, buf, in_math = [], "", False
    i = 0
    while i < len(par):
        ch = par[i]
        buf += ch
        if ch == "$":
            in_math = not in_math
        elif ch in ".!?" and not in_math and i + 1 < len(par) and par[i + 1] == " ":
            tail = buf[-5:].lower()
            if not (tail.endswith("e.g.") or tail.endswith("i.e.") or tail.endswith("vs.") or tail.endswith("cf.")):
                parts.append(buf.strip())
                buf = ""
        i += 1
    if buf.strip():
        parts.append(buf.strip())
    return parts


def parse_lesson(path):
    text = path.read_text()
    text = re.split(r"^## (?:Going deeper|Further reading|References)\b", text, flags=re.M)[0]
    headings = [clean_inline(h) for h in re.findall(r"^#{2,3} (.+)$", text, flags=re.M)]
    bolds = [clean_inline(b) for b in re.findall(r"\*\*(.+?)\*\*", text)]
    body = re.sub(r"^#.*$", "", text, flags=re.M)
    paras = [p.strip() for p in body.split("\n\n") if p.strip()]
    prose = [p for p in paras if not re.match(r"^(?:[-*|>]|\d+\.|\$\$|!\[|<)", p)]
    # Lessons often open with a short hook ("Not every die on a wafer works."),
    # so keep adding clean sentences of the opening paragraph until the
    # definition is long enough to say what the concept is.
    definition = ""
    for par in prose[:2]:
        picked = []
        for s in sentences(clean_inline(scrub(par))):
            if DENY.search(s) or len(s) <= 20:
                continue
            picked.append(s)
            if len(" ".join(picked)) >= 80 or len(picked) == 2:
                break
        if picked:
            definition = " ".join(picked)
            break
    seen, ideas = set(), []
    for h in headings:
        h = scrub(h)
        if h.lower() not in seen and not DENY.search(h):
            seen.add(h.lower())
            ideas.append(h)
    terms = []
    for b in bolds:
        b = scrub(b).strip(" .:,;")
        if not b or DENY.search(b) or len(b.split()) > 6:
            continue
        if " " not in b and b[0].isupper() and not ACRONYM.match(b):
            continue
        if b.lower() not in {t.lower() for t in terms}:
            terms.append(b)
    return definition, ideas, terms


def load_problems(path):
    if path:
        return json.loads(Path(path).read_text())
    ids = []
    for d in json.loads((ROOT / "docs/eda/curriculum.json").read_text())["domains"]:
        ids += [c["id"] for c in json.loads((ROOT / f"data/concepts/{d['domain']}.json").read_text())]
    res = subprocess.run(["go", "run", "./cmd/sitedump", "-per", "1"], input="\n".join(ids) + "\n",
                         cwd=ROOT, capture_output=True, text=True, check=True)
    return json.loads(res.stdout)


def slug(text):
    s = re.sub(r"[^a-z0-9\s-]", "", text.lower())
    return re.sub(r"[\s]+", "-", s.strip())


def anchor(cid):
    return "c-" + cid.replace(".", "-").replace("_", "-")


def blockquote(text):
    return "\n".join("> " + line if line.strip() else ">" for line in text.strip().split("\n"))


def build(problems):
    cur = json.loads((ROOT / "docs/eda/curriculum.json").read_text())
    courses = {c["id"]: c for c in json.loads((ROOT / "data/courses.json").read_text())}
    key_to_name = {d["key"]: courses[d["domain"]]["name"] for d in cur["domains"]}
    domains = sorted(cur["domains"], key=lambda d: d["layer"])  # stable: keeps curriculum order within a layer

    data = {"domains": [], "index": []}
    index = defaultdict(set)  # term -> concept ids
    labels = {}
    for d in domains:
        course = courses[d["domain"]]
        concepts = json.loads((ROOT / f"data/concepts/{d['domain']}.json").read_text())
        subs = {}
        for c in concepts:
            subs.setdefault(c["subdomain"], []).append(c)
        names = SUBDOMAINS[d["domain"]]
        dom = {
            "key": d["key"], "id": d["domain"], "name": course["name"], "layer": d["layer"],
            "description": scrub(course.get("description", "")),
            "upstream": [key_to_name.get(u, u.replace("_", " ")) for u in d.get("upstream", [])],
            "covers": [scrub(x) for x in d.get("covers", [])],
            "subdomains": [],
        }
        for sub_id, cs in subs.items():
            key = sub_id.split(".")[-1]
            title = FRONTIER if key == "frontier" else names.get(key, key.replace("_", " ").capitalize())
            sub = {"id": sub_id, "title": title, "concepts": []}
            for c in cs:
                definition, ideas, terms = parse_lesson(ROOT / f"data/lessons/teaching/{c['id']}.md")
                # First generated problem whose text passes the deny list.
                sample = None
                for p in problems.get(c["id"]) or []:
                    cand = {"question": scrub(p["q"]), "answer": scrub(p["a"]), "explanation": scrub(p.get("e", ""))}
                    if not any(DENY.search(v) for v in cand.values()):
                        sample = cand
                        break
                label = scrub(c["label"])
                sub["concepts"].append({"id": c["id"], "label": label, "definition": definition,
                                        "key_ideas": ideas, "terms": terms, "sample": sample})
                labels[c["id"]] = label
                index[label].add(c["id"])
                for t in terms:
                    index[t].add(c["id"])
            dom["subdomains"].append(sub)
        data["domains"].append(dom)

    # Merge index keys that differ only in case; prefer the form that appears most.
    merged = {}
    for term, ids in index.items():
        k = term.lower()
        if k in merged:
            merged[k]["ids"] |= ids
            if term[0].islower():
                merged[k]["term"] = term
        else:
            merged[k] = {"term": term, "ids": set(ids)}
    order = lambda c: list(labels).index(c)
    data["index"] = sorted(({"term": v["term"], "concepts": sorted(v["ids"], key=order)} for v in merged.values()),
                           key=lambda e: (re.sub(r"[^a-z0-9]", "", e["term"].lower()) or e["term"].lower(), e["term"]))
    return data, labels


def render(data, labels):
    doms = data["domains"]
    nsub = sum(len(d["subdomains"]) for d in doms)
    ncon = sum(len(s["concepts"]) for d in doms for s in d["subdomains"])
    out = ["# EDA curriculum glossary", "",
           "Every domain, subdomain and concept in the electronic design automation curriculum, in layer order. "
           "Each concept has a short definition taken from the opening of its lesson, the key ideas and terms the lesson covers, "
           "and one sample problem with its answer. An alphabetical index of terms is at the end.", "",
           f"{len(doms)} domains, {nsub} subdomains, {ncon} concepts, {len(data['index'])} index terms. "
           "Generated by `scripts/build_glossary.py`; do not edit by hand.", "",
           "## Contents", ""]
    layer = None
    for d in doms:
        if d["layer"] != layer:
            layer = d["layer"]
            out.append(f"- **Layer {layer}**")
        n = sum(len(s["concepts"]) for s in d["subdomains"])
        out.append(f"  - [{d['name']}](#{slug(d['name'])}) ({n} concepts)")
    out += ["- [Index of terms](#index-of-terms)", ""]

    for d in doms:
        out += ["---", "", f"## {d['name']}", "",
                f"Layer {d['layer']}. {d['description']}", ""]
        if d["upstream"]:
            out += ["Builds on: " + ", ".join(d["upstream"]) + ".", ""]
        out += ["Topics covered:", ""] + [f"- {x}" for x in d["covers"]] + [""]
        out += ["Subdomains: " + ", ".join(f"[{s['title']}](#{slug(d['name'] + ' ' + s['title'])})"
                                           for s in d["subdomains"]), ""]
        for s in d["subdomains"]:
            # Explicit anchor: subdomain titles such as "Overview" repeat across domains.
            out += [f'<a id="{slug(d["name"] + " " + s["title"])}"></a>', "", f"### {s['title']}", ""]
            for c in s["concepts"]:
                out += [f'<a id="{anchor(c["id"])}"></a>', "", f"#### {c['label']}", "", f"`{c['id']}`", ""]
                if c["definition"]:
                    out += [c["definition"], ""]
                if c["key_ideas"]:
                    out += ["Key ideas and worked examples:", ""] + [f"- {x}" for x in c["key_ideas"]] + [""]
                if c["terms"]:
                    out += ["Terms: " + ", ".join(c["terms"]) + ".", ""]
                if c["sample"]:
                    sm = c["sample"]
                    out += ["Sample problem:", "", blockquote(sm["question"]), "",
                            f"Answer: {sm['answer'].strip()}"]
                    if sm["explanation"].strip():
                        out += ["", f"Why: {sm['explanation'].strip()}"]
                    out.append("")
    out += ["---", "", "## Index of terms", ""]
    letter = None
    for e in data["index"]:
        first = re.sub(r"[^A-Za-z0-9]", "", e["term"])[:1].upper() or "#"
        first = "0-9" if first.isdigit() else first
        if first != letter:
            letter = first
            out += ["", f"### {letter}", ""]
        links = ", ".join(f"[{labels[c]}](#{anchor(c)})" for c in e["concepts"])
        if len(e["concepts"]) == 1 and labels[e["concepts"][0]] == e["term"]:
            out.append(f"- [{e['term']}](#{anchor(e['concepts'][0])})")
        else:
            out.append(f"- {e['term']}: {links}")
    return "\n".join(out).rstrip() + "\n"


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--problems", help="saved sitedump JSON; runs sitedump when omitted")
    args = ap.parse_args()
    data, labels = build(load_problems(args.problems))
    md = render(data, labels)
    js = json.dumps(data, ensure_ascii=False, indent=1)
    bad = sorted({m.group(0) for m in DENY.finditer(md)} | {m.group(0) for m in DENY.finditer(js)})
    if bad:
        sys.exit(f"denied names in output: {bad}")
    OUT_MD.write_text(md)
    OUT_JSON.write_text(js + "\n")
    doms = data["domains"]
    print(f"{OUT_MD.relative_to(ROOT)}: {len(md.encode())} bytes; {OUT_JSON.relative_to(ROOT)}: {len(js.encode())} bytes")
    print(f"domains {len(doms)}, subdomains {sum(len(d['subdomains']) for d in doms)}, "
          f"concepts {sum(len(s['concepts']) for d in doms for s in d['subdomains'])}, index terms {len(data['index'])}")


if __name__ == "__main__":
    main()
