# Setup and hold checks

Every flop-to-flop path is checked twice. The setup check asks whether the slowest data reaches the capture flop before the next capturing edge; the hold check asks whether the fastest data stays away long enough not to corrupt the value being captured on the current edge. Both use the same arrival and clock numbers, combined in opposite directions, and a timer reports each as a slack.

## The setup required time

Data launched by the edge at 0 must be captured by the next edge, one period $T$ later. The capture flop's clock pin sees that edge after its clock latency $L_{cap}$, and the data must be stable a setup time $t_{su}$ before it. The setup required time at the D pin is therefore $RAT = T + L_{cap} - t_{su}$. For an 800 ps period, a capture latency of 185 ps and a 29 ps setup time, data must arrive by $800 + 185 - 29 = 956$ ps. Clock uncertainty, when specified, is subtracted too. The launch side contributes nothing to the required time: it determines the arrival, not the deadline.

## Setup slack

The latest arrival at the D pin is the launch clock latency plus clock-to-Q plus the longest combinational delay: $AT = L_{launch} + t_{cq} + d_{max}$. Setup slack is the required time minus this arrival:
$$S_{setup} = (T + L_{cap} - t_{su}) - (L_{launch} + t_{cq} + d_{max}).$$
With $T = 1100$, $L_{cap} = 170$, $t_{su} = 36$, $L_{launch} = 140$, $t_{cq} = 74$ and $d_{max} = 410$ ps, the slack is $1234 - 624 = 610$ ps. The difference $L_{cap} - L_{launch}$ is the clock skew; positive skew adds setup margin. Setup failures are fixed by making the long path faster, lengthening the period, or shifting the clock edges.

## Hold slack

Hold compares the same capturing edge that the data was launched with. New data launched at time 0 must not reach the capture flop before the flop has held its old value for the hold time after its own clock arrival:
$$S_{hold} = (L_{launch} + t_{cq} + d_{min}) - (L_{cap} + t_h).$$
With $L_{launch} = 135$, $t_{cq} = 61$, $d_{min} = 26$, $L_{cap} = 125$ and $t_h = 28$ ps, the slack is $222 - 153 = 69$ ps. The period $T$ does not appear, so slowing the clock never fixes hold. Hold uses the shortest path and minimum delays, usually at the fast corner, and positive skew hurts it. The fix is to add delay to the short path, which is why hold fixing comes after setup fixing.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 8; Harris and Harris, "Digital Design and Computer Architecture", the section on timing of sequential logic.
