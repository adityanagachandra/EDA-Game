# Differential-algebraic equations in circuits

The equations a circuit simulator integrates are not ordinary differential equations. Nodes without capacitance and the branch currents of voltage sources give purely algebraic equations, so MNA produces a differential-algebraic system. Its index, its consistent initial conditions and the choice of integrator follow from that structure, and a simulator designer must handle all three.

## Singular capacitance matrices make DAEs

MNA writes $C\,x' + G\,x = u(t)$. Any row of $C$ that is zero, a node with no capacitor or a voltage-source current equation, is an algebraic constraint holding at every instant. When $C$ is singular the system cannot be rewritten as $x' = C^{-1}(\dots)$: it is a DAE.

The number of algebraic equations is $n - \text{rank}(C)$. A circuit with 10 MNA unknowns, capacitors at 6 nodes and no others has $\text{rank}(C) = 6$ and 4 algebraic equations. A resistive node between two capacitive nodes has its voltage fixed instantaneously by its neighbours, like an infinitely fast mode. Explicit ODE integrators cannot even be applied; implicit ones handle the algebraic rows naturally because they solve all equations at the new time point.

## Index and consistent initial conditions

The differentiation index counts how many times the algebraic equations must be differentiated to get an ODE for every variable. Typical circuits are index 1. Loops of capacitors and voltage sources, and cutsets of inductors and current sources, create index 2: the constraint restricts the differential variables themselves, hidden constraints appear, and numerical errors are amplified.

Initial values must satisfy every constraint. If they do not, an implicit integrator forces the state onto the constraint in the first step, and the derivative of that jump shows up as a spurious current impulse. SPICE obtains consistent values by solving the DC operating point first, where all derivatives are zero and every algebraic equation holds; the UIC option skips it and leaves consistency to the user.

## Why BDF methods handle index-1 DAEs

Backward Euler on the DAE gives $(C/h + G)\,x_{n+1} = (C/h)\,x_n + u_{n+1}$. Only the pencil $C/h + G$ must be nonsingular; $C$ itself never needs inverting. The algebraic rows of this system are just the constraints at $t_{n+1}$, so every accepted point satisfies them. BDF methods generalize this, and their damping of infinitely stiff components makes algebraic variables settle immediately.

The trapezoidal rule treats algebraic variables as infinitely stiff modes with amplification -1: an error in them alternates sign every step without decaying, which appears as ringing in currents of resistive nodes or sources. That is one more reason simulators insert backward Euler steps after breakpoints, and why general-purpose DAE codes such as DASSL and IDA are built on BDF.

## Going deeper

Brenan, Campbell and Petzold, "Numerical Solution of Initial-Value Problems in Differential-Algebraic Equations" (SIAM, 1996); Günther and Feldmann, "CAD-based electric-circuit modeling in industry", Surveys on Mathematics for Industry (1999).
