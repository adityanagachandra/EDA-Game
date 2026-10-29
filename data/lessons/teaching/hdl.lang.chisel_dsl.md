# Chisel and embedded hardware DSLs

Chisel is a hardware construction language embedded in Scala, and it stands for a broader idea: describing hardware with a program that generates it. Running the program elaborates the circuit, which is emitted as FIRRTL and compiled to Verilog. Embedded DSLs such as Chisel, Amaranth, SpinalHDL and Clash change where parameterization lives and what a hardware compiler receives as input.

## Generators embedded in a host language

A Chisel design is a Scala program. Executing it constructs a graph of hardware: modules, registers, wires and operators. Scala control flow runs at elaboration time, before any hardware exists, so a Scala for loop that creates eight adders is the analogue of a Verilog generate loop, and a Scala if decides which hardware to build at all. A Chisel when, by contrast, builds hardware (multiplexers) that chooses at run time. Confusing the two is the classic embedded-DSL bug.

Embedding gives generators the host language for free: types, functions, higher-order combinators, collections and libraries. A parameterized family of designs, such as a configurable cache or a systolic array, becomes an ordinary program, checked by the host compiler before any hardware is emitted.

## Connection semantics and widths

Chisel and FIRRTL use last-connect semantics: if a wire is connected unconditionally twice, x := a and later x := b, the second connection wins and x carries b. Inside a when block, a later connection overrides only under that condition, which is how conditional assignment becomes a multiplexer.

Reg creates a register on the module's implicit clock with no reset; RegInit adds a reset value. Widths may be left unspecified and are inferred. The plain + truncates to the larger operand width, while +& keeps the carry, so adding two 8-bit values with +& gives 9 bits. A wire that is never connected is reported as not fully initialized unless it is explicitly marked DontCare, which turns a forgotten connection into a compile error instead of a silent x.

## The landscape of hardware DSLs

Other embedded DSLs make similar trade-offs in different hosts. Amaranth (formerly nMigen) is embedded in Python; SpinalHDL in Scala; Clash takes another approach, compiling a subset of Haskell functions into circuits so that the functional program is the circuit itself. High-level synthesis tools sit further away, scheduling untimed C++ into RTL.

Every DSL ends in generated Verilog, because that is what simulators, synthesis and timing tools consume. Debugging therefore happens partly in the generated code: timing reports and waveforms name Verilog signals, so stable names and source locators in the emitted output are what connect a report back to the generator. A generator pays off when the design is a family of variants; for a single fixed block, the extra layer may cost more than it saves.

## Going deeper

Bachrach et al., "Chisel: Constructing Hardware in a Scala Embedded Language" (DAC 2012); the Chisel and Amaranth documentation.
