# Static CMOS gate structure

Every combinational standard cell in a digital library, from an inverter to an AOI222, is a static CMOS gate: an nMOS pull-down network and a dual pMOS pull-up network sharing one output. Synthesis maps logic onto these cells, extraction recognizes them in layout, and characterization measures them. A tool designer needs to go from a Boolean function to the two networks and back without hesitation, because cell generators, LVS and transistor-level verifiers all perform exactly this translation.

## The pull-down network from the function

A single-stage static CMOS gate computes an inverting function $F = \overline{E}$, where $E$ is built from the inputs with AND and OR. The pull-down network (PDN) must conduct exactly when $F = 0$, that is, when $E = 1$. Build it directly from $E$: every AND becomes nMOS devices in series, every OR becomes nMOS devices in parallel, and each literal is one transistor. For $F = \overline{AB + C}$ the PDN is (A series B) in parallel with C, which we write S(A, B) inside P(...): P(S(A, B), C). Because only un-complemented inputs appear, a function such as an AND or an XOR needs either an extra inverter stage or complemented inputs; a single static stage is always inverting.

## The pull-up network is the dual

The pull-up network (PUN) must conduct exactly when the PDN does not. By De Morgan, $\overline{E}$ in pMOS terms is obtained by exchanging series and parallel at every level of the PDN, keeping the same input on each device. For $F = \overline{AB + C}$, the PDN P(S(A, B), C) becomes the PUN S(P(A, B), C): A and B in parallel, in series with C. Check one vector: A = 1, B = 0, C = 0. The PDN has A on but B and C off, so it does not conduct; in the PUN, B's pMOS is on (gate low) and C's pMOS is on, so the output is pulled high, and indeed $F = \overline{0 + 0} = 1$. The duality must be applied at every level; swapping only the outermost connection is a common error.

## Why the structure is worth its cost

Static CMOS uses two transistors per literal, more than pass-transistor or dynamic styles, yet libraries are built almost entirely from it. The reasons are tool-friendly properties. The output is always driven by one network, so it is restored to a full rail regardless of input slew or noise. There is no static current between the rails apart from leakage. Every input sees only gate capacitance, so a cell's input is a fixed capacitance and its delay depends only on input slew and output load, which is exactly what a Liberty table can capture. And correctness is independent of timing: any input sequence gives the right output once inputs settle, so static timing analysis is sound. These properties are what let a library hide transistors behind arcs.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), sections 1.4 and 9.2; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), chapter 6.
