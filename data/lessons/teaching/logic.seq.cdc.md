# Clock domain crossing

Real chips have many clocks: per interface, per power domain, per performance mode. Wherever a signal launched by one clock is sampled by an unrelated clock, the synchronous guarantees stop and metastability, data incoherence and lost pulses become possible. Clock-domain-crossing (CDC) verification is a dedicated tool category because static timing analysis cannot check these paths; it treats them as false paths and relies on structure instead.

## Single-bit crossings

A single control bit crossing to an asynchronous domain goes through a two-flop (or longer) synchronizer, which bounds metastability. Two further rules apply. The signal must come straight from a flip-flop in the source domain, not from combinational logic, because a glitch could be sampled as a pulse. And the signal must be held long enough to be seen: a one-cycle pulse from a fast domain can fall entirely between two edges of a slow clock and vanish.

Pulses are therefore converted to level changes (toggle synchronizers) or stretched, and the receiver detects the edge after synchronization.

## Multi-bit crossings and reconvergence

Synchronizing each bit of a bus separately is wrong. Each synchronizer resolves independently, so on a cycle when several bits change, some may arrive one cycle later than others and the receiver sees a value that never existed in the source. The same problem appears when two separately synchronized control signals reconverge in the destination logic.

Three safe patterns exist. Gray-code the value so only one bit changes per step (works for counters and pointers). Hold the data stable and send a single synchronized qualifier (valid, or a request-acknowledge handshake); the receiver samples the data bus only after the qualifier arrives, by which time the bus is stable. Or use an asynchronous FIFO, which combines both.

## The asynchronous FIFO

An asynchronous FIFO has a dual-port memory written in one domain and read in the other. The write pointer is Gray-coded and synchronized into the read domain to compute "empty"; the read pointer is Gray-coded and synchronized into the write domain to compute "full". A synchronized pointer is a few cycles old, so the flags are pessimistic (the FIFO may look full or empty slightly longer than it is), which is safe.

The FIFO depth must cover rate mismatch over a burst: if a writer sends a burst of $B$ words at $f_w$ while the reader drains at $f_r < f_w$, about $B(1 - f_r/f_w)$ words accumulate, plus a few entries for synchronizer latency. CDC tools check that the pointers are Gray-coded, the synchronizers are present, and no other path crosses unguarded.

## Going deeper

Cummings, "Clock Domain Crossing (CDC) Design and Verification Techniques Using SystemVerilog" (SNUG 2008); Cummings, "Simulation and Synthesis Techniques for Asynchronous FIFO Design" (SNUG 2002).
