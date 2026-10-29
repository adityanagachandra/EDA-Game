# Golden and differential test harnesses

An agent that edits engines needs an automatic judge. A golden harness reruns a fixed suite of designs and compares the candidate build's metrics with stored baselines; a differential check proves that a change moved only what it was meant to move. This concept covers the arithmetic of those comparisons: relative change, tolerance gates and counting unintended movement, the numbers on which a promotion decision rests.

## Relative change against a golden baseline

For a metric such as total wirelength or cell area, the harness reports the relative change $\Delta = (x_{cand} - x_{base})/x_{base} \times 100\%$. Example: baseline area 29,900 square microns, candidate 29,617: $\Delta = -283/29900 \times 100 = -0.95\%$. The sign matters: for area and wirelength smaller is better, for slack larger is better, so a harness stores the preferred direction with each metric. Relative change makes designs of different sizes comparable, which is why it is the input to geometric-mean QoR summaries. For metrics near zero, such as WNS close to 0 ps, relative change explodes, so the harness uses absolute differences in ps instead. Golden values are stored with the engine version and seed that produced them, and are regenerated only on a reviewed baseline update.

## Counting tolerance violations

A QoR gate flags a design when its metric moves outside a tolerance band, in either direction, because an unexplained improvement is as suspicious as a regression. With a 1% band on wirelength and designs aes 76,900 → 78,669 (+2.30%), jpeg 85,500 → 86,782 (+1.50%), ibex 57,700 → 57,758 (+0.10%), uart 41,800 → 40,253 (−3.70%) and fft 84,700 → 82,921 (−2.10%), four designs are flagged: all but ibex. The band should sit above the engine's seed noise; if reruns of the baseline already vary by 0.8%, a 0.5% band flags noise. That is why deterministic engines matter: with bit-identical reruns the band can be tight, and any movement at all is real. Flagged designs go to the model or the human as feedback, with the per-design numbers.

## Proving only intended designs moved

A differential check asks a sharper question than a gate: did the change affect exactly what it should? A new derate plug-in meant to touch only fft is run on the suite, and WNS is compared with an absolute tolerance of 2 ps. Results: riscv32i −35 → −26, aes −37 → −33, gcd −12 → −10, fft −49 → −30, ibex −2 → −7. fft moved by 19 ps, as intended. Among the others, riscv32i (9 ps), aes (4 ps) and ibex (5 ps) exceed 2 ps, while gcd (2 ps) does not. Three unintended movers mean the plug-in's condition fires more broadly than claimed, so the change is rejected or its rule narrowed. The same idea applies at finer grain: compare per-path slacks and require that only paths through the targeted cells changed.

## Going deeper

The agent-first notes in this curriculum on golden and differential testing; the metrics and regression scripts of OpenROAD-flow-scripts.
