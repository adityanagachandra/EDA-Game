# Cell-aware fault models

Stuck-at and transition faults sit on cell pins, but many defects sit inside cells: a short between transistor terminals, an open in a diffusion contact, a resistive bridge in local interconnect. Such a defect may corrupt the cell's output only for one or two input combinations, which a pin-level test set rarely applies. Cell-aware test characterizes these defects once per library cell and hands ATPG a precise list of input conditions, closing a known gap in the defect level of advanced-node and automotive parts.

## Characterizing a cell

Cell-aware flows start from the cell **layout**. Extraction finds plausible defect locations: opens in each net segment and shorts between neighbouring shapes or transistor terminals. Each candidate defect is injected into a transistor-level netlist and simulated with SPICE for every input pattern of the cell (and, for delay effects, every pattern pair). The result is a table: for each defect, the input patterns whose output differs from the fault-free cell. That table is a **user-defined fault model** (UDFM) attached to the cell type, and it is computed once per library, then reused for every instance in every design.

## Finding the detecting pattern

Given a defect's faulty output column, the detecting patterns are exactly the rows where it differs from the fault-free column. Take an AOI21 cell, $Y = \overline{AB + C}$, whose fault-free column for ABC = 000 to 111 is 1 0 1 0 1 0 0 0. If characterization reports a defective column 1 0 1 0 1 0 1 0, only row ABC = 110 differs, so 110 is the single detecting pattern. ATPG must then justify that pin combination at the cell's inputs and propagate the cell output to an observation point. A stuck-at test set needs only a few patterns per cell and may never apply 110, which is how this defect escapes.

## Sequential and delay behaviour

Some cell defects need two patterns. A transistor **stuck-open** in a CMOS gate disconnects one pull-up or pull-down path, so for one input pattern the output floats and keeps its previous charge: the cell becomes sequential. Detection needs a first pattern that drives the output to the opposite value and a second that exposes the floating condition. Resistive internal bridges may only slow the output and are recorded as at-speed two-pattern conditions. Adding cell-aware faults raises pattern count, since their extra conditions constrain more cells per pattern and compaction merges fewer faults; a defect whose only exposing combination is impossible in context, such as tied inputs, is classified untestable.

## Going deeper

Hapke et al., "Cell-aware test," IEEE Transactions on Computer-Aided Design, 2014; Wang, Wu and Wen (2006), chapter 2 on transistor-level faults.
