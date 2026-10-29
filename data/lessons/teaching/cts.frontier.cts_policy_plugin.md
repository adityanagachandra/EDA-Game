# CTS policies as Tier 2 plug-ins

A CTS engine has a small correctness core and a large set of policies: how sinks are clustered, which topology is built, where gates are cloned, how latency targets are set. Customer requests land mostly in those policies. Organising the engine around the four-layer model and change tiers keeps the common change a parameter or a plug-in, never a core edit, and makes each change provable with golden and differential tests.

## Core and change surface of a CTS engine

The immutable core holds legal buffer placement, the delay and slew evaluation API, the tree data structure and its invariants, and the interface to the timer. The change surface holds the clustering heuristic, the topology policy (H-tree style, DME-based, greedy matching), gating clone and merge rules, buffer selection, and the order of processing clock domains, plus parameters such as skew, latency and slew targets and maximum fanout. Changing the delay and slew evaluation is a Tier 3 core edit: it underpins every decision and needs human review and full regression.

## Classifying requests by tier

A request to target 30 ps of skew instead of 50 ps is Tier 1, a parameter. A request to cluster sinks by timing criticality as well as location, which the clustering seam allows but no library policy implements, is Tier 2: a new clustering policy plug-in. So is a rule that clones every gate whose enable slack falls below 20 ps, written against a gating-policy seam that reads slack and calls the core's cloning service. A request to change how stage delay is computed is Tier 3. The seam is in the right place if Tier 3 essentially never happens for this kind of request.

## Writing and promoting a policy

A policy plug-in receives sink locations and loads, clock-domain and gating structure, blockages and the current targets; it should be called once per clock net or cluster, so the hot delay kernel stays in the compiled core. It must have no hidden global state, so results are deterministic and composable. Before promotion it must pass the engine's golden CTS tests and show skew, latency, power, buffer area and wirelength against the baseline on reference designs; identical timing with 3% more buffer area is a regression. Its manifest states which designs it suits, so an agent can select it, and per-design deltas are reported honestly.

## Going deeper

the OpenROAD documentation for the CTS module.
