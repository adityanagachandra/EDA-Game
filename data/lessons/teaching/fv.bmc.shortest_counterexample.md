# Counterexample depth in BMC

Because BMC checks depths in increasing order, the first depth at which it finds a counterexample is the length of the shortest one. Short counterexamples are easy to debug, which is a large part of BMC's practical value. Predicting that depth, from a state graph or from the arithmetic of a counter, tells you how deep a run must go and how many solver calls a bounded check will make.

## The first satisfiable depth

BMC at depth $k$ asks whether some path of $k$ transitions from an initial state ends in a bad state (or, in the any-frame form, hits one along the way). Running $k = 0, 1, 2, \dots$, the first satisfiable depth equals the breadth-first distance from the initial states to the nearest bad state. With s0 → s1, s1 → s2, s1 → s3, s3 → s4 and bad state s4, starting from s0: s1 is at distance 1, s2 and s3 at distance 2, s4 at distance 3. Depths 0, 1 and 2 are UNSAT and depth 3 returns the counterexample s0 s1 s3 s4. If the bad state is unreachable, every depth is UNSAT and BMC never terminates on its own, which is the gap completeness thresholds and induction fill.

## Counters and modular arithmetic

Datapath state often makes the depth computable by arithmetic. An $n$-bit register that starts at $s$ and adds $c$ modulo $2^n$ each cycle holds $s + ck \bmod 2^n$ after $k$ cycles. The assertion \"it never equals $t$\" first fails at the smallest $k \ge 0$ with $s + ck \equiv t \pmod{2^n}$. When $c$ is odd it is invertible modulo $2^n$, so a solution always exists and is $k \equiv (t - s) c^{-1}$. Example with 4 bits, $s = 3$, $c = 5$, $t = 0$: we need $5k \equiv 13 \pmod{16}$. The inverse of 5 modulo 16 is 13, since $5 \cdot 13 = 65 \equiv 1$, so $k \equiv 13 \cdot 13 = 169 \equiv 9$. Check: $3 + 45 = 48 \equiv 0$. BMC must reach depth 9. With a 32-bit counter the same question can need billions of frames, which BMC cannot reach.

## Stepping the bound

Scripts often check depths in steps of $d$ instead of 1, with each check covering all frames up to its depth (bad in any frame). The first SAT answer then comes at the first checked depth at or beyond the true depth $D$, which is $\lceil D/d \rceil \cdot d$, after $\lceil D/d \rceil + 1$ solver calls counting the depth-0 call. With $D = 7$ and $d = 3$, the checks are at depths 0, 3, 6 and 9, and the fourth call finds the bug: $\lceil 7/3 \rceil + 1 = 4$. Fewer calls save the per-call overhead, but each call is larger, and the returned trace may be longer than necessary unless the solver is asked for the earliest failing frame. Incremental solvers make stepping by 1 cheap, which is why most engines step one frame at a time.

## Going deeper

Biere et al., "Bounded Model Checking" (Advances in Computers, 2003); Een and Sorensson, "Temporal Induction by Incremental SAT Solving" (BMC Workshop, 2003).
