# JTAG instructions and registers

What the TAP controller shifts depends on the instruction loaded into each device's instruction register. IEEE 1149.1 defines a small set of mandatory and optional instructions, each selecting one data register and deciding who controls the pins. Board test programs, in-system programmers and debuggers are sequences of these instructions, and getting the path lengths right is essential for any tool that drives a chain.

## The standard instructions

**BYPASS** selects the 1-bit bypass register, so a device not under test adds one bit to the chain. **SAMPLE/PRELOAD** selects the boundary register while the core keeps control of the pins: it can snapshot pin values in a running system, or preload safe values. **EXTEST** selects the boundary register and lets its cells drive the output pins and capture inputs, for interconnect test. **IDCODE** selects the 32-bit identification register, which holds manufacturer, part and version fields with bit 0 always 1. Optional instructions include **INTEST** (test the core through its boundary cells), **CLAMP** (drive pins from the boundary register while selecting bypass), and **HIGHZ** (float all outputs). Vendors add private instructions to reach internal scan, BIST and debug logic.

## Path lengths

An IR scan always passes through every device's instruction register, since the instructions are what select each device's data path; there is no bypass for the IR. A chain with IR lengths 4, 8 and 5 needs a 17-bit IR scan. A DR scan passes through whichever data register each device selected: with $n$ devices, one in IDCODE and the rest in BYPASS, the path is $32 + (n - 1)$ bits; with one in SAMPLE/PRELOAD on a 300-cell boundary register, $300 + (n - 1)$.

## Conventions that tools rely on

During Capture-IR every instruction register loads a fixed pattern whose two least significant bits are 01. Shifting out the IR chain therefore reveals each device's IR position and lets a host check chain integrity. After Test-Logic-Reset, a device loads IDCODE if it has one, otherwise BYPASS; a blind DR scan right after reset therefore reads a sequence of 32-bit codes (starting with a 1) and single 0 bits from bypassed devices, which is how JTAG tools discover an unknown chain. Before EXTEST, the boundary register must be preloaded with safe values, because the cells drive the pins as soon as the instruction takes effect at Update-IR.

## Going deeper

IEEE Std 1149.1; Parker, "The Boundary-Scan Handbook" (2016).
