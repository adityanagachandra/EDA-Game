# Core wrappers and IJTAG networks

A modern SoC holds dozens of reused cores and hundreds of embedded instruments: BIST engines, sensors, PLL controls, debug blocks. IEEE 1500 standardizes a wrapper around each core so its tests can be reused and its boundary isolated, and IEEE 1687 (IJTAG) standardizes a reconfigurable scan network that reaches instruments through the chip's JTAG port. Together with 1149.1 they compose chip, core and instrument access into one hierarchy.

## IEEE 1500 wrappers

A 1500 wrapper surrounds a core with a **wrapper boundary register** of cells on the core's terminals, a **wrapper instruction register** and a **wrapper bypass**, controlled through a serial **wrapper serial port** driven by on-chip control signals rather than package pins, so there is no mandatory TAP inside. An optional **wrapper parallel port** feeds many internal scan chains at once for fast core test. In **internal test** mode the wrapper isolates the core and applies its own patterns; in **external test** mode the wrapper cells drive and capture the core's terminals to test the interconnect between cores. The core's test information is described in CTL, the core test language, so the integrator can reuse the vendor's patterns.

## IEEE 1687 networks

An IJTAG network is a scan path whose topology changes on the fly. The key element is the **SIB**, segment insertion bit: a 1-bit scan cell that, when set, splices a sub-network into the active path, and when clear bypasses it. The active path length is the sum of the SIB bits plus the segments under open SIBs, recursively. A SIB that is closed contributes 1 bit; an open SIB over a 24-bit instrument register contributes 25; an open SIB over two open nested SIBs on 49-bit and 61-bit registers contributes $1 + (1 + 49) + (1 + 61) = 113$. Reaching an instrument behind closed nested SIBs takes one scan per level to open the path, then the data access.

## Description and retargeting

IJTAG separates structure from procedure. **ICL** (instrument connectivity language) describes the network and each instrument's interface; **PDL** (procedural description language) describes operations at the instrument's own terminals, such as writing a register and reading a status bit. A **retargeting** tool combines them: it computes the SIB settings and scan operations at the chip's TAP that carry out a PDL procedure through the current network. This makes instrument procedures portable across chips, the same way STIL makes patterns portable across testers.

## Going deeper

IEEE Std 1500, "Embedded Core Test"; IEEE Std 1687, "Access and Control of Instrumentation Embedded within a Semiconductor Device"; Marinissen and Zorian, papers on IEEE 1500 core test.
