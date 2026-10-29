# Slew-limited buffer spacing

How far can one clock buffer drive a wire before the edge degrades too much? The answer sets the spacing of buffers along long clock routes and the size of the clusters a CTS engine builds. It comes from combining the Elmore delay of a driver plus wire plus load with a slew estimate, then solving for the length that just meets the transition limit. The quadratic growth of wire delay with length is what makes buffering necessary at all.

## Elmore delay of a driven wire

A buffer with output resistance $R_b$ drives a wire of length $L$ (resistance $r$, capacitance $c$ per unit length) into a load $C_{in}$, the next buffer's input. Its Elmore delay to the far end is
$$T = R_b(cL + C_{in}) + rL\left(\frac{cL}{2} + C_{in}\right).$$
The driver charges all capacitance; the wire's own resistance charges half of its capacitance and all of the load. For $R_b = 400\ \Omega$, $r = 1.3\ \Omega$/µm, $c = 0.24$ fF/µm, $C_{in} = 7$ fF and $L = 450$ µm, $T = 400 \cdot 115 + 585 \cdot 61 = 81\,685\ \Omega\cdot$fF, or 81.69 ps.

## From delay to slew

A common first-order estimate treats the response as a single pole with time constant equal to the Elmore delay, so the 10-90% slew is about $\ln 9 \cdot T \approx 2.197\,T$. For a wire with $T = 43.2$ ps the slew estimate is 95 ps. Real waveforms at the end of a resistive wire are not single-pole, and signoff timers use better models (effective capacitance, waveform propagation), but the estimate is accurate enough for deciding where buffers go. Its advantage is that it is a closed form in $L$.

## Solving for the longest wire

Setting $\ln 9 \cdot T = S$ for a limit $S$ gives a quadratic in $L$:
$$\frac{rc}{2}L^2 + (R_b c + r C_{in})L + R_b C_{in} - \frac{S}{\ln 9} = 0,$$
with $S/\ln 9$ in the same units as the other terms. For $R_b = 350\ \Omega$, $r = 1.5$, $c = 0.18$, $C_{in} = 2.5$ fF and $S = 60$ ps, the coefficients are 0.135, 66.75 and $875 - 27\,307$, and the positive root is about 260 µm. Beyond that the buffer must be repeated. The $L^2$ term is why doubling a wire more than doubles its delay, and why long routes are cut into buffered stages.

## Going deeper

Bakoglu, "Circuits, Interconnections, and Packaging for VLSI" (1990), the chapter on interconnect delay and repeaters; Kahng et al., "VLSI Physical Design" (2nd ed., 2022).
