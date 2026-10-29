# Delay of a distributed RC wire

A wire is not a lumped resistor and capacitor: its resistance and capacitance are spread along its length. Modeling it as a distributed line changes the Elmore delay by a factor of two and makes delay grow with the square of length. That quadratic growth is the reason long wires need repeaters, and the reason a router's detour can wreck a path that met timing after placement.

## The distributed line halves the lumped delay

Let a wire of length $L$ have resistance $r$ per unit length and capacitance $c$ per unit length. Lumped, it would be one resistor $rL$ charging one capacitor $cL$, with time constant $rcL^2$. Split it into many small segments instead: each bit of capacitance at position $x$ sees only the resistance $rx$ between it and the driver. The Elmore sum becomes $\int_0^L r x \, c \, dx = rcL^2/2$, half the lumped value. A 350 um wire with $r = 1$ ohm/um and $c = 0.15$ fF/um gives $1 \times 0.15 \times 350^2 / 2 = 9187.5$ ohm-fF, about 9.19 ps. The same result follows from a pi model with $cL/2$ at each end.

## Driver, wire and load together

A driver of resistance $R_d$ drives the wire into a load capacitance $C_L$. The driver resistance is shared by every capacitor downstream, so it charges the whole wire plus the load, $R_d (cL + C_L)$. The wire's own resistance charges half its own capacitance, as above, plus all of the load: $rL(cL/2 + C_L)$. The total Elmore delay is $T = R_d (cL + C_L) + rL\,(cL/2 + C_L)$. For $R_d = 800$ ohm, $L = 350$ um, $r = 1$ ohm/um, $c = 0.15$ fF/um and $C_L = 20$ fF, the driver term is 58000 ohm-fF and the wire term 16187.5 ohm-fF, so the delay is about 74.19 ps. One ohm times one femtofarad is one femtosecond.

## Quadratic growth and repeaters

Only the $rcL^2/2$ term is quadratic, but for long wires it dominates. Doubling a 400 um wire to 800 um, with $R_d = 400$ ohm, $r = 2$ ohm/um, $c = 0.2$ fF/um and $C_L = 17$ fF, takes the delay from 84.4 ps to 226 ps, much more than double. Cutting the wire into $k$ segments with repeaters makes each segment's quadratic term $rc(L/k)^2/2$, so the wire part of the total falls to $rcL^2/(2k)$ at the cost of $k$ repeater delays; at the optimum the total grows linearly with length. This is why routing detours matter: a net that was estimated as a short Steiner tree at placement and routed with a long detour can gain far more delay than its extra length suggests.

## Going deeper

Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed., 2003), chapter on wires; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on timing-driven routing and buffering.
