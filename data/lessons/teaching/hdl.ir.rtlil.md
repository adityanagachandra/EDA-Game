# Yosys RTLIL

RTLIL (RTL Intermediate Language) is the in-memory and textual IR of Yosys, the open-source synthesis framework. Every Yosys front end produces it, every pass transforms it and every back end writes from it, so it is the seam on which the whole tool is built. Understanding its objects is the prerequisite for reading Yosys logs, writing passes or plugging another front end into Yosys.

## Modules, wires, cells and connections

An RTLIL design is a set of modules. A module holds wires, cells, processes, memories and connections. A wire is just a named multi-bit signal with flags for being an input or output port; it stores nothing. A cell is an instance with a type, parameters and port connections: either an internal cell type, whose name begins with a dollar sign (such as the add, mux and dff cells), or a user module or black box.

Port connections use SigSpecs. A SigSpec is a concatenation of chunks, each a slice of a wire or a run of constant bits, so a cell input can connect to bits 7:4 of wire a followed by four constant zeros without any extra cells. Module-level connections join two SigSpecs directly, which is how plain assigns are represented.

## Coarse cells, fine cells and processes

Internal cells come in two grains. Coarse-grain cells are word-level operators parameterized by width and signedness, such as add with A_WIDTH, B_WIDTH, Y_WIDTH and A_SIGNED: one generic type covers every width, and passes read the parameters. Fine-grain cells are single-bit gates and flip-flops. Synthesis starts with coarse cells, then techmap lowers them to fine cells for ABC and technology mapping.

Two higher-level objects exist before lowering. Right after read_verilog, an always block is an RTLIL process: a decision tree of switch and case rules, plus sync rules for its clock edges. The proc passes later turn processes into multiplexer and flip-flop cells. Arrays become memory objects with read and write port cells, which the memory passes map to RAM resources or to flip-flops.

## Text format and attributes

write_rtlil prints a design in RTLIL text and read_rtlil reads it back exactly. This round trip is how pass development and debugging work: save the design before a suspicious pass, run the pass alone, diff the outputs. Regression tests for Yosys passes are often written as small RTLIL or Verilog inputs plus expected results.

Every object carries attributes. The src attribute records the source file and line, and it travels with cells through passes so that reports point back at RTL. The keep attribute tells cleanup passes not to remove an object even if nothing reads it; keep_hierarchy tells flatten to leave a module's boundary. Because RTLIL is the stable shared IR, a plug-in pass written against it works for designs from Verilog, from VHDL through GHDL or from SystemVerilog through slang or Surelog front ends.

## Going deeper

Wolf, the Yosys manual, chapters on RTLIL and internal cell types; the Yosys source, kernel/rtlil.h.
