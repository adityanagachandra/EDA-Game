# Maximum clock frequency

The fastest clock a synchronous design can use is set by its slowest register-to-register path. Computing it is the setup check solved for the period. Every architectural decision downstream (pipelining, retiming, adder choice) is ultimately judged by how it moves this number, and a timer's headline result, worst negative slack at a target frequency, is the same computation run in reverse.

## Minimum period of one path

Setting setup slack to zero gives the minimum period of a path:

$$T_{\min} = t_{cq} + t_{pd} + t_{su}, \qquad f_{\max} = \frac{1}{T_{\min}}.$$

With $t_{cq} = 50$ ps, logic delay $t_{pd} = 700$ ps and $t_{su} = 30$ ps, $T_{\min} = 780$ ps and $f_{\max} = 1/(780 \times 10^{-12}\ \text{s}) \approx 1282.1$ MHz. A handy conversion: $f$ in MHz is $10^6 / T$ with $T$ in ps. The register overhead $t_{cq} + t_{su}$ is paid on every path regardless of how little logic it contains.

## The worst path sets the clock

A design has many register-to-register paths, and they share one clock. The period must satisfy every path, so

$$T_{\min} = \max_{\text{paths}} (t_{cq} + t_{pd} + t_{su}).$$

If paths have logic delays 420, 610 and 580 ps with $t_{cq} = 60$ and $t_{su} = 40$ ps, the 610 ps path is critical, $T_{\min} = 710$ ps and $f_{\max} \approx 1408.5$ MHz. Speeding up the other two paths changes nothing; only the critical path, and any path within a margin of it, is worth optimizing. This is why timing reports sort paths by slack, and why a design whose top thousand paths are all nearly critical (a "timing wall") is much harder to speed up than one with a single long path.

## Skew and the effective period

With clock skew $s$ (capture clock later than launch by $s$), the constraint becomes $T \ge t_{cq} + t_{pd} + t_{su} - s$. Positive skew on a critical path lowers its required period; negative skew raises it.

For the 700 ps path above with $s = 40$ ps, $T_{\min} = 740$ ps, $f_{\max} \approx 1351.4$ MHz. The borrowed 40 ps is not free: the capture flip-flop launches the *next* stage 40 ps late, so the downstream path loses 40 ps, and every short path into the capture flop loses 40 ps of hold margin. Useful-skew scheduling is an optimization over all these constraints at once, a linear program over clock arrival times that tools solve during clock tree synthesis.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), section 3.5; Fishburn, "Clock Skew Optimization" (IEEE Transactions on Computers, 1990).
