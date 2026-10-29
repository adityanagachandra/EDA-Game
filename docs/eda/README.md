# The EDA curriculum

Mathua now teaches electronic design automation on top of its mathematics corpus. The EDA
curriculum is 33 domains and 1607 concepts, each with a three-step lesson and a generator that
produces a fresh problem on every draw. It sits on the existing university mathematics (linear
algebra, numerical analysis, graph theory, probability, ODEs, discrete mathematics and machine
learning), which it uses as prerequisites instead of re-teaching.

## Run it

```bash
cd web/next-app
NEXT_PUBLIC_API_URL="" NEXT_PUBLIC_SAME_ORIGIN_API=1 npm run build   # Node 24
cd ../..
go build -o bin/mathua ./cmd/mathua
DATABASE_URL=sqlite:mathua-eda.db ./bin/mathua --serve --port 8080
```

Open http://localhost:8080. Choose a destination on `/plan`: "EDA Tool Designer: the full ramp"
takes every domain in dependency order, and five narrower tracks cover physical design, synthesis
and verification, circuits and devices, agent-first architecture and boards, packages and FPGAs.
On `/onboard`, select the domains to learn: the diagnostic places you on the frontier, so the
mathematics you already know is credited instead of drilled.

## How the curriculum is built

| File | Role |
|---|---|
| `scripts/eda_curriculum.py` | The source of truth for the domain list: layer, upstream domains, target size, the seed topics each domain covers and the anchor concepts other domains may depend on. It writes `docs/eda/curriculum.json`. |
| `docs/eda/AUTHORING.md` | The contract every domain was written to. |
| `docs/eda/agent-first-concepts.md` | The source of the agent-first ideas taught in `agent_first_eda` and in every `<domain>.frontier` subdomain. |
| `data/concepts/<domain>.json` | Concepts, prerequisites, key prerequisites and encompasses. |
| `data/lessons/teaching/<id>.md`, `data/lessons/kp/<id>.json` | One lesson and one three-step shard per concept. |
| `internal/generator/eda/<key>/` | One generator package per domain, built on `internal/generator/eda/edakit`. |
| `scripts/eda_check.py` | The graph gate for one domain or `--all`. |
| `scripts/eda_wire.py` | Writes `internal/generator/all/eda.go`, which registers every EDA package. |
| `scripts/eda_integrate.py` | Maps lessons, and writes the EDA courses and tracks into `data/courses.json` and `data/destinations.json`. |

### Why the graph is acyclic by construction

Every domain has a layer. A concept may depend on its own domain, on the mathematics substrate,
and on the **anchors** of the domains listed as its upstream, which are always at a lower layer.
Anchors are named in `curriculum.json` before any content exists, so 33 domains could be written
in parallel without seeing each other. Frontier concepts may also depend on the anchors of
`agent_first_eda` (layer 13), and `agent_first_eda` never depends on a frontier concept, so the
one backward edge per domain cannot close a cycle.

### The layers

| Layer | Domains |
|---|---|
| 1 | EDA landscape, Boolean reasoning, graph algorithms, computational geometry, optimization, numerical methods, variability statistics, circuits and electromagnetics |
| 2 | Digital logic, semiconductor devices |
| 3 | CMOS circuits and cells, fabrication and DTCO, HDLs and compilers |
| 4 | Logic synthesis, formal verification, functional verification, circuit simulation |
| 5 | Technology mapping, high-level synthesis, static timing, test and DFT |
| 6 | Floorplanning and placement, analog and RF, FPGA CAD |
| 7 | Clock trees and power delivery, routing |
| 8 | Extraction and signoff, physical verification and DFM |
| 9 | PCB design |
| 10 | Packaging, chiplets and 3D integration |
| 11 | EDA software infrastructure (the open-source RTL-to-GDS flow hands-on) |
| 12 | Machine learning for EDA |
| 13 | Agent-first EDA architecture |

## Change it

Edit a domain's concept file, lessons or generators, then run:

```bash
python3 scripts/eda_check.py <domain>
go test ./internal/generator/eda/<key>/
python3 scripts/eda_integrate.py
```

A new domain needs an entry in `scripts/eda_curriculum.py`, then `python3 scripts/eda_curriculum.py`
and `python3 scripts/eda_wire.py`. After any change in corpus size, refresh the derived files:
`go run ./cmd/masteryunits`, `go run ./scripts/merge_concepts`, `python3 scripts/counts.py --write`
and `go test ./internal/generator/all -run TestConceptVarietyFloor -update-variety`.

`edakit.CheckDomain` draws 20000 problems per concept under a watchdog. Bound every redraw loop
with `generator.MaxDrawAttempts` and redraw every operand its condition reads: a loop that
re-rolled one operand hung once in about seven thousand draws before this check existed.
