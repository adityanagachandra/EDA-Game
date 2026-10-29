# Datapath and control

A processor, and most accelerators, splits into a datapath, which holds and transforms data (registers, ALUs, shifters, muxes, memories), and a control unit, which decides each cycle what the datapath does by driving its mux selects, register enables and operation codes. This split is how RTL is organized, how high-level synthesis builds hardware from a schedule, and how verification divides its effort.

## Datapath elements and control signals

The datapath is built from the components of the earlier subdomains: a register file, an ALU (adder, logic unit, shifter, comparator), multiplexers that route operands, and memory interfaces. Every mux select, register write enable, ALU operation code and memory read or write strobe is a control signal.

For one instruction, say `add r3, r1, r2`, control must select the register-file read addresses, set the ALU to add, route the ALU result (not the memory data) to the register file's write port, and enable the write to r3. A different instruction is a different setting of the same control signals on the same datapath.

## Single-cycle and multicycle control

In a single-cycle design every instruction completes in one long clock cycle, and control is pure combinational decoding of the instruction. The cycle must fit the slowest instruction (typically a load: instruction fetch, register read, address add, data memory, register write), so every other instruction wastes time.

A multicycle design breaks each instruction into steps of one short cycle each, reusing units across steps (one memory for instructions and data, one ALU for address and arithmetic). Control becomes a finite-state machine whose states are the steps, and instructions take different numbers of cycles. Pipelining then overlaps those steps across instructions. High-level synthesis builds the same structure automatically: its schedule assigns operations to control steps, its binding shares units, and its controller is an FSM.

## Hardwired versus microcoded control

Control can be hardwired, synthesized logic and FSMs derived from the instruction set, or microcoded, where each instruction maps to a sequence of microinstructions stored in a ROM and the control signals are read from it. Microcode is easier to change and suits complex instruction sets; hardwired control is faster and suits RISC designs and accelerators.

For verification, the split matters too. Datapath bugs are usually data-dependent (overflow, a wrong bit width) and suit word-level and arithmetic reasoning; control bugs are sequence-dependent (a missed stall, a wrong state) and suit model checking and coverage of FSM transitions. Tools for each, and their capacity limits, follow this boundary.

## Going deeper

Patterson and Hennessy, "Computer Organization and Design: The Hardware/Software Interface" (RISC-V ed., 2017), chapter 4; Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), chapter 7.
