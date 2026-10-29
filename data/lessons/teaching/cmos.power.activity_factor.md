# Signal probability and activity factors

The activity factor is the term in $\alpha C V^2 f$ that a power tool must estimate from the design itself. It is computed either from simulation (counting toggles in a waveform) or by propagating signal probabilities through the logic. The probabilistic method is the one tools use when no simulation exists, and it has known blind spots that a tool designer must understand.

## Signal probabilities through gates

Let $p$ be the probability that a signal is 1 in a cycle. For independent inputs, gate output probabilities follow from the truth table: AND2 gives $p_A p_B$, NAND2 gives $1 - p_A p_B$, OR2 gives $1 - (1 - p_A)(1 - p_B)$, NOR2 gives $(1 - p_A)(1 - p_B)$ and XOR2 gives $p_A(1 - p_B) + p_B(1 - p_A)$. Example: a NAND2 with $p_A = 0.5$ and $p_B = 0.8$ has output probability $1 - 0.4 = 0.6$. Propagating these probabilities through a netlist in topological order gives every node's static probability in one linear pass, which is how vectorless power estimation begins.

## From probability to transition probability

Dynamic power needs the probability of a rising transition, not of being 1. If successive cycles are independent, a node is 0 in one cycle and 1 in the next with probability
$$\alpha = p_0 p_1 = (1 - p)p.$$
This peaks at 0.25 when $p = 0.5$. In the example, $p_1 = 0.6$ gives $\alpha = 0.4 \times 0.6 = 0.24$. If the NAND2 output drives 10 fF at 0.9 V with a 1 GHz clock, its power is $0.24 \times 10 \times 0.81 \times 1 = 1.94$ microwatts, since fF times GHz is $10^{-6}$. A clock has $\alpha = 1$ under this definition, far above any data signal, which is why clock power deserves its own attention.

## Where the estimate goes wrong

Both independence assumptions fail in real logic. Reconvergent fanout makes gate inputs correlated: in $A \cdot \bar A$, independence would predict probability 0.25 instead of 0. Temporal correlation, such as a counter's high bit that changes rarely, makes $(1 - p)p$ wrong even when $p$ is right. And zero-delay probability ignores glitches: unequal path delays make outputs toggle several times per cycle before settling, which can add substantially to power in arithmetic logic. Tools handle these in layers: propagate correlations approximately, or annotate activity from RTL or gate-level simulation (for example through SAIF or VCD files), and use timing-aware simulation to capture glitch power. Knowing which estimate a report used tells you how far to trust it.

## Going deeper

Najm, "A survey of power estimation techniques in VLSI circuits" (IEEE Transactions on VLSI Systems, 1994); Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), section 6.2.1.
