# Objective weights as Tier 1 parameters

In an agent-first engine, an optimizer's objective weights are its most important change surface. Most user requests, favour timing, cut congestion, trade area for power, are weight changes: Tier 1 configuration with no recompile. Treating weights as first-class parameters lets an agent serve those requests in seconds, but only if it understands what a weight change can reach on the Pareto front and what evidence proves the change did what was asked.

## Which tier a request lands in

The four-layer model separates an optimizer's immutable core (legalization, the numerical kernel, correctness invariants) from its policies, parameters and observation hooks. A request to favour timing over wirelength is a parameter change: raise the timing weight. A request to penalize M2 vias twice as much lands in Tier 1 if the router already exposes per-layer via costs. A request for a cubic congestion penalty when only linear and quadratic forms exist is Tier 2: a new penalty function at the cost-function extension point, with the core untouched. Legality guarantees never move into parameters.

The design rule is to put the seam where common requests become configuration. If a team keeps asking for timing focus in one clock domain and the engine has only a global weight, the right answer is to extend the seam with per-path-group weights once, after which every such request is Tier 1.

## What a weight change selects

A weighted objective picks one point of the trade-off. With normalized run results (lower is better) Run 1: WNS penalty 6, wirelength 1; Run 2: 7, 8; Run 3: 7, 4; Run 4: 6, 9, and weights 0.6 for timing and 0.4 for wirelength, the scores are 4.0, 7.4, 5.8 and 7.2, so Run 1 is chosen. Shifting weight changes the score lines' slope and can move the choice to another Pareto point.

Two cautions apply. Weights act on normalized objectives, so the agent must use the engine's normalization, typically against a baseline run, or the weight numbers mean nothing. And weighted sums reach only the convex hull of the front: a design in a dent cannot be selected by any weights, and an agent that keeps sweeping weights to reach it should switch to an epsilon-constraint formulation, bounding one objective and optimizing the other.

## Proving a weight change worked

A Tier 1 change needs no rebuild, but it changes results, so the harness still runs the reference designs. The pass or fail checks are invariants: legality, no overlaps, cells on sites, fixed cells unmoved. The evidence is a differential QoR report against the baseline: raising the timing weight from 0.3 to 0.6 should improve WNS and TNS while wirelength and congestion change by a bounded amount.

The agent should report both sides of the trade, what was gained and what was paid, not only the improved metric. It must also respect noise: if two settings differ by 0.4% and seed-to-seed variation is about 1%, the difference is not a result until repeated seeds or paired comparisons separate it from noise. A successful change of this kind is a candidate for the skill flywheel: recorded with the design characteristics that made it work, so the next agent can select it.
