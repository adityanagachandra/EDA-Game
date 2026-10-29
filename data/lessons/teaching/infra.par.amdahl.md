# Amdahl's law and parallel speedup

EDA runs are long, and machines have many cores, so every engine author is asked how much faster the tool gets with more threads. Amdahl's law gives the first answer: the part of the run that stays serial, such as reading the database, building data structures or writing results, caps the speedup no matter how many threads are added. Knowing the law tells a tool designer where to spend effort: shrinking the serial fraction often matters more than parallelizing the kernel further.

## Speedup with a serial fraction

Split a single-threaded run of time $T_1$ into a serial fraction $s$ and a parallel fraction $1 - s$. On $p$ threads, ideally only the parallel part shrinks:
$$T_p = T_1\left(s + \frac{1 - s}{p}\right), \qquad S(p) = \frac{T_1}{T_p} = \frac{1}{s + (1 - s)/p}.$$
Example: a global-placement run is 10% serial. On 16 threads, $S = 1/(0.1 + 0.9/16) = 1/0.15625 = 6.4$. Sixteen threads give a speedup of 6.4, and the efficiency $S/p$ is only 40%. The formula is optimistic, since it ignores synchronization, memory bandwidth and load imbalance, so measured speedups are usually lower still.

## The ceiling of one over s

As $p \to \infty$ the parallel term $(1 - s)/p$ goes to zero and
$$\lim_{p \to \infty} S(p) = \frac{1}{s}.$$
Example: a single-threaded detailed-routing run takes 400 s, of which 20 s is serial: $s = 20/400 = 0.05$, so no machine can make it more than 20 times faster. With 32 threads the run reaches $1/(0.05 + 0.95/32) \approx 12.55$, already most of the way to the ceiling, and doubling to 64 threads gives only about 15.42.

The ceiling explains why engine teams attack serial phases: parsing Liberty and LEF/DEF in parallel, building the timing graph concurrently, overlapping file output with computation. Halving the serial fraction doubles the ceiling, which no number of extra cores can do.

## Serial fraction from a target speedup

Turn the law around to set a budget. If a product needs speedup $S$ on $p$ threads, solve $S = 1/(s + (1 - s)/p)$ for $s$:
$$s = \frac{p/S - 1}{p - 1}.$$
Example: a router must run 12 times faster on 16 threads. Then $s = (16/12 - 1)/15 = 0.3333/15 \approx 0.0222$: at most about 2.22% of the single-threaded run may be serial. That number becomes an engineering requirement for every phase outside the parallel kernel.

Amdahl assumes a fixed problem. Gustafson's view, that bigger machines are used for bigger designs, gives a gentler picture when the parallel work grows with the design while the serial setup does not. Both views agree that the serial part has to be measured, not guessed, before parallel work is planned.

## Going deeper

Amdahl, "Validity of the Single Processor Approach to Achieving Large Scale Computing Capabilities" (AFIPS 1967); Gustafson, "Reevaluating Amdahl's Law" (CACM 1988); Herlihy and Shavit, "The Art of Multiprocessor Programming", chapter 1.
