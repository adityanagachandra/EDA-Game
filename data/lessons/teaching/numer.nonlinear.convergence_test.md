# Newton convergence criteria

A Newton loop needs a precise rule for when to stop. Stop too early and the answer is wrong; demand too much and the simulator wastes iterations or never terminates because of roundoff. Circuit simulators use tolerances that mix relative and absolute parts on both node voltages and currents, and reading the iterate history tells a tool designer whether the Jacobian is right.

## Update tolerances: reltol plus abstol

SPICE accepts a node voltage when
$$|v_{k+1} - v_k| \le \text{reltol}\cdot\max(|v_{k+1}|, |v_k|) + \text{vntol},$$
with defaults reltol $= 10^{-3}$ and vntol $= 1$ µV; branch currents use abstol $= 1$ pA in place of vntol. The relative part scales with the signal; the absolute part keeps nodes near zero from demanding impossible relative accuracy.

Example: $v_k = 1.2000$ V, $v_{k+1} = 1.2009$ V. The tolerance is $10^{-3} \cdot 1.2009 + 10^{-6} = 1.2019$ mV, and the change is 0.9 mV, so the node passes. A node at 0.1 mV would be held to about 1.1 µV, dominated by vntol. Every node must pass, so the worst node decides.

## Residual tolerances on the equations

The update test can pass while KCL is still violated, for example when a step was heavily damped. So simulators also check the residual: at each node the net current must satisfy
$$\Big|\sum i\Big| \le \text{reltol}\cdot\max|i_{branch}| + \text{abstol}.$$
With reltol $= 10^{-3}$ and the largest branch current at a node of 20 µA, the allowed imbalance is 20 nA plus 1 pA. A residual of 12 nA passes; 35 nA fails.

Scaling the tolerance by the largest branch current, not by the net current, matters: the net is near zero by design, so a purely relative test on it would never pass. ngspice and Xyce both use variants of these two tests; the exact rules differ by simulator and are part of each tool's documented options.

## Reading quadratic convergence from the iterates

With errors $e_k$, the order of convergence $p$ satisfies $e_{k+1} \approx Ce_k^p$, so
$$p \approx \frac{\ln(e_{k+1}/e_k)}{\ln(e_k/e_{k-1})}.$$
Errors $10^{-2}$, $2 \times 10^{-4}$, $8 \times 10^{-8}$ give $p \approx \ln(4 \times 10^{-4})/\ln(0.02) = 7.82/3.91 = 2.0$: quadratic, as Newton with an exact Jacobian should be. Errors shrinking by a constant factor each step give $p \approx 1$.

Linear convergence where quadratic is expected is a diagnostic. Common causes are a Jacobian that does not match the residual (a device model whose derivatives are wrong or inconsistent with its current equations), a reused stale Jacobian, or heavy damping still active near the root. This is exactly the kind of check a tool designer adds when validating a new compact model.

## Going deeper

Quarles, "SPICE3 version 3C1 users guide" (UC Berkeley, 1989) on convergence options; Kelley, "Iterative Methods for Linear and Nonlinear Equations" (SIAM, 1995), chapter 5.
