# ATPG as a callable skill

In an agent-first flow, ATPG is not a program an engineer babysits but a skill that agents select and call: a compiled engine with golden tests and a manifest that states what it does, what it takes, what it returns and what may be tuned. Designing that skill well means deciding which parts of the ATPG engine are immutable core, which are policy and parameters, and which typed results downstream agents need.

## The manifest and the typed interface

An agent chooses a skill by reading manifests, not source. An ATPG manifest states the supported fault models (stuck-at, transition, cell-aware), required inputs (netlist, scan and clock definitions, fault list, constraints), outputs (patterns in STIL, a coverage report broken down by fault class, pattern count and shift length) and the tunable knobs (abort limit, compaction effort, fault ordering, X handling). Exposed as a tool with typed inputs and outputs, for example over MCP, the skill lets a test-time agent compute cost as patterns times shift cycles without parsing logs, and lets a planner choose between a fast structural engine and a SAT-based engine by their advertised strengths: the SAT engine for many aborted or suspected redundant faults.

## Core, policy and parameters

The **immutable core** is what coverage numbers rest on: fault simulation must never report a pattern as detecting a fault it does not detect, and five-valued evaluation and implication must be correct. **Policy** is how the search proceeds: fault ordering, backtrace heuristics, dynamic compaction strategy, and fault models plugged in at a defined extension point. **Parameters** are the backtrack limit, effort levels and pattern caps. A request for fewer patterns at the same coverage lands at Tier 1 (raise compaction effort, change ordering parameters); a new defect model from a cell-aware table lands at Tier 2 (a plug-in fault model); a change to the implication engine is Tier 3 and needs human review and full regression.

## Safety and reward

An ATPG skill sees the customer's netlist, so it runs sandboxed with no network: the IP boundary is part of the architecture. It must be deterministic and seedable, so a result can be compared with a baseline and attributed to a change. If agents are rewarded on coverage, they can learn to game it: test coverage rises when more faults are declared untestable, so the reward must check that untestable classifications are proven, that the fault list is unchanged, and that pattern count is held within bounds.

## Going deeper

The agent-first EDA notes in this curriculum (skills, manifests, change tiers); the Fault open-source ATPG and the OpenROAD flow for examples of scriptable DFT engines.
