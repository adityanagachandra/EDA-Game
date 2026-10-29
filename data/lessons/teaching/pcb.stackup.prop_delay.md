# Propagation delay and effective dielectric constant

A signal on a board trace moves at the speed of light divided by the square root of the effective dielectric constant the field sees. That gives a delay per inch for each layer, the number that converts lengths into time for every length-matching constraint and every timing budget at board level. Microstrip and stripline in the same laminate have different delays, so equal lengths are not equal times.

## Delay in a homogeneous dielectric

Light travels about 11.8 inches per ns in vacuum. In a medium with relative permittivity $\varepsilon_{eff}$, the speed is $c/\sqrt{\varepsilon_{eff}}$, so the delay per inch is $t_{pd} = \sqrt{\varepsilon_{eff}}/11.8$ ns. A stripline is fully embedded, so $\varepsilon_{eff} = \varepsilon_r$. Example: $\varepsilon_r = 4.0$ gives $2/11.8 = 0.1695$ ns, about 170 ps per inch, the familiar rule of roughly 6 inches per ns on FR-4. Delay grows only as the square root of the dielectric constant: going from 4.0 to 3.4 (a low-loss laminate) shortens delay by about 8%. The value to use is the dielectric constant at the signal's frequency, since laminate Dk falls slowly with frequency.

## Effective dielectric constant of microstrip

A microstrip's field is partly in the laminate and partly in the air above, so it sees an **effective** dielectric constant between 1 and $\varepsilon_r$. A common approximation for typical FR-4 geometries is $\varepsilon_{eff} \approx 0.475\,\varepsilon_r + 0.67$. Example: $\varepsilon_r = 4.2$ gives $\varepsilon_{eff} = 2.665$, so $t_{pd} = \sqrt{2.665}/11.8 \approx 138$ ps per inch, against 174 ps per inch for stripline in the same laminate. More exact forms depend on $w/h$; field solvers compute $\varepsilon_{eff}$ directly as the ratio of the line's capacitance with and without the dielectric. Solder mask over the trace raises $\varepsilon_{eff}$ slightly, and so does any coating applied after assembly.

## Delay of a mixed route

A routed net usually changes layers. Its delay is the sum over segments of length times that layer's delay, plus the vias. Example, in a laminate with $\varepsilon_r = 4.0$: 1.0 inch of microstrip at about 136 ps per inch and 3.0 inches of stripline at 169.5 ps per inch give $135.9 + 508.5 \approx 644$ ps. Two nets of equal total length but different layer mixes can differ by tens of ps, which exceeds many matching budgets. That is why board tools match **delay**, not raw length, once nets change layers: the constraint system stores a per-layer delay table derived from the stackup and converts each segment. A via adds its own delay, often modeled as an equivalent length, which tools let the designer specify.

## Going deeper

Johnson and Graham, "High-Speed Digital Design" (1993), chapter 4. Hall and Heck, "Advanced Signal Integrity for High-Speed Digital Designs" (2009), chapter 2.
