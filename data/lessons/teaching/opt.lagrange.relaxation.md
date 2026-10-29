# Lagrangian relaxation

Lagrangian relaxation moves hard coupling constraints into the objective with multipliers. In gate sizing the coupling constraints are the timing constraints on every arc of the timing graph; relaxing them leaves a problem that separates gate by gate. The landmark result is that the multiplier optimality conditions reduce to flow conservation on the timing graph, so the multipliers can be kept on a flow polytope and the per-gate subproblems solved in closed form. It remains a foundation of academic sizers and of modern sizing contests.

## Flow conservation of multipliers

Sizing minimizes area subject to arrival-time constraints $a_j \ge a_i + d_{ij}$ on every arc and $a_{sink} \le T$. Attach a multiplier $\lambda_{ij} \ge 0$ to each arc constraint. Setting the derivative of the Lagrangian with respect to each arrival time $a_i$ to zero gives
$$\sum_{\text{arcs into } i} \lambda_{ki} = \sum_{\text{arcs out of } i} \lambda_{ij}$$
at every internal node: the multipliers form a flow from the primary inputs to the outputs.

Example: a gate's incoming arcs carry 1.0, 2.4 and 0.3, total 3.7; two of its three outgoing arcs carry 1.5 and 0.4. The third must carry $3.7 - 1.9 = 1.8$. Restricting multipliers to this flow polytope eliminates the arrival times from the relaxed problem altogether. Critical paths carry large flow, and non-critical arcs carry little or none.

## The per-gate subproblem

With multipliers fixed on the flow polytope, the relaxed problem separates. For a gate of size x with area $\alpha x$ and delay $RC_L/x$ (drive resistance R at unit size, load $C_L$), and with $\mu$ the sum of multipliers on the arcs through the gate, the gate minimizes
$$\alpha x + \mu \frac{RC_L}{x}, \qquad x^* = \sqrt{\frac{\mu R C_L}{\alpha}}.$$
With $\alpha = 3$, $\mu = 1.7$, $R = 4$ and $C_L = 14$: $x^* = \sqrt{95.2/3} = 5.633$.

The multiplier acts as a price on delay. A gate carrying a lot of multiplier flow sits on critical paths and is sized up; a gate with $\mu = 0$ shrinks to its minimum size. In practice $C_L$ depends on the sizes of downstream gates, so the sizer sweeps the gates in topological order, a few local passes per multiplier update, until the sizes settle for that set of multipliers.

## Upstream load and the lower bound

A gate's input capacitance $c_{in}x$ also slows its driver: the driver's delay grows by $R_{up}c_{in}x$, weighted by the driver's multipliers $\mu_{up}$. The linear coefficient becomes $\alpha + \mu_{up}R_{up}c_{in}$, and
$$x^* = \sqrt{\frac{\mu RC_L}{\alpha + \mu_{up}R_{up}c_{in}}}.$$
With $\alpha = 2$, $\mu = 1.7$, $R = 6$, $C_L = 5$, $\mu_{up} = 0.9$, $R_{up} = 2$ and $c_{in} = 2$: the coefficient is $2 + 3.6 = 5.6$ and $x^* = \sqrt{51/5.6} = 3.018$. Ignoring the upstream term would give $\sqrt{25.5} = 5.05$, an oversized gate that slows its own driver.

For any multipliers the relaxed optimum is a lower bound on the minimum area that meets timing. The outer loop raises multipliers on violating paths and lowers them on slack ones, then re-projects them onto the flow polytope; the gap between the best timing-feasible sizing and the best bound measures convergence.

## Going deeper

Chen, Chu and Wong, "Fast and exact simultaneous gate and wire sizing by Lagrangian relaxation" (IEEE TCAD, 1999); Ozdal, Burns and Hu, "Algorithms for gate sizing and device parameter selection for high-performance designs" (IEEE TCAD, 2012), for discrete LR sizing; the ISPD gate-sizing contest papers.
