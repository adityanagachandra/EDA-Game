# Tie, filler, tap, endcap and decap cells

Some library cells have no logic function at all, or only a constant one: tie cells, fillers, well taps, endcaps, decoupling capacitors and spare cells. They exist because manufacturing, reliability and power integrity need things that logic does not. Placement, routing and signoff tools insert them at fixed stages of the flow and must know to exclude them from logic checks.

## Tie cells, fillers and spares

A constant input is driven by a tie-high or tie-low cell rather than wired straight to a rail. The tie cell connects the gate through a transistor, isolating the gate oxide from rail noise and electrostatic events, and it gives the netlist an ordinary driven net that tools can treat uniformly. Filler cells fill empty sites after placement and routing so that the n-well, implants and power rails stay continuous along each row; they have no logic. Spare cells are unused gates sprinkled across the die in advance, so that a later engineering change can be implemented by changing only metal masks; some libraries offer programmable gate-array fillers whose uncommitted transistors metal changes can wire into simple gates.

## Decoupling capacitance

When many cells switch together, they draw a burst of current that the power grid, with its resistance and inductance, cannot supply instantly, so the local supply droops. A decap cell is a capacitor between VDD and VSS, usually a transistor's gate oxide, placed near switching logic to supply that charge. If a cluster draws an extra $I$ for a time $\Delta t$ and the droop must stay below $\Delta V$, the decap must hold charge $I \Delta t$ while losing only $\Delta V$:
$$C = \frac{I \Delta t}{\Delta V}.$$
With 20 mA for 100 ps and a 40 mV budget, $C = 20 \times 100 / 40 = 50$ pF, since mA times ps over mV is pF. Decap costs area and gate leakage, so it is placed where power analysis shows droop, often as filler-replacement cells in gaps near hot spots.

## Taps, endcaps and flow placement

In tapless libraries, cells have no well or substrate contacts, so well-tap cells tie the n-well to VDD and the substrate to VSS at a maximum spacing set by latch-up rules. Taps are therefore inserted in a regular pattern before standard-cell placement, and the placer works around them. Endcap (boundary) cells terminate every row and surround blockages, satisfying rules on well and poly termination at row ends. The flow order is fixed: endcaps and taps after floorplanning, logic placement and routing, then fillers and decap fillers in the remaining gaps. Logic-equivalence checks and timing ignore physical-only cells (fillers, taps, endcaps, decaps), while tie cells, which carry constants, stay in the logical netlist.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on physical design flow; the OpenROAD documentation for tapcell and filler insertion.
