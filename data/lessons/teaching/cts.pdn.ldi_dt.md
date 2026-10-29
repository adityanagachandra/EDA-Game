# Inductive supply noise and L di/dt

Off-chip, the supply path runs through bumps, package planes and board, all of which have inductance. A fast change in chip current produces a voltage $L\,di/dt$ across that inductance, which subtracts from the on-die supply. Together with on-die decap the package forms a resonant circuit, and the worst droop often occurs at its resonance. This is why power-integrity analysis spans die, package and board.

## Inductive noise

The voltage across an inductance is $V = L\,di/dt$. For a package path of 170 pH and a current step of 500 mA in 350 ps, $V = 170 \cdot 500/350 = 242.9$ mV, a large fraction of a sub-1 V supply. Units: picohenries times milliamps per picosecond give millivolts. The noise depends on how fast the current changes, not on its average, so events such as clock gating that wakes many registers at once, or a power domain turning on, are the dangerous ones.

## Parallel bumps

Many power bumps in parallel reduce the effective inductance. Treating $N$ bumps as independent, each of inductance $L_b$, the total is $L_b/N$. To keep $L_b\Delta I/(N\Delta t)$ under a budget $V_b$, $N \ge L_b\Delta I/(\Delta t V_b)$: for $L_b = 120$ pH, a 3400 mA rise in 600 ps and a 25 mV budget, $N \ge 27.2$, so 28 bumps, and as many again for ground. Mutual inductance between neighbouring bumps makes real arrays less effective than this estimate, and interleaving VDD and VSS bumps reduces the loop inductance.

## Package and die resonance

The package inductance $L$ and the on-die decap $C$ form an LC tank with resonance $f = 1/(2\pi\sqrt{LC})$. For 260 pH and 290 nF, $f = 1/(2\pi\sqrt{7.54 \times 10^{-17}}) = 18.3$ MHz. Near this frequency the PDN impedance seen by the die peaks, and periodic activity at or near it, for instance a loop of bursts, builds up the largest droop. Designers specify a target impedance, roughly the allowed ripple divided by the transient current, and add on-die decap, package capacitors and damping until the impedance stays below it across frequency.

## Going deeper

Swaminathan and Engin, "Power Integrity Modeling and Design for Semiconductors and Systems" (2007); Popovich, Mezhiba and Friedman, "Power Distribution Networks with On-Chip Decoupling Capacitors" (2008).
