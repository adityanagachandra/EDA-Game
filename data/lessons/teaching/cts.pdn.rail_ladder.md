# IR drop along a resistive rail

The simplest IR-drop problem is a single rail feeding a row of cells: a resistor ladder. It already shows the two facts that dominate PDN design. Drop grows with the square of the distance from the supply, because each segment carries the current of everything beyond it. And feeding from both ends cuts the worst drop by about a factor of four. Every strap pitch rule is a consequence of these two facts.

## Current through each segment

Model a rail with $N$ taps, each drawing current $I$, joined by segments of resistance $R$, with one more segment to the supply. By KCL the first segment carries $NI$, the second $(N - 1)I$, and the last $I$. The drop at the far tap is the sum of the segment drops:
$$d_N = RI\sum_{k=1}^{N} k = RI\,\frac{N(N+1)}{2}.$$
For nine taps of 2.25 mA and $0.4\ \Omega$ segments, $d_9 = 0.4 \cdot 2.25 \cdot 45 = 40.5$ mV. Doubling the rail length roughly quadruples the drop.

## The drop at an intermediate tap

At tap $k$, only the first $k$ segments are in the path, carrying $N, N - 1, \ldots, N - k + 1$ tap currents:
$$d_k = RI\left(kN - \frac{k(k-1)}{2}\right).$$
For seven taps of 5 mA and $0.6\ \Omega$ segments, tap 4 sees $0.6 \cdot 5 \cdot 22 = 66$ mV. The profile is concave: the drop rises fastest near the supply, where segments carry the most current. This is why the first segments of a rail, and the vias that feed it, carry the highest current density and are where electromigration is checked first.

## Feeding from both ends

Add a second supply connection at the far end. The system is now $G d = i$ with a tridiagonal $G$: each tap has diagonal $2/R$ and off-diagonals $-1/R$. For five 2 mA taps and $0.55\ \Omega$ segments, the drops are 2.75, 4.4, 4.95, 4.4 and 2.75 mV, against 16.5 mV with one feed. The midpoint of a doubly fed rail sees roughly a quarter of the single-feed worst case, because each half is half as long and drop scales with length squared. Straps at regular pitch do exactly this for every rail.

## Going deeper

Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the chapter on power and ground routing.
