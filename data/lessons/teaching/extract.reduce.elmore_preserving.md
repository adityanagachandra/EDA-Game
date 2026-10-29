# Elmore-preserving elimination in an RC chain

The conductance-proportional split of an eliminated node's capacitance is not just a heuristic: on an RC chain it preserves the Elmore delay at every downstream node exactly. Working the algebra on a small chain shows why, quantifies how badly naive lumping fails, and gives a reducer designer a concrete invariant to test against.

## Splitting the eliminated capacitance

Take a chain: driver, $R_1$, node n1 with $C_1$, $R_2$, node n2 with $C_2$, and possibly more downstream. Eliminating n1 replaces $R_1$ and $R_2$ with one resistor $R_1 + R_2$. The conductances of the removed resistors are $1/R_1$ and $1/R_2$, so the share of $C_1$ moved to n2 is $(1/R_2)/(1/R_1 + 1/R_2) = R_1/(R_1 + R_2)$ and the share moved to the driver end is $R_2/(R_1+R_2)$. With $R_1 = 0.6$ kohm, $R_2 = 0.8$ kohm and $C_1 = 2$ fF, n2 gains $2 \times 0.6/1.4 = 0.857$ fF and the driver end gets 1.143 fF. The node closer through lower resistance gets the larger share, which matches intuition: n1's voltage tracks the neighbour it is more tightly tied to.

## The Elmore delay is unchanged

Before elimination the Elmore delay at n2 (with downstream capacitance $C_d$ beyond it) is $R_1(C_1 + C_2 + C_d) + R_2(C_2 + C_d)$. After it, the single resistor $R_1 + R_2$ charges n2's new capacitance plus everything downstream: $(R_1 + R_2)(C_2 + C_d + C_1 R_1/(R_1 + R_2)) = R_1 C_1 + (R_1 + R_2)(C_2 + C_d)$, identical. The capacitance placed at the driver end is behind no resistance, so it adds load but no wire delay. For $R_1 = 0.7$, $R_2 = 0.5$ kohm, $C_1 = 11$, $C_2 = 12$ fF, the delay is $0.7 \times 23 + 0.5 \times 12 = 22.1$ ps before and after. Total capacitance is also preserved, since the two shares sum to $C_1$. These two equalities are the invariants a reducer's tests should assert on every net.

## The cost of naive lumping

A naive reducer moves all of $C_1$ onto n2. Then $C_1$ sits behind $R_1 + R_2$ instead of behind $R_1$ alone, adding $R_2 C_1$ of delay that does not exist at n2 and every node beyond. With $R_2 = 0.4$ kohm and $C_1 = 11$ fF the error is 4.4 ps on a net whose true delay is 20.8 ps, about 21%. Dropping $C_1$ instead underestimates both delay and driver load. The conductance split is exact only for the first moment; higher moments, which affect slew, change slightly, and the error grows with $C_1$'s time constant. That is why elimination is restricted to fast nodes, and why a reducer's regression suite compares slews as well as delays against the unreduced network.

## Going deeper

Sheehan, "TICER: realizable reduction of extracted RC circuits", ICCAD (1999); Rubinstein, Penfield and Horowitz, "Signal delay in RC tree networks", IEEE Transactions on CAD (1983).
