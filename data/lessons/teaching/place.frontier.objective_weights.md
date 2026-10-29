# Objective weights as Tier 1 configuration

A placer's objective blends wirelength with congestion, timing and density terms through weights. Exposing those weights as configuration lets an agent answer most requests ("less congestion near the memories", "favour timing on this block") without recompiling anything. The agent's job then becomes quantitative: compute costs, compare runs and find the weight at which the preferred result changes.

## A weighted objective as configuration

Write the placement score as cost = HPWL + $\beta$ x overflow + $\gamma$ x |TNS|, or with congestion and power terms as needed. Each weight converts a metric into wirelength-equivalent units: $\beta$ = 10000 says one percent of overflow is worth 10000 um of wire. Run P with HPWL 854000 um, overflow 2% and TNS -391 ps costs $854000 + 20000 + 200 \cdot 391 = 952200$ with $\gamma$ = 200. The weights live in a configuration file or tool command, the engine reads them at start-up, and changing them is a Tier 1 change: no build, only re-runs and a QoR comparison against a baseline.

## Comparing runs

To choose between two runs, compute both costs and take the difference. Run P with HPWL 811000, overflow 4%, TNS -41 ps and run Q with 852000, 1% and -361 ps, under $\beta = 1000$ and $\gamma = 200$, cost 823200 and 925200: Q is worse by 102000, mostly from its TNS. Under a larger $\beta$ the comparison could flip, since Q has less overflow. A robust agent therefore reports not only which run wins but how sensitive the decision is to the weights, because the weights encode a judgement that the user, not the agent, should own. Runs must use the same legalizer, detailed placer and timer, or the comparison measures the tools rather than the settings.

## The break-even weight

The weight at which two runs cost the same follows from setting the costs equal. Holding $\gamma$ fixed, $\beta^* = \dfrac{C_Q - C_P}{\text{OF}_P - \text{OF}_Q}$, where $C$ is HPWL plus the TNS term. If P costs 897930 and Q 951650 without the overflow term, and their overflows are 12% and 10%, then $\beta^* = 53720/2 = 26860$: above it Q's lower overflow wins. An agent can compute such crossovers from a few runs and present them as the decision boundary, rather than silently picking a weight. If no weight setting gives what a request needs, for example a rule about specific cells, the request is not Tier 1 and needs a plug-in.

## Going deeper

The agent-first EDA notes in this curriculum, on change tiers and parameters; the OpenROAD gpl documentation for its density, timing and routability settings.
