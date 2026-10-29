# Nonlinear placement with a density penalty

Nonlinear analytical placement minimises smooth wirelength plus a weighted density penalty, $W(x) + \lambda D(x)$, by gradient methods, raising $\lambda$ until the cells are spread. APlace, NTUplace3, mPL and the electrostatic placers all follow this penalty method. Its behaviour is set by three things: the gradient step, the schedule of $\lambda$ and the initial balance between the two terms.

## One step on the penalised objective

For fixed $\lambda$, a placer takes gradient-type steps on $f(x) = W(x) + \lambda D(x)$. The gradient on a cell combines the wirelength pull, from the smooth model, and the density push, from the derivative of the smoothed bin densities: $\nabla f = \nabla W + \lambda \nabla D$. With step size $\eta$, plain gradient descent moves a cell to $x - \eta(\partial W/\partial x + \lambda\,\partial D/\partial x)$. A cell at x = 21 with $\partial W/\partial x = -16$, $\partial D/\partial x = -20$, $\lambda = 1$ and $\eta = 0.2$ moves to $21 + 0.2 \cdot 36 = 28.2$: both terms push it right. Real placers use conjugate gradient (NTUplace3) or Nesterov's method (ePlace) rather than plain descent, but the combined gradient is the same.

## Raising the penalty weight

Small $\lambda$ lets wirelength dominate and cells cluster; large $\lambda$ spreads them. The penalty method starts small and raises $\lambda$ over outer iterations, re-solving or continuing the inner optimisation each time, until overflow reaches its target. NTUplace3 doubles $\lambda$ per outer iteration; electrostatic placers multiply it by a factor slightly above 1 every iteration, adjusted by how much wirelength changed. With $\lambda_0 = 0.001$ and factor 1.1, reaching 1 takes $\lceil \ln 1000/\ln 1.1 \rceil = 73$ iterations. Raising $\lambda$ too fast freezes a poor arrangement; too slowly wastes runtime with cells hardly moving.

## Balancing the two forces at the start

The initial $\lambda$ should make the wirelength and density forces comparable, otherwise one term is ignored. A standard initialisation is the ratio of gradient norms, $\lambda_0 = \|\nabla W\|_1/\|\nabla D\|_1$, computed at the initial placement. If the wirelength gradients of four cells are 26, -26, -8, 21 (norm 81) and the density gradients are -263, 120, -76, 241 (norm 700), then $\lambda_0 = 0.1157$. The initial placement itself matters: placers often start from a quadratic solution so that the nonlinear phase begins with sensible relative positions. All of these numbers, the initial weight rule, the growth factor and the overflow target, are parameters an agent can tune without touching the solver.

## Going deeper

Chen et al., "NTUplace3: an analytical placer for large-scale mixed-size designs" (IEEE TCAD, 2008); Kahng and Wang, "Implementation and extensibility of an analytic placer" (IEEE TCAD, 2005).
