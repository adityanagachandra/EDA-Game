# The signoff checklist

Signoff is the set of checks a design must pass before tape-out, each run with the qualified tools at the worst conditions for what it checks. No single check covers another's failure modes: a design can meet timing and fail LVS, or pass DRC and fail in the field from electromigration. A tool designer working on any engine needs to know where its output lands in this checklist and which other checks it can silently break.

## The categories

Digital signoff covers five groups. Timing: setup and hold, with signal integrity, across all modes and corners, using signoff extraction. Power: total and per-domain power, and power-intent checks for multi-voltage designs. Reliability: electromigration on power and signal nets, and IR drop, static and dynamic. Physical verification: design-rule checking (DRC) against the foundry deck, layout versus schematic (LVS) to prove the layout implements the netlist, antenna checks for charge damage during processing, and density checks for polishing uniformity. Function: logical equivalence checking (LEC) between the final netlist and the reference, after all optimization and ECOs. Each check has its own tool and its own definition of clean, and tape-out requires all of them.

## Matching checks to conditions

Each check runs where it is hardest to pass. Setup at slow transistors, low voltage and the Cworst or RCworst wiring; hold at fast transistors, high voltage and Cbest or RCbest. Crosstalk needs coupling capacitors in the parasitics, not just totals. Timing with IR-drop awareness needs per-instance voltages from dynamic IR analysis. EM is worst at high temperature. Density matters because chemical-mechanical polishing dishes sparse regions and erodes dense ones, so fill is inserted, and that fill must be present in the final extraction because it adds capacitance. A block clean at the typical corner only has passed none of these; a missing corner is the most common signoff gap.

## The checklist after changes

Every late change re-opens the checklist. An ECO changes netlist and geometry, so timing, SI, EM/IR, DRC, LVS and LEC must all be re-run on its effects, typically incrementally. Metal fill after extraction invalidates parasitics. A change in temperature range or voltage can expose temperature inversion, where at low supply voltage cells are slowest when cold, so the slow corner moves. Signoff must also be run with the qualified extractor rather than the router's internal estimates, because only the former is calibrated to the foundry's golden solver. A disciplined flow records which version of every input each check used, so that the final signoff is provably consistent: the same netlist, the same layout and the same parasitics everywhere.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), the chapters on timing closure and physical verification; the OpenROAD flow scripts' final-report stage, which runs the open-source equivalents of these checks.
