# Combinational equivalence checking with miters

Every synthesis transformation must preserve function, and combinational equivalence checking proves it. The method builds a miter, a circuit that outputs 1 exactly when the two designs disagree, and proves that it can never output 1. Simulation, FRAIG-style sweeping and SAT make this scale to large designs, and it is the golden test for any new optimization.

## The miter

Given two circuits with matching inputs and outputs, connect corresponding inputs, XOR each pair of corresponding outputs, and OR the XORs together. The miter output is 1 on an input vector exactly when some output pair differs there. The designs are equivalent if and only if the miter is unsatisfiable. For $F = a(\bar{b} + \bar{c})$ and a rewrite $G = a(\bar{b} + \bar{b})$, the miter is 1 where they differ: with $a = 1$, $b = 0$ both are 1, with $b = 1$, $c = 0$ only $F$ is 1. Counting such vectors over the free inputs gives the number of counterexamples.

## Several outputs

For a multi-output circuit the miter ORs the per-output XORs, so a vector is a counterexample if any output differs. The counterexample count is the size of the union of the per-output difference sets, not their sum: two outputs that each differ on 4 vectors, 2 of them shared, give 6 counterexamples. In practice the checker does not count; it stops at the first satisfying assignment and reports it as a counterexample, which is often simulated to locate the bug. Checking outputs one at a time gives smaller SAT problems and points directly at the failing output.

## Making the proof scale

The miter of two cones with $A$ and $B$ AND nodes has about $A + B$ nodes plus three per output XOR and one per OR, so two outputs add 7. Solving it directly with SAT is hard when the two designs are structurally different. Real checkers first structurally hash the miter, which merges identical logic, then simulate to find candidate internal equivalences, prove them with SAT from the inputs upward (SAT sweeping), and merge them. Once internal points match, the output XORs often collapse to constant 0. ABC's cec command implements this. Retiming breaks register correspondence, so sequential transformations need sequential equivalence checking instead.

## Going deeper

Kuehlmann, Paruthi, Krohm and Ganai, "Robust Boolean reasoning for equivalence checking and functional property verification", IEEE TCAD (2002); Mishchenko, Chatterjee, Brayton and Een, "Improvements to combinational equivalence checking", ICCAD (2006).
