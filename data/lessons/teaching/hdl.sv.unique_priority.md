# unique and priority case

SystemVerilog's unique and priority modifiers on case and if statements state properties of the decision logic: that items never overlap, or that some item always matches. Synthesis uses the properties to build cheaper logic, and simulation checks them at run time. They replace the older full_case and parallel_case pragmas, which spoke only to synthesis and created silent mismatches.

## unique: mutually exclusive items

unique case asserts that exactly one item matches the selector, unless a default handles the remainder. Simulation reports a violation if two items match (overlap) or if none matches and there is no default. In casez with items 2'b1? and 2'b?1, the selector 2'b11 matches both, so a unique violation is reported, although the first match still executes.

The payoff in synthesis is structure. Mutually exclusive items let the tool OR the selected values together in a parallel multiplexer instead of building a priority chain, which is shallower and smaller. unique0 keeps the no-overlap check but allows no item to match. A default item satisfies the completeness check, but overlapping explicit items are still violations.

## priority: ordered but complete

priority case keeps the ordinary first-match-wins semantics, so overlap is allowed, and adds a promise that at least one item matches. A priority case with no default that sees an unmatched selector reports a violation in simulation.

Completeness is exactly what a latch would otherwise cover. Without the promise, a combinational case that leaves some selector values unassigned infers a latch; with priority, synthesis may treat the unlisted values as unreachable and build plain logic. The same keywords apply to if chains: unique if promises that the conditions are mutually exclusive, and priority if promises that some branch is taken.

## Why checked promises beat pragmas

A full_case or parallel_case pragma tells synthesis the same things, but simulation never sees it. If the promise is false, simulation holds an old value or takes the first match while the netlist does something else, and nothing reports the difference. unique and priority make both tools act on the same assumption and make simulation check it, so a broken promise appears as a violation message.

Run-time checks only cover the stimulus applied, though. A design that never violates unique in simulation can still fail in silicon on an input sequence that was never simulated. Proving the promise for all inputs needs formal tools, which can check the same properties as assertions. A common safe use is the decode of a one-hot state register, whose items are mutually exclusive by construction.

## Going deeper

IEEE Std 1800-2017, clauses 12.4.2 and 12.5.3; Sutherland, "SystemVerilog Saves the Day: the Evil Twins of Verilog Synthesis" (SNUG 2005).
