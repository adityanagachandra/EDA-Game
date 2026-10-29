# The gradient of log-sum-exp wirelength

A nonlinear placer moves each cell along the negative gradient of the objective, so it needs $\partial W/\partial x_i$ for every pin of every net, millions of times per iteration. For log-sum-exp wirelength the gradient is a difference of two softmax weights, a closed form that is cheap, bounded and easy to parallelise.

## Softmax weights

Differentiate $\gamma \ln \sum_j e^{x_j/\gamma}$ with respect to $x_i$: the chain rule gives
$$\frac{\partial}{\partial x_i}\,\gamma \ln \sum_j e^{x_j/\gamma} = \frac{e^{x_i/\gamma}}{\sum_j e^{x_j/\gamma}},$$
the softmax weight of pin $i$. These weights are positive and sum to 1, and the pins nearest the maximum get most of the weight. For x = 7, 5, 1, 7 with $\gamma$ = 5, the weight of the pin at 5 is $e^{1}/(e^{1.4}\cdot 2 + e^{1} + e^{0.2}) = 0.2256$. Interior pins get small but nonzero weights, unlike HPWL, whose gradient is nonzero only for the extreme pins.

## The full gradient of a net

The minimum term contributes a negative softmax over the negated coordinates, so
$$\frac{\partial W^{LSE}}{\partial x_i} = \frac{e^{x_i/\gamma}}{\sum_j e^{x_j/\gamma}} - \frac{e^{-x_i/\gamma}}{\sum_j e^{-x_j/\gamma}}.$$
Each component lies between -1 and 1. The rightmost pin has a gradient near +1, so descent pulls it left; the leftmost has a gradient near -1 and is pulled right; middle pins feel little. For x = 4, 10, 4 with $\gamma$ = 4, the pin at 4 gets $0.1543 - 0.4498 = -0.2955$. Summed over all pins of a net the gradient is exactly 0, because translating the whole net leaves its length unchanged: a useful unit test for any implementation.

## From pins to cells

A cell's gradient is the sum of the gradients of its pins over all nets it belongs to, each multiplied by the net weight; pin offsets inside the cell do not change the derivative with respect to the cell's position. A cell on two nets with pin gradients 0.12 and 0.15 has total gradient 0.27. Implementations compute the two exponential sums per net in one pass, then each pin's gradient in a second pass, so the cost is linear in the number of pins. On GPUs, one thread per net or per pin works well, with the per-cell sum done by atomic adds or a segmented reduction. The exponentials must be computed after subtracting the net's maximum (and for the other term its minimum) to avoid overflow when $\gamma$ is small.

## Going deeper

Lu et al., "ePlace: electrostatics-based placement using fast Fourier transform and Nesterov's method" (ACM TODAES, 2015); Lin et al., "DREAMPlace" (DAC 2019), on GPU gradient kernels.
