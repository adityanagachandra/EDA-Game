# Why analog automation lags digital

Digital implementation runs from RTL to GDS with little human intervention. Analog design still relies on experts who size transistors by hand, draw matched layouts and iterate with SPICE. For a tool designer, the reasons are not historical accident: they come from what analog circuits compute, how their quality is measured and how expensive it is to evaluate a candidate. Understanding the gap tells you which automation problems are tractable now and which are still research.

## No abstraction that hides the devices

Digital automation rests on a ladder of abstractions. Boolean logic hides voltages, standard cells hide transistors, and a Liberty table hides a cell's internals behind slew and load. Each gate restores its output to a clean logic level, so small errors in one stage do not propagate, and stages can be designed and checked independently. Analog circuits have no such ladder. An amplifier's gain, bandwidth, noise and offset depend continuously on every width, length and bias current, and on the parasitics of the wires connecting them. There is no small library of context-free building blocks: a current mirror behaves differently depending on its output voltage, its load and its layout. That is why analog synthesis cannot reduce to choosing and connecting cells, and why each block is in practice designed for its context.

## Many specs, expensive evaluation

A digital block's quality reduces to a few scalars: worst slack, power, area. An analog block is judged on a vector of interacting specifications: DC gain, gain-bandwidth product, phase margin, slew rate, noise, offset, output swing, power, often across a dozen PVT corners and under mismatch. Improving one usually costs another, so sizing is constrained multi-objective optimization. Each candidate is evaluated by a nonlinear SPICE simulation. One run takes seconds to minutes, but corners times Monte Carlo samples times optimizer iterations multiplies it: 50 corners and 200 samples is 10,000 simulations for one candidate. The design space is also hostile: performance is nonconvex in the sizes, devices can leave saturation abruptly, and feasible regions are thin. An optimizer can afford only thousands of evaluations in a space where naive search would need millions.

## Where automation works today

The mature parts of the analog flow are verification engines: SPICE simulation (ngspice, Xyce and commercial simulators), corner and Monte Carlo analysis, DRC, LVS and parasitic extraction. Synthesis-like tasks are partially automated. Sizing is helped by gm/Id lookup tables, equation-based optimizers and simulation-in-the-loop methods such as Bayesian optimization. Layout is helped by constraint-driven generators such as MAGICAL and ALIGN, which recognize matched structures and place them symmetrically. Topology selection and specification remain human decisions. The most automatable problems share a property: the topology is fixed and the targets are known, as in re-sizing across corners or porting a block to a new process. A tool designer should expect value from automating search and verification inside a designer's intent, not from replacing that intent.

## Going deeper

Gielen and Rutenbar, "Computer-aided design of analog and mixed-signal integrated circuits" (Proceedings of the IEEE, 2000); Razavi, "Design of Analog CMOS Integrated Circuits" (2nd ed., 2017), chapter 1.
