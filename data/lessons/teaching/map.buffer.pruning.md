# Dominance pruning of buffer candidates

Without pruning, van Ginneken's candidate lists would double at every buffer site. Pruning keeps only candidates that could still win: those not beaten in both load and required time. The pruned list has a simple monotone shape, which is what makes the algorithm fast and what lets the driver pick the best solution with one scan.

## Dominance

Candidate (C1, q1) **dominates** (C2, q2) if $C_1 \le C_2$ and $q_1 \ge q_2$: it presents no more load upstream and needs the signal no earlier. Whatever happens upstream, a dominated candidate can never give a better result, because every operation, wire, buffer or merge, is monotone in C and q. From (35, 410), (55, 360), (20, 325), (20, 375), (25, 350), the survivors are (20, 375) and (35, 410). Sorted by load, a pruned list always has strictly increasing required times, so it is a staircase: the Pareto front of load against required time.

## Choosing at the driver

At the source, the driver with resistance $R_d$ turns each candidate into a required time at the driver, $q - R_d C$. Scan the pruned list and keep the maximum, then trace back the stored decisions to recover the buffer placement. With candidates (13, 355), (25, 383), (41, 390), (49, 418), (64, 437) and $R_d$ = 2.9 kOhm, the values are 317.3, 310.5, 271.1, 275.9 and 251.4 ps, so the light-load candidate wins. A weaker driver favours candidates with small load, typically the buffered ones.

## Pruning after a buffer site

At a buffer site each candidate spawns a buffered one, and all buffered candidates share the load $C_b$. So at most one of them, the one with the largest required time, can survive, and the list grows by at most one per site. That is the source of van Ginneken's polynomial bound: for n sites and one buffer type, lists have O(n) entries and total work is $O(n^2)$; Shi and Li later reduced this to $O(n \log n)$. With b buffer types, each site adds up to b candidates and the lists grow accordingly, which motivates Lillis's extensions and further pruning by slew and area.

## Going deeper

van Ginneken, ISCAS 1990; Lillis, Cheng and Lin, "Optimal wire sizing and buffer insertion for low power and a generalized delay model", IEEE JSSC 1996; Shi and Li, "A fast algorithm for optimal buffer insertion", IEEE TCAD 2005.
