# Harmonic balance

Harmonic balance (HB) finds the periodic steady state in the frequency domain. Each circuit unknown is represented by its Fourier coefficients up to some harmonic, linear elements are stamped harmonic by harmonic exactly as in AC analysis, and nonlinear devices are evaluated in the time domain and transformed back with FFTs. Newton then balances the currents at every harmonic. HB is the workhorse of RF simulation, especially for mildly nonlinear circuits with frequency-domain models.

## Unknowns in the frequency domain

Truncate each waveform to a DC term and harmonics 1 to $K$: $x(t) = X_0 + \sum_{k=1}^{K} (X^c_k \cos k\omega t + X^s_k \sin k\omega t)$. Each of the circuit's $N$ unknowns then carries $2K + 1$ real coefficients, so the HB system has $N(2K + 1)$ real unknowns. A 100-unknown circuit with 10 harmonics has 2100. Linear elements stamp $G + jk\omega C$ at each harmonic independently, but every nonlinear device couples all harmonics of its terminals, so the Jacobian is block-dense in the harmonic dimension. That is why large HB problems are solved with GMRES and preconditioners rather than direct LU.

## Harmonics of a quadratic nonlinearity

A nonlinear device generates harmonics of its drive. For $i = a_1 v + a_2 v^2$ with $v = A\cos\omega t$, the identity $\cos^2 = \tfrac12 + \tfrac12\cos 2\omega t$ gives
$$i = \tfrac{a_2 A^2}{2} + a_1 A \cos\omega t + \tfrac{a_2 A^2}{2}\cos 2\omega t.$$
The quadratic term produces a DC shift and a second harmonic, each of amplitude $a_2 A^2/2$. With $a_1 = 3$ mS, $a_2 = 1.5$ mA/V$^2$ and $A = 0.2$ V, the fundamental is 0.6 mA and the DC and second harmonic are 0.03 mA. HB obtains these coefficients numerically: sample the waveform at $2K + 1$ or more points per period, evaluate the device, and FFT.

## Harmonics of a cubic nonlinearity

A cubic term $a_3 v^3$ uses $\cos^3 = \tfrac34\cos + \tfrac14\cos 3$:
$$i = \left(a_1 A + \tfrac34 a_3 A^3\right)\cos\omega t + \tfrac14 a_3 A^3 \cos 3\omega t.$$
With a negative $a_3$ the fundamental shrinks as the drive grows: gain compression. For $a_1 = 4.5$ mS, $a_3 = -0.1$ mA/V$^3$ and $A = 0.6$ V, the fundamental is $2.7 - 0.0162 = 2.684$ mA and the third harmonic is -0.0054 mA. With two input tones the same cubic term creates intermodulation products at $2f_1 - f_2$, which multitone HB captures directly. Truncating $K$ too low aliases the harmonics a strong nonlinearity generates, so HB suits mild nonlinearity best.

## Going deeper

Kundert, White and Sangiovanni-Vincentelli, "Steady-State Methods for Simulating Analog and Microwave Circuits" (Kluwer, 1990); Maas, "Nonlinear Microwave and RF Circuits" (2nd ed., Artech House, 2003).
