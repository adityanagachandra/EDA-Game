# Silicon debug and post-silicon validation

First silicon runs real software at full speed, billions of cycles per second, and finds bugs that pre-silicon simulation never reached. Post-silicon validation hunts those design and electrical bugs; silicon debug localizes them with very little visibility into the chip. Its tools, trace buffers, triggers, scan dumps and JTAG-based access, are design-for-debug structures that reuse much of the test infrastructure.

## Validation versus manufacturing test

Manufacturing test assumes the design is correct and screens each die for defects; its patterns are structural and its question is whether this die was built right. **Post-silicon validation** assumes the dies are mostly good and asks whether the design works: it runs operating systems, stress tests and random instruction streams across voltage, temperature and frequency. Failures that appear on every die point to design bugs; failures that depend on conditions, mapped with **shmoo plots** of pass and fail over voltage and frequency, point to speed paths, noise or power integrity. A validation escape is a bug in every part; a test escape is a defective part that passed test.

## Seeing inside: traces and triggers

The central difficulty is **observability**: a few internal signals can be watched while the chip runs at speed. An on-chip **trace buffer** records selected signals each cycle after a programmable **trigger** fires on a condition such as a bus value. Its capacity limits the window. With 1024 entries of 32 bits and 24 traced signals, each cycle's sample fits in one entry and 1024 cycles are captured; tracing 70 signals in 32-bit entries needs 3 entries per cycle and captures only 341. **Trace signal selection** chooses signals from which many untraced flop values can be reconstructed, and compression widens the window further.

## Scan dumps and reproduction

A **scan dump** stops the functional clocks at a trigger and shifts out all scan chains, giving the complete state at one instant at the cost of destroying execution. Combining dumps at successive stop points with traces localizes the failing window; then the window is reproduced in simulation or emulation, where visibility is complete, to find the root cause. Non-determinism from asynchronous clock crossings can stop a failure from reproducing, so debug flows include deterministic modes. Fixes use spare cells and metal-only ECOs or microcode and firmware patches when they exist, avoiding a full respin.

## Going deeper

Mishra, Morad, Ludwig and others, "Post-Silicon Validation and Debug" (2019); Abramovici et al., "A reconfigurable design-for-debug infrastructure for SoCs," Design Automation Conference, 2006.
