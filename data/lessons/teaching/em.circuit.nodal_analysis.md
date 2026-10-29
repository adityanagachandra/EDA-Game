# Nodal analysis

Nodal analysis chooses the node voltages as unknowns and writes KCL at every non-reference node. It is the method behind every SPICE-class simulator, because it can be assembled automatically from a netlist and produces a sparse, often symmetric, matrix. Solving two- and three-node networks by hand builds the intuition for what the simulator's matrix contains and why its structure matters.

## Writing one node equation

Pick a reference node (ground) and name the voltage of every other node relative to it. At each non-reference node, express every branch current leaving it in terms of node voltages, using $(V_a - V_b)/R$ for a resistor between $a$ and $b$, and add any current source. Then set the sum to zero.

Example: a 1.2 V source feeds node A through $R_1 = 2$ kohm; $R_2 = 4$ kohm runs from A to ground; a current source injects 0.3 mA into A. KCL at A: $(V - 1.2)/2 + V/4 = 0.3$. Multiplying out, $0.75V = 0.9$, so $V = 1.2$ V. Working in kohm and mA keeps the numbers in volts with no powers of ten. Every term has the form conductance times a voltage difference, which is what makes the method mechanical.

## Solving a two-node network

With two unknown nodes the KCL equations form a 2 by 2 linear system $G v = i$. The conductance matrix $G$ has the sum of conductances touching node $k$ on its diagonal and minus the conductance between nodes $j$ and $k$ off the diagonal; the right side holds source currents, including a voltage source seen through its series resistor as $V_s/R$.

Example: $V_s = 2$ V through 1 kohm to node 1, 1 kohm from node 1 to ground, 1 kohm from node 1 to node 2 and 1 kohm from node 2 to ground. In mS: $3V_1 - V_2 = 2$ and $-V_1 + 2V_2 = 0$. Then $V_1 = 2V_2$, so $5V_2 = 2$, $V_2 = 0.4$ V and $V_1 = 0.8$ V. Cramer's rule or elimination both work at this size; a simulator uses sparse LU.

## Three-node networks

A ladder of three nodes gives a tridiagonal 3 by 3 system: each node couples only to its neighbours. That sparsity is the norm in real circuits, where a node connects to a handful of others no matter how big the netlist, and it is why simulators scale to millions of nodes.

Example: a 3 V source drives 1 kohm to node 1, then 1 kohm to node 2, then 1 kohm to node 3, with 1 kohm from each node to ground. In mS the system is $3V_1 - V_2 = 3$, $-V_1 + 3V_2 - V_3 = 0$, $-V_2 + 2V_3 = 0$. Back-substituting $V_2 = 2V_3$ gives $V_1 = 5V_3$, and the first equation $15V_3 - 2V_3 = 3$ gives $V_3 = 0.231$ V, $V_2 = 0.462$ V, $V_1 = 1.154$ V. Gaussian elimination on a tridiagonal matrix costs linear time, the Thomas algorithm.

## Going deeper

Vlach and Singhal, "Computer Methods for Circuit Analysis and Design" (2nd ed.), chapters 2 and 3; Najm, "Circuit Simulation", chapter 2.
