# Tractable classes of SAT

SAT is NP-complete, but several syntactic classes are solvable in polynomial time, and a surprising share of the constraints that EDA tools generate fall into them. Recognizing a tractable fragment lets a tool skip the general solver or hand it a much easier problem. Schaefer's dichotomy theorem says that for Boolean constraint languages this is all there is: each language is either in one of a few polynomial classes or NP-complete. Knowing where the line falls tells a designer when a special-purpose algorithm exists.

## 2-SAT and the implication graph

A CNF in which every clause has at most two literals is solved in linear time. Each clause $(p \lor q)$ is equivalent to the implications $\neg p \to q$ and $\neg q \to p$; draw them as edges on $2n$ literal vertices. The formula is unsatisfiable exactly when some variable $x$ and $\neg x$ lie in the same strongly connected component, since then $x$ implies $\neg x$ and $\neg x$ implies $x$.

The clause $(a \lor \neg b)$ gives the edges $\neg a \to \neg b$ and $b \to a$, so a 2-CNF with $m$ clauses has $2m$ edges. Finding components with Tarjan's algorithm and reading off an assignment in reverse topological order takes $O(n + m)$. Binary clauses are common in practice, from equivalences and implications, and modern solvers store them in a separate implication structure for speed.

## Horn clauses and forward chaining

A Horn clause has at most one positive literal, such as $(\neg a \lor \neg b \lor c)$, which reads as the rule $a \wedge b \to c$. Horn-SAT is solved by starting with every variable false and setting a variable true only when some clause forces it. If a clause with no positive literal ends up with all its variables true, the formula is unsatisfiable; otherwise the forced set is a model.

That model is the least one in the componentwise order, because Horn models are closed under intersection. With linear-time bookkeeping (Dowling and Gallier), unit propagation alone decides Horn formulas, which is why Horn fragments are solved with no search at all. Dual-Horn, at most one negative literal, is symmetric.

## XOR constraints and where hardness begins

Constraints of the form $x_1 \oplus x_2 \oplus \dots \oplus x_k = c$ are linear equations over GF(2) and are solved by Gaussian elimination in polynomial time. In CNF without auxiliary variables a single such XOR needs $2^{k-1}$ clauses, and CDCL handles long XOR chains badly, so solvers aimed at cryptographic and arithmetic instances add native Gaussian elimination.

Allowing three literals per OR clause breaks all of these: 3-SAT is NP-complete. Schaefer (1978) showed the tractable Boolean constraint languages are exactly the 0-valid, 1-valid, Horn, dual-Horn, bijunctive (2-CNF) and affine (XOR) ones; everything else is NP-complete. A useful consequence for tool builders: if a problem's constraints are all binary or all Horn, a polynomial algorithm exists and a general SAT call is overkill.

## Going deeper

Aspvall, Plass and Tarjan, "A linear-time algorithm for testing the truth of certain quantified Boolean formulas", Information Processing Letters (1979). Dowling and Gallier, "Linear-time algorithms for testing the satisfiability of propositional Horn formulae", Journal of Logic Programming (1984). Schaefer, "The complexity of satisfiability problems", STOC (1978).
