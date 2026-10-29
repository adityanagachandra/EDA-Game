# Constant multiplication by signed digits

Multiplying by a constant is common in datapaths (filters, transforms, address arithmetic), and synthesis never builds a full multiplier for it. It expands the constant into shifts and adds. Canonical signed-digit recoding, which allows digits of -1, reduces the number of nonzero digits and therefore the number of adders.

## Shift and add from binary

$y = Kx$ with $K$ a constant is a sum of shifted copies of $x$, one per 1 bit of $K$. $75 = 1001011_2$ has four 1 bits, so $75x = (x \ll 6) + (x \ll 3) + (x \ll 1) + x$, which needs three adders. Shifts are wires, free in hardware. In general a constant with $p$ one bits needs $p - 1$ two-input adders. The cost model is adder count, with each adder sized to the result width, and the dominant cost of a constant multiplier comes from long runs of ones, such as $255 = 11111111_2$, which would need seven adders.

## Canonical signed digits

Allow digits -1, written $\bar{1}$. A run of ones $0111\dots1$ equals $1000\dots0\bar{1}$, so $255 = 2^8 - 1$ needs one subtractor. The **canonical signed-digit** (non-adjacent form) representation has no two adjacent nonzero digits, is unique, and has the fewest nonzero digits of any signed-digit form. Compute it from the least significant bit: if $K$ is odd, take digit $2 - (K \bmod 4)$, which is $+1$ or -1, and subtract it; then halve. For 267 the digits, most significant first, are $1\,0\,0\,0\,1\,0\,\bar{1}\,0\,\bar{1}$: four nonzero digits, against four 1 bits in binary.

## Savings and beyond single constants

With $z$ nonzero CSD digits the multiplier needs $z - 1$ adders or subtractors. CSD never uses more nonzero digits than binary, and on average uses about one third of the bit width against one half for binary. Some constants gain nothing: $134 = 10000110_2$ has three ones and CSD $1\,0\,0\,0\,1\,0\,\bar{1}\,0$ also has three digits. Larger savings come from sharing: in multiple constant multiplication (a filter multiplying one input by many constants), intermediate sums such as $3x$ and $5x$ are reused across constants, and algorithms build adder graphs that beat per-constant CSD. Synthesis tools and HLS compilers apply these recodings before bit-level optimization.

## Going deeper

Reitwiesner, "Binary arithmetic", Advances in Computers (1960), on the canonical recoding; Voronenko and Puschel, "Multiplierless multiple constant multiplication", ACM Transactions on Algorithms (2007).
