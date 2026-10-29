# MLIR and CIRCT

MLIR is a compiler infrastructure for building IRs at many abstraction levels in one framework; CIRCT (Circuit IR Compilers and Tools) applies it to hardware. Together they turn the classical ideas of this domain, SSA, dominance and passes, into reusable infrastructure for hardware compilers, and they are the main candidate substrate for compilers that agents can extend safely.

## Operations, regions and block arguments

In MLIR everything is an operation: an addition, a function and a whole module are all operations, with operands, results, attributes (compile-time constants such as widths or names) and optionally regions. A region is a list of blocks nested inside an operation, such as the body of a function, and each block is a list of operations ending in a terminator. Operations are grouped into dialects, namespaces of operations, types and attributes for one abstraction level or domain, such as arith, scf, llvm, or CIRCT's hw and comb.

MLIR is SSA, but it has no phi operation. Instead, blocks take arguments, and each branch passes values to its successor's arguments. A block argument is exactly a phi, placed by the same dominance-frontier reasoning, with the edge-to-value pairing written at each branch rather than collected at the join.

## Passes, patterns and progressive lowering

MLIR compilers are pipelines of passes. Many passes are built from rewrite patterns: local rules that match an operation and replace it. Canonicalization collects each operation's folders and simplification patterns into one cleanup pass applied to a fixed point. The dialect conversion framework lowers between dialects: the author supplies patterns plus a conversion target that declares which operations are legal afterwards, and the framework applies patterns until every illegal operation is gone or reports what could not be converted.

Progressive lowering is the resulting style: a program moves step by step through dialects of decreasing abstraction, each step small and testable. After every pass a verifier checks structural rules, such as operand types and dominance. Verification catches malformed IR; it does not prove that a pass preserved meaning, which needs tests or formal equivalence.

## CIRCT for hardware

CIRCT builds hardware compilers on this infrastructure. Its core dialects describe structure, combinational logic and state; other dialects cover FIRRTL, SystemVerilog emission, dataflow and high-level synthesis. Hardware module bodies are graph regions, where the dominance requirement of ordinary SSA regions is dropped, so an operation may use a value defined later in the block: feedback through registers is written directly.

Several tools ship with CIRCT. firtool compiles FIRRTL, as emitted by Chisel, into SystemVerilog. A SystemVerilog import built on slang brings Verilog into CIRCT dialects. arcilator compiles designs through the Arc dialect into fast native simulation code, and equivalence and bounded model checking tools work on the same IR. The shared infrastructure means a pass written once, with its tests, serves all of these flows.

## Going deeper

Lattner et al., "MLIR: Scaling Compiler Infrastructure for Domain Specific Computation" (CGO 2021); the CIRCT documentation and rationale documents.
