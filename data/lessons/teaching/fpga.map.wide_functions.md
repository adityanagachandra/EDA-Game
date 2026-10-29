# LUTs needed for a wide function

When a function has more inputs than a LUT, the mapper must build it from several LUTs. Simple bounds tell a tool designer how many LUTs and levels any mapping needs, and constructive decompositions such as Shannon expansion show what an arbitrary function costs. These numbers are the sanity checks on a mapper's output and the basis of LUT-count estimators.

## LUT count for associative functions

A $k$-LUT consumes up to $k$ signals and produces one, so it reduces the number of live signals by at most $k - 1$. Combining $n$ inputs into one output therefore needs at least $\lceil (n-1)/(k-1) \rceil$ LUTs. For associative, commutative functions such as AND, OR and XOR, a tree of LUTs meets the bound exactly. A 22-input AND on 3-LUTs needs $\lceil 21/2 \rceil = 11$ LUTs; on 6-LUTs it needs $\lceil 21/5 \rceil = 5$. The bound holds for any function that depends on all $n$ inputs, but non-decomposable functions need far more. It is still useful: a mapper that reports fewer LUTs than this bound for a function of $n$ inputs is wrong.

## Depth of a LUT tree

With $d$ levels of $k$-LUTs, at most $k^d$ primary inputs can reach the output, so a function that depends on $n$ inputs needs depth at least $\lceil \log_k n \rceil$. A 149-input AND on 4-LUTs needs 4 levels because $4^3 = 64 < 149 \le 256 = 4^4$. Larger $k$ reduces depth logarithmically: the same AND on 6-LUTs needs 3 levels since $6^2 = 36 < 149 \le 216$. This is one reason FPGA vendors moved from 4-LUTs to 6-LUTs: fewer levels means fewer slow routing hops. The depth bound is also the floor for the depth-optimal labels a mapper such as FlowMap computes; the labels can exceed it when the network's structure prevents a balanced tree.

## Shannon decomposition of arbitrary functions

An arbitrary function of $n > k$ inputs can be built by Shannon expansion on $n - k$ variables. Each of the $2^{n-k}$ cofactors is a function of $k$ inputs and fills one LUT. A multiplexer tree selected by the expanded variables then picks the right cofactor. A $k$-LUT holds a multiplexer with $s$ selects and $2^s$ data inputs when $s + 2^s \le k$, so a 6-LUT holds a 4:1 mux and a 4-LUT only a 2:1 mux. For a 6-input function on 4-LUTs: 4 cofactor LUTs, then 2:1 muxes in two levels (2 LUTs, then 1), 7 LUTs in total. Real mappers do much better on real logic because cofactors share structure, but the construction gives an upper bound of order $2^{n-k}$ for the worst case.

## Going deeper

Cong and Ding, "Combinational Logic Synthesis for LUT Based Field Programmable Gate Arrays", ACM TODAES (1996); Brown, Francis, Rose and Vranesic, "Field-Programmable Gate Arrays" (1992), chapter on technology mapping.
