# HLS tools and IRs: Vitis HLS, Bambu, CIRCT, Calyx

HLS is now built on compiler infrastructure. Commercial tools such as AMD Vitis HLS use LLVM-based front ends and expose the concepts of this domain as directives. Open tools and IRs, Bambu, CIRCT and its dialects, and the Calyx intermediate language, let a tool designer read and modify every stage. Knowing what each layer represents tells you where a new scheduling or binding idea should be implemented.

## Directives in a commercial flow

Vitis HLS exposes the main transformations as pragmas or directives. PIPELINE with an II target overlaps loop iterations. UNROLL replicates a loop body, fully or by a factor. ARRAY_PARTITION splits an array into banks, cyclic, block or complete, along a chosen dimension. DATAFLOW turns a function's sequence of sub-functions or loops into concurrent tasks connected by channels. INTERFACE chooses the protocol of each argument, such as a stream, an AXI4 master port or an AXI4-Lite register. INLINE controls whether a sub-function is merged into its caller. The report lists the achieved II per loop, its latency, resource estimates and the reason when a target II is missed. Every directive maps to a concept already covered here: modulo scheduling, unrolling, banking, dataflow and interface synthesis.

## Open tools and intermediate languages

Bambu, from the PandA project at Politecnico di Milano, is a complete open HLS framework: it takes C or C++ through a compiler front end and produces Verilog or VHDL, with its own scheduling, binding and memory allocation. CIRCT is an LLVM project that applies MLIR to hardware: a set of dialects at different abstraction levels, including handshake for dynamically scheduled dataflow circuits, Calyx for structured accelerators, and hw, comb and seq for RTL. Calyx, from Cornell, is an intermediate language whose components contain cells, groups of wire assignments, and a control program built from seq, par, if and while; its compiler turns the control program into finite-state machines. Lowering through several dialects lets each pass see the abstraction it needs.

## Static versus dynamic scheduling and MLIR-based HLS

Most HLS tools schedule statically: the compiler fixes every operation's cycle, and the controller is an FSM. Dynamically scheduled HLS, as in the Dynamatic tool and CIRCT's handshake dialect, builds elastic circuits in which each operation fires when its input tokens arrive through valid/ready channels. A loop whose branch arms take 2 or 20 cycles runs at the speed of the path actually taken, where a static schedule must budget for the slow arm; the price is handshake logic and buffers on every channel. MLIR-based flows such as ScaleHLS perform loop transformations in the affine dialect, where loop bounds and indices are affine functions so dependence analysis is exact, and then emit HLS C++ with directives. LLVM IR alone loses much of the loop and array structure these transformations need.

## Going deeper

Ferrandi et al., "Bambu: an open-source research framework for the high-level synthesis of complex applications" (DAC, 2021); Nigam et al., "A compiler infrastructure for accelerator generators" (ASPLOS, 2021); Josipović, Ghosal and Ienne, "Dynamically scheduled high-level synthesis" (FPGA, 2018); Ye et al., "ScaleHLS: a new scalable high-level synthesis framework on multi-level intermediate representation" (HPCA, 2022).
