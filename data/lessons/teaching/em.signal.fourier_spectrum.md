# Fourier spectra of digital signals

A digital signal is not a single frequency. Its spectrum, given by Fourier analysis, decides how much bandwidth every interconnect model, package model and measurement must cover. The harmonics of a clock, the knee frequency of an edge and the corners of a trapezoidal spectrum tell a tool designer what frequency range an extracted model must be accurate over, and where emissions will concentrate.

## Harmonics of a square wave

A periodic signal is a sum of harmonics at multiples of its fundamental $f_0 = 1/T$. An ideal 50% duty square wave from 0 to $V$ has a DC term $V/2$ and only odd harmonics, with amplitudes
$$c_n = \frac{2V}{n\pi}, \quad n = 1, 3, 5, \ldots$$
Even harmonics vanish because of the half-wave symmetry.

Example: a 1 V clock at 200 MHz has a 200 MHz component of $2/\pi = 0.637$ V, a 600 MHz component of 0.212 V and a 1 GHz component of 0.127 V, and nothing at 400 or 800 MHz. Amplitudes fall as $1/n$, slowly, so an ideal square wave has unbounded bandwidth. A duty cycle away from 50% brings the even harmonics back, which is one way emissions scans reveal a duty-cycle problem.

## The knee frequency

Real edges have finite rise time, which rolls off the high harmonics. A widely used rule of thumb from high-speed design places the knee frequency, below which most of an edge's energy lies, at
$$f_{knee} \approx \frac{0.5}{t_r},$$
with $t_r$ the 10-90 rise time. It is a guideline, not a theorem, and it is deliberately conservative relative to the 0.35/$t_r$ single-pole bandwidth.

Example: a 50 ps edge has $f_{knee} = 10$ GHz. A package or channel model that is accurate to 10 GHz will reproduce the edge's shape well; one fitted to 2 GHz will round the edge and miss reflections from small discontinuities. The clock frequency does not appear: a 100 MHz clock with 50 ps edges needs the same model bandwidth as a 2 GHz clock with 50 ps edges, which surprises people.

## The trapezoidal spectral envelope

A periodic trapezoidal pulse with pulse width $\tau$ and rise time $t_r$ has a spectral envelope with two corners. Below $1/(\pi\tau)$ the envelope is flat; between $1/(\pi\tau)$ and $1/(\pi t_r)$ it falls at 20 dB per decade; above $1/(\pi t_r)$ it falls at 40 dB per decade.

Example: a 5 ns pulse with 100 ps edges has corners at $1/(\pi \cdot 5 \text{ ns}) = 63.7$ MHz and $1/(\pi \cdot 100 \text{ ps}) = 3.18$ GHz. Halving the rise time moves the second corner to 6.37 GHz and raises every harmonic above the old corner by up to 6 dB. That is why electromagnetic compatibility problems often appear when a design is moved to a faster process without any change to its clock: the edges got faster. Emissions estimation tools use this envelope to bound spectra quickly.

## Going deeper

Paul, "Introduction to Electromagnetic Compatibility" (2nd ed.), chapter 3; Johnson and Graham, "High-Speed Digital Design: A Handbook of Black Magic", chapter 1; Ott, "Electromagnetic Compatibility Engineering", chapter 12.
