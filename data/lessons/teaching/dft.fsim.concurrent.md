# Concurrent fault simulation

Concurrent fault simulation keeps the fault-free circuit and, at each gate, a list of faulty copies of that gate that differ from it. Only where a fault's behaviour diverges does the simulator spend work. It extends naturally to sequential circuits, multi-valued logic and timing, and its idea, representing many machines as deltas from one reference, recurs across EDA in incremental timing and incremental simulation.

## Faulty gate records

At every gate the simulator stores the good gate (its input and output values) plus a **record** for each fault whose values at that gate differ from the good gate. A fault starts as a record at its own site, its **local** fault. When an input event reaches a gate, the good gate and every record there are re-evaluated. A record whose inputs and output now match the good gate has **converged** and is deleted, since the fault is invisible at that point. Because most faults agree with the good machine almost everywhere, the total number of records is far smaller than faults times gates.

## Divergence and scheduling

A record whose output differs from the good gate's output has **diverged**, and the fault's different value must be scheduled at every fanout gate, where it creates or updates a record. Take a 2-input AND whose good inputs are a = 0, b = 0, so the good output is 0. Faulty copies with inputs (1, 0) still output 0 and stop there, while copies with (1, 1) output 1 and diverge; only those propagate. When a diverged record reaches a primary output with a different value, its fault is detected and, with fault dropping, its records are removed everywhere.

## Where it fits

Concurrent simulation handles anything an event-driven simulator handles: X values, multiple strengths, timing, and state. In a sequential circuit a fault's record at a flip-flop carries a different state into the next cycle, so its effect can reappear later even with no difference at the outputs now. The cost is memory and list management: busy gates can hold long, changing lists. For full-scan designs, where the problem is combinational and two-valued, PPSFP is usually faster and leaner, which is why production stuck-at simulation uses bit-parallel engines while concurrent methods serve sequential, functional and timing-aware fault simulation.

## Going deeper

Ulrich and Baker, "Concurrent simulation of nearly identical digital networks," Computer, 1974; Bushnell and Agrawal (2000), chapter 5.
