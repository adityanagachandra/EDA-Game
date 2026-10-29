# Kirchhoff's current and voltage laws

Kirchhoff's two laws are the constraints every circuit simulator writes down before it solves anything. The current law says charge does not pile up at a node; the voltage law says potential is single-valued, so drops around any loop cancel. Nodal analysis, mesh analysis and the modified nodal analysis inside SPICE are all systematic ways to write these equations, so a tool builder must know exactly how many independent equations a network yields.

## Kirchhoff's current law at a node

Kirchhoff's current law (KCL) states that the currents leaving a node sum to zero at every instant: $\sum_k i_k = 0$, with entering currents counted negative. It is charge conservation applied to a node that stores no charge; any capacitance is modeled as a separate branch.

Example: at a node, 12 mA enters through branch 1, 5 mA leaves through branch 2 and 4 mA enters through branch 3. The net inflow is $12 - 5 + 4 = 11$ mA, so the fourth branch must carry 11 mA away from the node. If the arithmetic gives a negative number, that branch's current actually flows in. A simulator never chooses directions physically; it assigns a reference direction to every branch and lets the solution's sign report the truth.

## Kirchhoff's voltage law around a loop

Kirchhoff's voltage law (KVL) states that the voltage drops around any closed loop sum to zero. It follows from the electric field being conservative in lumped circuits: the potential at a point does not depend on the path taken to reach it. Walking around a loop, a source contributes a rise and each resistor a drop in the direction of its current.

Example: a loop holds a 3 V source and three elements dropping 0.8 V, 1.1 V and an unknown $V_x$ in the direction of the source's rise. KVL gives $3 = 0.8 + 1.1 + V_x$, so $V_x = 1.1$ V. A negative answer means the unknown element is actually a rise. KVL fails when changing magnetic flux threads the loop; then the loop's inductance must be modeled explicitly, which is why lumped models break at high frequency.

## Counting independent equations

A connected network with $n$ nodes and $b$ branches gives exactly $n - 1$ independent KCL equations: the $n$ node equations sum to zero, because every branch current leaves one node and enters another, so one equation is redundant. It gives $b - n + 1$ independent KVL equations: pick a spanning tree of $n - 1$ branches, and each of the remaining branches closes exactly one fundamental loop.

Example: 6 nodes and 9 branches give 5 independent KCL equations and $9 - 6 + 1 = 4$ independent loops. Together, $5 + 4 = 9$ equations match the 9 unknown branch currents once element laws relate currents and voltages. Nodal analysis picks the $n - 1$ node voltages as unknowns; mesh analysis picks the $b - n + 1$ loop currents. A simulator prefers nodal unknowns because they are easy to assign automatically and stay sparse.

## Going deeper

Desoer and Kuh, "Basic Circuit Theory", chapters 1 and 9; Chua, Desoer and Kuh, "Linear and Nonlinear Circuits", on graph-theoretic network equations.
