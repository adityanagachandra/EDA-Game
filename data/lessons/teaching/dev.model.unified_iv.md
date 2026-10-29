# Single-piece I-V interpolation

A compact model must describe weak inversion (exponential current) and strong inversion (square-law or velocity-limited current) with one smooth expression, so that Newton iterations never see a kink. Charge-based models do this with an interpolation function, $\ln^2(1 + e^x)$, that approaches each asymptote at the right end. Evaluating it and its limits shows concretely how a single-piece model replaces piecewise regions.

## A single-piece interpolation

EKV-style models write the forward current as
$$I_D = I_{spec}\,F(x), \qquad F(x) = \ln^2\!\left(1 + e^{x}\right), \qquad x = \frac{V_{GS} - V_t}{2n\,kT/q},$$
where $I_{spec}$ is a specific current proportional to $W/L$, mobility, oxide capacitance and $(kT/q)^2$. (The complete model also subtracts a reverse term in the drain voltage, which vanishes in saturation.)

Example: $n = 1.3$, overdrive 0 V. Then $x = 0$ and $F = (\ln 2)^2 = 0.480$. At overdrive 100 mV, $x = 0.1/(2 \times 1.3 \times 0.02585) = 1.488$, $\ln(1 + e^{1.488}) = 1.689$ and $F = 2.86$.

The function is infinitely differentiable everywhere, so the current, transconductance and all higher derivatives are continuous through threshold. There is no region switch for the simulator to stumble on, and the moderate-inversion region, where many low-power analog circuits operate, is described rather than patched over.

## The strong-inversion limit

For large positive $x$, $\ln(1 + e^x) \approx x$, so $F \approx x^2$. Substituting $x$ gives
$$I_D \approx I_{spec}\frac{(V_{GS} - V_t)^2}{(2n\,kT/q)^2},$$
which is the square law with an effective $k'/n$.

The approach is gradual. At overdrive 150 mV with $n = 1.5$, $x = 1.934$, $F = 4.282$ and $x^2 = 3.741$: the ratio is 1.14, still 14% above the square law. At 300 mV, $x = 3.87$ and the ratio falls to about 1.01.

That gradual approach is the moderate-inversion region, roughly a few $nkT/q$ either side of threshold, where neither asymptote is accurate. Piecewise models that switch from exponential to square law at a fixed point get this region wrong, which matters for $g_m/I_D$-based analog design and for near-threshold digital logic. Real models then add velocity saturation and mobility degradation on top of this core.

## The weak-inversion limit

For large negative $x$, $\ln(1 + e^x) \approx e^x$, so $F \approx e^{2x}$, and
$$I_D \approx I_{spec}\,e^{(V_{GS} - V_t)/(n\,kT/q)}.$$
This is the subthreshold exponential with swing $n(kT/q)\ln 10$.

Example: $n = 1.5$, overdrive -140 mV. $x = -1.805$, $F = 0.0232$ and $e^{2x} = 0.0270$; the ratio is 0.857. Deeper, at -300 mV, the ratio is about 0.98.

So one function covers both limits, with the transition placed by $V_t$ and its width set by $n\,kT/q$. For an extraction tool this structure is useful: subthreshold data constrain $n$ and $V_t$, strong-inversion data constrain $I_{spec}$ (mobility and $C_{ox}$), and the moderate region tests whether the interpolation is right. For a simulator developer it shows why compact models are smooth by construction rather than by careful patching of region boundaries.

## Going deeper

Enz, Krummenacher and Vittoz, "An analytical MOS transistor model valid in all regions of operation", Analog Integrated Circuits and Signal Processing (1995); Enz and Vittoz, "Charge-Based MOS Transistor Modeling" (2006).
