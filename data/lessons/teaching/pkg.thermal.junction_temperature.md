# Junction temperature and theta values

Datasheets summarise package thermal behaviour with theta values: thermal resistances between the junction and a reference such as the ambient air, the case or the board, measured under standard conditions. Junction temperature then follows from power and theta, and so does the maximum power a part can dissipate. Knowing what each theta does and does not mean prevents the most common thermal design errors.

## Junction from power and theta

The junction temperature is $T_j = T_a + P\,\theta_{ja}$, where $\theta_{ja}$ is the junction-to-ambient resistance in K/W (often written °C/W). A 15 W part with $\theta_{ja} = 3.2$ K/W in 45 °C air runs at $45 + 48 = 93$ °C. The junction temperature is what matters for timing, leakage and reliability: leakage grows steeply with temperature, so a hot die burns more power and runs hotter still. $\theta_{ja}$ is measured on a standard test board in still or moving air, so it is a comparison figure. A real system with a different board, airflow or neighbours gives a different value.

## Maximum allowed power

Inverting gives the power budget: $P_{max} = (T_{j,max} - T_a)/\theta_{ja}$. With $T_{j,max} = 105$ °C, $T_a = 50$ °C and $\theta_{ja} = 2.5$ K/W, $P_{max} = 55/2.5 = 22$ W. This number flows back into the architecture: it bounds the frequency and voltage the chip can run at, and dynamic thermal management throttles the chip when sensors report the junction approaching the limit. Leakage feedback makes the true limit lower than the formula when leakage is a large share of power, because the power itself depends on temperature. A designer who uses typical ambient temperature instead of the worst case gets a budget that fails in the field.

## Case and sink resistances

With a heat sink, the path splits into junction-to-case $\theta_{jc}$, case-to-sink $\theta_{cs}$ through the interface material, and sink-to-ambient $\theta_{sa}$, in series: $T_j = T_a + P(\theta_{jc} + \theta_{cs} + \theta_{sa})$. At 80 W with 0.15, 0.1 and 0.4 K/W and 35 °C air, $T_j = 35 + 80 \times 0.65 = 87$ °C. $\theta_{jc}$ is measured with essentially all heat forced through the case top, so it is valid only when that path dominates. Datasheets also give characterisation parameters such as $\Psi_{jt}$, junction-to-top: these are not resistances of a path but ratios that let an engineer infer $T_j$ from a measured case-top temperature in a real system.

## Going deeper

JEDEC JESD51 series on thermal test methods; Lasance and Poppe, "Thermal Management for LED Applications" (2014), on theta and psi.
