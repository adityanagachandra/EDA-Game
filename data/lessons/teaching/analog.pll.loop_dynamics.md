# Charge-pump PLL loop dynamics

The charge-pump PLL with a series RC loop filter is the workhorse clock synthesizer. Linearized, it is a second-order type-II feedback system whose natural frequency and damping factor set lock time, overshoot and noise filtering. Computing them from charge-pump current, VCO gain, divide ratio and filter values is the core of PLL design and the first check a PLL model must reproduce.

## Natural frequency

The phase-frequency detector and charge pump deliver an average current $I_{CP}\theta_e/(2\pi)$ for phase error $\theta_e$. The filter impedance is $R + 1/(sC)$, the VCO is $K_{VCO}/s$ in rad/s per volt and the divider $1/N$. Ignoring $R$ for the moment, the loop gain has two integrators and the closed loop has natural frequency $\omega_n = \sqrt{I_{CP}K_{VCO}/(2\pi NC)}$. With $K_{VCO}$ given in Hz/V, the $2\pi$ cancels: $\omega_n = \sqrt{I_{CP}K/(NC)}$. For $I_{CP} = 100\ \mu$A, $K = 500$ MHz/V, $N = 100$ and $C = 100$ pF, $\omega_n = \sqrt{5 \times 10^{12}} = 2.24 \times 10^6$ rad/s, or $f_n = 356$ kHz.

## Damping and the stabilizing zero

Two integrators alone give a loop with -180 degrees of phase and no damping: it would oscillate. The resistor $R$ adds a zero at $1/(RC)$ that restores phase, and the damping factor becomes $\zeta = (RC/2)\omega_n$. Continuing the example with $R = 10$ kilohms, $RC/2 = 5 \times 10^{-7}$ s and $\zeta = 1.12$. To hit a target, solve $R = 2\zeta/(C\omega_n)$: for $\zeta = 0.707$, $R = 6.3$ kilohms. Damping between about 0.7 and 1 gives a fast, well-behaved lock transient. In practice a small capacitor in parallel with the RC branch suppresses ripple on the control line, adding a third pole that must sit well above the loop bandwidth.

## What the loop bandwidth trades

The loop bandwidth, of the order of $\omega_n$ for moderate damping, decides who wins in the noise budget. Inside the bandwidth the output follows the reference, so reference, detector and divider noise are passed through (multiplied by $N$); outside it the VCO runs free, so its own phase noise dominates. A wide loop suppresses VCO noise and locks fast but passes more reference and charge-pump noise and fractional-N quantization noise. The loop is also a sampled system: the phase detector updates once per reference cycle, so the continuous-time model is only valid for bandwidths well below $f_{ref}$, a common rule being about a tenth. These trade-offs make the loop parameters a small optimization problem that PLL design tools solve before any transistor is sized.

## Going deeper

Gardner, "Charge-pump phase-lock loops" (IEEE Transactions on Communications, 1980); Razavi, "Design of CMOS Phase-Locked Loops" (2020).
