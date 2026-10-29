# Reward design: gates, normalization and shaping

A reward function is a specification. Whatever it fails to say, the policy will exploit. Designing a reward for an EDA task means deciding which conditions are absolute, how to make metrics comparable across designs, and how to combine several objectives into one number without letting one dominate. Each choice can be checked with small calculations.

## Legality gates first

Any output that violates a hard condition earns no credit, however good its metrics: a placement with overlaps, a routing with DRC violations, an edit that fails to build or fails a golden test, a netlist not equivalent to its specification. The gate is checked by an independent engine, not by the component being trained. Only outputs that pass all gates are scored on continuous metrics. Without gates, the first thing a policy learns is that an illegal placement has excellent wirelength. With them, of six candidate placements where two have overlaps and one has DRC violations after routing, only three are scored at all.

## Normalizing per design

Raw metrics scale with design size: a 2-million-cell design has far more wirelength than a 50,000-cell one, so a reward of minus wirelength lets large designs dominate training. The standard fix is relative improvement over a baseline on the same design. For a metric to minimize, $r = (x_{base} - x)/x_{base}$; for slack, which may be near zero or negative, divide the improvement by the clock period instead: $r = (\text{WNS} - \text{WNS}_{base}) / T_{clk}$. A wirelength drop from 2.00 m to 1.94 m gives $r = 0.03$ on any design size; a WNS change from -80 ps to -60 ps at a 1,000 ps period gives $r = 0.02$.

## Combining objectives

A composite reward is a weighted sum of normalized terms, $R = \sum_i w_i r_i$, applied to gated outputs. With weights 0.5 for wirelength, 0.3 for overflow and 0.2 for WNS, a candidate improving wirelength by 3%, overflow from 1,200 to 900 (25%) and WNS by 0.02 earns $0.015 + 0.075 + 0.004 = 0.094$. The overflow term dominates because its relative change is large, which may or may not be intended; inspecting each term's contribution is how weights get set. Terms with very different natural ranges may need clipping, so a huge improvement in one easy metric cannot buy a regression in another.

## Going deeper

Ng, Harada and Russell, "Policy invariance under reward transformations" (ICML 1999), on reward shaping; Sutton and Barto, "Reinforcement Learning: An Introduction" (2nd ed., 2018).
