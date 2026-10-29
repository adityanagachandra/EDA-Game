# Regression testing and QoR tracking

An EDA engine is never finished: every heuristic tweak, refactor or compiler upgrade can change results. Regression testing runs a fixed suite of designs through the tool on every change and compares quality of results (QoR), such as wirelength, worst slack, area, power and run time, against a stored baseline. The comparison must summarize many designs of very different sizes into one honest number, and the standard summary is the geometric mean of per-design ratios.

## Geometric mean of QoR ratios

For each design $i$, form the ratio $r_i = \text{new}_i / \text{old}_i$ of the candidate's metric to the baseline's. A ratio is scale free: a 1% change on a tiny design and on a huge one count equally, which raw differences would not. Summarize $n$ ratios by their geometric mean,
$$G = \left(\prod_{i=1}^{n} r_i\right)^{1/n} = \exp\left(\frac{1}{n}\sum_i \ln r_i\right).$$
Example: routed wirelength on three designs goes 1000 → 900, 500 → 550 and 2000 → 2000. The ratios are 0.9, 1.1 and 1.0. Their product is 0.99, so $G = 0.99^{1/3} \approx 0.997$. Computing through logarithms avoids overflow and underflow when the suite has hundreds of designs.

## Reporting a percent change

A geometric-mean ratio is reported as a signed percentage, $(G - 1) \times 100$. For the example, $(0.9967 - 1) \times 100 \approx -0.33\%$: for a lower-is-better metric, a 0.33% improvement on average. Regression dashboards show this per metric and per commit, so a slow drift becomes visible long before any single design looks broken.

A summary must not hide outliers. A suite whose geometric mean improves by 1% while one design gets 15% worse has a real problem, so dashboards also show the per-design ratios and flag any design beyond a tolerance. Run time needs its own baseline machine and load, since it varies with the host, and metrics with mixed directions (lower is better for area, higher for slack margin) should be put on a common orientation before they are combined. Worst slack, which can be zero or negative, is better tracked as a difference than as a ratio.

## Why the geometric mean is symmetric

Compare the candidate with the baseline, or the baseline with the candidate: the conclusion should be the same. The geometric mean has exactly this property, because the geometric mean of the reciprocals is the reciprocal of the geometric mean:
$$\left(\prod_i \frac{1}{r_i}\right)^{1/n} = \frac{1}{G}.$$
For the example, baseline over candidate gives $1/0.9967 \approx 1.003$, the same 0.33% read from the other side.

The arithmetic mean fails this test. Take two designs, one twice as large and one half as large: ratios 2 and 0.5. Their arithmetic mean is 1.25, which says the candidate is 25% worse, and the arithmetic mean of the reciprocals, 0.5 and 2, is also 1.25, which says the baseline is 25% worse. Both cannot be true. The geometric mean gives 1 both ways, the correct reading of a gain and a loss that cancel in ratio terms.

## Going deeper

Fleming and Wallace, "How Not to Lie with Statistics: The Correct Way to Summarize Benchmark Results" (CACM 1986); the OpenROAD-flow-scripts documentation on its metrics and regression checks.
