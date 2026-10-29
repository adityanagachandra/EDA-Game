# BGA fanout and escape routing

A ball grid array puts hundreds or thousands of connections under one package. Escape routing brings each ball out from under the package to a routing channel, through traces between balls and vias down to other layers. How many traces fit between balls, how many rows each layer can escape, and how many balls sit in the outer rows together decide the board's layer count. These are geometric counts a tool designer can compute before any routing.

## Traces between adjacent balls

Between two adjacent lands of pitch $p$ and land diameter $d$, the gap is $p - d$. Fitting $N$ traces of width $w$ with clearance $s$ to the lands and to each other needs $N w + (N + 1) s \le p - d$, so $N = \lfloor (p - d - s)/(w + s) \rfloor$. Example: a 1.0 mm pitch with 0.5 mm lands and 0.1 mm trace and space gives $(0.5 - 0.1)/0.2 = 2$ traces. At 0.8 mm pitch with 0.4 mm lands and the same rules, $(0.4 - 0.1)/0.2 = 1.5$, so one trace. Finer pitch quickly drops $N$ to zero, at which point traces cannot pass between balls on the surface and the design needs via-in-pad, microvias or finer trace and space, all of which cost money.

## Rows escaped per layer

Under a simple model, the outer row escapes directly outward, and each trace that fits between balls lets one more row escape on the same layer, so each signal layer escapes $N + 1$ rows, from the outside in. Inner rows reach other layers through vias placed between balls (dog-bone fanout) or in the lands. A full $n \times n$ array has $\lceil n/2 \rceil$ concentric rows. Example: a 20 by 20 array has 10 rows; with one trace between balls, each layer escapes 2, so 5 signal layers are needed, before counting the power and ground balls, which go straight to planes and do not need escape. Real escape planning removes those balls first and often does better than the model by routing diagonal channels.

## Counting balls in the outer rows

How many signals each layer carries follows from counting balls in rings. The outer $k$ rows of a full $n \times n$ array contain $n^2 - (n - 2k)^2$ balls. Example: a 20 by 20 array's outer 2 rows hold $400 - 256 = 144$ balls, and the next 2 rows hold $256 - 144 = 112$. Outer rings hold the most balls, so the top layer does the most work, and each deeper layer carries fewer signals. This is why chip and package designers place high-speed signals on outer rows, which escape on the top layers with short via stubs, and put power and ground in the center, where they reach planes through short vias. Ball-map planning that knows the board's escape capacity is a classic package and board co-design problem.

## Going deeper

Holden, "The HDI Handbook" (2009), chapters on BGA escape. Xilinx and Intel FPGA PCB design guides, chapters on BGA breakout, for vendor-specific patterns.
