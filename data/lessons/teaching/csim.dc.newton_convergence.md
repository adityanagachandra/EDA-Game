# Newton convergence tests

Newton never reaches the exact solution; a simulator must decide when to stop. SPICE's convergence tests combine relative and absolute tolerances on both node voltages and device currents, and their defaults are worth memorizing because they set the accuracy of every result. Understanding how quickly Newton converges near the solution also tells you how many iterations a healthy operating point should take.

## The voltage tolerance band

SPICE accepts a node voltage when successive iterates agree within a band:
$$|v_k - v_{k-1}| \le \text{RELTOL} \cdot \max(|v_k|, |v_{k-1}|) + \text{VNTOL},$$
with defaults RELTOL $= 10^{-3}$ and VNTOL $= 1$ uV. The relative term dominates for ordinary voltages; VNTOL keeps the test meaningful near 0 V, where a pure relative test could never pass. For iterates 1.2000 V and 1.1995 V, the band is $1.2 \times 10^{-3} + 10^{-6}$ V, about 1201 uV, and the change of 500 uV passes. Every node must pass, and the test is applied on consecutive iterates, so at least two iterations are always needed.

## Quadratic convergence in iterations

Near a simple root with a smooth Jacobian, Newton converges quadratically: $e_{k+1} \approx C e_k^2$, where $e_k$ is the error and $C$ depends on the curvature. Each iteration roughly doubles the number of correct digits. With $C = 2$ per volt and $e_0 = 0.1$ V, the errors are 0.02, $8 \times 10^{-4}$, $1.3 \times 10^{-6}$, $3.3 \times 10^{-12}$: four iterations reach picovolts. This is why a converging operating point typically finishes in a few iterations once it is close, and why a run that needs hundreds of iterations signals a problem: the iterate is far from the root, being damped or limited, or the Jacobian is wrong.

## The current tolerance

SPICE also checks device currents, because a voltage can settle while a steep exponential current is still far off. The test is $|i_k - i_{k-1}| \le \text{RELTOL} \cdot \max(|i_k|, |i_{k-1}|) + \text{ABSTOL}$, with ABSTOL $= 1$ pA by default. For a branch current near 0.124 mA, the band is $0.124 \times 10^{-3}$ mA plus 1 pA, about 124 nA. ABSTOL matters for circuits biased at tiny currents: a leakage-dominated node carrying 10 pA has a band of about 1 pA, dominated by ABSTOL. Raising ABSTOL to speed convergence on a low-current circuit quietly accepts currents that are wrong by more than their own size.

## Going deeper

The ngspice manual, the section on convergence options; Najm, "Circuit Simulation" (Wiley, 2010), chapter 4.
