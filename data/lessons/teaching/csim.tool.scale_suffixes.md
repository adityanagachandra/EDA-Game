# SPICE values and scale suffixes

SPICE numbers are written with scale suffixes: 2.2k, 470p, 1meg. The parser reads a number, then the letters after it, takes the first one or three letters as a scale factor and ignores the rest as decoration. That convention is compact, but it hides two traps that have cost designers many hours: M means milli, not mega, and a trailing F means femto. A simulator author implements this parser; a user must read it correctly.

## Reading scale suffixes

The SPICE scale factors are T ($10^{12}$), G ($10^{9}$), MEG ($10^{6}$), K ($10^{3}$), M ($10^{-3}$), U ($10^{-6}$), N ($10^{-9}$), P ($10^{-12}$) and F ($10^{-15}$). The value is the mantissa times the factor, in the base unit of the element: ohms for R, farads for C, henries for L. So 4.7k is 4700 ohms, 22n on a capacitor is $22 \times 10^{-9}$ F, which is 22000 pF, and 330n on an inductor is 0.33 uH. Converting into the unit a question asks for is a ratio of powers of ten: 100u expressed in nanofarads is $100 \times 10^{-6}/10^{-9} = 100000$ nF.

## The M and MEG trap

SPICE suffixes are case-insensitive, so M and m are the same suffix, and both mean milli. A one-megohm resistor must be written 1meg (or 1e6). A card R1 a b 2M is a 0.002 ohm resistor, nearly a short, and the circuit will simulate happily with absurd currents. The second trap is that letters after a recognized suffix are ignored: 10pF reads as 10p, 3.3uH as 3.3u, and 5MEGohm as 5 meg. But the letter F alone is the femto suffix, so a capacitor written 10F is 10 fF, not 10 farads. Likewise 1Mohm is one milliohm. A careful parser warns about such cases; a careful reader converts them before trusting a result.

## Computing with netlist values

Once parsed, values feed formulas directly, and a useful habit is to pair units so that products land in convenient units. Kilohms times picofarads is nanoseconds; ohms times femtofarads is femtoseconds; one over kilohms is millisiemens. For the cards R1 in out 2.2k and C1 out 0 470p, the time constant is $2.2 \times 10^{3} \times 470 \times 10^{-12} = 1.034 \times 10^{-6}$ s, which is 1034 ns. The same circuit with R1 written 2.2M would give 0.00103 ns: a factor of a billion, from one letter. Simulators echo parsed values in their logs, and checking them is cheap.

## Going deeper

The ngspice manual, the section on number formats and scale factors; Vladimirescu, "The SPICE Book" (Wiley, 1994).
