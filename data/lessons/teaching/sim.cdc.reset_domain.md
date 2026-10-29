# Reset synchronizers and reset-domain crossing

Resets cross domains too. An asynchronous reset asserts without a clock, which is useful, but its release is a timing event relative to each flop's clock, and an asynchronous reset in one part of a design can corrupt flops elsewhere that are not reset. Reset synchronizers handle the first problem; reset-domain crossing (RDC) analysis handles the second. Both are invisible to zero-delay RTL simulation, which is why they need dedicated checks.

## Reset synchronizers

A flop with an asynchronous reset has **recovery** and **removal** requirements: the reset must be released some time before the clock edge (recovery) and not too soon after it (removal), the reset-pin equivalents of setup and hold. If an external reset is released at an arbitrary time, some flops may leave reset one cycle earlier than others, or go metastable. A **reset synchronizer** fixes this: two flops clocked by the local clock, both asynchronously reset by the external reset, with the first flop's data input tied to the inactive level. Reset asserts immediately and asynchronously, through the flops' reset pins, and deasserts synchronously, rippling through the two flops on clock edges. Each clock domain needs its own synchronizer, because release must be synchronous to that domain's clock.

## Reset-domain crossings

A **reset-domain crossing** is a path from a flop reset by one reset to a flop reset by a different reset, or not reset at all, even when both are on the same clock. When the source's reset asserts, its output changes immediately, at a time unrelated to the destination clock. If the destination keeps capturing, it can sample that change inside its setup-hold window and go metastable, or capture a value the protocol never allowed. Typical cases are soft resets of one block while the rest of the chip runs, and power or debug resets. The standard fixes are to isolate the destination during the source reset, for example by holding its enable or clock-gating it, or to order the resets so the destination is reset first. Static timing analysis does not normally time these asynchronous paths against the destination clock, so dedicated RDC tools check them structurally.

## Verifying reset behavior

Reset bugs escape simulation in specific ways. Zero-delay simulation releases reset cleanly at one instant, so recovery and removal problems never appear. Two-state simulation starts every flop at a defined value, hiding flops that lack a reset. Four-state simulation helps with the second: flops that are not reset start at X, and X-propagation shows logic that depends on them before initialization, as long as X-optimistic constructs do not mask it. Verification of resets therefore combines four-state reset tests, structural reset and RDC checks, assertions on reset ordering (for example that one domain leaves reset only after another), and gate-level simulation with timing for the release sequence of critical blocks. Partial and soft resets need special scenarios: reset one block while traffic runs elsewhere and check that nothing outside is corrupted.

## Going deeper

Cummings, Mills and Golson, "Asynchronous and Synchronous Reset Design Techniques, Part Deux" (SNUG 2003); Ginosar, "Metastability and synchronizers: a tutorial", IEEE Design and Test of Computers (2011).
