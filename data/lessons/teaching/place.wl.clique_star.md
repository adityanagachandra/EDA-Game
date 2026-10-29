# Clique and star net models

Quadratic placement works only with two-pin connections, so every multi-pin net must be replaced by a graph of weighted two-pin edges. The clique and star models are the two classic substitutions. Their weights decide how a net pulls on its cells, and their edge counts decide whether the linear system stays sparse enough to solve.

## The clique model

The clique model joins every pair of a $k$-pin net's pins with an edge, $k(k-1)/2$ edges in total. Its quadratic wirelength in x is $w\sum_{i<j}(x_i - x_j)^2$. A common weight is $w = 1/(k-1)$, chosen so that large nets are not overweighted relative to small ones: a $k$-pin clique has many more edges than a two-pin net. For a three-pin net at x = 0, 4 and 10 the pairwise squared distances are 16, 100 and 36, summing to 152, and with weight $1/2$ the cost is 76. Other weightings, such as $2/k$, appear in the literature; the choice is a modeling decision, not a law.

## The star model

The star model adds one auxiliary movable node $s$ per net and connects it to each of the $k$ pins, $k$ edges in total. With unit weights its cost is $\sum_i (x_i - x_s)^2$, minimised when $x_s$ is the mean of the pins. For pins at 0, 4 and 10 the mean is 4.667 and the cost is $21.78 + 0.44 + 28.44 = 50.67$. A useful identity links the two models: $\sum_{i<j}(x_i - x_j)^2 = k\sum_i (x_i - \bar{x})^2$, so with the star node at the mean, a clique of weight $w$ equals a star of weight $kw$. In that sense the two are equivalent, and the star is the sparse one.

## Sparsity and choosing a model

A 100-pin net is 4950 clique edges but only 100 star edges plus one extra variable. Because quadratic placement solves a sparse linear system whose nonzeros are these edges, cliques on large nets would make the matrix dense and the solve slow. Practical placers therefore use cliques for small nets and stars for large ones, with a threshold of a few pins. Both models still share a deeper problem: a quadratic cost grows with the square of distance, so it over-penalises long nets and does not match HPWL. The bound-to-bound model fixes that by choosing edges and weights from the current placement so that the quadratic cost equals HPWL.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), section 4.3; Viswanathan and Chu, "FastPlace" (ISPD 2004), on hybrid net models.
