# Symbolic reachability by image computation

Symbolic reachability computes the set of reachable states without enumerating them, by repeatedly applying the transition relation to a set of states represented as a BDD. It was the breakthrough that let model checkers handle state spaces far beyond explicit search. Its iteration count, frontier sizes and backward variant are the basic cost quantities a model-checker designer reasons about.

## Image computation

Start with $R_0 = I$ and iterate $R_{i+1} = R_i \cup \text{Img}(R_i)$, where $\text{Img}(S)$ is the set of successors of states in $S$. The sequence grows until $R_{i+1} = R_i$, the fixpoint, which is the set of all reachable states. $R_i$ is exactly the set of states within $i$ steps of an initial state, so the number of image steps that add at least one new state equals the largest breadth-first distance from $I$ to a reachable state. With s0 → s1, s1 → s2, s2 → s0 and s2 → s3, starting at s0: $R_1$ adds s1, $R_2$ adds s2, $R_3$ adds s3, and $R_4 = R_3$. Three steps add states, and the fourth confirms the fixpoint. Long, thin state spaces such as counters need many iterations, which is a known weakness of breadth-first symbolic reachability.

## Counting the sets R_i

Because $R_i$ is the set of states at breadth-first distance at most $i$, its size can be read off a breadth-first search: count the states whose distance is at most $i$. With transitions s0 → s1, s2; s1 → s3; s2 → s3; s3 → s4; s4 → s0, starting at s0: distance 0 is s0, distance 1 is s1 and s2, distance 2 is s3, distance 3 is s4. So $|R_1| = 3$, $|R_2| = 4$ and $|R_3| = 5$. Symbolic tools often also keep the frontier $F_i = R_i \setminus R_{i-1}$, the states first reached at step $i$, and compute the next image from the frontier only, since images of older states add nothing new. The BDD of $R_i$ can be small even when $|R_i|$ is huge, and BDD size, not state count, is what limits the computation.

## Backward reachability and the BDD view

With BDDs over current-state variables $x$ and next-state variables $x'$, the image is the relational product $\text{Img}(S)(x') = \exists x.\, T(x, x') \wedge S(x)$, followed by renaming $x'$ to $x$. Packages compute the conjunction and the quantification in one pass (and-exists) and partition $T$ into conjuncts to keep intermediate BDDs small. The preimage $\text{Pre}(S)(x) = \exists x'.\, T(x, x') \wedge S(x')$ runs the same machinery backwards. Backward reachability starts from the bad states, $B_0 = \text{Bad}$ and $B_{i+1} = B_i \cup \text{Pre}(B_i)$, and the property fails exactly when the fixpoint contains an initial state. With s0 → s1, s1 → s2, s2 → s2, s3 → s2 and bad state s2, the backward fixpoint is $\{s2, s1, s3, s0\}$: 4 states, including s0, so the property fails. Counting it is breadth-first search on the reversed graph.

## Going deeper

Burch, Clarke, McMillan, Dill and Hwang, "Symbolic Model Checking: 10^20 States and Beyond" (LICS 1990); McMillan, "Symbolic Model Checking" (1993).
