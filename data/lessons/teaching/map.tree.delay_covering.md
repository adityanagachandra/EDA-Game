# Delay-optimal tree covering

Replace area by arrival time and the tree DP minimizes delay instead. With load-independent gate delays the recurrence is exact, because a node's best arrival depends only on what lies below it. The interesting part is the trade-off: the delay-optimal cover usually uses more area, and choosing among equally fast matches by area is the first step of area recovery.

## The arrival-time recurrence

Give each cell a fixed delay d(m) and each primary input an arrival time. The earliest arrival at node v over all covers of its subtree is
$$\text{arr}(v) = \min_{m} \Big( d(m) + \max_{u \in \text{inputs}(m)} \text{arr}(u) \Big).$$
The max replaces the sum: a gate fires when its latest input arrives. With INV = 1, NAND2 = 2, NAND3 = 3 and inputs at 0, the node NAND(INV(NAND(a, b)), c) can be NAND2 after INV after NAND2 (2 + 1 + 2 = 5) or one NAND3 (3), so arr = 3.

## Late inputs change the answer

Arrival times at the leaves matter. If input c arrives at 4 while a and b arrive at 0, the NAND3 cover gives 3 + 4 = 7, while NAND2(INV(NAND2(a, b)), c) gives 2 + max(3, 4) = 6: the late signal should enter as close to the output as possible, through a short path. The DP finds this automatically because it compares the maximum over each match's inputs. The same principle drives pin swapping and timing-driven decomposition later in the flow: restructure so that late signals see few levels.

## Area among delay-optimal covers

Many covers often reach the same arrival time. Breaking ties by area, keep the match with the smallest arrival and, among those, the smallest total area of its subtree cover. This lexicographic DP is still bottom-up and linear. It does not give the smallest-area cover that meets a deadline, because a node that is not critical may afford a slower but cheaper cover; that needs required times propagated from the outputs, which is the job of area recovery. The gap between the area of the delay-optimal cover and the area-optimal cover measures what timing costs.

## Going deeper

Touati, "Performance-Oriented Technology Mapping", PhD thesis, UC Berkeley, 1990; Cong and Ding, "On area/depth trade-off in LUT-based FPGA technology mapping", IEEE TVLSI 1994, for the same trade-off in LUT form.
