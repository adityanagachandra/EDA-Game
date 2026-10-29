# Latches and time borrowing

A level-sensitive latch is transparent while its clock is active, so data may arrive after the latch opens and still pass through. A path that arrives late at a latch borrows time from the next stage. Timers must model this to avoid reporting false violations on latch-based designs, and they must charge the borrowed time to the following path so that nothing is counted twice.

## Borrowing past the opening edge

Take a positive latch clocked with period $T$ and a 50 percent duty cycle, fed by a rising-edge flop that launches at 0. The latch's nominal capture point is its opening edge at $T$. If data arrives at $A \le T$, it waits for the latch to open: no borrowing. If it arrives later, while the latch is transparent, it passes straight through and the path borrows $A - T$. With $T = 1100$ ps and $A = 1550$ ps, the borrow is 450 ps. To the path ending at the latch the borrowing is not a violation; reports show it as time borrowed with zero slack, as long as the arrival stays within the borrowing limit.

## The maximum borrow

Data must still settle a setup time before the latch closes, at $T + H$ where $H$ is the high time. The latest legal arrival is $T + H - t_{su}$, so the most a path can borrow is the pulse width minus setup, $H - t_{su}$. With $H = 700$ ps and $t_{su} = 38$ ps, the limit is 662 ps. Beyond it the path truly fails setup at the closing edge. Timers allow the limit to be reduced, for example with set_max_time_borrow, so that a designer can keep borrowing within a budget. Clock skew and uncertainty at the latch adjust the opening and closing edges just as they adjust a flop's capture edge.

## Charging the next stage

Borrowed time is taken from the path that starts at the latch. When borrowing, data leaves the latch at $A + d_{DQ}$, the arrival plus the latch's D-to-Q delay, rather than at the opening edge plus clock-to-Q. The next path is timed from there. If the latch's output feeds 220 ps of logic to a flop capturing at $2T = 2000$ ps with a 50 ps setup, and $A = 1200$, $d_{DQ} = 60$ ps, the next stage's slack is $2000 - 50 - (1200 + 60 + 220) = 470$ ps. Borrowing passes slack along a pipeline: a long stage can take time from a short neighbor. A timer must propagate this across chains of latches, which makes latch-based timing iterative rather than a single pass.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 8 on time borrowing; Sakallah, Mudge and Olukotun, "Analysis and design of latch-controlled synchronous digital circuits" (DAC 1990).
