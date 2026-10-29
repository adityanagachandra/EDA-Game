# SCOAP testability measures

ATPG and test-point insertion need a cheap estimate of how hard it is to set a line to 0 or 1. SCOAP, the Sandia Controllability/Observability Analysis Program, computes integer controllabilities in one linear-time pass over the netlist. The numbers are heuristics, not probabilities, but they steer PODEM's backtrace, FAN's choices and the placement of test points in commercial tools.

## Combinational controllability

$CC0(x)$ and $CC1(x)$ estimate the number of line assignments needed to set $x$ to 0 or 1. Primary inputs have $CC0 = CC1 = 1$. For a gate, the cost of producing the controlled output is the cheapest single input at the controlling value plus 1, and the cost of the other output is the sum over all inputs at the non-controlling value plus 1. For AND: $CC0 = \min_i CC0(i) + 1$ and $CC1 = \sum_i CC1(i) + 1$. For OR the roles swap: $CC1 = \min_i CC1(i) + 1$, $CC0 = \sum_i CC0(i) + 1$. NAND and NOR swap the output values of AND and OR. An inverter gives $CC0(out) = CC1(in) + 1$ and $CC1(out) = CC0(in) + 1$. Fanout branches inherit the stem's values.

## Computing a value

Controllabilities are computed in topological order. With $g_1 = $ AND(b, c), $g_2 = $ AND(a, $g_1$), $g_3 = $ OR($g_1$, $g_2$), $g_4 = $ AND($g_2$, $g_3$): $g_1 = (2, 3)$; $g_2 = (\min(1, 2) + 1, 1 + 3 + 1) = (2, 5)$; $g_3 = (2 + 2 + 1, \min(3, 5) + 1) = (5, 4)$; $g_4 = (\min(2, 5) + 1, 5 + 4 + 1) = (3, 10)$. So setting $g_4$ to 1 is much harder than setting it to 0. For a 2-input XOR, $CC1 = \min(CC0_a + CC1_b, CC1_a + CC0_b) + 1$ and $CC0 = \min(CC0_a + CC0_b, CC1_a + CC1_b) + 1$.

## What the numbers mean

SCOAP counts assignments along a tree-like view of the circuit, so it ignores reconvergent fanout: the sum rule double-counts a shared input, and correlated signals can be easier or impossible to set than the number suggests. A line that is constant in reality still gets finite SCOAP values. Despite that, the measures rank lines usefully and cost one pass, which is why ATPG uses them to order decisions and why random-pattern-resistant regions show up as lines with large CC values. Probabilistic measures such as COP give signal probabilities instead and suit random-pattern analysis.

## Going deeper

Goldstein, "Controllability/observability analysis of digital circuits," IEEE Transactions on Circuits and Systems, 1979; Bushnell and Agrawal (2000), chapter 6.
