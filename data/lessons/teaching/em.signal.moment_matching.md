# Moments, Elmore delay and moment matching

Moments are the Taylor coefficients of a transfer function around $s = 0$. The first moment of an RC tree is the Elmore delay; matching the first few moments with a low-order rational function is the idea behind asymptotic waveform evaluation (AWE) and later Krylov-based model order reduction. Extractors and delay calculators use moments because they are cheap to compute on a tree and summarize its waveform.

## Moments of a transfer function

Expand $H(s) = m_0 + m_1 s + m_2 s^2 + \cdots$. Because $H(s) = \int_0^\infty h(t)e^{-st}dt$, the moments are $m_k = \frac{(-1)^k}{k!}\int_0^\infty t^k h(t)\,dt$. For an RC tree driven by an ideal source, $m_0 = 1$ at every node and
$$-m_1 = \sum_k R_{ik} C_k,$$
the Elmore delay, with $R_{ik}$ the resistance shared by the paths from the source to node $i$ and to capacitor $k$.

Example: $R_1 = 1$ kohm to node 1 with $C_1 = 20$ fF, then $R_2 = 2$ kohm to node 2 with $C_2 = 30$ fF. At node 1, $-m_1 = 1 \cdot (20 + 30) = 50$ ps. At node 2, $-m_1 = 1 \cdot 20 + 3 \cdot 30 = 110$ ps. The first moment is the mean of the impulse response, so it is a delay; nodes farther from the source have larger moments.

## The second moment

Higher moments follow a recursion on the tree. Each capacitor's current at order $k$ is driven by the previous-order moment at its own node:
$$m_k^{(i)} = -\sum_j R_{ij} C_j\, m_{k-1}^{(j)}.$$
So $m_2^{(i)} = \sum_j R_{ij} C_j \left(\sum_l R_{jl} C_l\right)$: the same Elmore-style sum, weighted by each capacitor's own first moment. The recursion is two tree traversals per order, linear in the size of the net.

Example, continuing: $-m_1^{(1)} = 50$ ps and $-m_1^{(2)} = 110$ ps. At node 2, $m_2 = R_{21}C_1 \cdot 50 + R_{22}C_2 \cdot 110 = 1 \cdot 20 \cdot 50 + 3 \cdot 30 \cdot 110 = 1000 + 9900 = 10900$ square ps. AWE computes $2q$ moments this way and fits a $q$-pole model to them.

## Spread of the impulse response

Since $\int t\,h = -m_1$ and $\int t^2 h = 2m_2$, the variance of the impulse response, viewed as a distribution, is
$$\sigma^2 = 2m_2 - m_1^2.$$
For a single pole, $H = 1/(1 + s\tau)$ gives $m_1 = -\tau$, $m_2 = \tau^2$ and $\sigma = \tau$.

Example, continuing: $\sigma^2 = 2 \cdot 10900 - 110^2 = 21800 - 12100 = 9700$, so $\sigma = 98.5$ ps against a mean of 110 ps. A spread close to the mean, as here and as at near-end nodes, means a long exponential-like tail: the 50% point sits well before the mean and the Elmore delay is pessimistic. A spread well below the mean, typical of far-end nodes, means a more symmetric response, where Elmore is accurate. Delay metrics such as D2M combine $m_1$ and $m_2$ for this reason. Matching many moments with explicit moments is numerically ill-conditioned, which is why PRIMA and other Krylov methods replaced plain AWE for large models.

## Going deeper

Pillage and Rohrer, "Asymptotic waveform evaluation for timing analysis", IEEE Transactions on Computer-Aided Design, 1990; Odabasioglu, Celik and Pileggi, "PRIMA: passive reduced-order interconnect macromodeling algorithm", 1998; Celik, Pileggi and Odabasioglu, "IC Interconnect Analysis".
