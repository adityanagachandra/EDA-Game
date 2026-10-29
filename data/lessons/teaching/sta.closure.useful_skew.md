# Useful skew

Zero skew is not the goal of a clock tree; meeting timing is. Deliberately delaying or advancing the clock of a flop moves time between the stages before and after it. Useful skew uses this to close a failing stage with slack borrowed from its neighbor, without touching the data path, and it is limited by the hold checks it erodes.

## Shifting slack between two stages

In a pipeline FF1, FF2, FF3, delay FF2's clock by $\delta$. The stage into FF2 captures later, so its setup slack grows by $\delta$; the stage out of FF2 launches later, so its setup slack shrinks by $\delta$. With slacks of -53 ps into FF2 and 93 ps out of it, a shift of $\delta = 95$ ps gives 42 and -2 ps: the shift was too large. The design's slack is the worse of the two, so the goal is to choose $\delta$ that maximizes the minimum. A negative $\delta$, advancing FF2's clock, moves time the other way.

## The balancing shift

The minimum of $s_A + \delta$ and $s_B - \delta$ is largest when the two are equal: $\delta^* = (s_B - s_A)/2$, which leaves both at $(s_A + s_B)/2$. For $s_A = -11$ and $s_B = 77$ ps, $\delta^* = 44$ ps and both stages end with 33 ps. If the average is negative, no skew can fix both stages; the pair needs data-path work. Across a whole design, choosing clock arrivals to maximize the worst slack over all flop pairs is a linear program in the arrival times, closely related to retiming and to finding cycles of minimum mean slack, and clock tree tools approximate it with targeted delays.

## The hold limit

Delaying FF2's clock also hurts hold on the stage into FF2: the capture edge moves later, so the hold slack of those paths falls by $\delta$. If that hold slack is $h_A$, the usable shift is at most $h_A$, so the best choice is $\delta = \min(\delta^*, h_A)$. With $s_A = -23$, $s_B = 111$ ps, $\delta^* = 67$ ps and $h_A = 78$ ps, the balance point is allowed and the worse setup slack becomes 44 ps; with $h_A = 30$ ps, the shift stops at 30 ps and the worse setup slack is 7 ps. Hold on the stage out of FF2 improves at the same time, which useful-skew optimizers exploit as well.

## Going deeper

Fishburn, "Clock skew optimization" (IEEE Transactions on Computers, 1990); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on clock tree synthesis.
