# Linear decompressors

A linear decompressor turns a few tester bits per cycle into many scan-cell values, each a GF(2) linear combination (an XOR) of the free variables the tester injected. Whether a test cube can be delivered becomes a question about a system of linear equations over GF(2). That framing gives the encoding algorithm, Gaussian elimination, and the capacity limit, the rank of the system.

## Cells as XORs of variables

During one pattern the tester injects free variables $x_1, x_2, \ldots$ into the decompressor, a few per shift cycle. Because the decompressor is built from flip-flops and XOR gates, every scan cell ends up holding a fixed XOR of some of those variables, determined by the decompressor structure and the cycle at which the cell was loaded. For example, $s_1 = x_1 \oplus x_3 \oplus x_4$, $s_2 = x_4$, $s_3 = x_1$, $s_4 = x_1 \oplus x_2$. Given variable values $x = (0, 0, 1, 1)$, the cells get the parities: $s_1 = 0 \oplus 1 \oplus 1 = 0$, $s_2 = 1$, $s_3 = 0$, $s_4 = 0$, so $s_1 s_2 s_3 s_4 = 0100$.

## Encoding a test cube

A cube's care bits become equations: each specified cell $s_i = b_i$ is a linear equation over GF(2) in the variables. Encoding solves the system by Gaussian elimination. If it is inconsistent (elimination yields $0 = 1$), the cube cannot be produced and ATPG must drop some target faults into another pattern. If it is consistent with rank $r$ over $n$ variables, exactly $2^{n - r}$ assignments produce it; any one will do, and the remaining freedom fills the don't-care cells pseudo-randomly, which often detects extra faults for free. Don't-care cells impose no equations, which is why sparse cubes encode easily.

## Capacity

A cube with $s$ care bits is guaranteed encodable whatever the care values if its $s$ coefficient vectors are linearly independent. For random coefficient vectors over $n$ variables, that probability is
$$\prod_{i=0}^{s-1} \left(1 - 2^{i - n}\right),$$
which is close to 1 while $s$ is well below $n$ and falls steeply as $s$ approaches $n$ (the last factor is 1/2 at $s = n$). So the number of injected variables per pattern bounds the care bits per pattern, and practical encoders keep a margin of several variables. Designs raise capacity with more channels or by continuously injecting variables during shift, as ring-generator schemes do.

## Going deeper

Rajski et al., "Embedded deterministic test," IEEE Transactions on CAD, 2004; Konemann, "LFSR-coded test patterns for scan designs," European Test Conference, 1991.
