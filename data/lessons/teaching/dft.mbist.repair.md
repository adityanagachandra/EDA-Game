# Redundancy allocation and memory repair

Large embedded memories rarely come out of the fab perfect, so they are built with spare rows and spare columns. After MBIST locates the faulty cells, a repair analysis chooses which spares replace which lines, and the result is programmed into fuses or other non-volatile storage. Repair converts many failing dies into good ones, and allocating spares well is a classic combinatorial problem.

## Must-repair analysis

The repair is a covering problem: every faulty cell must be covered by a replaced row or a replaced column. Some choices are forced. If a row has more faulty cells than there are spare columns, columns cannot cover them all, so that row **must** take a spare row; likewise a column with more faulty cells than spare rows must take a spare column. With 1 spare column and a row holding 3 faults, that row is must-repair. Must-repair lines are allocated first, their faults removed, and the analysis repeats, since the remaining spare counts have changed.

## Minimum cover

Without spare limits, the smallest number of lines that covers all faults is a minimum **vertex cover** of the bipartite graph whose vertices are rows and columns and whose edges are faulty cells. By **König's theorem**, in a bipartite graph that minimum equals the size of a maximum matching, so it is computable in polynomial time. Faults at (0, 2), (0, 5), (2, 6), (4, 2) and (6, 6) need 3 lines: for example row 0, column 2 and column 6. With the realistic constraint of at most $R$ spare rows and $C$ spare columns, deciding repairability is NP-complete, and built-in redundancy analysis uses must-repair passes plus bounded search.

## Repair with limited rows

A common sub-question: given at most $R$ spare rows, how many spare columns are needed? Try each choice of up to $R$ rows to replace, count the distinct columns among the remaining faults, and take the minimum. Built-in self-repair (BISR) logic runs this kind of analysis on chip after MBIST so that repair data can be computed at test time or even at power-up. The quality of the allocator shows up directly in **repair yield**, the fraction of dies made good, which feeds the memory's contribution to overall yield.

## Going deeper

Kuo and Fuchs, "Efficient spare allocation for reconfigurable arrays," IEEE Design and Test, 1987; Wang, Wu and Wen (2006), chapter 8.
