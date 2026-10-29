# NP-completeness

A problem is NP-complete when it is in NP and every problem in NP reduces to it in polynomial time. NP-complete problems are the hardest in NP: a polynomial algorithm for any one of them would give one for all, and none is known. Many core EDA problems, from SAT and test generation to graph coloring, partitioning and placement decisions, are NP-complete or NP-hard, and that fact is why EDA engines are built around heuristics, approximations, resource limits and exact solvers that are fast only on typical instances.

## The Cook-Levin theorem

Cook (1971) and Levin independently proved that SAT is NP-complete. The idea: any NP problem has a polynomial-time verifier that reads an instance and a certificate. For a fixed input size, that verifier's computation unrolls into a circuit of polynomial size whose free inputs are the certificate bits. The instance is a yes-instance exactly when some certificate makes the circuit accept, a circuit-satisfiability question.

Circuit-SAT then reduces to CNF-SAT by the Tseitin encoding, which adds one variable per gate and a few clauses per gate, linear in circuit size and satisfiable exactly when the circuit is. The canonical CNF would not do, since it can be exponential. Tseitin's formula is equisatisfiable, not equivalent: its models live over the inputs plus the gate variables.

## Proving a new problem NP-complete

Two steps are required. First, show the problem is in NP by naming a polynomial certificate and verifier. Second, reduce a known NP-complete problem to it in polynomial time. Karp's 1972 paper did this for 21 problems, including clique, vertex cover, graph coloring and set cover, and Garey and Johnson's catalog lists hundreds more.

Common mistakes are reducing the new problem to SAT (that proves only membership), using a reduction whose output is exponentially large (it proves nothing), or reasoning from benchmark runtimes (they describe instances, not the problem). Independent set is NP-complete because clique reduces to it by complementing the graph, a polynomial step, and a candidate set checks quickly.

## NP-hard, NP-complete and the problems a tool meets

NP-hard means at least as hard as every NP problem, with no requirement to be in NP. Optimization versions such as "find the minimum-wirelength placement" are NP-hard when their decision versions are NP-complete. Some problems sit higher still: deciding whether a formula has an equivalent formula with at most $k$ literals is complete for the second level of the polynomial hierarchy, because a guessed small formula must also be verified equivalent on every input.

Special cases are often easy. Graph 2-coloring is bipartiteness, solved by breadth-first search, while 3-coloring is NP-complete; 2-SAT is linear while 3-SAT is NP-complete; exact two-level minimization is hard because its covering step is set cover and the number of primes can be exponential. Locating the hard core of a tool's problem tells you where heuristics must go and where an exact algorithm is affordable.

## Going deeper

Cook, "The complexity of theorem-proving procedures", STOC (1971). Karp, "Reducibility among combinatorial problems" (1972). Garey and Johnson, "Computers and Intractability" (1979). Buchfuhrer and Umans, "The complexity of Boolean formula minimization", ICALP (2008).
