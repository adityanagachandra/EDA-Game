# Interfaces, modports and packages

Large SystemVerilog designs repeat the same buses and the same type declarations across dozens of modules. Interfaces bundle a bus into one named port, modports give each side its directions, and packages share types and constants. A front end has to resolve all three during compilation and elaboration, and getting their ordering rules wrong is a common source of tool failures on real code.

## Interfaces bundle signals

An interface declares a group of related signals, such as the 40 wires of an AXI channel, and modules take the interface as a single port. Connecting a bus is then one connection instead of forty, and adding a signal means editing the interface once rather than every port list and every instance, which removes a whole class of drift bugs.

Interfaces can hold more than signals: tasks, functions, assertions and small procedural blocks, so a protocol checker can live next to the bus it checks. They are an elaboration-time convenience. After elaboration and flattening, an interface port has become its individual signals, connected exactly as if they had been listed one by one; synthesis sees ordinary nets.

## Modports give directions

Inside an interface every signal is just a variable or net; directions depend on who is looking. A modport defines one view: in the master view, valid and data are outputs and ready is an input; in the slave view the reverse. A module port declared with a modport gets those directions, and a write to a signal the modport lists as input is a compile-time error.

To type-check a module with an interface port, a front end needs the interface definition (signal names and widths) and the modport (directions). Interfaces often take the clock as a port so that clocking blocks and synchronous helpers inside can refer to it, keeping the bus definition self-contained.

## Packages and compile order

A package is a namespace of shared declarations: types, parameters, constants and functions. A typedef struct packed declared in a package and used for the ports of two modules guarantees that both agree on the field layout and width. Modules import from the package, either explicitly (import pkg::opcode_t;) or with a wildcard (import pkg::*;), which is resolved lazily and can make a name's origin unclear or clash with local names.

Packages are compile-time dependencies. A module that imports a package can only be analyzed after the package, so file order (or a dependency-aware compile) matters. Module instantiation is different: an instance can refer to a module defined in a later file, because elaboration links instances to definitions after all files are parsed. Build systems and front ends that ignore this distinction fail on otherwise valid code.

## Going deeper

IEEE Std 1800-2017, clauses 25 (interfaces) and 26 (packages); Sutherland, "SystemVerilog for Design", chapter 10.
