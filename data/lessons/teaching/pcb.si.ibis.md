# IBIS models

IBIS (I/O Buffer Information Specification) is the standard behavioral model of a chip's I/O buffers for board-level signal integrity. Instead of a transistor netlist, the vendor publishes tables: current versus voltage for the pull-up, pull-down and clamp structures, voltage versus time for switching, die capacitance and package parasitics, each at typical, minimum and maximum corners. Board simulators combine IBIS buffers with transmission-line models of the routed nets. The format is a clean seam between the silicon vendor's IP and the board designer's analysis.

## What an IBIS model contains

An IBIS file is plain text organized by component and pin. For each buffer model it gives **V-I tables**: [Pullup] and [Pulldown] for the output stage, and [POWER Clamp] and [GND Clamp] for the protection diodes that conduct when the pin is driven above the supply or below ground. **V-t tables** ([Rising Waveform], [Falling Waveform]) give the output's switching behavior into specified test loads, and [Ramp] gives edge rates. **C_comp** is the die capacitance at the pad. The [Package] section gives lumped R_pkg, L_pkg and C_pkg, and per-pin values or package models can refine it. Every table has **typ, min and max** columns for process, voltage and temperature corners. Because it describes only pin behavior, it hides the transistor design and simulates far faster than SPICE.

## Reading drive strength from a V-I table

Near 0 V the [Pulldown] curve is nearly linear, and its slope gives the output resistance the line sees: $R = \Delta V / \Delta I$. Example: 0.2 V at 8 mA and 0.4 V at 16 mA give $0.2\ \text{V}/8\ \text{mA} = 25$ ohms. Farther out the transistor saturates, and the ratio of voltage to current at a saturated point (1.5 V at 30 mA, say, 50 ohms) is a large-signal number, not the small-signal resistance that sets the source reflection coefficient $\Gamma_S = (R - Z_0)/(R + Z_0)$. With $R = 25$ ohms on a 50 ohm line, $\Gamma_S = -1/3$, and a series terminator of about 25 ohms would match the source. The simulator uses the full nonlinear curve; the slope is what a designer uses to choose termination.

## Corners, power and the limits of IBIS

A simulation at the typ corner alone is incomplete: the fast-strong corner (max column) gives the worst overshoot, ringing and crosstalk, and the slow-weak corner (min) the worst timing, so SI runs bracket both. Classic IBIS assumes an ideal supply; when many outputs switch at once, the local rail sags and the buffers slow, so power-aware features tie output behavior to supply current for simultaneous switching noise analysis. Lumped package R, L and C are too coarse when package traces are electrically long or pins couple, and are replaced by per-pin or S-parameter package models. Finally, tables cannot express adaptive equalization in multi-gigabit SerDes, which is why IBIS-AMI adds executable algorithmic models.

## Going deeper

The IBIS specification, maintained by the IBIS Open Forum. Hall and Heck, "Advanced Signal Integrity for High-Speed Digital Designs" (2009), chapter on I/O buffer modeling.
