# Transient equations as a DAE

Transient analysis integrates the circuit equations in time, but those equations are not an ordinary differential equation. Nodes without capacitance, voltage sources and inductor constraints contribute algebraic equations, so the system is a differential-algebraic equation (DAE). Its structure explains why SPICE starts from a DC operating point, why it integrates charges rather than capacitances, and why it uses implicit methods that solve a Newton problem at every time point.

## Resistive and reactive parts

For a linear circuit, MNA gives $G x + C \dot{x} = b(t)$. The matrix $G$ holds the memoryless stamps: conductances, the $\pm 1$ incidence entries of voltage sources and inductors, controlled sources. The matrix $C$ holds the coefficients of time derivatives: capacitances in node rows and $-L$ in inductor branch rows, whose equation is $v_+ - v_- - L\, di/dt = 0$. For nonlinear devices the general form is $f(x) + \tfrac{d}{dt} q(x) = b(t)$, where $q$ collects device charges and inductor fluxes. Integrating charge directly conserves it: the discretized $C(v)\, dv/dt$ of a nonlinear capacitor does not equal the exact change of charge, so capacitance-based integration can create or destroy charge numerically.

## Algebraic constraints and singular C

A row of $C$ is zero wherever no derivative appears: a node with no capacitor, the branch row of a voltage source. Then $C$ is singular and those equations are algebraic, holding at every instant. Initial values therefore cannot be arbitrary: the algebraic equations must hold at $t = 0$, which is why SPICE computes a consistent starting point, the DC operating point, before integrating. Loops of capacitors and voltage sources, or cutsets of inductors and current sources, raise the DAE index: some unknowns then depend on derivatives of the inputs, and initialization and integration become ill-posed. Inverting $C$ to obtain an ODE is impossible when $C$ is singular and would destroy sparsity anyway.

## Stiffness and why SPICE is implicit

Circuits mix time constants across many decades: a 1 ps parasitic next to a 1 us bias filter. An explicit method such as forward Euler is stable only for steps of the order of the smallest time constant, so it would take $10^{7}$ steps to simulate 10 us, long after the fast mode has died out. Backward Euler and the trapezoidal rule are A-stable: every decaying mode stays bounded for any step size, so the step is chosen by accuracy alone. The price is that each step is an implicit nonlinear system, solved by Newton with the same sparse LU machinery as the DC operating point, with companion models replacing the derivatives.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 5; Brenan, Campbell and Petzold, "Numerical Solution of Initial-Value Problems in Differential-Algebraic Equations" (SIAM, 1996).
