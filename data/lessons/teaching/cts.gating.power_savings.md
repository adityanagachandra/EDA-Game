# Clock-gating power savings

Gating pays when the clock it removes outweighs the cost of the gates. The savings depend on how often each enable is true, how much capacitance sits behind each gate, and how much always-toggling capacitance each gate adds on the ungated side. These are quick computations, and they explain where synthesis inserts gates and why it refuses to gate single flops.

## Power of a gated subtree

A gated subtree toggles only in enabled cycles, so its average power is $p_{en} C V^2 f$, where $p_{en}$ is the fraction of cycles in which the enable is true. A gate driving 160 pF of subtree, enabled 40% of the time, at 1.75 GHz and 0.7 V, averages $0.4 \cdot 160 \cdot 0.49 \cdot 1.75 = 54.9$ mW instead of 137.2 mW. The enable probability comes from simulation or from activity propagation; a pessimistic power estimate assumes enables are true more often than they really are.

## Net savings with gate overhead

Each ICG presents its clock input and internal latch to the ungated clock, which toggles every cycle. With 160 gates of 4 fF each the overhead is 0.64 pF. If they gate 165 pF in total and the enables are true 60% of the time, the removed average load is $0.4 \cdot 165 = 66$ pF, so the net saving at 0.5 GHz and 0.7 V is $(66 - 0.64) \cdot 0.49 \cdot 0.5 = 16.0$ mW. The enable logic itself also switches and is part of a full accounting, as are the clock buffers behind the gate, which stop toggling too, so gating higher in the tree saves more.

## The break-even bank width

Gating $n$ flops with clock pins of $C_{pin}$ removes $(1 - p_{en})\,n\,C_{pin}$ of average load and adds the gate's overhead $C_g$. Gating saves power only when $n > C_g/((1 - p_{en})C_{pin})$. For $C_g = 2$ fF, $C_{pin} = 0.9$ fF and $p_{en} = 0.3$, $n > 3.17$, so at least four flops must share the gate. Real thresholds are higher, because the enable logic and the gate's area count too. Synthesis tools therefore have a minimum bank width for gating, and they merge banks with identical enables under one gate.

## Going deeper

Keating, Flynn, Aitken, Gibbons and Shi, "Low Power Methodology Manual" (2007); Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed., 2003).
