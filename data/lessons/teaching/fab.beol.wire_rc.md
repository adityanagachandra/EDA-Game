# Wire RC per unit length and scaling

Every routed net is, electrically, a distributed resistance and capacitance. Resistance per unit length comes from resistivity and cross-section; capacitance per unit length from geometry and dielectric. Their product sets how wire delay grows with length, and its behaviour under scaling is why interconnect, not transistors, often limits performance at advanced nodes.

## Resistance per unit length

A rectangular wire of width $w$ and thickness $t$ with resistivity $\rho$ has resistance per length $r = \rho/(w t)$. With $\rho = 4$ micro-ohm cm (40 ohm nm), $w = 20$ nm and $t = 40$ nm: $r = 40 \times 1000/800 = 50$ ohm per um. For narrow copper, $\rho$ is an effective value well above the bulk 1.7 micro-ohm cm because electrons scatter at surfaces and grain boundaries, and the barrier steals cross-section. Extraction rule files therefore tabulate resistance against drawn width and sometimes spacing. Capacitance per length, typically 0.1 to 0.25 fF/um for local wires, comes from area and fringe coupling to the layers above and below and to neighbours on the same layer.

## Distributed RC delay

A uniform line of length $L$ with $r$ and $c$ per unit length has total resistance $rL$ and capacitance $cL$, but each piece of capacitance sees only the resistance between it and the driver. The Elmore delay of the distributed line is $r c L^2/2$. For $r = 50$ ohm/um, $c = 0.2$ fF/um and $L = 200$ um: $0.5 \times 50 \times 0.2 \times 40000 = 200{,}000$ ohm fF, which is 200 ps. The square law means doubling length quadruples delay, which is why long wires are broken by repeaters: with optimal buffering, delay grows linearly in length.

## What scaling does to wires

Shrink width and thickness by a factor $s < 1$. Resistance per length rises by $1/s^2$; capacitance per length stays roughly constant because it depends on ratios of dimensions. A route of fixed length therefore gets slower by $1/s^2$, and more when size effects raise resistivity. With $s = 0.7$ and resistivity up by 1.2, a 20 ps route becomes $20 \times 1.2/0.49 = 49$ ps. Local wires shrink in length along with the cells, which partly compensates, but global wires do not. This is the interconnect bottleneck behind buffering, layer assignment, low-k dielectrics and backside power.

## Going deeper

Bakoglu, "Circuits, Interconnections, and Packaging for VLSI" (1990); Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits", chapter 4.
