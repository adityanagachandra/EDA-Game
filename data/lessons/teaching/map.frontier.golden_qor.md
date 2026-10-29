# Equivalence and QoR deltas as golden tests

When agents edit mapping code or cost functions, every change needs an automatic, trustworthy verdict. For mapping that verdict has two parts: the mapped netlist must be equivalent to its input, and its quality must not regress beyond agreed budgets on a fixed suite of designs. Getting the order, the baseline and the statistics right is what makes the harness safe to automate.

## Equivalence before quality

Run a combinational equivalence check (ABC's cec, or a SAT-based checker) between the pre-mapping network and the mapped netlist, on every design in the suite. Only a proved result passes. A counterexample means the change is wrong, whatever its QoR. A timeout is not a pass: hold the change, rerun with more effort or another method, and keep the gate closed. Then apply the QoR rule, for example: delay may regress by at most 1 percent and area must not grow unless delay improves. A candidate with area -6.8 percent and delay +3.1 percent fails that rule on delay.

## Computing and reporting deltas

Report each design's percentage change, (candidate - baseline) / baseline x 100, and summarize with the mean of the per-design percentages (or a geometric mean of ratios), so one large design does not dominate. For areas 1041 to 1061, 826 to 821 and 1614 to 1559, the changes are +1.92, -0.61 and -3.41 percent, mean -0.70 percent. Always show the per-design list next to the mean: a 2 percent mean gain can hide one design that got 9 percent worse. Measure delay with the reference timer on the mapped netlist, not the mapper's internal estimate.

## What makes a golden suite trustworthy

The baseline is the currently promoted recipe on the same designs, library and constraints; nothing else may differ. Runs must be deterministic and seeded, or noise looks like progress. Run time and memory are tracked metrics with their own budgets, since they are part of a tool's QoR. A targeted change should have a targeted effect: if only one cell class's cost changed, the differential report should show only matches involving that class moving. Designs are never removed from the suite to make a change pass; outliers are investigated, and a real regression on one design blocks promotion until it is understood.

## Going deeper

Brayton and Mishchenko, "ABC: an academic industrial-strength verification tool", CAV 2010, for cec; the EPFL combinational benchmark suite, a standard golden set for mapping QoR.
