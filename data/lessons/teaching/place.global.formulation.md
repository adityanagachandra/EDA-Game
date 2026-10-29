# The global placement problem

Global placement decides roughly where every standard cell goes. Stated exactly it is a huge, NP-hard, discrete problem; every practical placer solves a relaxation of it. Knowing the decision variables, the objective, the density constraint and the relaxations each placer family makes is the map a tool designer needs before reading any particular algorithm.

## Variables, objective and constraint

The decision variables are the coordinates $(x_i, y_i)$ of every movable cell; macros and I/O pins are usually fixed. The objective is total weighted wirelength, $\sum_e w_e \,\mathrm{WL}_e(x, y)$, with HPWL as the usual measure. The constraint is density: for every bin $b$ of a grid over the placement region, the movable cell area inside it must not exceed a target fraction $\rho_t$ of the bin's free area. Without the constraint, minimising wirelength collapses connected cells onto one point. Legality (rows, sites, no overlap) is not imposed here; it is deferred to legalization. A target density below 1, for example 0.7, also leaves room for later buffering and routing.

## Three families of relaxation

Min-cut placers relax geometry into partitioning: recursively bisect the netlist and the region, minimising nets cut at each level. Quadratic placers replace wirelength by a sum of squared two-pin distances, whose minimum solves a sparse linear system, and handle density by separate spreading steps. Nonlinear analytical placers smooth HPWL with log-sum-exp or weighted-average functions and add a differentiable density penalty, $\min_x W(x) + \lambda D(x)$, solved by gradient methods. Electrostatics-based placers such as ePlace, RePlAce and DREAMPlace belong to the last family, with density modeled as electric potential energy. Each family trades accuracy of the wirelength model, quality of the spreading model and speed.

## Stopping and handing off

Global placement ends before density is perfectly satisfied. Placers track overflow, the fraction of cell area that sits above target in over-full bins, and stop around a small value such as 0.1, because squeezing out the last overflow costs wirelength that legalization can recover more cheaply with small moves. The handoff has failure modes. A low average overflow can hide local hot spots that force long legalization moves. Fixed pins and macros are essential: in quadratic placement they anchor the linear system, and without them it is singular. And because the objective is wirelength, timing and routability enter only through net weights, density targets and inflation, which must be designed in, not hoped for.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 4; Markov, Hu and Kim, "Progress and challenges in VLSI placement research" (Proceedings of the IEEE, 2015).
