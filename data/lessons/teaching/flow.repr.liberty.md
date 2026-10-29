# Liberty cell libraries

Liberty is the format that tells synthesis, timing and power tools how each library cell behaves. It replaces a cell's transistors with tables: delay and output slew as functions of input slew and output load, power per transition, leakage, pin capacitances, logic function and timing checks for sequential cells. It is the abstraction that makes static timing analysis of a billion-transistor chip possible, and its accuracy limits the accuracy of everything that reads it.

## What a Liberty file contains

A Liberty library declares units (for example time in ns, capacitance in pF), the operating condition it was characterized at (one process, voltage and temperature), and then one cell group per cell. Each cell lists pins with direction and capacitance, a function string such as `!(A & B)` for a NAND2, and timing arcs from each input to each output. Each arc holds lookup tables, the nonlinear delay model (NLDM), indexed by input transition and output load. Sequential cells add setup, hold and clock-to-output arcs. A separate Liberty file exists for each characterized corner.

## Reading and interpolating a table

A timing table is a grid. Suppose a NAND2's rise delay row at a given input slew has loads 0.01, 0.02, 0.04 pF with delays 0.050, 0.068, 0.101 ns. At a load of exactly 0.02 pF the delay is the table value, 0.068 ns. At 0.03 pF the tool interpolates linearly between the neighbouring entries:
$$d = 0.068 + \frac{0.03 - 0.02}{0.04 - 0.02}(0.101 - 0.068) = 0.0845\ \text{ns}.$$
With both slew and load between grid points the tool interpolates bilinearly. Outside the table it extrapolates, which is less accurate and which timers flag.

## Power, leakage and what tables miss

Liberty also holds internal switching energy per arc and leakage power per cell, sometimes per input state. Total leakage of a block is the sum over instances: 4000 cells at 2.5 nW and 1500 at 6 nW leak $10{,}000 + 9{,}000 = 19{,}000$ nW, or 19 $\mu$W. Tables have limits. NLDM models the output as a ramp, which is poor for long resistive wires and very small geometries, so current-source models (CCS, ECSM) store output current waveforms instead. Statistical extensions (the Liberty variation format) add per-arc sigma for on-chip variation. A library missing a corner, or with coarse tables, propagates its error into every timing and power number downstream.

## Going deeper

The Liberty user guide and reference manual (Synopsys, distributed through its open library program); the OpenSTA documentation on Liberty support.
