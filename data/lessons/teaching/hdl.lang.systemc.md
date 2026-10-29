# SystemC and transaction-level modeling

SystemC is a C++ class library with a discrete-event simulation kernel, standardized as IEEE 1666. It is used for architecture models, virtual platforms that run software before silicon exists, and as an input to high-level synthesis. For a tool designer it is a useful contrast with Verilog and VHDL: the same evaluate-update semantics, delivered as a library rather than a language.

## Modules and processes in C++

A SystemC model is an ordinary C++ program compiled with a standard compiler and linked with the SystemC kernel. SC_MODULE(adder) declares a class derived from sc_module, whose constructor creates ports, signals and processes. Processes come in two kinds. An SC_METHOD runs to completion each time an event in its sensitivity list fires; it has no stack of its own, so triggering it is a plain function call. An SC_THREAD has its own stack and can call wait() to suspend in the middle of its body, which makes sequential behavior easy to write but costs a context switch at each suspension.

A clocked SC_METHOD sensitive to clk.pos() corresponds to always @(posedge clk). Because it is all C++, models can use any library, templates and data structures, which is both SystemC's power and the reason only a subset is synthesizable.

## Evaluate-update and delta cycles

An sc_signal follows evaluate-update semantics. A process that calls s.write(5) and then s.read() sees the old value: the write is applied in the update phase, after all runnable processes have executed in the evaluate phase. One evaluate phase plus one update phase is a delta cycle, and it does not advance simulated time. Updated signals can wake more processes, so the kernel repeats delta cycles until nothing is runnable, then advances time to the next scheduled event.

This is the same mechanism as VHDL signals and Verilog nonblocking assignment, adopted for the same reason: processes that run at the same time see consistent values regardless of the order in which the kernel happens to execute them.

## Transaction-level modeling

TLM-2.0 abstracts communication away from pins. Instead of toggling address, data and handshake signals for a bus read, an initiator calls a function such as b_transport with a generic payload describing the whole transaction. A memory model answers in one call. Virtual platforms built this way run fast enough to boot an operating system.

Loosely timed models add temporal decoupling: an initiator runs ahead of simulated time and synchronizes with the kernel only once per quantum, cutting context switches at the cost of timing accuracy; approximately timed models add phases for more accurate contention. TLM models are not formally equivalent to RTL; teams connect the two through transactors, which convert transactions to cycle-level signals so that one testbench can drive either. A synthesizable SystemC subset, with bit-accurate types such as sc_uint, is a common input to high-level synthesis.

## Going deeper

IEEE Std 1666-2011; Black, Donovan, Bunton and Keist, "SystemC: From the Ground Up" (2nd ed.); Aynsley, the OSCI TLM-2.0 user manual.
