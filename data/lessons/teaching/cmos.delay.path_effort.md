# Path effort and minimum path delay

A path of gates has an optimum: there is a minimum delay, achieved when every stage bears the same effort. Logical effort turns this into closed-form arithmetic on path effort $F = GBH$. It is the theory behind buffer chains, the first estimate any sizer makes, and a clean example of how a product constraint leads to an equal-allocation optimum.

## Path logical, electrical and branching effort

For an $N$-stage path, the path logical effort is the product of the stages' logical efforts, $G = \prod g_i$, and the path electrical effort is $H = C_{out}/C_{in}$, the final load over the first gate's input. Where a stage's output also drives gates off the path, define its branching effort $b = (C_{on} + C_{off})/C_{on}$, and the path branching effort $B = \prod b_i$. The path effort is
$$F = G B H.$$
Example: an inverter, NAND2 and NOR2 path with $\gamma = 2$ has $G = 1 \cdot \tfrac{4}{3} \cdot \tfrac{5}{3} = 2.22$. With a 5 fF input, a 200 fF load and no branching, $H = 40$ and $F = 88.9$. If the first output also drives off-path gates of twice the on-path input capacitance, $b = 3$ and $F = 267$.

## Equal stage effort minimizes delay

The path delay is $D = \sum (g_i h_i) + P$ with $P = \sum p_i$. The product of the stage efforts $f_i = g_i h_i$ is fixed at $F$, because the intermediate capacitances cancel. Minimizing a sum of positive terms whose product is fixed gives equal terms (by the AM-GM inequality, or a Lagrange multiplier), so each stage should bear
$$\hat f = F^{1/N}, \qquad D_{min} = N F^{1/N} + P.$$
For the three-stage example with $F = 88.9$: $\hat f = 4.46$ and $P = 1 + 2 + 2 = 5$, so $D = 13.4 + 5 = 18.4$ units of $\tau$. The minimum is reached by sizing, not by changing gates, so it is the best this topology can do.

## What the optimum assumes

The result assumes a fixed topology, a fixed number of stages, a fixed input capacitance and a known load, and continuous sizes. Real libraries offer discrete drive strengths, wires add fixed capacitance that does not scale with gate size, and branching often varies across a cone. Optimizers therefore use the logical-effort optimum as a starting point or a lower bound, then snap to library sizes and verify with a timer. The equal-effort principle still guides many heuristics, for example in buffer insertion and fanout trees, and it explains a common observation: if one stage in a timing report has a much larger delay than its neighbours, the path is unbalanced and resizing that stage and its driver will usually help.

## Going deeper

Sutherland, Sproull and Harris, "Logical Effort: Designing Fast CMOS Circuits" (1999), chapters 1 to 3.
