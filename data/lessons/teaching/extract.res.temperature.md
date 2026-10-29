# Temperature dependence of wire resistance

Metal resistance rises with temperature, roughly linearly over the operating range. An extractor reports resistance at a stated temperature, and signoff corners run from cold to hot, so the same geometry yields different resistances per corner. The temperature coefficient is the knob, and a tool designer has to know which numbers scale with it and which, like the driver resistance characterized in a cell library, are handled elsewhere.

## The linear temperature coefficient

Over the range chips operate in, metal resistance follows $R(T) = R_0[1 + \alpha (T - T_0)]$, where $R_0$ is the value at a reference temperature $T_0$ (often 25 C) and $\alpha$ is the temperature coefficient of resistance, on the order of 3 to $4 \times 10^{-3}$ per C for copper interconnect. A wire of 160 ohm at 25 C with $\alpha = 0.003$ per C measures $160 \times (1 + 0.003 \times 85) = 200.8$ ohm at 110 C, a 25% increase. Some technology files also carry a quadratic term; the linear form is the one most often used. Vias and barrier materials have their own coefficients, so an extractor applies them per layer and per via type.

## Moving between corners

Extraction results are often stored at one temperature and rescaled for another. To move from a hot value to a cold one, first recover the reference value, then scale again. A net reported at 620.4 ohm at 125 C with $\alpha = 0.0032$ per C has $R_{25} = 620.4/(1 + 0.0032 \times 100) = 470$ ohm, and at -40 C it becomes $470 \times (1 - 0.0032 \times 65) = 372.24$ ohm. The two-step route matters: scaling directly by the ratio of the two bracketed factors gives the same answer, but scaling by $\alpha \Delta T$ relative to the hot value does not, because the coefficient is defined against $R_0$. A tool that writes SPEF per corner bakes this in; a tool that rescales on the fly must carry the reference temperature with every value.

## What scales and what does not

Only the interconnect resistance in the parasitic netlist follows the wire coefficient. A cell's drive strength is characterized in the library at each corner, and the timer picks the library that matches the corner, so the extractor does not scale it. Consider a net whose Elmore delay at 25 C splits into a 26 ps driver term and a 27 ps wire term. With $\alpha = 0.0032$ per C at 125 C, only the wire term grows, to 35.64 ps, giving 61.64 ps. Wire-dominated nets, the long ones, are the most temperature sensitive. Because transistors may get faster or slower with temperature depending on supply voltage, the hot corner is not automatically the slowest for every path, and signoff runs both ends of the range.

## Going deeper

The temperature sections of an Interconnect Technology Format description in a process design kit; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed., 2003), chapter 4.
