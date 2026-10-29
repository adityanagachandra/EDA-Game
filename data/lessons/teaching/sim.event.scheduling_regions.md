# Verilog scheduling regions

The Verilog and SystemVerilog standard does not leave the order of events within a time step entirely to the simulator. It divides each time step into named **regions** that run in a fixed order. Blocking assignments, nonblocking updates, assertion evaluation and testbench code each have their place. This structure is what makes nonblocking assignments model flip-flops correctly, what makes concurrent assertions race-free, and what a simulator implementer must reproduce exactly.

## The region order of a time step

IEEE 1800 defines, for each time step: **Preponed** (sample values for assertions), **Active** (run processes, execute blocking assignments, evaluate the right-hand sides of nonblocking assignments), **Inactive** (processes suspended with `#0`), **NBA** (apply the left-hand updates of nonblocking assignments), **Observed** (evaluate concurrent assertions), **Reactive** (run program blocks and assertion action blocks), **Re-Inactive** and **Re-NBA** (the testbench counterparts of Inactive and NBA), and **Postponed** (read-only: the strobe and monitor system tasks report here). When NBA updates wake new processes, the scheduler loops back to Active within the same time step. A `#0` delay only moves a process from Active to Inactive; it does not advance time.

## Blocking versus nonblocking results

Start with $a = 3$, $b = 5$ and run on one clock edge `a = b; b = a;`. The first blocking assignment sets $a = 5$ immediately, and the second reads the new $a$, so $b = 5$: both registers hold 5 and the old $a$ is lost. Now run `a <= b; b <= a;`. Both right-hand sides are read in Active (5 and 3), and both updates land in NBA, so $a = 5$ and $b = 3$: the values swap, exactly what two flip-flops do. The same logic explains shift registers: `b <= a; c <= b;` makes $c$ receive the old $b$, a two-stage pipeline, while `b = a; c = b;` makes the old $a$ fall straight through to $c$. Mixed code follows the same rule: a blocking read after a nonblocking assignment still sees the old value, because NBA has not run yet.

## Races between processes

Within the Active region, the standard does not specify the order in which ready processes run. Two `always @(posedge clk)` blocks, one doing `q1 = d;` and the other `q2 = q1;`, race: if the first runs first, $q_2$ gets the new $q_1$; otherwise it gets the old one. Both results are legal, so two compliant simulators can disagree, and one simulator can change behavior after an unrelated edit. The cure is discipline, not delay: nonblocking assignments for every register written in clocked blocks, blocking assignments only in combinational blocks, and testbenches that drive after the edge and sample before it, for example through clocking blocks. Sprinkling `#0` only moves the race into the Inactive region, where two `#0` processes race again.

## Going deeper

IEEE Std 1800-2017, clause 4; Cummings, "Nonblocking Assignments in Verilog Synthesis, Coding Styles That Kill!" (SNUG 2000).
