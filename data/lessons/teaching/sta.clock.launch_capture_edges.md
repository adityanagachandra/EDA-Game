# Setup and hold relationships between clocks

When launch and capture flops are clocked by different but synchronous clocks, there is no single period to plug into the setup check. The timer must find, among all pairs of launch and capture edges, the most restrictive pair for setup, and then derive the hold check from it. This edge arithmetic underlies every multi-frequency path and every multicycle exception.

## Integer-multiple periods

Take two clocks from one source, both rising at 0. For each launch edge, the setup capture edge is the nearest capture edge strictly later. If the launch clock has period 4 ns and the capture clock 16 ns, the launch edges at 0, 4, 8 and 12 all pair with the capture edge at 16, with gaps of 16, 12, 8 and 4. The most restrictive setup relationship is 4 ns, from the last launch edge before capture. From slow to fast, launch at 16 ns and capture at 4 ns, the launch at 0 pairs with the capture at 4, and the setup relationship is 4 ns. In both directions the faster clock's period sets the setup relationship unless a multicycle exception says the design waits longer.

## General periods over the common period

For periods that are not multiples, the edge pattern repeats after the least common multiple. With launch period 9 ns and capture 6 ns, the common period is 18 ns, and the launch edges at 0 and 9 pair with capture edges at 6 and 12: gaps of 6 and 3, so the setup relationship is 3 ns. With periods 10 and 4 the common period is 20 ns, launches at 0 and 10 pair with 4 and 12, and the relationship is 2 ns. Clocks whose periods have no small common multiple produce tiny, unrealistic relationships; that usually means they should be declared asynchronous rather than timed.

## Deriving the hold relationship

Hold is derived from the most restrictive setup pair, launch at $t$ and capture at $c$. Two checks follow. Data launched at $t$ must not be captured by the capture edge one period before, at $c - T_c$, giving $(c - T_c) - t$. Data launched by the next launch edge, at $t + T_l$, must not be captured by $c$, giving $c - (t + T_l)$. The larger of the two is the hold relationship. For launch period 3 and capture period 8, the tightest setup pair is launch at 15 and capture at 16; the two values are $8 - 15 = -7$ and $16 - 18 = -2$, so hold is checked with -2 ns. For equal periods both values are 0, the familiar same-edge hold.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 8 on multi-frequency paths; the PrimeTime and OpenSTA documentation on clock relationships.
