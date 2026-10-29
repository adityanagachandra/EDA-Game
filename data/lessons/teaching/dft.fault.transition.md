# Transition faults

Many real defects do not fix a line at a constant: a resistive via or a weak transistor makes a signal switch too slowly. The transition fault model captures this as a gross delay at one line, slow-to-rise or slow-to-fall, large enough that a transition launched at speed misses the capture edge anywhere it propagates. It is the standard at-speed fault model in production flows because it reuses stuck-at machinery while targeting timing defects.

## Slow-to-rise and slow-to-fall

Each line gets two transition faults. A **slow-to-rise** fault means that a 0-to-1 transition on the line arrives too late to be captured; a **slow-to-fall** fault is the same for 1-to-0. The model assumes the extra delay is large, so the line effectively still holds its old value when the capture edge arrives, no matter which path the transition takes. That makes the fault list exactly as long as the stuck-at list, two per line counting fanout branches, and it lets fault simulation reuse stuck-at logic. A circuit with 12 lines has 24 transition faults.

## Two-pattern tests

A transition fault needs a pair $\langle V_1, V_2 \rangle$. $V_1$ is the **initialization** vector: it sets the line to the starting value, 0 for slow-to-rise and 1 for slow-to-fall. $V_2$ is the **launch** vector, applied at functional speed: it must drive the opposite value and propagate the line's value to an output. During $V_2$ a slow-to-rise line still shows 0, so $V_2$ must be a test for the line stuck-at-0. In short: a slow-to-rise test on $x$ is $V_1$ with $x = 0$ followed by a stuck-at-0 test for $x$. A pair whose first vector already sets $x = 1$ launches no transition, and a pair whose second vector does not sensitize a path from $x$ never shows the delay.

## What the model misses

The transition model is a gross-delay model. A small-delay defect, say 50 ps on a path with 200 ps of slack, fails only if the test happens to propagate through a long enough path. ATPG for transition faults tends to choose the easiest, often shortest, path, so small delays escape. Timing-aware ATPG uses slack from static timing to prefer long paths, and path-delay tests target specific critical paths. In scan designs the two vectors cannot be chosen freely, which leads to the launch-on-capture and launch-on-shift schemes.

## Going deeper

Waicukauski, Lindbloom, Rosen and Iyengar, "Transition fault simulation," IEEE Design and Test, 1987; Bushnell and Agrawal (2000), chapter 12.
