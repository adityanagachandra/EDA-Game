# Skew of a given clock tree

Before building trees, a CTS engine must evaluate them. Given a candidate topology with wire resistances and capacitances, it computes the delay to every sink and the skew. Elmore delay is the standard first estimate, cheap enough to call inside every merge and every move. This concept applies Elmore delay to clock trees and introduces the simplest skew repair, a dummy load at the faster sink, which shows how capacitance moves delay.

## Elmore delay at a clock sink

In an RC tree, the Elmore delay to a node is the sum over edges on its root path of the edge resistance times the capacitance downstream of that edge. Compute downstream capacitance bottom-up first. Take a root that drives node u through $100\ \Omega$ ($C_u = 10$ fF) and node w through $150\ \Omega$ ($C_w = 20$ fF); u feeds sinks of 30 and 20 fF through $100\ \Omega$ each, and w feeds a 30 fF sink through $400\ \Omega$. Downstream of w is 50 fF, so the delay to w's sink is $150 \cdot 50 + 400 \cdot 30 = 19\,500\ \Omega\cdot$fF, which is 19.5 ps.

## Skew from the sink delays

Repeat the computation for every sink and subtract the smallest delay from the largest. In a tree where u carries 60 fF and w carries 20 fF, sinks might sit at 7, 14 and 4 ps, giving a skew of 10 ps. Two effects drive the imbalance: a sink at the end of a long or resistive branch, and a heavy sibling that loads the shared edge above it. A timer computes all sink delays in one pre-order pass after one post-order pass, so the cost is linear in the tree size. CTS engines evaluate thousands of candidate trees this way.

## Balancing with a dummy load

Adding capacitance $\Delta C$ at a sink slows every node above it. With a trunk $R_0$ and two branches $R_a$ and $R_b$, a load at sink b raises both sinks' delay by $R_0 \Delta C$ through the shared trunk, and sink b's by a further $R_b \Delta C$. The skew therefore closes by $R_b \Delta C$, so a 5 ps skew with $R_b = 200\ \Omega$ needs $\Delta C = 5000/200 = 25$ fF at the faster sink. Dummy loads are cheap to add but cost power and latency. CTS engines usually prefer to re-balance through tapping points, wire snaking or buffer sizing.

## Going deeper

Elmore, "The transient response of damped linear networks with particular regard to wideband amplifiers", Journal of Applied Physics (1948); Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the clock chapter.
