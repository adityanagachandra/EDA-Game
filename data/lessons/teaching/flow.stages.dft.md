# Design for test and scan insertion

A fabricated chip can have manufacturing defects that no design review could find, so each die is tested before it ships. Design for test (DFT) adds hardware that makes testing possible in seconds: scan chains, test compression, built-in self-test for memories and boundary scan for board-level tests. Scan insertion runs during or right after synthesis, and it changes the netlist that every later stage sees.

## Scan chains

A sequential circuit is hard to test because its internal state is hard to set and observe. Scan replaces each flip-flop with a scan flip-flop that has a multiplexer at its input. In test mode the flops are connected into shift registers, called scan chains. The tester shifts a pattern into the chains, pulses the functional clock once to capture the combinational logic's response, and shifts the result out. The sequential test problem becomes a combinational one, which automatic test pattern generation (ATPG) can solve. With $F$ scan flops split as evenly as possible into $C$ chains, the longest chain has $L = \lceil F/C \rceil$ flops.

## Test time

Test time is dominated by shifting. With $N$ patterns and chain length $L$, the first pattern is shifted in over $L$ cycles; then each pattern takes one capture cycle followed by $L$ cycles that shift its response out while the next pattern shifts in. The total is $(N + 1)L + N$ cycles. For $F = 20{,}000$ flops in $C = 16$ chains, $L = 1250$. With $N = 2000$ patterns, the test takes $2001 \cdot 1250 + 2000 = 2{,}503{,}250$ cycles, which at a 50 MHz shift clock is about 50.07 ms.

## What DFT costs the rest of the flow

Every scan flop is slightly larger and slower than a plain flop, and the scan connections add wires. Placement reorders the scan chains to follow cell locations, since a chain stitched in netlist order would zig-zag across the die. More chains shorten test time but need more tester pins, which is why compression logic expands a few external channels into many internal chains. Testers are expensive per second, so test time is product cost. DFT is therefore a flow-wide trade: area and timing on one side, test cost and defect coverage on the other.

## Going deeper

Bushnell and Agrawal, "Essentials of Electronic Testing for Digital, Memory and Mixed-Signal VLSI Circuits" (2000), chapters on scan design.
