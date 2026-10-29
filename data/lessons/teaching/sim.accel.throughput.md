# Simulation versus emulation throughput

Verification platforms differ by orders of magnitude in speed, and they also differ in setup cost: an emulator or prototype must be compiled before it runs. Comparing them is a throughput calculation. The numbers decide which tests run where, how long a regression takes, and when a long compile is worth paying. They also explain why simulator speed per cycle, not per gate, is the figure that matters for system-level verification.

## Simulation speed in cycles per second

A software simulator's speed is measured in simulated clock cycles per second of wall time. It falls roughly in proportion to design size and activity, so a large SoC that simulates at a few hundred cycles per second is normal, while a small block may run at hundreds of thousands. Wall time for a test is cycles divided by speed. A 20-million-cycle boot sequence at 500 cycles per second takes $2 \times 10^7 / 500 = 40{,}000$ seconds, about 11.11 hours. A billion-cycle software workload would take more than three weeks. These numbers explain the push toward compiled simulation, parallel simulation and hardware acceleration: system-level scenarios are long, and simulation time is linear in their length.

## Emulation time including compile

An emulator runs much faster, in the range of a megahertz or more for large designs, but each new RTL version must be compiled first, which can take hours. Total time is compile plus run: with a 6-hour compile and a 2 MHz emulation speed, running $50 \times 10^9$ cycles takes $5 \times 10^{10} / 2 \times 10^6 = 25{,}000$ seconds, about 6.94 hours, for a total of 12.94 hours. For short tests, compile dominates; for long ones, run time does. The same compiled image can serve many tests, so the compile cost is shared across a regression, which is why emulation is organized around stable RTL snapshots rather than every edit.

## The break-even test length

Emulation beats simulation for tests longer than a break-even length $N^*$ where the times are equal: $N^*/f_s = C + N^*/f_e$, so $N^* = C / (1/f_s - 1/f_e)$. With $f_s = 200$ cycles per second, $f_e = 1$ MHz and $C = 4$ hours (14,400 s): $1/f_s - 1/f_e = 0.005 - 0.000001 = 0.004999$ s per cycle, so $N^* \approx 2.88$ million cycles. Tests shorter than about 2.9 million cycles finish sooner in simulation, longer ones in emulation. Because $f_e \gg f_s$, the formula is close to $N^* \approx C f_s$: the break-even is the number of cycles the simulator could run during the compile. This is a per-test view; across a regression, shared compiles push the break-even down further.

## Going deeper

Hauck and DeHon (eds.), "Reconfigurable Computing: The Theory and Practice of FPGA-Based Computation" (2008), the chapter on logic emulation; Wile, Goss and Roesner, "Comprehensive Functional Verification" (2005), the chapter on simulation acceleration.
