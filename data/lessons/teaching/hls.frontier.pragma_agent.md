# Pragma-search agents

Choosing HLS directives is a search over a huge discrete space with an expensive, noisy evaluator, exactly the kind of loop an agent can run. A pragma-search agent proposes directive settings, runs HLS, reads the report and decides the next experiment. Done well, it reaches good Pareto points in far fewer runs than blind search; done badly, it games estimates and ships unverified designs.

## The search problem

The directive space multiplies: a few loops and arrays with a handful of options each give thousands of configurations, and each HLS run takes minutes to hours. An agent therefore works at Tier 1: it edits directive values such as II targets, unroll factors and partition factors, never the engine. Before spending a run it applies cheap legality and resource checks: complete partitioning of a million-element array would need a million registers and can be rejected immediately, and an unroll factor beyond the trip count is equivalent to full unrolling. After each run it reads the structured report, the achieved II and latency per loop, resource use, and the stated reason a loop missed its target, rather than only a pass or fail code. The report is the observation hook the agent's decisions depend on.

## Reading the report to choose the next move

The most valuable information in an HLS report is the binding constraint. A loop at II = 4 limited by the memory ports of array A calls for partitioning A cyclically or restructuring its accesses; adding multipliers would change nothing. A loop at II = 5 limited by a carried dependence through $s$ is RecMII-bound: partitioning arrays will not lower it, and the fix is to reassociate or break the recurrence, which is a source change outside the agent's Tier 1 remit and must go to a human or a reviewed transformation. Each proposal is judged as a delta against a stored baseline for that kernel. When two proposals give identical reports, the varied directives did not touch the binding constraint, and a good policy moves its search elsewhere instead of repeating similar points.

## Rewards, proxies and verification

An agent rewarded only on reported latency learns to unroll and partition everything, spending unlimited area; the objective must include resource limits or be posed as a Pareto search. HLS reports are estimates that miss routing delay and congestion, so a design with excellent reported numbers may miss timing after implementation; the loop must recurse into the accurate engines, logic synthesis and place and route, for the finalists. A result that never ran C/RTL co-simulation is unverified, however good its QoR, since directives like dataflow or interface changes can introduce protocol or ordering bugs. And some moves must stay out of the agent's hands entirely: editing the C source to delete a dependence the agent judges unnecessary can change the function computed, and needs human review.

## Going deeper

Schafer and Wang, "High-level synthesis design space exploration: past, present, and future" (IEEE TCAD, 2020); the agent_first_eda lessons on the harness loop and fast versus accurate engines.
