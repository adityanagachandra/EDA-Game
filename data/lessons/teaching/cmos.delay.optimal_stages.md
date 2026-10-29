# The optimal number of stages

When a small gate must drive a large load, adding inverter stages can make the path faster, up to a point. The optimal number of stages balances fewer, harder-working stages against the parasitic delay each extra stage costs. It sets the length of clock and output drivers and the shape of buffer trees that synthesis and physical tools build every day.

## The trade-off in one function

For a chain of $N$ inverters delivering path effort $F$ (here $F = H$, the load over the first input), the delay is
$$D(N) = N F^{1/N} + N p_{inv}.$$
Few stages means a large stage effort $F^{1/N}$; many stages means a small effort but a parasitic delay $p_{inv}$ paid at every stage. With $F = 256$ and $p_{inv} = 1$: $D(2) = 2 \cdot 16 + 2 = 34$, $D(4) = 4 \cdot 4 + 4 = 20$, $D(6) = 6 \cdot 2.52 + 6 = 21.1$. The minimum is near 4 stages. Adding inverters to an existing path is legal only in pairs if the polarity must be preserved, which is a discrete constraint the formula does not see.

## The optimal stage effort

Treat $N$ as continuous and set $dD/dN = 0$. With $\rho = F^{1/N}$ the condition becomes
$$\rho(1 - \ln \rho) + p_{inv} = 0.$$
For $p_{inv} = 0$ the solution is $\rho = e \approx 2.72$; for $p_{inv} = 1$ it is $\rho \approx 3.59$. The best number of stages is then $\hat N = \log_\rho F$. The rule of thumb $\hat N \approx \log_4 F$ uses a stage effort of 4, slightly more than the optimum. With $F = 256$ it gives exactly 4. The delay curve is flat near its minimum: stage efforts anywhere from about 2.4 to 6 stay within roughly 15 percent of the best, so rounding $N$ to a nearby integer, or to an even number for polarity, loses little.

## Choosing the integer and its uses

Because $N$ must be an integer, evaluate $D$ at the integers around $\log_\rho F$ and keep the smallest. For $F = 1000$ with $p_{inv} = 1$, $\log_4 1000 = 4.98$; $D(5) = 5 \cdot 3.98 + 5 = 24.9$ and $D(6) = 6 \cdot 3.16 + 6 = 25.0$, so 5 stages win by a hair. The flatness of the optimum is good news for tools: buffer-insertion and clock-tree algorithms can choose stage counts and sizes from a small library without needing the exact optimum. It also warns against over-buffering. A path with stage efforts of 1.5 has too many stages, and removing a pair of inverters makes it faster; timing optimizers check for exactly this when they consider buffer removal.

## Going deeper

Sutherland, Sproull and Harris, "Logical Effort: Designing Fast CMOS Circuits" (1999), chapter 3; Weste and Harris, "CMOS VLSI Design" (4th ed.), section 4.5.4.
