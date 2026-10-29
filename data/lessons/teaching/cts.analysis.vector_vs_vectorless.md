# Vector-based and vectorless power analysis

Power analysis needs activity. Vector-based analysis takes it from simulation of real workloads; vectorless analysis assumes it at the inputs and registers and propagates it through the logic. Both are used throughout a project, and they answer different questions: vector-based analysis estimates the power of a known scenario, while vectorless analysis covers scenarios nobody simulated. The same split appears in dynamic IR analysis, where the worst switching pattern matters more than the average one.

## Two sources of activity

Vector-based analysis reads activity from simulation output. A VCD (or similar waveform format) records every time-stamped value change; a SAIF file aggregates per-net statistics, such as toggle counts and the time spent at 0 and 1, which is enough for average power and much smaller. Vectorless analysis assumes default toggle rates and probabilities at primary inputs and register outputs, then propagates them through the logic with probabilistic models. It needs only the netlist and libraries, so it runs early, before testbenches exist. Clock nets are the exception that needs neither: their activity follows from the frequency.

## Power from toggle counts

Each rising transition of a net draws $CV^2$ from the supply and the falling one draws none, so on average each toggle costs $\tfrac{1}{2}CV^2$. Over a window $T$ with $N$ toggles, average switching power is $P = \tfrac{1}{2}NCV^2/T$. A net that toggled 47 000 times in 250 µs, driving 8 fF at 0.75 V, averages $\tfrac{1}{2} \cdot 47\,000 \cdot 8 \cdot 0.5625/250$ fJ per µs, which is 0.423 µW. Tools sum this over nets and add internal cell power from Liberty tables and leakage. Forgetting the half double-counts switching energy.

## Coverage and its pitfalls

Vector-based results describe only the simulated workload, and only the nets that were annotated: after synthesis many RTL names no longer match the netlist, and unannotated nets fall back to propagated or default activity, so a report with low annotation coverage is partly vectorless. Zero-delay simulation misses glitch power. Average power over a long test hides short high-activity windows, which drive dynamic IR and di/dt. For IR signoff, vectorless methods search for a pessimistic simultaneous-switching scenario that the testbench may never reach. Long realistic workloads are often captured in emulation and mapped onto the netlist.

## Going deeper

Najm, "A survey of power estimation techniques in VLSI circuits", IEEE Transactions on VLSI Systems (1994); IEEE Std 1364 for the VCD format and the SAIF specification for switching activity files.
