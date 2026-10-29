# Linear-region drain current

Below pinch-off, a MOSFET with an inverted channel behaves like a voltage-controlled resistor. This linear (triode) region governs pass gates, transmission gates, the final part of every output transition and analog switches. Its current equation is the first thing a compact model must reproduce, and the region check that decides between it and saturation is the first branch every hand analysis takes.

## The on-resistance of a switch

With an inverted channel and a tiny drain voltage, the channel charge is uniform, $Q_{inv} = C_{ox}(V_{GS} - V_t)$ per unit area. The current is charge times drift velocity, $I_D = W Q_{inv}\mu V_{DS}/L$, so
$$I_D \approx k'\frac{W}{L}(V_{GS} - V_t)V_{DS}, \qquad k' = \mu C_{ox},$$
and the device is a resistor with
$$R_{on} = \frac{1}{k'(W/L)(V_{GS} - V_t)}.$$

Example: $k' = 300$ µA/V$^2$, $W/L = 4$, $V_{GS} - V_t = 0.6$ V. $R_{on} = 1/(300 \times 10^{-6} \times 4 \times 0.6) = 1.39$ kilohms.

This is the resistance that switch-level timing models and RC delay estimates assign to a turned-on transistor. It falls with overdrive, so pass gates passing a voltage near $V_{DD}$, where $V_{GS}$ shrinks, become slow. That is why transmission gates pair an NMOS with a PMOS.

## The full linear-region current

As $V_{DS}$ grows, the channel charge falls from source to drain, since the local gate-to-channel voltage is $V_{GS} - V(y)$. Integrating the current along the channel gives the long-channel linear-region equation:
$$I_D = k'\frac{W}{L}\left[(V_{GS} - V_t)V_{DS} - \frac{V_{DS}^2}{2}\right], \qquad V_{DS} < V_{GS} - V_t.$$

Example: $k' = 200$ µA/V$^2$, $W/L = 5$, $V_t = 0.35$ V, $V_{GS} = 1.0$ V, $V_{DS} = 0.3$ V. The bracket is $0.65 \times 0.3 - 0.045 = 0.15$, so $I_D = 200 \times 5 \times 0.15 = 150$ µA.

The quadratic term is the price of the falling channel charge. It makes the output characteristic bend over and flatten exactly where $V_{DS} = V_{GS} - V_t$, the drain saturation voltage. This simple equation ignores the variation of depletion charge along the channel (bulk-charge effect), which real compact models include.

## Choosing the region

Before computing a current, decide the region by comparing $V_{DS}$ with the overdrive $V_{ov} = V_{GS} - V_t$.
If $V_{GS} < V_t$, the device is in subthreshold (ideally off).
If $V_{DS} < V_{ov}$, it is linear and the triode equation applies.
If $V_{DS} \ge V_{ov}$, the channel is pinched off at the drain, and the current stays at its value at $V_{DS} = V_{ov}$: $I_D = \tfrac12 k'(W/L)V_{ov}^2$.

Example: $k' = 250$ µA/V$^2$, $W/L = 4$, $V_{ov} = 0.6$ V. At $V_{DS} = 0.4$ V the device is linear, $I_D = 1000 \times (0.24 - 0.08) = 160$ µA. At $V_{DS} = 0.9$ V it is saturated, $I_D = 0.5 \times 1000 \times 0.36 = 180$ µA.

The two formulas meet with equal value and slope at $V_{DS} = V_{ov}$, which is the minimum smoothness a simulator needs. Real compact models go further and blend the regions with smoothing functions, so that higher derivatives are continuous too.

## Going deeper

Sedra and Smith, "Microelectronic Circuits", chapter on MOSFETs; Tsividis and McAndrew, "Operation and Modeling of the MOS Transistor" (3rd ed.), chapter 4.
