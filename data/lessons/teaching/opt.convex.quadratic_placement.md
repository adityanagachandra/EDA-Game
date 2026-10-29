# Quadratic placement as a convex QP

Quadratic placement models every two-pin connection as a spring and minimizes the total squared length. The objective is a convex quadratic, its minimum is the solution of one sparse symmetric positive definite linear system per coordinate, and that system can be built and solved for millions of cells. It is the starting point of most analytical placers, from early force-directed methods through FastPlace and SimPL, and the cleanest example of an EDA problem solved exactly by linear algebra.

## One movable cell between fixed pins

With one movable cell at x connected to fixed pins at $p_k$ with weights $w_k$, the cost is $\sum_k w_k (x - p_k)^2$. Its derivative $2\sum_k w_k (x - p_k)$ vanishes at the weighted mean
$$x^* = \frac{\sum_k w_k p_k}{\sum_k w_k}.$$
A cell connected to pins at 10 (weight 1), 40 (weight 3) and 90 (weight 1) lands at $(10 + 120 + 90)/5 = 44$ um.

Two consequences shape practice. Squared length pulls hard on long nets, so quadratic placement overweights them compared with HPWL; placers reweight each net by roughly the inverse of its current length to imitate linear wirelength. And with no overlap constraint cells pile up near the center of mass of their pins, which is why spreading forces, density penalties or look-ahead legalization must follow.

## Coupled cells and the linear system

With several movable cells the gradient conditions couple them. Cell 1 is tied to a pad at $P_1$ with weight $w_1$, cell 2 to a pad at $P_2$ with weight $w_2$, and the cells to each other with weight $w_{12}$. Setting both partial derivatives to zero gives
$$(w_1 + w_{12}) x_1 - w_{12} x_2 = w_1 P_1, \qquad -w_{12} x_1 + (w_2 + w_{12}) x_2 = w_2 P_2.$$
With $P_1 = 0$, $P_2 = 100$ and all weights 1: $2x_1 - x_2 = 0$ and $-x_1 + 2x_2 = 100$, so $x_1 = 33.33$ and $x_2 = 66.67$. The cells split the distance evenly, as three equal springs in series must.

The matrix is symmetric with positive diagonal and is positive definite whenever every connected group of cells touches at least one pad. The x and y systems are independent, so a placer solves two of them.

## Building A and b for a netlist

For a netlist of movable cells, $A = L + D$, where L is the weighted Laplacian of the movable-cell graph and D adds each cell's pad weights on the diagonal; b collects pad weight times pad position. So $A_{ii}$ is the total weight of every net at cell i including pads, $A_{ij} = -w_{ij}$ for movable neighbours i and j, and $b_i = \sum_{\text{pads } k} w_{ik} P_k$.

Example: cell 2 connects to cell 1 with weight 3, to cell 3 with weight 2, and to a pad at 40 with weight 1. Then $A_{22} = 3 + 2 + 1 = 6$, $A_{21} = -3$, $A_{23} = -2$ and $b_2 = 40$. Multi-pin nets are first expanded into two-pin edges by a clique or star model, or by the bound-to-bound model of Kraftwerk2 and SimPL, which reproduces HPWL exactly at the current placement. The matrix is sparse, with a few nonzeros per row, and is solved by preconditioned conjugate gradient.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 4 on analytical placement; Viswanathan and Chu, "FastPlace: efficient analytical placement using cell shifting, iterative local refinement and a hybrid net model" (IEEE TCAD, 2005); Kim, Lee and Markov, "SimPL: an effective placement algorithm" (IEEE TCAD, 2012).
