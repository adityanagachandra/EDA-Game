# SymbiYosys and open formal tools

SymbiYosys (sby) is the open-source front end for formal hardware verification built on Yosys. It reads a short configuration file, runs Yosys to turn the design and its properties into a formal model, and dispatches the model to engines such as yosys-smtbmc and ABC. Knowing its file format, modes and engines is the practical way into every concept of this domain, because each engine corresponds to an algorithm studied here.

## The .sby file

A .sby file is divided into sections. [options] sets the mode, for example mode prove, and parameters such as depth 20. [engines] lists one or more engines to run. [script] holds the Yosys commands that read and prepare the design, typically read -formal for each source file, which enables the assert, assume and cover constructs, followed by prep -top with the top module's name. [files] lists the source files to copy into the run directory. An optional [tasks] section names several configurations, such as bmc and prove, and option lines prefixed with a task name apply only to that task, so one file can drive a bounded bug hunt and an unbounded proof of the same design. Results land in a directory per task, with a status file and trace files.

## Modes and engines

Modes select the question. bmc checks the assertions up to the depth. prove tries to establish them for all depths. cover searches for traces that reach each cover statement. live checks liveness properties with engines that support them. Engines select the algorithm. smtbmc drives an SMT solver, for example smtbmc yices, smtbmc boolector or smtbmc z3; in prove mode it runs k-induction, a BMC base case plus an induction step of the configured depth, both of which must pass. abc pdr runs ABC's IC3 implementation, which builds an inductive invariant itself and often succeeds where k-induction's step keeps failing. abc bmc3 is ABC's incremental bounded model checker. Several engines can be listed, and the first conclusive answer wins, which is a simple portfolio strategy.

## Reading the results

PASS in bmc mode means no assertion fails within the depth, not forever. FAIL comes with a counterexample trace, written as a VCD waveform and often as a testbench, that starts at reset and obeys all assumptions. In prove mode with smtbmc, a failing base case is a real bug, while a passing base case with a failing induction step produces an induction trace that starts from an arbitrary state satisfying the assertions for $k$ cycles; it is a counterexample to induction, usually fixed by a helper assertion or a larger depth, not by changing the RTL. A cover run reports each cover as reached, with a witness trace, or unreached. When a failing trace shows inputs the real environment never produces, the fix is an assume statement describing the legal behavior.

## Going deeper

The SymbiYosys documentation (YosysHQ); the Yosys manual, the chapters on formal verification and the write_smt2 and write_aiger backends; the ABC documentation for pdr and bmc3.
