# Fixpoints on the subset lattice

Many verification questions are answered by iterating a set operator until it stops changing. Reachable states, states from which a property can fail, live variables in a compiler, all are fixpoints of monotone operators on the lattice of subsets of a finite set. The Knaster-Tarski theorem guarantees that least and greatest fixpoints exist, and on a finite lattice plain iteration finds them. Symbolic model checkers run these iterations on BDDs or with SAT, and their cost is the number of iterations times the cost of one image.

## Iterating a monotone operator

An operator $\tau$ on subsets is monotone when $X \subseteq Y$ implies $\tau(X) \subseteq \tau(Y)$. Starting from the bottom and iterating, $\emptyset \subseteq \tau(\emptyset) \subseteq \tau^2(\emptyset) \subseteq \dots$, the chain must stop on a finite lattice, and where it stops is the least fixpoint $\mu Z.\, \tau(Z)$.

Reachability uses $\tau(Z) = I \cup \mathrm{Img}(Z)$, or equivalently $R_0 = I$ and $R_{k+1} = R_k \cup \mathrm{Img}(R_k)$. With transitions $s_0 \to s_1$, $s_1 \to s_2$, $s_2 \to s_0$, $s_3 \to s_2$ and $I = \{s_0\}$: $R_1 = \{s_0, s_1\}$, $R_2 = \{s_0, s_1, s_2\}$, and $R_3 = R_2$. Two steps added states; the third confirmed convergence. The number of productive steps equals the largest breadth-first distance from $I$, so deep counters make reachability slow.

## The least fixpoint is the reachable set

The set the iteration stops at is the smallest set that contains $I$ and is closed under successors, which is precisely the reachable set: here $\{s_0, s_1, s_2\}$, 3 states, with $s_3$ unreachable. A safety property holds exactly when this set avoids the bad states.

Symbolic implementations represent $R_k$ as a BDD and compute $\mathrm{Img}(R) = \exists x, i.\, R(x) \wedge T(x, i, x')$, then rename $x'$ to $x$. Only the frontier $R_k \setminus R_{k-1}$ needs imaging, which saves work. Bounded model checking replaces the fixpoint by unrolling $k$ steps into one SAT instance, and IC3 replaces it by a sequence of over-approximations that converge from above.

## Greatest fixpoints

Greatest fixpoints start from the top and remove elements. The states from which some infinite path stays inside a set $P$ forever are $\nu Z.\, P \cap \mathrm{Pre}(Z)$, where $\mathrm{Pre}(Z)$ is the set of states with a successor in $Z$. Start with $Z_0 = P$ and repeatedly delete states with no successor inside $Z$.

Add a self-loop $s_3 \to s_3$ to the system above and take $P = \{s_1, s_2, s_3\}$. In $Z_0$, $s_2$ has its only successor $s_0$ outside, so it goes. Then $s_1$, whose successor $s_2$ is gone, goes too. $s_3$ keeps its self-loop. The fixpoint is $\{s_3\}$, the CTL set $EG\, P$. Liveness and fairness checks nest least and greatest fixpoints, which is why their cost can be quadratic in the number of iterations.

## Going deeper

Clarke, Grumberg, Kroening, Peled and Veith, "Model Checking" (2nd ed., 2018), chapters on fixpoint characterizations and symbolic algorithms. Tarski, "A lattice-theoretical fixpoint theorem and its applications", Pacific Journal of Mathematics (1955).
