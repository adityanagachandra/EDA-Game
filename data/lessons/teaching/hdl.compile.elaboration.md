# Elaboration

Elaboration turns a set of parsed modules into one concrete design. It picks the top, links every instance to its module, evaluates parameters and widths per instance, unrolls generate constructs and builds the instance tree. It is the phase where HDL compilers differ most and where SystemVerilog is hardest to implement, which is why shared elaborators such as slang and Surelog exist.

## What elaboration builds

Parsing checks each file against the grammar on its own. Elaboration starts from the top module and builds the design: it resolves each instance to its module definition, binds parameter values, evaluates every constant expression in declarations (so a port declared [W-1:0] gets a number), unrolls generate loops and selects generate-if branches. The result is a tree of instances in which nothing is symbolic any more.

Only what is reachable from the top is elaborated. A module that is parsed but never instantiated under the chosen top is not part of the design, although some tools treat every uninstantiated module as an extra top, which is why flows name the top explicitly. The elaborated design is what simulation and synthesis consume.

## Parameters resolve top down

A child's parameters can depend on values its parent passes in, so elaboration resolves the hierarchy top down: the parent's values are known before the child is specialized. Two instances of fifo with W = 8 and W = 16 become two specializations of fifo, one per distinct parameter set; tools name them separately (Yosys uses a paramod prefix). A generate-if on a parameter is decided per instance, so only one implementation exists for each.

defparam complicates this, because it can override a parameter from anywhere in the hierarchy, making the result depend on resolution order. Modern guidelines and tools discourage it in favor of inline #() overrides, which keep every override local to the instance it affects.

## Errors that only elaboration can see

Some errors are invisible until the design is assembled. Every file can parse cleanly while an instance names a module that no file defines. A port connection's width mismatch is only detectable once a parameter fixes the port width. Recursive instantiation terminates only if a parameter-dependent generate condition eventually fails, as when N is halved down to 1.

Because elaboration is where SystemVerilog's semantics live, building a complete elaborator is a large project. Surelog and slang each build a full elaborated model, and Surelog writes it out as UHDM so that other tools (a Yosys plug-in, Verilator, CIRCT) can consume one shared interpretation instead of implementing their own. That sharing matters because two elaborators that disagree produce two different designs from the same source.

## Going deeper

IEEE Std 1800-2017, clause 23.10 (elaboration) and clause 27 (generate constructs); the slang and Surelog documentation.
