# Depletion width and junction capacitance

The region around a junction that is swept free of mobile carriers is the depletion region. Its width sets the junction capacitance that loads every source and drain node, the field that drives breakdown and tunneling leakage, and the depletion charge that enters the MOSFET threshold. Parasitic extraction and compact junction models both reduce to the same few formulas, derived by integrating Poisson's equation across this region.

## Width of the depletion region

In the depletion approximation the region from $-x_p$ to $x_n$ holds only ionized dopants, and charge neutrality requires $N_A x_p = N_D x_n$. Integrating Poisson's equation twice, with the total potential drop equal to $V_{bi} + V_R$ under reverse bias $V_R$, gives
$$W = x_n + x_p = \sqrt{\frac{2\varepsilon_{Si}(V_{bi} + V_R)}{q}\left(\frac{1}{N_A} + \frac{1}{N_D}\right)}.$$

Example: $N_A = N_D = 10^{16}$ cm$^{-3}$, $V_{bi} = 0.71$ V, zero bias. With $\varepsilon_{Si} = 1.036 \times 10^{-12}$ F/cm, $W = \sqrt{2 \times 1.036 \times 10^{-12} \times 0.71 / 1.602 \times 10^{-19} \times 2 \times 10^{-16}} = 4.3 \times 10^{-5}$ cm $= 430$ nm.

The field is triangular and peaks at the metallurgical junction, at $E_{max} = 2(V_{bi} + V_R)/W$. Width grows as the square root of the total voltage and shrinks as the square root of the doping.

## Capacitance per unit area

A change in reverse bias adds or removes depletion charge at the edges, so the junction behaves as a parallel-plate capacitor of width $W$:
$$C_j = \frac{\varepsilon_{Si}}{W} \quad \text{per unit area}.$$

For a one-sided $n^+$-p junction, such as a source on a p-well, the heavily doped side barely depletes and $W \approx \sqrt{2\varepsilon_{Si}(V_{bi} + V_R)/(qN_A)}$.

Example: $N_A = 10^{18}$ cm$^{-3}$, $V_{bi} = 0.9$ V and $V_R = 0$. Then $W = 3.41 \times 10^{-6}$ cm and $C_j = 1.036 \times 10^{-12}/3.41 \times 10^{-6} = 3.04 \times 10^{-7}$ F/cm$^2$. Since 1 F/cm$^2$ equals $10^7$ fF/µm$^2$, $C_j = 3.04$ fF/µm$^2$.

Heavier well doping raises junction capacitance, which loads every diffusion node. Extraction tools combine an area term like this with a sidewall term per unit perimeter, each with its own parameters in the model card.

## Capacitance versus bias

Because $W \propto (V_{bi} + V_R)^{1/2}$ for an abrupt junction, capacitance falls with reverse bias as
$$C_j(V_R) = \frac{C_{j0}}{(1 + V_R/V_{bi})^{m}}, \qquad m = \tfrac12.$$
Graded junctions have $m$ closer to $1/3$. This is the form SPICE junction models use, with $C_{j0}$, $V_{bi}$ (often called PB) and $m$ (MJ) as parameters, for both the area and sidewall terms.

Example: $V_{bi} = 0.8$ V. The ratio $C_j(0.5\ \text{V})/C_j(3.5\ \text{V}) = \sqrt{4.3/1.3} = 1.82$.

The same relation underlies C-V doping extraction: $1/C_j^2$ is linear in $V_R$ with slope $2/(q\varepsilon_{Si}N)$, so the slope gives the doping and the intercept gives $V_{bi}$. A timing engine sees the consequence as a nonlinear load: a diffusion node's capacitance depends on its voltage during the transition, which characterization averages into effective values.

## Going deeper

Sze and Ng, "Physics of Semiconductor Devices" (3rd ed.), section 2.2; the ngspice manual, section on the junction diode model.
