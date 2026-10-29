# Input and output delays

A block is timed in isolation, but its ports connect to logic outside it. Input and output delays describe that outside logic as time budgets relative to a clock, so paths that enter or leave the block get proper setup checks. Choosing these budgets well is how a chip-level timing problem is split into block-level problems that can be closed independently.

## Input paths

set_input_delay 410 -clock clk [get_ports din] says data at din arrives 410 time units after the clk edge that launched it outside the block. Inside, the path from din to a flop is timed like any flop path, with the input delay in place of the launch flop's clock-to-Q and external logic: $S = T - t_{su} - (\text{input delay} + d)$ with ideal clocks. For a 1400 ps clock, an input delay of 410 ps, 180 ps of internal logic and a 33 ps setup, the slack is 777 ps. A -min input delay sets the earliest arrival for hold checks on the same port, and -clock_fall references the falling edge.

## Output paths

set_output_delay 300 -clock clk [get_ports dout] says the logic and setup outside the block need 300 units of the cycle after dout. The port's required time is the capture edge minus the output delay. A flop inside drives dout through logic: $S = T - \text{output delay} - (t_{cq} + d)$. With a 900 ps clock, a 300 ps output delay, a 44 ps clock-to-Q and 180 ps of logic, dout is required at 600 ps and the slack is 376 ps. Output delay is a requirement, not an arrival, and confusing the two is a common mistake when reading interface specifications, which often give setup and hold times at a device's pins.

## Budgeting a feedthrough

A purely combinational path from an input port to an output port is bounded by both budgets: it may use only $T - \text{input delay} - \text{output delay}$. For a 500 ps clock with an input delay of 190 ps and an output delay of 110 ps, the feedthrough may take 200 ps. At the chip level these budgets come from timing budgeting: the top-level path is split among the blocks it crosses, and each block gets I/O delays consistent with its share. If budgets are too generous to one block, another block inherits an impossible constraint. This is a classic failure-propagation point in hierarchical flows, and budgets are refined as the blocks are implemented.

## Going deeper

Gangadharan and Churiwala, "Constraining Designs for Synthesis and Timing Analysis" (2013), chapters on I/O constraints; Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 7.
