# Most probable failure point

Where in parameter space does a circuit most likely fail? In standard-normal coordinates the answer is the point of the failure region closest to the origin, because the Gaussian density falls off with distance. Finding that point turns a probability question into an optimization, gives a first-order estimate of the failure rate, and tells importance sampling where to aim.

## Distance to a linear boundary

Whiten the parameters so they are independent standard normal: $x_i = \Delta p_i/\sigma_i$. Suppose failure occurs when $a \cdot x > b$, a hyperplane. The most probable failure point minimizes $\|x\|^2$ subject to $a \cdot x = b$; a Lagrange multiplier gives $x^* = b\,a/\|a\|^2$, at distance
$$\beta = \frac{b}{\|a\|}$$
from the origin. $\beta$ is called the reliability index. For $3x_1 + 4x_2 > 25$, $\|a\| = 5$ and $\beta = 5$: the cell fails at a 5-sigma point in the direction $(0.6, 0.8)$. Every other failure point is farther away and individually less likely.

## The failure probability

For a linear boundary, the projection $a \cdot x$ of a standard normal vector is $N(0, \|a\|^2)$, so the failure probability is exactly
$$P = P(a \cdot X > b) = Q(\beta).$$
In the example, $P = Q(5) = 2.87 \times 10^{-7}$. For a curved boundary, $Q(\beta)$ is the first-order reliability method (FORM) estimate, exact for a hyperplane and approximate otherwise. A convex failure region (curving away from the origin) has less probability than $Q(\beta)$; a concave one has more. The estimate also misses separate failure regions entirely: a bitcell can fail in one direction by read upset and in another by write failure, and each region has its own most probable point.

## Scaling and the search in practice

Physical parameters are not unit-variance. With $\Delta V_1$ and $\Delta V_2$ of sigmas $\sigma_1$ and $\sigma_2$ and failure when $a_1\Delta V_1 + a_2\Delta V_2 > b$, substitute $\Delta V_i = \sigma_i x_i$ to get coefficients $a_i\sigma_i$ and $\beta = b/\sqrt{a_1^2\sigma_1^2 + a_2^2\sigma_2^2}$. Forgetting to whiten weights the directions wrongly and finds the wrong point. A real circuit has no closed-form boundary, so engines search for the point with gradient or sampling-based optimization, each evaluation a SPICE run, or fit a surrogate of the margin first. The point found then becomes the centre of an importance-sampling proposal, which corrects the curvature error that FORM alone leaves.

## Going deeper

Ditlevsen and Madsen, "Structural Reliability Methods" (1996), chapters on FORM; Dolecek, Qazi, Shah and Chandrakasan, "Breaking the simulation barrier: SRAM evaluation through norm minimization," ICCAD 2008.
