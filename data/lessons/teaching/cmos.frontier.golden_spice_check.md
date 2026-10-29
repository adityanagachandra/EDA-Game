# Golden SPICE checks for library tables

A golden SPICE check is the test that decides whether a characterization change is allowed to ship. It compares a candidate's Liberty tables with tables produced by a trusted reference flow, entry by entry, against a tolerance. Getting the error metric and the tolerance right decides whether the check catches real regressions without drowning the harness in false alarms.

## Relative error at a grid point

The basic metric is the relative error of a candidate value $c$ against the golden value $g$:
$$e = \frac{|c - g|}{g} \times 100\%.$$
A candidate delay of 52 ps where the golden SPICE value is 50 ps has $e = 4\%$. Relative error is natural for delays that span orders of magnitude across a table, from a few picoseconds at small loads to hundreds at large ones. The check runs on every arc, edge, table type and grid point of the golden cells, and for transitions, constraints, capacitances and power as well as delay, because a characterizer can be right on delay and wrong on slew, and slew propagates.

## Maximum error over a table

A table passes or fails on its worst entry, so the summary statistic is the maximum relative error over the grid, not the mean. A mean can hide one bad corner of the table, and that corner, often the smallest load or the fastest slew, may be exactly where real paths operate. Small golden values amplify fixed absolute errors: 1 ps of error on a 10 ps entry is 10%, while on a 200 ps entry it is 0.5%. A report that lists the maximum and where it occurs, by row and column, tells the engineer or the agent whether the problem is a measurement threshold, an interpolation region or a convergence failure in a few simulations.

## Tolerances with an absolute floor

A pure relative tolerance fails tiny entries for errors no timer would notice, while a pure absolute tolerance is too loose for large ones. Practical checks combine them: an entry passes if
$$|c - g| \le \max(r \cdot g,\ a),$$
for a relative tolerance $r$ and an absolute floor $a$. With $r = 2\%$ and $a = 1$ ps, a 12 ps entry may be off by 1 ps and a 150 ps entry by 3 ps. The check is a gate in the skill flywheel: nothing is promoted without it, and its own definition must be versioned, since loosening a tolerance is itself a change that can hide regressions. It complements a differential check against the previous library, which catches changes in cells the goldens do not cover.

## Going deeper

The agent-first EDA material in this curriculum (golden and differential testing); the Liberty Reference Manual.
