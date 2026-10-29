# Kogge-Stone versus Brent-Kung

Kogge-Stone and Brent-Kung are the two ends of the prefix-adder design space. Kogge-Stone has minimum depth and fanout 2 but many nodes and long wires. Brent-Kung has nearly the minimum possible node count but almost twice the depth. Counting their nodes and levels exactly is the basis for any architecture comparison, whether by a human, a datapath synthesizer, or an agent exploring the space.

## Kogge-Stone: minimum depth, fanout two

In Kogge-Stone, at level $\ell = 0, 1, \dots, \log_2 n - 1$, every position $i \ge 2^\ell$ combines its current span with the span ending at $i - 2^\ell$. Level $\ell$ therefore has $n - 2^\ell$ nodes, and the total is

$$\sum_{\ell=0}^{\log_2 n - 1} (n - 2^\ell) = n \log_2 n - n + 1.$$

For $n = 16$: depth 4 and $64 - 16 + 1 = 49$ nodes. Every node drives at most two others, so no buffering is needed, but the wire from position $i - 2^\ell$ to $i$ spans $2^\ell$ bit pitches, so the last level has wires reaching half the adder's width. Wiring, not logic, is Kogge-Stone's cost.

## Brent-Kung: few nodes, twice the depth

Brent-Kung computes the prefixes in two sweeps. The up-sweep builds a binary tree over the inputs: $n/2 + n/4 + \dots + 1 = n - 1$ nodes in $\log_2 n$ levels, producing the prefixes at positions $2^k - 1$. The down-sweep fills in the remaining prefixes with $n - 1 - \log_2 n$ nodes in $\log_2 n - 1$ levels. Totals:

$$\text{nodes} = 2n - 2 - \log_2 n, \qquad \text{depth} = 2\log_2 n - 1.$$

For $n = 16$: 26 nodes in 7 levels, against Kogge-Stone's 49 nodes in 4 levels. Fanout stays at most 2 and wires are short, so Brent-Kung is attractive when timing is loose and area or power matters.

## Comparing the two

At $n = 64$, Kogge-Stone has depth 6 and $384 - 64 + 1 = 321$ nodes; Brent-Kung has depth 11 and $128 - 2 - 6 = 120$ nodes, a difference of 201 nodes for 5 extra levels. If each level costs one prefix-node delay, the Kogge-Stone adder is almost twice as fast in logic, but its long wires add RC delay that the node count does not show, and its area and power are more than double.

Intermediate networks interpolate: Han-Carlson runs Kogge-Stone on odd positions only and adds one level to fix the even ones; Ladner-Fischer and Knowles families trade fanout against node count. Datapath synthesis typically starts from a fast structure and relaxes it where slack allows. Comparing candidates on depth and node count alone is a fast first filter; only placement and routing reveal the wire cost.

## Going deeper

Kogge and Stone, "A Parallel Algorithm for the Efficient Solution of a General Class of Recurrence Equations" (IEEE Transactions on Computers, 1973); Brent and Kung, "A Regular Layout for Parallel Adders" (IEEE Transactions on Computers, 1982); Harris, "A Taxonomy of Parallel Prefix Networks" (2003).
