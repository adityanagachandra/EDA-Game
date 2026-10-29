# MLIR and CIRCT as an agent-editable compiler

MLIR and CIRCT, and Yosys with its plug-in passes, are compilers designed to be extended. Mapping them onto the four-layer engine model shows where an agent's changes belong: configuration, plug-in passes and patterns, or the rarely touched core of IR semantics. This concept makes that mapping concrete, with change tiers, seams and the checks that promote an agent-written pass into a reusable skill.

## The four layers of a hardware compiler

The immutable core of a CIRCT-based compiler is the dialect definitions, the meaning of each operation and the verifier; every pass and lowering is written against them, so changing them is a core change. The strategy or policy layer is the passes and rewrite patterns, which decide how a design is transformed and are chosen, ordered and added without changing the IR's meaning. The parameter layer is pass options and the pipeline description: thresholds, modes, which passes run in which order. The observation layer is pass instrumentation and IR dumps before and after chosen passes, where an agent watches the compiler's state.

Yosys maps the same way: RTLIL data structures and internal cell semantics are the core, passes are policy, synthesis script options (whether to flatten, which ABC script to run) are parameters, and logs and write_rtlil dumps are observation hooks.

## Which tier a request lands in

A request to reorder existing passes is Tier 1, configuration: MLIR pass pipelines can be written textually, so no recompile is needed. A request for a new peephole rewrite, such as turning a multiplexer of two constants into an XOR, is Tier 2: a plug-in pattern registered at the dialect's extension point, leaving operation definitions untouched. Adding a Yosys pass as a loadable plug-in rather than patching Yosys sources is also Tier 2.

A request that changes what an operation means, such as the overflow behavior of comb.add, is Tier 3: a core edit, needing human approval and full regression of every pass that relied on the old meaning. Customer-specific policy, such as mapping wide multipliers to a particular DSP macro, belongs in a plug-in lowering pattern selected by configuration, never in the core lowering of the multiply operation, so the core stays shared and stable.

## Seams, checks and promotion

A rewrite pattern is a good seam for agent-written changes: small, local, independently testable and applied by a trusted driver that keeps the IR well-formed. What it may break is confined to the operations it matches. The differential check for a new pattern must prove two things: that each rewrite preserves the function of the matched operations, for example with an equivalence check over the before and after IR, and that quality on reference designs moves only as intended.

Once those checks pass, the pattern can be promoted from an experiment into a library skill. What makes it reusable is a manifest stating what it does, its inputs and outputs and its tunable options, plus its golden and differential tests. Other agents then select it by reading the manifest, without reading its source, and the tests keep it trustworthy as the compiler evolves.

## Going deeper

The MLIR documentation on pattern rewriting, dialect conversion and pass pipelines; the Yosys manual on writing and loading plug-in passes; the agent-first EDA notes in this curriculum.
