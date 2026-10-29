# Agent-first EDA: concepts for the frontier subdomains

This file is the concept source for the `agent_first_eda` domain and for the `<domain>.frontier`
subdomain in every other EDA domain. It distills a modular, agent-editable engine architecture
into general engineering ideas, and adds the public state of agentic EDA as of late 2026. Teach these ideas as engineering, not as product
marketing. Do not name people, customers, dollar figures or internal schedules in lessons.

## 1. The four-layer engine model

Every engine (placer, router, timer, extractor, synthesizer, simulator) is built as four layers,
ordered from hardest to change to easiest:

1. **Immutable core.** Correctness, legality, invariants, the numerical kernel. Examples:
   legalization and overlap-free slot assignment in a placer; the DRC-correct path model and
   grid/layer legality in a router; graph propagation and the delay-model API in a timer; the
   field-solver mathematics in an extractor. Changing it needs human approval and full
   regression.
2. **Strategy or policy.** The heuristics that decide *how*: cost functions, move sets,
   orderings, rip-up strategies. These are selected from a library of compiled policies, and new
   ones are added as plug-ins.
3. **Parameters.** Weights, thresholds, modes. Configuration only, with no recompile. Most change
   requests land here.
4. **Observation hooks.** Where an agent reads live state (timing, congestion, density, IR drop)
   and injects behavior.

Design rule: keep the core small and rarely touched. Push everything that varies into policy,
parameters and hooks, so that the common change never recompiles the core.

Per-engine split of core versus change surface:

| Engine | Core (fixed) | Change surface |
| --- | --- | --- |
| Placer | Legalization, overlap-free slots, row and site rules | Objective weights (wirelength, congestion, timing), move set, density targets, rules such as "abut regular cells" |
| Router | DRC-correct path model, grid and layer legality | Cost and penalty functions, rip-up and reroute strategy, layer directives |
| Timer | Graph propagation, delay-model API | Which corners and paths to report, pruning, fast or surrogate modes, derating schemes |
| Extractor | Field-solver mathematics, netlist correctness | Accuracy-versus-speed mode, net and layer selection, pattern library, calibration |

## 2. Skills, manifests and selection

A **skill** is an engine, or a piece of one, packaged in a standard folder: source code, a
build, golden tests (known-good inputs and outputs) and a short plain-language manifest that
says what the skill does, its inputs and outputs, and which parts are meant to be tuned.

An agent does not read engine source to choose a skill. It reads the manifests, like catalog
entries. It matches the design's characteristics (cell regularity, corner, clock structure, the
current objective) against what each manifest advertises, and calls the compiled binary.
Selection is a cheap text decision. The expensive compute stays in the engines.

## 3. The harness loop

The harness automates the edit, compile and test cycle that engine developers already run:

1. A request arrives in plain language.
2. The model proposes a small code edit.
3. The harness writes the file and runs an incremental build (seconds, with a compiler cache).
4. The harness runs the golden tests and a QoR comparison against a baseline.
5. A build or test failure goes back to the model as feedback, and the model tries again.
6. When the tests pass, the change is shown to a human as a reviewable diff.

## 4. Change tiers

The model edits at the smallest safe tier:

- **Tier 1, configuration.** A value or table change. No recompile.
- **Tier 2, plug-in.** A small implementation at a defined extension point. The core is untouched.
- **Tier 3, core edit.** A change to the core source. Rare, and only with a human-approved diff.

A seam is in the right place when Tier 3 edits essentially never happen for that kind of change.

## 5. Seams and extension points: the derating example

Timing derating is the clearest example. The propagation mathematics (arrival, required, slack,
path enumeration) is settled, but every team margins differently: flat OCV, AOCV tables by logic
depth and distance, POCV with LVF sigma, plus custom margin stacks from silicon correlation.

- **Core:** timing-graph propagation and the delay-model API. When the core needs an arc delay it
  asks one question at one seam: what derate applies to this arc in this context?
- **Seam:** a `DerateProvider` interface. Given the context (arc, cell and Vt class, early or
  late, rise or fall, logic depth, distance, slew, corner, clock versus data, path group) it
  returns a factor (OCV, AOCV) or a mean and sigma adjustment (POCV).
- **Skills:** `FlatOCV`, `AOCV`, `POCV` providers and a `CustomMargin` layer that stacks on top.
- **Tier 1:** choose the provider, the sigma count, global early and late derates, table paths,
  per-path-group margins.
- **Tier 2:** a new provider, for example one that adds 3% late derate on data paths through HVT
  cells with logic depth under 5, and falls through to AOCV for clocks.
- **Proof:** a differential timing report on a reference design shows that only the intended
  paths moved, by the intended amount.

Code-structure rules for seams: narrow and stable interfaces (abstract base plus factory), a
frozen extension ABI, the hot kernel inside the core, and coarse plug-in calls (per batch, not per
cell) so the compiler can still inline and vectorize. Pass plain data structures, keep no hidden
global state, and make everything deterministic and seedable.

## 6. Verification of engine changes

Every engine ships with golden tests and a QoR baseline for instant differential checks. Nothing
is promoted without passing them. Skills run in a sandbox with CPU, memory and time limits and no
network. Every change is versioned, reviewable and reversible. Design and process data stay
on-site: the IP boundary is part of the architecture, not a policy note.

## 7. The skill flywheel and the brain model

A skill found for one design (by the model, an agent, a field engineer or a customer) goes
through four steps: identify, iterate quickly as a plug-in, validate with golden and
differential tests, promote into a compiled library skill. The promoted skill lands in a shared
library, where any agent can select it at compiled speed, and in the internal domain model (the
"brain"), which learns the skill and when to pick it. The know-how flows back. The customer's
design does not.

## 8. Fast and accurate engines in a continuous loop

Engines exist in a fast version and an accurate version. The fast versions (fast placer, fast
timer, fast extraction, fast congestion estimator) run in a continuous modeling loop: place,
route and analyze on the fly. The model watches the metrics and their variance to decide which
skill to run next, which skills need customizing, and where to recurse into the accurate
engines. Surrogate models are one way to build a fast mode.

## 9. Models and learning

- **Three roles a base model can play:** drive the harness (any strong tool-calling model);
  run the continuous skill-selection loop (close to the data, fine-tuned); become the on-site
  domain model (needs open weights). A closed API model can only fill the first role.
- **Model abstraction layer:** the harness stays model-agnostic behind one tool-call interface,
  with a thin adapter per model (chat template, tool-call format). Mixing models is a design
  goal.
- **Tool-grounded reinforcement learning:** deterministic solvers and EDA tools score a model's
  outputs against objectives, constraints and ground truth, and that score is the reward. Reward
  design, reward hacking (gaming a fast proxy), and the gap between fast-proxy reward and signoff
  QoR are the central risks.
- **Evaluation:** measure compile-pass rate and skill-edit correctness on real engines, plus QoR
  deltas against baselines. Generic coding benchmarks under-separate models on this task.

## 10. Agents in the flow

- **MCP and tool interfaces:** engines exposed as tools with typed inputs and outputs, so other
  AI systems can call them.
- **Agentic signoff and ECO loops:** an agent reads timing, DRC and IR reports, proposes fixes,
  runs them through the engines, and iterates, with a human approving the final ECO.
- **Multi-agent orchestration:** specialized agents (RTL, verification, implementation, signoff)
  coordinated by a planner, with shared state in the design database.
- **Design-space exploration agents:** parallel experiments over flow parameters, with a model
  choosing the next experiment.
- **Industry context (public, late 2026):** commercial vendors ship "super-agent" products for
  front-end design and verification, digital implementation, analog, and PCB and packaging, and
  domain-specialized model portfolios trained with tool feedback. Vendors report gains as "up to"
  figures against unnamed baselines, so treat them as claims, not results.

## What a frontier subdomain should teach

Each `<domain>.frontier` subdomain (4 to 6 concepts) applies these ideas to its own domain. Good
frontier concepts are concrete and checkable. For example:

- which layer of the four-layer model a given change belongs in, for this domain's engine
- where the seam belongs for this domain's most customer-specific behavior, and what context it
  must expose
- what the golden and differential tests for this engine must check
- what the fast mode of this engine trades away, and when an agent must recurse to the accurate one
- what a tool-grounded reward for this domain looks like, and how it can be gamed
- what this domain's agentic ECO or closure loop reads, proposes and verifies
