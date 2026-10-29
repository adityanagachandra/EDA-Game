# IR drop analysis

Static IR-drop analysis answers one question: with each cell drawing its average current, how far below the supply does each node sit? It is one sparse linear solve, $G d = i$, over a grid with up to billions of nodes. The answer gates signoff, because cells on a drooping supply are slower than the timer assumed. Solving small instances by hand makes the structure concrete before scaling it up.

## Series paths

When the grid is a tree, currents can be added by KCL and drops read off directly. A pad feeds node a through $0.9\ \Omega$ and a feeds b through $1.7\ \Omega$; a draws 18 mA and b draws 18 mA. The first resistor carries 36 mA and the second 18 mA, so b drops $0.9 \cdot 36 + 1.7 \cdot 18 = 63$ mV. Ohms times milliamps give millivolts, a convenient unit system. This is the same as solving the 2 by 2 system $G d = i$, but the tree structure lets back-substitution happen by inspection.

## A doubly fed rail

With supplies at both ends, current splits and the system must be solved. For pad, $0.7\ \Omega$, n1, $2\ \Omega$, n2, $0.7\ \Omega$, n3, $1.7\ \Omega$, pad, with loads of 3, 9 and 4 mA, $G$ is tridiagonal with diagonal $(1/0.7 + 1/2,\ 1/2 + 1/0.7,\ 1/0.7 + 1/1.7)$ and off-diagonals $-1/2$ and $-1/0.7$. Solving gives drops of 5.71, 16.02 and 13.33 mV. The heaviest drop is not always at the heaviest load; it is where load meets high resistance to every supply.

## A small mesh

In a 2 by 2 grid with branches n1-n2 $2\ \Omega$, n2-n4 $1\ \Omega$, n3-n4 $4\ \Omega$ and n1-n3 $1\ \Omega$, supplies at n1 and n4 through $0.25\ \Omega$ each, and loads of 9, 7, 10 and 8 mA, the 4 by 4 system gives drops of 4.71, 8.76, 12.53 and 3.79 mV. The worst node, n3, is the one whose path to both supplies runs through the resistive n3-n4 branch. Production tools solve the same system with millions of nodes using sparse Cholesky or preconditioned conjugate gradient, and then map each cell to its local drop for timing.

## Going deeper

Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the chapter on power and ground routing; Najm, "Circuit Simulation" (2010).
