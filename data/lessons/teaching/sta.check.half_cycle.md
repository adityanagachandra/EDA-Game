# Half-cycle paths between opposite edges

Not every path runs between rising edges. A path from a rising-edge flop to a falling-edge flop, or the reverse, has only part of a cycle for setup, and an unusually easy hold check. Such half-cycle paths appear in DDR interfaces, in negative-edge logic and around latches, and a timer derives their checks from the clock waveform rather than from the period alone.

## Time between launch and capture edges

Let a clock rise at 0 and fall at $H$, with period $T$. For a rise-launch, fall-capture path, the setup capture edge is the first falling edge after launch, at $H$: the path has $H$ instead of $T$. For a fall-launch, rise-capture path, data leaves at $H$ and is captured by the rising edge at $T$, a window of $T - H$. With a 50 percent duty cycle both are $T/2$. With $T = 1600$ ps and $H = 784$ ps, a fall-to-rise path gets 816 ps. Any duty-cycle distortion in the clock tree directly shrinks one of these windows, so half-cycle paths are sensitive to clock waveform quality.

## Setup slack of a half-cycle path

The check is the usual setup check with the shorter window: $S = W - t_{su} - (t_{cq} + d_{max})$, where $W$ is the launch-to-capture window. With $W = 600$ ps, $t_{su} = 47$, $t_{cq} = 50$ and $d_{max} = 240$ ps, the slack is 263 ps. The same logic between two rising-edge flops would have had 600 ps more. Half-cycle paths therefore carry less logic, and a designer who moves a flop to the opposite edge to fix hold or to retime a signal must check that the remaining logic fits the half cycle. Clock latencies and uncertainty enter exactly as in a full-cycle check.

## Hold is easy for half-cycle paths

Hold is checked against the capturing edge just before the setup capture edge. For rise-to-fall that is the falling edge at $H - T$, before the launch at 0. The hold slack is $(t_{cq} + d_{min}) - ((H - T) + t_h)$, which is large: with $T = 1300$, $H = 585$, $t_{cq} = 63$, $d_{min} = 35$ and $t_h = 28$ ps, it is $98 - (-715 + 28) = 785$ ps. For fall-to-rise, the hold edge is the rise at 0, $H$ before the launch. This is why designers sometimes insert a negative-edge flop or latch to repair hold on paths with large clock skew: the half-cycle offset turns a hold problem into a setup budget.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 8 on half-cycle paths.
