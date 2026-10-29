# Voltage sources in MNA

Plain nodal analysis cannot express a voltage source: its current is not a function of its terminal voltages. Modified nodal analysis fixes this by adding the source's current as a new unknown and its defining equation as a new row. This augmentation, which gives MNA its name, is also how inductors, controlled voltage sources and ammeters enter the system. It creates zeros on the diagonal, which shapes how SPICE must factor the matrix.

## Why a voltage source needs a branch current

A voltage source fixes $v_+ - v_- = V$ for any current. In KCL at its terminals that current appears, but no conductance can produce it. MNA therefore adds the branch current $i_V$ as an unknown and the constraint as an equation. Every independent voltage source adds one unknown and one row; resistors, capacitors (open at DC) and current sources add none. A circuit with 3 non-ground nodes and 3 voltage sources has a 6 by 6 system. Ordering matters for reading entries: SPICE-style tools put the node voltages first, in node order, then the branch currents in the order the sources appear.

## The augmented stamp

The stamp for V between n+ and n- has two parts. In the branch-current column, put $+1$ in row n+ and -1 in row n-: the current $i_V$ leaves n+ into the source and arrives at n-, which enters KCL at those nodes. In the branch row, put $+1$ in column n+ and -1 in column n-, and the value $V$ on the right-hand side. The block is the transpose of the column part, and the diagonal of the branch row is zero. For nodes 1 and 2 with V1 between node 1 and ground, entries (1,3) and (3,1) are $+1$, entry (3,3) is 0 and $b_3 = V$. A source between two nodes puts four $\pm 1$ entries.

## The sign of the branch current

SPICE reports the branch current with a fixed convention: positive current flows into the + terminal, through the source, and out of the - terminal. A source that delivers power pushes current out of its + terminal into the circuit, so its reported current is negative. For V1 1 0 1 driving 4 k$\Omega$ plus 1 k$\Omega$ in series, and also 0.25 k$\Omega$ directly to ground, the source delivers $1/5 + 1/0.25 = 4.2$ mA and reports -4.2 mA. Engineers reading .op output see supply currents as negative numbers; a positive supply current means the source is absorbing power, which is a useful check for a wrongly connected circuit.

## Going deeper

Ho, Ruehli and Brennan, "The modified nodal approach to network analysis" (IEEE TCAS, 1975); Najm, "Circuit Simulation" (Wiley, 2010), chapter 2.
