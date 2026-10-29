# Transfer-function moments of RC networks

Interconnect networks extracted from layout have millions of RC elements, far too many to simulate with every gate. Model order reduction replaces them by small models with nearly the same input-output behaviour. Its language is the moments of the transfer function: the Taylor coefficients of $H(s)$ around $s = 0$. The first moment is the Elmore delay; higher moments refine it.

## Expanding H(s) in powers of s

A linear RC network driven at one port obeys $(G + sC)x = bu$ in the Laplace domain, and an output $y = l^Tx$ has transfer function $H(s) = l^T(G + sC)^{-1}b$. Expanding around $s = 0$,
$$H(s) = m_0 + m_1s + m_2s^2 + \cdots.$$
For a single RC section, $H(s) = 1/(1 + s\tau)$ with $\tau = RC$, so $m_k = (-\tau)^k$. With $R = 2$ kΩ and $C = 0.5$ pF, $\tau = 1$ ns: $m_0 = 1$, $m_1 = -1$ ns, $m_2 = 1$ ns$^2$.

$m_0$ is the DC gain, $-m_1$ the mean of the impulse response (the Elmore delay), and $m_2$ relates to its spread. Matching the first few moments matches the low-frequency, and hence the slow time-domain, behaviour.

## Computing moments by repeated solves

Expanding $(G + sC)^{-1}$ gives the recursion
$$Gm_0 = b, \qquad Gm_k = -Cm_{k-1},$$
for the moment vectors at all nodes. Each moment costs one solve with the same matrix $G$, factored once.

Example: a unit source drives $R_1 = 1$ kΩ into node 1 ($C_1 = 1$ pF), then $R_2 = 2$ kΩ into node 2 ($C_2 = 0.5$ pF). At DC there is no current, so $m_0 = (1, 1)$. Solving $Gm_1 = -Cm_0$ means injecting currents $-C_1$ and $-C_2$ into the resistor tree: node 1 sees $-R_1(C_1 + C_2) = -1.5$ ns, node 2 sees $-1.5 - R_2C_2 = -2.5$ ns. So $m_1 = (-1.5, -2.5)$ ns, the negated Elmore delays, from one solve.

## Moments, Elmore delay and a one-pole fit

A one-pole model $H(s) \approx m_0/(1 + s\tau)$ has moments $m_0, -m_0\tau, m_0\tau^2, \dots$. Matching $m_0$ and $m_1$ gives $\tau = -m_1/m_0$, the Elmore estimate. Matching $m_1$ and $m_2$ instead gives
$$\tau = -\frac{m_2}{m_1}.$$
For node 2 above, $m_2$ comes from one more solve: $m_2 = (2.75, 5.25)$ ns$^2$, so $\tau = 5.25/2.5 = 2.1$ ns, against the Elmore 2.5 ns. The two-moment fit captures the dominant pole better; the Elmore value is the classical upper bound on the 50% delay of an RC tree.

Asymptotic waveform evaluation (AWE) matches $2q$ moments with $q$ poles through a Padé approximant. It works for small $q$ but becomes numerically unreliable beyond about eight poles, because explicitly computed moments are nearly linearly dependent. Krylov projection fixes that.

## Going deeper

Pillage and Rohrer, "Asymptotic waveform evaluation for timing analysis", IEEE Trans. CAD (1990); Antoulas, "Approximation of Large-Scale Dynamical Systems" (SIAM, 2005), chapter 11.
