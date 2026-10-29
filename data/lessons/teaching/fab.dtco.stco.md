# System-technology co-optimization and chiplets

System-technology co-optimization asks how to split a system across dies and packages. Splitting a large die into chiplets raises yield, since yield falls exponentially with area, but adds interface area, assembly cost and assembly yield loss. The trade is computed with the same yield and cost arithmetic as a single die, and its answer depends on defect density and die size.

## The monolithic baseline

Start with a monolithic die of area $A$. On a 300 mm wafer costing $C_w$, with dies per wafer from the edge-loss formula and Poisson yield $e^{-A D_0}$, the cost per good die is $C_w/(DPW \cdot Y)$. A 600 mm^2 die has $DPW = 117.8 - 27.2 = 90$ and, at 0.2 defects per cm^2, $Y = e^{-1.2} = 0.301$. At 15,000 dollars per wafer that is $15000/(90 \times 0.301) = 553.4$ dollars per good die. Large dies on a young process are punished by yield, which is the motivation for splitting.

## Splitting into chiplets

Split the design into $k$ chiplets, each $A/k$ plus an overhead fraction for die-to-die interfaces. With $k = 4$ and 10% overhead, each chiplet is 165 mm^2, giving 376 per wafer and $Y = e^{-0.33} = 0.719$, so a good chiplet costs $15000/(376 \times 0.719) = 55.5$ dollars. If chiplets are tested before assembly (known good die), a system needs four good chiplets and survives assembly with yield, say, 0.97: $4 \times 55.5/0.97 = 228.8$ dollars. The silicon cost falls by a factor of $553.4/228.8 = 2.42$, before counting the extra cost of the package and the interface power.

## When splitting pays

The advantage grows with die area and defect density, and shrinks as the process matures: at 0.05 defects per cm^2 the same split saves far less, and the package cost may dominate. Overhead area, the cost of advanced packaging, die-to-die latency and energy per bit, and test cost all push the other way. STCO also lets each chiplet use the best technology, for example logic on a leading node and I/O or analog on an older, cheaper one. For tools, chiplets add partitioning, die-to-die timing, multi-die extraction and thermal analysis across a package.

## Going deeper

Stow, Xie, Siddhu and Loh, "Cost-effective design of scalable high-performance systems using active and passive interposers", ICCAD (2017); Hennessy and Patterson, "Computer Architecture: A Quantitative Approach", chapter 1.
