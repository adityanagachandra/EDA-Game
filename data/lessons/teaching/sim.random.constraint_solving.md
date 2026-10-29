# Constraint solvers for stimulus

Behind every call to `randomize()` is a constraint solver that must find a legal assignment quickly and, ideally, sample uniformly over all legal assignments. These two goals pull in different directions. BDD-based solvers sample exactly but can blow up in memory; SAT and SMT solvers find solutions fast but are not naturally uniform. The trade-offs shape what constraint writers can express efficiently and are an active research area in stimulus generation.

## BDD-based and SAT-based solvers

A **BDD-based** solver builds a binary decision diagram of the conjunction of the constraints over the bits of the random variables. Counting the satisfying paths below each node takes one pass; sampling then walks from the root, taking each branch with probability proportional to its count, which yields an exactly uniform solution in time linear in the number of bits. The weakness is size: for multiplication and some other arithmetic, every variable order gives an exponentially large BDD. A **SAT** or **SMT** solver instead bit-blasts the constraints into clauses (or reasons at word level) and searches for one satisfying assignment. Search is fast in practice even on large problems, but the solution it returns is biased by its heuristics. Near-uniform sampling needs extra machinery, for example adding random XOR (parity) constraints that partition the solution space into small cells and sampling within a random cell.

## The cost of rejection sampling

The naive alternative is rejection: draw every random variable uniformly and retry until the constraints hold. If the space has $T$ assignments and $L$ are legal, each draw succeeds with probability $L/T$, and the number of draws is geometric with mean $T/L$. Example: two 4-bit variables with $x + y < 4$: $T = 256$ and $L = 10$, so on average 25.6 draws per accepted solution. Rejection is exactly uniform and trivial to implement, but the ratio $T/L$ explodes as constraints tighten: equality constraints on 32-bit fields make $L/T$ about $2^{-32}$. That is why production solvers propagate constraints instead of guessing and checking, and why rejection survives only for loosely constrained fields.

## Performance pitfalls for constraint writers

Solvers partition random variables into groups linked by constraints and solve each group separately, so unnecessary links (one constraint touching many unrelated fields) create one large hard problem. Nonlinear arithmetic, wide multiplications, and modulo on wide values are expensive for every solver. Function calls inside constraints impose an order: arguments are solved first and the function's result is then a constant, which can make otherwise satisfiable sets fail. Non-random state variables are treated as constants, which shrinks the problem. `randc` variables cycle through a random permutation of their range before repeating, which requires the solver to remember history. When `randomize()` becomes slow, profiling which constraint group dominates usually shows a single badly written constraint.

## Going deeper

Yuan, Pixley and Aziz, "Constraint-Based Verification" (2006); Chakraborty, Meel and Vardi, "A scalable and nearly uniform generator of SAT witnesses", CAV (2013).
