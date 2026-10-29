# Package substrates: organic and glass

A package substrate is the small wiring board under the die. Most are organic: a glass-fibre core with build-up dielectric layers and copper traces, made by semi-additive plating. Glass-core substrates are an emerging alternative, offering flatter panels, finer lines and a coefficient of thermal expansion that can be tuned closer to silicon. The substrate's line and space, layer count and CTE decide how many signals escape the die and how much the assembly warps.

## Wiring capacity per layer

Substrate wiring is specified by line width $w$ and space $s$. One layer carries $\lfloor 1000/(w+s) \rfloor$ traces per millimetre of edge, with $w$ and $s$ in micrometres. Organic build-up layers on mainstream substrates run around 15/15 to 8/8 um line and space; leading-edge organic and glass substrates aim lower. Across a 10 mm die edge on four signal layers at 10/10 um, capacity is $10 \times 4 \times \lfloor 1000/20 \rfloor = 2000$ traces. Wiring capacity, not bump count, is often the bottleneck: if a die needs more signals out of one edge than the substrate can carry, the designer must add layers, tighten the line and space, or move to an interposer.

## Layers needed for a signal count

Invert the capacity formula to size the stackup. If $S$ signals must cross an edge of length $E$ mm at pitch $w+s$, the number of signal layers is $\lceil S / (E \lfloor 1000/(w+s) \rfloor) \rceil$. For 3000 signals across 12 mm at 12/12 um, one layer carries $12 \times 41 = 492$, so seven layers are needed. Every build-up layer adds cost and a yield loss, and layers come in pairs on both sides of the core to keep the substrate balanced against warpage. Halving line and space roughly halves the layer count, which is the economic case for finer-line substrates and for glass, whose smooth flat surface supports finer lithography than a woven-glass organic core.

## CTE and dimensional stability

Substrate materials expand with temperature. A point at distance $D$ from the neutral point moves by $\alpha D \Delta T$, so a substrate and a die move apart by $|\alpha_{sub} - \alpha_{Si}| D \Delta T$. With an organic core at 15 ppm/K, silicon at 2.6 ppm/K, $D = 15$ mm and a 100 K swing, the mismatch is $12.4 \times 10^{-6} \times 15000 \times 100 = 18.6$ um, a large fraction of a 40 um microbump pitch. Underfill and compliant joints absorb it on flip chip, but fine-pitch die-to-die joints cannot. Glass cores can be formulated with CTE in the low single digits, much closer to silicon, and are stiffer and flatter. Their open problems are brittleness, through-glass via formation and handling of large thin panels.

## Going deeper

Tummala, "Fundamentals of Microsystems Packaging", chapter on package substrates; the IEEE Heterogeneous Integration Roadmap chapter on substrates.
