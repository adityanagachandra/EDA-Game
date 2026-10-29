# Distributed RC lines and wire delay

An on-chip wire is not one resistor and one capacitor; its resistance and capacitance are spread along its length. Modeled as a distributed RC line, its voltage obeys a diffusion equation, its delay grows with the square of length, and its 50% delay differs from the lumped estimate. These facts drive repeater insertion and every wire-delay model a placer or router uses.

## A wire as a distributed RC line

A wire of length $L$ with resistance $r$ and capacitance $c$ per unit length has totals $R = rL$ and $C = cL$. Cut into many small segments it becomes an RC ladder, and in the limit the voltage obeys $\partial v/\partial t = \frac{1}{rc}\,\partial^2 v/\partial x^2$, the heat equation. Signals diffuse rather than propagate.

The Elmore delay at the far end sums, for each element of capacitance $c\,dx$ at position $x$, the resistance between it and the driver, $rx$: $T = \int_0^L r x\, c\, dx = rcL^2/2 = RC/2$. Example: 1000 um of wire at 0.5 ohm/um and 0.2 fF/um has $R = 500$ ohm and $C = 200$ fF, so $RC = 100$ ps and the Elmore delay is 50 ps, half the lumped product, because the near-end capacitance sees little resistance.

## The 0.38 RC and 0.69 RC rules

Solving the diffusion equation for an ideal step at the near end gives a far-end 50% delay of about $0.38RC$, below the Elmore value $0.5RC$ (Elmore is an upper bound on the 50% delay of an RC tree's step response). A single lumped RC, with all resistance in front of all capacitance, gives $0.69RC$.

Example: the 500 ohm, 200 fF wire has $RC = 100$ ps; its 50% delay is about 38 ps, whereas lumping it would claim 69 ps. Lumped models therefore overestimate wire delay by nearly a factor of two. The 10-90 rise time of the distributed line is about $0.9RC$, against $2.2RC$ for the lumped stage. Practical delay calculators drive the wire with a driver resistance too, and then the total is close to $0.69R_dC + 0.38RC$ plus a cross term, which is why wire models are fitted rather than guessed.

## Quadratic growth with length

Because both $R$ and $C$ grow linearly with length, the wire's intrinsic delay grows as $L^2$: doubling a wire quadruples its delay. Splitting the wire into $k$ equal segments with ideal repeaters gives each segment $RC/k^2$, so the total wire delay is $k \cdot 0.38RC/k^2 = 0.38RC/k$, linear in length once the segment length is fixed.

Example: a 3000 um wire with $r = 0.4$ ohm/um and $c = 0.2$ fF/um has $RC = 1200 \cdot 600$ ohm fF $= 720$ ps and an unbuffered delay of about 274 ps. Three ideal repeaters' worth of segments cut the wire part to 91 ps. Real repeaters add their own delay, so an optimizer balances repeater delay against wire delay; the optimum segment length depends on the driver's resistance and input capacitance and on the wire's r and c.

## Going deeper

Bakoglu, "Circuits, Interconnections, and Packaging for VLSI", chapter 5; Sakurai, "Closed-form expressions for interconnection delay, coupling, and crosstalk in VLSIs", IEEE Transactions on Electron Devices, 1993.
