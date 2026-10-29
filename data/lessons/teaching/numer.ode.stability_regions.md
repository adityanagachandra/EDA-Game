# Stability regions of integration methods

Whether an integration method is safe on a circuit depends on the eigenvalues of the circuit's linearization and on the method's stability region. The scalar test equation $y' = \lambda y$ captures this exactly. Comparing forward Euler, backward Euler and trapezoidal on it explains why simulators are implicit and how each method distorts decaying and oscillating modes.

## The test equation and the amplification factor

Apply a one-step method to $y' = \lambda y$ with step $h$. Every method in this family gives $y_{n+1} = R(z)\,y_n$ with $z = h\lambda$ and an amplification factor $R$:
- forward Euler: $R(z) = 1 + z$;
- backward Euler: $R(z) = 1/(1 - z)$;
- trapezoidal: $R(z) = (1 + z/2)/(1 - z/2)$.

The stability region is the set of $z$ with $|R(z)| \le 1$, where the numerical solution does not grow. At $z = -3$: forward Euler gives $R = -2$ (unstable, oscillating growth), backward Euler 0.25, trapezoidal -0.2. For a linear circuit $x' = Ax$, the same analysis applies to each eigenvalue of $A$, so every circuit mode must sit inside the region.

## Comparing forward Euler, backward Euler and trapezoidal

Forward Euler's region is the disk $|1 + z| \le 1$. For a real decaying mode, $\lambda < 0$, it needs $h \le 2/|\lambda|$: a mode with $\lambda = -20$ per ns limits the step to 100 ps, however slowly the waveform of interest changes.

Backward Euler is stable everywhere outside the disk $|1 - z| < 1$, which contains the whole left half-plane: it is A-stable, and since $R \to 0$ as $z \to -\infty$ it is also L-stable. The trapezoidal rule is stable exactly on the closed left half-plane: A-stable, but $R \to -1$ for very stiff modes, so not L-stable.

Example: $\lambda = -10$ per ns, $h = 0.5$ ns, $z = -5$. Forward Euler: $|R| = 4$, divergent. Backward Euler: $1/6 = 0.167$. Trapezoidal: $|{-1.5}/3.5| = 0.429$.

## Oscillatory modes and numerical damping

An undamped LC tank has $\lambda = \pm i\omega$, so $z = i\omega h$ on the imaginary axis. Then $|R_{FE}| = \sqrt{1 + (\omega h)^2} > 1$: forward Euler makes every oscillation grow. $|R_{BE}| = 1/\sqrt{1 + (\omega h)^2} < 1$: backward Euler damps it artificially. $|R_{TR}| = 1$: trapezoidal preserves the amplitude exactly (with a phase error).

With $\omega h = 0.2$ over 50 steps, backward Euler shrinks the amplitude to $(1.04)^{-25} = 0.375$ of its true value, a serious error when simulating an oscillator or a PLL, while trapezoidal keeps it at 1. This is the trade-off a simulator designer weighs: trapezoidal for oscillators and accuracy, backward Euler or BDF where stiff modes and ringing matter more, with step control keeping $\omega h$ small either way.

## Going deeper

Hairer and Wanner, "Solving Ordinary Differential Equations II" (2nd ed., Springer, 1996), chapter IV.2; LeVeque, "Finite Difference Methods for Ordinary and Partial Differential Equations" (SIAM, 2007), chapter 7.
