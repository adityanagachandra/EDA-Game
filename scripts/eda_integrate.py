#!/usr/bin/env python3
"""
eda_integrate.py, wire the authored EDA domains into the shipped corpus.

Idempotent. For every EDA domain file present in data/concepts/ it:

  1. maps each concept to its lesson in data/lessons/lessons.json (source teaching/<id>.md)
  2. writes one course per domain into data/courses.json (course id = domain name), placed
     ahead of the mathematics courses, plus a "Mathematics for EDA" course whose targets are the
     substrate concepts the EDA graph actually depends on
  3. writes the EDA learning tracks into data/destinations.json, ahead of the mathematics ones

Run after the domains pass scripts/eda_check.py: python3 scripts/eda_integrate.py
"""
import json
import os

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))


def load(p):
    with open(os.path.join(ROOT, p), encoding="utf-8") as f:
        return json.load(f)


def dump(p, data):
    with open(os.path.join(ROOT, p), "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2, ensure_ascii=False)
        f.write("\n")


COURSE_NAMES = {
    "eda_landscape": ("The EDA Landscape", "What electronic design automation is, the RTL-to-GDSII flow stage by stage, the representations each stage consumes and produces, PPA, the industry and its history."),
    "boolean_reasoning": ("Boolean Reasoning and Complexity", "Boolean functions, cofactors, Shannon expansion, covers and clauses, and the complexity theory that explains why every EDA problem needs heuristics."),
    "eda_graph_algorithms": ("Graph Algorithms for EDA", "Netlists as hypergraphs, partitioning (KL, FM, multilevel, spectral), shortest paths and A*, flows, Steiner trees and the graph machinery under placement, routing and timing."),
    "computational_geometry": ("Computational Geometry for Layout", "Rectilinear geometry, polygon Booleans, plane sweep and spatial indexes: the kernels under DRC, extraction, fill and routing."),
    "optimization_for_eda": ("Optimization for EDA", "LP, ILP, convex and nonconvex optimization, Lagrangian relaxation, metaheuristics, dynamic programming and Bayesian optimization, each tied to the EDA problem it solves."),
    "numerical_methods_for_eda": ("Numerical Methods for EDA", "Sparse linear algebra, Newton's method, stiff ODE integration, Krylov solvers, PDE discretization and model order reduction for simulators and analyzers."),
    "variability_statistics": ("Variability, Yield and Reliability Statistics", "Monte Carlo and rare-event estimation, process variation, statistical timing mathematics, yield models and lifetime statistics."),
    "circuits_and_electromagnetics": ("Circuits and Electromagnetics", "Circuit laws, RC and RLC dynamics, transfer functions, transmission lines, S-parameters and the field physics behind parasitics and signal integrity."),
    "semiconductor_devices": ("Semiconductor Devices and Compact Models", "Band theory, carrier transport, junctions, MOSFET physics, FinFETs and nanosheets, reliability and the compact models EDA tools consume."),
    "digital_logic_design": ("Digital Logic Design", "Combinational and sequential logic, flip-flop timing, FSMs, adders and multipliers, pipelines and memories, for a mathematician entering hardware."),
    "fabrication_and_dtco": ("Fabrication and Design-Technology Co-optimization", "Process steps, lithography, OPC and multiple patterning, the interconnect stack, PDKs and design rules, DTCO and yield."),
    "cmos_circuits_and_cells": ("CMOS Circuits and Cell Libraries", "Static CMOS, logical effort, power, standard-cell layout and the Liberty abstraction every digital tool consumes."),
    "hdl_and_compilers": ("HDLs and Hardware Compilers", "Verilog, SystemVerilog and VHDL semantics, elaboration, ASTs, IRs, SSA and dataflow analysis, RTLIL, MLIR and CIRCT."),
    "logic_synthesis": ("Logic Synthesis", "Two-level and multi-level minimization, don't cares, AIGs and rewriting, FSM synthesis, retiming, e-graphs, Yosys and ABC."),
    "formal_verification": ("Formal Verification", "SAT and CDCL, BDDs, equivalence checking, model checking, BMC, k-induction, IC3, SMT and assertions."),
    "functional_verification": ("Functional Verification and Simulation", "Event-driven and compiled simulation, testbenches, constrained random, coverage, UVM, emulation and CDC."),
    "circuit_simulation": ("Circuit Simulation", "Modified nodal analysis, Newton with companion models, transient integration and step control, AC, PSS, Monte Carlo and model order reduction."),
    "technology_mapping": ("Technology Mapping and Sizing", "Tree and DAG covering, cuts, LUT mapping, Boolean matching, gate sizing and buffer insertion."),
    "high_level_synthesis": ("High-Level Synthesis", "CDFGs, scheduling, allocation and binding, pipelining and modulo scheduling, from algorithm to RTL."),
    "timing_analysis": ("Static Timing Analysis", "Timing graphs, delay models, clocks and constraints, OCV, AOCV, POCV and SSTA, incremental timing and timing closure."),
    "test_and_dft": ("Test and Design for Test", "Fault models, fault simulation, ATPG, scan and compression, BIST, JTAG, diagnosis and yield learning."),
    "floorplanning_and_placement": ("Floorplanning and Placement", "Floorplan representations, wirelength models, quadratic and electrostatic placement, legalization, detailed and timing-driven placement."),
    "analog_mixed_signal_rf": ("Analog, Mixed-Signal and RF Design Automation", "Small-signal design, op-amps, data converters and PLLs, sizing as optimization, matching and analog layout, RF and mixed-signal simulation."),
    "fpga_eda": ("FPGA CAD", "FPGA architecture, packing, annealing placement, PathFinder routing and open FPGA toolchains."),
    "clock_and_power": ("Clock Trees and Power Delivery", "Clock tree synthesis, zero-skew and useful skew, clock gating, power grids, IR drop, decap, EM and low-power intent."),
    "routing": ("Routing", "Maze and A* routing, pattern and global routing, negotiated congestion, layer assignment, detailed and channel routing."),
    "extraction_and_signoff": ("Extraction, Signal Integrity and Signoff", "Parasitic extraction, SPEF, RC reduction, crosstalk and noise, ECOs and the signoff checklist."),
    "physical_verification_dfm": ("Physical Verification and DFM", "GDSII and OASIS, DRC engines, LVS as graph matching, antenna and fill, hotspots and tape-out."),
    "pcb_design": ("PCB Design", "Schematic to layout, stackups and controlled impedance, differential pairs, vias, manufacturing outputs and board-level SI."),
    "packaging_and_integration": ("Packaging, Chiplets and 3D Integration", "Flip chip and BGA, interposers, TSVs, hybrid bonding, UCIe, thermal, mechanical and reliability physics."),
    "eda_software_infrastructure": ("Building EDA Software", "Design databases, format parsers, incremental computation, deterministic parallelism, regression and QoR, and the open-source RTL-to-GDS flow hands-on."),
    "ml_for_eda": ("Machine Learning for EDA", "GNNs on netlists, congestion and timing prediction, RL placement and its debate, surrogates and Bayesian optimization, LLMs for RTL."),
    "agent_first_eda": ("Agent-First EDA Architecture", "Engines built for agents: the four-layer model, skills and harness loops, change tiers and seams, golden tests, the skill flywheel, tool-grounded RL and agentic signoff."),
}

TRACKS = [
    ("eda-tool-designer", "EDA Tool Designer: the full ramp",
     "Every EDA domain in dependency order, from the flow and its mathematics to agent-first engine architecture. The aggressive path to designing modern EDA tools.",
     None),
    ("eda-physical-design", "Physical Design Engines",
     "Partitioning, placement, clock trees, routing, extraction and physical verification, with the timing and optimization they rest on.",
     ["eda_graph_algorithms", "computational_geometry", "optimization_for_eda", "timing_analysis", "floorplanning_and_placement",
      "clock_and_power", "routing", "extraction_and_signoff", "physical_verification_dfm"]),
    ("eda-synthesis-verification", "Synthesis and Verification Engines",
     "Boolean reasoning, HDL compilers, logic synthesis, technology mapping, HLS, simulation, formal verification and test.",
     ["boolean_reasoning", "hdl_and_compilers", "logic_synthesis", "technology_mapping", "high_level_synthesis",
      "functional_verification", "formal_verification", "test_and_dft"]),
    ("eda-circuits-devices", "Circuits, Devices and Simulation",
     "Electromagnetics, devices, fabrication, CMOS cells, circuit simulation and analog design automation.",
     ["circuits_and_electromagnetics", "numerical_methods_for_eda", "semiconductor_devices", "fabrication_and_dtco",
      "cmos_circuits_and_cells", "circuit_simulation", "analog_mixed_signal_rf"]),
    ("eda-agent-first", "Agent-First EDA Architect",
     "The landscape, the software infrastructure of real tools, machine learning for EDA, and agent-first engine architecture.",
     ["eda_landscape", "eda_software_infrastructure", "ml_for_eda", "agent_first_eda"]),
    ("eda-systems", "Boards, Packages and FPGAs",
     "FPGA CAD, PCB design, packaging, chiplets and 3D integration.",
     ["fpga_eda", "pcb_design", "packaging_and_integration"]),
]


def main():
    cur = load("docs/eda/curriculum.json")
    order = sorted(cur["domains"], key=lambda d: d["layer"])
    present = []
    concepts = {}
    for d in order:
        p = os.path.join("data", "concepts", d["domain"] + ".json")
        if os.path.exists(os.path.join(ROOT, p)):
            present.append(d)
            concepts[d["domain"]] = load(p)

    # 1. lessons.json
    lessons = load("data/lessons/lessons.json")
    mapped = {e["concept_id"] for e in lessons}
    added = 0
    for d in present:
        for c in concepts[d["domain"]]:
            if c["id"] not in mapped and os.path.exists(os.path.join(ROOT, "data/lessons/teaching", c["id"] + ".md")):
                lessons.append({"concept_id": c["id"], "source": "teaching/" + c["id"] + ".md"})
                mapped.add(c["id"])
                added += 1
    # Drop only mappings whose concept exists nowhere in the corpus. Filtering by EDA prefix
    # would be wrong: variability_statistics shares the "stat." prefix with the mathematics
    # statistics domain, and a prefix filter deleted that domain's lesson mappings.
    all_ids = set()
    for name in os.listdir(os.path.join(ROOT, "data/concepts")):
        if name.endswith(".json") and name != "enrichment.json":
            all_ids |= {c["id"] for c in load(os.path.join("data/concepts", name))}
    lessons = [e for e in lessons if e["concept_id"] in all_ids]
    lessons.sort(key=lambda e: e["concept_id"])
    dump("data/lessons/lessons.json", lessons)

    # 2. courses.json
    courses = [c for c in load("data/courses.json") if c.get("grade") != "EDA"]
    math_ids = set()
    for name in os.listdir(os.path.join(ROOT, "data/concepts")):
        if name.endswith(".json") and name != "enrichment.json":
            for c in load(os.path.join("data/concepts", name)):
                if c["domain"] in cur["math_domains"]:
                    math_ids.add(c["id"])
    substrate = sorted({p for d in present for c in concepts[d["domain"]]
                        for p in c.get("prerequisites", []) if p in math_ids})
    eda_courses = [{
        "id": "eda_math",
        "name": "Mathematics for EDA",
        "grade": "EDA",
        "description": "The mathematics the EDA curriculum builds on: discrete mathematics and logic, graph theory, linear algebra, numerical analysis, probability, ODEs and optimization, as far as the EDA graph needs them.",
        "targets": substrate,
    }]
    for d in present:
        name, desc = COURSE_NAMES[d["domain"]]
        eda_courses.append({"id": d["domain"], "name": name, "grade": "EDA", "description": desc,
                            "targets": [c["id"] for c in concepts[d["domain"]]]})
    dump("data/courses.json", eda_courses + courses)

    # 3. destinations.json
    have = {d["domain"] for d in present}
    dests = [x for x in load("data/destinations.json") if not x["id"].startswith("eda-")]
    eda_dests = []
    for tid, name, desc, members in TRACKS:
        ids = [d["domain"] for d in order] if members is None else members
        ids = [i for i in ids if i in have]
        if ids:
            eda_dests.append({"id": tid, "name": name, "description": desc, "courses": ids})
    dump("data/destinations.json", eda_dests + dests)

    total = sum(len(v) for v in concepts.values())
    print(f"OK: {len(present)} EDA domains, {total} concepts, {added} lessons mapped, "
          f"{len(eda_courses)} courses, {len(eda_dests)} tracks, {len(substrate)} substrate targets")


if __name__ == "__main__":
    main()
