# case, casez and casex matching

Verilog has three case statements that differ only in how they treat x and z bits: case, casez and casex. The differences decide which item executes when a selector carries unknown bits, so they determine whether simulation hides a bug or exposes it. Synthesis reads all three as multiplexers, so a tool designer must know exactly where simulation and hardware diverge.

## casez and wildcard items

All three forms try the items in order and execute the first one that matches; later matches are ignored, which makes a case statement a priority structure unless the items are mutually exclusive. In casez, a z bit, which may be written as ?, is a wildcard: it matches anything in that position.

Take sel = 4'b1100 and the items 4'b??0? (y = 1), 4'b1?0? (y = 2), 4'b001? (y = 3). The first item checks only bit 1, which is 0, so it matches and y = 1 even though the second item also matches. Casez with ? is the standard way to write priority decoders and instruction decoders with don't-care fields.

## Unknown bits in the selector

casez also treats z in the selector as a wildcard. If sel = 4'b000z, an item 4'b0000 matches, because the z lane matches anything. That is surprising, since the selector really is unknown in that lane. Plain case is stricter: it compares like ===, so every bit including x and z must match literally. In plain case, ? is simply another spelling of z and matches only a z on the selector.

The consequence is that a selector driven by an undriven or floating net may silently match a casez item in simulation, while plain case falls through to the default. Choosing case for decoders whose selector should never be unknown keeps that failure visible.

## casex and x-masking

casex extends the wildcard to x as well as z, in both the item and the selector. With sel = 4'b100x, the item 4'b1xx0 matches because its middle lanes are wildcards and the selector's last lane is x, which also matches anything. Worse, a selector that is entirely x, such as an uninitialized register, matches the first item, so simulation proceeds as if a valid value were present.

This is x-optimism: the simulator chooses a definite branch where real hardware might do anything. Coding guidelines therefore discourage casex in RTL and prefer casez with explicit ? wildcards, or SystemVerilog's case inside, which treats only the item's x and z (and ?) as wildcards and never the selector's. For synthesis all forms become the same multiplexer logic, so the choice affects only how honestly simulation reports unknowns.

## Going deeper

IEEE Std 1800-2017, clauses 12.5.1 and 12.5.4; Mills, "Yet Another Latch and Gotchas Paper" and Sutherland, "Gotcha Again" (SNUG).
