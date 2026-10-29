# Behavioral macromodels

Not every block can be simulated at transistor level or reduced from an exact netlist. Packages and boards arrive as measured S-parameters; op-amps and regulators are replaced by behavioral models to speed system simulation. Macromodeling turns such data or abstractions into models a simulator can integrate. The engineering questions are always the same: is the model stable, passive and causal, and over what range is it valid?

## Rational fits of measured data

Tabulated frequency data define a time response only through convolution, which is costly at every step and sensitive to the measured band's limits. Vector fitting instead fits a rational function $H(s) \approx d + \sum_i r_i/(s - p_i)$ by iteratively relocating poles, giving a state-space model or an equivalent circuit that a simulator integrates directly. The fitted poles must lie in the left half plane: a pole with positive real part is a growing mode that makes transient diverge, so fitting algorithms flip or constrain such poles. Like a PRIMA model, a fitted model is stamped as a small dense block, because its states couple all ports.

## Passivity and causality

A model of a passive structure must itself be passive: for an admittance model, the Hermitian part of $Y(j\omega)$ must be positive semidefinite at every frequency. A stable but non-passive model can make an interconnection of stable parts oscillate, because it can supply energy. Passivity is checked by testing eigenvalues of the real part over frequency or with Hamiltonian-matrix tests, and enforced by small perturbations of the residues. Causality requires the impulse response to vanish before $t = 0$; band-limited or inconsistent data can produce acausal fits. PRIMA guarantees passivity by construction, through congruence; data-driven fits must check and enforce it.

## Behavioral models in a netlist

Behavioral macromodels, often written in Verilog-A, replace transistor-level blocks with equations: an op-amp as gain, a pole and output limits; a digital block as a table of supply current versus time for power-grid analysis. They run far faster, but they model only what they encode, so a fixed current signature does not respond to supply droop and a one-pole op-amp does not slew realistically. Numerically, behavioral code must keep currents, charges and their derivatives continuous: an abrupt if-statement creates a discontinuity that makes Newton oscillate and forces tiny steps. Every macromodel should state its validity range of temperature, load and frequency.

## Going deeper

Gustavsen and Semlyen, "Rational approximation of frequency domain responses by vector fitting" (IEEE Transactions on Power Delivery, 1999); Kundert and Zinke, "The Designer's Guide to Verilog-AMS" (Kluwer, 2004).
