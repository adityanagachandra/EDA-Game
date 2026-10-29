# Generalized scaling and the power wall

Generalized scaling separates the shrink factor of dimensions from that of voltage. It lets a device or architecture engineer ask what happens when voltage scales more slowly than geometry, the situation after Dennard scaling ended. The resulting factors for delay and power density quantify the power wall, and they are the back-of-the-envelope models behind DVFS policies and power-performance trade-offs.

## Generalized delay scaling

Shrink all dimensions, including oxide thickness, by $1/\kappa$, and the supply voltage by $1/s$, with $1 \le s \le \kappa$. Using the long-channel current $I \propto (W/L)C_{ox}V^2$:
$W/L$ is unchanged, $C_{ox}$ rises by $\kappa$, so current scales by $\kappa/s^2$;
gate capacitance $\propto WL/t_{ox}$ scales by $1/\kappa$;
delay $CV/I$ scales by $(1/\kappa)(1/s)/(\kappa/s^2) = s/\kappa^2$.

Example: $\kappa = 1.4$ and $s = 1.1$. Delay scales by $1.1/1.96 = 0.561$, faster than Dennard's $1/\kappa = 0.714$.

With $s = \kappa$ this reduces to constant-field scaling. With $s < \kappa$ the field grows by $\kappa/s$, which gives more speed but stresses reliability. The long-channel current law is an assumption; with velocity saturation the current scales less favorably, so these factors are upper bounds on the speed benefit.

## Generalized power density

Power per device is $VI \propto (1/s)(\kappa/s^2) = \kappa/s^3$ when frequency follows the delay. Device density rises by $\kappa^2$. So power density scales by
$$\frac{\kappa^3}{s^3}.$$

Example: $\kappa = 1.4$, $s = 1.1$. Power density scales by $(1.4/1.1)^3 = 2.06$ per generation.

Only $s = \kappa$, full voltage scaling, keeps power density constant. Every bit of voltage that fails to scale costs power density cubically under these assumptions. After the threshold floor stopped voltage scaling ($s$ near 1), raising frequency with the delay would have multiplied power density by nearly $\kappa^3$ per generation, far beyond what packaging can remove. This is the quantitative content of the power wall, and why frequency stopped following device speed.

## Power density at fixed frequency

Industry responded by freezing clock frequency and spending density on more cores and wider datapaths. With frequency fixed, dynamic power per device is $CV^2f \propto (1/\kappa)(1/s^2)$, and with density $\kappa^2$:
$$\text{power density} \propto \frac{\kappa}{s^2}.$$

Example: $\kappa = 1.4$, $s = 1.1$: $1.4/1.21 = 1.16$ per generation. With no voltage scaling at all ($s = 1$), power density still grows by $\kappa$.

So even a frozen clock does not hold power flat without voltage scaling. Not every transistor can switch at full rate within the power budget, a situation sometimes described as dark silicon. The engineering responses are what EDA tools now optimize: clock gating and power gating to idle unused logic, multiple voltage domains and DVFS, near-threshold operation for efficiency, multi-threshold libraries, and power-aware synthesis and placement. Each one is a way of making $\alpha C V^2 f$ smaller where it matters.

## Going deeper

Baccarani, Wordeman and Dennard, "Generalized scaling theory and its application to a 1/4 micrometer MOSFET design", IEEE Transactions on Electron Devices (1984); Taur and Ning, "Fundamentals of Modern VLSI Devices" (2nd ed.), section 4.1.
