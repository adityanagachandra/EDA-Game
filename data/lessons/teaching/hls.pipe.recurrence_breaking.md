# Breaking recurrences with partial sums

When RecMII is the binding bound, more hardware cannot help: the loop must be rewritten so that its dependence cycle is shorter or spans more iterations. The classic case is a floating-point accumulation whose adder latency sets the II. Interleaved partial sums, a form of reassociation, raise the iteration distance of the recurrence and recover II = 1, at the cost of a final reduction and changed rounding.

## What sits on the recurrence

Identify which operations lie on the loop-carried cycle. In $s = s + x[i]$, only the addition is on it, so with a pipelined adder of latency $L_a$ and distance 1, $\mathrm{RecMII} = L_a$. In $s = s + a[i] \times b[i]$, the multiplication reads only array values and is off the cycle; it can be computed ahead in the pipeline, so RecMII is still $L_a$. In $s = s \times c + x[i]$, both the multiplication and the addition are on the cycle, so $\mathrm{RecMII} = L_m + L_a$: with a 6-cycle multiplier and a 5-cycle adder, II = 11. Moving every operation that does not depend on the carried value off the cycle is the first step of any recurrence fix, and HLS tools do much of it automatically during scheduling.

## Interleaved partial sums

Keep $P$ partial sums and let iteration $i$ update partial sum $i \bmod P$. Each partial sum is now updated every $P$ iterations, so its recurrence has distance $P$, and the bound becomes $\lceil L_a / P \rceil$. To reach a target II, choose $P \ge \lceil L_a / \mathrm{II} \rceil$: an 8-cycle adder reaches II = 3 with $P = 3$, and II = 1 with $P = 8$. The rewrite is a reassociation of the sum. On integers it is exact. On floating point it changes the order of additions and therefore the rounding, so tools do it only when the user permits, and the verification must use a tolerance. Unrolling an accumulation without this rewrite does not help: the unrolled body still has a single chain through $s$, now even longer.

## Total cost including the final reduction

After the loop, the $P$ partial sums must be combined. A balanced tree takes $\lceil \log_2 P \rceil$ levels of additions, each costing $L_a$ cycles. With $N$ values, the main loop at II = 1 issues its last addition at cycle $N - 1$ and it completes $L_a$ cycles later, so the total is $(N - 1) + L_a + \lceil \log_2 P \rceil L_a$. For $N = 270$, $L_a = 6$ and $P = 6$, that is $269 + 6 + 3 \times 6 = 293$ cycles, against $270 \times 6 = 1620$ for the original loop at II = 6. The extra cost is $P - 1$ registers and a short reduction, negligible for long loops. The same idea, splitting one long dependence chain into several shorter ones, is behind tree-height reduction in straight-line code.

## Going deeper

Kastner, Matai and Neuendorffer, "Parallel Programming for FPGAs" (2018), chapter on loop-carried dependences; Muchnick, "Advanced Compiler Design and Implementation" (1997), on reassociation.
