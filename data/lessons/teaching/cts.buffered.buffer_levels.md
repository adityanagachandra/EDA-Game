# Buffer levels and fanout in a clock tree

An unbuffered wire tree works only over short spans. A real clock tree is a fanout tree of buffers, each driving a bounded number of loads, with a root buffer at the top. The fanout limit comes from slew and load constraints, and it sets the depth of the tree, its buffer count and its latency. Choosing the fanout is the first trade-off a buffered CTS makes: shallow trees with slow stages, or deep trees with fast ones.

## Levels from fanout

If no buffer may drive more than $F$ loads, $L$ levels of buffers reach at most $F^L$ sinks, so a tree for $N$ sinks needs $L = \lceil \log_F N \rceil$ levels. For 8092 flops and $F = 6$, $6^5 = 7776 < 8092 \le 6^6 = 46\,656$, so six levels are needed. The level count matters beyond latency: every level adds a stage whose delay varies independently on each branch, so deep trees accumulate more uncorrelated on-chip variation between sinks.

## Counting buffers

Working up from the sinks, the bottom level needs $\lceil N/F \rceil$ buffers, the next $\lceil \lceil N/F \rceil / F \rceil$, and so on until one root remains. For 2549 flops and $F = 3$, the counts are 850, 284, 95, 32, 11, 4, 2, 1, a total of 1279 buffers. The total is close to $N/(F - 1)$, so halving the fanout from 6 to 3 roughly multiplies the buffer count by 2.5. Buffers are clock power: each toggles every cycle, so buffer count enters the power objective directly.

## The fanout that minimises latency

A buffer's delay grows with its load. Model a stage driving $F$ loads as $a + bF$ ps, where $a$ is its intrinsic delay and $b$ the cost per load. Latency is $\lceil \log_F N \rceil (a + bF)$. For $N = 1000$, $a = 10$ and $b = 6$: fanout 2 needs 10 levels at 22 ps (220 ps), fanout 4 needs 5 levels at 34 ps (170 ps), fanout 8 needs 4 levels at 58 ps (232 ps). The optimum sits at a moderate fanout, the same result logical effort gives for optimal buffer chains, where the best stage effort is around 4. The ceiling makes the curve jagged, so tools evaluate candidates rather than solve for $F$.

## Going deeper

Sutherland, Sproull and Harris, "Logical Effort: Designing Fast CMOS Circuits" (1999); Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the clock chapter.
