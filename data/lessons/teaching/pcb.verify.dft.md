# Design for test on boards

A built board must be tested for assembly defects: opens, shorts, missing or wrong parts. Testers reach nets through physical test points, probed by a bed-of-nails fixture or flying probes, and through boundary scan (IEEE 1149.1, JTAG), which drives and senses device pins from inside the chips. Design for test plans that access during layout. Its quantities are coverage, the fraction of nets testable, and the length of the JTAG scan chain, both straightforward counts.

## Test access and coverage

In-circuit test (ICT) presses spring probes onto **test points**, pads of a minimum size and spacing on one side of the board. Flying-probe testers move a few probes over the board, needing no fixture but running slower. Either way, a net is testable only if a probe can reach it or a device can drive and sense it. **Coverage** is the testable fraction: with 1,000 nets, 600 with test points and 250 more reachable only through boundary-scan cells, coverage is $(600 + 250)/1000 = 85\%$. BGAs hide most of their nets, so boundary scan is often the only access. Adding test points costs area and capacitance on fast nets, so DFT trades coverage against density and signal integrity, net by net.

## The JTAG chain and instruction registers

Boundary-scan devices are linked in a **chain**: TDI enters the first device, each device's TDO feeds the next, and TCK and TMS are shared. Each device has an **instruction register** (IR) and, among its data registers, a **boundary-scan register** with one cell per I/O. Loading instructions shifts through all IRs in series, so the Shift-IR length is the sum of the IR lengths: devices with 4, 6 and 10 bit IRs need 20 bits. The first bring-up test reads every device's IDCODE, which proves the chain is wired and powered. The board netlist plus each device's BSDL file (Boundary Scan Description Language) tells test software which pins each cell drives and senses, so it can generate interconnect tests automatically.

## Bypass and test pattern length

To test through one device, the others are put in **BYPASS**, which selects a 1-bit bypass register in each, keeping the chain short. The data-register path is then the target's boundary-scan register plus one bit per bypassed device. Example: a 4-device chain whose FPGA has an 800-cell boundary register needs $800 + 3 = 803$ TCK cycles in Shift-DR per pattern. Interconnect tests (EXTEST) drive patterns from one device's output cells and capture them at another's input cells, detecting opens and shorts on nets between scan devices; a short pattern set such as a counting sequence distinguishes every net. Chain order, TCK routing and reset (TRST) are layout decisions, so DFT must be planned before the board is routed.

## Going deeper

IEEE 1149.1, "Standard for Test Access Port and Boundary-Scan Architecture". Parker, "The Boundary-Scan Handbook" (4th ed., 2016).
