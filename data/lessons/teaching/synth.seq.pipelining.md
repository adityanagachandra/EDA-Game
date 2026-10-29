# Pipelining as retiming

Pipelining adds registers to a feed-forward datapath to raise throughput, at the cost of latency. In synthesis it is implemented as retiming: the designer (or tool) adds $k$ register stages at the inputs or outputs, and retiming spreads them through the logic. Seeing pipelining this way explains both its power and its limits around feedback loops.

## Pipelining versus retiming

Retiming preserves the input-output behavior cycle for cycle; it only moves registers. Pipelining changes the behavior by a fixed delay: outputs appear $k$ cycles later. Adding $k$ registers on every input-to-output path of a feed-forward circuit and then retiming yields a pipeline with up to $k + 1$ stages. Ideally, a 12 ns path with two added stages splits into three 4 ns segments. Synthesis tools support register balancing: write the registers after the logic in RTL, and retiming moves them backward into it. The environment, testbench included, must expect the new latency.

## Loops limit pipelining

In a loop, the number of registers is part of the behavior: adding a register to an accumulator's feedback path changes what it computes. Retiming keeps each loop's register count, so the period is bounded by the loop's delay per register. A loop with 18 ns of delay and 3 registers can never run faster than 6 ns. Pipelining therefore works on feed-forward datapaths and on the feed-forward parts around loops. Processors with feedback through architectural state need design changes, such as forwarding and hazard logic, which are not synthesis transformations at all.

## Diminishing returns

Each pipeline stage pays fixed sequencing overhead: clock-to-q, setup time and clock skew. As stages get shorter, that overhead becomes a large fraction of the period, so throughput gains fall off. Node delays are indivisible in the retiming graph, so a path of logic depth 30 split into 5 stages may leave one stage at depth 9 rather than 6. Every added register also loads the clock tree and toggles every cycle, adding clock power and area. Verification uses sequential equivalence with a latency offset. The right number of stages is an optimization of throughput against overhead, power and latency.

## Going deeper

Leiserson and Saxe, "Retiming synchronous circuitry", Algorithmica (1991), on systolic conversion; Harris, "Skew-Tolerant Circuit Design" (2001), on sequencing overhead.
