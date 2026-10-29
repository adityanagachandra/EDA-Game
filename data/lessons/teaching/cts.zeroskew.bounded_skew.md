# Bounded-skew clock trees

Zero skew is stricter than timing needs. If the timer tolerates a skew budget $B$, a merge no longer needs one exact tapping point: any tap whose merged skew stays within $B$ will do. That freedom is spent on shorter wire and fewer detours. Bounded-skew trees generalise zero-skew merging by tracking each subtree's minimum and maximum sink delay and by replacing merging segments with merging regions.

## Delay ranges instead of one delay

A bounded-skew subtree is summarised by the interval $[a, b]$ of its sink delays from its root, with $b - a \le B$. Merge subtree 1 ($[a_1, b_1]$) and subtree 2 ($[a_2, b_2]$) on a wire of length $L$, with the tap at $x$ from root 1, under path-length delay. From the tap, the latest sink is at $\max(x + b_1, L - x + b_2)$ and the earliest at $\min(x + a_1, L - x + a_2)$. The merged skew is their difference. With subtree delays $[120, 135]$ and $[90, 100]$, $L = 360$ and $x = 80$, the latest sink is at 380 and the earliest at 200, a skew of 180 µm, far outside any reasonable bound.

## The feasible tapping interval

Requiring merged skew at most $B$ gives two binding inequalities, slow side 2 against fast side 1 and slow side 1 against fast side 2:
$$\frac{L + b_2 - a_1 - B}{2} \le x \le \frac{L + a_2 - b_1 + B}{2},$$
clamped to $[0, L]$. With $[100, 115]$, $[80, 105]$, $L = 140$ and $B = 30$, the lower end is $(140 + 105 - 100 - 30)/2 = 57.5$ µm. When not clamped, the interval's width is $B - \tfrac{1}{2}[(b_1 - a_1) + (b_2 - a_2)]$: each subtree's own skew uses up part of the budget. A zero-skew merge of zero-skew subtrees is the special case of width zero.

## Spending the freedom

In the plane, the feasible taps form a merging region rather than a segment, and bounded-skew DME keeps these regions bottom-up and embeds top-down as before, choosing points that shorten the parent's edge. With $B = 0$ the method reduces to zero-skew DME; with $B$ unbounded it approaches a minimum Steiner tree, which ignores skew. Practical CTS uses a budget derived from the timing constraints, since the same skew that saves wire is taken from setup or hold margins. The region representation is more complex than a segment, which is the price of the extra freedom.

## Going deeper

Cong, Kahng, Koh and Tsao, "Bounded-skew clock and Steiner routing", ACM Transactions on Design Automation of Electronic Systems (1998).
