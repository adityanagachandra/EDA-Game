# RC network reduction

A detailed extraction can produce tens of millions of RC nodes, far more than a timer needs. RC reduction replaces each net's network with a much smaller one that behaves the same at the pins: same total capacitance, same Elmore delay to each sink, and ideally the same higher moments. Reduction sits between extraction and analysis, and its choices (what to preserve, how aggressively to eliminate, how to keep the model passive) decide whether the timer is both fast and right.

## Why reduce

Delay calculation, crosstalk analysis and power analysis all run per net, with cost growing in the number of nodes. A detailed network from fracturing has a node at every segment boundary, most of which add accuracy no analysis can use: two nodes 0.1 um apart on a wire settle almost simultaneously. Reduction keeps the port behaviour, the response seen at the driver and sinks, while removing internal nodes. The minimum requirements are that the total capacitance (the driver's load) and the first moment, or Elmore delay, to each sink are preserved. Deleting small capacitors is not reduction: it changes both. Real reduction moves charge storage and resistance onto the nodes that remain.

## Choosing the reduced model

Consider a chain: driver, $R_1 = 250$ ohm, node n1 with $C_1 = 2$ fF, $R_2 = 400$ ohm, sink with $C_2 = 3$ fF. The Elmore delay at the sink is $250 \times 5 + 400 \times 3 = 2450$ ohm-fF. Eliminating n1 leaves one 650 ohm resistor and two capacitors, one at the driver end and one at the sink. Splitting $C_1$ in proportion to the conductances of the removed resistors puts $C_1 R_2/(R_1 + R_2) = 1.23$ fF at the driver end and $C_1 R_1/(R_1+R_2) = 0.77$ fF at the sink, so the sink sees $650 \times 3.77 = 2450$ ohm-fF: preserved. Moving all of $C_1$ to the sink overestimates; dropping it loses load; splitting it in half is right only when $R_1 = R_2$.

## Accuracy, passivity and locality

Preserving only the first moment keeps Elmore delay but can distort slew and the 50% delay of resistively shielded nets, which depend on higher moments. Two families of method address this. Node elimination, such as TICER, removes nodes one at a time by a star-mesh transform, choosing nodes whose time constant is small compared with signal transitions; it is local, keeps a physical RC netlist, and suits incremental updates after an ECO. Krylov projection, such as PRIMA, projects the whole network onto a small subspace that matches several moments and is guaranteed passive, which matters because a non-passive model can generate energy and destabilize simulation. Both must treat coupling capacitors like ground capacitors during redistribution, keeping their far terminals on the aggressors.

## Going deeper

Sheehan, "TICER: realizable reduction of extracted RC circuits", ICCAD (1999); Odabasioglu, Celik and Pileggi, "PRIMA: passive reduced-order interconnect macromodeling algorithm", IEEE Transactions on CAD (1998).
