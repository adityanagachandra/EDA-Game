# Emulation and FPGA prototyping

Software simulation of a full SoC runs far too slowly to boot an operating system or run real workloads. Hardware-assisted verification runs the design on hardware instead: **emulators** built from custom processors or FPGAs, and **FPGA prototypes**. They change what verification can reach, from millions of cycles to billions, and they impose their own constraints on the design, the testbench and debugging.

## Emulators and prototypes

An **emulator** compiles the synthesizable design onto dedicated hardware, either arrays of special-purpose processors that evaluate logic in parallel or large arrays of FPGAs, and runs it orders of magnitude faster than a software simulator. Emulators are built for verification: they capture signal traces deeply for debug, support fast recompiles after a change, and serve many users from a shared machine. An **FPGA prototype** maps the design onto commercial FPGA boards and runs faster still, closer to real system speed, but with limited visibility into internal signals and longer bring-up, because the design must be partitioned, constrained and timed for the FPGAs. Both require synthesizable RTL: behavioral delays, unsynthesizable testbench code and four-state X semantics have no direct hardware meaning.

## Testbenches and transactors

A conventional testbench runs on a host computer. If the emulated DUT had to exchange signal values with the host on every clock, the link latency would cap its speed near that of simulation. The standard solution is **transaction-based acceleration**, or co-modeling: the pin-level parts of the testbench, the drivers and monitors, are written in synthesizable form as **transactors** that run inside the emulator, and only transactions cross the link to the host-side testbench, in batches. The host still runs sequences, scoreboards and coverage, so the verification environment is largely reused. **Speed bridges** connect real external devices, such as a network or memory interface running at full speed, to the slower emulated design by buffering and rate-matching their traffic. Debug visibility is limited by trace storage, so captures are triggered and windowed.

## Choosing the platform

Each platform fits a stage. **Simulation** gives instant turnaround and complete visibility, so it is the right tool while RTL changes daily and blocks are being debugged. **Emulation** runs fast enough to boot an operating system and run driver software on the full SoC before tape-out, while still offering hardware-assisted debug, and it accelerates long regressions. **FPGA prototypes** suit firmware and software teams who need many copies of the design running near real speed for weeks. Speed does not replace checking: an emulator running a billion cycles finds nothing unless scoreboards, assertions and coverage judge what happened. Teams therefore carry the simulation environment forward and measure coverage on every platform.

## Going deeper

Hauck and DeHon (eds.), "Reconfigurable Computing: The Theory and Practice of FPGA-Based Computation" (2008), the chapters on logic emulation and multi-FPGA systems; the Accellera SCE-MI standard (Standard Co-Emulation Modeling Interface) for transaction-level co-emulation.
