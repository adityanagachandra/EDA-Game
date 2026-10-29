# Equation-based sizing

Equation-based sizing writes each performance as an explicit function of the design variables and solves the resulting system or optimization directly. In its simplest form it is the designer's hand procedure scripted: start from GBW and load, derive transconductances, then currents and widths. In its most formal form, posynomial models make the problem a geometric program. Its speed makes it the natural inner model for exploration and a starting point for simulation-based refinement.

## From GBW to the input pair

In a Miller two-stage op-amp, $\mathrm{GBW} = g_{m1}/(2\pi C_c)$. Choosing $C_c$ (often a fraction of the load capacitance, for stability and noise) fixes the input-pair transconductance: $g_{m1} = 2\pi \cdot \mathrm{GBW} \cdot C_c$. For GBW = 40 MHz and $C_c = 2$ pF, $g_{m1} = 2\pi \cdot 40 \cdot 2 = 502.7\ \mu$S, using the unit identity that MHz times pF gives microsiemens. The procedure works backwards from specs to device parameters, and each step uses an equation that is exact within its model. A sizing script can evaluate the whole chain in microseconds, so it can sweep $C_c$ or other free choices to explore the space before any simulation.

## From transconductance to current and size

With $g_{m1}$ known, the gm/Id methodology gives the bias. Choosing $g_m/I_D = 15\ \mathrm{V}^{-1}$ for the input pair (high efficiency and low noise, at the cost of linear range) gives $I_D = 502.7/15 = 33.5\ \mu$A per side and a tail current of 67 microamperes. The width follows from the lookup table's current density at that $g_m/I_D$ and the chosen $L$. Each device is sized the same way from its own role: mirror loads at low $g_m/I_D$ for matching, the second-stage driver from the stability condition below. Because every quantity comes from a model, the procedure is only as accurate as its device data, which is why lookup tables replaced square-law equations in modern equation-based flows.

## Stability as a sizing equation

Phase margin becomes an equation by placing the non-dominant pole: require $g_{m6}/(2\pi C_L) \ge k\cdot\mathrm{GBW}$ with $k$ around 2.2 for roughly 65 degrees. For GBW = 40 MHz and $C_L = 5$ pF, $g_{m6} \ge 2.2 \cdot 2\pi \cdot 40 \cdot 5 = 2765\ \mu$S. Note the asymmetry: the second stage needs several times the input pair's transconductance, so it dominates power. Equation-based optimizers collect such relations (gain, poles, slew, swing, noise) and either solve them in sequence, as above, or hand them to a solver. Geometric programming fits each performance as a posynomial in the variables, which makes the problem convex after a log change of variables, so the global optimum of the model is found reliably. The weakness is model error: short-channel devices are not posynomial, so the result must be verified, and usually refined, in SPICE.

## Going deeper

del Mar Hershenson, Boyd and Lee, "Optimal design of a CMOS op-amp via geometric programming" (IEEE TCAD, 2001); Jespers and Murmann, "Systematic Design of Analog CMOS Circuits Using Pre-Computed Lookup Tables" (2017).
