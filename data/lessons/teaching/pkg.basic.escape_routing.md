# BGA escape routing

Escape routing is the problem of bringing every signal out from inside a dense ball or bump array to the array's edge, so it can be routed onward. Inner balls can only leave through the channels between outer balls, so the channel capacity sets how many rows escape on each layer. This one inequality drives layer counts in package substrates and printed circuit boards alike, and it is a classic constrained routing problem for package routers.

## Traces through one channel

Between two adjacent pads of diameter $d$ at pitch $p$ the free gap is $p - d$. Placing $n$ traces of width $w$ in it needs $n$ widths and $n+1$ spaces $s$ (one to each pad and between traces), so $n w + (n+1) s \le p - d$, which gives $n = \lfloor (p - d - s)/(w + s) \rfloor$. For a 1000 um pitch with 500 um pads and 75/75 um line and space, $n = \lfloor (1000 - 500 - 75)/150 \rfloor = 2$ traces per channel. Shrinking the pad or the line and space adds channel capacity; shrinking the pitch removes it fast, which is why fine-pitch BGAs need fine-line boards with via-in-pad or microvias.

## Rows per layer and layer count

The outer row of balls escapes directly. Each channel between outer-row balls carries $n$ traces, and on each side of the array there is about one channel per ball in a row, so $n$ further rows can escape through the channels. One layer therefore escapes $n + 1$ rows. An array with $R$ rows from edge to centre needs $\lceil R/(n+1) \rceil$ routing layers, ignoring power and ground balls that drop straight down through vias. With two traces per channel and 15 rows to escape, five layers are needed. This simple count is the first estimate a planner makes before running a router, and it explains why depopulating rows or assigning inner balls to power reduces the board layer count.

## Balls escaped on the first layer

For an $N \times N$ full array, the first layer escapes the outer $k = n+1$ rings. The balls inside them number $(N - 2k)^2$ when $N > 2k$, so the escaped count is $N^2 - (N - 2k)^2$. For $N = 30$ and $k = 3$ that is $900 - 576 = 324$ balls on the first layer. The second layer, reached through vias, escapes the next $k$ rings, and so on. The rings shrink toward the centre, so later layers carry fewer balls each; a planner exploits this by putting power and ground at the centre where they need vias anyway, keeping signal rings on few layers. Package routers solve escape as a flow problem that respects these channel capacities before detailed routing.

## Going deeper

Yan and Wong, "Recent research development in PCB layout", ICCAD (2010); Fang, Lee and Chang, "A network-flow-based RDL routing algorithm for flip-chip design", IEEE TCAD (2007).
