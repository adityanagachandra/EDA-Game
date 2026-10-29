# Electromigration and IR-drop signoff

Extracted resistance feeds two reliability and power checks. Electromigration (EM) signoff checks that current densities in wires and vias are low enough for the chip to last its rated lifetime. IR-drop signoff checks that the power grid delivers voltage within budget to every cell. Both are computed from the extracted power and signal networks, and both can be broken by an innocent-looking timing ECO.

## Current density

Electromigration is the slow drift of metal atoms pushed by the electron flow, which eventually opens voids or forms shorts. Its rate depends on current density $J = I/(WT)$. A power strap 0.1 um wide and 0.2 um thick carrying 5.6 mA has $J = 5.6/0.02 = 280$ mA/um$^2$; against a limit of 200 mA/um$^2$ it fails. EM rules come from the foundry per layer, per via type and per temperature, with separate limits for DC (unidirectional, as in power nets), average, RMS and peak current on signal nets. EM signoff computes currents from the extracted network and activity, divides by the extracted cross-section, and flags violations. The fix is wider wires, more vias or less current per wire.

## IR drop along a strap

Current flowing through grid resistance drops voltage: $\Delta V = IR$. A strap with sheet resistance 0.02 ohm/sq, 2 um wide and 750 um long has 375 squares, 7.5 ohm, and delivering 5 mA through it costs 37.5 mV, 4.2% of a 0.9 V supply. Real grids are meshes of many layers, solved as large sparse linear systems (nodal analysis with current sources at cell taps), statically with average currents and dynamically with time-varying switching. Signoff budgets static plus dynamic drop at a few percent of the supply. The drop matters for timing because cells slow down at lower voltage, so signoff timers take per-instance voltages from IR analysis.

## Lifetime scaling

Black's equation models EM lifetime as $\text{MTTF} \propto J^{-n} e^{E_a/kT}$, with a current exponent $n$ typically between 1 and 2 and an activation energy $E_a$ specific to the metal and its interfaces. At a fixed temperature, raising $J$ by a factor of 2 with $n = 2$ cuts lifetime to $2^{-2} = 0.25$ of its old value. Temperature enters exponentially, so hot spots age fastest. The design consequence: a timing ECO that upsizes a driver raises the current in its output wire and the vias under it, possibly past the EM limit, even though timing improved. Every timing ECO must be followed by EM and IR re-checks, which is one of the dependencies an automated ECO loop must respect.

## Going deeper

Black, "Electromigration: a brief survey and some recent results", IEEE Transactions on Electron Devices (1969); Lienig and Thiele, "Fundamentals of Electromigration-Aware Integrated Circuit Design" (2018).
