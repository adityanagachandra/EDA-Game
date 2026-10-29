# Modified nodal analysis and element stamps

Plain nodal analysis cannot express an ideal voltage source or an inductor's DC short, so SPICE-class simulators use modified nodal analysis (MNA): node voltages plus selected branch currents as unknowns. The matrix is built element by element, each device adding a small pattern of entries called its stamp. Stamps are the interface between device models and the linear solver, which makes them the core data structure a simulator writer must master.

## Stamping conductances

A conductance $g$ between nodes $a$ and $b$ adds $+g$ to entries $(a,a)$ and $(b,b)$ and $-g$ to $(a,b)$ and $(b,a)$. A conductance from $a$ to ground stamps only $+g$ at $(a,a)$, because the ground row and column are deleted. Summing all stamps gives the nodal conductance matrix, which for a resistive network is a weighted graph Laplacian with the ground row removed: symmetric, diagonally dominant, and positive definite when every node has a path to ground.

Example: 1 kohm from node 1 to ground, 2 kohm between nodes 1 and 2 and 4 kohm from node 2 to ground. In mS: $G_{11} = 1 + 0.5 = 1.5$, $G_{22} = 0.5 + 0.25 = 0.75$ and $G_{12} = G_{21} = -0.5$. An entry between two nodes with no element between them is zero, which is where the sparsity comes from.

## Stamping voltage sources and inductors

An ideal voltage source fixes $V_a - V_b = E$ but does not say what current flows, so MNA adds its branch current $i_E$ as a new unknown. The stamp puts $+1$ and -1 in the new column at rows $a$ and $b$ (the current enters KCL there), and a new row $V_a - V_b = E$. An inductor is treated the same way, with its branch current as an unknown, which handles its DC short and its $v = L\,di/dt$ relation cleanly.

So the MNA system has one unknown per non-ground node plus one per voltage source and per inductor. Example: 11 nodes including ground, 3 voltage sources and 2 inductors give $10 + 3 + 2 = 15$ unknowns. Resistors, capacitors and current sources add no unknowns. The extra rows have zeros on the diagonal, so the matrix is no longer positive definite and the solver must pivot.

## Companion models for capacitors

In transient analysis the integration method turns each capacitor into a resistive companion model at every time step. Backward Euler approximates $i = C\,dv/dt$ by $i_n = (C/h)(v_n - v_{n-1})$: a conductance $G_{eq} = C/h$ in parallel with a current source $I_{eq} = (C/h) v_{n-1}$ that carries the history. The simulator stamps $G_{eq}$ like a resistor and $I_{eq}$ into the right-hand side, then solves a purely algebraic system.

Units on chip are convenient: femtofarads over picoseconds give millisiemens. Example: 100 fF at $h = 2$ ps gives $G_{eq} = 50$ mS; with $v_{n-1} = 0.6$ V the history source is 30 mA. Small steps make $G_{eq}$ large, which dominates the diagonal. Changing the step changes the matrix, which is why simulators refactor when they adapt $h$; the trapezoidal rule uses $2C/h$ with a different history term.

## Going deeper

Ho, Ruehli and Brennan, "The modified nodal approach to network analysis", IEEE Transactions on Circuits and Systems, 1975; Najm, "Circuit Simulation", chapters 2 and 4; the ngspice manual on transient analysis.
