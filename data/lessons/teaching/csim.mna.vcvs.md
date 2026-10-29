# VCVS stamps and controlled branch rows

A voltage-controlled voltage source forces an output voltage proportional to a controlling voltage. Like an independent voltage source it needs a branch-current unknown, and like a VCCS it couples to controlling nodes. Its stamp combines both patterns, and it is the building block of ideal amplifiers, behavioral macromodels and many textbook test circuits. Reading its entries and its branch current is a direct test of the MNA conventions.

## The VCVS branch row

The card E1 n+ n- nc+ nc- mu imposes $v_{n+} - v_{n-} - \mu(v_{nc+} - v_{nc-}) = 0$. MNA adds its branch current as an unknown and this constraint as a row. The row has $+1$ at column n+, -1 at n-, $-\mu$ at nc+ and $+\mu$ at nc-, with a zero right-hand side. The branch column, as for any voltage source, has $+1$ in row n+ and -1 in row n-. Ground entries are dropped. With three nodes and the cards V1 2 0 1 then E1 1 0 2 0 2, unknowns are $v_1, v_2, v_3$, then the V1 and E1 currents; the E1 row (row 5) has $+1$ at column 1 and -2 at column 2.

## Counting unknowns with controlled sources

The unknown count follows one rule: anything that forces a voltage needs its current. Node voltages count first; then add one branch current for each independent voltage source, each VCVS and each inductor (SPICE keeps inductor currents as unknowns). VCCS, current sources, resistors and capacitors add none. A circuit with 5 nodes, 2 voltage sources, 1 VCVS, 2 inductors, any number of resistors, capacitors and current sources, and a VCCS has $5 + 2 + 1 + 2 = 10$ unknowns. Current-controlled sources add a further branch only for the CCVS, which forces a voltage; the CCCS reuses the controlling source's current.

## Solving an ideal amplifier stage

Consider V1 1 0 Vin, E1 2 0 1 0 mu and a load from node 2 to ground. The E1 row forces $v_2 = \mu v_1 = \mu V_{in}$ at once, and the load current is $v_2/R_L$. E1 delivers that current out of its + terminal, so with SPICE's convention (positive current flows into the + terminal) its branch current is $-\mu V_{in}/R_L$. With $V_{in} = 0.5$ V, $\mu = 2$ and a load of 0.5 k$\Omega$ plus 2.5 k$\Omega$ in series, $v_2 = 1$ V, the source delivers $1/3$ mA and reports -0.333 mA. An ideal VCVS has zero output resistance, so the load never changes $v_2$; real op-amp macromodels add an output resistance.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 2; Vlach and Singhal, "Computer Methods for Circuit Analysis and Design" (2nd ed., 1994).
