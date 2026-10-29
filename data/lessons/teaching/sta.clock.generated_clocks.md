# Generated clocks

Many clocks on a chip are derived from others: divided by a counter, multiplied by a PLL, or shaped by selecting particular master edges. SDC describes them with create_generated_clock, which ties the new clock to a master. That tie lets the timer compute the generated clock's edges and its latency through the generating logic, so paths between master and generated domains are timed with the correct relationship.

## Dividing a master clock

create_generated_clock -name div2 -source [get_ports clk] -divide_by 2 [get_pins U1/Q] declares that the clock at the divider flop's output has twice the master's period. In general -divide_by $N$ gives period $N T$ and -multiply_by $M$ gives $T/M$. A 4 ns master divided by 8 is a 32 ns clock. Defining the divider output with a separate create_clock instead would create an unrelated clock: paths between clk and the divided domain would lose their phase relationship, and the latency through the divider would not be computed. The generated clock's edges are aligned to the master edges that produce them, which create_generated_clock derives automatically.

## Building a waveform from master edges

The -edges option builds the generated waveform from numbered master edges: edge 1 is the master's first rise, edge 2 its first fall, edge 3 its second rise, and so on. A list {e1 e2 e3} means the generated clock rises at master edge e1, falls at e2 and rises again at e3. For a 4 ns master with waveform {0 2}, edges are at 0, 2, 4, 6, 8, 10, 12. Then -edges {1 4 7} rises at 0, falls at 6 and rises at 12: period 12 ns, high for 6 ns, a divide-by-3 with 50 percent duty cycle. The high time is $t(e_2) - t(e_1)$ and the period $t(e_3) - t(e_1)$. -edge_shift can offset each edge further.

## Latency through the generating logic

With propagated clocks, a generated clock's source latency is the master's source latency plus the delay from the master's definition point through its network to the divider and through the divider itself. A sink in the generated domain then adds the generated clock's own network. With a master source latency of 200 ps, 380 ps of master network to the divider flop, a 68 ps clock-to-Q and 390 ps of generated network, the edge reaches the sink at 1038 ps. If a master-domain flop near the divider sees its clock at about 580 ps, a path from it to that sink has about 458 ps of skew, which the timer computes only because the generated clock was declared relative to its master.

## Going deeper

The Synopsys Design Constraints format specification, create_generated_clock; Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 7.
