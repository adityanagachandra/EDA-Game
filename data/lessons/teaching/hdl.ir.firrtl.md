# FIRRTL and its lowering

FIRRTL (Flexible Intermediate Representation for RTL) is the IR between Chisel and Verilog. Chisel elaborates a circuit and emits FIRRTL; a FIRRTL compiler infers widths and resets, resolves conditional connections, flattens aggregate types and emits Verilog. Its passes are a compact case study in hardware-compiler lowering, now implemented in CIRCT's firtool.

## Circuits, connections and widths

A FIRRTL circuit contains modules with ports, wires, registers, nodes and memories, and statements that connect them. Conditional connections appear in when blocks; after lowering, each becomes a multiplexer selecting between the connections made inside and outside the block. Aggregate types (bundles and vectors) describe structured ports, and a flipped field in a bundle flows in the opposite direction to the rest, as ready does in a valid-ready interface.

Widths may be left unspecified. Width inference computes the minimal widths that satisfy every expression and connection, using the primitive operations' width rules: add is lossless, so adding an 8-bit and a 5-bit UInt gives $\max(8, 5) + 1 = 9$ bits, and multiplication gives the sum of the widths.

## The lowering passes

A handful of passes take FIRRTL from Chisel's output to netlist form. InferWidths fixes every unspecified width. InferResets resolves abstract reset types to synchronous or asynchronous resets from how they are driven. ExpandWhens implements last-connect semantics: it turns nested conditional connections into one multiplexer expression per signal, so that afterwards every sink has a single driving expression. LowerTypes splits bundles and vectors into individual ground-type signals with derived names, because Verilog ports and wires are ground types.

After these passes the circuit is essentially a netlist of primitive operations, registers and memories, which maps directly onto CIRCT's hw, comb and seq dialects for optimization and Verilog export.

## Initialization, invalidation and annotations

FIRRTL enforces that every sink is driven on every path. A wire that is never connected, or connected only under some conditions, is reported as not fully initialized, which turns a forgotten connection into a compile error rather than a silent x in simulation. Marking a signal invalid (DontCare in Chisel) is an explicit don't-care: the compiler may choose any value, which can simplify logic, much like assigning x in Verilog for synthesis.

Annotations carry metadata alongside the circuit, targeting specific components: don't-touch markers that protect signals from optimization, memory configuration, or hooks for verification and test. Current Chisel uses the CIRCT-based firtool, the MLIR FIRRTL compiler, which replaced the original Scala FIRRTL compiler; the passes above now live in CIRCT's FIRRTL dialect.

## Going deeper

Izraelevitz et al., "Reusability is FIRRTL Ground: Hardware Construction Languages, Compiler Frameworks, and Transformations" (ICCAD 2017); the FIRRTL specification; the CIRCT FIRRTL dialect documentation.
