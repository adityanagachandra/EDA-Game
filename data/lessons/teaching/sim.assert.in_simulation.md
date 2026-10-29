# Assertions in simulation

Assertions state properties a design must satisfy, written next to the design and checked continuously while it simulates. SystemVerilog Assertions (SVA) come in two kinds: immediate assertions, procedural checks executed like statements, and concurrent assertions, clocked temporal properties evaluated by the simulator's scheduler. Knowing when an assertion samples its operands and when it evaluates them explains most surprising assertion results and how simulators keep assertions race-free.

## Immediate, deferred and concurrent assertions

An **immediate** assertion, `assert (a || b);`, is a procedural statement: it checks its expression at the moment it executes, inside an always or initial block. In combinational code that can fire on a glitch that settles later in the same time step. A **deferred** immediate assertion, `assert #0` or `assert final`, postpones reporting until values settle for the time step, which removes such false failures. A **concurrent** assertion, `assert property (@(posedge clk) req |-> ##[1:3] ack);`, is clocked and temporal: it starts an attempt at every clock and can check behavior over many cycles. Concurrent assertions live in modules, interfaces or checkers, and `bind` attaches them to a design instance without editing its source. `cover property` counts matches of the same kinds of properties.

## Sampling and evaluation regions

Concurrent assertions follow fixed scheduling rules. Their operands are **sampled in the Preponed region**, before anything changes in the time step. The property is **evaluated in the Observed region**, after the design's Active and NBA activity, using those sampled values. **Action blocks** (the pass and fail code) run in the **Reactive** region. The consequence: an assertion clocked on posedge clk sees every signal as it was just before that edge. If a flop updates q at the edge, the assertion at that edge sees the old q, and a waveform showing q rising at that edge while the assertion reports q was 0 is correct behavior. Functions such as rose, fell, stable and past compare successive sampled values. Sampling in Preponed is what makes assertions immune to races with the design.

## Engineering assertions for simulation

Assertions turn silent misbehavior into an immediate failure at the right place and time, which shortens debug from the scoreboard's late symptom to the protocol violation itself. The same properties serve formal tools, which try to prove them for all inputs, so writing them once pays twice. They are not free in simulation: every concurrent assertion starts an attempt each clock, and long or unbounded sequences keep many attempts alive, so thousands of complex assertions slow a simulator noticeably. They also prove nothing about scenarios that never happen: an overflow assertion on a FIFO that never fills passes trivially. Pair important assertions with cover properties on their preconditions, so coverage shows the assertion was really exercised.

## Going deeper

IEEE Std 1800-2017, clause 16 (assertions); Cerny, Dudani, Havlicek and Korchemny, "SVA: The Power of Assertions in SystemVerilog" (2nd ed., 2015).
