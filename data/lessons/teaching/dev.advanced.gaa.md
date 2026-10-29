# Gate-all-around nanosheet transistors

Gate-all-around nanosheet transistors stack several thin horizontal silicon sheets and wrap the gate completely around each one. Surrounding the channel on all sides gives better electrostatic control than a FinFET's three sides, and the sheet width is a continuous layout choice, which undoes the width quantization that fins imposed. Nanosheets are the device architecture that follows FinFETs at leading-edge logic nodes.

## Effective width of a nanosheet stack

Each sheet of width $W$ and thickness $T$ conducts along its full perimeter, since the gate surrounds it. The effective width of a stack of $N$ sheets is
$$W_{eff} = N \cdot 2(W + T).$$

Example: 3 sheets, each 30 nm wide and 6 nm thick, give $W_{eff} = 3 \times 72 = 216$ nm.

Sheet thickness is kept to a few nanometres so that the gate controls the whole cross-section, and sheets are separated vertically by the gate stack itself. The fabrication sequence grows alternating silicon and silicon-germanium layers, patterns them, then removes the SiGe sacrificial layers to release the sheets before depositing the gate around them. Inner spacers between the sheets isolate the gate from the source and drain, and their quality matters for parasitic capacitance.

## Nanosheet versus FinFET in one footprint

For a fixed cell footprint, the comparison is effective width per footprint and electrostatic control. A FinFET adds width by adding fins side by side; a nanosheet adds width by stacking sheets vertically and widening them.

Example: two fins of height 50 nm and width 6 nm give $2 \times 106 = 212$ nm. Three sheets of width 35 nm and thickness 5 nm give $3 \times 80 = 240$ nm, a ratio of 1.13, with the gate wrapping every surface.

The ratio depends heavily on sheet width and count, so it can go either way. The real gains are better electrostatics (lower swing and DIBL at the same gate length, enabling shorter gates or lower supply) and design flexibility. Costs include more complex fabrication, higher parasitic capacitance between sheets and gate, and new variability sources such as sheet thickness and sheet-to-sheet differences.

## Continuous width as a design knob

Unlike fin count, sheet width can be drawn at any value within design rules. To reach a target effective width with $N$ sheets of thickness $T$,
$$W = \frac{W_{eff}}{2N} - T.$$

Example: a target of 200 nm with 3 sheets of 5 nm gives $W = 200/6 - 5 = 28.3$ nm.

This flexibility returns a knob that FinFETs took away. Standard-cell libraries can offer sheet widths tuned per cell, and SRAM designers can set pull-up, pull-down and pass-gate strengths with non-integer ratios. It also makes the device model's width dependence matter again: compact models for gate-all-around devices (BSIM-CMG supports GAA geometries) take sheet width, thickness and count as instance parameters, and extraction must fit across widths. Very narrow and very wide sheets behave differently, the former through stronger edge effects and the latter through weaker electrostatics at sheet edges.

## Going deeper

Loubet and others, "Stacked nanosheet gate-all-around transistor to enable scaling beyond FinFET", Symposium on VLSI Technology (2017); the BSIM-CMG technical manual, gate-all-around sections.
