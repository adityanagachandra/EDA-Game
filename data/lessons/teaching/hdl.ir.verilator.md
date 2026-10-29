# Verilator's compiled cycle model

Verilator is the fastest open-source Verilog simulator because it is a compiler: it turns RTL into C++ that evaluates the design in a statically chosen order with two-state values. That design choice explains both its speed and its differences from event-driven simulators, and it makes Verilator a natural engine for automated compile-and-simulate loops.

## A compiled, two-state model

An event-driven simulator such as Icarus keeps an event queue and evaluates processes as their inputs change. Verilator instead analyzes the whole design, orders its logic statically and emits C++ classes whose eval function computes the next state of the design in that fixed order. There is no event-queue overhead, and the generated code is compiled with full optimization.

Values are two-state by default: each bit is 0 or 1, mapping directly onto machine integers, and x and z are resolved to concrete values. In the classic flow the user writes a C++ (or SystemC) main loop that sets inputs, toggles the clock and calls eval; newer releases can also build a standalone binary from a Verilog testbench. With --lint-only, Verilator runs its front end and checks without building a model, which makes it a popular linter.

## Scheduling, warnings and threads

Static ordering must still respect the language. Verilator orders the computation so that all clocked blocks read pre-edge values before any register is updated, reproducing nonblocking semantics with ordering and shadow variables. Combinational logic is sorted by dependency. When the dependency graph appears cyclic, often because different bits of one vector feed each other, Verilator warns UNOPTFLAT and must iterate that part of the logic until it converges, which slows the model.

For speed on large designs, Verilator partitions the logic into tasks with dependencies and schedules them across threads. With its timing option, it supports delays and event controls by compiling processes into C++ coroutines, so ordinary Verilog testbenches with delays can run; without it, such constructs are not modeled.

## Where Verilator differs from four-state simulators

Two-state modeling can hide bugs. A register with no reset that reaches an output shows x in a four-state simulator but a concrete value in Verilator, so a design can pass in Verilator and fail elsewhere. Teams randomize initial values across regression runs (with the x-initial and x-assign options) so that logic silently depending on power-up zeros is caught.

Racy code is the other source of differences. The standard leaves the order of same-time processes open, so Verilator's static order and an event-driven simulator's order can both be legal and still disagree; the fix belongs in the code. For automated RTL loops Verilator is attractive: strict diagnostics and very fast runs, with the C++ compile step as the main latency to manage, for example with a compiler cache and incremental builds.

## Going deeper

Snyder, the Verilator manual; the Verilator documentation on scheduling and the timing option.
