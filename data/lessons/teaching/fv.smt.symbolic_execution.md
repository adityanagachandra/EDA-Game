# Symbolic execution

Symbolic execution runs a program, or an RTL model, on symbolic inputs instead of concrete ones. Each branch splits the run into paths, each path carries a path condition over the inputs, and an SMT solver decides which paths are feasible and produces inputs that drive them. It underlies test generation, bug finding and some hardware analysis tools, and its central problem, path explosion, is a counting problem.

## Paths through a program

A program with $k$ independent if-statements in sequence has $2^k$ syntactic paths: each branch can go either way, independently of the others. With three such statements there are 8 paths, with 30 over a billion. An if / else-if chain with $k$ conditions and a final else has only $k + 1$ paths, because exactly one arm executes. Loops multiply paths by the number of iterations considered, which is why symbolic executors bound loop unrolling. Hardware has the same structure: one clock cycle of an always_comb block with $k$ independent conditionals has $2^k$ paths, and $c$ cycles of symbolic simulation multiply them again. Counting syntactic paths is the first estimate of the cost, before the solver prunes any.

## Path conditions and feasibility

Along each path the executor collects a path condition: the conjunction of the branch conditions taken, each negated when the else side is taken. At each branch it asks an SMT solver whether the extended path condition is satisfiable; unsatisfiable paths are infeasible and pruned. Take a 4-bit unsigned input $x$ and the statements if ($x < 8$) ... ; if ($x \ge 12$) ... . The path that takes both branches needs $x < 8 \wedge x \ge 12$, which is unsatisfiable, so only 3 of the 4 syntactic paths are feasible. For small input domains the feasible paths are exactly the distinct branch-outcome patterns that some input produces, which can be counted by enumerating the inputs. A satisfying assignment of a feasible path condition is a concrete test input for that path, which is how symbolic executors generate tests.

## Counting the inputs on a path

The number of inputs that follow a path is the model count of its path condition. In the example, the path that takes neither branch has $\neg(x < 8) \wedge \neg(x \ge 12)$, so $8 \le x \le 11$: 4 inputs. Intersect the set allowed by each conjunct; for bit tests such as $x[0] = 1$, keep the values with that bit set. Model counts measure how likely a path is under random inputs, which explains why random testing misses rare paths and symbolic execution finds them. In practice, path explosion is managed by merging paths at join points into one state with ite expressions, by concolic execution (running concretely and negating one branch condition at a time), and by search heuristics that prioritize uncovered code. Tools such as KLEE apply these ideas to software, and symbolic simulation applies them to RTL.

## Going deeper

King, "Symbolic Execution and Program Testing" (Communications of the ACM, 1976); Cadar, Dunbar and Engler, "KLEE: Unassisted and Automatic Generation of High-Coverage Tests for Complex Systems Programs" (OSDI 2008).
