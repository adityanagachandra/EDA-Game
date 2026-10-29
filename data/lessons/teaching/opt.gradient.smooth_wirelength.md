# Smoothing nonsmooth objectives

HPWL, the half-perimeter wirelength, is convex but not differentiable: it is a max minus a min of pin coordinates. Gradient-based placers need a smooth surrogate, and the standard choices are log-sum-exp (LSE) and weighted-average (WA) models with a smoothing parameter gamma. The parameter trades accuracy against smoothness, and its schedule is one of the most consequential knobs in an analytical placer.

## The log-sum-exp maximum

The smooth maximum of pin coordinates $x_1, \dots, x_n$ is
$$\text{LSE}_\gamma(x) = \gamma \ln \sum_i e^{x_i / \gamma}.$$
It is convex, infinitely differentiable, and always at least $\max_i x_i$. Its gradient is a softmax: pin i gets weight $e^{x_i/\gamma}/\sum_j e^{x_j/\gamma}$, so the extreme pin dominates while nearby pins share the pull.

Example: pins at 16, 10 and 35 um with $\gamma = 7$ um. The exponentials are $e^{2.286} = 9.83$, $e^{1.429} = 4.17$ and $e^{5} = 148.4$, summing to 162.4, so $\text{LSE} = 7 \ln 162.4 = 35.63$ um, slightly above the true maximum 35. To avoid overflow, implementations subtract the maximum first: $\max_i x_i + \gamma \ln \sum_i e^{(x_i - \max) / \gamma}$.

## LSE wirelength against HPWL

The LSE wirelength of a net in x is a smooth max plus a smooth version of minus the min:
$$W_\gamma = \gamma \ln \sum_i e^{x_i/\gamma} + \gamma \ln \sum_i e^{-x_i/\gamma}.$$
For pins at 19, 6 and 24 with $\gamma = 6$: the first term is 26.37 and the second is -5.087, so $W_\gamma = 21.28$ um against an HPWL of 18. The y direction is handled the same way, and the total is summed over nets.

Every pin receives a gradient, not only the two extreme ones, which spreads the force smoothly as pins cross each other. The weighted-average (WA) model, $\sum x_i e^{x_i/\gamma} / \sum e^{x_i/\gamma}$ minus its mirror for the min, is also smooth and never above HPWL, and for the same gamma it is usually closer to HPWL. NTUplace3 uses LSE; ePlace and its successors use WA.

## Error bounds and the gamma schedule

Since $\max_i x_i \le \text{LSE}_\gamma \le \max_i x_i + \gamma \ln n$ for n pins, the LSE wirelength overestimates HPWL by at most $2\gamma \ln n$ per net. For a 3-pin net with $\gamma = 3$ that is $6 \ln 3 = 6.59$ um. To guarantee an error of at most $\varepsilon$ on every n-pin net, choose $\gamma \le \varepsilon / (2\ln n)$; for 4-pin nets within 2 um, $\gamma \le 0.72$ um.

Small gamma is accurate but makes the objective nearly nonsmooth: the gradient changes abruptly and the curvature, roughly proportional to $1/\gamma$, makes the problem stiff and the stable step small. Placers therefore start with gamma comparable to a few bin widths and shrink it as density overflow falls, so early iterations move cells smoothly and late iterations see an accurate wirelength. In ePlace and RePlAce gamma is computed from the current density overflow.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 4; Chen et al., "NTUplace3: an analytical placer for large-scale mixed-size designs" (IEEE TCAD, 2008), for log-sum-exp in practice.
