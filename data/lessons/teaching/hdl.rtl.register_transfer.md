# Register-transfer level design

Register-transfer level (RTL) is the abstraction at which almost all digital hardware is written and at which the EDA flow begins. An RTL description says which values are stored in registers and how each register's next value is computed from the current state and inputs, once per clock cycle. Everything downstream (synthesis, equivalence checking, timing, placement) consumes or preserves this contract, so a tool designer must know exactly what RTL fixes and what it leaves open.

## Registers and next-state logic

An RTL design is a synchronous state machine written at word level. The state is the set of registers $s$; the inputs are $x$; on every active clock edge each register takes the value of a next-state function, $s_{n+1} = f(s_n, x_n)$, and outputs are $y_n = g(s_n, x_n)$. The functions $f$ and $g$ are combinational: words of bits combined by operators such as +, comparisons, multiplexers and bitwise logic.

The transfer acc ← acc + x describes two components: a register holding acc and an adder computing acc + x from its output. A synthesizer builds exactly that, a register and the adder that feeds it. Three registers with no inputs and a reset therefore evolve deterministically: their values after ten cycles are $f$ applied ten times to the reset state. Nothing about gates, delays or placement appears in the description.

## The cycle-accurate abstraction

RTL ignores propagation delay inside the combinational logic. That is legitimate only because of synchronous design: timing closure guarantees that every path settles before the next capturing edge, so only the settled values at the edges matter. This is also why RTL simulation is fast. A 32-bit addition is one operation per cycle at RTL, but dozens of gate evaluations, each a potential event, in a gate-level simulation.

The abstraction also explains what two designs have in common. A ripple-carry adder and a prefix adder that produce the same sum are the same RTL; they differ in delay and area, which RTL semantics cannot see. Writing a * b fixes the function and the cycle in which the product appears, and leaves the architecture (array, Booth, Wallace tree) to synthesis, which picks it from constraints.

## What later tools may change

Synthesis must preserve the observable behavior: the sequence of outputs for every input sequence from reset. It need not preserve internal registers. Re-encoding a state machine from binary to one-hot, or retiming registers across logic, changes which values are stored where but not what leaves the block. This has a direct consequence for verification. When registers correspond one to one, equivalence reduces to comparing the combinational logic between matched registers; after re-encoding or retiming, correspondence is lost and a sequential method is required.

RTL is also the handoff point between front end and back end. High-level synthesis outputs RTL because the whole implementation flow consumes it. The cost of an RTL bug grows with the stage at which it is found: a functional fix after layout reruns synthesis, placement, routing and timing, unless an engineering change order (ECO) can patch the netlist locally.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture", chapters 3 and 4; Micheli, "Synthesis and Optimization of Digital Circuits", chapter 1.
