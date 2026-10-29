# Hardware-software co-design

A system's function can run in software on a processor or in dedicated hardware. Hardware-software co-design decides that split together with the architecture, instead of building the hardware first and fitting software on afterwards. The decision is quantitative: how much of the runtime moves to hardware, how fast that hardware is, and what it costs in area and power. Amdahl's law is the first computation every architect runs.

## Partitioning a function

Profile the software and find the kernels that dominate runtime: a filter, a matrix multiply, a cryptographic round. Each kernel is a candidate for an accelerator. Moving it to hardware gains speed and energy but costs area, design effort and flexibility, and adds communication between processor and accelerator. Co-design tools and methods (system-level models in SystemC, high-level synthesis from C++, virtual platforms that run real software on a model of the hardware) let both sides be evaluated before RTL is frozen. Getting the split wrong is expensive: an accelerator for a kernel that turns out to be 5% of runtime is wasted silicon.

## Amdahl's law

If a fraction $f$ of the original runtime is accelerated by a factor $s$ and the rest is unchanged, the overall speedup is
$$S = \frac{1}{(1 - f) + f/s}.$$
With $f = 0.8$ and $s = 10$, $S = 1/(0.2 + 0.08) = 3.57$. A program that took 100 ms now takes 28 ms. The bound as $s \to \infty$ is $1/(1 - f) = 5$: no accelerator, however fast, beats the part that stays in software. In the example, 80% acceleration bounds the speedup at 5 no matter how much hardware is spent.

## Sizing the accelerator

Invert the law to size the hardware. To reach an overall speedup $S$ with fraction $f$ accelerated, the accelerator needs
$$s = \frac{f}{1/S - (1 - f)},$$
which exists only when $S < 1/(1 - f)$. For $f = 0.8$ and a target of $S = 4$, $s = 0.8/(0.25 - 0.2) = 16$. Raising the target to 4.5 needs $s = 0.8/(0.2222 - 0.2) = 36$. Near the bound, each extra bit of speedup costs a much faster accelerator, which is a strong argument for accelerating more of the program rather than accelerating one kernel harder.

## Going deeper

De Micheli and Gupta, "Hardware/Software Co-Design", Proceedings of the IEEE (1997); Hennessy and Patterson, "Computer Architecture: A Quantitative Approach", chapter 1 on Amdahl's law.
