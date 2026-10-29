# Clock uncertainty in setup and hold

set_clock_uncertainty is the SDC knob that reserves margin on clock edges. Setup uncertainty makes the capturing edge effectively earlier; hold uncertainty makes the hold requirement larger. What the number should contain changes during the flow: before clock tree synthesis it stands in for skew that does not exist yet, and after it only jitter and margin remain.

## Building a pre-CTS uncertainty

Before CTS, clocks are ideal and every flop sees the same edge, so the timer knows nothing about skew. Teams therefore budget setup uncertainty as the sum of the expected jitter, an estimate of the skew the tree will add, and a design margin. With 10 ps of jitter, an 85 ps skew estimate and 15 ps of margin, the command is set_clock_uncertainty -setup 110 [get_clocks clk]. After CTS, clocks become propagated and the timer computes real skew, so the skew estimate must be removed from the uncertainty; leaving it in counts skew twice and makes the design look worse than it is. Uncertainty can also be set between specific clock pairs, for crossings with extra jitter.

## Setup slack with uncertainty

Setup uncertainty $U_s$ is subtracted from the required time: $S = T - U_s - t_{su} - (t_{cq} + d_{max})$ with ideal clocks. For $T = 700$, $U_s = 41$, $t_{su} = 39$, $t_{cq} = 43$ and $d_{max} = 310$ ps, the slack is 267 ps. Path reports list the uncertainty as its own line in the required-time section, so it can be audited. Because it applies to every path of a clock, a too-large uncertainty is an expensive mistake: it makes the optimizer upsize and buffer thousands of paths that would have passed, costing area and power.

## Hold slack with uncertainty

Hold uncertainty $U_h$ is added to the hold requirement: $S = (t_{cq} + d_{min}) - (t_h + U_h)$ with ideal clocks. For $t_{cq} = 30$, $d_{min} = 99$, $t_h = 42$ and $U_h = 13$ ps, the slack is 74 ps. Hold uncertainty is usually smaller than setup uncertainty, because cycle-to-cycle jitter does not apply to a same-edge check; what remains covers the uncertainty of one edge at two sinks and modeling error. Pre-CTS, a hold skew budget is sometimes included, but hold fixing is generally postponed until after CTS, when real skew is known, so buffers are not inserted against skew that never materializes.

## Going deeper

The Synopsys Design Constraints format specification, set_clock_uncertainty; Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 7.
