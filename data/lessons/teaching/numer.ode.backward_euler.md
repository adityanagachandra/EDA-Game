# Backward Euler

Transient circuit simulation integrates the circuit equations in time, and the workhorse integrator is backward Euler. It evaluates the derivative at the end of the step, so every step is an equation to solve rather than an explicit update. That costs a linear or Newton solve per step but buys stability on the stiff systems that circuits always are.

## The implicit update

For $y' = f(y, t)$ with step $h$, backward Euler sets
$$y_{n+1} = y_n + h\,f(y_{n+1}, t_{n+1}).$$
The unknown $y_{n+1}$ appears on both sides. For the linear equation $y' = -ay + b$ it can be solved directly:
$$y_{n+1} = \frac{y_n + hb}{1 + ah}.$$
With $a = 2$, $b = 4$, $y_0 = 5$ and $h = 0.3$: $y_1 = (5 + 1.2)/1.6 = 3.875$, moving toward the equilibrium $b/a = 2$.

Forward Euler would use the slope at $y_n$ instead, $y_1 = 5 + 0.3(-10 + 4) = 3.2$. For a nonlinear circuit, the implicit equation is solved by Newton's method at each timestep, so a transient run is a sequence of Newton solves, each one a sequence of sparse LU factorizations.

## One step for an RC circuit

An RC low-pass with input $u$ obeys $\tau v' = u - v$, $\tau = RC$. Backward Euler gives
$$v_{n+1} = \frac{v_n + (h/\tau)\,u_{n+1}}{1 + h/\tau}.$$
With $R = 2$ kΩ and $C = 0.5$ pF, $\tau = 1$ ns. From $v_0 = 0$ with a 1 V step input and $h = 0.25$ ns: $v_1 = 0.25/1.25 = 0.2$ V. The exact value is $1 - e^{-0.25} = 0.2212$ V, so one step loses about 0.02 V, the first-order error of the method.

The update factor $1/(1 + h/\tau)$ lies between 0 and 1 for every positive $h$, however large. A huge step does not blow up; it just jumps to the steady state $v = u$. That unconditional stability is what simulators need.

## Backward Euler on a coupled system

For a linear system $x' = Ax + b$, backward Euler gives
$$(I - hA)\,x_{n+1} = x_n + hb.$$
With $A = \begin{pmatrix} -3 & 1 \\ 1 & -2 \end{pmatrix}$, $b = 0$, $x_0 = (1, 0)$ and $h = 0.5$:
$$I - hA = \begin{pmatrix} 2.5 & -0.5 \\ -0.5 & 2 \end{pmatrix}, \quad \det = 4.75,$$
so $x_1 = (2/4.75,\ 0.5/4.75) = (0.4211, 0.1053)$.

In MNA form the system is $Cx' + Gx = u$, and backward Euler becomes $(C/h + G)x_{n+1} = (C/h)x_n + u_{n+1}$. The matrix $C/h + G$ has the same pattern at every step, so the symbolic factorization is reused, and when $h$ is unchanged and the circuit is linear even the numeric factors are reused.

## Going deeper

Hairer and Wanner, "Solving Ordinary Differential Equations II: Stiff and Differential-Algebraic Problems" (2nd ed., Springer, 1996), chapter IV.3; Najm, "Circuit Simulation" (Wiley, 2010), chapter 4.
