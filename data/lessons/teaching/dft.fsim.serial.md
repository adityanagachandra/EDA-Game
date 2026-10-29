# Serial fault simulation and fault dropping

The simplest fault simulator copies the circuit, injects one fault, simulates every pattern, and repeats for every fault. Serial fault simulation is the reference against which every faster method is checked, and it shows clearly where the cost goes. Its single most important optimization, fault dropping, also appears in every production simulator, so the counting here carries over directly.

## The serial algorithm

Serial fault simulation simulates the fault-free circuit once per pattern to get the expected responses. Then, for each fault, it simulates the faulty circuit on the patterns and compares outputs. One **run** is one faulty circuit on one pattern. Without any optimization, $F$ faults and $P$ patterns cost $F \times P$ runs plus $P$ fault-free runs. Each run is an ordinary logic simulation, so a compiled or levelized simulator can be reused unchanged; that reuse is the main virtue of the serial method. Its weakness is that it repeats work: most of a faulty circuit behaves exactly like the fault-free one.

## Fault dropping

When the goal is coverage rather than diagnosis, a fault can be **dropped** as soon as one pattern detects it. A fault first detected by pattern $k$ then costs $k$ runs, and an undetected fault costs all $P$. With 6 faults and 5 patterns where three faults are first detected at patterns 2, 1 and 2 and three are never detected, the runs are $2 + 1 + 2 + 5 + 5 + 5 = 20$, against 30 without dropping. Because coverage curves rise steeply, most faults are detected by the first few patterns, and dropping typically removes the bulk of the work. The remaining cost is dominated by the hard and undetectable faults, which must be simulated against every pattern.

## Coverage as a by-product

The first-detection table also gives the **coverage curve**: the fraction of faults detected after the first $k$ patterns. Plotting it shows diminishing returns, which is how test engineers decide where a pattern set stops paying for its tester time. Dropping is turned off for **diagnosis**, which needs the full pass/fail behaviour of every fault on every pattern to build a fault dictionary, and for **n-detect** test, which wants each fault detected several times by different patterns to raise the chance of catching unmodeled defects.

## Going deeper

Bushnell and Agrawal (2000), chapter 5; Abramovici, Breuer and Friedman (1990), chapter 5.
