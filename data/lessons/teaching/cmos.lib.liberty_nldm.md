# Liberty NLDM timing tables

The nonlinear delay model (NLDM) is the classic Liberty timing model: for each arc and output edge, a two-dimensional table of delay and a table of output transition, indexed by input transition and output load. It is the measured, nonlinear replacement for the RC and logical-effort models, and the core lookup every static timer performs millions of times per run.

## Two axes, four tables

An NLDM arc carries cell_rise and cell_fall tables for delay and rise_transition and fall_transition tables for output slew. A lu_table_template names the axes; in the common convention index_1 is the input net transition and index_2 the total output net capacitance. Each table row corresponds to one input transition and each column to one load. Example: index_1 = 10, 30, 90 ps and index_2 = 1, 4, 16 fF, with the row for 30 ps reading 22, 35, 86. At an input transition of 30 ps and a load of 4 fF, the delay is 35 ps. At exactly a grid point, the lookup is just reading the entry; between grid points a timer interpolates, which is taught with timing analysis.

## Choosing the table from the edges

The tables are named for the output edge. For a negative_unate arc, such as an inverter, a rising input produces a falling output, so the delay is read from cell_fall and the output slew from fall_transition. For a positive_unate arc, such as a buffer, a rising input gives a rising output and uses cell_rise. For a non_unate arc both are possible, so the timer computes both. Liberty also distinguishes the transition on the input, which only selects the row, from the transition on the output, which comes from the transition table and is passed on. Mixing up these two, or reading the input edge's table instead of the output edge's, is the most common bug in hand-written lookups.

## Chaining lookups along a path

Delay calculation along a path chains these lookups. Stage 1's input slew and load give its delay and its output transition; that transition becomes stage 2's input slew, which together with stage 2's load gives stage 2's delay, and so on. If inverter U1 receives a rising input of 30 ps transition and drives 4 fF, its cell_fall and fall_transition at (30 ps, 4 fF) might be 35 ps and 28 ps; then U2, an inverter whose output rises, uses its cell_rise at an input transition of 28 ps and its own load. The load each stage sees is the sum of its wire capacitance and the input pin capacitances of its receivers, taken from Liberty. Every delay therefore depends on upstream slews and downstream loads, which is what makes incremental timing non-trivial.

## Going deeper

The Liberty Reference Manual, timing groups and lookup tables; Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 3.
