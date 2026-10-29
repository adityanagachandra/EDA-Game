# Manufacturing, lithography and DFM

Manufacturing is the branch where geometry becomes silicon, and its tools work on the far side of tape-out: mask synthesis, lithography simulation, process modelling and yield analysis. Design for manufacturability (DFM) pushes some of that knowledge back into the design flow. As features shrank below the wavelength of the light that prints them, these tools became as essential as placement and routing.

## Printing below the wavelength

Leading processes print features much smaller than the wavelength of the light used, 193 nm for deep ultraviolet and 13.5 nm for extreme ultraviolet (EUV) lithography. Diffraction blurs and distorts drawn shapes: corners round, line ends pull back, and dense and isolated lines print at different widths. Resolution enhancement techniques counter this. Optical proximity correction (OPC) adds serifs and biases to each polygon so the printed shape matches the intended one. Sub-resolution assist features help isolated lines print. Inverse lithography computes mask shapes by optimization. Multiple patterning splits one layer onto two or more masks when pitches are too tight for a single exposure, which turns layout into a graph-colouring problem.

## Design for manufacturability

DFM moves manufacturing awareness upstream. Design rules are the first form: minimum width and spacing, density rules that require metal coverage within a range so that chemical-mechanical polishing leaves a flat surface. Fill insertion adds dummy metal to meet them. Recommended rules and litho-friendly layout patterns improve yield beyond the minimum rules. Critical-area analysis computes how much of a layout is vulnerable to a random particle of a given size, which feeds yield models. Redundant via insertion replaces single vias with pairs where space allows, so that one failed via does not break a net.

## Process and device modelling

Technology computer-aided design (TCAD) simulates fabrication processes and device physics: how implants diffuse, how a transistor's current depends on its structure. It is used to develop processes and to generate the compact models in a PDK. Design-technology co-optimization (DTCO) evaluates process options by their effect on standard-cell area, power and performance, so process and library are developed together. Manufacturing data flows back as yield learning: failing locations on wafers are correlated with layout patterns, and the patterns become new rules or OPC fixes. Each of these loops shortens the time between a new process and its first high-yield products, which is why manufacturing tools sit on the same critical path as design tools.

## Going deeper

Wong, "Resolution Enhancement Techniques in Optical Lithography" (2001); Mack, "Fundamental Principles of Optical Lithography" (2007).
