# Monotone functions

A monotone function is an order-preserving map from $B^n$ to $\{0, 1\}$: raising any input never lowers the output. Equivalently, it is positive unate in every variable, and it can be written with AND and OR only, no inversions. Monotone functions are the ones a tool can reason about with the order alone: their prime implicants are read straight off the minimal true points, their minimal CNF off the maximal false points, and many hard problems, such as SAT, become trivial for them.

## Monotone means order-preserving

$f$ is monotone when $x \le y$ implies $f(x) \le f(y)$. Its ON-set is then an up-set: if a vector is ON, every vector above it is ON. So the ON-set is determined by its minimal elements, the minimal true points, and is the union of the up-sets they generate.

For $f(a, b, c) = ab + c$ the minimal true points are 110 and 001. The up-set of 001 is every vector with $c = 1$, four rows; the up-set of 110 adds 110 itself (111 is already counted). So $f$ has 5 ON rows. The number of monotone functions of $n$ variables is the Dedekind number: 2, 3, 6, 20, 168 and 7581 for $n = 0$ to 5, with no closed form known.

## Minimal true points are the primes

For a monotone function every prime implicant is a product of positive literals, and the primes correspond one to one with the minimal true points: the point with ones exactly at positions $S$ gives the prime $\prod_{i \in S} x_i$. The minimal true points form an antichain, a set of pairwise incomparable vectors.

So the minimum SOP of a monotone function is just its set of primes, and every prime is essential: the minimal true point of a prime is covered by no other prime. In $f = ab + c$, the points 110 and 001 give the primes $ab$ and $c$. Two-level minimization, NP-hard in general, is easy here, which is one reason unate and monotone special cases are detected early.

## Maximal false points give the clauses

Dually, the OFF-set of a monotone function is a down-set, determined by its maximal false points. Each maximal false point with zeros exactly at positions $Z$ gives the clause $\bigvee_{i \in Z} x_i$, and the conjunction of these clauses is the minimal CNF.

For $f = ab + c$ the OFF rows are 000, 010 and 100. 000 lies below 010, so it is not maximal; 010 and 100 are maximal, because every vector above them is ON. They give the clauses $(a + c)$ and $(b + c)$, and indeed $(a + c)(b + c) = ab + c$. Converting between the two antichains is monotone dualization, a problem solvable in quasi-polynomial time (Fredman and Khachiyan) and not known to be polynomial.

## Going deeper

Crama and Hammer, "Boolean Functions: Theory, Algorithms, and Applications" (2011), chapters on monotone functions and dualization. Fredman and Khachiyan, "On the complexity of dualization of monotone disjunctive normal forms", Journal of Algorithms (1996).
