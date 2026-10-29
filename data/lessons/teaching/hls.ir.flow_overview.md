# From algorithm to RTL: the HLS flow

High-level synthesis (HLS) compiles an untimed algorithm, usually C or C++, into a register-transfer-level design: a datapath of functional units, registers and multiplexers, and a controller that sequences it. It sits one level above logic synthesis, and its output is the RTL that Yosys or a commercial tool then maps to gates. A tool designer cares because HLS makes the decisions RTL designers otherwise make by hand: which cycle each operation runs in and which hardware it shares.

## What HLS turns into what

The input is a function whose semantics are sequential and untimed: statements, loops, arrays, arithmetic. The output is a clocked machine that computes the same function. Everything the source states, the arithmetic and the data dependences, is preserved; what the source leaves open is time and space. HLS chooses the clock cycle of every operation (scheduling), how many adders and multipliers exist (allocation), which operation runs on which unit and which register holds which value (binding), and how the block talks to the outside world (interface synthesis). The same C can therefore produce a one-multiplier design that takes 64 cycles or a 64-multiplier design that takes one, depending on directives such as unroll and pipeline. That freedom is the whole point, and also why the result must be verified against the C.

## The steps of the flow

A front end parses the C into an intermediate representation, typically SSA inside a control-data-flow graph, and runs compiler optimizations: inlining, constant propagation, dead-code removal, bit-width narrowing, loop transformations. Then the synthesis core runs. Allocation picks the number of units of each type. Scheduling assigns each operation a control step, respecting dependences and the allocation. Binding maps operations to unit instances and values to registers, which fixes the multiplexers in front of each unit. Finally the tool emits RTL: the datapath, a finite-state machine whose states are the control steps, and the interface logic. A value produced in step 2 and consumed in step 5 must survive two clock edges, so scheduling directly creates register demand. The order matters: binding needs a schedule, and resource-constrained scheduling needs an allocation.

## Where HLS fits and where it struggles

HLS shines on algorithmic blocks with regular loops over arrays: filters, transforms, linear algebra, image pipelines, network packet processing. It struggles where the cycle behavior is the specification, as in a bus protocol defined by a timing diagram, and where hardware cannot be statically sized, as in unbounded recursion or malloc. Its reports of latency, initiation interval and resources are estimates from its own operator models, made before logic synthesis and routing, so the final clock frequency can differ. And because the tool changed the timing, the generated RTL must be checked with C/RTL co-simulation, running the C testbench against the RTL. The practical result is a short design loop: change a directive, re-run HLS in minutes, and compare quality of results, a loop that later invites automation.

## Going deeper

Coussy and Morawiec (eds.), "High-Level Synthesis: From Algorithm to Digital Circuit" (2008); De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapters 4 to 6; Kastner, Matai and Neuendorffer, "Parallel Programming for FPGAs" (2018).
