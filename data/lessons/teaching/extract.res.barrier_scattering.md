# Barrier layers and size effects in copper

A modern copper wire is not all copper. A thin barrier liner, needed to stop copper diffusing into the dielectric, coats its sides and bottom and conducts poorly. And as wires shrink toward the electron mean free path, surface and grain-boundary scattering raise copper's resistivity. Both effects grow as wires narrow, so an extractor's resistance model at advanced nodes must use the conducting cross-section and a width-dependent resistivity rather than the drawn rectangle and a constant.

## The conducting cross-section

In a damascene copper process the barrier, often tantalum-based, lines both sidewalls and the bottom of the trench. If the barrier is $b$ thick and conducts negligibly, the copper occupies $(W - 2b)(T - b)$ of the drawn $W \times T$. A wire 0.05 um wide and 0.04 um thick with a 0.004 um barrier keeps $0.042 \times 0.036 = 0.001512$ um$^2$ of copper out of 0.002 um$^2$: only 75.6%. The barrier does not shrink at the same rate as the wire, because it must stay thick enough to block diffusion, so its share of the cross-section grows at each node. This is one of the main reasons resistance per unit length has risen much faster than simple scaling predicts.

## Resistance penalty of the liner

Resistance goes as one over the conducting area, so the liner raises it by the ratio of drawn to copper area. For a wire 0.03 um wide and 0.08 um thick with a 0.003 um barrier, the copper area is $0.024 \times 0.077 = 0.001848$ um$^2$ against 0.0024 drawn, a ratio of 1.2987: the barrier alone adds 29.9% to the resistance. Narrow wires suffer most, because the sidewall term $2b$ is a larger fraction of a small width. An extractor that uses drawn dimensions with bulk resistivity would under-report resistance by this factor. In practice the technology file either describes the barrier explicitly or folds it into an effective resistivity per width, and the extractor must apply whichever convention the file uses, never both.

## Width-dependent resistivity

When the copper width approaches the electron mean free path (about 40 nm in bulk copper at room temperature), electrons scatter off surfaces and grain boundaries more often, and resistivity rises. A simple model writes $\rho = \rho_0 (1 + K/W_c)$ with $W_c$ the copper width. With $\rho_0 = 0.018$ ohm-um and $K = 0.02$ um, a wire drawn 0.04 um wide with a 0.003 um barrier has $W_c = 0.034$ um and $\rho = 0.0286$ ohm-um. For 40 um of length and 0.057 um of copper height, $R = 0.0286 \times 40 / 0.001938 = 590$ ohm. Real models (Fuchs-Sondheimer for surfaces, Mayadas-Shatzkes for grains) are more detailed, and foundries tabulate the result per layer and width. The point for tool design is that resistivity is a function of geometry the extractor must look up, not a constant.

## Going deeper

Gall, "Electron mean free path in elemental metals", Journal of Applied Physics (2016); the IRDS interconnect chapter on resistivity scaling and alternative metals.
