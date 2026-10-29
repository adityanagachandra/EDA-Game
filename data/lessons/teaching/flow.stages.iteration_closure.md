# Why design flows iterate

The flow diagram is a straight line, but no real project runs it once. Timing, power, area and verification are checked at many stages with better information each time, and when a later check fails, an earlier stage re-runs. Design closure is the process of converging that loop. For a tool designer, the shape of the loop sets what matters: fast estimates early, accurate analysis late, and incremental engines everywhere.

## The same question asked many times

Timing is estimated in synthesis with wire-load models or virtual placement, re-estimated after placement with HPWL-based wires, after CTS with real clock skew, after routing with extracted parasitics, and finally signed off across all corners. Power and area follow a similar ladder. Each estimate is cheaper and less accurate than the next. A stage optimizes against the estimate it has, so each estimate's error becomes a later stage's surprise. Closure is fast when early estimates correlate well with signoff, and slow when they do not. Correlation between stages is measured and tracked like any other QoR metric.

## Loops and their costs

When a check fails, the cheapest fix is the most local one. A failing path at signoff is first tried with an ECO: resize, buffer, reroute. If that fails, the block may be re-placed with new constraints. If placement cannot fix it, the floorplan or even the RTL changes. Each step back re-runs every stage after it, so the loops are nested and their costs multiply. Typical loops are the timing-closure loop (placement through signoff), the congestion loop (placement and routing), and the functional loop (RTL and verification), which can reopen everything downstream.

## Convergence and when it fails

Iteration helps only if it converges. Two fixes can fight: upsizing gates for setup adds area and power and can cause hold violations; spreading cells to relieve congestion lengthens critical wires. A flow converges when each loop reduces the violations without creating as many elsewhere. Practical safeguards include margins in early stages (over-constraining the clock in synthesis), correlation work that calibrates fast estimates against signoff, and incremental tools that re-analyze only what changed, so each iteration costs minutes rather than a full run. When none of these works, the team changes a constraint or the architecture instead of iterating forever.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 8 on timing closure.
