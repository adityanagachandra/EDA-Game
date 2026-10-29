# Recovery and removal checks

Asynchronous set and reset pins act immediately, without the clock, but releasing them is not instantaneous in effect: if a reset deasserts too close to a clock edge, the flop may not leave reset cleanly on that edge, and different flops may leave reset on different cycles. Recovery and removal are the setup and hold checks for the deasserting edge of an asynchronous control.

## The recovery required time

Recovery is the setup-like check. The deassertion must reach the asynchronous pin at least the recovery time $t_{rec}$ before the active clock edge at that flop. With a reset synchronized and launched by one clock edge, the next capturing edge arrives at $T + L_{cap}$, so the required time is $T + L_{cap} - t_{rec}$. For a 1000 ps clock, a capture latency of 110 ps and a 73 ps recovery time, the deassertion must arrive by 1037 ps. Liberty models these checks as timing arcs from the clock pin to the reset pin, with timing types recovery_rising and removal_rising, characterized like setup and hold.

## Recovery slack

The deassertion arrives at $L_{launch} + t_{cq} + d_{max}$, where $d_{max}$ is the longest path of the reset net, often a large buffered tree reaching thousands of flops. Recovery slack is required minus arrival. With a required time of 1127 ps and an arrival of $60 + 44 + 680 = 784$ ps, the slack is 343 ps. Large reset trees make recovery a real constraint: the reset fan-out can rival a clock tree in size, but it is buffered as data. If recovery fails, flops leave reset on different cycles, and a state machine can start in an illegal state. Designs therefore deassert reset synchronously through a reset synchronizer and time the tree like any data path.

## Removal slack

Removal is the hold-like check. The deassertion must not arrive so soon after the clock edge that it falls inside the removal window: $S_{rem} = (L_{launch} + t_{cq} + d_{min}) - (L_{cap} + t_{rem})$. With $140 + 79 + 85 = 304$ ps against $120 + 52 = 172$ ps, the slack is 132 ps. As with hold, the period does not appear. Only deassertion is checked. Assertion takes effect asynchronously regardless of the clock, so a timer has nothing to compare it against; whatever logic consumes the reset must tolerate its asynchronous arrival. Recovery and removal appear in reports alongside setup and hold, in their own path groups.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 8 on recovery and removal; Cummings and Mills, "Synchronous Resets? Asynchronous Resets? I am so confused!" (SNUG 2002).
