# Functional verification in the flow

Functional verification asks whether the design does what the specification says. It runs from the first RTL to the final netlist and consumes most of a chip project's engineering effort. Because exhaustive testing is impossible, verification is a measurement problem: generate stimulus, check results against a reference, and measure how much of the behaviour has been exercised. A tool designer meets it in simulators, formal engines and coverage databases.

## Simulation, formal and emulation

Simulation drives the RTL with test stimulus and compares outputs to a reference model. Constrained-random testbenches (for example with the UVM methodology) generate many legal stimuli automatically. Formal verification proves a property for all inputs, using SAT, BDD or model-checking engines, but scales only to bounded problems or blocks. Emulation and FPGA prototyping map the design onto special hardware and run it millions of times faster than software simulation, so real software can boot. Equivalence checking is a separate formal step: it proves two versions of the design, for example RTL and the synthesized netlist, compute the same function.

## Coverage closes the loop

Since tests cannot be exhaustive, teams measure coverage. Code coverage counts which lines, branches and toggles were exercised. Functional coverage counts which specified situations occurred, as bins: for a packet size field with 8 bins and a packet type with 4 bins, the cross coverage has $8 \cdot 4 = 32$ bins. If tests have hit 26 of them, functional coverage is $26/32 = 81.25\%$. Coverage holes direct the next tests. Full coverage does not prove correctness, only that every listed situation was seen at least once; the list itself can be incomplete.

## Why late bugs cost more

A bug found in RTL is fixed by editing code and re-running tests. Found after synthesis, it may need an engineering change order on the netlist and a new round of equivalence checking. Found after layout, it reopens placement and routing. Found in silicon, it needs a respin. A common rule of thumb models each later stage as multiplying the fix cost by about ten. Under that assumption a bug that costs 2 engineer-days in RTL costs about $2 \cdot 10^3 = 2000$ engineer-days if it escapes three stages to silicon. That multiplier, more than any other fact, explains verification's share of effort.

## Going deeper

Wile, Goss and Roesner, "Comprehensive Functional Verification: The Complete Industry Cycle" (2005); the Accellera UVM user guide.
