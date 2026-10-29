# Dynamic programming

Dynamic programming solves a problem by combining optimal solutions of overlapping subproblems, computed once and stored. It applies when the problem has optimal substructure: the best completion from a state does not depend on how the state was reached. In EDA it is exact and fast wherever the structure is a chain or a tree: discrete sizing along a path, tree-covering technology mapping, van Ginneken buffer insertion, and channel and maze computations.

## The Bellman recurrence

Choose one of several sizes for each gate in a chain. Each choice has its own delay cost, and each adjacent pair adds a loading cost that depends on both sizes. Define the cost-to-go $V_i(s)$ as the cheapest total for gates $i$ through the end, given that gate i has size s. The last gate's cost-to-go is its own cost, and earlier ones satisfy
$$V_i(s) = c_i(s) + \min_{s'}\bigl[t_i(s, s') + V_{i+1}(s')\bigr].$$
Example: gate 3 has $V_3(S) = 7$ and $V_3(M) = 9$; gate 2 at size M costs 3, with transition costs $t(M, S) = 5$ and $t(M, M) = 6$. Then $V_2(M) = 3 + \min(5 + 7, 6 + 9) = 3 + 12 = 15$. The best way to finish from gate 2 at size M is to size gate 3 small, whatever gates before 2 did.

## Solving the whole chain

Compute the cost-to-go backward from the last gate to the first, then take the minimum over the first gate's sizes. Example with sizes S and M: gate costs 9 and 6, then 9 and 1, then 2 and 2; transitions from gate 1 to 2: SS 6, SM 3, MS 6, MM 3; from gate 2 to 3: SS 3, SM 6, MS 6, MM 2. Backward: $V_3 = (2, 2)$. $V_2(S) = 9 + \min(3 + 2, 6 + 2) = 14$ and $V_2(M) = 1 + \min(6 + 2, 2 + 2) = 5$. $V_1(S) = 9 + \min(6 + 14, 3 + 5) = 17$ and $V_1(M) = 6 + \min(6 + 14, 3 + 5) = 14$. The optimum is 14: M, M, M.

Storing the argmin at every state lets a forward traceback recover the assignment. The same table is a shortest path in a layered graph, and the recurrence is the Bellman-Ford relaxation on a DAG in reverse topological order.

## Counting the work

Brute force evaluates every assignment: with k sizes for each of s gates, $k^s$ of them, $3^4 = 81$ for four gates with three sizes and about $3.5 \times 10^9$ for twenty. DP evaluates each transition once: $(s - 1)k^2$ evaluations, 27 for the four-gate example and 171 for twenty gates. The exponential search becomes linear in the chain length because the state, the current gate's size, captures everything the future needs.

That last condition is the whole art. If the future cost also depends on, say, the accumulated arrival time, the state must include it, and the table grows with the number of distinct arrival values; van Ginneken's buffer insertion keeps that growth in check by pruning dominated states. On a general netlist DAG with reconvergent fanout no small state suffices, which is why exact DP sizing works on chains and trees and only heuristically on whole circuits.

## Going deeper

Bellman, "Dynamic Programming" (1957); Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms" (4th ed., 2022), chapter 14; Bertsekas, "Dynamic Programming and Optimal Control" (4th ed., 2017).
