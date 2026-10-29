# Process corners

A process corner is a set of device models with every global parameter moved to a chosen statistical extreme. Corners are the oldest and still the most used way digital and analog flows account for variation. A tool designer must understand what a corner represents, what it leaves out, and why the number of corners in a modern signoff run explodes.

## What a corner represents

Foundries deliver models at named corners such as TT, SS, FF, SF and FS. The first letter describes the NMOS and the second the PMOS: SS is slow N and slow P, SF is slow N with fast P. Each corner shifts the global parameters, threshold voltage, mobility, oxide thickness, together to roughly a 3-sigma point of their joint die-to-die distribution. Every transistor on the die moves identically. Corners therefore model global variation, and only global variation: within-die mismatch, which makes neighbouring devices differ, is invisible to them. Skewed corners SF and FS matter for ratioed circuits, such as SRAM cells, level shifters and sense amplifiers, whose behaviour depends on the relative strength of N and P devices.

## Choosing and counting corners

A signoff corner combines process, supply voltage, temperature and interconnect: for example SS at low voltage and high temperature with an RC-worst extraction. Setup is usually checked at slow corners and hold at fast ones. The combinations multiply: 5 process corners, 3 voltages, 3 temperatures and 5 interconnect corners give 225 analysis points, before multiple operating modes are added. Two effects prevent pruning by intuition. Temperature inversion makes low-voltage delay at advanced nodes worst at low temperature, so the worst temperature depends on voltage. And interconnect corners trade resistance against capacitance, so RC-worst suits long resistive nets while C-worst suits short ones. Flows prune to a dominant subset by analysis, not by assumption.

## Where corners mislead

Corners stack worst cases. Setting several independent parameters at 3 sigma at once produces a point far rarer than 3 sigma of the actual performance, which over-margins the design. At the same time corners can be optimistic: they ignore local variation, so a design clean at every corner can fail on paths whose devices vary against each other, which is why flows add on-chip variation derates or statistical timing. A corner built to be worst for drive current need not be worst for leakage, offset or noise margin. Corners are the right tool when global variation dominates and the metric is monotone in the corner parameters; elsewhere they are a starting point for statistical analysis, not a substitute.

## Going deeper

Weste and Harris, "CMOS VLSI Design: A Circuits and Systems Perspective" (4th ed.), the section on process corners and variability; Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), the timing closure chapter.
