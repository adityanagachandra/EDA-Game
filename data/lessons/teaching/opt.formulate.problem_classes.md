# Continuous versus combinatorial optimization

Once the variables are fixed, the problem falls into a class: linear, quadratic, convex, integer, mixed-integer, nonconvex, combinatorial or black-box. The class decides which algorithms apply, what guarantees exist and how the run time scales. EDA uses every class, often in sequence: a continuous relaxation for global structure, then a combinatorial repair. Recognizing the class of a subproblem is how a tool designer chooses between a sparse linear solve, an ILP solver, a dynamic program and a local search.

## Counting a discrete decision space

A combinatorial problem has a finite but enormous set of candidate solutions, and its size is the first thing to compute. If each of $n$ gates picks one of $k$ sizes independently, there are $k^n$ sizings: 10 gates with 3 sizes give $3^{10} = 59049$. Assigning $n$ cells to $n$ legal slots is a permutation, $n!$ of them: $8! = 40320$. Splitting $n$ cells into two unlabeled halves of $n/2$ gives $\binom{n}{n/2}/2$ bisections, because swapping the halves gives the same split: 8 cells give $70/2 = 35$.

The classic slips are confusing $k^n$ with $n^k$, forgetting that halves are unlabeled, or counting subsets ($2^n$) when the sizes are fixed. These counts grow exponentially or factorially, so a few dozen variables already put exhaustive search out of reach. The space is also unstructured: there is no gradient between two permutations.

## Naming the problem class

Read the class off the variables, objective and constraints. Continuous variables with a linear objective and linear constraints give a linear program (LP). Add a convex quadratic objective and it is a quadratic program (QP): quadratic placement, $\min \sum w_{ij}(x_i - x_j)^2$ with some cells fixed, is an unconstrained convex QP solved by one sparse linear system. Binary selection with linear covering rows, such as choosing test patterns to cover all faults, is a 0-1 integer linear program (ILP). Mixing continuous coordinates with binary relative-position variables, as in big-M floorplanning, is a mixed-integer linear program (MILP).

Continuous sizing with posynomial Elmore delays is a geometric program, convex after the substitution $x = e^u$. A constraint $(x_i - x_j)^2 + (y_i - y_j)^2 \ge d^2$ describes the outside of a disk, so it makes the model a nonconvex quadratically constrained program (QCQP).

## Why EDA alternates continuous and discrete

Legal placement is discrete (cells on rows and sites), but global placers solve a continuous relaxation first because gradients and sparse solvers scale to millions of cells. Legalization then repairs the small discrete violations. Continuous sizing returns values like 2.37x for a library of 1x, 2x, 4x and 8x, so a discretization step snaps sizes and repairs the timing that rounding breaks.

Some problems stay combinatorial: detailed routing chooses paths on a discrete track grid, and technology mapping chooses matches that cover a subject graph, solved exactly on trees by dynamic programming. Others are black boxes: tuning ten flow knobs where each evaluation is a six-hour run with no gradient calls for surrogate-based search. The structural difference for the designer is simple. Continuous problems offer gradients and local models; discrete ones need enumeration, bounds, relaxations or neighborhoods.

## Going deeper

Papadimitriou and Steiglitz, "Combinatorial Optimization: Algorithms and Complexity" (1982), chapter 1; Boyd and Vandenberghe, "Convex Optimization" (2004), chapter 4 for the LP, QP, QCQP and geometric-program forms.
