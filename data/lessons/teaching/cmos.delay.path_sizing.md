# Sizing a path at equal stage effort

Knowing the optimal stage effort is not yet a design: each gate needs a size. Working backward from the load, each gate's input capacitance follows from its logical effort and the stage effort. This back-propagation is the core of analytical gate sizing and a useful mental model for what continuous sizers in synthesis and physical optimization compute.

## The optimal stage effort for the path

Start from the path: gates with logical efforts $g_1, \dots, g_N$, first input $C_{in}$, final load $C_{out}$ and, for simplicity, no branching. Compute $G = \prod g_i$, $H = C_{out}/C_{in}$, $F = GH$ and the stage effort $\hat f = F^{1/N}$. Example: an inverter, NAND2 and NOR2 with $\gamma = 2$, a 4 fF input and a 120 fF load give $G = 2.22$, $H = 30$, $F = 66.7$ and $\hat f = 4.05$. Every stage will bear this effort when sized optimally, so each gate's electrical effort is fixed at $h_i = \hat f / g_i$: the inverter's $h$ is 4.05, the NAND2's 3.04 and the NOR2's 2.43. Gates with larger logical effort get smaller electrical effort.

## Back-propagate input capacitances

Since $\hat f = g_i C_{out,i}/C_{in,i}$, each gate's input capacitance is
$$C_{in,i} = \frac{g_i\, C_{out,i}}{\hat f}.$$
Start at the load and work toward the input. In the example the NOR2 drives 120 fF, so its input is $\tfrac{5}{3} \cdot 120/4.05 = 49.3$ fF. That is the NAND2's load, so the NAND2's input is $\tfrac{4}{3} \cdot 49.3/4.05 = 16.2$ fF. That is the inverter's load, so its input is $1 \cdot 16.2/4.05 = 4.0$ fF, recovering the given first input: a built-in check on the arithmetic. From each input capacitance the transistor widths follow, using the gate's equal-drive ratios.

## From continuous sizes to a real netlist

The back-propagated sizes are continuous, while a library offers a handful of drive strengths such as X1, X2, X4 and X8. A sizer snaps to the nearest available size, then re-evaluates, since snapping one gate changes its driver's load. Branching changes the recursion: a stage's output load includes off-path capacitance, so $C_{in,i} = g_i(C_{on} + C_{off})/\hat f$ where the branching effort is part of $F$. Wire capacitance adds a fixed load that does not scale with gate size, which breaks the clean product structure and is one reason production sizers use numerical optimization, for example Lagrangian relaxation over a timing graph, rather than closed forms. The closed form remains the right first guess and the right sanity check.

## Going deeper

Sutherland, Sproull and Harris, "Logical Effort: Designing Fast CMOS Circuits" (1999), chapters 1 and 2; Chen, Chu and Wong, "Fast and exact simultaneous gate and wire sizing by Lagrangian relaxation" (IEEE TCAD, 1999).
