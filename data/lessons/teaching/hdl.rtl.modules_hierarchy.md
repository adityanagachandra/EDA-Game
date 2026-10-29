# Modules, ports and design hierarchy

A module is the unit of definition in Verilog and SystemVerilog, and instances of modules form the design hierarchy. Elaboration builds that hierarchy, synthesis may keep or flatten it, and every report names objects by their hierarchical path. Tool builders need a precise picture of how ports bind, what graph the hierarchy forms and what a tool may do with module boundaries.

## Ports and connections

A module declares ports with directions: input, output or inout (bidirectional, typically tristate pads). An instance binds expressions to those ports in one of two ways. Named binding, adder u1 (.y(s), .a(x), .b(z)), connects by port name, so order does not matter. Positional binding, adder u2 (s, x, z), connects the first expression to the first declared port, and so on; if the declaration is (input a, input b, output y), s drives a, which is almost certainly a bug.

A module is a definition; an instance is one use of it, with its own state and its own connections. Width mismatches between a port and the signal connected to it are legal in Verilog and follow assignment rules (truncation or extension), which is why lint tools flag them: they rarely mean what the designer intended.

## The instance tree

Instances form a rooted tree: each instance has exactly one parent and the top module is the root. Module definitions form a directed acyclic graph instead, because one definition may be used under several parents. A register is named by its hierarchical path, such as top.cpu.alu.acc, and two instances of the same module have distinct registers with distinct paths: hardware is never shared between instances.

Counting is a tree computation. If top has 3 instances of core and each core has 2 instances of lane, there are $3 \times 2 = 6$ lanes. The definition graph must be acyclic, because a module that instantiates itself describes infinite hardware. Recursion is legal only when a parameter-dependent generate condition stops it during elaboration.

## What tools do with boundaries

A tool may keep hierarchy or flatten it. Keeping it allows separate compilation, per-block reuse and smaller problems for each run, at the cost of optimization across boundaries: a constant tied to a submodule's input cannot propagate inside while the boundary remains. Flattening removes the boundaries and exposes that optimization.

Two further mechanisms come from boundaries. A black box is a module known only by its ports, used for IP and macros implemented elsewhere; the tool connects to it but cannot see inside. Uniquification is needed when instances of one module diverge, for example because their inputs are tied to different constants and optimization treats them differently: each divergent instance needs its own definition before it can be optimized on its own.

## Going deeper

IEEE Std 1800-2017, clause 23 (modules and hierarchy); the Yosys manual entries for hierarchy, flatten and uniquify.
