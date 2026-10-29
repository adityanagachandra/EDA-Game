# Timing windows and aggressor filtering

Assuming every aggressor switches against every victim at the worst moment is safe but very pessimistic. A timer knows, for each net, the window of time in which it can switch. An aggressor whose window does not overlap the victim's cannot affect the victim's transition. Filtering aggressors by window overlap is the main tool for removing crosstalk pessimism, and it couples timing and SI into an iteration.

## Window overlap

Static timing computes, for each net, the earliest and latest times it can switch: its timing window $[e, l]$. An aggressor with window $[e_a, l_a]$ can affect a victim with window $[e_v, l_v]$ only if the intervals overlap, $e_a \le l_v$ and $l_a \ge e_v$. For a victim switching in [200, 280] ps, aggressors in [350, 400], [220, 260] and [210, 250] ps give two overlaps; the first aggressor has finished switching long before or will start long after, as far as this victim's transition is concerned. Non-overlapping aggressors are treated as quiet, Miller factor 1, rather than removed: their coupling still loads the victim. Some tools pad windows by the aggressor's slew to account for the finite transition.

## Pessimism removed

Without windows, setup analysis must assume all aggressors oppose, adding $0.69 R_d \sum C_c$ of delay over the quiet case. With windows, only overlapping aggressors count. The pessimism removed is $0.69 R_d \sum_{\text{non-overlap}} C_c$. For a 2.3 kohm victim driver with aggressors of 0.9, 2 and 0.6 fF, of which only the 0.9 fF one overlaps, the full delta is $0.69 \times 2.3 \times 3.5 = 5.55$ ps and the windowed delta is 1.43 ps: 4.13 ps of pessimism removed. On a chip with millions of coupled nets, this is often the difference between hundreds of false SI violations and a handful of real ones. Logical correlation, such as two outputs that can never switch in opposite directions, removes more.

## The windowed delta and the iteration

The reported delta delay counts only overlapping aggressors, each at factor 2: for a 1.7 kohm driver and overlapping couplings totalling 6.5 fF, $0.69 \times 1.7 \times 6.5 = 7.62$ ps. But the delta changes the victim's arrival, which moves its window, and every aggressor's window depends on its own crosstalk delays in turn. SI-aware timing therefore iterates: compute windows without crosstalk, compute deltas, recompute windows, repeat until windows stop moving. Zhou, Shenoy and Nicholls showed that this is a fixpoint computation on a lattice, which explains when and to what it converges; in practice it settles in a few rounds. An incremental timer must re-run this locally after every ECO, since moving one wire changes couplings and windows for its neighbours.

## Going deeper

Sapatnekar, "Timing" (2004), the chapter on crosstalk-aware timing; Zhou, Shenoy and Nicholls, "Timing analysis with crosstalk is a fixpoint on a complete lattice", IEEE Transactions on CAD (2003).
