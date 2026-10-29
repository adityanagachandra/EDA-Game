# Finite element analysis of packages

Finite element analysis (FEA) solves the equations of elasticity, and often plasticity and creep, over a meshed model of the package. It is the tool that predicts warpage, joint strain, die stress and delamination risk where closed-form models fail: at edges, corners, non-uniform layers and nonlinear materials. For a tool designer the important questions are what the model includes, how it is meshed and loaded, and how its answers can mislead.

## Modelling choices

A package FEA model starts with geometry, materials and loads. Geometry is simplified: thousands of solder joints are often replaced by a homogenized layer except in a submodel of the critical corner joints. Materials need temperature-dependent moduli and CTEs, plus nonlinear behaviour: solder creeps at room temperature, mold compound and underfill are viscoelastic and change sharply at their glass transition temperature. The load is a temperature history starting from the stress-free state, not a single temperature. Symmetry is used to model a quarter or an eighth of the package when geometry and loads allow, cutting the cost. Each simplification is a choice that must be justified against what the analysis must predict.

## Mesh, elements and convergence

FEA discretizes displacement with elements: bricks or tetrahedra in 3-D, shells for thin layers. Thin layers such as die-attach films and RDL dielectrics are hard to mesh because elements must stay well shaped while being very thin, so analysts use layered shell or solid-shell elements, or many thin bricks. Accuracy is checked by refining the mesh and confirming that the quantity of interest converges. At sharp re-entrant corners between bonded materials, linear elasticity predicts a stress singularity, so the peak stress keeps rising with refinement and never converges. Analysts then compare quantities that do converge, such as stress at a fixed distance or an energy release rate, rather than the peak value.

## Validation and use in the flow

FEA results must be validated against measurements: warpage across the reflow profile from shadow moiré, strain from digital image correlation, and cycles to failure from thermal cycling. A model calibrated this way can then explore designs faster than building them. In an EDA flow, FEA is the accurate engine: chip-package interaction checks use it to map stress from bumps into the die's back-end stack, and warpage simulation guides substrate stackups. Because full models are slow, tools build reduced models or surrogates from FEA runs, and the main risk is using a surrogate outside the range of designs it was fitted on.

## Going deeper

Zienkiewicz, Taylor and Zhu, "The Finite Element Method: Its Basis and Fundamentals" (7th ed., 2013); Lau (ed.), "Thermal Stress and Strain in Microelectronics Packaging" (1993).
