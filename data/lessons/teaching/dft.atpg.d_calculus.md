# D-calculus and five-valued logic

Automatic test pattern generation needs to reason about two circuits at once, the fault-free one and the faulty one, while most lines are still unassigned. Roth's D-calculus packs both machines into one value per line. Every structural ATPG algorithm, from the D-algorithm to PODEM and FAN, computes with these five values, so evaluating gates in D-calculus is the arithmetic of test generation.

## The five values

Each line carries a pair (fault-free value, faulty value). The pairs (0, 0) and (1, 1) are written 0 and 1: both machines agree. The pair (1, 0) is written **D**: the good circuit has 1 and the faulty circuit 0. The pair (0, 1) is written $\overline{D}$. A line not yet assigned is **X**. A fault is **activated** when its site carries D or $\overline{D}$: a line stuck-at-0 whose fault-free value is 1 carries D, a line stuck-at-1 whose fault-free value is 0 carries $\overline{D}$. If the fault-free value equals the stuck value, the site carries a plain 0 or 1 and the fault is not activated.

## Evaluating gates

To evaluate a gate, evaluate the good components and the faulty components separately with ordinary three-valued logic, then recombine. AND(D, 1) has good inputs (1, 1) and faulty inputs (0, 1), giving (1, 0) = D: the error passes because 1 is non-controlling for AND. AND(D, 0) gives (0, 0) = 0: the controlling side input blocks it. AND(D, $\overline{D}$) gives good AND(1, 0) = 0 and faulty AND(0, 1) = 0, so 0: the two errors cancel, which is how reconvergent fanout can mask a fault. NAND and NOR invert, so NAND(D, 1) = $\overline{D}$. XOR propagates through any side value: XOR(D, 1) = $\overline{D}$. When a component is unknown, AND(D, X) has good value AND(1, X) = X, so the result is recorded as X.

## Multi-gate evaluation

Values compose gate by gate in topological order. With $g_1 = $ OR(a, b), $g_2 = $ NOR($g_1$, c), a = D, b = 0 and c = 0: $g_1 = $ OR(D, 0) = D, and $g_2 = $ NOR(D, 0) = $\overline{D}$, so the error reaches $g_2$ inverted. If b had been X instead, $g_1$ would have good value OR(1, X) = 1 but faulty value OR(0, X) = X; the five-valued system collapses such partly known pairs to X, which loses some information. Nine-valued and other richer algebras keep the two components separately to avoid that loss, at the cost of more cases per gate.

## Going deeper

Roth, "Diagnosis of automata failures: a calculus and a method," IBM Journal of Research and Development, 1966; Bushnell and Agrawal (2000), chapter 7.
