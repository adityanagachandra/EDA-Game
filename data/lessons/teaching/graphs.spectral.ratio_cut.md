# Ratio cut and normalized objectives

A hard balance constraint is one way to stop a partitioner from cutting off a tiny block. Another is to fold balance into the objective: divide the cut by a measure of block sizes, so small blocks are penalized smoothly. Ratio cut and its relatives are the objectives spectral methods actually approximate, and they are natural for clustering, where the right block sizes are not known in advance.

## The ratio cut

Wei and Cheng's ratio cut of a bisection is

$$R(A, B) = \frac{\text{cut}(A, B)}{|A| \cdot |B|}.$$

The product $|A||B|$ is largest when the sides are equal, so a balanced partition with the same cut scores better. With cut weight 6 and sides of 4 and 4, $R = 6/16 = 0.375$. Splitting off a single vertex with cut 2 out of 8 vertices gives $R = 2/7 \approx 0.286$, so ratio cut still allows lopsided splits when they are much sparser, which is the point: it finds natural clusters rather than forced halves.

## Comparing candidate partitions

Because ratio cut trades cut against balance, comparing partitions needs the full ratio, not the cut alone. Partition $P$ with cut 12 and sides 3 and 3 has $R = 12/9 \approx 1.333$; partition $Q$ with cut 18 and sides 2 and 4 has $R = 18/8 = 2.25$. $P$ wins on both counts here, but a partition with a smaller cut and very unequal sides can lose to a balanced one with a larger cut.

Minimizing ratio cut is NP-hard, like bisection. Its value for tool design is that it has no tolerance parameter to tune and that it is the quantity spectral methods bound. Clustering steps in placement and in hierarchical floorplanning use ratio-style objectives to decide where natural boundaries lie.

## RatioCut and the Fiedler vector

Spectral clustering uses $\text{RatioCut}(A, B) = \text{cut}/|A| + \text{cut}/|B| = n \cdot \text{cut}/(|A||B|)$, the same objective scaled by $n$. Set $x_i = \sqrt{|B|/|A|}$ on $A$ and $-\sqrt{|A|/|B|}$ on $B$. Then $x \perp \mathbf{1}$, $x^T x = n$ and $x^T L x = n \cdot \text{RatioCut}(A, B)$, so the Rayleigh quotient of $x$ equals the RatioCut.

Minimizing over all real $x \perp \mathbf{1}$ instead of these two-valued vectors gives $\lambda_2$, so $\lambda_2 \le \text{RatioCut}$ for every bisection, a lower bound computable in polynomial time. Hagen and Kahng used this to justify spectral ratio-cut partitioning: compute the Fiedler vector, then pick the best split along its sorted order. Normalized cut replaces $|A|$ by the total degree of $A$ and corresponds to the normalized Laplacian.

## Going deeper

Wei and Cheng, "Ratio cut partitioning for hierarchical designs" (IEEE TCAD, 1991); Hagen and Kahng, "New spectral methods for ratio cut partitioning and clustering" (IEEE TCAD, 1992).
