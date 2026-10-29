# Authoring an EDA domain

This repository is Mathua, an adaptive mastery-learning engine (Go back end, Next.js front end)
modeled on Math Academy. It holds about 1250 mathematics concepts. We are adding a full
**electronic design automation (EDA)** curriculum on top of that mathematics: 33 domains and
about 1,450 concepts. Each domain is authored by one agent, in parallel. This file is that
agent's contract. Read it end to end before you write anything.

## Who the learner is

The learner is an applied-mathematics undergraduate with strong technical maturity. They want
the fastest possible ramp to being a **designer of EDA tools**: someone who can read, modify and
eventually write placers, routers, timers, synthesizers, simulators and verifiers, and who
understands where agent-first, AI-native EDA is going. Teach:

- **mechanisms and first principles**: what the engine computes and why
- **formulations**: decision variables, objective, constraints, why exact optimization is hard,
  and how practical methods relax, approximate, decompose or search
- **representations**: what each data structure or file format preserves, abstracts away, and
  lets a tool compute
- **trade-offs and failure propagation**: how a poor decision here hurts a later stage
- **implementation reality**: complexity, incremental updates, data structures, the open-source
  tool that implements the idea (Yosys, ABC, OpenROAD, OpenSTA, KLayout, Magic, ngspice, Xyce,
  Verilator, SymbiYosys, nextpnr, VPR and so on)

Do not teach tool button-pushing. Do not over-index on machine learning. Be technically exact:
a wrong formula or a confidently stated false claim is the worst possible defect. If you are not
sure of a fact (a date, a version, a vendor claim), leave it out.

## What you produce

For your domain `<domain>` with prefix `<prefix>` and generator package key `<key>` (all three
are in `docs/eda/curriculum.json`):

| File | What |
| --- | --- |
| `data/concepts/<domain>.json` | The concept list (schema below) |
| `data/lessons/teaching/<id>.md` | One lesson per concept |
| `data/lessons/kp/<id>.json` | One knowledge-point shard per concept: exactly 3 steps |
| `internal/generator/eda/<key>/*.go` | The generators, a `Register(reg *generator.Registry)` function |
| `internal/generator/eda/<key>/generators_test.go` | Calls `edakit.CheckDomain(t, Register, "<domain>")` |

**Do not edit any other file.** Do not touch `internal/generator/all`, `data/lessons/lessons.json`,
`data/concepts/enrichment.json`, `internal/generator/eda/edakit`, another domain's files or the
mathematics corpus. The coordinator wires everything in afterwards. If you need a helper the kit
lacks, write it inside your own package. Exception: `timing_analysis` already contains two
reference concepts (`sta.delay.elmore`, `sta.graph.slack_meaning`); the `sta` author keeps them
and builds the rest of the domain around them.

## The concept list

```json
{
  "id": "sta.delay.elmore",
  "label": "Elmore delay of an RC tree",
  "domain": "timing_analysis",
  "subdomain": "timing_analysis.delay",
  "grading_type": "numeric",
  "prerequisites": ["em.circuit.rc_time_constant", "graphs.tree.rooted_traversal"],
  "key_prerequisites": ["em.circuit.rc_time_constant"],
  "encompasses": ["em.circuit.rc_time_constant"],
  "mastery_threshold": {"streak": 5, "avg_time_seconds": 120}
}
```

Rules:

- **ids** are `<prefix>.<subdomain-short>.<name>`, lower-case snake segments, for example
  `place.global.electrostatic`. Every anchor listed for your domain in `curriculum.json` must
  exist with exactly that id. Use the anchor's label or a close improvement.
- **labels** are plain text (no `$`, no backslashes), at most about 60 characters, phrased as the
  topic ("Kernighan-Lin partitioning"), not as a question.
- **subdomains** are `<domain>.<name>`. Use 5 to 9 subdomains, ordered from foundational to
  advanced, plus the mandatory `<domain>.frontier` (see below).
- **granularity**: one concept is one thing a learner can fail at independently, small enough to
  master in one sitting. Split "SAT solving" into encoding, unit propagation, conflict analysis,
  clause learning, decision heuristics, restarts and so on. Hit the domain's `target` count
  (within 15%). Cover every item in the domain's `covers` list. If a covers item overlaps with
  another domain, teach it from your domain's angle and depend on the other domain's anchor for
  the shared foundation rather than re-teaching it.
- **grading_type** is `numeric` or `multiple_choice`, nothing else.
- **mastery_threshold**: numeric concepts use `streak` 5 and `avg_time_seconds` 60 to 180 (sized
  to the computation). Multiple-choice concepts use `streak` 6 and `avg_time_seconds` 30 to 75.

### The graph is the product

The scheduler will never show a concept until its prerequisites are mastered, so the
prerequisite graph *is* the course. Build it like a careful textbook author:

- **prerequisites**: 2 to 4 per concept. List what the learner must already know, and nothing
  that is merely related. The first concepts of a subdomain lean on upstream anchors and the
  mathematics substrate, later ones on earlier concepts of your own domain.
- **key_prerequisites**: 1 or 2 of the prerequisites, the ones the concept leans on hardest
  (when a learner fails, the engine remediates these first). Required on every concept.
- **encompasses**: lower-level concepts that this concept *practices implicitly*. Solving it
  exercises them, so a correct answer earns them repetition credit. Each encompassed id must be a
  prerequisite or an ancestor through your own domain. At least 60% of your concepts declare one
  or more. Example: computing a full-path setup slack encompasses arrival-time propagation and
  NLDM interpolation.
- **mean in-degree** (prerequisites plus encompasses, per concept) must be at least 2.3.
- **What you may reference** (enforced by `scripts/eda_check.py`):
  1. your own domain's concepts;
  2. the mathematics substrate, listed in `docs/eda/math-substrate-ids.tsv` (id, domain,
     label). Use it: linear algebra, numerical analysis, graph theory, discrete mathematics,
     probability and statistics, ODEs, calculus and machine learning are already taught there;
  3. the **anchors** of the domains in your `upstream` list (`curriculum.json`). Only anchors.
     They are guaranteed to exist. Non-anchor concepts of other domains may not exist yet;
  4. frontier concepts only: the anchors of `agent_first_eda`.
  
  The domain layers make this acyclic by construction. Never depend on a lower-layer domain's
  non-anchor concept, and never on a domain that is not upstream of you.

### The frontier subdomain

Every domain except `agent_first_eda` ends with a subdomain `<domain>.frontier` of **4 to 6
concepts** with ids `<prefix>.frontier.<name>`. They teach the AI-native, agent-first future of
*this* domain: read `docs/eda/agent-first-concepts.md` for the ideas (the four-layer engine
model, skills and manifests, the harness loop, change tiers, seams and extension points, golden
and differential testing, the skill flywheel, fast versus accurate engines, tool-grounded RL,
agentic signoff and ECO loops). Apply them concretely to your domain's engines. Example for
timing: the `DerateProvider` seam, fast surrogate timers in a modeling loop, an agentic timing
ECO loop, and what a differential timing check must prove. Frontier concepts depend on your own
core concepts plus one or two `agentic.*` anchors. Make them scenario-based and checkable ("this
customer request: which tier does it land in, and why?"), not vague futurism.

`agent_first_eda` itself has no frontier subdomain: the whole domain is the frontier. Its
concepts may not depend on any `*.frontier.*` concept.

## The lesson

`data/lessons/teaching/<id>.md`:

```markdown
# <Label, plain text>

<An opening paragraph of 40 to 100 words: what the concept is, where it sits in the flow, why a
tool designer cares.>

## <Knowledge point 1 title>

<90 to 200 words. One idea. Concrete. Include a small worked example with numbers wherever the
concept computes anything.>

## <Knowledge point 2 title>

...

## <Knowledge point 3 title>

...

## Going deeper

<Optional. One to three canonical sources by name, for example: Kahng, Lienig, Markov and Hu,
"VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 4;
the OpenROAD documentation for gpl; Biere et al., "Handbook of Satisfiability". No URLs.>
```

- Exactly three `## ` knowledge-point sections, then the optional `## Going deeper`. No other
  `## ` headings. `### ` subheadings inside a section are allowed but rarely needed.
- The three section titles must equal the three `section` values of the shard, character for
  character. Keep them plain text: no `$`, no backticks.
- At least 250 words in total. Each section must make sense alone, because the app shows one
  knowledge point at a time.
- Mathematics in `$...$` (inline) or `$$...$$` (display). KaTeX, so standard LaTeX only. Write
  units in text (`ps`, `fF`, `ohm` or `$\Omega$`).
- No emoji, no circled letters such as (a) written as a circled glyph, no pictographs. Arrows and
  mathematical symbols are fine.
- Write the knowledge points in the order the learner climbs: definition and mechanism, then the
  computation or decision procedure, then the subtlety, edge case or engineering consequence.

## The shard

`data/lessons/kp/<id>.json`, exactly three entries:

```json
[
  {"label": "Required minus arrival", "section": "Required minus arrival",
   "subgoals": ["Define arrival time as the latest time the signal settles",
                "Define required time as the latest time the capture flop can accept it",
                "Subtract arrival from required to get setup slack"]}
]
```

`label` equals `section`. Each entry has exactly three subgoals: short imperative steps that name
the actual move ("Sum the products of shared resistance and capacitance"). Never use filler such
as "Review ...", "Understand ...", "Follow the worked example", "Apply the rule", "Check the
result", "Verify".

## The generators

Read `internal/generator/eda/edakit/edakit.go` (the kit) and the reference package
`internal/generator/eda/sta/generators.go`. Problems are **generated, not stored**: a learner
can answer dozens of questions on one concept, and none may repeat inside a session.

- **Computational concepts** (`edakit.Func`): draw a random instance and compute the answer.
  Examples: Elmore delay, slack along a path, HPWL of a net, FM gain of a cell, cut size,
  Lee-wavefront distance, left-edge track count, number of BDD nodes, unit propagations, fault
  coverage, LFSR next state, ASAP schedule length, Newton iterate, MNA matrix entry, IR drop at a
  node, logical-effort delay, NLDM bilinear interpolation, yield from defect density. Aim for
  **at least 40% of your concepts to be computational** (numeric, or a multiple-choice question
  whose options are computed from the drawn instance). Many EDA ideas compute; find the
  computation.
  - Answers are exact integers, or decimals with the rounding stated in the stem ("Round to 2
    decimal places."). Format with `edakit.Num(x, places)` or `edakit.IntStr(n)`. The numeric
    grader accepts any value within half a unit of the last decimal place shown.
  - A numeric answer contains only digits, sign, decimal point and commas. Never a unit, a
    letter or a matrix. A comma-separated answer is graded as an **unordered** list.
  - State units in the stem and say what unit to answer in.
  - Scale instance size with `ctx.Difficulty` (`edakit.Scaled`). Where the three knowledge points
    differ in kind, branch on `ctx.KPOrdinal` (1, 2, 3; 0 means unguided).
  - Write an `Explanation` that shows the computation with this instance's numbers, and set
    `Facts` to the key intermediate values.
- **Conceptual concepts** (`edakit.Bank`): an authored bank of `edakit.Item`. Each item is a
  distinct situation, not a rephrasing. **At least 12 items per bank**, at least 3 per knowledge
  point (set `KP` to 1, 2 or 3), so guided practice on each step has supply. Each item has one
  unambiguous correct option and three or four distractors that are *plausible misconceptions*
  (the wrong formula, the wrong direction, the confusion with a neighbouring concept), similar in
  length and style to the answer. `Why` explains the reasoning in one or two sentences. Prefer
  scenario questions ("a placer reports... what happened?", "which representation lets the tool
  compute X?") over definition recall.
- **Mixed concepts** (`edakit.Mix`): combine a computational part and a bank, for example
  computing a gain at KP 2 and reasoning about it at KP 3.
- Use `$...$` in stems when it helps. A literal dollar sign must not appear in prose. Write one
  backslash per LaTeX command in the rendered string (in a Go interpreted string literal that
  means `"\\Omega"`; a raw string literal needs a single backslash).
- Never put the answer text inside the stem.
- Generators use the global `math/rand` (through the kit). Do not seed it yourself.
- Keep files reasonable: split the package into several files by subdomain
  (`partition.go`, `fm.go`, `frontier.go`, ...).

## Validate before you report

Run, from the repository root (Go is at `/opt/homebrew/bin/go`; put it on `PATH`):

```bash
export PATH=/opt/homebrew/bin:$PATH
gofmt -l internal/generator/eda/<key>/        # must print nothing
go vet ./internal/generator/eda/<key>/
go test ./internal/generator/eda/<key>/       # the domain contract
python3 scripts/eda_check.py <domain>          # the graph contract
```

All four must pass. Run only your own package's tests: other domains are being written at the
same time, so `go test ./...` will see their half-finished state.

Then read back a sample of your own output as a learner would: run a small Go test or program
that prints three generated problems per concept, and read them. Fix anything confusing, wrong or
repetitive.

## Report

When you finish, report in under 200 words: the concept count, the subdomain list, the share of
computational concepts, any covers items you placed in a different domain and why, and any
external references that you wanted but could not make because they were not anchors.
