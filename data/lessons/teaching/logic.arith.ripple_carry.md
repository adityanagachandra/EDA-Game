# Ripple-carry adders

The ripple-carry adder chains $n$ full adders, each carry-out feeding the next carry-in. It is the smallest adder and the slowest: its delay grows linearly with width. It is the baseline against which every faster adder is measured, and synthesis tools still choose it when timing allows, because area and power are lowest.

## The worst-case carry chain

In the worst case a carry generated in column 0 must ripple through every column to the top. Model each full adder's carry-out as settling $d_c$ gate delays after its carry-in (or after its operand inputs, for column 0, whose carry-in is ready at time 0), and its sum as settling $d_s$ gate delays after its carry-in. Then $c_{i+1}$ settles at $(i + 1) d_c$, the carry-out of an $n$-bit adder at $n\,d_c$, and the most significant sum bit at

$$t_{s_{n-1}} = (n - 1)\, d_c + d_s.$$

With $d_c = 2$ and $d_s = 2$, a 32-bit ripple adder's top sum settles after $62 + 2 = 64$ gate delays, and its carry-out after 64.

## Converting to time and sizing to a budget

Gate delays become time once a technology is fixed. If a gate delay is 15 ps, the 32-bit adder above takes $64 \cdot 15 = 960$ ps, too slow for a 1 GHz clock once register overhead is added.

Turn the formula around to find the widest ripple adder that fits a budget $B$: the largest $n$ with $(n - 1) d_c + d_s \le B$. With $d_c = 2$, $d_s = 3$ and a budget of 40 gate delays, $(n - 1) \cdot 2 \le 37$ gives $n = 19$. Synthesis tools make exactly this kind of decision when they pick an adder implementation per instance: ripple where the slack allows, a faster structure where it does not.

## Average versus worst case

The worst case is rare for random operands. A carry ripples only through a run of propagate columns ($a_i \ne b_i$) that starts at a generating column. For uniformly random operands the expected longest such run is about $\log_2 n$, so a typical addition settles long before the worst case.

Synchronous design cannot exploit this: the clock period must cover the worst case, because any operand pair might occur and static timing analysis is input-independent. Self-timed adders with completion detection do exploit it, at the cost of extra logic and an asynchronous handshake. In synchronous logic the response is architectural: shorten the worst-case carry path itself, which is what carry-lookahead, carry-select and prefix adders do.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), section 5.2.1; Koren, "Computer Arithmetic Algorithms" (2nd ed., 2002), chapter 5.
