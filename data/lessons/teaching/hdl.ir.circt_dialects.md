# CIRCT dialects: HW, Comb, Seq and lowering

CIRCT represents hardware with a family of MLIR dialects, each capturing one concern: structure, combinational logic, state, SystemVerilog emission and several higher-level models. A hardware compiler built on CIRCT is a sequence of lowerings between them. Knowing what each dialect holds is what lets you read CIRCT IR, choose where a transformation belongs and follow a design from Chisel or SystemVerilog to emitted Verilog.

## HW, Comb and Seq

Three core dialects describe a synchronous circuit. hw holds structure: hw.module with ports, hw.instance, hw.output and hw.constant, plus aggregate types such as arrays and structs. comb holds side-effect-free combinational operations on integer types: add, mul, mux, icmp, concat, extract and the bitwise operators. seq holds state: register operations such as seq.compreg and seq.firreg, and the clock type.

Comb integers are signless, following MLIR's convention: an i8 is just 8 bits. Signedness lives in the operations that care about it, such as the predicate of icmp (slt versus ult) or the choice between signed and unsigned division. This removes a whole class of Verilog sizing and signedness ambiguities from the IR.

## SV, ExportVerilog and the lowering path

The sv dialect sits just above Verilog text. It has operations for SystemVerilog-specific constructs, such as always blocks, regs, wires, assigns and ifdef, used to control what the emitted code looks like. ExportVerilog prints hw, comb, seq and sv operations as SystemVerilog, choosing names and expression layout so that downstream simulators, linters and synthesis tools accept it.

A Chisel design flows firrtl, then hw with comb and seq, then sv, then Verilog text. Registers stay as seq operations rather than being lowered early to always blocks, because analyses and optimizations can reason about a register operation directly without recognizing a procedural idiom; only near the end do they become sv constructs for printing.

## Higher-level and simulation dialects

Other dialects capture other models of hardware. handshake describes latency-insensitive dataflow with ready-valid channels, used by flows that compile software-like programs into circuits. calyx is the IR of the Calyx compiler for accelerators: components with structural groups and an explicit control schedule. llhd (low-level hardware description) has signals, processes and time, matching event-driven HDL simulation semantics. moore captures SystemVerilog semantics coming out of the slang-based importer before they are lowered to the core dialects. arc models state transfers for the arcilator simulator.

Each sits at a different point between source semantics and netlist, and progressive lowering connects them. Choosing the dialect at which to implement a transformation is the main architectural decision when extending CIRCT: as high as the information it needs, as low as the simplicity it wants.

## Going deeper

The CIRCT documentation for the HW, Comb, Seq, SV, Handshake, Calyx, LLHD, Moore and Arc dialects; Nigam et al., "A Compiler Infrastructure for Accelerator Generators" (ASPLOS 2021).
