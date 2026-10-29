# BSIM

BSIM, the Berkeley Short-channel IGFET Model family, is the most widely used set of industry-standard MOSFET compact models. Its members cover planar bulk devices (BSIM4), silicon on insulator (BSIM-SOI), multi-gate FinFETs and gate-all-around devices (BSIM-CMG) and independent-gate FD-SOI (BSIM-IMG). Foundry PDKs ship BSIM model cards, standard-cell characterization runs SPICE with them, and every timing number a design sees passes through one.

## The BSIM family

BSIM models are developed at UC Berkeley and standardized through the Compact Model Coalition, which qualifies models and releases so that one foundry model file runs consistently across simulators.

BSIM3 and its successor BSIM4 model planar bulk MOSFETs. BSIM4 added effects that became important in deep submicron nodes, such as gate tunneling current, stress and layout-dependent effects, and more detailed parasitics.

BSIM-SOI models partially and fully depleted SOI, including the floating body node and self-heating.

BSIM-CMG (common multi-gate) models devices whose gates are tied together around the channel: FinFETs and, in later versions, gate-all-around nanosheets. Its instance parameters include fin count, fin height and fin thickness.

BSIM-IMG (independent multi-gate) models devices with separately biased front and back gates, such as ultra-thin-body FD-SOI with back biasing.

Picking the model is the first decision in any extraction or simulation flow: equations, parameters and instance geometry all depend on it.

## Inside a BSIM model card

BSIM4 is threshold-voltage based: it builds a threshold expression with short-channel, narrow-width and body-bias terms, then a unified current expression across regions with smoothing functions. BSIM-CMG is surface-potential based: it solves the electrostatics of the multi-gate body for the surface potential and derives charge and current from it, adding corrections for quantum confinement in thin bodies (a threshold shift and a charge-centroid offset).

A model card holds hundreds of parameters covering threshold, mobility and its degradation, velocity saturation, DIBL, output conductance, parasitic resistance and capacitance, leakage components, noise and temperature dependence. Only a subset is tuned for a given process.

Foundry cards often use binning: parameters are fitted separately over ranges of $W$ and $L$, each with size-dependent terms of the form $P = P_0 + P_L/L + P_W/W + P_P/(WL)$, and the simulator chooses the bin by instance geometry. Bin boundaries must be continuous, or a sizing optimizer that changes width across a boundary sees a jump.

## BSIM in tools

Reference BSIM releases are distributed as Verilog-A source as well as being implemented inside commercial simulators. A Verilog-A compiler path lets a simulator load the exact reference model, with derivatives generated automatically, instead of hand-porting the equations into C, where implementations can diverge.

In open source, ngspice includes BSIM3 and BSIM4 natively and can load compiled Verilog-A models, including BSIM-CMG, through its OSDI interface, with OpenVAF as one compiler that produces OSDI binaries. This combination is what lets open-source flows simulate FinFET PDKs.

In a PDK flow, the BSIM card drives standard-cell characterization: SPICE runs per cell, arc, input slew, output load and PVT corner build the Liberty timing and power tables that static timing and power analysis consume.

Model versions are pinned exactly. Parameters are extracted against one version's equations, and running them under another version can shift currents and charges, so PDKs specify both the model and its version.

## Going deeper

Chauhan and others, "FinFET Modeling for IC Simulation and Design: Using the BSIM-CMG Standard" (2015); the BSIM4 and BSIM-CMG technical manuals from the BSIM group; the ngspice manual, OSDI chapter.
