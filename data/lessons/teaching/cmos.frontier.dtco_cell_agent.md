# Agents choosing cell architectures in DTCO

Design-technology co-optimization (DTCO) chooses the cell architecture, track height, fins or sheets per device, pin layers and similar template parameters, jointly with the process. Each candidate must be laid out, characterized and pushed through place and route on benchmark blocks before it can be judged, so exploration is expensive. Agents that choose the next experiment and run the engines as tools make that exploration faster, and the hard part is scoring candidates honestly.

## The architecture choice and its scoring

A cell architecture fixes every cell's layout template: fewer tracks give density but fewer fins per device and less room for in-cell routing and pin access; more tracks give drive and routability at the cost of area. A candidate cannot be scored on single cells, because an area win can be lost to routing congestion and pin-access failures once real nets are routed. The agent's loop is therefore the harness loop applied to exploration: propose an architecture, generate and characterize cells, run place and route on benchmark blocks through typed tool interfaces, read structured QoR results, compare against a baseline, and choose the next candidate. Because each evaluation costs hours, the choice of the next candidate matters, and Bayesian optimization is a natural method.

## Pareto sets, not single winners

Area, delay and power trade against each other, so the output of exploration is a Pareto set: the candidates not dominated by any other. A candidate dominates another if it is at least as good on every metric and strictly better on one. Given candidates scored as (area, delay) of (20, 18), (16, 24), (22, 17), (21, 19) and (16, 22), the point (21, 19) is dominated by (20, 18) and (16, 24) by (16, 22), leaving three Pareto-optimal candidates. Product goals then select from the set: a mobile part weights power and area, a server part weights delay. The agent's job is to map the front accurately; the commitment to one architecture, which every later design inherits, remains a human decision.

## Proxies and the danger of overfitting

Fast proxies, such as cell area and FO4 delay, let the agent prune candidates cheaply, but they miss pin access, routability and wire load, which only appear after routing. An architecture can also overfit: an 8% area gain on one benchmark, at one corner and one utilization, may vanish elsewhere. Robust evaluation spans several benchmarks, corners and utilizations, and treats a reward measured on a narrow setting as an invitation to gaming. Exposing the engines as tools with typed inputs and outputs, rather than scraping logs, keeps experiments reproducible and comparable, which is what lets a model learn from many runs which architectures to try next.

## Going deeper

The agent-first EDA material in this curriculum (harness loop, engines as tools, design-space exploration); published DTCO studies of track-height scaling in IEDM and VLSI Symposium proceedings.
