# Differentiable placement

Analytical global placement minimizes a smooth wirelength plus a density penalty over cell coordinates. DREAMPlace observed that this is the same computation as training a neural network: coordinates are the parameters, wirelength and density are the loss, and the gradients come from operators that a deep-learning toolkit runs on a GPU. The result is the same algorithm much faster, not a learned placer.

## Weighted-average wirelength

HPWL, the span of a net's pins, is not differentiable at the extreme pins and has zero gradient at interior pins. The weighted-average model replaces max and min by soft versions: $\mathrm{WA}_x = \frac{\sum x_i e^{x_i/\gamma}}{\sum e^{x_i/\gamma}} - \frac{\sum x_i e^{-x_i/\gamma}}{\sum e^{-x_i/\gamma}}$. For pins at $x = 0$, 4 and 10 um with $\gamma = 2$ um, the soft max is near 9.6 and the soft min near 0.5, and $\mathrm{WA}_x = 9.12$ um against an HPWL span of 10 um. Every pin now has a nonzero gradient, so the optimizer moves interior pins as well.

## The smoothing gap

WA underestimates the true span, and the gap shrinks with $\gamma$. For the same pins, $\gamma = 2$ gives a gap of 0.88 um and $\gamma = 1$ a gap of 0.09 um. A large $\gamma$ gives smooth, informative gradients early, when cells are far from their final places; a small one makes the objective match HPWL near the end. Placers therefore shrink $\gamma$ over the iterations, a continuation method that solves a sequence of easier problems on the way to the real one. The same idea appears in the log-sum-exp wirelength model, which overestimates instead.

## Placement as training

With loss $\mathrm{WL} + \lambda D$, where $D$ is an electrostatic density penalty, a step updates each cell by $x \leftarrow x - \eta(\partial \mathrm{WL}/\partial x + \lambda\, \partial D/\partial x)$. A cell at $x = 12$ um with gradients -1.4 and 1.9, $\lambda = 4$ and step 0.5 moves to $12 - 0.5(-1.4 + 7.6) = 8.9$ um. DREAMPlace implements wirelength and density as custom operators with hand-written backward kernels, uses Nesterov's method with adaptive steps, and raises $\lambda$ over iterations so that density increasingly dominates. The output is a global placement that still needs legalization and detailed placement.

## Going deeper

Lin et al., "DREAMPlace: Deep Learning Toolkit-Enabled GPU Acceleration for Modern VLSI Placement" (DAC 2019); Lu et al., "ePlace: Electrostatics-Based Placement Using Fast Fourier Transform and Nesterov's Method," ACM TODAES (2015).
