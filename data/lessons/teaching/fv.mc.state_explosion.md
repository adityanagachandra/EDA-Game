# State-space explosion

The number of states of a system is the product of the numbers of states of its parts, so it grows exponentially with the number of components or bits. This state-space explosion is the central obstacle of model checking, and every major technique in the field, from BDDs and SAT-based engines to abstraction, partial-order reduction and symmetry reduction, is a response to it. Computing the sizes involved makes the problem concrete.

## Products of components

When components run together, a global state records the local state of each one, so the global state space is the Cartesian product of the local ones. $n$ copies of a component with $k$ states give $k^n$ global states: 10 processes of a 5-state protocol already give $5^{10} \approx 9.8$ million states, and 20 give about $9.5 \cdot 10^{13}$. Not every product state need be reachable, but explicit-state model checkers must store each reachable one, and in concurrent systems a large fraction often is. The growth is in the number of components, not their size, which is why verifying a protocol for 2 or 3 nodes is easy and for 20 is not.

## Counters and FSMs together

Hardware adds bits. A free-running $r$-bit counter alone has $2^r$ states; combined with independent FSMs of sizes $a$, $b$ and $c$, the design has $a \cdot b \cdot c \cdot 2^r$ states. A block with FSMs of 4, 6 and 5 states and a 10-bit counter has $120 \cdot 1024 = 122{,}880$ states, and a 32-bit counter multiplies everything by over four billion. Real designs have thousands of state bits, so the state space is astronomically larger than any explicit search could visit. Symbolic methods work only because BDD or SAT representations do not grow with the number of states: a 32-bit counter's reachable set has a tiny BDD. They still fail when the reachable set has no compact representation, which is why abstraction, cone-of-influence reduction (dropping logic the property does not depend on) and bounded methods matter.

## Symmetry reduction

When processes are identical and interchangeable, global states that differ only by a permutation of processes behave the same, and a model checker can keep one representative per class. For $n$ identical processes with $k$ local states each, a class is determined by how many processes sit in each local state, a multiset of size $n$ over $k$ kinds. By stars and bars there are $\binom{n + k - 1}{k - 1}$ such classes. For $n = 6$ and $k = 4$ that is $\binom{9}{3} = 84$ representatives, against $4^6 = 4096$ states. Symmetry must actually hold, including in the property. Partial-order reduction attacks a different source of explosion: interleavings of independent concurrent actions, of which it explores only one representative order. Both are mainly used in explicit-state software and protocol checkers; hardware engines rely more on symbolic representations and abstraction.

## Going deeper

Clarke, Grumberg, Kroening, Peled and Veith, "Model Checking" (2nd ed., 2018), the chapters on symmetry and partial-order reduction; Baier and Katoen, "Principles of Model Checking" (2008), chapter 8.
