# Fixpoint computation of CTL operators

Symbolic model checkers compute CTL operators as fixpoints of set transformers: start from a set of states, apply one step of EX and a union or intersection, and repeat until the set stops changing. The number of iterations governs run time, and the distinction between least and greatest fixpoints is exactly the distinction between eventually and forever. These iterations are also the template for reachability and many other symbolic algorithms.

## EF as a least fixpoint

EF $q$ is the least fixpoint of $Z \mapsto [[q]] \cup \text{EX}\, Z$. Compute it by iterating from the target states: $Z_0 = [[q]]$ and $Z_{i+1} = Z_i \cup \text{EX}\, Z_i$, which adds, at step $i + 1$, the states with a successor in $Z_i$. $Z_i$ is the set of states that can reach $q$ within $i$ steps. Stop at the smallest $i$ with $Z_{i+1} = Z_i$. With transitions s0 → s1, s1 → s2, s2 → s2, s3 → s0 and $q$ only at s2: $Z_0 = \{s2\}$, $Z_1 = \{s1, s2\}$, $Z_2 = \{s0, s1, s2\}$, $Z_3 = \{s0, s1, s2, s3\}$, and $Z_4 = Z_3$. The smallest $i$ with $Z_{i+1} = Z_i$ is 3. The count equals the largest backward distance to a $q$-state among states that can reach one. Starting from $[[q]]$ rather than from the empty set saves one iteration and is the convention used here; the fixpoint itself, the set of states satisfying EF $q$, is the same either way, and here it contains every state.

## EG as a greatest fixpoint

EG $p$ is the greatest fixpoint of $Z \mapsto [[p]] \cap \text{EX}\, Z$. Iterate downward: $Z_0 = [[p]]$ and $Z_{i+1} = Z_i \cap \text{EX}\, Z_i$, which removes states that have no successor inside the current set. What survives is the set of $p$-states with an infinite path inside the $p$-states. With s0 $\{p\}$ → s1 $\{p\}$ → s2 $\{\}$, s3 $\{p\}$ → s3, and s2 → s2: $Z_0 = \{s0, s1, s3\}$; s1's only successor s2 is outside, so $Z_1 = \{s0, s3\}$; now s0's only successor s1 is outside, so $Z_2 = \{s3\}$; and $Z_3 = Z_2$. The answer is $i = 2$, and EG $p$ holds only at s3. Starting from all states and intersecting is what makes it a greatest fixpoint; starting from the empty set would give the wrong answer, the empty set.

## Until and the lattice argument

E[$p$ U $q$] is the least fixpoint of $Z \mapsto [[q]] \cup ([[p]] \cap \text{EX}\, Z)$: iterate $Z_0 = [[q]]$ and $Z_{i+1} = Z_i \cup ([[p]] \cap \text{EX}\, Z_i)$, adding only $p$-states that can step into the set. Why must these iterations stop, and why at the right set? The subsets of $S$ form a finite lattice under inclusion, and each transformer is monotone: a larger $Z$ gives a larger result. By the Knaster-Tarski theorem a monotone function on a complete lattice has least and greatest fixpoints; iterating from the bottom (or from $[[q]]$, which lies below the fixpoint) climbs to the least one, and iterating down from $[[p]]$ reaches the greatest one. Each step adds or removes at least one state until it stops, so at most $|S|$ iterations are needed. In BDD-based checkers each iteration is a preimage computation, so the iteration count multiplies the cost.

## Going deeper

McMillan, "Symbolic Model Checking" (1993); Clarke, Grumberg, Kroening, Peled and Veith, "Model Checking" (2nd ed., 2018), the chapter on symbolic model checking.
