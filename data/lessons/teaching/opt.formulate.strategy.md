# Relax, decompose, approximate or search

Faced with an NP-hard model, a tool designer has four moves: relax a hard constraint to get a bound and a starting point, decompose the problem into smaller coupled pieces, approximate by rounding or a simpler model, or search the discrete space directly. Every major EDA engine is a recognizable combination of these moves, and naming the combination tells you what the engine can guarantee, where it will fail, and which knob to turn when it does.

## Relaxation and its bound

A relaxation enlarges the feasible set so the problem becomes easy. Replacing $x \in \{0, 1\}$ by $0 \le x \le 1$ turns an ILP into an LP. Dropping overlap and site constraints turns placement into a quadratic or smooth program. Because the relaxed set contains every feasible point, the relaxed optimum bounds the original: a lower bound for minimization.

Relaxation is rarely the end. Quadratic placement relaxes overlap and then spreads and legalizes. LP-based global routing relaxes integrality and then rounds the fractional route choices. Lagrangian relaxation is a different flavour: it moves the coupling constraints, such as timing, into the objective with multipliers, so the relaxed problem separates into per-gate subproblems. In every case a repair or coordination step must restore feasibility, and the bound tells you how much quality might remain.

## Decomposition and coordination

Decomposition splits one large problem into smaller ones. Min-cut placement recursively bisects the netlist and the die, so each region becomes an independent smaller problem. Hierarchical flows optimize blocks separately against budgets and then the top level against block abstracts. The risk is coupling: cross-block timing or congestion can be missed unless budgets are good or the levels iterate.

Formal decompositions coordinate the pieces through duality. Benders decomposition keeps the complicating integer decisions in a master problem and returns cuts from an LP subproblem's dual. ADMM splits an objective across partitions that share boundary variables and drives the copies to agree with scaled dual variables. Lagrangian relaxation, seen from this side, is price-based coordination: the multipliers are prices on shared resources. Decomposition is exact only when the coupling is handled exactly, so measure what the split costs on reference designs.

## Search, approximation and certificates

Search explores the discrete space directly with moves and an acceptance rule: simulated annealing on Polish expressions for floorplans, Fiduccia-Mattheyses passes for partitioning, local swaps in detailed placement. Search scales and adapts to odd constraints but carries no certificate. Approximation rounds or simplifies with a guarantee or at least a known error: randomized rounding of an LP, a greedy cover, an Elmore model in place of a full timer.

Exact methods survive where structure allows. Van Ginneken's buffer insertion on a routing tree and tree-covering technology mapping are dynamic programs that are optimal for their models. A strong practical pairing is search for the solution plus a relaxation for the bound: any feasible point gives an upper bound for minimization, the relaxation gives a lower bound, and the gap between them is the certificate. When a requirement such as $10^6$-cell detailed placement in minutes rules out exact methods, that pairing still tells the user what is being left on the table.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapters 2 to 5 for partition, floorplan, placement and routing strategies; Bertsekas, "Nonlinear Programming" (3rd ed., 2016), chapter on duality and decomposition.
