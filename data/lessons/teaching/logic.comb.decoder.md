# Decoders and address decoding

A decoder turns an $n$-bit binary code into $2^n$ one-hot lines: exactly one output is asserted, the one whose index equals the input. Decoders select memory rows, enable registers in a register file, steer bus transactions to peripherals, and generate minterms. Every SRAM macro contains a row decoder, and its depth and size scale directly with the address width.

## One output per code

An $n$-to-$2^n$ decoder asserts output $Y_i$ when the input $A_{n-1} \dots A_0$ equals $i$, and only then (when an enable input, if present, is active). Each output is one minterm: $Y_5$ of a 3-to-8 decoder is $A_2 \bar{A_1} A_0$. A decoder is therefore an AND-plane that generates every minterm, and any function of the inputs is an OR of selected outputs.

With input $A_3 A_2 A_1 A_0 = 1101$, a 4-to-16 decoder asserts $Y_{13}$. With active-low outputs (common, because NAND gates are faster than AND in CMOS) the same input drives $\bar{Y}_{13}$ low and leaves the other fifteen high.

## Building large decoders from small ones

A wide decoder is built as a tree. To make an $n$-to-$2^n$ decoder from $k$-to-$2^k$ decoders with enables, the top $k$ address bits drive one decoder whose outputs enable $2^k$ decoders at the next level, and so on. For $n$ a multiple of $k$ the count is

$$1 + 2^k + 2^{2k} + \dots + 2^{n-k}.$$

A 4-to-16 decoder from 2-to-4 decoders needs $1 + 4 = 5$; a 6-to-64 needs $1 + 4 + 16 = 21$ of them. Memory row decoders use a related idea, predecoding: split the address into groups of 2 or 3 bits, decode each group into one-hot lines, then combine one line from each group with a small AND per row. This keeps each final gate's fan-in small, which keeps it fast.

## Address decoding and partial decoding

In a memory map, a device is selected when the upper address bits match its base pattern; the lower bits pick a location inside the device. With a 16-bit address bus and a device selected by $A_{15} A_{14} A_{13} A_{12} = 1010$, the device occupies $2^{12} = 4096$ consecutive addresses, from 0xA000 to 0xAFFF.

A cheaper partial decoder ignores some address bits. If the select logic looks only at $A_{15}A_{14} = 10$ for a device that uses only 12 address bits internally, the two ignored bits $A_{13}A_{12}$ make the device appear at $2^2 = 4$ aliases, and its footprint grows to $2^{14}$ addresses. Aliasing is harmless until another device is mapped into the aliased range; a verification tool checking an address map for overlaps must model exactly which bits each decoder compares.

## Going deeper

Mano and Ciletti, "Digital Design" (5th ed., 2013), chapter 4; Weste and Harris, "CMOS VLSI Design" (4th ed., 2010), chapter 12 on SRAM row decoders.
