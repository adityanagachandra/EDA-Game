# Propagation delay through a gate network

Every gate takes time to respond. The delay of a combinational block is set by its slowest input-to-output path, and its fastest path matters too, for hold checks. Computing both is a longest and shortest path problem on a directed acyclic graph, the same propagation that static timing analysis performs at scale. Learning it on a ten-gate network makes the timer's core loop concrete.

## Arrival times in topological order

Model the netlist as a DAG: a node per gate, an edge from each driver to each gate it feeds. Primary inputs arrive at time 0 (or at stated times). A gate's output arrival time is

$$t_{out} = \max_{\text{inputs } i} t_i + d_{gate}.$$

Process gates in topological order so each input is known before it is used. With $g_1 = \text{NAND}(a, b)$ of delay 20 ps, $g_2 = \text{NOR}(g_1, c)$ of delay 30 ps and $g_3 = \text{XOR}(g_1, g_2)$ of delay 40 ps: $t_1 = 20$, $t_2 = \max(20, 0) + 30 = 50$, $t_3 = \max(20, 50) + 40 = 90$ ps. The critical path is $a \to g_1 \to g_2 \to g_3$.

## Inputs that arrive late

Primary inputs rarely all arrive at time 0: some come from slow upstream logic or from another block's output. The propagation rule is unchanged; only the initial values differ. In the network above, if $c$ arrives at 35 ps then $t_2 = \max(20, 35) + 30 = 65$ ps and $t_3 = \max(20, 65) + 40 = 105$ ps, so the critical path now starts at $c$.

This is why a timer stores an arrival time per node rather than a single block delay: the critical path depends on the context the block sits in. Hierarchical timing models summarize a block as input-to-output delays per pin pair for exactly this reason.

## The shortest path and why it matters

The output does not hold its old value until the critical path settles. It can start to change as soon as the fastest path delivers a new input: the contamination delay, computed with $\min$ in place of $\max$:

$$t^{\min}_{out} = \min_i t^{\min}_i + d^{\min}_{gate}.$$

For the example, the shortest path is $a \to g_1 \to g_3$: $20 + 40 = 60$ ps, against a longest path of 90 ps. A register fed by this block sees its input change from 60 ps onward and settle at 90 ps. Setup checks use the maximum; hold checks use the minimum. A network with a large gap between the two forces a designer to add delay on short paths without slowing long ones, which is why a timer carries both.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), section 2.9; Cormen et al., "Introduction to Algorithms" (4th ed., 2022), the section on shortest paths in DAGs.
