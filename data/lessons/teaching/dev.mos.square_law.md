# The MOSFET square law

The square law, $I_D = \tfrac12\mu C_{ox}(W/L)(V_{GS} - V_t)^2$, is the long-channel saturation current of a MOSFET. Modern short-channel devices deviate from it, but it remains the reference model: hand sizing, first-order delay estimates, the language of overdrive and transconductance, and the asymptote that every compact model must approach for long devices all come from it.

## Saturation current

When $V_{DS} \ge V_{GS} - V_t$ the inversion charge vanishes at the drain end: the channel is pinched off. Further drain voltage drops across a short depleted region near the drain, and the current stays at the value it reached at $V_{DS} = V_{GS} - V_t$:
$$I_{D,sat} = \frac{1}{2}k'\frac{W}{L}(V_{GS} - V_t)^2.$$

Example: $k' = 200$ µA/V$^2$, $W/L = 10$, $V_t = 0.4$ V, $V_{GS} = 1.0$ V. $I_{D,sat} = 0.5 \times 200 \times 10 \times 0.36 = 360$ µA.

The current depends on the square of overdrive, so a 10% overdrive change gives about 21% current change. That quadratic sensitivity is why threshold variation dominates delay variation at low supply voltage. In short-channel devices velocity saturation reduces the exponent toward 1, which the alpha-power law captures.

## Sizing for a target current

Inverting the square law gives the aspect ratio needed for a current:
$$\frac{W}{L} = \frac{2I_D}{k'(V_{GS} - V_t)^2}.$$

Example: 300 µA from a device with $k' = 250$ µA/V$^2$ at overdrive 0.5 V needs $W/L = 600/(250 \times 0.25) = 9.6$.

Standard-cell libraries implement this idea in discrete drive strengths: X1, X2, X4 cells scale width (or fin count) so that current scales proportionally. A gate sizing engine, such as the one in a timing-driven optimizer, picks among these sizes trading delay against input capacitance and leakage, both of which also scale with width. In FinFET processes $W$ is quantized by fins, so the continuous answer must be rounded up to a whole number of fins.

## Transconductance in saturation

The transconductance measures how strongly the gate controls the current:
$$g_m = \frac{\partial I_D}{\partial V_{GS}} = k'\frac{W}{L}(V_{GS} - V_t) = \frac{2I_D}{V_{GS} - V_t}.$$

Example: $I_D = 100$ µA at overdrive 0.2 V gives $g_m = 1000$ µA/V $= 1$ mS. The same current at overdrive 0.4 V gives only 0.5 mS.

The form $g_m = 2I_D/V_{ov}$ shows the analog trade-off: low overdrive gives more transconductance per unit current but needs a wider, slower device. The ratio $g_m/I_D = 2/V_{ov}$ is the figure of merit of the $g_m/I_D$ design method. It cannot exceed the subthreshold limit $1/(n\,kT/q)$, which is where the square law breaks down at low overdrive; a good compact model passes smoothly between the two.

## Going deeper

Sedra and Smith, "Microelectronic Circuits"; Jespers and Murmann, "Systematic Design of Analog CMOS Circuits" (2017), for the gm/ID method.
