# Verilator's compilation model

Verilator is the most widely used open-source fast simulator for SystemVerilog. It is a compiler, not an interpreter: it translates synthesizable RTL into C++ (or SystemC) classes, which a host compiler turns into an executable model. Knowing its compilation model, its two-state value system and its threading scheme tells a tool designer what a compiled simulator can and cannot promise, and why teams use it alongside a four-state event-driven simulator.

## From RTL to a C++ model

Verilator parses and elaborates the design, optimizes it, orders the logic statically, and emits C++ classes. The user writes a harness, traditionally in C++, that instantiates the top-level model, sets input ports, calls `eval()` to make the model compute, and reads outputs; toggling a clock means setting the clock input and calling `eval()` again. The model only computes when asked. Recent versions can also run SystemVerilog testbenches with delays and event controls when built with the `--timing` option, which implements suspended processes with C++ coroutines. Either way the heavy lifting is compiled, statically scheduled code with no general runtime event queue for most of the design, which is where the speed comes from.

## Two-state values and X handling

Verilator's ordinary signals are two-state: X and Z do not exist at run time. Where the source assigns X explicitly, or a variable is never initialized, Verilator has to choose a concrete value. The `--x-assign` option controls explicit X assignments and `--x-initial` controls initial values; the choices include fixed zeros, a fast unspecified choice, and unique values that can be randomized at run time with a seed. Randomizing is the useful one for verification: a missing reset that a fixed 0 would hide shows up as seed-dependent behavior. What two-state simulation cannot do is show X propagation, so reset and initialization bugs that a four-state simulator reports as X on an output may pass silently in Verilator. Teams therefore run reset sequences in a four-state simulator too.

## Threads, lint and circular logic

Verilator can generate a multithreaded model. It partitions the statically ordered evaluation graph into coarse macro-tasks and schedules them across threads, because per-signal synchronization would cost more than it saves. Its front end is also a strict linter: width mismatches, undriven or unused signals, and constructs it cannot schedule are reported, which is why many teams run Verilator lint even when sign-off uses another simulator. One warning matters for performance: UNOPTFLAT says a signal appears to be part of a combinational loop after flattening. The logic may be acyclic bit by bit (different bits of one vector feed each other), but at the signal level Verilator cannot order it and falls back to slower iterative evaluation. Splitting the vector usually removes the warning and restores static ordering.

## Going deeper

The Verilator manual (sections on the C++ harness, x-assign and x-initial, multithreading and warnings).
