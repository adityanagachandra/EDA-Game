# Differential parasitic checks for ECOs

When an agent claims that an ECO rerouted three nets, the claim is checkable: re-extract, compare every net's parasitics before and after, and confirm that only the claimed nets and their geometric neighbours moved. This differential parasitic check is the golden test for extraction-level changes. It turns an agent's narrative into a verified diff and catches the side effects that make late ECOs dangerous.

## Which nets moved

Extract the design before and after the ECO with the same extractor version, corner and calibration tables, so that the only difference is the layout change. For each net, compare total capacitance (and, for a thorough check, resistance and each coupling). A net has moved if its change exceeds a tolerance chosen above numerical noise, for example 0.02 fF. With six nets changing $8.1 \to 6$, $2.9 \to 2.9$, $5.8 \to 6.4$, $1.5 \to 0.1$, $1 \to 1$ and $2.1 \to 2.1$ fF, three moved. Determinism matters: if extracting the same layout twice gives different values, the tolerance must absorb that noise and small real effects disappear. That is why extractor seams must be deterministic and seedable.

## Unexpected movers

The ECO's claim defines the expected set: the rerouted nets plus their geometric neighbours, whose coupling legitimately changes when a wire moves beside them. A moved net outside that set is an unexpected mover, evidence that the ECO touched more than it claimed, for example an accidental reroute, a deleted via, or a fill change. Each must be reported with its geometry and change so that an engineer or agent can tell a real side effect from an incomplete expectation, such as a neighbour the claim forgot. The neighbour set itself is computed geometrically by the extractor, from the old and new routes, so that the agent cannot widen it to hide a side effect.

## The pass rule

Two kinds of failure count. Every unexpected mover is one. Every claimed net that did not move is another: if the ECO says it rerouted a net and its parasitics are unchanged, the reroute did not happen or did nothing, and the claim is wrong. With four unexpected movers and one claimed net unchanged, the check reports five failures. The ECO is accepted only at zero. This mirrors the differential timing check for timing engines, and the two compose: the parasitic diff proves the change was confined, and a differential timing report on the new parasitics proves its effect on slack was the intended one. Together they are the evidence an agent hands a human reviewer.

## Going deeper

The golden and differential testing practice described in the agent-first EDA material of this curriculum; McKeeman, "Differential testing for software", Digital Technical Journal (1998).
