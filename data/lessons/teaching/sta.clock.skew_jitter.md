# Clock skew and jitter

Two effects move clock edges away from their ideal times: skew, a fixed spatial difference in arrival between two flops, and jitter, a random temporal variation of the edges themselves. They enter the setup and hold checks differently, and confusing them, or counting either twice, is a classic source of wrong slack.

## Signed skew

Skew is the capture clock arrival minus the launch clock arrival, $\text{skew} = L_{cap} - L_{launch}$. It is signed. With a launch arrival of 190 ps and a capture arrival of 85 ps, the skew is -105 ps: the capture flop sees the edge early. Positive skew gives the data more time for setup and less margin for hold; negative skew does the reverse. Skew is a property of a launch and capture pair, not of a single flop, so one flop can have positive skew with respect to one neighbor and negative skew with respect to another. Global skew, the largest difference across all sinks, is a tree-quality metric but does not enter any single check.

## Setup slack with skew and jitter

Cycle-to-cycle jitter is the variation of one clock period relative to the previous one. A setup check spans one cycle, from a launch edge to the next capture edge, so jitter $J$ can shorten that cycle and is subtracted:
$$S_{setup} = T - J + \text{skew} - t_{su} - t_{cq} - d_{max}.$$
With $T = 800$, $J = 37$, skew -35, $t_{su} = 53$, $t_{cq} = 46$ and $d_{max} = 260$ ps, the slack is 369 ps. In SDC, jitter is usually folded into set_clock_uncertainty rather than modeled separately, together with a margin and, before CTS, an estimate of skew.

## Hold slack with skew

Hold compares launch and capture on the same edge, so the period and the cycle-to-cycle jitter between consecutive edges do not enter it:
$$S_{hold} = t_{cq} + d_{min} - t_h - \text{skew}.$$
With $t_{cq} = 79$, $d_{min} = 75$, $t_h = 33$ and a skew of -15 ps, the slack is 136 ps. Positive skew subtracts directly from hold slack, which is why flops with short paths between them, such as shift registers, are sensitive to skew: a capture clock that arrives late lets the new data race through. Hold uncertainty in SDC adds margin for effects such as the jitter of a single edge seen differently at two sinks.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 7; Friedman, "Clock distribution networks in synchronous digital integrated circuits" (Proceedings of the IEEE, 2001).
