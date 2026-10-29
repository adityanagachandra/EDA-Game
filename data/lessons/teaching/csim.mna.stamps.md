# Modified nodal analysis stamps

Modified nodal analysis (MNA) is how SPICE turns a netlist into a linear system. Each element adds a small, fixed pattern of entries, its stamp, to the matrix and the right-hand side, so the matrix is assembled by walking the element list once. Stamping is the heart of every simulator's load phase: Newton, transient and AC analyses all rebuild the same structure with new values. Getting a single entry right, by hand, is the test of whether you understand the formulation.

## The resistor stamp

Write KCL at each non-ground node: the currents leaving through elements equal the source currents injected. A resistor with conductance $g = 1/R$ between nodes $a$ and $b$ carries $g(v_a - v_b)$ out of $a$ and the negative of that out of $b$. Its stamp therefore adds $+g$ at $(a,a)$ and $(b,b)$ and $-g$ at $(a,b)$ and $(b,a)$. A diagonal entry $G_{ii}$ is the sum of the conductances of every resistor touching node $i$, including resistors to ground. With 1 k$\Omega$ from node 1 to ground and 2 k$\Omega$ from node 1 to node 2, $G_{11} = 1 + 0.5 = 1.5$ mS. Working in kilohms gives conductances in millisiemens, currents in milliamperes and voltages in volts.

## Off-diagonal entries and the ground row

The off-diagonal entry $G_{ij}$ is minus the total conductance connected directly between nodes $i$ and $j$, and zero if none is. Parallel resistors between the same pair add. Ground is the reference: its row and column are deleted, because its voltage is known, so a resistor from node $a$ to ground stamps only its diagonal $+g$ at $(a,a)$. For a resistive network the result is symmetric and diagonally dominant. With 4 k$\Omega$ and 1 k$\Omega$ in parallel between nodes 2 and 3, $G_{23} = G_{32} = -(0.25 + 1) = -1.25$ mS. A missing connection gives a structural zero, which the sparse solver never stores.

## Assembling the whole matrix

The assembled matrix is the sum of all stamps, which suggests a check. Each resistor between two non-ground nodes contributes $+g, +g, -g, -g$, summing to zero; each resistor to ground contributes a single $+g$. So the sum of all entries of $G$ equals the total conductance from the circuit to ground. A network with 0.5 k$\Omega$ from node 1 to ground, 2 k$\Omega$ from node 3 to ground and any number of resistors between nodes 1, 2 and 3 has entry sum $2 + 0.5 = 2.5$ mS. Simulators assemble in exactly this order: zero the matrix, loop over elements, add each stamp. Because entries only accumulate, the loop parallelizes by element with care for collisions on shared nodes.

## Going deeper

Ho, Ruehli and Brennan, "The modified nodal approach to network analysis" (IEEE Transactions on Circuits and Systems, 1975); Najm, "Circuit Simulation" (Wiley, 2010), chapter 2.
