# Moments of an RC tree

The moments of a network's impulse response are the coefficients of its transfer function's Taylor expansion in $s$. The first moment is the Elmore delay; the second carries information about spread and shape. Reduction methods match moments, delay metrics are built from them, and for RC trees they have a simple path-tracing formula that a tool can evaluate in linear time.

## The first moment

Expand the transfer function from the driver to node $i$ as $H_i(s) = 1 - m_1(i) s + m_2(i) s^2 - \dots$. For an RC tree, $m_1(i) = \sum_k R_{ik} C_k$, where $R_{ik}$ is the resistance shared by the paths from the root to $i$ and to $k$, and the sum runs over every capacitor. This is the Elmore delay, the mean of the impulse response when it is nonnegative, as it is for RC trees. In a tree with 0.2 kohm from the driver to node 1 (10 fF), 0.5 kohm from node 1 to node 2 (3 fF), 0.1 kohm from node 1 to node 3 (5 fF), and 0.3 kohm from node 2 to node 4 (3 fF), node 2 shares 0.2 kohm with everything and 0.7 kohm with itself and node 4: $m_1(2) = 0.2 \times 21 + 0.5 \times 6 = 7.2$ ps.

## The second moment

Higher moments follow the same pattern recursively: $m_q(i) = \sum_k R_{ik} C_k m_{q-1}(k)$, with $m_0 = 1$. Each capacitor's contribution is weighted by its own previous moment, so the second moment needs the first moments of every node first. For a chain with 0.5 kohm to node 1 (5 fF) and 0.4 kohm to node 2 (3 fF), the first moments are $m_1(1) = 0.5 \times 8 = 4$ ps and $m_1(2) = 4 + 0.4 \times 3 = 5.2$ ps. Then $m_2(2) = 0.5 \times 5 \times 4 + 0.9 \times 3 \times 5.2 = 10 + 14.04 = 24.04$ ps$^2$. Computing all moments to order $q$ takes $q$ passes over the tree, each linear in its size, by accumulating downstream weighted capacitance.

## Why lumping fails

Moments make the cost of crude models precise. Lumping a chain into one resistor $R_{tot}$ and one capacitor $C_{tot}$ gives first moment $R_{tot} C_{tot}$, which puts every capacitor behind all the resistance. For a chain of 0.1, 0.5 and 0.4 kohm with 3, 5 and 10 fF, the true $m_1$ at the end is $0.1 \times 18 + 0.5 \times 15 + 0.4 \times 10 = 13.3$ ps, while the lump gives $1 \times 18 = 18$ ps, 35% high. A distributed uniform line has $m_1 = RC/2$, half the lump. Delay metrics such as the D2M metric combine the first two moments to approximate the 50% delay better than Elmore alone, which is an upper bound on the 50% delay of an RC tree's step response.

## Going deeper

Rubinstein, Penfield and Horowitz, "Signal delay in RC tree networks", IEEE Transactions on CAD (1983); Alpert, Devgan and Kashyap, "RC delay metrics for performance optimization", IEEE Transactions on CAD (2001).
