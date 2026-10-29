# ABC optimization scripts

ABC exposes its transformations as commands and combines them into scripts. resyn, resyn2 and compress2 are the classic recipes: sequences of balancing, rewriting and refactoring that each improve the AIG a little and set up the next pass. A script is a strategy in the engine's four-layer sense, and choosing or learning scripts is a natural job for an agent.

## The core commands

strash builds a structurally hashed AIG from the current network. balance (b) rebuilds AND supergates for minimum depth. rewrite (rw) replaces 4-input cuts by better precomputed structures. refactor (rf) collapses larger cuts and re-factors them. resub (rs) re-expresses nodes using existing divisors. The z variants, rwz and rfz, also accept zero-gain moves. fraig merges functionally equivalent nodes found by simulation and SAT. For mapping, map targets a standard-cell library and if maps to K-input LUTs. cec proves combinational equivalence of two networks. Each command preserves function; they differ in what they optimize and how much they cost.

## Reading resyn2

resyn2 is b; rw; rf; b; rw; rwz; b; rfz; rwz; b. It alternates area moves (rewrite and refactor) with depth moves (balance), and ends with zero-cost variants that reshape the graph so later passes, or a second run, find more gains. Each pass is greedy and local, so running resyn2 twice usually helps a little more. compress2 is the same pattern with the -l option on each command, preserving logic levels, for timing-critical designs. dch computes structural choices before mapping, so the mapper can pick among equivalent structures, which reduces the structural bias that any single AIG imposes.

## Scripts as strategy

Script quality depends on the design. Arithmetic-heavy logic responds differently from control logic, and long scripts trade run time for QoR. The commands are compiled, verified transformations; the script is policy that selects and orders them, so changing the script is a configuration or plug-in change, not a core edit. That makes script choice a natural search space: design-space exploration tries variants, and learned policies choose the next command from AIG statistics. Whatever chooses the script, every result must still pass equivalence checking, and proxies such as node count must be confirmed by mapped area and timing.

## Going deeper

Brayton and Mishchenko, "ABC: an academic industrial-strength verification tool", CAV (2010); Mishchenko, Chatterjee and Brayton, "DAG-aware AIG rewriting", DAC (2006); the ABC command documentation.
