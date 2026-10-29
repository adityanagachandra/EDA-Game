# Tape-out, fabrication, packaging and test

Tape-out is the moment the design team releases final layout to the foundry. After it, the design leaves the EDA world for the factory: mask making, wafer fabrication, wafer test, dicing, packaging and final test. Months pass between tape-out and working parts in hand, and any bug that reaches this point costs a respin. Knowing what happens after tape-out explains why signoff is so conservative.

## Tape-out and mask data preparation

At tape-out the team delivers GDSII or OASIS layout plus documentation. The foundry and the mask shop then run mask data preparation: fill insertion if not already done, optical proximity correction (OPC), which pre-distorts shapes so that they print correctly through lithography, and fracturing of the polygons into the form a mask writer exposes. One mask is made per lithography layer (and more under multiple patterning), so an advanced process needs a set of dozens of masks, which is a large share of a respin's cost. Mask data preparation is itself a large compute job.

## Fabrication

Wafer fabrication builds the transistors first (front-end-of-line) and then the metal interconnect stack (back-end-of-line) through hundreds of process steps: lithography, etching, deposition, implantation, polishing. Fabrication takes weeks to months. The foundry monitors the process with test structures in the scribe lines between dies. Engineering lots, the first wafers of a new design, may be run faster at higher cost to get silicon for debug sooner. Process variation between wafers and across a wafer is why silicon is tested and binned rather than assumed identical. For the design team, the fabrication wait is when test programs, bring-up boards and software are prepared, so that first silicon can be debugged immediately.

## Packaging and production test

Wafer sort probes every die on the wafer and marks failing ones. Good dies are diced, assembled into packages (wire-bond or flip-chip, possibly with other dies in an advanced package) and tested again in final test, often across temperature. Production test runs the scan patterns and memory built-in self-test that DFT made possible, plus parametric checks and speed binning, in which parts are sorted by the frequency they reach. Test escapes, defective parts that pass test, are counted as defective parts per million shipped. When first silicon fails, engineers debug it with scan diagnosis and physical failure analysis, then decide between a metal-only fix, which changes a few masks, and a full respin.

## Going deeper

Plummer, Deal and Griffin, "Silicon VLSI Technology: Fundamentals, Practice and Modeling" (2000); Bushnell and Agrawal, "Essentials of Electronic Testing".
