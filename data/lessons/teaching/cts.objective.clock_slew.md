# Clock slew and the transition limit

Clock edges must be sharp. A slow clock transition makes flop timing unpredictable, increases short-circuit power in every clock pin and makes the tree more sensitive to noise and variation. Libraries and design rules therefore impose a maximum transition (slew) on clock nets, usually tighter than on data nets. The slew limit is what forces a CTS engine to insert buffers at all: it caps how much capacitance and wire one driver may see.

## The single-pole estimate

Model a buffer driving a lumped load as a resistance $R$ charging a capacitance $C$. The output is $v(t) = V(1 - e^{-t/\tau})$ with $\tau = RC$. It crosses 10% of the swing at $t_{10} = \tau \ln(10/9)$ and 90% at $t_{90} = \tau \ln 10$, so the 10-90% transition time is
$$t_{10\text{-}90} = \tau \ln 9 \approx 2.197\,RC.$$
For $R = 700\ \Omega$ and $C = 70$ fF, $RC = 49\,000\ \Omega\cdot$fF $= 49$ ps and the slew is about 107.7 ps. The estimate ignores the nonlinear driver and the wire's distributed resistance, but it is the right first-order tool for sizing and fanout decisions.

## Maximum load for a slew limit

Invert the estimate. A limit $S$ allows at most $C_{max} = S/(R \ln 9)$. With $S = 150$ ps and $R = 350\ \Omega$, $C_{max} = 150\,000/(2.197 \cdot 350) = 195.1$ fF. In a library this is what the max_capacitance and max_transition attributes encode for each cell, and CTS reads them as hard constraints. Remember the unit bookkeeping: ohms times femtofarads give femtoseconds, so divide by 1000 for picoseconds.

## From load budget to fanout

Each sink costs its clock-pin capacitance plus a share of wire. If a buffer of $650\ \Omega$ must meet 110 ps, its budget is $110\,000/(2.197 \cdot 650) = 77.0$ fF; at 2.5 fF of pin and 4 fF of wire per sink it can drive $\lfloor 77.0/6.5 \rfloor = 11$ sinks. That number, the maximum fanout, sets how wide each tree level can be and so how many levels the tree needs. Stronger buffers (lower $R$) raise the fanout but present more input capacitance to the stage above, which is the trade-off that logical effort makes precise.

## Going deeper

Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits: A Design Perspective" (2nd ed., 2003), chapters on delay and interconnect; the Liberty reference manual's descriptions of max_transition and max_capacitance.
