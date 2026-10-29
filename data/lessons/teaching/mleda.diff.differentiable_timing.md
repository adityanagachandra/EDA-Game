# Differentiable static timing

A differentiable timer propagates smoothed arrival times through the timing graph so that a timing objective, such as total negative slack, has a gradient with respect to every delay and, through delay models, with respect to cell positions or sizes. It lets a gradient-based placer optimize timing directly instead of through net weights. The chain rule does the work, and the approximations decide whether it helps signoff.

## Smoothed arrival at a gate

At gate G with inputs arriving at $a_i$ through arcs of delay $d_i$, STA computes $A = \max_i (a_i + d_i)$. A differentiable timer computes $A = \mathrm{LSE}_\gamma(a_i + d_i)$. With arc sums 132 and 128 ps and $\gamma = 4$ ps, $A = 132 + 4\ln(1 + e^{-1}) = 133.25$ ps. Propagation runs in topological order exactly as in STA, so it costs one pass over the graph, and the same pass in reverse computes gradients. Arc delays themselves come from a differentiable delay model, for example Elmore delay on an estimated wire, which links $A$ to cell positions.

## The chain rule through a soft max

The derivative of $A$ with respect to an arc delay is the softmax weight of that input: $\partial A/\partial d_1 = e^{(a_1 + d_1)/\gamma}/\sum_j e^{(a_j + d_j)/\gamma}$. In the example it is 0.731. Gradients compose along paths: if an endpoint's arrival depends on $A$ with coefficient 1, it inherits the same 0.731. Backpropagation computes all such derivatives in one reverse pass, the same way a neural network's gradients are computed, which is why such timers are written as operators in a deep-learning toolkit and run on a GPU. The cost of the backward pass is the same order as the forward pass, linear in the size of the timing graph.

## Gradient of TNS

Total negative slack is $\mathrm{TNS} = \sum_j \max(0, -\mathrm{slack}_j)$, summed over endpoints. Each violating endpoint contributes the derivative of its arrival; each passing endpoint contributes zero. If G feeds two endpoints, one violating and one passing, then $\partial \mathrm{TNS}/\partial d_1 = 0.731$; if both violate, it is $2 \times 0.731 = 1.46$, since reducing $d_1$ helps both. The hinge is itself non-smooth at zero slack, so implementations smooth it too or accept a subgradient. The result guides placement; signoff STA on the routed design still measures whether timing improved. A WNS objective, by contrast, gives gradient only through the single worst endpoint unless it too is smoothed.

## Going deeper

Guo and Lin, "Differentiable-Timing-Driven Global Placement" (DAC 2022); Lin et al., "DREAMPlace" (DAC 2019).
