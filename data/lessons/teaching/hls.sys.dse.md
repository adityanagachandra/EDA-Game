# Design-space exploration of HLS directives

One C kernel admits thousands of microarchitectures through its directives: unroll factors, pipeline settings, partition factors, inlining, resource limits. Design-space exploration (DSE) searches that space for designs that trade latency against area well. The space is large, each evaluation costs a synthesis run, and the answer is a Pareto front rather than a single best point.

## The size of the directive space

Each directive choice multiplies the space. A kernel with three loops, each taking one of four unroll factors and pipelined or not, has $(4 \times 2)^3 = 512$ loop configurations; three arrays with three partitioning choices each multiply that by $3^3 = 27$, giving 13,824 points. With HLS runs taking minutes to hours, exhaustive search is out of reach even for small kernels. Many points are also illegal or pointless: a partition factor larger than the array, an unroll factor that exceeds the trip count, or a pipeline inside a fully unrolled loop. A practical explorer first prunes with cheap rules, then samples the remaining space, and spends synthesis runs only on configurations whose estimates look promising.

## Pareto-optimal designs

With two objectives to minimize, latency and area, a design dominates another when it is at least as good in both and strictly better in one. The Pareto front is the set of designs no other design dominates. Among (20 cycles, 1700 LUTs), (80, 900), (140, 1600), (40, 1100), (50, 1900), (100, 600) and (10, 1900), the design (140, 1600) is dominated by (80, 900), and (50, 1900) by (40, 1100), leaving five Pareto-optimal designs. A designer picks from the front according to the area budget or latency target; points off the front are never the right choice. Sorting by latency and sweeping while tracking the best area so far finds the front in $O(n \log n)$ for two objectives.

## Measuring a front with hypervolume

To compare explorers, measure the region a front dominates. The hypervolume indicator is the area of the set of points that are dominated by at least one front design and lie below a reference point $(L_{ref}, A_{ref})$. Sort the front by latency; between consecutive front latencies the dominated height is $A_{ref}$ minus the area of the earlier point. For a single front point (20, 100) and reference (220, 2200), the hypervolume is $200 \times 2100 = 420000$ cycle-LUTs. A larger hypervolume means a front that is closer to the ideal corner and better spread. Explorers based on Bayesian optimization or evolutionary search report hypervolume against the number of synthesis runs spent, which is the honest measure of exploration efficiency.

## Going deeper

Schafer and Wang, "High-level synthesis design space exploration: past, present, and future" (IEEE TCAD, 2020); Zitzler and Thiele, "Multiobjective evolutionary algorithms: a comparative case study and the strength Pareto approach" (IEEE TEC, 1999).
