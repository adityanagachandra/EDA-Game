# Mead-Conway, SPICE and the Berkeley CAD tools

Modern EDA grew from a few landmark ideas of the 1970s and 1980s. Carver Mead and Lynn Conway made chip design teachable as structured design with simple rules. SPICE at Berkeley made circuit simulation a shared tool. Berkeley's logic synthesis and layout programs, led by researchers such as Robert Brayton, Alberto Sangiovanni-Vincentelli and Richard Newton, became the basis of commercial tools and of today's open-source ones.

## Mead and Conway

Carver Mead at Caltech and Lynn Conway at Xerox PARC wrote "Introduction to VLSI Systems" (1980), which taught system designers, not only device experts, to design chips. Its key ideas were scalable design rules expressed in units of a parameter $\lambda$, so a layout could move between processes, and a structured, hierarchical design style. Conway taught a course at MIT in 1978 whose students' designs were fabricated together on one wafer. That multi-project-wafer idea grew into MOSIS, a service that gave universities affordable fabrication and seeded a generation of designers and CAD researchers.

## SPICE at Berkeley

SPICE (Simulation Program with Integrated Circuit Emphasis) was developed at UC Berkeley in the early 1970s by Laurence Nagel under Donald Pederson, who is often called the father of SPICE. Berkeley released it freely, and SPICE2 became the reference circuit simulator in universities and industry. Its algorithms, modified nodal analysis, sparse matrix solution, Newton-Raphson iteration for nonlinear devices and implicit integration for time steps, are still the core of commercial simulators and of ngspice and Xyce. Richard Newton, later dean of engineering at Berkeley, worked on circuit simulation in the same group before turning to layout and design methodology.

## Berkeley logic synthesis and layout tools

Berkeley's CAD group produced a sequence of influential programs. Espresso, by Brayton and colleagues, minimized two-level logic. MIS and then SIS performed multi-level logic synthesis, and ABC, from Alan Mishchenko and Brayton, continues that line with AIG-based optimization. Magic, by John Ousterhout and students, was an interactive layout editor with continuous design-rule checking. Students and faculty from this group founded or shaped major EDA companies: Sangiovanni-Vincentelli was a co-founder of both Cadence and Synopsys, and many core commercial engines descend from Berkeley ideas. Because many of these programs were distributed with source code, they trained generations of students and became the reference implementations that commercial tools were measured against.

## Going deeper

Mead and Conway, "Introduction to VLSI Systems" (1980); Nagel and Pederson, "SPICE (Simulation Program with Integrated Circuit Emphasis)", Berkeley memorandum (1973); Brayton, Hachtel, McMullen and Sangiovanni-Vincentelli, "Logic Minimization Algorithms for VLSI Synthesis" (1984).
