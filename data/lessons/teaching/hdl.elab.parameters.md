# Parameters and parameter overrides

Parameters make one module definition describe a family of hardware: a FIFO of any depth, a bus of any width. Elaboration binds a value to every parameter of every instance and evaluates every expression that depends on them, including port widths. Front ends must get the override rules exactly right, because a wrong width propagates silently into synthesis and simulation.

## Defaults and named overrides

A module declares parameters with defaults: module buffer #(parameter DEPTH = 16, parameter W = 8). A localparam such as localparam PW = W + 2*DEPTH; is derived and cannot be overridden from outside, which makes it the right place for values that must stay consistent with the others. Port widths can use both, as in output [PW-1:0] bus.

An instance overrides parameters by name: buffer #(.W(12)) u0 (...); sets W to 12 and leaves DEPTH at its default. Elaboration then evaluates the localparam for this instance, $PW = 12 + 2 \times 16 = 44$, and the bus port of u0 is 44 bits wide. Every instance is evaluated separately, so two instances with different overrides have different port widths.

## Values passed down the hierarchy

Overrides can be expressions over the parent's own parameters. If core #(parameter N = 2) instantiates buffer #(.DEPTH(N*2), .W(N+3)), and top instantiates core #(.N(6)), elaboration resolves top down: u_core gets N = 6, then u_buf gets DEPTH = 12 and W = 9, so PW = $9 + 2 \times 12 = 33$.

The order is forced: a child cannot be specialized until its parent's values are known. This is why elaboration is a top-down tree walk, and why the older defparam statement, which can override a parameter from anywhere in the hierarchy, makes the result depend on resolution order and is discouraged. Each distinct set of parameter values produces one specialized module, which later stages treat as a separate definition.

## Positional overrides

An instance can also override parameters by position: buffer #(6, 10) u0 (...); assigns the values in the order the parameters are declared in the module header. With the header #(parameter DEPTH = 4, parameter W = 8) this sets DEPTH = 6 and W = 10, giving $PW = 10 + 2 \times 6 = 22$; if the header declared W first, the same instance would set W = 6 and DEPTH = 10, giving $6 + 20 = 26$.

Positional binding therefore couples every instance to the declaration order of the module, and reordering parameters in a header silently changes every positional instance. Named overrides avoid this, which is why style guides and linters prefer them. A front end must also reject attempts to override a localparam and report overrides of parameters that do not exist, both common symptoms of a mismatch between an instance and its module.

## Going deeper

IEEE Std 1800-2017, clause 6.20 (constants and parameters) and clause 23.10 (overriding module parameters).
