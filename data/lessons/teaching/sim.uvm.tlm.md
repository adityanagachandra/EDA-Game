# Transaction-level modeling and TLM ports

Transaction-level modeling (TLM) describes communication as function calls that carry whole transactions, not as signal changes on pins. Inside a UVM testbench, TLM ports connect components. In SystemC, TLM-2.0 is the standard for building fast virtual platforms of processors, buses and peripherals. Both rest on the same idea: abstracting time and pins away makes simulation orders of magnitude faster, at a cost in timing accuracy that has to be managed explicitly.

## Ports, exports and analysis ports

A TLM **port** is where a component calls an interface method such as `put`, `get` or `write`; an **export** forwards the interface through a hierarchy level; an **imp** (implementation) is where the method is actually implemented. Connecting port to imp means the caller invokes the callee's method directly. Put and get can block, so they suit producer-consumer pairs, often decoupled by a `uvm_tlm_fifo`. The **analysis port** is special: its `write` method is a function, so it takes zero simulation time and never blocks, and it broadcasts each transaction to every connected subscriber, one call per subscriber. With zero subscribers the write is still legal and simply goes nowhere. That is why monitors publish through analysis ports: a monitor never waits on, or even knows about, its consumers.

## SystemC TLM-2.0: loosely and approximately timed

TLM-2.0 defines a **generic payload** (command, address, data pointer, length, byte enables, response status) and two transport styles. **Loosely timed** models use `b_transport`, a blocking call that completes the whole transaction and returns with a timing annotation: the target adds its delay to a time argument instead of waiting. This is fast and is the basis of virtual platforms that boot software before RTL exists. **Approximately timed** models use `nb_transport` calls with explicit phases (BEGIN_REQ, END_REQ, BEGIN_RESP, END_RESP), so pipelining, outstanding transactions and contention can be modeled. Neither style is cycle-accurate; checking an arbiter's cycle-level fairness needs RTL. A scoreboard comparing a TLM reference model against RTL must compare transaction content and order, not cycle timing.

## Temporal decoupling and the quantum

A loosely timed initiator does not need to synchronize with the SystemC kernel on every transaction. With **temporal decoupling**, it keeps a local time offset and runs ahead of global simulation time, yielding to the kernel only when the offset reaches the **global quantum**. Over $T$ of simulated time with quantum $q$, it synchronizes about $T/q$ times. Example: with $q = 500$ ns over 100 us of simulated time, it synchronizes 200 times instead of once per transaction. Raising the quantum to 5 us cuts synchronizations by a factor of 10 and speeds simulation, because context switches into the kernel are expensive. The price is accuracy: initiators can disagree about time by up to one quantum, so interactions between them, such as a shared flag polled by two processors, become less precise.

## Going deeper

IEEE Std 1666-2011 (SystemC, including TLM-2.0); Accellera, "UVM 1.2 User's Guide" (2015), chapter on TLM; Aynsley, "OSCI TLM-2.0 Language Reference Manual" (2009).
