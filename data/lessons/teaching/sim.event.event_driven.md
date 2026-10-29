# Event-driven logic simulation

A logic simulator answers one question millions of times: given a change on some input, what changes next, and when? An event-driven simulator answers it by doing work only where something changed. It keeps a queue of pending value changes (events), and evaluates only the gates whose inputs just changed. This is the engine inside Icarus Verilog and inside every commercial four-state simulator, and its cost model explains why simulation speed depends on activity, not on design size.

## Events and evaluations

An **event** is a change in the value of a net at a particular time. An **evaluation** is the work of recomputing one gate (or one process) because an input changed. The two are different. When net $n$ changes, the simulator looks up the fanout list of $n$ and evaluates each gate on it. If the gate's newly computed output equals its current output, nothing is scheduled: the evaluation cost time but produced no event. If the output differs, the simulator schedules an event on the output net at the current time plus the gate delay. The classic two-list loop repeats: take all events due at the earliest time, apply them, collect the set of fanout gates, evaluate each once, schedule the resulting changes. Nets that never change cost nothing.

## Counting events after an input toggle

Take inputs $a = 1$, $b = 0$, $c = 1$ and unit-delay gates $n_1 = \text{AND}(a, b)$, $n_2 = \text{OR}(n_1, c)$, $n_3 = \text{XOR}(a, n_2)$. Initially $n_1 = 0$, $n_2 = 1$, $n_3 = 0$. At time 0 input $b$ rises: 1 event. Its fanout is $n_1$, so evaluate $n_1$: it becomes 1, scheduled for time 1 (1 evaluation). At time 1, $n_1$ changes (event 2); evaluate $n_2$: $1 \lor 1 = 1$, unchanged, so nothing is scheduled (evaluation 2). The activity stops there. In total: 2 events and 2 evaluations, and $n_3$ was never touched. Counting both numbers tells you what an engine actually pays: evaluations cost CPU time, events cost queue operations.

## Activity and why event-driven wins

The **activity factor** is the fraction of gates evaluated per unit of work, for example per clock cycle. In a typical synchronous design only a few percent of nets change in a given cycle, so an event-driven simulator touches a small fraction of the gates, while a simulator that evaluates everything pays for every gate. The overhead per event is real, though: inserting into a time-ordered queue, walking a fanout list and dispatching an evaluation costs far more than the Boolean operation itself. Glitches make it worse. Reconvergent paths with unequal delays can make a net change twice in one cycle, and each glitch spawns more events downstream. When activity is high, a levelized cycle-based engine that evaluates every gate once without a queue can be faster.

## Going deeper

Abramovici, Breuer and Friedman, "Digital Systems Testing and Testable Design", chapter 3 (logic simulation); the Icarus Verilog source, in particular the vvp runtime.
