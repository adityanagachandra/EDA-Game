# Clock source and network latency

A clock edge does not reach every flop at the same time. It first travels from its true origin, an oscillator or PLL, to the point where the clock is defined, and then through the on-chip clock tree to each sink. These two delays, source latency and network latency, decide when each flop actually launches and captures, and their difference between two flops is skew.

## Clock arrival at a sink

Source latency is the delay from the clock's origin to its definition point, set in SDC with set_clock_latency -source. Network latency is the delay from the definition point through the clock buffers and wires to a flop's clock pin. With propagated clocks, the timer computes it from the actual tree; with ideal clocks, it uses a value from set_clock_latency. The clock arrival at a sink is the sum. For a source latency of 70 ps, shared buffers of 37 and 29 ps and branch buffers of 29 and 32 ps, the edge launched at 0 reaches the flop at $70 + 37 + 29 + 29 + 32 = 197$ ps. In OpenSTA, clocks stay ideal until set_propagated_clock is applied.

## Skew between two sinks

Skew is the capture clock arrival minus the launch clock arrival. The source latency and any shared trunk of the tree appear in both, so they cancel: only the branches after the split point matter. If the launch branch adds 48 ps after the split and the capture branch 36 ps, the skew is -12 ps, whatever the source latency and trunk delays. This is why clock tree synthesis tries to balance branch delays, and why the location of the split point matters for variation: a long shared trunk is harmless for skew, and later the common path is also credited back by clock reconvergence pessimism removal.

## Propagated clocks in the setup check

With propagated clocks, the setup check uses each flop's actual clock arrival: $S = T + L_{cap} - t_{su} - (L_{launch} + t_{cq} + d)$. With a 600 ps period, launch and capture arrivals of 126 and 145 ps, $t_{cq} = 53$, $d = 470$ and $t_{su} = 37$ ps, the slack is 59 ps. Before clock tree synthesis there is no tree to propagate through: the clock net is a single huge unbuffered net, so timers use ideal clocks with an estimated latency and add a skew budget to clock uncertainty. After synthesis the clocks are switched to propagated, and the skew budget is dropped from the uncertainty, so that skew is not counted twice.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 7 on clock latency; the OpenSTA documentation for set_propagated_clock.
