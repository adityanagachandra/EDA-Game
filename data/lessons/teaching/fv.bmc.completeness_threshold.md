# Diameters and completeness thresholds

BMC proves nothing beyond its bound unless the bound is known to be large enough. A completeness threshold is a depth beyond which no new counterexample can appear, so BMC up to it is a proof. For simple safety properties the threshold is a diameter of the state graph. Computing diameters by hand shows when this route to a proof is practical and why it usually is not for real designs.

## The reachability diameter

For an invariant (\"the bad state is never reached\"), every reachable state has a shortest path from an initial state, and the longest of those shortest paths is the reachability diameter from the initial states. Any reachable bad state is reached within that many steps, so BMC that checks the bad signal in every frame up to the diameter is complete. To compute it, run breadth-first search from the initial states and take the largest distance. With s0 → s1, s2; s1 → s3; s2 → s1; s3 → s0, starting at s0: s1 and s2 are at distance 1 and s3 at distance 2, so the diameter is 2, and BMC to depth 2 proves any invariant of this system. The catch is that computing the diameter of a real design is as hard as reachability itself.

## The recurrence diameter

The recurrence diameter from the initial states is the length of the longest loop-free path from an initial state. Any path longer than it must revisit a state, and a path that revisits a state can be shortened by cutting out the loop, so it also serves as a completeness threshold, and it is at least the reachability diameter. Its advantage is that it can be checked with SAT: BMC can ask whether a path of length $k$ with pairwise distinct states exists, and once none does, the bound is complete. In the example above, the path s0 → s2 → s1 → s3 has 3 transitions and visits no state twice, while every longer path repeats one, so the recurrence diameter is 3, against a reachability diameter of 2. On some graphs the gap is exponential, which limits this approach.

## Diameters of counters

Counters show why completeness thresholds rarely help in hardware. A single $n$-bit counter that resets to 0 and increments every cycle visits its states in one cycle of length $2^n$, so its farthest state, $2^n - 1$, is $2^n - 1$ steps from reset: the diameter is $2^n - 1$, over four billion for 32 bits. Two counters modulo $a$ and $b$ that step together from $(0, 0)$ return to $(0, 0)$ after $\text{lcm}(a, b)$ cycles and visit each of those pairs once, so the diameter is $\text{lcm}(a, b) - 1$; for $a = 4$ and $b = 6$ it is 11, and only 12 of the 24 pairs are reachable. Real designs contain timers, FIFOs and pointers, so their diameters are huge. That is why proofs use induction, k-induction or IC3, which do not depend on depth.

## Going deeper

Biere, Cimatti, Clarke and Zhu, "Symbolic Model Checking without BDDs" (TACAS 1999); Kroening and Strichman, "Efficient Computation of Recurrence Diameters" (VMCAI 2003).
