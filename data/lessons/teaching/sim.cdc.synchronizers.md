# CDC verification and synchronizers

Whenever a signal crosses from one clock domain to an unrelated one, the receiving flop can sample it while it changes and go metastable. Clock-domain crossing (CDC) verification checks that every such crossing uses a correct synchronization scheme. It sits between simulation and static analysis, because ordinary RTL simulation is blind to the failure: it never models metastability. A verification tool designer needs both the circuit facts and the reasons simulation misses them.

## Metastability and the two-flop synchronizer

If a flop's data input changes inside its setup-hold window, the flop can enter a metastable state and take an unbounded, though usually short, time to resolve to 0 or 1. Across asynchronous clocks this will happen eventually, so it cannot be prevented, only made harmless. The standard **two-flop synchronizer** feeds the asynchronous signal into a flop clocked by the destination clock, then into a second flop. The first flop may go metastable; it has almost a full destination clock period to resolve before the second flop samples it, so the probability that metastability propagates beyond the second flop is tiny. The synchronizer does not fix the latency: the destination sees a change one cycle earlier or later depending on when it was sampled, and designs must tolerate that.

## Structural CDC checks

Static CDC tools find every path between domains and classify the scheme used. Common violations: **combinational logic before the synchronizer**, whose glitches can be captured as false pulses, so a crossing must start from a flop; **fanout of the first synchronizer stage** to several loads, which can resolve differently; **reconvergence**, where related bits pass through separate synchronizers and recombine, arriving in different cycles and forming combinations the source never produced; and **fast-to-slow pulses**, where a one-cycle pulse in a fast domain can fall between slow-domain edges and vanish, requiring pulse stretching, a toggle scheme or a handshake. Multi-bit values need Gray coding, a handshake that holds data stable, or an asynchronous FIFO. The tool reports each crossing as synchronized, unsynchronized or using a recognized scheme.

## Why simulation misses CDC bugs

In zero-delay RTL simulation, a flop always captures a clean value, and the synchronizer's latency is always exactly two cycles. Metastability never happens, and the one-cycle uncertainty never appears, so logic that only works with one particular latency passes every test. This is why CDC is checked structurally and with formal protocol checks, and why simulation adds **metastability injection**: a model that randomly adds or removes one cycle of delay at each synchronizer output, emulating unpredictable resolution. Under injection, reconvergence bugs and fragile handshakes show up as test failures. Assertions on handshake signals, such as "data is stable while the request is pending", catch protocol errors that structure alone cannot see: a structurally correct crossing can still be used incorrectly.

## Going deeper

Ginosar, "Metastability and synchronizers: a tutorial", IEEE Design and Test of Computers (2011); Cummings, "Clock Domain Crossing (CDC) Design and Verification Techniques Using SystemVerilog" (SNUG 2008).
