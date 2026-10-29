# Cycle-based and compiled simulation

Event-driven simulation pays for flexibility: every change goes through a queue. For synchronous designs, most of that flexibility is unused, because only the values at clock edges matter. **Cycle-based** simulation evaluates each combinational gate once per cycle in levelized order, and **compiled** simulation turns that order into straight-line host code. Together they are why Verilator is fast, and their cost model explains when an event-driven engine is still the better choice.

## Evaluating once per cycle in levelized order

A cycle-based simulator assumes a synchronous design: registers update on clock edges and combinational logic only needs its settled value before the next edge. Each cycle it updates the registers, then evaluates every combinational gate once, in level order, and stops. There is no notion of gate delay, no intra-cycle timing, and no glitch: the simulator computes the final value directly. Multiple clocks are handled by scheduling each clock's edges in time order and evaluating the logic that each edge affects. Logic whose behavior depends on delays, such as a pulse generator built from inverters, cannot be represented, and neither can timing checks. Because values only matter at edges, most cycle-based engines also use two-state values packed into machine words.

## Cost model and break-even activity

Let a design have $N$ gates. An event-driven engine evaluates only the active fraction $a$ of them per cycle, at cost $c_e$ per evaluation including queue overhead, so it spends $a N c_e$ per cycle. A cycle-based engine evaluates all $N$ at a much smaller cost $c_c$ each, spending $N c_c$. They break even at $a^* = c_c / c_e$. With $c_e = 50$ ns and $c_c = 5$ ns, $a^* = 10\%$: below 10% activity event-driven wins, above it cycle-based wins. With 20% activity the event-driven engine is $0.2 \times 50 / 5 = 2$ times slower per cycle. Real activity in large designs is often a few percent, but per-evaluation overhead ratios are large too, which is why both engines survive.

## What compilation adds and what it loses

An interpreted simulator walks data structures at run time: fetch a gate, look up its type, read its inputs. A **compiled** simulator generates host code, one statement per operation in levelized order, and compiles it with an optimizing compiler. Evaluation becomes straight-line code over machine words, with registers allocated by the compiler, no dispatch, and bitwise operations on 64 signals at once. The price is compile time after every change, which can dominate for short tests, and the loss of what the cycle model already dropped: delays, timing checks, detailed X behavior. Compiled cycle-based simulation therefore suits long runs on stable RTL, such as software on a processor model, while debugging fresh RTL with timing or X questions still favors an event-driven four-state simulator.

## Going deeper

The Verilator manual, sections on internals and scheduling; Abramovici, Breuer and Friedman, "Digital Systems Testing and Testable Design", chapter 3 (compiled and event-driven simulation).
