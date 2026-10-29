# Launch-on-capture and launch-on-shift

A transition test needs two vectors applied at speed, but a scan design can only load one state at a time. Two schemes create the second vector. Launch-on-shift makes it with one more shift; launch-on-capture makes it with a functional clock. They differ in fault coverage, in design cost and in how realistic the tested transitions are, and every at-speed scan flow chooses between them or mixes them.

## Launch-on-shift

In **launch-on-shift** (LOS), also called skewed-load, the last shift of the chain is the launch event. After shifting in $V_1$, one more shift with SE = 1 produces $V_2$: $Q_1$ takes the scan-in bit and each other flop takes its predecessor's value. Then SE drops and an at-speed capture clock follows. For a 3-flop chain with $V_1 = Q_1 Q_2 Q_3 = 110$ and scan-in holding 1, $V_2 = 111$. Because $V_2$ is just a shifted copy of a freely loaded state, ATPG controls both vectors almost independently and LOS reaches high transition coverage. The cost: SE must switch from 1 to 0 within one functional cycle, so it must be timed and buffered like a clock.

## Launch-on-capture

In **launch-on-capture** (LOC), also called broadside, $V_1$ is shifted in slowly, SE drops, and two at-speed functional clocks follow. The first, the launch clock, makes each flop capture its next-state function of $V_1$, so $V_2$ is the circuit's own response. The second captures the result. With next-state functions $D_1 = \overline{Q_2 Q_3}$, $D_2 = Q_1 \oplus Q_2$, $D_3 = Q_1 + Q_3$ and $V_1 = 110$, the launch state is $V_2 = 101$. SE is static during the at-speed pair, so it needs no fast distribution. ATPG must solve a two-time-frame problem: choose $V_1$ so that both $V_1$ and its response satisfy the fault's conditions.

## Choosing between them

LOS usually gives higher coverage but needs at-speed scan enable, and it can test transitions from states the functional circuit never reaches, which may reject good parts (over-testing). LOC is easier to implement and its second vector is a genuine functional response of the first, but it leaves some transition faults untestable. Many flows use LOC as the default, add LOS patterns or hybrid schemes where the design allows, and rely on an on-chip clock controller to generate the precise launch and capture pulses from a PLL, since low-cost testers cannot. Primary inputs are typically held constant across launch and capture.

## Going deeper

Savir and Patil, "Broad-side delay test," IEEE Transactions on CAD, 1994; Savir, "Skewed-load transition test," International Test Conference, 1992; Wang, Wu and Wen (2006), chapter 6.
