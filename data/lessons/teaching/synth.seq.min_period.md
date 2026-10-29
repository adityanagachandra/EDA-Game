# Minimum-period retiming

Minimum-period retiming finds the lags that make the clock period as small as possible. It combines the period computation, the W and D matrices and a feasibility test into a search over candidate periods. The answer is bounded below by the largest node delay and by the delay per register of every loop, and it tells the designer how much a register-only change can buy.

## The starting period

Before retiming, the period is the longest path through register-free edges, summing node delays. In a three-node ring $v_1 \to v_2 \to v_3 \to v_1$ with delays 3, 2 and 6 and registers 1, 1, 2, every edge carries a register, so every path between registers is a single node and the period is the largest node delay, 6 ns. That is already optimal: no retiming can split a node. When registers are bunched on a few edges, long register-free chains appear and the period is much larger than any node delay. That gap between the current period and the bounds is what retiming can recover.

## Searching for the minimum

The minimum period is one of the values $D(u, v)$, since the period of any retimed circuit is the delay of some register-free path, and each such path achieving a period is some $D$. Sort the distinct $D$ values, and binary-search for the smallest $c$ for which the constraint system ($r(u) - r(v) \le w(e)$ per edge, $r(u) - r(v) \le W(u, v) - 1$ per pair with $D(u, v) > c$) is feasible. Each test is a Bellman-Ford run, giving $O(|V|^3 \log |V|)$ overall; the FEAS algorithm of Leiserson and Saxe avoids building the full matrices. On tiny circuits one can also enumerate all lag vectors within the bounds $-W(v_1, v) \le r(v) \le W(v, v_1)$, which checks the same answer.

## Lower bounds and the register cost

Two bounds limit what retiming can achieve. The period is at least the largest node delay. And every cycle keeps its register count, so a cycle with total delay $D_C$ and $W_C$ registers forces a period of at least $D_C / W_C$; the maximum over cycles of this ratio is the iteration bound of the loop. Many lag vectors usually reach the minimum period, and they can differ in register count. Picking the cheapest is minimum-area retiming under a period constraint: for a three-node ring with minimum period 5 ns, the cheapest retiming meeting it keeps 5 registers. Production flows retime for period first, then for area at that period.

## Going deeper

Leiserson and Saxe, "Retiming synchronous circuitry", Algorithmica (1991), algorithms OPT1, OPT2 and FEAS; Shenoy and Rudell, ICCAD (1994).
