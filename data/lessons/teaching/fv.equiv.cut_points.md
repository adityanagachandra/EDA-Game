# Cut points and false negatives

Proving equivalence of large cones in one SAT call can be too hard, so checkers cut the problem at internal points already proved equivalent. Cut points make each query small, but they abstract away the logic in front of them, and the abstraction can produce counterexamples that the real circuits never exhibit. Managing these false negatives is a small, complete instance of abstraction refinement that every equivalence engine must get right.

## Cutting at proved equivalences

Suppose sweeping has proved that node $u$ in design $A$ equals node $u'$ in design $B$. A cut point replaces both by one fresh free variable $w$ in all later checks. Logic downstream of the cut now reads $w$ on both sides, and the solver no longer reasons about the fan-in cones of $u$ and $u'$. Each output miter then only contains the logic between the nearest cuts and the output, so a large cone splits into a chain of small problems. The approach is sound for proofs: $w$ can take every value $u$ can take and more, so if the outputs agree for every value of $w$, they agree for the values the real logic produces. Sweeping engines place cuts at every proved-equivalent pair, which is why proving bottom-up keeps each SAT call small.

## False negatives

The free variable forgets correlations. If two cut variables $w_1$ and $w_2$ stand for signals that in the real circuit always satisfy $w_1 = \neg w_2$, the solver may set $w_1 = w_2$, a combination that never occurs, and find an output difference under it. That counterexample is a false negative: a reported difference between designs that are actually equivalent. False negatives are most likely when cuts are placed at correlated signals, such as a signal and a re-encoded copy of it, or two outputs of the same small function. They never make a wrong proof, only a wrong alarm, but a checker that reports them as bugs is useless in practice, so every cut-based engine must detect them.

## Resolving false negatives

When a miter with cuts is satisfiable, the checker takes the counterexample's primary-input values and simulates the full designs without cuts. If the outputs really differ, the bug is real and is reported. If they agree, the counterexample relied on impossible cut values, and the checker refines: it removes the offending cut, or moves it back toward the inputs, so that the correlated logic is visible again, and re-solves. This is the CEGAR pattern applied to combinational logic: over-approximate, check, test the counterexample on the concrete system, and refine only when it is spurious. The refinement may make a query larger, so engines try cheap fixes first, such as removing only the cuts that the counterexample's difference actually depends on.

## Going deeper

Kuehlmann and Krohm, "Equivalence Checking Using Cuts and Heaps" (DAC 1997); Kuehlmann et al., "Robust Boolean Reasoning for Equivalence Checking and Functional Property Verification" (IEEE TCAD, 2002).
