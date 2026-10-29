# The processor performance equation

Execution time factors into instruction count, cycles per instruction (CPI) and clock period. Every microarchitecture and implementation decision moves one of these terms, often at the expense of another, and the equation is how the effect of a timing closure result (a faster clock) is weighed against the effect of an architectural choice (fewer stalls). Amdahl's law bounds what any single improvement can buy.

## Time equals count times CPI times period

$$T_{exec} = IC \times CPI \times T_{clk} = \frac{IC \times CPI}{f}.$$

A program of $2 \times 10^{9}$ instructions with CPI 1.4 on a 2.5 GHz clock takes $2\times10^9 \cdot 1.4 / 2.5\times10^9 = 1.12$ s. The three factors belong to different owners: the compiler and instruction set fix $IC$, the microarchitecture fixes $CPI$, and the circuit implementation and timing closure fix $f$. A deeper pipeline raises $f$ but usually raises $CPI$ through more stall cycles, so only the product decides.

## CPI from an instruction mix

Different instruction classes take different numbers of cycles. The average CPI is the frequency-weighted mean:

$$CPI = \sum_k f_k \cdot CPI_k.$$

With 50% ALU instructions at 1 cycle, 30% loads and stores at 2 cycles and 20% branches at 3 cycles, $CPI = 0.5 + 0.6 + 0.6 = 1.7$. Measurements of CPI by class come from simulation or hardware counters on representative workloads; architectural simulators such as gem5 produce exactly these numbers before any RTL exists.

## Amdahl's law

If an improvement speeds up a fraction $p$ of the execution time by a factor $s$ and leaves the rest unchanged, the overall speedup is

$$S = \frac{1}{(1 - p) + p/s}.$$

Accelerating 60% of the time by a factor of 10 gives $1 / (0.4 + 0.06) \approx 2.17$, far below 10, and even an infinite $s$ caps the speedup at $1/(1 - p) = 2.5$. The same law governs EDA flows: parallelizing a router step that is 70% of runtime across many cores cannot make the flow more than about 3.3 times faster, which is why serial bottlenecks such as timing updates draw so much engineering attention.

## Going deeper

Hennessy and Patterson, "Computer Architecture: A Quantitative Approach" (6th ed., 2017), chapter 1; Amdahl, "Validity of the Single Processor Approach to Achieving Large Scale Computing Capabilities" (AFIPS, 1967).
