# Multi-Vt assignment trade-offs

A multi-Vt library gives every cell several threshold flavours with identical footprints. Optimizers start with fast, leaky cells where timing demands them and slow, frugal cells elsewhere, then swap cells with spare slack to higher Vt in a leakage-recovery pass. The arithmetic of that trade, leakage saved against delay spent, is what the optimizer evaluates for every candidate swap.

## Mixed-flavour leakage totals

Block leakage is the sum over cells of each cell's leakage. If HVT cells leak $r$ times less than LVT cells, a block with $N_L$ LVT cells of leakage $P_L$ and $N_H$ HVT cells leaks
$$P = N_L P_L + N_H P_L / r.$$
Example: 8,000 LVT cells at 10 nW and 92,000 HVT cells with $r = 10$ give $80 + 92 = 172$ microwatts. The LVT cells are 8% of the count but nearly half the leakage. This skew is typical: a small fraction of fast cells on critical paths dominates leakage, which is why the choice of which cells stay low-Vt matters so much. Real per-cell leakage also depends on the cell type and input state, which libraries record per state.

## Leakage recovery

Leakage recovery swaps cells to a higher Vt wherever the timing slack allows. Swapping $m$ of $N$ LVT cells to HVT saves the fraction
$$\frac{m}{N}\left(1 - \frac{1}{r}\right)$$
of the block's leakage, if all cells leak alike. With $N = 50{,}000$, $m = 35{,}000$ and $r = 10$ the saving is $0.7 \times 0.9 = 63\%$. Because the HVT version of a cell has the same footprint and pin positions, the swap needs no re-placement or re-routing, only an updated timing analysis, which is why recovery is one of the cheapest late-flow optimizations. Footprint-compatible flavours are a library design decision that tools rely on.

## Spending slack on a path

Each swap costs delay. If an HVT gate is $x\%$ slower than its LVT version, a path of $k$ LVT gates of delay $d$ with setup slack $s$ can absorb $\lfloor s / (d x / 100) \rfloor$ swaps, at most $k$. With 12 gates of 20 ps, 50 ps of slack and a 20% penalty, each swap costs 4 ps and 12 swaps fit. But gates are shared: a cell on many paths can only be swapped if every path through it has the slack, so real optimizers order candidates by leakage saved per unit of slack consumed and update timing incrementally as they go. Swapping also affects hold, since slower cells help hold, and it must be checked at all corners, because the delay penalty of HVT grows at low voltage.

## Going deeper

Wei, Roy and others on dual-threshold design, for example "Design and optimization of dual-threshold circuits for low-voltage low-power applications" (IEEE Transactions on VLSI Systems, 1999); the OpenROAD documentation for the repair_timing and leakage-recovery commands.
