# Pipelining

Pipelining splits a long combinational computation into stages separated by registers, so several operations are in flight at once. The clock period shrinks to the slowest stage plus register overhead, throughput rises, and latency rises slightly. It is the most important timing technique in digital design, used in processors, DSP datapaths and accelerators, and retiming tools automate the placement of pipeline registers.

## Cycle time of a pipeline

Cutting logic into stages with delays $d_1, \dots, d_k$ and inserting a register after each, the clock period is set by the slowest stage plus the register overhead $t_{reg} = t_{cq} + t_{su}$:

$$T_{clk} = \max_i d_i + t_{reg}.$$

Logic of 900 ps split into stages of 250, 200, 300 and 150 ps, with $t_{reg} = 50$ ps, runs at $T_{clk} = 350$ ps instead of $900 + 50 = 950$ ps unpipelined. The imbalance matters: the 300 ps stage sets the clock, and the 150 ps stage wastes half its cycle. Ideal balancing would give $900/4 + 50 = 275$ ps.

## Latency and throughput

A $k$-stage pipeline finishes its first result after $k$ cycles (latency $k\, T_{clk}$) and then one result per cycle. Processing $N$ items takes

$$T_N = (k + N - 1)\, T_{clk}.$$

With $k = 4$ and $T_{clk} = 350$ ps, one item takes 1400 ps, longer than the 950 ps unpipelined, but 1000 items take $1003 \cdot 350 = 351050$ ps, about 351 ns, against 950 ns unpipelined. Throughput approaches $1/T_{clk}$ as $N$ grows; latency is always worse. Which one matters depends on the workload: a streaming filter cares about throughput, a dependent loop cares about latency.

## Speedup and its limits

The speedup on $N$ items is the unpipelined time over the pipelined time, $N (D + t_{reg}) / ((k + N - 1) T_{clk})$, where $D$ is the total logic delay. For the example with $N = 1000$, $950000 / 351050 \approx 2.71$, well below the 4 a naive count suggests, because of imbalance and register overhead.

Deeper pipelines have diminishing returns. As $k$ grows, $T_{clk} \to D/k + t_{reg}$ approaches $t_{reg}$, so speedup saturates near $(D + t_{reg})/t_{reg}$; clock skew and the hazards of the next topic eat further into the gain. Retiming moves registers across logic to balance stages without changing behavior, and an HLS tool chooses the number of stages from the target clock period.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), section 3.6; Leiserson and Saxe, "Retiming Synchronous Circuitry" (Algorithmica, 1991).
