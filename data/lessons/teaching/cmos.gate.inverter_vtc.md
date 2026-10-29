# The CMOS inverter transfer characteristic

The voltage transfer characteristic (VTC) of the CMOS inverter, output voltage against input voltage, is the reference curve for every static gate. Its switching threshold sets where a gate decides between 0 and 1, its slope sets noise immunity, and its shape explains why CMOS restores levels. Characterization flows and noise-analysis tools measure points on this curve, and sizing decisions start from it.

## Five regions of the curve

Sweep the input from 0 to $V_{DD}$. Below $V_{tn}$ the nMOS is off and the pMOS pulls the output to $V_{DD}$. Just above $V_{tn}$ the nMOS conducts in saturation while the pMOS stays in its linear region, and the output begins to fall. Near the middle both devices are saturated; in an ideal square-law device this region is almost vertical, so a small input change produces a large output change: the gain is high. Further on, the pMOS saturates while the nMOS enters its linear region, and above $V_{DD} - |V_{tp}|$ the pMOS is off and the output sits at 0. The output stays at a rail over wide input ranges and moves through a steep transition in between. That steep, high-gain middle is what makes a chain of inverters restore a degraded logic level back to a clean rail.

## Computing the switching threshold

The switching threshold $V_M$ is the input voltage at which $V_{out} = V_{in}$. There both devices are saturated, and with the long-channel square law their currents are equal:
$$\frac{k_n}{2}(V_M - V_{tn})^2 = \frac{k_p}{2}(V_{DD} - V_M - |V_{tp}|)^2.$$
Taking square roots with $r = \sqrt{k_p / k_n}$ gives
$$V_M = \frac{V_{tn} + r(V_{DD} - |V_{tp}|)}{1 + r}.$$
Each $k$ is proportional to mobility times $W/L$, so $k_p/k_n = (\mu_p W_p)/(\mu_n W_n)$. Example: $V_{DD} = 1.2$ V, $V_{tn} = |V_{tp}| = 0.4$ V and $k_p = k_n$ give $r = 1$ and $V_M = (0.4 + 0.8)/2 = 0.6$ V, exactly midway. With $W_p = W_n$ and $\mu_n = 2\mu_p$, $r = \sqrt{0.5} = 0.707$ and $V_M = (0.4 + 0.566)/1.707 = 0.566$ V: the stronger nMOS pulls the threshold down.

## Sizing to place the threshold

Inverting the formula gives the strength ratio needed for a target $V_M$: $r = (V_M - V_{tn})/(V_{DD} - V_M - |V_{tp}|)$, then $W_p/W_n = r^2 \mu_n/\mu_p$. For a symmetric threshold with equal thresholds, $r = 1$ and $W_p/W_n = \mu_n/\mu_p$, about 2 to 3 in older planar processes. Two cautions matter for tools. First, short-channel devices are closer to velocity-saturated than square-law, so real current is nearer linear in overdrive and the required ratio is smaller; characterization measures the true curve rather than trusting the formula. Second, $V_M$ is weakly sensitive to sizing: doubling $W_p$ moves it only modestly. That is why libraries size for equal rise and fall delays rather than for an exact $V_M$, and accept wherever the threshold lands.

## Going deeper

Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), chapter 5; Weste and Harris, "CMOS VLSI Design" (4th ed.), section 2.5.
