# Pareto optimality

Chip design never optimizes one number. Power, performance and area (PPA) pull against each other, as do wirelength and congestion, or yield and speed. Pareto optimality makes the trade-off precise: a design is Pareto-optimal if no other design is at least as good in every objective and strictly better in one. A flow sweep or design-space exploration should report the Pareto front, the set of such designs, so the designer chooses the trade-off instead of the tool hiding it.

## Dominance between two designs

With every objective minimized, design A dominates design B when A is no worse than B in every objective and strictly better in at least one. Equal values in all objectives do not dominate, and dominance is a partial order: many pairs are incomparable, each better in a different objective.

From designs D1 (delay 3, power 4), D2 (6, 2), D3 (12, 7), D4 (4, 4), D5 (6, 2), D6 (9, 8) and D7 (1, 5): which designs dominate D4 (4, 4)? D1 has delay 3 and power 4, no worse in both and better in delay, so it dominates D4. D2 and D5 have lower power but higher delay; D7 has lower delay but higher power. Only D1 dominates D4. D5 is dominated by nobody: D2 ties it exactly, which is not domination.

## Finding the Pareto front

A design is Pareto-optimal if no other design dominates it. Check each design against every other, $O(n^2)$ comparisons for n designs; for two objectives, sorting by the first objective and sweeping while tracking the best second objective finds the front in $O(n \log n)$.

Example: D1 (3, 1), D2 (6, 5), D3 (9, 2), D4 (10, 4), D5 (3, 10), D6 (10, 1), D7 (7, 11), D8 (6, 7). Sorted by delay, the best power so far is 1 after D1 (with D5 at the same delay but power 10, which D1 dominates). Every later design has power above 1, so all are dominated by D1, including D6 (10, 1), which ties D1's power with worse delay. The front is D1 alone: a design that is best in both objectives dominates everything. Real sweeps rarely have such a design, and their fronts show the trade-off curve.

## More objectives, larger fronts

With three objectives, such as power, delay and area, domination requires being no worse in all three, which is harder, so a larger fraction of designs is nondominated. Example: D2 (power 6, delay 4, area 7), D3 (3, 10, 5), D5 (10, 1, 12), D6 (1, 12, 8), D7 (1, 11, 5) and two worse designs. D7 dominates D6 (same power, better delay and area) but no design dominates D2, D3, D5 or D7, so the front has 4 designs.

As the number of objectives grows, almost every design becomes nondominated, and the front stops being a useful filter. That is the practical argument for keeping the objective count small, for turning some objectives into constraints (area at most a budget), and for preference information such as weights or reference points. Multi-objective search methods such as NSGA-II rank by dominance level and spread designs along the front with a crowding measure.

## Going deeper

Deb, "Multi-Objective Optimization Using Evolutionary Algorithms" (2001), chapters 2 and 6; Miettinen, "Nonlinear Multiobjective Optimization" (1999).
