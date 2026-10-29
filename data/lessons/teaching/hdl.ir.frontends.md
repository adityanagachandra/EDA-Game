# Open-source HDL front ends

Open-source EDA has several HDL front ends, each with a different scope: event-driven simulators, compiled simulators, VHDL tools, full SystemVerilog compilers and shared elaborated databases. A tool designer must know which front end to build on, what each gets right and how they connect, because a flow's language coverage is only as good as the front end it uses.

## The main tools and their roles

Icarus Verilog compiles Verilog with iverilog into an intermediate form, vvp code, that the vvp runtime simulates event by event with four-state values; it is close to the standard's scheduling model. Verilator translates Verilog and SystemVerilog into a C++ (or SystemC) model that is compiled with a testbench into a fast, statically scheduled, two-state simulator; it is also a widely used linter. GHDL analyzes, compiles and simulates VHDL with several code-generation back ends, and offers synthesis through a Yosys plug-in.

slang is a C++ library and compiler front end for SystemVerilog that aims at complete parsing, type checking and elaboration, with precise diagnostics; other tools use it as their front end. Surelog parses and elaborates SystemVerilog and writes the result as UHDM. Yosys has its own read_verilog, which covers Verilog-2005 and part of SystemVerilog.

## Front-end architecture and shared models

Every Verilog front end runs the same phases: preprocessing (macros, includes and ifdef), lexing and parsing, then elaboration (hierarchy, parameters and generate), then type checking and lowering into the tool's IR. Elaboration is where SystemVerilog is hardest and where tools diverge most.

UHDM (Universal Hardware Data Model) attacks that duplication. It is a serialized object model of a fully elaborated SystemVerilog design, based on the standard's VPI object model. Surelog produces it, and consumers such as a Yosys plug-in read it, so each consumer gets one shared elaborated design instead of parsing SystemVerilog itself. A Yosys plug-in built on slang does the same job by another route, translating slang's elaborated model into RTLIL; CIRCT's SystemVerilog import also builds on slang.

## Choosing and cross-checking front ends

Front ends differ widely in language coverage, and a shared conformance suite, sv-tests, makes the gaps visible by running thousands of small SystemVerilog examples through each tool and recording which parse, elaborate or simulate correctly. Verible is a separate family of parser-based tools for lint, formatting and indexing, which work on the syntax tree without full elaboration.

Running RTL through two independent front ends is cheap insurance. When Icarus and Verilator disagree, or a lint warning about width appears in one tool only, the cause is usually real: a truncated expression, a race or behavior that depends on x. Differential testing across independent implementations is the same idea as differential testing of an engine change, applied to the language itself.

## Going deeper

The documentation of Icarus Verilog, Verilator, GHDL, slang, Surelog and UHDM; the CHIPS Alliance sv-tests dashboard.
