# Compact-model parameter extraction

Parameter extraction fits a compact model's parameters to measured or TCAD-simulated device data. It starts with simple, physically motivated extractions of threshold and transconductance factor, then fits short-channel, mobility, parasitic and leakage terms region by region, and ends with a global least-squares polish. The procedure is well defined and checkable, which is why it is a natural task for an automated or agent-driven flow.

## Threshold by linear extrapolation

In the linear region at small $V_{DS}$, $I_D \approx k'(W/L)(V_{GS} - V_t - V_{DS}/2)V_{DS}$ is linear in $V_{GS}$. Real curves bend because of mobility degradation and series resistance, so the standard method takes the tangent at the point of maximum transconductance and extrapolates it to zero current:
$$V_t = V_{GS,0} - \frac{I_{D,0}}{g_{m,max}} - \frac{V_{DS}}{2},$$
where $V_{GS,0}$ and $I_{D,0}$ are the bias and current at maximum $g_m$.

Example: $V_{DS} = 50$ mV, maximum $g_m = 250$ µA/V at $V_{GS} = 0.7$ V where $I_D = 40$ µA. $V_t = 0.7 - 0.16 - 0.025 = 0.515$ V.

Other definitions are common: constant current (the $V_{GS}$ at which $I_D = I_0 \cdot W/L$, for a chosen $I_0$), and second-derivative peak. They give different numbers for the same device, so an extraction report must state its method, and comparing thresholds across sources requires the same definition.

## Threshold from saturation data

In saturation the long-channel law $I_D = \tfrac12 K(V_{GS} - V_t)^2$, with $K = k'W/L$, makes $\sqrt{I_D}$ linear in $V_{GS}$. Two points determine the line:
$$V_t = \frac{V_2\sqrt{I_1} - V_1\sqrt{I_2}}{\sqrt{I_1} - \sqrt{I_2}}.$$

Example: $I_1 = 32$ µA at 0.8 V and $I_2 = 128$ µA at 1.2 V. $\sqrt{I_1} = 5.657$ and $\sqrt{I_2} = 11.314$, so $V_t = (1.2 \times 5.657 - 0.8 \times 11.314)/(5.657 - 11.314) = 0.4$ V.

In practice an extraction fits a line through many points rather than two, which is a linear least-squares problem, and checks the residuals. Curvature in $\sqrt{I_D}$ reveals velocity saturation and mobility degradation, which the square law ignores. That is why staged extraction uses long-channel devices for these first parameters: their behavior is closest to the simple law, so the fitted values are physical.

## The transconductance factor and staged fitting

The slope of $\sqrt{I_D}$ against $V_{GS}$ is $\sqrt{K/2}$, so
$$K = 2\left(\frac{\sqrt{I_2} - \sqrt{I_1}}{V_2 - V_1}\right)^2.$$
Example: the same points give slope $5.657/0.4 = 14.14$ and $K = 400$ µA/V$^2$.

Production extraction is staged. Long-wide devices fix threshold, body effect and low-field mobility. Short devices add roll-off, DIBL, velocity saturation and output conductance. Narrow devices add width effects. Then come parasitic resistance, capacitances from C-V data, leakage components and temperature coefficients. Each stage fits the parameters that dominate in a chosen bias and geometry region, holding earlier ones fixed, which keeps each fit well conditioned and the values physical.

A final global optimization, usually weighted least squares on log current for subthreshold data and relative error elsewhere, polishes everything together. Parameter bounds keep the optimizer from trading physical meaning for fit quality.

## Going deeper

Schroder, "Semiconductor Material and Device Characterization" (3rd ed.), chapter 4; Ortiz-Conde and others, "A review of recent MOSFET threshold voltage extraction methods", Microelectronics Reliability (2002).
