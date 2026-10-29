# Path delay faults

A path delay fault says that the cumulative delay along one specific path from an input to an output exceeds the clock period, even if no single gate is grossly slow. It is the most faithful model of distributed timing defects and of process-induced slow paths, and the hardest to use: the number of paths in a netlist is exponential. Tool designers meet it in timing-aware ATPG, path selection from static timing, and at-speed test of critical paths.

## Counting paths

A **path** runs from a primary input (or scan flop output) through a chain of gates to a primary output (or scan flop input). Each path has two path delay faults, one for a rising and one for a falling transition at its input. Paths are counted by dynamic programming over the DAG in topological order: the number of paths reaching a gate is the sum over its inputs of the paths reaching them, with each primary input contributing 1. The total is the sum at the outputs. The number of paths through one internal node $g$ is (paths from inputs to $g$) times (paths from $g$ to outputs). For a node reached by 3 paths and reaching 2 outputs, that is 6 paths and 12 path delay faults.

## Exponential growth

Reconvergent fanout multiplies paths. A cascade of $k$ stages, each fanning out to $w$ parallel gates that reconverge, has $w^k$ end-to-end paths. Ten stages of width 2 give 1024 paths and 2048 path delay faults from a structure of only 30 gates. Real designs have far more paths than any test can cover, so tools test a **selected** set: the longest or most critical paths reported by static timing, often the paths with the smallest slack, plus paths through suspected weak spots. This is the same path-explosion problem that drives path-based analysis in timing.

## Robust and non-robust tests

A path test launches a transition at the path input and checks that it arrives in time at the output. A **robust** test detects the fault regardless of delays on other paths: side inputs along the path are held stable at non-controlling values when the on-path signal moves to its controlling value. A **non-robust** test only requires the side inputs to be non-controlling in the second vector, so a late side-input transition can mask the fault. Many paths have no robust test at all; some have no test of either kind because of the logic structure, and those paths are called functionally unsensitizable.

## Going deeper

Smith, "Model for delay faults based upon paths," International Test Conference, 1985; Bushnell and Agrawal (2000), chapter 12; Krstic and Cheng, "Delay Fault Testing for VLSI Circuits" (1998).
